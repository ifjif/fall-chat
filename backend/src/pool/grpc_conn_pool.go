package pool

import (
	"context"
	"fmt"
	"maps"
	"sync"
	"sync/atomic"
	"time"
	"zzc/fall-chat/src/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ConnFactory func(addr string) (*grpc.ClientConn, error)

func WithDefaultFactory() ConnFactory {
	return func(addr string) (*grpc.ClientConn, error) {
		conn, err := grpc.NewClient(addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("dial failed: %w", err)
		}

		return conn, nil
	}
}

type pooledConn struct {
	*grpc.ClientConn
	lastUsed time.Time
}

type ConnPool struct {
	mu          sync.RWMutex
	conns       map[string]*pooledConn // -> {"ip:port":conn}
	closed      atomic.Uint32
	factory     ConnFactory
	idleTimeout time.Duration
}

type ConnPoolConfig struct {
	Factory     ConnFactory
	IdleTimeout time.Duration
}

func NewConnPool(cfg ConnPoolConfig) (*ConnPool, error) {
	if cfg.Factory == nil {
		cfg.Factory = WithDefaultFactory()
	}
	if cfg.IdleTimeout <= 0 {
		cfg.IdleTimeout = 6 * time.Minute
	}

	cp := &ConnPool{
		conns:       make(map[string]*pooledConn),
		factory:     cfg.Factory,
		idleTimeout: cfg.IdleTimeout,
	}

	go cp.cleaner()

	return cp, nil
}

func (c *ConnPool) dial(addr string) (*pooledConn, error) {
	conn, err := c.factory(addr)
	if err != nil {
		return nil, err
	}

	return &pooledConn{
		ClientConn: conn,
		lastUsed:   time.Now(),
	}, nil
}

func (c *ConnPool) createAndStore(addr string) (*grpc.ClientConn, error) {
	c.mu.RLock()
	if conn, ok := c.conns[addr]; ok {
		c.mu.RUnlock()
		return conn.ClientConn, nil
	}
	c.mu.RUnlock()

	newConn, err := c.dial(addr)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if conn, ok := c.conns[addr]; ok {
		newConn.Close()
		return conn.ClientConn, nil
	}

	newConn.lastUsed = time.Now()
	c.conns[addr] = newConn
	return newConn.ClientConn, nil
}

// 调用方不需要 加锁
func (c *ConnPool) removeConnAndClose(addr string, target *pooledConn) {
	if target == nil {
		return
	}
	c.mu.Lock()
	conn := c.conns[addr]
	if conn == nil || conn != target {
		c.mu.Unlock()
		return // 已被别处拿走
	}
	delete(c.conns, addr)
	c.mu.Unlock()
	conn.Close()
}

// 长时间没有使用
func (c *ConnPool) isTimeout(cc *pooledConn) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Since(cc.lastUsed) > c.idleTimeout
}

// 仅判断状态 是否 connectivity.Shutdown
func (c *ConnPool) isValid(cc *pooledConn) bool {
	if cc == nil {
		return false
	}

	return cc.GetState() != connectivity.Shutdown
}

// 发送 health 请求
func isHealthy(ctx context.Context, cc *pooledConn) bool {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	client := grpc_health_v1.NewHealthClient(cc)
	resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return false // 网络/超时/RST
	}

	return resp.Status == grpc_health_v1.HealthCheckResponse_SERVING
}

func (c *ConnPool) cleaner() {
	ticker := time.NewTicker(c.idleTimeout / 2)
	defer ticker.Stop()

	for range ticker.C {
		if c.closed.Load() == 1 {
			return
		}

		c.cleanIdle()
	}
}

func (c *ConnPool) cleanIdle() {
	var remove []string

	idles := make(map[string]*pooledConn)
	c.mu.Lock()
	maps.Copy(idles, c.conns)
	c.mu.Unlock()

	for k, v := range idles {
		if !c.isValid(v) {
			remove = append(remove, k)
			continue
		}

		if c.isTimeout(v) {
			remove = append(remove, k)
			utils.LogWarn("超时进行移出")
			continue
		}

		/*
			*TODO: 暂时不对每个发
								if !isHealthy(context.Background(), v) {
									remove = append(remove, k)
								}
		*/
	}

	for _, s := range remove {
		c.removeConnAndClose(s, idles[s])
	}

	utils.LogDebug("剩余：%v", c.conns)
}

func (c *ConnPool) Get(addr string) (*grpc.ClientConn, error) {
	if c.closed.Load() == 1 {
		return nil, fmt.Errorf("连接池已关闭")
	}
	c.mu.RLock()
	conn, ok := c.conns[addr]
	c.mu.RUnlock()

	if !ok {
		return c.createAndStore(addr)
	}

	if !c.isValid(conn) {
		c.removeConnAndClose(addr, conn)
		return c.createAndStore(addr)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	conn.lastUsed = time.Now()

	return conn.ClientConn, nil
}

func (c *ConnPool) Close() {
	if !c.closed.CompareAndSwap(0, 1) {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for addr, cc := range c.conns {
		cc.Close()
		delete(c.conns, addr)
	}
}

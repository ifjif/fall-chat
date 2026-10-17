package pool

import "log"

var GrpcConnPool *ConnPool

func init() {
	p, err := NewConnPool(ConnPoolConfig{
		Factory: WithDefaultFactory(),
	})
	if err != nil {
		log.Fatalf("grpc conn pool create failed: %v", err)
	}
	GrpcConnPool = p
}

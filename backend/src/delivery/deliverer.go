package delivery

import (
	"context"
	"fmt"
	"zzc/fall-chat/src/conns"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/pool"
	"zzc/fall-chat/src/proto/chatpt"
)

type Deliverer interface {
	Deliver(ctx context.Context, userID string, payload []byte) error
}

type DelivererDeps struct {
	Hub      *conns.Hub
	Store    store.WSStore
	GrpcPool *pool.ConnPool
}

type delivererImpl struct {
	hub      *conns.Hub
	store    store.WSStore
	grpcPool *pool.ConnPool
}

func NewDeliverer(deps DelivererDeps) Deliverer {
	return &delivererImpl{
		hub:      deps.Hub,
		store:    deps.Store,
		grpcPool: deps.GrpcPool,
	}
}

func (d *delivererImpl) Deliver(ctx context.Context, userID string, payload []byte) error {
	// 1.本地命中，直接投递
	if client, ok := d.hub.GetClient(userID); ok {
		client.SendMessage(payload)
		return nil
	}

	// 2.本地未命中，查路由表做跨节点转发
	addr, err := d.store.WSCache().GetGrpcAddr(ctx, userID)
	if err != nil {
		return fmt.Errorf("get grpc addr from cache: %w", err)
	}
	if addr == "" {
		// 3. 路由表也没有，说明对方不在线或节点已下线
		// return d.saveOffline(ctx, userID, payload)
		return nil
	}

	conn, err := d.grpcPool.Get(addr)
	if err != nil {
		return fmt.Errorf("get grpc conn %s: %w", addr, err)
	}

	stub := chatpt.NewChatClient(conn)
	_, err = stub.Chat(ctx, &chatpt.ChatpMessage{
		Content: payload,
		ToID:    userID,
	})
	if err != nil {
		return fmt.Errorf("remote chat rpc %s: %w", addr, err)
	}

	return nil
}

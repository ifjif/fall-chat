package grpcs

import (
	"zzc/fall-chat/src/proto/chatpt"
	"zzc/fall-chat/src/server"
)

func Init(s *server.GrpcServer) {
	// 关键：把实现注册进去，内部会用 xx_ServiceDesc 完成路由
	chatpt.RegisterChatServer(s.App, NewChatServer())
}

/*
func callClient() {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer secret")
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}
	defer conn.Close()
	// 关键：用生成的 NewGreeterClient 包一层连接
	c := hello.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.SayHello(ctx, &hello.HelloRequest{Name: "Alice"})
	if err != nil {
		log.Fatalf("call failed: %v", err)
	}

	utils.LogInfo(resp.Message)

	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				utils.LogError("recv err: %v", err)
				return
			}
			utils.LogDebug("from server: %s", msg.Content)
		}
	}()

	for i := range 5 {
		if err := stream.Send(&hello.ChatMessage{
			Content: fmt.Sprintf("msg %d", i),
		}); err != nil {
			utils.LogError("send err: %v", err)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err := stream.CloseSend(); err != nil {
		utils.LogError("close send err:", err)
	}
}
*/

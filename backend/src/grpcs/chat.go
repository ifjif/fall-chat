package grpcs

import (
	"context"
	"zzc/fall-chat/src/proto/chatpt"
	"zzc/fall-chat/src/service"
	"zzc/fall-chat/src/utils"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatServer struct {
	chatpt.UnimplementedChatServer
	wsService *service.Ws
}

func NewChatServer() *ChatServer {
	s := &ChatServer{}
	s.Inject()
	return s
}

func (s *ChatServer) Inject() {
	wsService := service.GlobalService[service.WS_SERVICE_NAME].(*service.Ws)
	s.wsService = wsService
}

func (s *ChatServer) Chat(ctx context.Context, mes *chatpt.ChatpMessage) (*emptypb.Empty, error) {
	utils.LogInfo("收到grpc消息：to=%s ", mes.ToID)
	s.wsService.SendBytes(ctx, mes.ToID, mes.Content)
	return &emptypb.Empty{}, nil
}

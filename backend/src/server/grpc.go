package server

import (
	"context"
	"log"
	"net"
	"time"
	"zzc/fall-chat/src/host"
	"zzc/fall-chat/src/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	App      *grpc.Server
	listener net.Listener
}

func (s *GrpcServer) Serve() {
	utils.LogInfo("server listening on %s", host.GRPCAllAddr)
	s.App.Serve(s.listener)
}

var GRPC *GrpcServer

func InitGrpc() {
	lis, err := net.Listen("tcp", host.GRPCAllAddr)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logginInterceptor,
			authInterceptor,
		),
	)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	GRPC = &GrpcServer{
		App:      s,
		listener: lis,
	}
}

func logginInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handle grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	resp, err := handle(ctx, req)
	dur := time.Since(start)

	if err != nil {
		utils.LogError("[grpc] method=%s dur=%v err=%v", info.FullMethod, dur, err)
	} else {
		utils.LogInfo("[grpc] method=%s dur=%v", info.FullMethod, dur)
	}

	return resp, err
}

func authInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	tokens := md.Get("authorization")
	utils.LogInfo("[grpc] 获取token=%v", tokens)

	return handler(ctx, req)
}

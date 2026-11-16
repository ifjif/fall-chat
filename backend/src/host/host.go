package host

import (
	"log"
	"zzc/fall-chat/src/utils"
)

var (
	GRPCPort    = ""
	HOSTIP4     = ""
	GRPCAllAddr = ":" + GRPCPort
	GRPCAddr    = HOSTIP4 + ":" + GRPCPort
)

func init() {
	grpcPort, exists := utils.GetEnv("GRPC_PORT")
	if !exists {
		utils.LogError("未提供grpc port")
		log.Fatal()
	}
	hostIP4, exists := utils.GetEnv("HOST_IP4")
	if !exists {
		utils.LogError("未提供host ip4")
		log.Fatal()

	}
	GRPCPort = grpcPort
	HOSTIP4 = hostIP4
}

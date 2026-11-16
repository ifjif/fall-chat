package mq

import (
	"log"
	"zzc/fall-chat/src/utils"
)

func init() {
	v, exsits := utils.GetEnv("KAFKA_BROKER")
	if !exsits {
		utils.LogError("未提供kafka broker address,kafka不可用")
		log.Fatal()
	}
	KafkaBroker = v
	kafkaInit()
}

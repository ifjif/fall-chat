package binevent

import (
	"encoding/binary"
	"fmt"
)

// 在外部维护一个表映射缓存 (实际生产中需考虑并发安全，如使用 sync.Map)
var (
	TableMapCache        = make(map[uint64]*TableMapInfo)
	BinlogVersion uint16 = 0 // 默认值，从 FORMAT_DESCRIPTION_EVENT 中更新
)

// Binlog Event Header 固定19字节
type EventHeader struct {
	Timestamp uint32 // 4B
	EventType uint8  // 1B
	ServerID  uint32 // 4B
	EventLen  uint32 // 4B
	NextPos   uint32 // 4B
	Flags     uint16 // 2B
}

// 用于缓存的表映射信息
type TableMapInfo struct {
	SchemaName string
	TableName  string
	ColTypes   []byte   // 列类型数组
	ColMetas   []uint16 // 列元数据数组
}

func ParseHeader(data []byte) (*EventHeader, error) {
	if len(data) < 19 {
		return nil, fmt.Errorf("event data too short for header")
	}
	header := &EventHeader{
		Timestamp: binary.LittleEndian.Uint32(data[0:4]),
		EventType: data[4],
		ServerID:  binary.LittleEndian.Uint32(data[5:9]),
		EventLen:  binary.LittleEndian.Uint32(data[9:13]),
		NextPos:   binary.LittleEndian.Uint32(data[13:17]),
		Flags:     binary.LittleEndian.Uint16(data[17:19]),
	}
	return header, nil
}

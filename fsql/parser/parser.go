package parser

import (
	"encoding/binary"
	"fmt"
	"time"
)

// 从字节切片中解析 Packed Integer
// 返回值：解析出的整数值，消耗的字节数，错误信息
func ReadPackedInteger(data []byte) (uint64, int, error) {
	if len(data) == 0 {
		return 0, 0, fmt.Errorf("data too short for packed integer")
	}

	switch data[0] {
	case 251: // 0xfb,表示NULl(在某些上下文中)
		return 0, 1, nil
	case 252: // 0xfc: 后跟2字节
		if len(data) < 3 {
			return 0, 0, fmt.Errorf("data too short for 2-byte packed integer")
		}
		return uint64(binary.LittleEndian.Uint16(data[1:3])), 3, nil
	case 253: // 0xfd: 后跟3字节
		if len(data) < 4 {
			return 0, 0, fmt.Errorf("data too short for 3-byte packed integer")
		}
		val := uint64(data[0]) | uint64(data[1])<<8 | uint64(data[2])<<16
		return val, 4, nil
	case 254: // 0xfe:后更8字节
		if len(data) < 9 {
			return 0, 0, fmt.Errorf("data too short for 8-byte packed integer")
		}
		return binary.LittleEndian.Uint64(data[1:9]), 9, nil

	default: // 0x00 - 0xfa: 1字节
		return uint64(data[0]), 1, nil
	}
}

// isNullBitmapSet 检查在 NULL Bitmap 中，指定列索引是否为 NULL
// bitmap: NULL Bitmap 字节切片
// colIndex: 列的索引 (从 0 开始)
// offset: 偏移量 (在 Row Event 中，通常 offset = 2)
func IsNullBitmapSet(bitmap []byte, colIndex int, offset int) bool {
	// 计算对应的字节位置 (从 0 开始)
	bytePos := (colIndex + offset) / 8
	// 计算该字节中的位位置（从0开始)
	bitPos := (colIndex + offset) % 8

	if bytePos >= len(bitmap) {
		return false // 越界保护
	}

	// 检查该位置是否为1（表示NULL）
	return (bitmap[bytePos] & (1 << bitPos)) != 0
}

// DecodeDatetime2 解析 MySQL DATETIME2 格式的字节数组
// data: 从 binlog 读取的字节切片 (长度 4~7 字节)
// fsp: 小数秒精度 (0~6)
func DecodeDatetime2(data []byte, fsp int) (time.Time, error) {
	if len(data) < 4 {
		return time.Time{}, fmt.Errorf("datetime2 data too short")
	}

	// 1. 读取前 4 字节，转换为 uint32
	intPart := uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])

	// 2. 按照位掩码提取 年、月、日、时、分、秒
	year := int(intPart >> 23 & 0x1FFFF) // 17 bits
	month := int(intPart >> 19 & 0x0F)   // 4 bits
	day := int(intPart >> 14 & 0x1F)     // 5 bits
	hour := int(intPart >> 9 & 0x1F)     // 5 bits
	minute := int(intPart >> 3 & 0x3F)   // 6 bits
	second := int(intPart & 0x3F)        // 6 bits

	// 3. 解析小数秒 (微秒)
	// 根据 fsp 决定读取多少个字节来补充微秒
	var microsecond int64
	if fsp > 0 {
		fracBytes := (fsp + 1) / 2
		if len(data) < 4+fracBytes {
			return time.Time{}, fmt.Errorf("datetime2 fractional part too short")
		}
		var frac int64
		for i := 0; i < fracBytes; i++ {
			frac = frac<<8 | int64(data[4+i])
		}
		// 将读取到的字节转换为微秒
		// 如果读取的字节数不足3字节，需要补齐
		for i := 0; i < 3-fracBytes; i++ {
			frac *= 256
		}
		microsecond = frac
	}

	return time.Date(year, time.Month(month), day, hour, minute, second, int(microsecond*1000), time.UTC), nil
}

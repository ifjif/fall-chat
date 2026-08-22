package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// readUntil 辅助函数：从 bytes.Reader 中逐字节读取，直到遇到指定的分隔符
func readUntil(reader *bytes.Reader, delim byte) ([]byte, error) {
	var buf bytes.Buffer
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return buf.Bytes(), err
		}
		if b == delim {
			break
		}
		buf.WriteByte(b)
	}
	return buf.Bytes(), nil
}

// HandshakeV10 结构体，用于存储解析出的关键字段
type HandshakeV10 struct {
	ProtocolVersion uint8
	ServerVersion   string
	ConnectionID    uint32
	AuthPluginData  []byte // 完整的盐值 (Scramble)
	AuthPluginName  string
	Capabilities    uint32
}

// ParseHandshakeV10 解析 MySQL HandshakeV10 握手包
func ParseHandshakeV10(data []byte) (*HandshakeV10, error) {
	reader := bytes.NewReader(data)
	hs := &HandshakeV10{}

	// 1. 读取协议版本 (1 byte)
	if err := binary.Read(reader, binary.LittleEndian, &hs.ProtocolVersion); err != nil {
		return nil, fmt.Errorf("读取协议版本失败: %w", err)
	}

	// 2. 读取服务器版本 (NUL-terminated string)
	serverVersion, err := readUntil(reader, 0x00)
	if err != nil {
		return nil, fmt.Errorf("读取服务器版本失败: %w", err)
	}
	hs.ServerVersion = string(serverVersion)

	// 3. 读取连接ID (4 bytes)
	if err := binary.Read(reader, binary.LittleEndian, &hs.ConnectionID); err != nil {
		return nil, fmt.Errorf("读取连接ID失败: %w", err)
	}

	// 4. 提取第一部分盐值 (auth-plugin-data-part-1, 8 bytes)
	part1 := make([]byte, 8)
	if _, err := io.ReadFull(reader, part1); err != nil {
		return nil, fmt.Errorf("读取盐值第一部分失败: %w", err)
	}

	// 5. 跳过填充字节 (filler, 1 byte, 固定为 0x00)
	if _, err := reader.ReadByte(); err != nil {
		return nil, fmt.Errorf("读取填充字节失败: %w", err)
	}

	// 6. 读取能力标志位低2字节 (capability_flags_1, 2 bytes)
	var capFlags1 uint16
	if err := binary.Read(reader, binary.LittleEndian, &capFlags1); err != nil {
		return nil, fmt.Errorf("读取能力标志位低2字节失败: %w", err)
	}

	// 7. 跳过默认字符集 (1 byte)
	if _, err := reader.ReadByte(); err != nil {
		return nil, fmt.Errorf("读取字符集失败: %w", err)
	}

	// 8. 跳过服务器状态标志 (2 bytes)
	if _, err := io.CopyN(io.Discard, reader, 2); err != nil {
		return nil, fmt.Errorf("读取服务器状态标志失败: %w", err)
	}

	// 9. 读取能力标志位高2字节 (capability_flags_2, 2 bytes)
	var capFlags2 uint16
	if err := binary.Read(reader, binary.LittleEndian, &capFlags2); err != nil {
		return nil, fmt.Errorf("读取能力标志位高2字节失败: %w", err)
	}
	// 合并完整的能力标志位
	capFlags := uint32(capFlags1) | (uint32(capFlags2) << 16)

	// 10. 检查是否支持 CLIENT_PLUGIN_AUTH (标志位 0x00200000)
	var authPluginDataLen uint8 = 0
	if capFlags&0x00200000 != 0 {
		// 读取 auth-plugin-data-len (1 byte)
		if err := binary.Read(reader, binary.LittleEndian, &authPluginDataLen); err != nil {
			return nil, fmt.Errorf("读取盐值总长度失败: %w", err)
		}
	} else {
		// 如果不支持，读取固定的 0x00
		if _, err := reader.ReadByte(); err != nil {
			return nil, fmt.Errorf("读取固定0字节失败: %w", err)
		}
	}

	// 11. 跳过保留字段 (reserved, 10 bytes, 全0)
	if _, err := io.CopyN(io.Discard, reader, 10); err != nil {
		return nil, fmt.Errorf("读取保留字段失败: %w", err)
	}

	// 12. 提取第二部分盐值 (auth-plugin-data-part-2)
	// 长度计算规则: MAX(13, auth_plugin_data_len - 8)
	var part2 []byte
	if capFlags&0x00200000 != 0 && authPluginDataLen > 0 {
		part2Len := uint32(authPluginDataLen) - 8
		if part2Len < 13 {
			part2Len = 13
		}
		part2 = make([]byte, part2Len)
		if _, err := io.ReadFull(reader, part2); err != nil {
			return nil, fmt.Errorf("读取盐值第二部分失败: %w", err)
		}
		// 协议规范 part-2包含一个末尾的\0, 需截断
		if len(part2) > 0 && part2[len(part2)-1] == 0x00 {
			part2 = part2[:len(part2)-1]
		}
	}
	// 增加能力标志
	hs.Capabilities = capFlags

	// 13. 拼接完整的盐值 (Scramble)
	hs.AuthPluginData = append(part1, part2...)

	// 14. 读取认证插件名称 (NUL-terminated string)
	if capFlags&0x00200000 != 0 {
		pluginName, err := readUntil(reader, 0x00)
		if err == nil {
			hs.AuthPluginName = string(pluginName)
		}
	}

	return hs, nil
}

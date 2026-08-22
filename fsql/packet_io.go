package main

import (
	"fmt"
	"io"
	"net"
)

// 封装 mysql协议底层网络读写
type PacketIO struct {
	conn     net.Conn
	sequence uint8 // 序列号，每次发包/收包 自增
}

func NewPacketIO(conn net.Conn) *PacketIO {
	return &PacketIO{conn: conn}
}

// 读取一个完整的mysql 网络包
func (p *PacketIO) ReadPacket() ([]byte, error) {
	// 读4字节包头
	header := make([]byte, 4)
	if _, err := io.ReadFull(p.conn, header); err != nil {
		return nil, err
	}

	// 解析payload 长度和序列号
	length := int(uint32(header[0]) | uint32(header[1]<<8) | uint32(header[2]<<16))
	seq := header[3]

	// 简单校验序列号
	if seq != p.sequence {
		return nil, fmt.Errorf("MySQL 序列号不匹配: 期望 %d, 实际收到 %d", p.sequence, seq)
	}

	p.sequence = seq + 1

	payload := make([]byte, length)
	if _, err := io.ReadFull(p.conn, payload); err != nil {
		return nil, err
	}

	return payload, nil
}

// ResetSequence 重置序列号为 0（在发送新的 Command 前必须调用）
func (p *PacketIO) ResetSequence() {
	p.sequence = 0
}

// 发送一个完整的mysql包
func (p *PacketIO) WritePacket(data []byte) error {
	// 构造 4b 包头
	length := len(data)
	header := []byte{
		byte(length),
		byte(length >> 8),
		byte(length >> 16),
		p.sequence,
	}
	p.sequence++
	// ✅ 正确做法：先写包头，再写包体，绝不使用 append！
	if _, err := p.conn.Write(header); err != nil {
		return err
	}
	if _, err := p.conn.Write(data); err != nil {
		return err
	}

	return nil
}

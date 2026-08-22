package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
)

// caching_sha2_password 协议常量
const (
	AuthStatusOK               byte = 0x03
	AuthStatusFastAuthFail     byte = 0x04
	AuthSwitchRequestPacket    byte = 0xfe
	AuthStatusRequestPublicKey byte = 0x02
	AuthStatusPublicKey        byte = 0x01
)

// CachingSha2Authenticator 封装认证状态机
type CachingSha2Authenticator struct {
	Pio        *PacketIO
	Password   string
	Salt       []byte
	ServerCaps uint32
}

// Authenticate 执行完整的 caching_sha2_password 认证流程
func (c *CachingSha2Authenticator) Authenticate(user string) error {
	// 6. Auth Response Data
	authResp := CalcCachingSha2Password(c.Salt, c.Password)
	resp := BuildHandshakeResponse41(authResp, user, c.ServerCaps)

	// 2. 发送初始认证包
	if err := c.Pio.WritePacket(resp); err != nil {
		return fmt.Errorf("发送初始认证包失败: %w", err)
	}

	fmt.Println("开始读取响应")
	// 3. 读取服务端响应
	resp, err := c.Pio.ReadPacket()
	fmt.Println("开始读取响应完成")
	if err != nil {
		return fmt.Errorf("读取认证响应失败: %w", err)
	}

	if len(resp) == 0 {
		fmt.Println("空的认证响应")
		return fmt.Errorf("收到空的认证响应")
	}

	fmt.Printf("收到服务端响应: %s\n", string(resp))

	// 4. 状态机处理
	switch resp[0] {
	case AuthStatusOK:
		fmt.Println("认证成功")
		// 服务端缓存命中，直接认证成功
		return nil

	case AuthStatusFastAuthFail:
		fmt.Println("缓存未命中")
		// 缓存未命中，需要执行 Full Authentication
		return c.performFullAuthentication()

	case AuthSwitchRequestPacket:
		fmt.Println("认证切换请求")
		return c.handleAuthSwitch(resp, user)

	default:
		return fmt.Errorf("未知的认证状态码: 0x%02x", resp[0])
	}
}

func (c *CachingSha2Authenticator) handleAuthSwitch(packet []byte, user string) error {
	// 1. 解析新插件名 (跳过首字节 0xfe，读取 NUL 结尾字符串)
	reader := bytes.NewReader(packet[1:])
	pluginName, err := readUntil(reader, 0x00)
	if err != nil {
		return fmt.Errorf("解析 AuthSwitch 插件名失败: %w", err)
	}

	// 2. 读取剩余的所有字节作为新的 Salt
	newSalt, err := readUntil(reader, 0x00)
	if err != nil {
		return fmt.Errorf("读取 AuthSwitch 新盐值失败: %w", err)
	}

	fmt.Printf("切换到插件: %s, 新盐值: %x\n", pluginName, newSalt)

	// 3. 根据新插件计算 auth_response
	var newAuthResponse []byte
	switch string(pluginName) {
	case "mysql_native_password":
		newAuthResponse = CalcMysqlNativePassword(newSalt, c.Password)
	default:
		return fmt.Errorf("不支持的 AuthSwitch 插件: %s", pluginName)
	}

	// 4. 发送 AuthSwitchResponse (直接发送 payload，不需要包头，Pio 会自动处理)
	if err := c.Pio.WritePacket(newAuthResponse); err != nil {
		return fmt.Errorf("发送 AuthSwitchResponse 失败: %w", err)
	}

	// 5. 读取最终的 OK 包
	finalResp, err := c.Pio.ReadPacket()
	if err != nil {
		return fmt.Errorf("读取 AuthSwitch 最终响应失败: %w", err)
	}

	if finalResp[0] == 0x00 {
		fmt.Println("认证成功 (AuthSwitch)")
		return nil
	}

	return fmt.Errorf("AuthSwitch 认证失败: %s", string(finalResp))
}

// performFullAuthentication 处理缓存未命中的情况
func (c *CachingSha2Authenticator) performFullAuthentication() error {
	// 5. 发送请求公钥指令 (0x02)
	if err := c.Pio.WritePacket([]byte{AuthStatusRequestPublicKey}); err != nil {
		return fmt.Errorf("请求RSA公钥失败: %w", err)
	}

	// 6. 接收服务端下发的 RSA 公钥包 (首字节为 0x01)
	pubKeyPacket, err := c.Pio.ReadPacket()
	if err != nil {
		return fmt.Errorf("接收RSA公钥失败: %w", err)
	}
	if len(pubKeyPacket) < 1 || pubKeyPacket[0] != AuthStatusPublicKey {
		return fmt.Errorf("期望收到公钥包(0x01)，实际收到: 0x%02x", pubKeyPacket[0])
	}

	// 7. 解析 PEM 格式的 RSA 公钥
	pubKey, err := ParseRSAPublicKey(pubKeyPacket[1:])
	if err != nil {
		return fmt.Errorf("解析RSA公钥失败: %w", err)
	}

	// 8. 加密密码并发送
	encryptedPwd, err := EncryptPassword(c.Password, c.Salt, pubKey)
	if err != nil {
		return fmt.Errorf("RSA加密密码失败: %w", err)
	}

	if err := c.Pio.WritePacket(encryptedPwd); err != nil {
		return fmt.Errorf("发送加密密码失败: %w", err)
	}

	// 9. 读取最终的认证结果
	finalResp, err := c.Pio.ReadPacket()
	if err != nil {
		return fmt.Errorf("读取最终认证结果失败: %w", err)
	}

	if len(finalResp) > 0 && finalResp[0] == AuthStatusOK {
		return nil
	}
	return fmt.Errorf("Full Authentication 失败，服务端响应: %v", finalResp)
}

/*
* 向主库发送一个新版的 Dump 命令：COM_BINLOG_DUMP_GTID (命令码 0x1e)。
* 这个命令在协议层面明确告诉主库：“我是一个支持新特性（包括 CRC32 Checksum）的现代从库”。
*
* 主库不支持 GTID模式，不要使用
 */
func (c *CachingSha2Authenticator) SendBinlogDumpGTID(filename string, pos uint64, serverID uint32) error {
	buf := &bytes.Buffer{}

	// 1. Command: COM_BINLOG_DUMP_GTID (0x1e)
	buf.WriteByte(0x1e)

	// 2. Flags (2字节, Little Endian)
	// 0x04: BINLOG_THROUGH_GTID
	var flagsBuf [2]byte
	binary.LittleEndian.PutUint16(flagsBuf[:], 0x04)
	buf.Write(flagsBuf[:])

	// 3. Server ID (4字节, Little Endian)
	var idBuf [4]byte
	binary.LittleEndian.PutUint32(idBuf[:], serverID)
	buf.Write(idBuf[:])

	// 4. Binlog Filename (Length-Encoded String)
	// ⚠️ 注意：这里绝对不能加 \0 结尾！
	writeLenencString(buf, filename)

	// 5. Binlog Position (8字节, Little Endian)
	// ⚠️ 注意：GTID 模式下 Position 必须是 8 字节 (uint64)！
	var posBuf [8]byte
	binary.LittleEndian.PutUint64(posBuf[:], pos)
	buf.Write(posBuf[:])

	// 6. GTID Data (Length-Encoded String)
	// 暂时不使用 GTID，传空字符串（writeLenencString 会写入 1 字节的 0x00）
	writeLenencString(buf, "")

	return c.Pio.WritePacket(buf.Bytes())
}

// 辅助函数：写入 MySQL 协议的 Length-Encoded String
func writeLenencString(buf *bytes.Buffer, s string) {
	length := len(s)
	switch {
	case length <= 250:
		buf.WriteByte(byte(length))
	case length <= 0xffff:
		buf.WriteByte(0xfc)
		buf.WriteByte(byte(length))
		buf.WriteByte(byte(length >> 8))
	case length <= 0xffffff:
		buf.WriteByte(0xfd)
		buf.WriteByte(byte(length))
		buf.WriteByte(byte(length >> 8))
		buf.WriteByte(byte(length >> 16))
	default:
		buf.WriteByte(0xfe)
		// 8字节长度，这里省略，binlog文件名不可能这么大
	}
	buf.WriteString(s)
}

// 发送binlogdump
func (c *CachingSha2Authenticator) SendBinlogDump(filename string, pos uint32, serverID uint32) error {
	buf := &bytes.Buffer{}

	// 1. Command: COM_BINLOG_DUMP (1字节)
	buf.WriteByte(0x12)

	// 2. Position (4字节, Little Endian)
	// 使用 PutUint32 确保绝对只写入 4 个字节
	// binary.Write(buf, binary.LittleEndian, pos) // 如果 pos 确定是 uint32 可以保留，但推荐下面这种：
	// 更安全的写法：
	var posBuf [4]byte
	binary.LittleEndian.PutUint32(posBuf[:], pos)
	buf.Write(posBuf[:])

	// 3. Flags (2字节, Little Endian)
	var flagsBuf [2]byte
	binary.LittleEndian.PutUint16(flagsBuf[:], 0x00)
	buf.Write(flagsBuf[:])

	// 4. Server ID (4字节, Little Endian)
	var idBuf [4]byte
	binary.LittleEndian.PutUint32(idBuf[:], serverID)
	buf.Write(idBuf[:])

	// 5. Filename (以 \0 结尾)
	buf.WriteString(filename)
	// buf.WriteByte(0x00)

	return c.Pio.WritePacket(buf.Bytes())
}

// ==================== 底层密码学工具函数 ====================

// CalcCachingSha2Password 计算 caching_sha2_password 的初始 Auth Response
func CalcCachingSha2Password(salt []byte, password string) []byte {
	// SHA256(password)
	hash1 := sha256Sum([]byte(password))
	// SHA256(SHA256(password))
	hash2 := sha256Sum(hash1)
	// SHA256(SHA256(SHA256(password)) + salt)
	hash3 := sha256Sum(append(hash2, salt...))
	// XOR(SHA256(password), SHA256(SHA256(SHA256(password)) + salt))
	return xorBytes(hash1, hash3)
}

// EncryptPassword 使用 RSA 公钥加密密码 (XOR + OAEP)
func EncryptPassword(password string, salt []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	// 1. 将密码转为字节，末尾补 \0
	pwdBytes := append([]byte(password), 0x00)

	// 2. 与盐值进行循环 XOR
	xorPwd := make([]byte, len(pwdBytes))
	for i := 0; i < len(pwdBytes); i++ {
		xorPwd[i] = pwdBytes[i] ^ salt[i%len(salt)]
	}

	// 3. 使用 RSA-OAEP 加密 (SHA-1 作为哈希函数，这是 MySQL 协议规范)
	ciphertext, err := rsa.EncryptOAEP(nil, rand.Reader, pubKey, xorPwd, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// ParseRSAPublicKey 解析 MySQL 返回的 PEM 格式 RSA 公钥
func ParseRSAPublicKey(pemData []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("PEM解码失败")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("不是RSA公钥")
	}
	return rsaPub, nil
}

// sha256Sum 计算 SHA-256 哈希
func sha256Sum(data []byte) []byte {
	// 实际项目中请使用 crypto/sha256
	// 这里为了代码简洁省略了 import，实际使用时请补全：
	// h := sha256.New()
	// h.Write(data)
	// return h.Sum(nil)
	return nil // 占位，请替换为真实的 sha256 实现
}

// xorBytes 对两个字节切片进行异或操作
func xorBytes(a, b []byte) []byte {
	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i%len(b)]
	}
	return result
}

// writeLenencInt 按照 MySQL 协议规范写入长度编码整数
func WriteLenencInt(buf *bytes.Buffer, n uint64) {
	switch {
	case n <= 250:
		buf.WriteByte(byte(n))
	case n <= 0xffff:
		buf.WriteByte(0xfc)
		binary.Write(buf, binary.LittleEndian, uint16(n))
	case n <= 0xffffff:
		buf.WriteByte(0xfd)
		// 注意：这里需要按小端序写入 3 个字节
		buf.WriteByte(byte(n))
		buf.WriteByte(byte(n >> 8))
		buf.WriteByte(byte(n >> 16))
	default:
		buf.WriteByte(0xfe)
		binary.Write(buf, binary.LittleEndian, n)
	}
}

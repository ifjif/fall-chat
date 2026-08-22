package main

import (
	"bytes"
	"encoding/binary"
)

func BuildHandshakeResponse41(authResponse []byte, username string, serverCaps uint32) []byte {
	buf := &bytes.Buffer{}
	// 定义客户端支持的所有能力标志
	clientSupportedCaps := uint32(
		0x00000200 | // CLIENT_PROTOCOL_41
			0x00008000 | // CLIENT_SECURE_CONNECTION
			0x00080000 | // CLIENT_PLUGIN_AUTH
			0x00200000 | // CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA
			// 0x00000800 | // CLIENT_SSL (极其重要！)
			0x00010000 | // CLIENT_MULTI_STATEMENTS
			0x00020000 | // CLIENT_MULTI_RESULTS
			0x00040000) // CLIENT_PS_MULTI_RESULTS

	// 1. client_flag (int<4>) - 必须包含 CLIENT_PROTOCOL_41 等基础能力
	// clientFlags := uint32(0x0000a20d) // 示例标志，实际需根据服务端能力动态计算
	// 动态计算客户端能力标志
	// clientCaps := serverCaps & (0x00000200 | // CLIENT_PROTOCOL_41
	//	0x00008000 | // CLIENT_SECURE_CONNECTION
	//	0x00080000 | // CLIENT_PLUGIN_AUTH
	//	0x00200000) // CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA
	// clientCaps := serverCaps & clientSupportedCaps
	clientCaps := clientSupportedCaps
	binary.Write(buf, binary.LittleEndian, clientCaps)

	// 2. max_packet_size (int<4>)
	binary.Write(buf, binary.LittleEndian, uint32(4194304))

	// 3. character_set (int<1>) - 例如 utf8mb4 (45)
	buf.WriteByte(45)

	// 4. filler (string[23]) - 23个0(包含了前面的character_set)
	buf.Write(make([]byte, 23))

	// 5. username (string<NUL>)
	buf.WriteString(username)
	buf.WriteByte(0x00)

	// 6. auth_response (string<length>) - 假设启用了 CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA
	// buf.WriteByte(byte(len(authResponse))) // 长度编码 (简化版，实际需用 lenenc 格式)
	WriteLenencInt(buf, uint64(len(authResponse)))
	buf.Write(authResponse)

	// 7. client_plugin_name (string<NUL>) - 假设启用了 CLIENT_PLUGIN_AUTH
	buf.WriteString("caching_sha2_password")
	buf.WriteByte(0)

	return buf.Bytes()
}

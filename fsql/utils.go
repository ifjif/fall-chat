package main

import (
	"crypto/sha1"
	"fmt"
)

// 解析 MySql 握手包
func ParseHandshake(data []byte) (salt []byte, authPlugin string, err error) {
	fmt.Println(len(data))
	if data[0] != 10 {
		return nil, "", fmt.Errorf("unsupported protocol version: %d", data[0])
	}

	// 1. 提取第一部分盐值 (固定 8 bytes, offset 33)
	salt = append(salt, data[8:16]...)

	// 2. 提取第二部分盐值 (固定 12 bytes, offset 52)
	if len(data) >= 64 {
		salt = append(salt, data[25:37]...)
	}

	// 3. 提取 Auth Plugin Name
	// offset 64 是第二部分盐值后面的 \0 终止符
	// offset 65 才是插件名的第一个字符
	if len(data) > 65 {
		pluginStart := 37
		pluginEnd := pluginStart

		// 寻找插件名末尾的 \0 终止符
		for pluginEnd < len(data) && data[pluginEnd] != 0 {
			pluginEnd++
		}

		// 确保截取到正确的字符串
		authPlugin = string(data[pluginStart:pluginEnd])
	}

	return salt, authPlugin, nil
}

// 计算 mysql_native_password 的加密结果
// 公式 xor(sha1(pass), sha1(salt + sha1(sha1(pass))))
func CalcMysqlNativePassword(salt []byte, password string) []byte {
	if len(password) == 0 {
		return nil
	}

	// sha1(pass)
	crypt := sha1.New()
	crypt.Write([]byte(password))

	stage1 := crypt.Sum(nil)

	// sha1(sha1(pass))
	crypt.Reset()
	crypt.Write(stage1)
	stage2 := crypt.Sum(nil)

	// sha1(salt + sha1(sha1(pass)))
	crypt.Reset()
	crypt.Write(salt)
	crypt.Write(stage2)
	stage3 := crypt.Sum(nil)

	// xor(sha1(pass) , sha1(salt + sha1(sha1(pass))))
	result := make([]byte, len(stage1))
	for i := range stage1 {
		result[i] = stage1[i] ^ stage3[i]
	}

	return result
}

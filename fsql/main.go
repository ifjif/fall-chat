package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
	binevent "zzc/fsql/bin_event"
	"zzc/fsql/parser"
)

func main() {
	conn, err := net.Dial("tcp", "mysql8:3306")
	if err != nil {
		log.Fatalln(err)
	}

	pio := NewPacketIO(conn)

	handshake, err := pio.ReadPacket()
	if err != nil {
		log.Fatalln(err)
	}

	hs, err := ParseHandshakeV10(handshake)
	if err != nil {
		log.Fatalf("解析握手包失败: %v", err)
	}
	fmt.Printf("协议版本: %d\n", hs.ProtocolVersion)
	fmt.Printf("服务器版本: %s\n", hs.ServerVersion)
	fmt.Printf("连接ID: %d\n", hs.ConnectionID)
	fmt.Printf("认证插件: %s\n", hs.AuthPluginName)
	fmt.Printf("提取到的完整盐值(Scramble): %x\n", hs.AuthPluginData)

	auth := &CachingSha2Authenticator{
		Pio:        pio,
		Password:   "12345678",
		Salt:       hs.AuthPluginData,
		ServerCaps: hs.Capabilities,
	}

	if hs.AuthPluginName == "mysql_native_password" {
		// todo
	} else if hs.AuthPluginName == "caching_sha2_password" {
		fmt.Println("进行caching_sha2_password认证")
		err = auth.Authenticate("fsql")
		if err != nil {
			log.Fatalf("认证失败：%v", err)
		}
	} else {
		log.Fatalf("不支持的认证插件: %s", hs.AuthPluginName)
	}

	fmt.Println("认证成功")

	// 重置序列号
	auth.Pio.ResetSequence()

	fmt.Println("发送dump命令，进入事件监听循环")
	// 使用 show master status 来查看filename和pos
	// SELECT @@global.gtid_mode; 判断主库是否支持 GTID复制模式
	// err = auth.SendBinlogDumpGTID("mysql-bin.000089", uint64(7411), uint32(100))
	auth.SendBinlogDump("mysql-bin.000132", uint32(919), uint32(100))
	if err != nil {
		log.Fatal("发送 Binlog Dump 失败:", err)
	}

	fmt.Println("开始监听 Binlog 事件流...")

	// 持续读取事件
	for {
		packet, err := auth.Pio.ReadPacket()
		if err != nil {
			log.Fatal("读取 Binlog 包失败:", err)
		}

		// 检查第一个字节：0x00 表示 OK，0xFF 表示服务端报错
		if packet[0] == 0xFF {
			// 解析 Error Packet 并打印错误信息
			log.Fatal("服务端返回错误:", string(packet[1:]))
		}

		// 跳过首字节的 OK 标志，剩下的就是完整的 Binlog Event
		eventData := packet[1:]
		// 【关键】如果主库开启了 CRC32，每个 Event 的末尾都有 4 字节的校验和
		// 在解析 Event 之前，必须把这 4 个字节切掉，否则 Header 里的 Length 会和实际对不上
		// 在[mysqld] binlog_checksum = NONE(关闭checksum)
		// if len(eventData) >= 4 {
		// 取出末尾的 4 字节作为 checksum (你可以选择校验它，也可以直接丢弃)
		// checksum := eventData[len(eventData)-4:]

		// 截断数据，去掉末尾的 4 字节
		//eventData = eventData[:len(eventData)-4]
		//}

		// TODO: 在这里调用你的 Binlog Event 解析器
		// parseBinlogEvent(eventData)

		fmt.Printf("收到一个 Binlog 事件，长度: %d, 数据：%v\n", len(eventData), eventData)

		// 至少需要 19 字节才能解析 Event Header
		if len(eventData) < 19 {
			continue
		}
		header, err := binevent.ParseHeader(eventData)
		if err != nil {
			fmt.Printf("解析失败：%v\n", err)
			continue
		}

		fmt.Printf("解析成功：%+v\n", header)

		switch header.EventType {
		case 15: // FORMAT_DESCRIPTION_EVENT
			// 格式: 2 bytes (binlog-version) + 50 bytes (server-version) + 4 bytes (created) + 1 byte (header-length)
			fdePayload := eventData[19:]
			if len(fdePayload) >= 2 {
				binevent.BinlogVersion = binary.LittleEndian.Uint16(fdePayload[0:2])
				fmt.Printf("👉 [FORMAT_DESC] 捕获到 Binlog 版本: %d\n", binevent.BinlogVersion)
			}
		case 19: // TABLE_MAP_EVENT
			// Post-Header 长度通常为 8 字节 (table_id: 6, flags: 2)
			// 注意：实际长度可能因 MySQL 版本而异，这里按常见情况处理
			payload := eventData[19:]
			if len(payload) < 8 {
				continue
			}
			tableID := binary.LittleEndian.Uint64(append(payload[0:6], 0, 0)) // 6字节转uint64
			// 跳过 flags (2 bytes)
			pos := 8

			// 解析 Schema Name (长度前缀字符串)
			schemaLen := int(payload[pos])
			pos++
			schemaName := string(payload[pos : pos+schemaLen])
			pos += schemaLen + 1 // +1 跳过结尾的 \0

			// 解析 Table Name
			tableLen := int(payload[pos])
			pos++
			tableName := string(payload[pos : pos+tableLen])
			pos += tableLen + 1 // +1 跳过结尾的 \0

			// 跳过 column-count (packed integer)
			colCount, colCountBytes, err := parser.ReadPackedInteger(payload[pos:])
			if err != nil {
				continue
			}
			pos += colCountBytes // 正确跳过 column-count 占用的字节数

			// 提取列类型数组
			colTypes := payload[pos : pos+int(colCount)]
			pos += int(colCount)

			// 【关键】解析 Column Metadata
			// 1. 读取 Metadata 的总字节数 (Packed Integer)
			metaLen, metaLenBytes, err := parser.ReadPackedInteger(payload[pos:])
			if err != nil {
				continue
			}
			pos += metaLenBytes

			// 2. 提取 Metadata 字节切片
			metaBytes := payload[pos : pos+int(metaLen)]
			pos += int(metaLen)

			// 3. 将 Metadata 字节转换为 uint16 数组
			colMetas := make([]uint16, colCount)
			metaPos := 0

			for i := 0; i < int(colCount); i++ {
				switch colTypes[i] {
				case 15: // MYSQL_TYPE_VARCHAR
					if metaPos+2 <= len(metaBytes) {
						colMetas[i] = binary.LittleEndian.Uint16(metaBytes[metaPos : metaPos+2])
						metaPos += 2
					}
				case 18: // MYSQL_TYPE_DATETIME / MYSQL_TYPE_DATETIME2
					// DATETIME 的 metadata 占 1 字节（小数秒精度）
					if metaPos+1 <= len(metaBytes) {
						colMetas[i] = uint16(metaBytes[metaPos])
						metaPos += 1
					}
				case 246: // MYSQL_TYPE_DECIMAL
					// DECIMAL 的 metadata 占 2 字节
					if metaPos+2 <= len(metaBytes) {
						colMetas[i] = binary.LittleEndian.Uint16(metaBytes[metaPos : metaPos+2])
						metaPos += 2
					}
				case 252: // MYSQL_TYPE_BLOB / TEXT
					// BLOB 的 metadata 占 1 字节
					if metaPos+1 <= len(metaBytes) {
						colMetas[i] = uint16(metaBytes[metaPos])
						metaPos += 1
					}
				default:
					// INT, BIGINT, TINYINT 等确实没有 metadata，什么都不做
				}
			}

			binevent.TableMapCache[tableID] = &binevent.TableMapInfo{
				SchemaName: schemaName,
				TableName:  tableName,
				ColTypes:   colTypes,
				ColMetas:   colMetas,
			}
			fmt.Printf("👉 [TABLE_MAP] 缓存表: %s.%s (ID: %d), 列数: %d\n", schemaName, tableName, tableID, colCount)

		case 30, 31, 32: // WRITE_ROWS_EVENT, UPDATE_ROWS_EVENT, DELETE_ROWS_EVENT
			rowPayload := eventData[19:]
			if len(rowPayload) < 6 {
				continue
			}

			// 1. 动态获取 Post-Header 长度
			postHeaderLen := 8 // 默认 MySQL 5.6+
			if binevent.BinlogVersion < 2 {
				postHeaderLen = 6
			}

			// 2. 正确提取 tableID
			var tableID uint64
			if postHeaderLen == 8 {
				tableID = binary.LittleEndian.Uint64(append(rowPayload[0:6], 0, 0))
			} else {
				tableID = uint64(binary.LittleEndian.Uint32(rowPayload[0:4]))
			}

			tblInfo, ok := binevent.TableMapCache[tableID]
			if !ok {
				fmt.Printf("⚠️ 警告: 找不到 TableID %d 的映射\n", tableID)
				continue
			}

			pos := postHeaderLen
			colCount := len(tblInfo.ColTypes)
			// 【核心】Bitmap 长度计算绝对不能 +2
			nullBitmapLen := (colCount + 7) / 8

			imageCount := 1
			if header.EventType == 31 { // UPDATE
				imageCount = 2
			}

			for img := 0; img < imageCount; img++ {
				if len(rowPayload) < pos+nullBitmapLen {
					fmt.Println("⚠️ 警告: 数据长度不足以读取 Null Bitmap")
					goto NEXT_EVENT
				}
				currentNullBitmap := rowPayload[pos : pos+nullBitmapLen]
				pos += nullBitmapLen

				imageType := "BEFORE_IMAGE"
				if img == 1 {
					imageType = "AFTER_IMAGE"
				}
				fmt.Printf("👉 [%s] 目标表: %s.%s\n", imageType, tblInfo.SchemaName, tblInfo.TableName)

				parseFailed := false
				for i := 0; i < colCount; i++ {
					// 检查是否为 NULL (offset=2 是因为 Row Event Bitmap 前2位保留)
					if parser.IsNullBitmapSet(currentNullBitmap, i, 2) {
						fmt.Printf("   列 %d: NULL\n", i)
						continue // NULL 不占用数据空间，直接跳过
					}

					switch tblInfo.ColTypes[i] {
					case 8: // MYSQL_TYPE_LONGLONG (BIGINT)
						if len(rowPayload) < pos+8 {
							parseFailed = true
							break
						}
						val := binary.LittleEndian.Uint64(rowPayload[pos : pos+8])
						pos += 8
						fmt.Printf("   列 %d (BIGINT): %d\n", i, val)

					case 3: // MYSQL_TYPE_INT
						if len(rowPayload) < pos+4 {
							parseFailed = true
							break
						}
						val := binary.LittleEndian.Uint32(rowPayload[pos : pos+4])
						pos += 4
						fmt.Printf("   列 %d (INT): %d\n", i, val)

					case 1: // MYSQL_TYPE_TINY
						if len(rowPayload) < pos+1 {
							parseFailed = true
							break
						}
						val := int8(rowPayload[pos])
						pos += 1
						fmt.Printf("   列 %d (TINYINT): %d\n", i, val)

					case 15: // MYSQL_TYPE_VARCHAR
						maxBytes := tblInfo.ColMetas[i]
						var lenBytes int
						var strLen uint64
						if maxBytes <= 255 {
							if len(rowPayload) < pos+1 {
								parseFailed = true
								break
							}
							strLen = uint64(rowPayload[pos])
							lenBytes = 1
						} else {
							if len(rowPayload) < pos+2 {
								parseFailed = true
								break
							}
							strLen = uint64(binary.LittleEndian.Uint16(rowPayload[pos : pos+2]))
							lenBytes = 2
						}
						pos += lenBytes
						if len(rowPayload) < pos+int(strLen) {
							parseFailed = true
							break
						}
						strVal := string(rowPayload[pos : pos+int(strLen)])
						pos += int(strLen)
						fmt.Printf("   列 %d (VARCHAR): %s\n", i, strVal)

					case 18: // MYSQL_TYPE_DATETIME2 (MySQL 5.6+)
						fsp := int(tblInfo.ColMetas[i])
						dtLen := 4 + (fsp+1)/2
						if len(rowPayload) < pos+dtLen {
							parseFailed = true
							break
						}

						// 调用解码函数
						dtBytes := rowPayload[pos : pos+dtLen]
						t, err := parser.DecodeDatetime2(dtBytes, fsp)
						if err != nil {
							fmt.Printf("   列 %d (DATETIME2): 解析失败 %v, 原始字节: %v\n", i, err, dtBytes)
						} else {
							// 格式化为与数据库一致的字符串
							var formatted string
							if fsp == 0 {
								formatted = t.Format("2006-01-02 15:04:05")
							} else {
								formatted = t.Format("2006-01-02 15:04:05") + fmt.Sprintf(".%0*d", fsp, t.Nanosecond()/int(time.Microsecond))
							}
							fmt.Printf("   列 %d (DATETIME2): %s\n", i, formatted)
						}
						pos += dtLen

					default:
						fmt.Printf("   列 %d: 未实现解析的类型 (%d)\n", i, tblInfo.ColTypes[i])
						parseFailed = true
					}

					if parseFailed {
						fmt.Printf("⚠️ 警告: 解析列 %d 时数据越界，跳过当前 Image 剩余列\n", i)
						goto NEXT_IMAGE
					}
				}
			NEXT_IMAGE:
			}
			goto NEXT_EVENT // 无论解析是否成功，都跳到下一个 Event，避免 pos 错乱影响后续

		NEXT_EVENT:

			// 3. 根据事件类型打印提示
			var action string
			switch header.EventType {
			case 30:
				action = "INSERT"
			case 31:
				action = "UPDATE"
			case 32:
				action = "DELETE"
			}
			fmt.Printf("👉 [ROW_EVENT] 动作: %s, 目标表: %s.%s\n", action, tblInfo.SchemaName, tblInfo.TableName)

		default:
			fmt.Printf("👉 收到其他事件 (Type: %d)\n", header.EventType)
		}

	}
}

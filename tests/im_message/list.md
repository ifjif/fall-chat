# 拉取会话消息(基于seq cursor的分页)

找已有消息，cursor分页

## 请求

      -X POST

      url = http://localhost:8080/api/v1/imMessage/list

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkyNjYzNTgzLCJuYmYiOjE3OTI2NjE3ODMsImlhdCI6MTc5MjY2MTc4M30.vcEao35UO5mf6d-pBaDXaGNlERalTr7XtDpMoRpk71s"

      -d @list-req.json
      -o list-res.json

- **json**

      {
        "conv_id": "7a810027038bdbc8",
        "next": "",
        "sort_key": "seq_desc" // seq_desc, seq_asc
      }

## 错误1

    {
      "code": 2003,
      "message": "未知排序行为",
      "status": 400
    }

## 错误2(maybe next解码失败)

    {
      "code": 9999,
      "message": "服务器内部错误",
      "status": 500
    }

## 成功无消息

    {
      "code": 200,
      "message": "",
      "data": {
        "data": [],
        "next": "eyJ2IjoxLCJwIjp7Imhhc19uZXh0IjpmYWxzZSwiU2VxIjotMSwiSUQiOjEwNzU3NTQ2NzIwNDE1MzM0NH19",
        "has_next": false
      }
    }

## 成功有消息

      {
        "code": 200,
        "message": "",
        "data": {
          "data": [
            {
              "id": "7a4b0a97aacbdbca",
              "conv_id": "7a810027038bdbc8",
              "sender_id": "7aeabd8c21cbdbca",
              "msg_type": 4,
              "content": "{\"url\":\"files/fc6020cc11695345e5c228d2008a7a12/Nacos安装指南.md\",\"name\":\"Nacos安装指南.md\",\"size\":3898,\"type\":\"text/markdown\",\"text\":\"这是一段文件\"}",
              "status": "success",
              "send_time": "2026-10-26T12:21:32.164Z",
              "seq": 143
            },
            {
              "id": "7a4b7474654bdbca",
              "conv_id": "7a810027038bdbc8",
              "sender_id": "7aeabdf9ad0bdbca",
              "msg_type": 4,
              "content": "{\"url\":\"files/e94a15a085e3cc83a79c23e1b53e0acb/new 6.txt\",\"name\":\"new 6.txt\",\"size\":5731,\"type\":\"text/plain\",\"text\":\"\"}",
              "status": "success",
              "send_time": "2026-10-26T12:14:29.563Z",
              "seq": 142
            },
            {
              "id": "7a4b73e0828bdbca",
              "conv_id": "7a810027038bdbc8",
              "sender_id": "7aeabdf9ad0bdbca",
              "msg_type": 3,
              "content": "{\"url\":\"files/706804877f7eefe8a63865d1ef8a0e71/video_1793014846078_5ycvke2\",\"cover_url\":\"files/a46f71229d41f5daf485198b36414491/cover_video_1793014846078_5ycvke2\",\"width\":2520,\"height\":1384,\"size\":998699,\"type\":\"video/webm\",\"name\":\"video_1793014846078_5ycvke2\",\"duration\":26.289,\"text\":\"\"}",
              "status": "success",
              "send_time": "2026-10-26T11:41:22.598Z",
              "seq": 141
            },
            {
              "id": "7a4b6097530bdbca",
              "conv_id": "7a810027038bdbc8",
              "sender_id": "7aeabdf9ad0bdbca",
              "msg_type": 3,
              "content": "{\"url\":\"files/79cf9212469828a077a2c92e5f57ff27/video_1793011521524_d2kz05x\",\"cover_url\":\"files/d0ce4da49ddda4a3605133ff49f800d2/cover_video_1793011521524_d2kz05x\",\"width\":2520,\"height\":1384,\"size\":2171085,\"type\":\"video/webm\",\"name\":\"video_1793011521524_d2kz05x\",\"duration\":12.642,\"text\":\"\"}",
              "status": "success",
              "send_time": "2026-10-26T10:45:25.219Z",
              "seq": 140
            }
          ],
          "next": "eyJ2IjoxLCJwIjp7Imhhc19uZXh0Ijp0cnVlLCJTZXEiOjE0MCwiSUQiOjEwODI3NDgwNTI0NDk1NjY3Mn19",
          "has_next": true
        }
      }

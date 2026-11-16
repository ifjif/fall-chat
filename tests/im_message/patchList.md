# 找seq范围内的消息(已拉取最新后 之后最新的)

在已有最新消息后，补后来的最新消息(最多50,超过不应该走它而是走list重新从最新拉取)

## 请求

      -X POST

      url = http://localhost:8080/api/v1/imMessage/patchList

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkzMjg1ODQwLCJuYmYiOjE3OTMyODQwNDAsImlhdCI6MTc5MzI4NDA0MH0.5ZKeN6nsdty3u43LBmpGWyqzL70fuykR-86x_khTx2Q"

      -d @patchList-req.json
      -o patchList-res.json

- **json**

      {
        "conv_id": "7a810027038bdbc8",
        "min_seq": 130,
        "max_seq": 141
      }

## 错误1(max<min / (max-min) > 50)

    {
      "code": 400,
      "message": "Bad Request",
      "status": 400
    }

## 成功

      {
        "code": 200,
        "message": "",
        "data": [
          {
            "id": "7a4b6097530bdbca",
            "conv_id": "7a810027038bdbc8",
            "sender_id": "7aeabdf9ad0bdbca",
            "msg_type": 3,
            "content": "{\"url\":\"files/79cf9212469828a077a2c92e5f57ff27/video_1793011521524_d2kz05x\",\"cover_url\":\"files/d0ce4da49ddda4a3605133ff49f800d2/cover_video_1793011521524_d2kz05x\",\"width\":2520,\"height\":1384,\"size\":2171085,\"type\":\"video/webm\",\"name\":\"video_1793011521524_d2kz05x\",\"duration\":12.642,\"text\":\"\"}",
            "status": "success",
            "send_time": "2026-10-26T10:45:25.219Z",
            "seq": 140
          },
          {
            "id": "7a4b62941f8bdbca",
            "conv_id": "7a810027038bdbc8",
            "sender_id": "7aeabdf9ad0bdbca",
            "msg_type": 3,
            "content": "{\"url\":\"files/852f60d8f1095758745eb266979a4192/video_1793010985301_fiqx5jz\",\"cover_url\":\"files/26d2f9c63073c30008fdd901eabb9011/cover_video_1793010985301_fiqx5jz\",\"width\":640,\"height\":480,\"size\":1247525,\"type\":\"video/webm\",\"name\":\"video_1793010985301_fiqx5jz\",\"duration\":8.419,\"text\":\"这是录制的视频\"}",
            "status": "success",
            "send_time": "2026-10-26T10:36:44.241Z",
            "seq": 139
          }
        ]
      }

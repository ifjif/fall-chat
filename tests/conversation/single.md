# 请求，需auth

      -X GET

      url = http://localhost:8080/api/v1/conversation/single

      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg4ODY4OTUyLCJuYmYiOjE3ODg4NjcxNTIsImlhdCI6MTc4ODg2NzE1Mn0.5Xgby0NQRfZ7jDX-vBZ6SyCl_OR9f9gbGjFDYc8xZJ8"

      -G
      -d "conv_id=7a895f6ce1cbdbca"

# 成功，存在

      {
        "code": 200,
        "message": "",
        "data": {
          "id": "7a895f6ce1cbdbca",
          "name": "我是李四",
          "avatar": "files/80fea4a632bb159ef84be63a4ae25a2d/avatar55.webp",
          "type": "peer",
          "peer_id": "7aeabd8c21cbdbca",
          "last_msg_id": 90799988350652416,
          "last_msg_preview": "hello,hello,hello,hello,hello",
          "last_msg_sender_id": 81476227113684992,
          "last_msg_time": "2026-09-08T05:26:44.205Z",
          "last_msg_seq": 2,
          "last_read_seq": 2,
          "unread_count": 0
        }
      }

# 成功，不存在

      {
        "code": 200,
        "message": "",
        "data": null
      }

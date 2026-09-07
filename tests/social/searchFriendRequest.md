# 请求 POST 需auth

      -X POST

      url = http://localhost:8080/api/v1/social/searchFriendRequest

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3NTY2MjA0LCJuYmYiOjE3ODc1NjQ0MDQsImlhdCI6MTc4NzU2NDQwNH0.dIeq5rF9iSxr-xzSia0BZavX-qzcEqvHYLfU6T1HD7E"

- **json**

      {
        "filter_status": [] // 0:待处理 1:已通过 2:已拒绝 3:已过期
      }

# 成功

      {
        "code": 200,
        "message": "",
        "data": [
          {
            "id": "7ae5843e850bdbca",
            "from_user_id": "7aeabdf9ad0bdbca",
            "to_user_id": "7aeabd8c21cbdbca",
            "sender": true,
            "message": "你好李四，我是张三",
            "status": "pending",
            "avatar": "files/80fea4a632bb159ef84be63a4ae25a2d/avatar55.webp",
            "nickname": "我是李四",
            "apply_type": "search",
            "created_at": "2026-08-23T11:30:37.5Z"
          }
        ]
      }

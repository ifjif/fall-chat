# 请求 POST 需auth

      -X POST

      url = http://localhost:8080/api/v1/social/applyAddFriend

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3NDg1NTUyLCJuYmYiOjE3ODc0ODM3NTIsImlhdCI6MTc4NzQ4Mzc1Mn0.ole6tV-rzHPJEZKCK0vJekkTscRukKyNqCK1Bd2fbj0"

- **json**

      {
        "target_id": "7aeabd8c21cbdbca",
        "greeting": "你好李四，我是张三",
        "apply_type": 0
      }

# 成功

      {
        "code": 200,
        "message": "发送成功"
      }

# 已存在（如待处理/已同意）

      {
        "code": 4802,
        "message": "好友申请已存在",
        "status": 400
      }

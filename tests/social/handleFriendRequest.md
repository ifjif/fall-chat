# 请求 POST

      -X POST

      url = http://localhost:8080/api/v1/social/handleFriendRequest

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkOGMyMWNiZGJjYSIsInVzZXJuYW1lIjoibGlzaGkiLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3NTgzNDQzLCJuYmYiOjE3ODc1ODE2NDMsImlhdCI6MTc4NzU4MTY0M30.XoFMBD0dSk_-NapTp-ROnYMYMeuFLqWB821YhSnowls"

- **json**

      {
        "id": "7ae456846e4bdbca",
        "action": 1 // 1:同意 2:拒绝
      }

# 成功

      {
        "code": 200,
        "message": "处理成功"
      }

# 错误1

      {
        "code": 4803,
        "message": "好友申请不存在",
        "status": 400
      }

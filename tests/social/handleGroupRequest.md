# 请求 POST 需auth

      -X POST

      url = http://localhost:8080/api/v1/social/handleGroupRequest

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkOGMyMWNiZGJjYSIsInVzZXJuYW1lIjoibGlzaGkiLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3ODQwMDg2LCJuYmYiOjE3ODc4MzgyODYsImlhdCI6MTc4NzgzODI4Nn0.IeX1UR5S3yN1hiDhTFWDFdX8SgS6fivXJfgFt6Oz0nU"

- **json**

      {
        "id": "7afad9e97d0bdbca",
        "action": 2 // 1: approved 2: rejected
      }

# 成功

      {
        "code": 200,
        "message": "操作成功"
      }

# 错误1

      {
        "code": 4814,
        "message": "无法处理自己的群申请",
        "status": 400
      }

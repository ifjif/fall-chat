# 请求 POST 需auth

      -X POST

      url = http://localhost:8080/api/v1/groupInfo/groups

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3NjU5MTQ1LCJuYmYiOjE3ODc2NTczNDUsImlhdCI6MTc4NzY1NzM0NX0.MQwssmuLfIIrvEicoRZDsnFXEPzGklY78cmRpSq5Pzo"

- **json**

      {
        "filter": [] // 0:普通 1:管理员 2:群主
      }

# 成功

      {
        "code": 200,
        "message": "",
        "data": [
          {
            "id": "7afb2381c90bdbca",
            "name": "go语言交流群1",
            "avatar": "files/b455771725f50432a89fc253e2b755ef/20250612151218.jpg",
            "description": "讨论go语言的交流群"
          }
        ]
      }

# 请求 需auth

-X POST

url = <http://localhost:8080/api/v1/file/chunkStatus>

-H "Content-Type: application/json"
-H "Accept: application/json"
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3MzkwOTY4LCJuYmYiOjE3ODczODkxNjgsImlhdCI6MTc4NzM4OTE2OH0.GWSjkYWncUIirVNxN8xeBh9m7L_j9CClmMODJg6YdQI"

# 成功

      {
        "code": 200,
        "message": "",
        "data": [
          "0",
          "1"
        ]
      }

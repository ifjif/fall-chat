# 获取一个表情包的信息

## 请求

      -X GET

      url = http://localhost:8080/api/v1/emoji/pack/7a49da1daf4bdbca

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkzMTA3MjYwLCJuYmYiOjE3OTMxMDU0NjAsImlhdCI6MTc5MzEwNTQ2MH0.4db7U9xknaea97qG4E2FBhy8If6s8oZIRmG-F4FTKiU"

      -o packDetail-res.json

## 成功

      {
        "code": 200,
        "message": "",
        "data": {
          "id": "7a49da1daf4bdbca",
          "row": "张三表情包1",
          "cover_url": "",
          "owner_id": "7aeabdf9ad0bdbca",
          "soruce": "customer",
          "price_type": "free",
          "price": 0,
          "status": "draft",
          "created_at": "0001-01-01T00:00:00Z"
        }
      }

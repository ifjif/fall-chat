# 列出所有表情包（基于游标的分页）

## 请求

      -X POST

      url = http://localhost:8080/api/v1/emoji/listPack

      -H "Content-type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkzMTA3MjYwLCJuYmYiOjE3OTMxMDU0NjAsImlhdCI6MTc5MzEwNTQ2MH0.4db7U9xknaea97qG4E2FBhy8If6s8oZIRmG-F4FTKiU"

      -d @listPack-req.json
      -o listPack-res.json

- **json**

      {
        "next": "",
        "sort_key": "created_desc", // created_desc,created_asc
        "price_type": 1,            // 1免费 2付费
        "status": 1                 // 0:下架 1:草稿 2:上架 3:审核中
      }

## 错误1

      {
        "code": 2003,
        "message": "未知排序行为",
        "status": 400
      }

## 错误2 (可能是 next解码失败)

      {
        "code": 9999,
        "message": "服务器内部错误",
        "status": 500
      }

## 成功没有数据

      {
        "code": 200,
        "message": "",
        "data": {
          "data": null,
          "next": "eyJ2IjoxLCJwIjp7Imhhc19uZXh0IjpmYWxzZSwiaWQiOjEwODY2ODk1MzUxNzEwMTA1Nn19",
          "has_next": false
        }
      }

## 成功有数据

      {
        "code": 200,
        "message": "",
        "data": {
          "data": [
            {
              "id": "7a49da1daf4bdbca",
              "name": "张三表情包1",
              "cover_url": "",
              "owner_id": "7aeabdf9ad0bdbca",
              "soruce": "customer",
              "price_type": "free",
              "price": 0,
              "status": "added",
              "created_at": "0001-01-01T00:00:00Z"
            }
          ],
          "next": "eyJ2IjoxLCJwIjp7Imhhc19uZXh0IjpmYWxzZSwiaWQiOjEwODY2ODk1MzUxNzEwMTA1Nn19",
          "has_next": false
        }
      }

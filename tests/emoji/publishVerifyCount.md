# 发布表情包 的item数量验证

## 请求

      -X POST

      url = http://localhost:8080/api/v1/emoji/publish/verifyCount

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkzODc0MjkwLCJuYmYiOjE3OTM4NzI0OTAsImlhdCI6MTc5Mzg3MjQ5MH0._GUH9BcfGuFYbOrBB0wH4ZOfmQ9x2dSif8IGjNvqVyY"

      -d @publishVerifyCount-req.json
      -o publishVerifyCount-res.json

- **json**

      {
        "pack_id": "7a46b43dcf8bdbca",
        "item_count": 5
      }

## 成功

      {
        "code": 200,
        "message": ""
      }

## 错误1

      {
        "code": 5001,
        "message": "推送量太多:最多4",
        "status": 400
      }

## 错误2(不是自己的，或已推送)

      {
        "code": 5000,
        "message": "表情包不存在",
        "status": 400
      }

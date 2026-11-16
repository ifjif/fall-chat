# 创建表情包

## 请求 POST

      -X POST

      url = http://localhost:8080/api/v1/emoji/createPack

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkzMDkzMDI1LCJuYmYiOjE3OTMwOTEyMjUsImlhdCI6MTc5MzA5MTIyNX0.zmPU5BwE1T7KcDWJxCpRG0uEroVYD_8IK77nWzI11v8"

      -d @createPack-req.json
      -o createPack-res.json

- **json**

      {
        "name": "xx",
        "source": 1,      // 1系统 2用户
        "price_type":1,  // 1免费 2付费
        "price": 0,
        "cover_url": ""
      }

## 成功

    {
      "code": 200,
      "message": "创建成功"
    }

## 失败1，请求参数错误

    {
      "code": 4000,
      "message": "请求参数错误",
      "status": 400
    }

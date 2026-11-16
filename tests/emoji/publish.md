# 发布一个表情包

需要前置 数量验证

## 请求

      -X POST

      url = http://localhost:8080/api/v1/emoji/publish

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkzMTU2MjQ0LCJuYmYiOjE3OTMxNTQ0NDQsImlhdCI6MTc5MzE1NDQ0NH0.kENPuBifZG_Oh_uhnZp8ltHU-hBUV4dCAnXlRjO7iQk"

      -d @publish-req.json
      -o publish-res.json

- **json**

      {
        "id": "7a49da1daf4bdbca",
        "items": [
          {
            "media_type": 4,        // 1静态图 2动态图 3字体 4unicode
            "font_name": "1F600",
            "media_url": "",
            "thumb_url": "",
            "width": 0,
            "height": 0,
            "file_size": 0,
            "sort_no": 0
          },
          {
            "media_type": 4,
            "font_name": "1F601",
            "media_url": "",
            "thumb_url": "",
            "width": 0,
            "height": 0,
            "file_size": 0,
            "sort_no": 1
          }
        ]
      }

## 成功

      {
        "code": 200,
        "message": "发布成功"
      }

## 错误1(不是自己表情包，或状态已是发布)

      {
        "code": 5000,
        "message": "表情包不存在",
        "status": 400
      }

## 错误2(数量验证失败)

    {
      "code": 5002,
      "message": "推送量不匹配",
      "status": 400
    }

# 好友申请

## 请求

    -X POST

    url = http://localhost:8080/api/v1/friendRequest/apply

    -H "Content-Type: application/json"
    -H "Accept: application/json"
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkOTkzMTBiZGJjYSIsInVzZXJuYW1lIjoid2FuZ3d1IiwiZGV2aWNlX2lkIjoiIiwianRpIjoiIiwidHlwZSI6ImFjY2VzcyIsImlzcyI6ImZhbGwtY2hhdCIsImV4cCI6MTc5MjQ5OTY2NywibmJmIjoxNzkyNDk3ODY3LCJpYXQiOjE3OTI0OTc4Njd9.VUl6UkmZ0E20ITVoXufibFg19VV-g1GsETgb7duBQ4o"

-d @apply-req.json
-o apply-res.json

- **json**

      {
        "target_id": "7aeabd8c21cbdbca",
        "greeting": "你好李四，我是张三,再次申请",
        "apply_type": 0
      }

## 成功

      {
        "code": 200,
        "message": "发送成功"
      }

## 已存在（如待处理/已同意）

      {
        "code": 4802,
        "message": "好友申请已存在",
        "status": 400
      }

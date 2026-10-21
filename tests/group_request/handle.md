# 入群申请处理

## 请求

    -X POST

    url = http://localhost:8080/api/v1/groupRequest/handle

    -H "Content-Type: application/json"
    -H "Accept: application/json"
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkOGMyMWNiZGJjYSIsInVzZXJuYW1lIjoibGlzaGkiLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkyNDk5MDUxLCJuYmYiOjE3OTI0OTcyNTEsImlhdCI6MTc5MjQ5NzI1MX0._eGsNf5w-Kh0E3rc85tHBhA1zBZTj3fYVHYjqtS0SY0"

    -d @handle-req.json
    -o handle-res.json

- **json**

      {
        "id": "7afad9e97d0bdbca",
        "action": 2 // 1: approved 2: rejected
      }

## 成功

      {
        "code": 200,
        "message": "操作成功"
      }

## 错误1

      {
        "code": 4814,
        "message": "无法处理自己的群申请",
        "status": 400
      }

## 错误2

    {
      "code": 4812,
      "message": "群申请不存在",
      "status": 400
    }

## 错误3

    {
      "code": 4815,
      "message": "未加入群，无法处理",
      "status": 400
    }

## 错误4

    {
      "code": 4813,
      "message": "群申请已处理或过期",
      "status": 400
    }

## 错误5

    {
      "code": 4816,
      "message": "不是管理员，无法处理",
      "status": 400
    }

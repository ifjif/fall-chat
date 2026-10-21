# 好友申请处理

## 请求

    -X POST

    url = http://localhost:8080/api/v1/friendRequest/handle

    -H "Content-Type: application/json"
    -H "Accept: application/json"
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkOGMyMWNiZGJjYSIsInVzZXJuYW1lIjoibGlzaGkiLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkyNTAwNTA4LCJuYmYiOjE3OTI0OTg3MDgsImlhdCI6MTc5MjQ5ODcwOH0.8QOR963-1S1vs4zIp8wdPGVMJB7lsgStl83QNEwLUJY"

    -d @handle-req.json
    -o handle-res.json

- **json**

      {
        "id": "7ae456846e4bdbca",
        "action": 1 // 1:同意 2:拒绝
      }

## 成功

      {
        "code": 200,
        "message": "处理成功"
      }

## 错误1

      {
        "code": 4803,
        "message": "好友申请不存在",
        "status": 400
      }

## 错误2

    {
      "code": 4804,
      "message": "好友申请已处理或过期",
      "status": 400
    }

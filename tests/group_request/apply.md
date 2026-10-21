# 入群申请

## 请求

      -X POST

      url = http://localhost:8080/api/v1/social/applyAddGroup

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3NjcxMDY1LCJuYmYiOjE3ODc2NjkyNjUsImlhdCI6MTc4NzY2OTI2NX0.JVs_LXK2bEOTTliFObHd5e2MIxzKQ4GVrvMUIOP1GVQ"

- **json**

      {
        "target_id": "7afb2381c90bdbca",
        "greeting": "我是张三",
        "invited_id": "",   // 存在，生成方充当邀请者
        "apply_type": 0 // 0:search 1:qrcode 2:invited
      }

## 成功

    {
      "code": 200,
      "message": "发送成功"
    }

## 错误1

      {
        "code": 4811,
        "message": "无法申请已加入的群",
        "status": 400
      }

## 错误2

      {
        "code": 4805,
        "message": "错误的群申请类型",
        "status": 400
      }

## 错误3

      {
        "code": 4806,
        "message": "群申请无法邀请自己",
        "status": 400
      }

## 错误4

      {
        "code": 4818,
        "message": "邀请人未加入群，无法处理",
        "status": 400
      }

## 错误5

    {
      "code": 4807,
      "message": "群申请已存在",
      "status": 400
    }

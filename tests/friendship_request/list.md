# 找我的好友申请信息(基于cursor分页找)

## 请求

      -X POST

      url = http://localhost:8080/api/v1/friendRequest/list

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkyNDkxOTQ3LCJuYmYiOjE3OTI0OTAxNDcsImlhdCI6MTc5MjQ5MDE0N30.APpS8ozGAeqon_hM0zR7eb2s_AaGS3DXGwIJmsqe42E"

      -d @list-req.json
      -o list-res.json

- **json**

      {
        "status": [] // 0:待处理 1:已通过 2:已拒绝 3:已过期
        "next": "",
        "sort_key": "created_asc" //created_asc,created_desc
      }

## 错误1

      {
        "code": 2003,
        "message": "未知排序行为",
        "status": 400
      }

## 错误2 (可能是next解码失败)

    {
      "code": 9999,
      "message": "服务器内部错误",
      "status": 500
    }

## 成功 无数据

      {
        "code": 200,
        "message": "",
        "data": {
          "data": null,
          "next": "eyJ2IjoxLCJwIjp7Imhhc19uZXh0IjpmYWxzZSwiSUQiOjg1NzYwNzA3OTE1NjgxNzkyfX0=",
          "has_next": false
        }
      }

## 成功 有数据

      {
        "code": 200,
        "message": "",
        "data": {
          "data": [
            {
              "id": "7afb652318cbdbca",
              "from_user_id": "7aeabd99310bdbca",
              "to_user_id": "7aeabdf9ad0bdbca",
              "sender": false,
              "message": "请求添加好友",
              "status": "approved",
              "avatar": "files/4e572590aeed9869812e9cd7eeb014a1/avatar33.webp",
              "nickname": "王五xx",
              "apply_type": "search",
              "created_at": "2026-08-25T07:42:26.124Z"
            },
            {
              "id": "7afb656a358bdbca",
              "from_user_id": "7aeabd99310bdbca",
              "to_user_id": "7aeabdf9ad0bdbca",
              "sender": false,
              "message": "我是王五",
              "status": "rejected",
              "avatar": "files/4e572590aeed9869812e9cd7eeb014a1/avatar33.webp",
              "nickname": "王五xx",
              "apply_type": "search",
              "created_at": "2026-08-25T07:41:13.593Z"
            }
          ],
          "next": "eyJ2IjoxLCJwIjp7Imhhc19uZXh0Ijp0cnVlLCJJRCI6ODU3NjA0MDM2OTg2MTgzNjh9fQ==",
          "has_next": true
        }
      }

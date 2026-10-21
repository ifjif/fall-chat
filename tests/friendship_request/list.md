# 找我的好友申请信息

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
        "filter_status": [] // 0:待处理 1:已通过 2:已拒绝 3:已过期
      }

# 成功

      {
        "code":200,
        "message":"",
        "data":
          [{
            "id":"7afb652318cbdbca",
            "from_user_id":"7aeabd99310bdbca",
            "to_user_id":"7aeabdf9ad0bdbca",
            "sender":false,
            "message":"请求添加好友",
            "status":"approved",
            "avatar":"files/4e572590aeed9869812e9cd7eeb014a1/avatar33.webp",
            "nickname":"王五xx",
            "apply_type":"search",
            "created_at":"2026-08-25T07:42:26.124Z"
        }]
      }

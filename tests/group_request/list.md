# 找用户的群申请记录

## 请求

    -X POST

    url = http://localhost:8080/api/v1/groupRequest/list

    -H "Content-type: application/json"
    -H "Accept: application/json"
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkyNDkxOTQ3LCJuYmYiOjE3OTI0OTAxNDcsImlhdCI6MTc5MjQ5MDE0N30.APpS8ozGAeqon_hM0zR7eb2s_AaGS3DXGwIJmsqe42E"

    -d @list-req.json
    -o list-res.json

- **json**

      {
        "filter_status": [] // 0:pending 1:approved 2:rejected 3:expired
      }

# 成功

      {
        "code": 200,
        "message": "",
        "data": [
          {
            "id": "7afad9e97d0bdbca",
            "group_id": "7afada2fab8bdbca",
            "group_name": "rust交流群1",
            "group_avatar": "files/4d652c4788c247c2e6c9f07bca6e2273/avatar48.webp",
            "from_user_id": "7aeabdf9ad0bdbca",
            "from_user_name": "张三信息",
            "from_user_avatar": "files/eaa4463110b66813962913723317e941/avatar13.webp",
            "invite_user_id": "",
            "invite_user_name": "",
            "invite_user_avatar": "",
            "handle_by_user_id": "",
            "handle_by_user_name": "",
            "handle_by_user_avatar": "",
            "message": "我是张三",
            "status": "pending",
            "sender": true,
            "apply_type": "search",
            "created_at": "2026-08-25T14:55:57.66Z"
          }
        ]
      }

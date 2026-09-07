# 请求 POST 需auth

      -X POST

      url = http://localhost:8080/api/v1/social/searchGroupRequest

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3ODEwNzkyLCJuYmYiOjE3ODc4MDg5OTIsImlhdCI6MTc4NzgwODk5Mn0.CyD89dgzXz9AtAAAa6nUR9aLMrbGgngbwbppYSwOgn4"

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

# 找用户的群申请记录(基于cursor分页)

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
        "status": [] // 0:pending 1:approved 2:rejected 3:expired
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

## 成功，无数据

      {
        "code": 200,
        "message": "",
        "data": {
          "data": null,
          "next": "eyJ2IjoxLCJwIjp7Imhhc19uZXh0IjpmYWxzZSwiSUQiOjEwNjE0NTI4MjAxMzUzMjE2MH19",
          "has_next": false
        }
      }

## 成功有数据

      {
        "code": 200,
        "message": "",
        "data": {
          "data": [
            {
              "id": "7ab2cb35a78bdbca",
              "group_id": "7ab330fe310bdbca",
              "group_name": "lua交流群1",
              "group_avatar": "files/4d652c4788c247c2e6c9f07bca6e2273/avatar48.webp",
              "from_user_id": "7aeabd8c21cbdbca",
              "from_user_name": "我是李四",
              "from_user_avatar": "files/329c5e9ef9190d4655aecaa3dff8d2c4/截图 2026-10-07 12-51-16.png",
              "invite_user_id": "7aeabdf9ad0bdbca",
              "invite_user_name": "张三信息",
              "invite_user_avatar": "files/60356a57260f33ed5fbc475b8f3cec43/截图 2026-08-29 21-12-02.png",
              "handle_by_user_id": "7aeabdf9ad0bdbca",
              "handle_by_user_name": "张三信息",
              "handle_by_user_avatar": "files/60356a57260f33ed5fbc475b8f3cec43/截图 2026-08-29 21-12-02.png",
              "message": "我是张三",
              "status": "approved",
              "sender": false,
              "apply_type": "invited",
              "created_at": "2026-10-20T11:51:42.258Z"
            },
            {
              "id": "7ab33766680bdbca",
              "group_id": "7ab330fe310bdbca",
              "group_name": "lua交流群1",
              "group_avatar": "files/4d652c4788c247c2e6c9f07bca6e2273/avatar48.webp",
              "from_user_id": "7aeabd8c21cbdbca",
              "from_user_name": "我是李四",
              "from_user_avatar": "files/329c5e9ef9190d4655aecaa3dff8d2c4/截图 2026-10-07 12-51-16.png",
              "invite_user_id": "7aeabdf9ad0bdbca",
              "invite_user_name": "张三信息",
              "invite_user_avatar": "files/60356a57260f33ed5fbc475b8f3cec43/截图 2026-08-29 21-12-02.png",
              "handle_by_user_id": "7aeabdf9ad0bdbca",
              "handle_by_user_name": "张三信息",
              "handle_by_user_avatar": "files/60356a57260f33ed5fbc475b8f3cec43/截图 2026-08-29 21-12-02.png",
              "message": "我是张三",
              "status": "rejected",
              "sender": false,
              "apply_type": "invited",
              "created_at": "2026-10-20T11:32:50.959Z"
            }
          ],
          "next": "eyJ2IjoxLCJwIjp7Imhhc19uZXh0Ijp0cnVlLCJJRCI6MTA2MTEyNDEzOTUwMDIxNjMyfX0=",
          "has_next": true
        }
      }

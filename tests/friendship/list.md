# 找我的好友

## 请求

      -X POST

      url = http://localhost:8080/api/v1/friendship/list

      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkyNDkxOTQ3LCJuYmYiOjE3OTI0OTAxNDcsImlhdCI6MTc5MjQ5MDE0N30.APpS8ozGAeqon_hM0zR7eb2s_AaGS3DXGwIJmsqe42E"

      -o list-res.json

## 成功

        {
          "code": 200,
          "message": "",
          "data": [
            {
              "id": "7aeabd8c21cbdbca",
              "nickname": "我是李四",
              "avatar": "files/80fea4a632bb159ef84be63a4ae25a2d/avatar55.webp",
              "remark": "",
              "source": "search"
            },
            {
              "id": "7aeabd99310bdbca",
              "nickname": "王五xx",
              "avatar": "files/4e572590aeed9869812e9cd7eeb014a1/avatar33.webp",
              "remark": "",
              "source": "search"
            }
          ]
        }

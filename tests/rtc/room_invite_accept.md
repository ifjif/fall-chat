# rtc 接受房间邀请

## 请求

      -X POST

      url = http://localhost:8080/api/v1/rtc/room/invite/accept

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkOGMyMWNiZGJjYSIsInVzZXJuYW1lIjoibGlzaGkiLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzk0NzU0OTYxLCJuYmYiOjE3OTQ3NTMxNjEsImlhdCI6MTc5NDc1MzE2MX0.LJwy9Q__S_cf9dDK1TRyA0OYNio6nduukybFjDTfZKg"

      -d @room_invite_accept-req.json
      -o room_invite_accept-res.json

- **json**

      {
        "room_id": "room_81476227113684992_459822249",
        "inviter_id": "7aeabdf9ad0bdbca"
      }

## 成功

      {
        "code": 200,
        "message": ""
      }

## 错误1(已激活或不存在)

      {
        "code": 5401,
        "message": "房间不存在",
        "status": 400
      }

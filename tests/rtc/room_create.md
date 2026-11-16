# 创建一个 rtc 流房间

## 请求

      -X POST

      url = http://localhost:8080/api/v1/rtc/room/create

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzk0NzUyNTA2LCJuYmYiOjE3OTQ3NTA3MDYsImlhdCI6MTc5NDc1MDcwNn0._aaC1qB9BsU47ST0v9xMWlw53b6xFh4pINOuSV8SUbk"

      -d @room_create-req.json
      -o room_create-res.json

- **json**

      {
        "type": "one_to_one" // one_to_one / group
      }

## 成功

      {
        "code": 200,
        "message": "",
        "data": "room_81476227113684992_459822249"
      }

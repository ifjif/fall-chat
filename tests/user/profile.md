# 请求

      -X GET


      url = http://localhost:8080/api/v1/user/profile?token="xx"

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3MjExMDI4LCJuYmYiOjE3ODcyMDkyMjgsImlhdCI6MTc4NzIwOTIyOH0.2tF70dSbEKm5vpBtKH04sBM9PfIPP0c3RnTmDF__Qlc"

请求头携带 accessToken,用来刷新缓存

# 成功

        {
          "code": 200,
          "message": "",
          "data": {
            "id": "7aeabdf9ad0bdbca",
            "avatar": "",
            "username": "zhangsan",
            "nickname": "user7aeabdf9ad0bdbca",
            "email": "",
            "phone": "",
            "created_at": "2026-08-13T11:57:26.534Z",
            "status": 1,
            "search_flag": 1
          }
        }

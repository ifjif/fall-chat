# 请求 POST

      {
        "username": "zhangsan",   // min=3,max=50
        "password": "12345678"    // min=8,max=20
      }

# 成功

      {
        "code": 200,
        "data": {
          "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FkZTE2MWM3OTBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiJ1bmtub3duIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg1OTM1NDE1LCJuYmYiOjE3ODU4NDkwMTUsImlhdCI6MTc4NTg0OTAxNX0.gQOQahIk0QRjSRZGMVqlmYqpDfe2st-ksXf5pI8t9Fs",
          "expires_in": 86399,
          "user": {
            "id": "7ade161c790bdbca",
            "username": "zhangsan",
            "email": "",
            "created_at": "2026-08-04T11:48:30.259Z"
          }
        },
        "message": "登录成功"
      }

# 失败

      {
        "code": 4001,
        "message": "用户名或密码错误",
        "status": 404
      }

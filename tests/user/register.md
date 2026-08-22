# 请求 POST

      {
        "username": "zhangsan",      //min=3,max=50 唯一
        "password": "12345678"   //min=8,max=20
      }

# 成功

      {
        "id": "7ade161c790bdbca"
      }

# 失败

      {
        "code": 4003,
        "message": "用户已存在",
        "status": 400
      }

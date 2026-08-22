# 请求 POST

      {
        "username": "test",
        "email": "test@test.com",
        "new_pass": "test",
        "code": "test",
        "code_id": "test"
      }

# code_id 解码失败

      {
        "code": 9999,
        "message": "服务器内部错误",
        "status": 500
      }

# code 无效

      {
        "code": 4200,
        "message": "验证码失效",
        "status": 400
      }

# 邮箱错误

      {
        "code": 4003,
        "message": "邮箱错误",
        "status": 404
      }

# 用户名不正确

      {
        "code": 4000,
        "message": "用户名或密码错误",
        "status": 404
      }

# 成功

      {
        "code": 200,
        "message": "找回成功"
      }

# 请求 POST

      {
        "email": "test@test.com",
        "old_pass": "test",
        "new_pass": "test",
        "code": "test",
        "code_id": "test"
      }

# 未认证

      {
        "code": 401,
        "message": "缺少认证凭证",
        "status": 401
      }

# 不是 Bearer 认证格式

      {
        "code": 401,
        "message": "凭证格式错误,要求Bearer",
        "status": 401
      }

# 凭证无效

      {
        "code": 401,
        "message": "凭证无效",
        "status": 401
      }

# code id解码失败

      {
        "code": 9999,
        "message": "服务器内部错误",
        "status": 500
      }

# code 不正确

      {
        "code": 4200,
        "message": "验证码失效",
        "status": 400
      }

# 密码错误

      {
        "code": 4000,
        "message": "用户名或密码错误",
        "status": 404
      }

# 成功

      {
        "code": 200,
        "message": "更新成功"
      }

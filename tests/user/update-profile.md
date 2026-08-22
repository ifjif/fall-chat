# 请求 POST

      {
        "email": "test@test.com",
        "nickname": "test",
        "password": "test"
      }

# 未认证

      {
        "code": 401,
        "message": "缺少认证凭证",
        "status": 401
      }

# 非 Bearer方案

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

# 请求Body

    Name        string `json:"name" validate:"min=1,max=100"`
    Description string `json:"description" validate:"max=255"`
    Avatar      string `json:"avatar" validate:"max=255"`
    MaxMembers  int    `json:"max_members" validate:"min=2"`
    SearchFlag  int    `json:"search_flag" validate:"oneof=0 1"`
    JoinMode    int    `json:"join_mode" validate:"oneof=0 1 2"`

# 参数错误

      {
        "code": 4000,
        "message": "请求参数错误",
        "status": 400
      }

# 缺少认证凭证

      {
        "code": 401,
        "message": "缺少认证凭证",
        "status": 401
      }

# 非Bearer认证方案

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

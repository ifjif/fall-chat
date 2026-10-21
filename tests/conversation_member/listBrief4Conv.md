# 获取会话的基本信息(已弃用)

```
{
id
名字
在群中的名字
头像
}
```

# 请求 需auth

      -X GET

      url = http://localhost:8080/api/v1/conversationMember/listBrief4Conv

      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg4OTU3NTY5LCJuYmYiOjE3ODg5NTU3NjksImlhdCI6MTc4ODk1NTc2OX0.CiIAL9a1Rj7K8j-m3iJNx2HHKrWF5UoL0_Bvb05su5I"

      -G
      -d "conv_id=7a8fc1ebdfcbdbca"

# 成功

      {
        "code": 200,
        "message": "",
        "data": {
          "7aeabd8c21cbdbca": {
            "ID": "7aeabd8c21cbdbca",
            "OriginNickname": "我是李四",
            "InGropNickname": "",
            "Avatar": "files/80fea4a632bb159ef84be63a4ae25a2d/avatar55.webp"
          },
          "7aeabdf9ad0bdbca": {
            "ID": "7aeabdf9ad0bdbca",
            "OriginNickname": "张三信息",
            "InGropNickname": "",
            "Avatar": "files/60356a57260f33ed5fbc475b8f3cec43/截图 2026-08-29 21-12-02.png"
          }
        }
      }

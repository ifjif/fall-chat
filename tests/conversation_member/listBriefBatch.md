# 获取用户在聊天中的基本信息

```
通过批次用户id [id1,id2,...]
返回
{
  id
  名字
  在群中的名字
  头像
}
```

# 请求 需auth

      -X POST

      url = http://localhost:8080/api/v1/conversationMember/listBriefBatch

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg5MDQzNDE2LCJuYmYiOjE3ODkwNDE2MTYsImlhdCI6MTc4OTA0MTYxNn0.5-nX-EYOOi2JDDLrQub_kLwbAfxcqs7UPnUMshS5JaI"

**json**

      {
        "conv_id": "7a8fc1ebdfcbdbca",
        "member_ids": [
          "7aeabdf9ad0bdbca",
          "7aeabd8c21cbdbca"
        ]
      }

# 成功

      {
        "code": 200,
        "message": "",
        "data": {
          "7aeabd8c21cbdbca": {
            "id": "7aeabd8c21cbdbca",
            "nickname": "我是李四",
            "in_group_nickname": "",
            "avatar": "files/80fea4a632bb159ef84be63a4ae25a2d/avatar55.webp"
          },
          "7aeabdf9ad0bdbca": {
            "id": "7aeabdf9ad0bdbca",
            "nickname": "张三信息",
            "in_group_nickname": "",
            "avatar": "files/60356a57260f33ed5fbc475b8f3cec43/截图 2026-08-29 21-12-02.png"
          }
        }
      }

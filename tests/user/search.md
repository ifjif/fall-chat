# 请求 POST 需auth

-X POST

url = <http://localhost:8080/api/v1/user/search>

-H "Content-Type: application/json"
-H "Accept: application/json"
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3Mzk3NDM1LCJuYmYiOjE3ODczOTU2MzUsImlhdCI6MTc4NzM5NTYzNX0.NNzMdvxP_Hn-DZUQBBhjkEpvOgZFh556LNFBXKZgn2U"

json:

    {
      "name": "sh",
      "next_cursor": ""
    }

# 成功

## 双边数据

      {
        "code": 200,
        "message": "",
        "data": {
          "result": [
            {
              "type": "user",
              "id": "7aeabd8c21cbdbca",
              "name": "lishi",
              "avatar": "",
              "nickname": "user7aeabd8c21cbdbca",
              "relation_status": "can_add"
            },
            {
              "type": "group",
              "id": "7ae7347c168bdbca",
              "name": "sh2group",
              "avatar": "files/89b8f0a444f1db552df33bc8746e2d06/veer-124990665.webp",
              "description": "fsdfds",
              "join_mode": 1,
              "relation_status": "can_join"
            }
          ],
          "next_cursor": "7aeabd8c21cbdbca_7ae7347c168bdbca"
        }
      }

## 单边数据

      {
        "code": 200,
        "message": "",
        "data": {
          "result": [
            {
              "type": "user",
              "id": "7aeabd8c21cbdbca",
              "name": "lishi",
              "avatar": "",
              "nickname": "user7aeabd8c21cbdbca",
              "relation_status": "can_add"
            }
          ],
          "next_cursor": "7aeabd8c21cbdbca_end"
        }
      }

## 没有数据了

      {
        "code": 200,
        "message": "",
        "data": {
          "result": null,
          "next_cursor": "end_end"
        }
      }

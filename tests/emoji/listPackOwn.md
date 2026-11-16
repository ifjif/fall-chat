# 找我拥有的表情包 (游标分页)

## 请求

      -X POST

      url = http://localhost:8080/api/v1/emoji/listPack/own

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkzNjIzOTYyLCJuYmYiOjE3OTM2MjIxNjIsImlhdCI6MTc5MzYyMjE2Mn0.QvsM9brQsHSndxEPV5PS7EJnliFwwR7xVxXUheZPbxw"

      -d @listPackOwn-req.json
      -o listPackOwn-res.json

- **json**

      {
        "next": "",
        "source": 1, // 1:系统默认 2: 自己上传 3:购买
        "status": 2 // 0:下架 1:草稿 2:上架 3:审核中
      }

## 成功，没有数据

      {
        "code": 200,
        "message": "",
        "data": {
          "data": [],
          "next": "eyJ2IjoxLCJwIjp7Imhhc19uZXh0IjpmYWxzZSwiaWQiOjExMDgzNjg4NzUzNDQ0MDQ0OH19",
          "has_next": false
        }
      }

## 成功，有数据

      {
        "code": 200,
        "message": "",
        "data": {
          "data": [
            {
              "id": "7a420e44164bdbca",
              "name": "张三emoji1",
              "cover_url": "files/bb80b3b088aa283edbd7fd4729c9d362/default.jpg",
              "owner_id": "7aeabdf9ad0bdbca",
              "soruce": "customer",
              "price_type": "free",
              "price": 0,
              "status": "draft",
              "created_at": "2026-11-02T12:26:13.239Z"
            }
          ],
          "next": "eyJ2IjoxLCJwIjp7Imhhc19uZXh0IjpmYWxzZSwiaWQiOjExMDgzNjg4NzUzNDQ0MDQ0OH19",
          "has_next": false
        }
      }

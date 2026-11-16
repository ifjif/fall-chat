# 获取 标签包中的所有 item

## 请求

      -X POST

      url = http://localhost:8080/api/v1/emoji/listItem

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzkzMTg5ODY0LCJuYmYiOjE3OTMxODgwNjQsImlhdCI6MTc5MzE4ODA2NH0.hr56rulCHY_vajiNkHzfv2UHW6u01sgphsUtLsE67hc"

      -d @listItem-req.json
      -o listItem-res.json

- **json**

      {
        "pack_id": "7a49da1daf4bdbca",
        "next": "",
        "sort_key": "sort_no_desc"   // sort_no_asc,sort_no_desc
      }

## 成功没有数据

    {
      "code": 200,
      "message": "",
      "data": {
        "data": null,
        "next": "eyJ2IjoxLCJwIjp7InNvcnRfbm8iOjAsImlkIjoxMDg4NzQ5ODEzMTk5MDUyODAsImhhc19uZXh0IjpmYWxzZX19",
        "has_next": false
      }
    }

## 错误1

      {
        "code": 2003,
        "message": "未知排序行为",
        "status": 400
      }

## 错误2 (可能是next解码失败)

      {
        "code": 9999,
        "message": "服务器内部错误",
        "status": 500
      }

## 成功带数据

      {
        "code": 200,
        "message": "",
        "data": {
          "data": [
            {
              "id": "7a4906fc314bdbc8",
              "pack_id": "7a49da1daf4bdbca",
              "media_type": "unicode",
              "media_url": "",
              "font_name": "1F602",
              "thumb_url": "",
              "width": 0,
              "height": 0,
              "file_size": 0,
              "sort_no": 2
            },
            {
              "id": "7a4906fc314bdbcb",
              "pack_id": "7a49da1daf4bdbca",
              "media_type": "unicode",
              "media_url": "",
              "font_name": "1F601",
              "thumb_url": "",
              "width": 0,
              "height": 0,
              "file_size": 0,
              "sort_no": 1
            },
            {
              "id": "7a4906fc314bdbca",
              "pack_id": "7a49da1daf4bdbca",
              "media_type": "unicode",
              "media_url": "",
              "font_name": "1F600",
              "thumb_url": "",
              "width": 0,
              "height": 0,
              "file_size": 0,
              "sort_no": 0
            }
          ],
          "next": "eyJ2IjoxLCJwIjp7InNvcnRfbm8iOjAsImlkIjoxMDg4NzQ5ODEzMTk5MDUyODAsImhhc19uZXh0IjpmYWxzZX19",
          "has_next": false
        }
      }

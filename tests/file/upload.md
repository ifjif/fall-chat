# 请求

-X POST

url = <http://localhost:8080/api/v1/file/upload>

-H "Content-Type: multipart/form-data"
-H "Accept: application/json"
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3MzI3MzM2LCJuYmYiOjE3ODczMjU1MzYsImlhdCI6MTc4NzMyNTUzNn0.oaZ7S-hxTSEAJBjWC7axipU_nr_0iIfJhpfpuBolDxM"

-F "file=@test-small5m.png"
-F "file_id=testfileid"

# 成功

      {
        "code": 200,
        "message": "上传成功",
        "data": "files/testfileid/test-small5m.png"
      }

# 文件大于 5MB

      {
        "code": 4401,
        "message": "文件太大请分片",
        "status": 400
      }

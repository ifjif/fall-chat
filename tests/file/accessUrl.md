# 请求 GET 需auth

-X GET

url = <http://localhost:8080/api/v1/file/accessUrl>

-H "Accept: application/json"
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3MzkyOTI1LCJuYmYiOjE3ODczOTExMjUsImlhdCI6MTc4NzM5MTEyNX0.WarAQktnl4W1CmIZZY--SoyuKxmU_Vpwy6yvHUKcUeU"

-G
-d "file=files/testchunksmerge/test-merge.zip"

# 成功

      {
        "code": 200,
        "message": "",
        "data": "/api/v1/file/fetch?file=files%2Ftestchunksmerge%2Ftest-merge.zip\u0026token=Bearer+eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiZmlsZV9hY2Nlc3MiLCJpc3MiOiJmYWxsLWNoYXQiLCJleHAiOjE3ODczOTIwOTYsIm5iZiI6MTc4NzM5MTc5NiwiaWF0IjoxNzg3MzkxNzk2fQ.y2PvHwvyGhrGzDJPeWWA-lM2kU8t9SirinHU2lXFfiQ"
      }

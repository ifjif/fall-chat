# 请求 GET

      -X GET

      url = http://localhost:8080/api/v1/ws/accessUrl

      -H "Content-Type: application/json"
      -H "Accept: application/json"
      -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3OTA5Nzk4LCJuYmYiOjE3ODc5MDc5OTgsImlhdCI6MTc4NzkwNzk5OH0.1FU25k7FX3ZMEqOwaAWWEVqqO6dBXs2WngDh31XqUPs"

# 成功

      {
        "code": 200,
        "message": "",
        "data": "ws://localhost:8080/api/v1/ws/7aeabdf9ad0bdbca?token=Bearer+eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2FlYWJkZjlhZDBiZGJjYSIsInVzZXJuYW1lIjoiemhhbmdzYW4iLCJkZXZpY2VfaWQiOiIiLCJqdGkiOiIiLCJ0eXBlIjoiYWNjZXNzIiwiaXNzIjoiZmFsbC1jaGF0IiwiZXhwIjoxNzg3OTEwMTM4LCJuYmYiOjE3ODc5MDgzMzgsImlhdCI6MTc4NzkwODMzOH0.Sx6uuNHb1tmjP6nCrVFshFZPE5lM4KbF-GeSg0MdfXI"
      }

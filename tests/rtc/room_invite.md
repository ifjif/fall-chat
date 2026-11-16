# rtc 邀请进入 房间

## 成功

    邀请方响应
      {
        "code": 200,
        "message": ""
      }

    被邀请方接受到ws消息
      {
        "type": "invite",
        "timestamp": 1794752517432,
        "data": {
          "room_id": "room_81476227113684992_459822249",
          "peer_id": "7aeabd8c21cbdbca",
          "inviter_id": "7aeabdf9ad0bdbca"
        }
      }

## 错误1

      {
        "code": 5404,
        "message": "请勿重复邀请",
        "status": 400
      }

## 错误2 (没有这个房间)

      {
        "code": 5403,
        "message": "不是房间创建者",
        "status": 400
      }

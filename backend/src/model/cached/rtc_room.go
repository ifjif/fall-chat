package cached

import (
	"time"
	"zzc/fall-chat/src/types"
)

type RtcRoomState string

func (rr RtcRoomState) MarshalBinary() ([]byte, error) {
	return []byte(rr), nil
}

const (
	RtcRoomStatePending RtcRoomState = "pending"
	RtcRoomStateWaiting RtcRoomState = "waiting"
	RtcRoomStateActive  RtcRoomState = "active"
	RtcRoomStateClosed  RtcRoomState = "closed"
)

type RtcRoomType string

func (rt RtcRoomType) MarshalBinary() (data []byte, err error) {
	return []byte(rt), nil
}

const (
	RtcRoomTypeOne2One RtcRoomType = "one_to_one"
	RtcRoomTypeGroup   RtcRoomType = "group"
)

type RtcRoom struct {
	RoomID      string          `redis:"room_id"`
	CreatorID   int64           `redis:"creator_id"`
	Type        RtcRoomType     `redis:"type"`
	MaxMember   int             `redis:"max_member"`
	CreatedTime types.RedisTime `redis:"created_time"`
	ExpiredTime types.RedisTime `redis:"expired_time"`
	State       RtcRoomState    `redis:"state"`
}

func (r *RtcRoom) IsWaiting() bool {
	return r.State == RtcRoomStateWaiting
}

func (r *RtcRoom) IsPending() bool {
	return r.State == RtcRoomStatePending
}

func (r *RtcRoom) IsExpire() bool {
	return r.ExpiredTime.Before(time.Now())
}

package model

import "time"

/*
*

	id bigint auto_increment primary key
	ip varchar(45),
	status tinyint default 1 comment '1:占用, 0: 释放',
	last_heartbeat datetime
	create_at datetime
	update_at datetime

*
*/
type MachineIdPool struct {
	ID            int64  `gorm:"primaryKey"`
	IP            string `gorm:"size:45"`
	Status        int8   `gorm:"default:1"`
	LastHeartbeat time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

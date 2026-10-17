package jobs

import (
	"time"
	"zzc/fall-chat/src/data"
)

type UpdateGrpcaddrExpireJob struct {
	CommonJob
	Duration time.Duration
	ID       string
	dao      *data.UpdateGrpcaddrExpireJobDao
}

func NewUpdateGrpcaddrExpireJob() *UpdateGrpcaddrExpireJob {
	return &UpdateGrpcaddrExpireJob{
		dao: data.GetUpdateGrpcaddrExpireJobDao(),
	}
}

func (u *UpdateGrpcaddrExpireJob) Key() string {
	return u.ID
}

func (u *UpdateGrpcaddrExpireJob) Delay() time.Duration {
	return u.Duration
}

func (u *UpdateGrpcaddrExpireJob) Execute() error {
	u.dao.UpdateGrpcAddrExpire2Redis(u.ID)
	return nil
}

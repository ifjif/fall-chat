package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/utils"
)

type Email struct {
	BaseService
	store store.EmailStore
}

func NewEmail(name string) *Email {
	return &Email{
		BaseService: BaseService{
			name: name,
		},
	}
}

func (es *Email) Inject() {
	es.store = data.GetBaseDao()
}

func (es *Email) PassResetCode(ctx context.Context, email string) (id string, err error) {
	// 随机6位 code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	duration := 15 * time.Minute

	prc := model.PasswordRestCode{}
	err_ := model.InjectID(&prc.Model)
	if err != nil {
		return id, es.Wrap("生成ID", err_, "email", email)
	}

	prc.Code = code
	prc.ExpireAt = time.Now().Add(duration)

	err = es.store.
		EmailRepos().
		Email.
		CreateOne(ctx, prc)
	if err != nil {
		return id, es.Wrap("插入重置密码验证码记录", err, "username", es.name)
	}

	go func() {
		utils.SendResetCodeEmail(email, code)
	}()

	id = prc.PublicID
	return id, nil
}

func (es *Email) VerifyPasswordRestCode(ctx context.Context, codeID, code string) error {
	id, err := model.DecodeID(codeID)
	if err != nil {
		return es.Wrap("codeID解码", err, "code_id", codeID)
	}

	prc, err := es.store.
		EmailRepos().
		Email.
		FindByID(ctx, id)
	if err != nil {
		return es.Wrap("获取重置密码验证码", err, "code_id", codeID)
	}

	if prc.Used || prc.Code != code || time.Now().After(prc.ExpireAt) {
		return es.Wrap("验证码验证", errs.EmailErrInvalid, "code_id", codeID)
	}

	// 设置为使用
	err = es.store.
		EmailRepos().Email.MarkUsed(ctx, prc)
	if err != nil {
		return es.Wrap("更新使用状态", err, "code_id", codeID)
	}

	return nil
}

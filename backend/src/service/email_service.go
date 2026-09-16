package service

import (
	"fmt"
	"math/rand"
	"time"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/utils"
)

type EmailService struct {
	BaseService
}

func NewEmailService(name string) *EmailService {
	return &EmailService{
		BaseService: BaseService{
			name: name,
		},
	}
}

func (es *EmailService) PassResetCode(email string) (id string, err error) {
	// 随机6位 code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	duration := 15 * time.Minute

	id_, err_ := es.sf.NextID()
	if err_ != nil {
		return id, es.Wrap("生成ID", err_, "email", email)
	}

	pID := utils.SF.Encode(id_)

	prc := model.NewPasswordRestCode(id_)
	prc.Code = code
	prc.PublicID = pID
	prc.ExpireAt = time.Now().Add(duration)

	result := es.db.Create(prc)
	if result.Error != nil {
		return id, es.Wrap("插入重置密码验证码记录", result.Error, "username", es.name)
	}

	go func() {
		utils.SendResetCodeEmail(email, code)
	}()

	id = pID
	return id, nil
}

func (es *EmailService) VerifyPasswordRestCode(codeID, code string) error {
	id, err := es.sf.Decode(codeID)
	if err != nil {
		return es.Wrap("codeID解码", err, "code_id", codeID)
	}

	var prc model.PasswordRestCode
	result := es.db.Where("id = ?", id).Take(&prc)

	if result.Error != nil {
		return es.Wrap("获取重置密码验证码", err, "code_id", codeID)
	}

	if prc.Used || prc.Code != code || time.Now().After(prc.ExpireAt) {
		return es.Wrap("验证码验证", errs.EmailErrInvalid, "code_id", codeID)
	}

	// 设置为使用
	result = es.db.Model(&prc).Update("used", true)
	if result.Error != nil {
		return es.Wrap("更新使用状态", result.Error, "code_id", codeID)
	}

	return nil
}

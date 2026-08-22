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
		es.LogWarn("PassResetCode", err_,
			"email", email,
			"reason", "生成ID失败",
		)
		err = err_
		return id, err
	}

	pID := utils.SF.Encode(id_)

	prc := model.NewPasswordRestCode(id_)
	prc.Code = code
	prc.PublicID = pID
	prc.ExpireAt = time.Now().Add(duration)

	result := es.db.Create(prc)
	if result.Error != nil {
		es.LogWarn("PassResetCode", result.Error,
			"username", es.name,
			"reason", "插入重置密码验证码记录失败",
		)
		err = result.Error
		return id, err
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
		es.LogWarn("VerifyPasswordRestCode", err,
			"code_id", codeID,
			"reason", "codeID解码失败",
		)
		return err
	}

	var prc model.PasswordRestCode
	result := es.db.Where("id = ?", id).Take(&prc)

	if result.Error != nil {
		es.LogWarn("VerifyPasswordRestCode", err,
			"code_id", codeID,
			"reason", "获取重置密码验证码失败",
		)
		return result.Error
	}

	if prc.Used || prc.Code != code || time.Now().After(prc.ExpireAt) {
		es.LogWarn("VerifyPasswordRestCode", errs.EmailErrInvalid,
			"code_id", codeID,
		)
		return errs.EmailErrInvalid
	}

	// 设置为使用
	result = es.db.Model(&prc).Update("used", true)
	if result.Error != nil {
		es.LogWarn("VerifyPasswordRestCode", result.Error,
			"code_id", codeID,
			"reason", "更新使用状态失败",
		)

		return result.Error
	}

	return nil
}

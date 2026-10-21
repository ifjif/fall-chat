package service

import (
	"context"
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/utils"
)

type User struct {
	BaseService
	emailService *Email
	store        store.UserStore
}

func NewUser(name string) *User {
	return &User{
		BaseService: BaseService{
			name: name,
		},
	}
}

func (us *User) Inject() {
	us.BaseService.Inject()
	es := GlobalService[EMAIL_SERVICE_NAME].(*Email)
	us.emailService = es
	us.store = data.GetBaseDao()
}

func (us *User) RegisterUser(ctx context.Context, username, password string) (string, error) {
	var pID string

	count, err := us.store.
		UserRepos().
		User.
		CountByUsername(ctx, username)
	if err != nil {
		return pID, us.Wrap("RegisterUser", err)
	}

	if count > 0 {
		return pID, us.Wrap("RegisterUser", errs.UserErrExistUser,
			"username", username,
			"reason", "用户已存在",
		)
	}

	id, err := us.sf.NextID()
	if err != nil {
		return pID, us.Wrap("RegisterUser", err,
			"username", username,
			"reason", "生成ID失败",
		)
	}

	pID = us.sf.Encode(id)

	user := model.NewUser(id)
	user.PublicID = pID
	user.Username = username
	user.Nickname = "user" + pID
	user.Password = utils.EncryptPassword(password)

	err = us.store.
		UserRepos().
		User.
		CreateOne(ctx, *user)
	if err != nil {
		return "", us.Wrap("RegisterUser", err)
	}

	return pID, nil
}

func (us *User) Login(ctx context.Context, username, password string) (*response.LoginResult, error) {
	user, err := us.store.
		UserRepos().
		User.
		FindByUsername(ctx, username)
	if err != nil {
		return nil, us.Wrap("Login", err)
	}
	err = utils.DecryptPassword(user.Password, password)
	if err != nil {
		return nil, us.Wrap("", errs.UserErrUsernameOrPassword.Meger(err), "username", username, "reason", "密码不匹配")
	}

	if *user.Status == 0 {
		return nil, us.Wrap("", errs.UserErrAccountDisabled, "username", username)
	}

	// arg1 jti存redis
	_, refreshToken, refreshExpiresIn, err := utils.GenerateRefreshToken(user.PublicID, user.Username, "")
	if err != nil {
		return nil, us.Wrap("JWT refresh token 生成", err,
			"username", user.Username)
	}

	accessToken, accessExpiresIn, err := utils.GenerateAccessToken(user.PublicID, user.Username, "")
	if err != nil {
		return nil, us.Wrap("JWT access token 生成", err,
			"username", user.Username)
	}

	result := &response.LoginResult{
		User:             user,
		RefreshToken:     refreshToken,
		RefreshExpiresIn: refreshExpiresIn,
		AccessToken:      accessToken,
		AccessExpiresIn:  accessExpiresIn,
	}

	return result, nil
}

func (us *User) ResetPassword(ctx context.Context, cmd *command.RestPassword) error {
	err := us.emailService.VerifyPasswordRestCode(ctx, cmd.CodeID, cmd.Code)
	if err != nil {
		return us.Wrap("verify code faild", err,
			"username", cmd.Username, "code_id", cmd.CodeID)
	}

	user, err := us.store.
		UserRepos().
		User.
		FindByUsername(ctx, cmd.Username)
	if err != nil {
		return err
	}

	err = utils.DecryptPassword(user.Password, cmd.OldPass)
	if err != nil {
		return us.Wrap("密码匹配", err, "username", cmd.Username)
	}

	if cmd.Email != user.Email {
		return us.Wrap("", errs.UserErrEmail,
			"username", cmd.Username,
		)
	}

	if *user.Status == 0 {
		return us.Wrap("", errs.UserErrAccountDisabled,
			"username", cmd.Username,
		)
	}

	newPass := utils.EncryptPassword(cmd.NewPass)

	err = us.store.
		UserRepos().
		User.UpdatePassword(ctx, int64(user.ID), newPass)
	if err != nil {
		return us.Wrap("update password", err,
			"username", us.name,
			"reason", "更新密码失败",
		)
	}

	return nil
}

func (us *User) ForgetPassword(ctx context.Context, cmd *command.ForgetPassword) error {
	err := us.emailService.VerifyPasswordRestCode(ctx, cmd.CodeID, cmd.Code)
	if err != nil {
		return us.Wrap("verify code", err,
			"username", cmd.Username, "email", cmd.Email)
	}

	user, err := us.store.
		UserRepos().
		User.
		FindByUsername(ctx, cmd.Username)
	if err != nil {
		return us.Wrap("找用户", err, "username", cmd.Username)
	}

	if user.Email != cmd.Email {
		return us.Wrap("", errs.UserErrEmail,
			"username", cmd.Username,
		)
	}

	newPass := utils.EncryptPassword(cmd.NewPass)

	err = us.store.
		UserRepos().
		User.UpdatePassword(ctx, int64(user.ID), newPass)
	if err != nil {
		return us.Wrap("更新密码", err,
			"username", cmd.Username)
	}

	return nil
}

func (us *User) UpdateProfile(ctx context.Context, cmd *command.UpdateProfile) error {
	user, err := us.store.
		UserRepos().
		User.FindByID(ctx, cmd.ID)
	if err != nil {
		return us.Wrap("找用户", err, "id", cmd.ID)
	}

	if err := utils.DecryptPassword(user.Password, cmd.Password); err != nil {
		return us.Wrap("", errs.UserErrUsernameOrPassword.Meger(err),
			"username", user.Username)
	}

	if *user.Status == 0 {
		return us.Wrap("", errs.UserErrAccountDisabled, "username", user.Username)
	}

	searchFlag := int8(cmd.SearchFlag)
	updateUser := model.User{
		Email:      cmd.Email,
		Phone:      cmd.Phone,
		Nickname:   cmd.Nickname,
		Avatar:     cmd.Avatar,
		SearchFlag: &searchFlag,
	}
	updateUser.ID = user.ID

	err = us.store.
		UserRepos().
		User.
		UpdateOne(ctx, updateUser)
	if err != nil {
		return us.Wrap("更新profile", err,
			"username", user.Username)
	}

	return nil
}

func (us *User) Profile(ctx context.Context, id int64) (*model.User, error) {
	user, err := us.store.
		UserRepos().
		User.FindByID(ctx, id)
	if err != nil {
		return nil, us.Wrap("找用户", err, "id", id)
	}

	return user, nil
}

package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
)

func ToResetPasswordCommand(req *request.ResetPassword) *command.RestPassword {
	return &command.RestPassword{
		ID:       req.ID,
		Username: req.Username,
		Email:    req.Email,
		OldPass:  req.OldPass,
		NewPass:  req.NewPass,
		Code:     req.Code,
		CodeID:   req.CodeID,
	}
}

func ToUpdateProfileCommand(req *request.UpdateProfile) *command.UpdateProfile {
	return &command.UpdateProfile{
		ID:         req.ID,
		Username:   req.Username,
		Email:      req.Email,
		Nickname:   req.Nickname,
		Password:   req.Password,
		SearchFlag: req.SearchFlag,
		Avatar:     req.Avatar,
		Phone:      req.Phone,
	}
}

func ToLoginUserResponse(result *response.LoginResult) *response.LoginUser {
	return &response.LoginUser{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.AccessExpiresIn,
		User:        *ToUserResponse(result.User),
	}
}

func ToUserResponse(model *model.User) *response.User {
	return &response.User{
		ID:        model.PublicID,
		Username:  model.Username,
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
	}
}

func ToForgetPasswordCommand(req *request.ForgetPassword) *command.ForgetPassword {
	return &command.ForgetPassword{
		Username: req.Username,
		Email:    req.Email,
		NewPass:  req.NewPass,
		Code:     req.Code,
		CodeID:   req.CodeID,
	}
}

func ToUserProfileResponse(m *model.User) *response.UserProfile {
	return &response.UserProfile{
		ID:         m.PublicID,
		Avatar:     m.Avatar,
		Username:   m.Username,
		Nickname:   m.Nickname,
		Email:      m.Email,
		Phone:      m.Phone,
		Status:     int(*m.Status),
		SearchFlag: int(*m.SearchFlag),
		CreatedAt:  m.CreatedAt,
	}
}

package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToResetPasswordCommand(req *request.ResetPasswordRequest) *command.RestPasswordCommand {
	return &command.RestPasswordCommand{
		ID:       req.ID,
		Username: req.Username,
		Email:    req.Email,
		OldPass:  req.OldPass,
		NewPass:  req.NewPass,
		Code:     req.Code,
		CodeID:   req.CodeID,
	}
}

func ToUpdateProfileCommand(req *request.UpdateProfileRequest) *command.UpdateProfileCommand {
	return &command.UpdateProfileCommand{
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

func ToLoginUserResponse(result *response.LoginResult) *response.LoginUserResponse {
	return &response.LoginUserResponse{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.AccessExpiresIn,
		User:        *ToUserResponse(result.User),
	}
}

func ToUserResponse(model *model.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        model.PublicID,
		Username:  model.Username,
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
	}
}

func ToForgetPasswordCommand(req *request.ForgetPasswordRequest) *command.ForgetPasswordCommand {
	return &command.ForgetPasswordCommand{
		Username: req.Username,
		Email:    req.Email,
		NewPass:  req.NewPass,
		Code:     req.Code,
		CodeID:   req.CodeID,
	}
}

func ToSearchQuery(req *request.SearchRequest) *query.SearchQuery {
	return &query.SearchQuery{
		ID:         req.ID,
		Username:   req.Username,
		Name:       req.Name,
		NextCursor: req.NextCuror,
	}
}

func ToUserProfileResponse(m *model.User) *response.UserProfileResponse {
	return &response.UserProfileResponse{
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

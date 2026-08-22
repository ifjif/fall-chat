package service

import (
	"errors"
	"strings"
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/utils"

	"gorm.io/gorm"
)

type UserService struct {
	BaseService
	emailService *EmailService
}

func NewUserService(name string) *UserService {
	return &UserService{
		BaseService: BaseService{
			name: name,
		},
	}
}

func (us *UserService) Inject() {
	us.BaseService.Inject()
	es := GlobalService[EMAIL_SERVICE_NAME].(*EmailService)
	us.emailService = es
}

func (us *UserService) FindByID(id int64) (*model.User, error) {
	user := &model.User{}
	user.ID = uint64(id)

	result := us.db.Take(user)
	if result.Error != nil {
		us.LogWarn("FindByID", result.Error, "id", id)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.UserErrNotExistUser
		}

		return nil, result.Error
	}

	return user, nil
}

func (us *UserService) FindByUsername(username string) (*model.User, error) {
	user := &model.User{}
	result := us.db.Where("username = ?", username).Take(user)
	if result.Error != nil {
		us.LogWarn("FindByUsername", result.Error, "username", username)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.UserErrUsernameOrPassword
		}

		return nil, result.Error
	}

	return user, nil
}

func (us *UserService) CountByUsername(username string) (int64, error) {
	var count int64
	result := us.db.Model(&model.User{}).Where("username = ?", username).Count(&count)
	if result.Error != nil {
		us.LogWarn("CountByUsername", result.Error, "username", username)
		return 0, result.Error
	}

	return count, nil
}

func (us *UserService) createUser(user *model.User) error {
	result := us.db.Create(user)
	if result.Error != nil {
		us.LogWarn("createUser", result.Error)
		return result.Error
	}

	return nil
}

func (us *UserService) RegisterUser(username, password string) (string, error) {
	var pID string

	count, err := us.CountByUsername(username)
	if err != nil {
		return pID, err
	}

	if count > 0 {
		us.LogWarn("RegisterUser", errs.UserErrExistUser,
			"username", username,
			"reason", "用户已存在",
		)
		return pID, errs.UserErrExistUser
	}

	id, err := us.sf.NextID()
	if err != nil {
		us.LogWarn("RegisterUser", err,
			"username", username,
			"reason", "生成ID失败",
		)
		return pID, err
	}

	pID = us.sf.Encode(id)

	user := model.NewUser(id)
	user.PublicID = pID
	user.Username = username
	user.Nickname = "user" + pID
	user.Password = utils.EncryptPassword(password)

	err = us.createUser(user)

	return pID, err
}

func (us *UserService) Login(username, password string) (*response.LoginResult, error) {
	user, err := us.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	err = utils.DecryptPassword(user.Password, password)
	if err != nil {
		us.LogWarn("Login", err, "username", username, "reason", "密码不匹配")
		return nil, errs.UserErrUsernameOrPassword
	}

	if *user.Status == 0 {
		us.LogWarn("Login", errs.UserErrAccountDisabled, "username", username)
		return nil, errs.UserErrAccountDisabled
	}

	// arg1 jti存redis
	_, refreshToken, refreshExpiresIn, err := utils.GenerateRefreshToken(user.PublicID, user.Username, "")
	if err != nil {
		us.LogWarn("Login", err,
			"username", user.Username,
			"reason", "JWT refresh token 生成失败",
		)
		return nil, errs.CommonErrJWTTokenGeneration
	}

	accessToken, accessExpiresIn, err := utils.GenerateAccessToken(user.PublicID, user.Username, "")
	if err != nil {
		us.LogWarn("Login", err,
			"username", user.Username,
			"reason", "JWT access token 生成失败",
		)
		return nil, errs.CommonErrJWTTokenGeneration
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

func (us *UserService) ResetPassword(cmd *command.RestPasswordCommand) error {
	err := us.emailService.VerifyPasswordRestCode(cmd.CodeID, cmd.Code)
	if err != nil {
		us.LogError("ResetPassword", err,
			"username", cmd.Username,
			"code_id", cmd.CodeID,
			"reason", "verify_code_faild",
		)
		return err
	}

	user, err := us.FindByUsername(cmd.Username)
	if err != nil {
		return err
	}

	err = utils.DecryptPassword(user.Password, cmd.OldPass)
	if err != nil {
		us.LogWarn("ResetPassword", err,
			"username", cmd.Username,
			"reason", "密码不匹配",
		)
		return errs.UserErrUsernameOrPassword
	}

	if cmd.Email != user.Email {
		us.LogWarn("ResetPassword", errs.UserErrEmail,
			"username", cmd.Username,
		)
		return errs.UserErrEmail
	}

	if *user.Status == 0 {
		us.LogWarn("ResetPassword", errs.UserErrAccountDisabled,
			"username", cmd.Username,
		)
		return errs.UserErrAccountDisabled
	}

	newPass := utils.EncryptPassword(cmd.NewPass)

	result := us.db.Model(user).Update("password", newPass)
	if result.Error != nil {
		us.LogWarn("ResetPassword", result.Error,
			"username", us.name,
			"reason", "更新密码失败",
		)
	}
	return result.Error
}

func (us *UserService) ForgetPassword(cmd *command.ForgetPasswordCommand) error {
	err := us.emailService.VerifyPasswordRestCode(cmd.CodeID, cmd.Code)
	if err != nil {
		us.LogError("ForgetPassword", err,
			"username", cmd.Username,
			"email", cmd.Email,
			"reason", "verify_code_faild",
		)

		return err
	}

	user, err := us.FindByUsername(cmd.Username)
	if err != nil {
		return err
	}

	if user.Email != cmd.Email {
		us.LogWarn("ForgetPassword", errs.UserErrEmail,
			"username", cmd.Username,
		)

		return errs.UserErrEmail
	}

	newPass := utils.EncryptPassword(cmd.NewPass)

	result := us.db.Model(user).Update("password", newPass)
	if result.Error != nil {
		us.LogWarn("ForgetPassword", result.Error,
			"username", cmd.Username,
			"reason", "更新密码失败",
		)

		return result.Error
	}

	return nil
}

func (us *UserService) UpdateProfile(cmd *command.UpdateProfileCommand) error {
	user, err := us.FindByID(cmd.ID)
	if err != nil {
		return err
	}

	if err := utils.DecryptPassword(user.Password, cmd.Password); err != nil {
		us.LogWarn("UpdateProfile", err,
			"username", user.Username,
			"reason", "密码不匹配",
		)
		return errs.UserErrUsernameOrPassword
	}

	if *user.Status == 0 {
		us.LogWarn("UpdateProfile", errs.UserErrAccountDisabled, "username", user.Username)
		return errs.UserErrAccountDisabled
	}

	searchFlag := int8(cmd.SearchFlag)
	updateUser := model.User{
		Email:      cmd.Email,
		Phone:      cmd.Phone,
		Nickname:   cmd.Nickname,
		Avatar:     cmd.Avatar,
		SearchFlag: &searchFlag,
	}

	result := us.db.Model(user).
		Updates(updateUser)

	if result.Error != nil {
		us.LogWarn("UpdateProfile", result.Error,
			"username", user.Username,
			"reason", "更新profile失败",
		)

		return result.Error
	}

	return nil
}

// 模糊匹配
//
// user -> 10
// group - 10
// user_cursor = last_user_id    infinity end
// group_cursor = last_group_id  infinity end
//
//	already_friend	显示"发消息"	friendships 中 status=1
//	request_pending	显示"等待验证"	friend_requests 中有待处理申请
//	can_add	显示"添加好友"	无任何关系记录
//	already_joined	显示"进入群聊"	group_members 中有记录
//	need_approval	显示"申请加入"	join_mode=1 且未加入
//	cannot_join	显示"不可加入"	join_mode=2
//	can_join 显示加入群聊 join_mode=1
func (us *UserService) Search(qu *query.SearchQuery) (*response.SearchResult, error) {
	// 先找人
	id := qu.ID
	searchName := qu.Name
	keyword := "%" + searchName + "%"
	// userID_groupID
	nextCursor := qu.NextCursor
	cursors := strings.Split(nextCursor, "_")
	userCursor := "infinity"
	groupCursor := "infinity"
	last_user_id := int64(0)
	last_group_id := int64(0)
	if len(cursors) >= 2 {
		userCursor = cursors[0]
		groupCursor = cursors[1]
	}

	if userCursor != "infinity" && userCursor != "end" {
		last_user_id, _ = us.sf.Decode(userCursor)
	}
	if groupCursor != "infinity" && groupCursor != "end" {
		last_group_id, _ = us.sf.Decode(groupCursor)
	}

	finalResult := &response.SearchResult{}
	var results []response.SearchDataResult

	// 找人
	var userResults []response.SearchDataResult
	if userCursor != "end" {
		result := us.db.Raw(`
		select u.public_id as id,
					 u.username as name,
					 u.nickname as nickname,
					 u.avatar as avatar,
					 '' as description,
					 case 
							when fs.status = 1 then 'already_friend'
							when fsq.status = 0 then 'request_pending'
							else 'can_add'
						end as relation_status 
		from users u
		left join friendships fs
				on u.id = fs.friend_id and fs.status = 1 and fs.user_id = ?
		left join friendship_requests fsq
				on u.id = fsq.to_user_id and fsq.status = 0 and fsq.from_user_id = ?
		where u.search_flag = 1
			and u.id != ?
			and (u.username like ? or u.nickname like ?)
			and u.id > ?
		order by u.id asc
		limit 10
		`, id, id, id, keyword, keyword, last_user_id).
			Scan(&userResults)

		if result.Error != nil {
			us.LogWarn("search", result.Error,
				"username", qu.Username,
				"searchName", searchName,
				"reason", "模糊找用户失败",
			)

			return nil, result.Error
		}
	}
	for i := range userResults {
		userResults[i].Type = "user"
	}
	results = append(results, userResults...)
	nextUserID := "end"
	if len(userResults) != 0 {
		newID := userResults[len(userResults)-1].ID
		if newID != userCursor {
			nextUserID = newID
		}
	}

	// 找群
	var groupResults []response.SearchDataResult

	if groupCursor != "end" {
		result := us.db.Raw(`
		select  g.public_id as id,
						g.group_name as name,
						g.avatar as avatar,
						g.description as description,
						g.join_mode as join_mode,
						'' as nickname,
						case 
							when gm.user_id is not null then 'already_joined'
							when g.join_mode = 2 then 'cannot_join'
							when g.join_mode = 1 then 'need_approval'
							else 'can_join'
						end as relation_status
		from group_infos g
			left join group_members gm on
			g.id = gm.group_id and gm.user_id = ?
		where g.search_flag = 1 
			and g.group_name like ?
			and g.id > ?
		order by g.id asc
		limit 10
		`, id, keyword, last_group_id).
			Scan(&groupResults)

		if result.Error != nil {
			us.LogWarn("search", result.Error,
				"username", qu.Username,
				"searchName", searchName,
				"reason", "模糊找群失败",
			)
			return nil, result.Error
		}
	}
	for i := range groupResults {
		groupResults[i].Type = "group"
	}

	results = append(results, groupResults...)
	nextGroupID := "end"
	if len(groupResults) != 0 {
		newID := groupResults[len(groupResults)-1].ID
		if newID != groupCursor {
			nextGroupID = newID
		}
	}

	cursor := nextUserID + "_" + nextGroupID

	finalResult.Result = results
	finalResult.NextCursor = cursor

	return finalResult, nil
}

func (us *UserService) Profile(id int64) (*model.User, error) {
	user, err := us.FindByID(id)
	if err != nil {
		return nil, err
	}

	return user, err
}

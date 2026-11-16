package errs

import "github.com/gofiber/fiber/v3"

var (
	SocialErrBlacklisted                   = NewBizError(SocialStartCode, fiber.StatusBadRequest, "对方已把你拉黑")
	SocialErrBlacklist                     = NewBizError(SocialStartCode+1, fiber.StatusBadRequest, "你拉黑了对方")
	SocialErrFriendRequestExists           = NewBizError(SocialStartCode+2, fiber.StatusBadRequest, "好友申请已存在")
	SocialErrFriendRequestNotExists        = NewBizError(SocialStartCode+3, fiber.StatusBadRequest, "好友申请不存在")
	SocialErrFriendRequestHandledOrExpired = NewBizError(SocialStartCode+4, fiber.StatusBadRequest, "好友申请已处理或过期")
	SocialErrGroupRequestApplyType         = NewBizError(SocialStartCode+5, fiber.StatusBadRequest, "错误的群申请类型")
	SocialErrGroupRequestInviteSelf        = NewBizError(SocialStartCode+6, fiber.StatusBadRequest, "群申请无法邀请自己")
	SocialErrGroupRequestExists            = NewBizError(SocialStartCode+7, fiber.StatusBadRequest, "群申请已存在")
	SocialErrGroupNotExists                = NewBizError(SocialStartCode+8, fiber.StatusBadRequest, "申请的群不存在")
	SocialErrGroupDisabled                 = NewBizError(SocialStartCode+9, fiber.StatusBadRequest, "申请的群已禁用")
	SocialErrGroupNotJoin                  = NewBizError(SocialStartCode+10, fiber.StatusBadRequest, "申请的群不可加入")
	SocialErrGroupMemberExists             = NewBizError(SocialStartCode+11, fiber.StatusBadRequest, "无法申请已加入的群")
	SocialErrGroupRequestNotExists         = NewBizError(SocialStartCode+12, fiber.StatusBadRequest, "群申请不存在")
	SocialErrGroupRequestHandledOrExpired  = NewBizError(SocialStartCode+13, fiber.StatusBadRequest, "群申请已处理或过期")
	SocialErrGroupRequestHandleSelf        = NewBizError(SocialStartCode+14, fiber.StatusBadRequest, "无法处理自己的群申请")
	SocialErrGroupMemberNotExists          = NewBizError(SocialStartCode+15, fiber.StatusBadRequest, "未加入群，无法处理")
	SocialErrGroupMemberRoleOrdinary       = NewBizError(SocialStartCode+16, fiber.StatusBadRequest, "不是管理员，无法处理")
	SocialErrGroupMemberOverflow           = NewBizError(SocialStartCode+17, fiber.StatusBadRequest, "群成员已满")
	SocialErrGroupMemberInvitorNotExists   = NewBizError(SocialStartCode+18, fiber.StatusBadRequest, "邀请人未加入群，无法处理")
	SocialErrGroupConvNotExists            = NewBizError(SocialStartCode+19, fiber.StatusBadRequest, "群会话不存在")
	SocialErrGroupConvMembNotExists        = NewBizError(SocialStartCode+20, fiber.StatusBadRequest, "群会话成员不存在")
)

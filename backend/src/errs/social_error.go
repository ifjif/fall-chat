package errs

import "github.com/gofiber/fiber/v3"

var (
	SocialErrBlacklisted                   = NewBizError(4800, fiber.StatusBadRequest, "对方已把你拉黑")
	SocialErrBlacklist                     = NewBizError(4801, fiber.StatusBadRequest, "你拉黑了对方")
	SocialErrFriendRequestExists           = NewBizError(4802, fiber.StatusBadRequest, "好友申请已存在")
	SocialErrFriendRequestNotExists        = NewBizError(4803, fiber.StatusBadRequest, "好友申请不存在")
	SocialErrFriendRequestHandledOrExpired = NewBizError(4804, fiber.StatusBadRequest, "好友申请已处理或过期")
	SocialErrGroupRequestApplyType         = NewBizError(4805, fiber.StatusBadRequest, "错误的群申请类型")
	SocialErrGroupRequestInviteSelf        = NewBizError(4806, fiber.StatusBadRequest, "群申请无法邀请自己")
	SocialErrGroupRequestExists            = NewBizError(4807, fiber.StatusBadRequest, "群申请已存在")
	SocialErrGroupNotExists                = NewBizError(4808, fiber.StatusBadRequest, "申请的群不存在")
	SocialErrGroupDisabled                 = NewBizError(4809, fiber.StatusBadRequest, "申请的群已禁用")
	SocialErrGroupNotJoin                  = NewBizError(4810, fiber.StatusBadRequest, "申请的群不可加入")
	SocialErrGroupMemberExists             = NewBizError(4811, fiber.StatusBadRequest, "无法申请已加入的群")
	SocialErrGroupRequestNotExists         = NewBizError(4812, fiber.StatusBadRequest, "群申请不存在")
	SocialErrGroupRequestHandledOrExpired  = NewBizError(4813, fiber.StatusBadRequest, "群申请已处理或过期")
	SocialErrGroupRequestHandleSelf        = NewBizError(4814, fiber.StatusBadRequest, "无法处理自己的群申请")
	SocialErrGroupMemberNotExists          = NewBizError(4815, fiber.StatusBadRequest, "未加入群，无法处理")
	SocialErrGroupMemberRoleOrdinary       = NewBizError(4816, fiber.StatusBadRequest, "不是管理员，无法处理")
	SocialErrGroupMemberOverflow           = NewBizError(4817, fiber.StatusBadRequest, "群成员已满")
	SocialErrGroupMemberInvitorNotExists   = NewBizError(4818, fiber.StatusBadRequest, "邀请人未加入群，无法处理")
)

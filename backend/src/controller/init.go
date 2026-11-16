package controller

import (
	"github.com/gofiber/fiber/v3"
)

type (
	InitOption func(*initConfig)
	initConfig struct {
		beforeInject   func()
		beforeSetRoute func() fiber.Router
	}
)

func WithBeforeSetRoute(hook func() fiber.Router) InitOption {
	return func(c *initConfig) {
		c.beforeSetRoute = hook
	}
}

func WithBeforeInjectHook(hook func()) InitOption {
	return func(c *initConfig) {
		c.beforeInject = hook
	}
}

var GlobalController = map[string]Controller{}

func Init(app *fiber.App, opts ...InitOption) {
	cfg := &initConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	doRegister()

	group := app.Group("")
	if cfg.beforeSetRoute != nil {
		group = cfg.beforeSetRoute()
	}

	doSetRoute(group)
	if cfg.beforeInject != nil {
		cfg.beforeInject()
	}
	doInject()
}

func register(controller Controller) {
	GlobalController[controller.Name()] = controller
}

func doSetRoute(r fiber.Router) {
	for _, c := range GlobalController {
		c.SetRoute(r)
	}
}

func doInject() {
	for _, c := range GlobalController {
		c.Inject()
	}
}

func doRegister() {
	register(NewUser("UserController"))
	register(NewEmail("EmailController"))
	register(NewAuth("AuthController"))
	register(NewGroupInfo("GroupInfoController"))
	register(NewFile("FileController"))
	register(NewWS("WSController"))
	register(NewConversation("ConversationController"))
	register(NewConversationMember("ConversationMemberController"))
	register(NewSearch("SearchController"))
	register(NewFriendship("FriendshipController"))
	register(NewFriendshipRequest("FriendshipRequestController"))
	register(NewGroupRequest("RequestController"))
	register(NewIMMessage("IMMessageController"))
	register(NewEmoji("EmojiController"))
	register(NewRTCController("RTCController"))
}

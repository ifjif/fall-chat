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

func doRegister() {
	register(NewUserController("UserController"))
	register(NewEmailController("EmailController"))
	register(NewAuthController("AuthController"))
	register(NewGroupInfoController("GroupInfoController"))
	register(NewFileController("FileController"))
	// register(NewSocialController("SocialController"))
	register(NewWSController("WSController"))
	register(NewConversationController("ConversationController"))
	register(NewConversationMemberController("ConversationMemberController"))
	register(NewSearchController("SearchController"))
	register(NewFriendshipController("FriendshipController"))
	register(NewFriendshipRequestController("FriendshipRequestController"))
	register(NewGroupRequestController("RequestController"))
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

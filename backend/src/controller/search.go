package controller

import (
	"zzc/fall-chat/src/bind"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/middleware"
	"zzc/fall-chat/src/resp"
	"zzc/fall-chat/src/service"

	"github.com/gofiber/fiber/v3"
)

type Search struct {
	Base
	searchService *service.Search
}

func NewSearch(name string) *Search {
	bs := Base{name: name}
	return &Search{
		Base: bs,
	}
}

func (sc *Search) Inject() {
	s := service.GlobalService[service.SEARCH_SERVICE_NAME].(*service.Search)
	sc.searchService = s
}

func (sc *Search) SetRoute(r fiber.Router) {
	group := r.Group("/search")
	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Post("/userOrGroup", sc.userOrGroup)
}

// /search/userOrGroup 需auth
func (sc *Search) userOrGroup(c fiber.Ctx) error {
	userID := ctxkey.GetUserID(c)
	username := ctxkey.GetUsername(c)
	req, err := bind.Body[request.UserOrGroup](c)
	if err != nil {
		return sc.Wrap("Body绑定", err,
			"username", username,
		)
	}

	req.Username = username
	req.ID = userID

	qu := converter.ToUserOrGroupQuery(req)

	results, err := sc.searchService.UserOrGroup(c.Context(), qu)
	if err != nil {
		return err
	}

	return resp.
		New(c).
		OK().
		JSON(results)
}

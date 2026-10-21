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

type SearchController struct {
	BaseController
	searchService *service.Search
}

func NewSearchController(name string) *SearchController {
	bs := BaseController{name: name}
	return &SearchController{
		BaseController: bs,
	}
}

func (sc *SearchController) Inject() {
	s := service.GlobalService[service.SEARCH_SERVICE_NAME].(*service.Search)
	sc.searchService = s
}

func (sc *SearchController) SetRoute(r fiber.Router) {
	group := r.Group("/search")
	authGroup := group.Group("")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.Post("/userOrGroup", sc.userOrGroup)
}

// /search/userOrGroup 需auth
func (sc *SearchController) userOrGroup(c fiber.Ctx) error {
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

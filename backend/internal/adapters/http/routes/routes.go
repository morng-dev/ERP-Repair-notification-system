package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/morng-dev/erp/internal/adapters/http/handler"
	"github.com/morng-dev/erp/internal/adapters/http/middleware"
)

type Routes struct {
	authMW        *middleware.AuthMiddleware
	authHandler   *handler.AuthHandler
	profesHandler *handler.ProfessionHandler
}

func NewRoutes(
	authMW *middleware.AuthMiddleware,
	authHandler *handler.AuthHandler,
	profesHandler *handler.ProfessionHandler,

) *Routes {
	return &Routes{
		authHandler:   authHandler,
		authMW:        authMW,
		profesHandler: profesHandler,
	}
}

const (
	PermissionCreateProfession = "create_profession"
	PermissionReadProfession   = "read_profession"
	PermissionUpdateProfession = "update_profession"
	PermissionDeleteProfession = "delete_profession"
)

func (r *Routes) SetUpRoute(app *fiber.App) {

	api := app.Group("/api/v1")
	//auth
	auth := api.Group("/auth")
	auth.Post("/register", r.authHandler.Register)
	auth.Post("/login", r.authHandler.Login)

	profession := app.Group("/profession", r.authMW.AuthRequire())
	profession.Post("/", r.authMW.PermissionRequire(PermissionCreateProfession))
}

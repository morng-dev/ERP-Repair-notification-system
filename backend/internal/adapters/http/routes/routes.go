package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/morng-dev/erp/internal/adapters/http/handler"
	"github.com/morng-dev/erp/internal/adapters/http/middleware"
)

type Routes struct {
	authMW             *middleware.AuthMiddleware
	authHandler        *handler.AuthHandler
	profesHandler      *handler.ProfessionHandler
	permissionsHandler *handler.PermissionHandler
}

func NewRoutes(
	authMW *middleware.AuthMiddleware,
	authHandler *handler.AuthHandler,
	profesHandler *handler.ProfessionHandler,
	permissionsHandler *handler.PermissionHandler,

) *Routes {
	return &Routes{
		authHandler:        authHandler,
		authMW:             authMW,
		profesHandler:      profesHandler,
		permissionsHandler: permissionsHandler,
	}
}

const (
	PermissionCreateProfession = "profession:create"
	PermissionReadProfession   = "profession:read"
	PermissionUpdateProfession = "profession:update"
	PermissionDeleteProfession = "profession:delete"
)

func (r *Routes) SetUpRoute(app *fiber.App) {

	api := app.Group("/api/v1")
	//auth
	auth := api.Group("/auth")
	auth.Post("/register", r.authHandler.Register)
	auth.Post("/login", r.authHandler.Login)
	auth.Post("/forget-password", r.authHandler.ForgotPassword)
	auth.Post("/reset-password", r.authHandler.ResetPassword)
	auth.Post("/change-password", r.authMW.AuthRequire(), r.authHandler.Changepassword)
	//profession
	profession := app.Group("/profession", r.authMW.AuthRequire())
	profession.Post("/", r.authMW.PermissionRequire(PermissionCreateProfession))
	//permissions
	permission := app.Group("/permission", r.authMW.AuthRequire())
	permission.Post("/", r.permissionsHandler.CreatePermission)
	permission.Put("/", r.permissionsHandler.UpdatePermission)
}

package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/morng-dev/erp/internal/adapters/http/handler"
	"github.com/morng-dev/erp/internal/adapters/http/middleware"
)

type Routes struct {
	authMW             *middleware.AuthMiddleware
	ratelimitMW        *middleware.RatelimitMW
	authHandler        *handler.AuthHandler
	profesHandler      *handler.ProfessionHandler
	permissionsHandler *handler.PermissionHandler
	assetHandler       *handler.AssetHandler
}

func NewRoutes(
	authMW *middleware.AuthMiddleware,
	ratelimitMW *middleware.RatelimitMW,
	authHandler *handler.AuthHandler,
	profesHandler *handler.ProfessionHandler,
	permissionsHandler *handler.PermissionHandler,
	assetHandler *handler.AssetHandler,

) *Routes {
	return &Routes{
		authHandler:        authHandler,
		authMW:             authMW,
		profesHandler:      profesHandler,
		permissionsHandler: permissionsHandler,
		assetHandler:       assetHandler,
		ratelimitMW:        ratelimitMW,
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
	auth.Post("/login", r.ratelimitMW.LimitRequest(45, time.Second), r.authHandler.Login)
	auth.Post("/forget-password", r.authHandler.ForgotPassword)
	auth.Post("/reset-password", r.authHandler.ResetPassword)
	auth.Post("/change-password", r.authMW.AuthRequire(), r.authHandler.Changepassword)
	//profession
	profession := api.Group("/profession", r.authMW.AuthRequire())
	profession.Post("/", r.authMW.PermissionRequire(PermissionCreateProfession), r.profesHandler.CreateProfession)
	//permissions
	permission := api.Group("/permission", r.authMW.AuthRequire())
	permission.Post("/", r.permissionsHandler.CreatePermission)
	permission.Put("/", r.permissionsHandler.UpdatePermission)

	asset := api.Group("/asset")
	asset.Post("/", r.authMW.AuthRequire(), r.assetHandler.CreateAsset)
}

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
	userHandler        *handler.UserHandler
	profesHandler      *handler.ProfessionHandler
	permissionsHandler *handler.PermissionHandler
	assetHandler       *handler.AssetHandler
}

func NewRoutes(
	authMW *middleware.AuthMiddleware,
	ratelimitMW *middleware.RatelimitMW,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
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
const (
	login          = "login"
	Register       = "register"
	forgetpassword = "forgetpassword"
	resetpassword  = "resetpassword"
	changepassword = "changepassword"
)

func (r *Routes) SetUpRoute(app *fiber.App) {

	api := app.Group("/api/v1")
	//auth
	auth := api.Group("/auth")
	auth.Post("/register", r.ratelimitMW.RatelimitRequest(5, Register, time.Second), r.authHandler.Register)
	auth.Post("/login", r.ratelimitMW.RatelimitRequest(10, login, time.Second), r.authHandler.Login)
	auth.Post("/forget-password", r.ratelimitMW.RatelimitRequest(10, forgetpassword, time.Hour), r.authHandler.ForgotPassword)
	auth.Post("/reset-password", r.authHandler.ResetPassword)
	auth.Post("/change-password", r.ratelimitMW.RatelimitRequest(5, changepassword, 30*time.Minute), r.authMW.AuthRequire(), r.authHandler.Changepassword)
	//user
	user := api.Group("/user")
	user.Get("/", r.userHandler.GetUser)
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

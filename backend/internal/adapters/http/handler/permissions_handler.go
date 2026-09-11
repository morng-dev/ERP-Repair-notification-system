package handler

import "github.com/morng-dev/erp/internal/core/domain/ports/services"

type PermissionHandler struct {
	permissionService services.PermissionsService
}

func NewPermissionHandler(permissionService services.PermissionsService) *PermissionHandler {
	return &PermissionHandler{permissionService: permissionService}
}

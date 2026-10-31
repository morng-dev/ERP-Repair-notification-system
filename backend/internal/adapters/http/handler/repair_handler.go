package handler

import "github.com/morng-dev/erp/internal/core/domain/ports/services"

type RepairHandler struct {
	RepairService services.RepairService
}

func NewRepairHandler(RepairService services.RepairService) *RepairHandler {
	return &RepairHandler{RepairService: RepairService}
}

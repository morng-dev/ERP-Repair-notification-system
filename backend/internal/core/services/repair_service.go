package services

import (
	"context"
	"errors"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/kafkaService/kafkaproducer"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type RepairService struct {
	repairRepo repositories.RepairRepository
	userRepo   repositories.UserRepository
	publisher  kafkaproducer.NotificationKafkaService
}

func NewRepairService(
	repairRepo repositories.RepairRepository,
	userRepo repositories.UserRepository,
	publisher kafkaproducer.NotificationKafkaService,
) services.RepairService {
	return &RepairService{
		repairRepo: repairRepo,
		userRepo:   userRepo,
		publisher:  publisher,
	}
}

func (s *RepairService) CreateRepair(ctx context.Context, req *entities.RepairRequest) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	repair := &entities.Repair{
		ID:         id,
		Title:      req.Title,
		AssetID:    req.AssetID,
		ReportedBy: req.ReportedBy,
		AssignedTo: req.AssignedTo,
		Status:     entities.RepairStatusPending,
		Priority:   entities.RepairPriorityNormal,
	}
	if err := s.repairRepo.Create(ctx, repair); err != nil {
		return err
	}
	title, content := getRepairStatusMessage(repair.Status)
	s.notify(ctx, repair.ReportedBy, repair.ID, "repair_created", title, content)

	return nil
}

func (s *RepairService) UpdateStatusRepair(ctx context.Context, repairID uuid.UUID, status entities.RepairStatus) error {
	if !status.IsValid() {
		return errors.New("invalid repair status")
	}
	repair, err := s.repairRepo.GetByID(ctx, repairID)
	if err != nil {
		return err
	}
	if repair.Status == status {
		return errors.New("repair already has this status")
	}
	if !entities.CanTransition(repair.Status, status) {
		return errors.New("invalid status transition")
	}
	if err := s.repairRepo.UpdateStatus(ctx, repairID, status); err != nil {
		return err
	}
	title, content := getRepairStatusMessage(status)
	s.notify(ctx, repair.ReportedBy, repair.ID, "repair_status_updated", title, content)
	return nil
}

func (s *RepairService) GetAllRepair(ctx context.Context, page, limit int) ([]*entities.Repair, *entities.PaginationResponse, error) {
	repairs, total, err := s.repairRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, nil, err
	}
	totalPage := int(math.Ceil(float64(total) / float64(limit)))
	pagination := &entities.PaginationResponse{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPage,
	}
	return repairs, pagination, nil
}

func (s *RepairService) notify(ctx context.Context, userID uuid.UUID, repairID uuid.UUID, typ, title, content string) {
	id, err := uuid.NewV7()
	if err != nil {
		log.Printf("notify: gen id: %v", err)
		return
	}

	notif := &entities.Notification{
		ID:          id,
		Title:       title,
		ReferenceID: repairID,
		UserID:      userID,
		Type:        typ,
		Content:     content,
		Status:      "unread",
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.publisher.Publish(ctx, notif); err != nil {
		log.Printf("publish notification failed (user=%s): %v", userID, err)
	}
}

func getRepairStatusMessage(status entities.RepairStatus) (string, string) {
	switch status {
	case entities.RepairStatusInProgress:
		return "กำลังดำเนินการซ่อม", "เจ้าหน้าที่กำลังดำเนินการซ่อมอุปกรณ์ของคุณ"

	case entities.RepairStatusCompleted:
		return "ซ่อมเสร็จเรียบร้อย", "การดำเนินการซ่อมอุปกรณ์ของคุณเสร็จสิ้นแล้ว"

	case entities.RepairStatusCancelled:
		return "ยกเลิกการซ่อม", "คำขอแจ้งซ่อมของคุณถูกยกเลิก"

	default:
		return "อัปเดตสถานะการซ่อม", "สถานะคำขอแจ้งซ่อมของคุณมีการเปลี่ยนแปลง"
	}
}

package services

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/kafkaService/kafkaproducer"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type RepairService struct {
	repairRepo repositories.RepairRepository
	publisher  kafkaproducer.NotificationKafkaService
}

func NewRepairService(repairRepo repositories.RepairRepository, publisher kafkaproducer.NotificationKafkaService) services.RepairService {
	return &RepairService{
		repairRepo: repairRepo,
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
		Status:     "pending",
		Priority:   "normal",
	}
	if err := s.repairRepo.Create(ctx, repair); err != nil {
		return err
	}

	s.notify(ctx, repair.ReportedBy, repair.ID, "repair", "รับเรื่องแจ้งซ่อมแล้ว", "ระบบบันทึกคำขอแจ้งซ่อมของคุณเรียบร้อย รอเจ้าหน้าที่ดำเนินการ")

	return nil
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

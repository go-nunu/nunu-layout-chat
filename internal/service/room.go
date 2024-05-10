package service

import (
	"context"
)

type RoomService interface {
	GetRoom(ctx context.Context, id int64) error
}

func NewRoomService(service *Service) RoomService {
	return &roomService{
		Service: service,
	}
}

type roomService struct {
	*Service
}

func (s *roomService) GetRoom(ctx context.Context, id int64) error {
	return nil
}

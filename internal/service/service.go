package service

import (
	"nunu-layout-chat/internal/repository"
	"nunu-layout-chat/pkg/jwt"
	"nunu-layout-chat/pkg/log"
	"nunu-layout-chat/pkg/sid"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewService(tm repository.Transaction, logger *log.Logger, sid *sid.Sid, jwt *jwt.JWT) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
	}
}

package service

import (
	"github.com/go-kratos/kratos/v2/log"

	v1 "kratos-realworld/api/backend/v1"
	"kratos-realworld/internal/biz"
)

type BackendService struct {
	v1.UnimplementedBackendServer
	bc  *biz.BackendUsecase
	uc  *biz.UserUsecase
	fc  *biz.FileUsecase
	sc  *biz.ProfileUsecase
	log *log.Helper
}

func NewBackendService(bc *biz.BackendUsecase, uc *biz.UserUsecase, fc *biz.FileUsecase) *BackendService {
	return &BackendService{bc: bc, uc: uc, fc: fc}
}

package service

import (
	"github.com/edrlab/pubstore/pkg/stor"
)

type Service struct {
	stor *stor.Stor
}

func Init(s *stor.Stor) *Service {
	return &Service{stor: s}
}

func (service *Service) GetLicenceIdTransaction(publication *stor.Publication, user *stor.User) string {

	if publication == nil || user == nil {
		return ""
	}

	userID := user.ID
	pubID := publication.ID
	transaction, _ := service.stor.GetTransactionByUserAndPublication(userID, pubID)

	if transaction == nil {
		return ""
	}

	licenceUUID := transaction.LicenceId
	return licenceUUID
}

package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
)

type ServiceUsecase interface {
	ClearService() (err error)
	GetService() (status *models.Status, err error)
}

type ServiceUsecaseImpl struct {
	repoService repositories.ServiceRepository
}

func MakeServiceUseCase(service repositories.ServiceRepository) ServiceUsecase {
	return &ServiceUsecaseImpl{repoService: service}
}

func (serviceUsecase *ServiceUsecaseImpl) ClearService() (err error) {
	return serviceUsecase.repoService.ClearService()
}

func (serviceUsecase *ServiceUsecaseImpl) GetService() (status *models.Status, err error) {
	return serviceUsecase.repoService.GetService()
}

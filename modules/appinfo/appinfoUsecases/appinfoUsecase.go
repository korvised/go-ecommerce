package appinfoUsecases

import (
	"github.com/korvised/go-ecommerce/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
}

type appinfoUsecase struct {
	appinfoRepository appinfoRepositories.IAppinfoRepository
}

func AppinfoUsecase(appinfoRepository appinfoRepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepository,
	}
}

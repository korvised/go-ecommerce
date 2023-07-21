package appinfoHandlers

import (
	"github.com/korvised/go-ecommerce/config"
	"github.com/korvised/go-ecommerce/modules/appinfo/appinfoUsecases"
)

type IAppinfoHandler interface {
}

type appinfoHandler struct {
	cfg            config.IConfig
	appinfoUsecase appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, appinfoUsecase appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:            cfg,
		appinfoUsecase: appinfoUsecase,
	}
}

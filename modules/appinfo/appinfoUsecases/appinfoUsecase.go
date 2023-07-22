package appinfoUsecases

import (
	"github.com/korvised/go-ecommerce/modules/appinfo"
	"github.com/korvised/go-ecommerce/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategories(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(categoryId int) error
}

type appinfoUsecase struct {
	appinfoRepository appinfoRepositories.IAppinfoRepository
}

func AppinfoUsecase(appinfoRepository appinfoRepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepository,
	}
}

func (u *appinfoUsecase) FindCategories(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	return u.appinfoRepository.FindCategories(req)
}
func (u *appinfoUsecase) InsertCategory(req []*appinfo.Category) error {
	return u.appinfoRepository.InsertCategory(req)
}

func (u *appinfoUsecase) DeleteCategory(categoryId int) error {
	return u.appinfoRepository.DeleteCategory(categoryId)
}

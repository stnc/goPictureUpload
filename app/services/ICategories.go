package services

import (
	"stncCms/app/domain/entity"
)

//CatAppInterface service
type CatAppInterface interface {
	Save(*entity.Categories) (*entity.Categories, map[string]string)
	GetByID(uint64) (*entity.Categories, error)
	GetAll() ([]entity.Categories, error)

	Update(*entity.Categories) (*entity.Categories, map[string]string)
	Delete(uint64) error
}

//CatApp struct  init
type CatApp struct {
	request CatAppInterface
}

var _ CatAppInterface = &CatApp{}

//Save service init
func (f *CatApp) Save(Cat *entity.Categories) (*entity.Categories, map[string]string) {
	return f.request.Save(Cat)
}

//GetAll service init
func (f *CatApp) GetAll() ([]entity.Categories, error) {
	return f.request.GetAll()
}

//GetByID service init
func (f *CatApp) GetByID(CatID uint64) (*entity.Categories, error) {
	return f.request.GetByID(CatID)
}

//Update service init
func (f *CatApp) Update(Cat *entity.Categories) (*entity.Categories, map[string]string) {
	return f.request.Update(Cat)
}

//Delete service init
func (f *CatApp) Delete(CatID uint64) error {
	return f.request.Delete(CatID)
}

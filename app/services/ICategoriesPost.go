package services

import (
	"stncCms/app/domain/entity"
)

//CatPostAppInterface service
type CatPostAppInterface interface {
	Save(*entity.CategoryPosts) (*entity.CategoryPosts, map[string]string)
	GetAllforPostID(uint64) ([]entity.CategoryPosts, error)
	GetAllforCatID(uint64) ([]entity.CategoryPosts, error)
	GetAll() ([]entity.CategoryPosts, error)
	Update(*entity.CategoryPosts) (*entity.CategoryPosts, map[string]string)
	Delete(uint64) error
	DeleteForPostID(uint64) error
	DeleteForCatID(uint64) error
}

//CatPostApp struct  init
type CatPostApp struct {
	fr CatPostAppInterface
}

var _ CatPostAppInterface = &CatPostApp{}

//Save service init
func (f *CatPostApp) Save(Cat *entity.CategoryPosts) (*entity.CategoryPosts, map[string]string) {
	return f.fr.Save(Cat)
}

//GetAll service init
func (f *CatPostApp) GetAll() ([]entity.CategoryPosts, error) {
	return f.fr.GetAll()
}

//GetAllforPostID service init
func (f *CatPostApp) GetAllforPostID(PostID uint64) ([]entity.CategoryPosts, error) {
	return f.fr.GetAllforPostID(PostID)
}

//GetAllforCatID service init
func (f *CatPostApp) GetAllforCatID(CatID uint64) ([]entity.CategoryPosts, error) {
	return f.fr.GetAllforCatID(CatID)
}

//Update service init
func (f *CatPostApp) Update(Cat *entity.CategoryPosts) (*entity.CategoryPosts, map[string]string) {
	return f.fr.Update(Cat)
}

//Delete service init
func (f *CatPostApp) Delete(ID uint64) error {
	return f.fr.Delete(ID)
}

//DeleteForPostID service init
func (f *CatPostApp) DeleteForPostID(PostID uint64) error {
	return f.fr.DeleteForPostID(PostID)
}

//DeleteForCatID service init
func (f *CatPostApp) DeleteForCatID(CatID uint64) error {
	return f.fr.DeleteForCatID(CatID)
}

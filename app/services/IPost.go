package services

import (
	"stncCms/app/domain/entity"
)

//PostAppInterface interface
type PostAppInterface interface {
	Save(*entity.Post) (*entity.Post, map[string]string)
	GetByID(uint64) (*entity.Post, error)
	GetAll() ([]entity.Post, error)
	GetAllP(int, int) ([]entity.Post, error)
	Update(*entity.Post) (*entity.Post, map[string]string)
	Count(*int64)
	Delete(uint64) error
}
type postApp struct {
	request PostAppInterface
}

var _ PostAppInterface = &postApp{}

func (f *postApp) Count(postTotalCount *int64) {
	f.request.Count(postTotalCount)
}

func (f *postApp) Save(post *entity.Post) (*entity.Post, map[string]string) {
	return f.request.Save(post)
}

func (f *postApp) GetAll() ([]entity.Post, error) {
	return f.request.GetAll()
}

func (f *postApp) GetAllP(postsPerPage int, offset int) ([]entity.Post, error) {
	return f.request.GetAllP(postsPerPage, offset)
}

func (f *postApp) GetByID(postID uint64) (*entity.Post, error) {
	return f.request.GetByID(postID)
}

func (f *postApp) Update(post *entity.Post) (*entity.Post, map[string]string) {
	return f.request.Update(post)
}

func (f *postApp) Delete(postID uint64) error {
	return f.request.Delete(postID)
}

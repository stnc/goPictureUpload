package services

import (
	"stncCms/app/domain/entity"
)

//WebArchiveAppInterface interface
type WebArchiveLinksAppInterface interface {
	Save(*entity.WebArchiveLinks) (*entity.WebArchiveLinks, map[string]string)
	GetByID(uint64) (*entity.WebArchiveLinks, error)
	GetAll() ([]entity.WebArchiveLinks, error)
	GetAllP(int, int) ([]entity.WebArchiveLinks, error)
	Update(*entity.WebArchiveLinks) (*entity.WebArchiveLinks, map[string]string)
	Count(*int64)
	Delete(uint64) error
}

//WebArchiveLinksApp struct
type WebArchiveLinksApp struct {
	request WebArchiveLinksAppInterface
}

var _ WebArchiveLinksAppInterface = &WebArchiveLinksApp{}

//Count counter
func (f *WebArchiveLinksApp) Count(postTotalCount *int64) {
	f.request.Count(postTotalCount)
}

//Save data
func (f *WebArchiveLinksApp) Save(webarchive *entity.WebArchiveLinks) (*entity.WebArchiveLinks, map[string]string) {
	return f.request.Save(webarchive)
}

//GetAll all list
func (f *WebArchiveLinksApp) GetAll() ([]entity.WebArchiveLinks, error) {
	return f.request.GetAll()
}

//GetAllP all list for pagination
func (f *WebArchiveLinksApp) GetAllP(postsPerPage int, offset int) ([]entity.WebArchiveLinks, error) {
	return f.request.GetAllP(postsPerPage, offset)
}

//GetByID single row for id
func (f *WebArchiveLinksApp) GetByID(ID uint64) (*entity.WebArchiveLinks, error) {
	return f.request.GetByID(ID)
}

//Update data
func (f *WebArchiveLinksApp) Update(webarchive *entity.WebArchiveLinks) (*entity.WebArchiveLinks, map[string]string) {
	return f.request.Update(webarchive)
}

//Delete data
func (f *WebArchiveLinksApp) Delete(ID uint64) error {
	return f.request.Delete(ID)
}

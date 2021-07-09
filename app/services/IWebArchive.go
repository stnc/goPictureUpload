package services

import (
	"stncCms/app/domain/entity"
)

//WebArchiveAppInterface interface
type WebArchiveAppInterface interface {
	Save(*entity.WebArchive) (*entity.WebArchive, map[string]string)
	GetByID(uint64) (*entity.WebArchive, error)
	GetAll() ([]entity.WebArchive, error)
	GetAllP(int, int) ([]entity.WebArchive, error)
	Update(*entity.WebArchive) (*entity.WebArchive, map[string]string)
	Count(*int64)
	Delete(uint64) error
}

//WebArchiveApp struct
type WebArchiveApp struct {
	request WebArchiveAppInterface
}

var _ WebArchiveAppInterface = &WebArchiveApp{}

//Count counter
func (f *WebArchiveApp) Count(postTotalCount *int64) {
	f.request.Count(postTotalCount)
}

//Save data
func (f *WebArchiveApp) Save(webarchive *entity.WebArchive) (*entity.WebArchive, map[string]string) {
	return f.request.Save(webarchive)
}

//GetAll all list
func (f *WebArchiveApp) GetAll() ([]entity.WebArchive, error) {
	return f.request.GetAll()
}

//GetAllP all list for pagination
func (f *WebArchiveApp) GetAllP(perPage int, offset int) ([]entity.WebArchive, error) {
	return f.request.GetAllP(perPage, offset)
}

//GetByID single row for id
func (f *WebArchiveApp) GetByID(ID uint64) (*entity.WebArchive, error) {
	return f.request.GetByID(ID)
}

//Update data
func (f *WebArchiveApp) Update(webarchive *entity.WebArchive) (*entity.WebArchive, map[string]string) {
	return f.request.Update(webarchive)
}

//Delete data
func (f *WebArchiveApp) Delete(ID uint64) error {
	return f.request.Delete(ID)
}

package services

import (
	"stncCms/app/domain/entity"
)

//MediaAppInterface interface
type MediaAppInterface interface {
	Save(*entity.Media) (*entity.Media, map[string]string)
	GetByID(uint64) (*entity.Media, error)
	GetAll(modulID int, contentID int) ([]entity.Media, error)
	GetAllforModul(modulID int, hayvanID int) ([]entity.Media, error)
	GetAllP(int, int) ([]entity.Media, error)
	Update(*entity.Media) (*entity.Media, map[string]string)
	Count(modulID int, contentID int, MediaTotalCount *int)
	Delete(uint64) error
}
type MediaApp struct {
	request MediaAppInterface
}

var _ MediaAppInterface = &MediaApp{}

func (f *MediaApp) Count(modulID int, contentID int, MediaTotalCount *int) {
	f.request.Count(modulID, contentID, MediaTotalCount)
}

func (f *MediaApp) Save(Media *entity.Media) (*entity.Media, map[string]string) {
	return f.request.Save(Media)
}

func (f *MediaApp) GetAll(modulID int, contentID int) ([]entity.Media, error) {
	return f.request.GetAll(modulID, contentID)
}
func (f *MediaApp) GetAllforModul(modulID int, hayvanID int) ([]entity.Media, error) {
	return f.request.GetAllforModul(modulID, hayvanID)
}

func (f *MediaApp) GetAllP(MediasPerPage int, offset int) ([]entity.Media, error) {
	return f.request.GetAllP(MediasPerPage, offset)
}

func (f *MediaApp) GetByID(ID uint64) (*entity.Media, error) {
	return f.request.GetByID(ID)
}

func (f *MediaApp) Update(Media *entity.Media) (*entity.Media, map[string]string) {
	return f.request.Update(Media)
}

func (f *MediaApp) Delete(ID uint64) error {
	return f.request.Delete(ID)
}

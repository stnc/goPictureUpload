package services

import (
	"stncCms/app/domain/entity"
)

//LanguageAppInterface interface
type LanguageAppInterface interface {
	Save(*entity.Languages) (*entity.Languages, map[string]string)
	GetByID(uint64) (*entity.Languages, error)
	GetAll() ([]entity.Languages, error)
	Update(*entity.Languages) (*entity.Languages, map[string]string)
	Delete(uint64) error
}
type langApp struct {
	request LanguageAppInterface
}

var _ LanguageAppInterface = &langApp{}

func (l *langApp) Save(lang *entity.Languages) (*entity.Languages, map[string]string) {
	return l.request.Save(lang)
}

func (l *langApp) GetAll() ([]entity.Languages, error) {
	return l.request.GetAll()
}

func (l *langApp) GetByID(langID uint64) (*entity.Languages, error) {
	return l.request.GetByID(langID)
}

func (l *langApp) Update(lang *entity.Languages) (*entity.Languages, map[string]string) {
	return l.request.Update(lang)
}

func (l *langApp) Delete(langID uint64) error {
	return l.request.Delete(langID)
}

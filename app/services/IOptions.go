package services

import (
	"stncCms/app/domain/entity"
)

//OptionsAppInterface service
type OptionsAppInterface interface {
	GetAll() ([]entity.Options, error)
	GetOption(value string) string
	SetOption(name, value string)
	DeleteOptionID(id uint64) error
	DeleteOption(value string) error
}

//OptionsApp struct  init
type optionsApp struct {
	request OptionsAppInterface
}

var _ OptionsAppInterface = &optionsApp{}

//GetAll service init
func (f *optionsApp) GetAll() ([]entity.Options, error) {
	return f.request.GetAll()
}

//GetOption update
func (f *optionsApp) GetOption(value string) string {
	return f.request.GetOption(value)
}

//UpdateOption update
func (f *optionsApp) SetOption(name, value string) {
	f.request.SetOption(name, value)
}

//DeleteOptionID update
func (f *optionsApp) DeleteOptionID(id uint64) error {
	return f.request.DeleteOptionID(id)
}

//DeleteOption update
func (f *optionsApp) DeleteOption(value string) error {
	return f.request.DeleteOption(value)
}

package services

import (
	"stncCms/app/domain/dto"
	"stncCms/app/domain/entity"
)

//OdemelerAppInterface interface
type OdemelerAppInterface interface {
	Save(*entity.Odemeler) (*entity.Odemeler, map[string]string)
	SaveBatch([]entity.Odemeler) ([]entity.Odemeler, map[string]string)
	GetByID(uint64) (*entity.Odemeler, error)
	GetAll() ([]entity.Odemeler, error)
	Update(*entity.Odemeler) (*entity.Odemeler, map[string]string)
	SetUpdateHesaplar(*dto.Odemeler) (*dto.Odemeler, map[string]string)
	Delete(uint64) error
	LastPrice(uint64) (*dto.OdemelerSonFiyat, error)
	OdemelerToplami(kurbanID uint64) float64
	KurbanSonKalanUcret(kurbanID uint64) float64
	KasaBorcu(kurbanID uint64) float64
	GetOdemeRelation(odemeID uint64) (*dto.OdemeMakbuzu, error)
}
type odemelerApp struct {
	request OdemelerAppInterface
}

var _ OdemelerAppInterface = &odemelerApp{}

func (f *odemelerApp) Save(odemeler *entity.Odemeler) (*entity.Odemeler, map[string]string) {
	return f.request.Save(odemeler)
}

func (f *odemelerApp) SaveBatch(odemeler []entity.Odemeler) ([]entity.Odemeler, map[string]string) {
	return f.request.SaveBatch(odemeler)
}

func (f *odemelerApp) GetAll() ([]entity.Odemeler, error) {
	return f.request.GetAll()
}

func (f *odemelerApp) GetByID(odemelerID uint64) (*entity.Odemeler, error) {
	return f.request.GetByID(odemelerID)
}

func (f *odemelerApp) Update(odemeler *entity.Odemeler) (*entity.Odemeler, map[string]string) {
	return f.request.Update(odemeler)
}

func (f *odemelerApp) SetUpdateHesaplar(odemeler *dto.Odemeler) (*dto.Odemeler, map[string]string) {
	return f.request.SetUpdateHesaplar(odemeler)
}

func (f *odemelerApp) Delete(odemelerID uint64) error {
	return f.request.Delete(odemelerID)
}

func (f *odemelerApp) LastPrice(kurbanID uint64) (*dto.OdemelerSonFiyat, error) {
	return f.request.LastPrice(kurbanID)
}

func (f *odemelerApp) OdemelerToplami(kurbanID uint64) float64 {
	return f.request.OdemelerToplami(kurbanID)
}

func (f *odemelerApp) KurbanSonKalanUcret(kurbanID uint64) float64 {
	return f.request.KurbanSonKalanUcret(kurbanID)
}

func (f *odemelerApp) KasaBorcu(kurbanID uint64) float64 {
	return f.request.KasaBorcu(kurbanID)
}

func (f *odemelerApp) GetOdemeRelation(odemeID uint64) (*dto.OdemeMakbuzu, error) {
	return f.request.GetOdemeRelation(odemeID)
}

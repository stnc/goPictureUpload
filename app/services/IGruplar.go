package services

import (
	"stncCms/app/domain/dto"
	"stncCms/app/domain/entity"
)

//GruplarAppInterface interface
type GruplarAppInterface interface {
	Save(*entity.Gruplar) (*entity.Gruplar, map[string]string)
	GetByID(uint64) (*entity.Gruplar, error)
	GetByIDAllRelations(uint64) (*entity.Gruplar, error)
	GetByIDAllRelationsHayvanOlmayanlar(uint64) (*entity.Gruplar, error)
	GetByAllRelations() ([]dto.GruplarExcelandIndex, error)
	GetAll() ([]entity.Gruplar, error)
	GetAllFindDurum(int) ([]entity.Gruplar, error)
	GetAllFindDurumAndAgirlikTipi(int, int) ([]entity.Gruplar, error)
	ToplamOdemeler(uint64) float64
	KalanBorclar(uint64) float64
	KasaBorcu(uint64) float64
	GetAllP(int, int) ([]entity.Gruplar, error)
	Update(*entity.Gruplar) (*entity.Gruplar, map[string]string)
	Delete(uint64) error
	KurbanFiyati(uint64) float64
	SatisFiyatTuru(uint64) int
	Count(*int64)
	SetGrupID(GrupID uint64, kurbanID uint64, grupLideri int)
}

type gruplarApp struct {
	request GruplarAppInterface
}

var _ GruplarAppInterface = &gruplarApp{}

func (f *gruplarApp) Save(Kurban *entity.Gruplar) (*entity.Gruplar, map[string]string) {
	return f.request.Save(Kurban)
}

func (f *gruplarApp) GetAll() ([]entity.Gruplar, error) {
	return f.request.GetAll()
}

func (f *gruplarApp) GetAllP(postsPerPage int, offset int) ([]entity.Gruplar, error) {
	return f.request.GetAllP(postsPerPage, offset)
}

func (f *gruplarApp) GetByID(kurbanID uint64) (*entity.Gruplar, error) {
	return f.request.GetByID(kurbanID)
}

func (f *gruplarApp) GetByIDAllRelations(kurbanID uint64) (*entity.Gruplar, error) {
	return f.request.GetByIDAllRelations(kurbanID)
}

func (f *gruplarApp) GetByIDAllRelationsHayvanOlmayanlar(kurbanID uint64) (*entity.Gruplar, error) {
	return f.request.GetByIDAllRelationsHayvanOlmayanlar(kurbanID)
}

func (f *gruplarApp) GetAllFindDurum(durum int) ([]entity.Gruplar, error) {
	return f.request.GetAllFindDurum(durum)
}

func (f *gruplarApp) GetAllFindDurumAndAgirlikTipi(status int, agirilikTipi int) ([]entity.Gruplar, error) {
	return f.request.GetAllFindDurumAndAgirlikTipi(status, agirilikTipi)
}
func (f *gruplarApp) GetByAllRelations() ([]dto.GruplarExcelandIndex, error) {
	return f.request.GetByAllRelations()
}

func (f *gruplarApp) Update(Kurban *entity.Gruplar) (*entity.Gruplar, map[string]string) {
	return f.request.Update(Kurban)
}

func (f *gruplarApp) Delete(kurbanID uint64) error {
	return f.request.Delete(kurbanID)
}
func (f *gruplarApp) KurbanFiyati(kurbanID uint64) float64 {
	return f.request.KurbanFiyati(kurbanID)
}

func (f *gruplarApp) ToplamOdemeler(GrupID uint64) float64 {
	return f.request.ToplamOdemeler(GrupID)
}

func (f *gruplarApp) KalanBorclar(GrupID uint64) float64 {
	return f.request.KalanBorclar(GrupID)
}
func (f *gruplarApp) KasaBorcu(GrupID uint64) float64 {
	return f.request.KasaBorcu(GrupID)
}

func (f *gruplarApp) SatisFiyatTuru(kurbanID uint64) int {
	return f.request.SatisFiyatTuru(kurbanID)
}
func (f *gruplarApp) Count(postTotalCount *int64) {
	f.request.Count(postTotalCount)
}

func (f *gruplarApp) SetGrupID(grupID uint64, kurbanID uint64, grupLideri int) {
	f.request.SetGrupID(grupID, kurbanID, grupLideri)
}

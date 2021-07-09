package services

import (
	"stncCms/app/domain/dto"
	"stncCms/app/domain/entity"

	"github.com/gin-gonic/gin"
)

//KisilerAppInterface service
type KisilerAppInterface interface {
	Save(*entity.Kisiler) (*entity.Kisiler, map[string]string)
	GetByID(uint64) (*entity.Kisiler, error)
	GetByIDRel(uint64) (*entity.Kisiler, error)
	GetByReferansID(uint64) (*entity.Kisiler, error) // //TODO: tam olarak verecek belli değil kullanımda değil
	GetAll() ([]entity.Kisiler, error)
	GetAllP(int, int) ([]entity.Kisiler, error)
	Count(*int64)
	HasKisiKurban(id uint64, postTotalCount *int)
	HasKisiReferans(id uint64, postTotalCount *int)
	Search(string) (map[string]string, []entity.Kisiler, error)
	Update(*entity.Kisiler) (*entity.Kisiler, map[string]string)
	Delete(uint64) error
	ListDataTable(c *gin.Context) (dto.KisilerDataTableResult, error)
}

//KisilerApp struct  init
type KisilerApp struct {
	request KisilerAppInterface
}

var _ KisilerAppInterface = &KisilerApp{}

//Save service init
func (f *KisilerApp) Save(Kisiler *entity.Kisiler) (*entity.Kisiler, map[string]string) {
	return f.request.Save(Kisiler)
}

//GetAll service init
func (f *KisilerApp) GetAll() ([]entity.Kisiler, error) {
	return f.request.GetAll()
}

//GetByID service init
func (f *KisilerApp) GetByID(kisiID uint64) (*entity.Kisiler, error) {
	return f.request.GetByID(kisiID)
}

//GetByIDRel service init
func (f *KisilerApp) GetByIDRel(kisiID uint64) (*entity.Kisiler, error) {
	return f.request.GetByIDRel(kisiID)
}

//GetByID service init  //TODO: tam olarak verecek belli değil kullanımda değil
func (f *KisilerApp) GetByReferansID(referansID uint64) (*entity.Kisiler, error) {
	return f.request.GetByReferansID(referansID)
}

//Update service init
func (f *KisilerApp) Update(Kisiler *entity.Kisiler) (*entity.Kisiler, map[string]string) {
	return f.request.Update(Kisiler)
}

//Delete service init
func (f *KisilerApp) Delete(KisilerID uint64) error {
	return f.request.Delete(KisilerID)
}

func (f *KisilerApp) Search(value string) (map[string]string, []entity.Kisiler, error) {
	return f.request.Search(value)
}

func (f *KisilerApp) Count(postTotalCount *int64) {
	f.request.Count(postTotalCount)
}

func (f *KisilerApp) HasKisiKurban(id uint64, postTotalCount *int) {
	f.request.HasKisiKurban(id, postTotalCount)
}

func (f *KisilerApp) HasKisiReferans(id uint64, postTotalCount *int) {
	f.request.HasKisiReferans(id, postTotalCount)
}

func (f *KisilerApp) GetAllP(postsPerPage int, offset int) ([]entity.Kisiler, error) {
	return f.request.GetAllP(postsPerPage, offset)
}
func (f *KisilerApp) ListDataTable(c *gin.Context) (dto.KisilerDataTableResult, error) {
	return f.request.ListDataTable(c)
}

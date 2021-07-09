package services

import (
	"stncCms/app/domain/dto"
	"stncCms/app/domain/entity"

	"github.com/gin-gonic/gin"
)

//HayvanSatisYerleriAppInterface service
type HayvanSatisYerleriAppInterface interface {
	Save(*entity.HayvanSatisYerleri) (*entity.HayvanSatisYerleri, map[string]string)
	GetByID(uint64) (*entity.HayvanSatisYerleri, error)
	GetAll() ([]entity.HayvanSatisYerleri, error)
	ListDataTable(c *gin.Context) (dto.HayvanSatisYerleriDataTableResult, error)
	GetAllP(int, int) ([]entity.HayvanSatisYerleri, error)
	Update(*entity.HayvanSatisYerleri) (*entity.HayvanSatisYerleri, map[string]string)
	Count(*int64)
	Delete(uint64) error
}

//HayvanSatisYerleriApp struct  init
type hayvanSatisYerleriApp struct {
	request HayvanSatisYerleriAppInterface
}

var _ HayvanSatisYerleriAppInterface = &hayvanSatisYerleriApp{}

//Save service init
func (f *hayvanSatisYerleriApp) Save(Cat *entity.HayvanSatisYerleri) (*entity.HayvanSatisYerleri, map[string]string) {
	return f.request.Save(Cat)
}

//GetAll service init
func (f *hayvanSatisYerleriApp) GetAll() ([]entity.HayvanSatisYerleri, error) {
	return f.request.GetAll()
}

func (f *hayvanSatisYerleriApp) Count(totalCount *int64) {
	f.request.Count(totalCount)
}

func (f *hayvanSatisYerleriApp) GetAllP(perPage int, offset int) ([]entity.HayvanSatisYerleri, error) {
	return f.request.GetAllP(perPage, offset)
}

//GetByID service init
func (f *hayvanSatisYerleriApp) GetByID(CatID uint64) (*entity.HayvanSatisYerleri, error) {
	return f.request.GetByID(CatID)
}

//Update service init
func (f *hayvanSatisYerleriApp) Update(data *entity.HayvanSatisYerleri) (*entity.HayvanSatisYerleri, map[string]string) {
	return f.request.Update(data)
}

//Delete service init
func (f *hayvanSatisYerleriApp) Delete(id uint64) error {
	return f.request.Delete(id)
}

//ListDataTable service init
func (f *hayvanSatisYerleriApp) ListDataTable(c *gin.Context) (dto.HayvanSatisYerleriDataTableResult, error) {
	return f.request.ListDataTable(c)
}

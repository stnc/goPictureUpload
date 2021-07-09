package repository

import (
	"errors"
	"fmt"
	"os"
	"stncCms/app/domain/dto"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stnchelper"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//HayvanSatisYerleriRepo struct
type HayvanSatisYerleriRepo struct {
	db *gorm.DB
}

//HayvanSatisYerleriRepositoryInit initial
func HayvanSatisYerleriRepositoryInit(db *gorm.DB) *HayvanSatisYerleriRepo {
	return &HayvanSatisYerleriRepo{db}
}

//HayvanSatisYerleriRepo implements the repository.HayvanSatisYerleriRepository interface
// var _ interfaces.HayvanSatisYerleriAppInterface = &HayvanSatisYerleriRepo{}

//Save data
func (r *HayvanSatisYerleriRepo) Save(data *entity.HayvanSatisYerleri) (*entity.HayvanSatisYerleri, map[string]string) {
	dbErr := map[string]string{}

	var err error
	err = r.db.Debug().Create(&data).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "data title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return data, nil
}

//Update upate data
func (r *HayvanSatisYerleriRepo) Update(data *entity.HayvanSatisYerleri) (*entity.HayvanSatisYerleri, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&data).Error
	//db.Table("libraries").Where("id = ?", id).Update(dataData)

	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return data, nil
}

//Count fat
func (r *HayvanSatisYerleriRepo) Count(dataTotalCount *int64) {
	var data entity.HayvanSatisYerleri
	var count int64
	r.db.Debug().Model(data).Count(&count)
	*dataTotalCount = count
}

//Delete data
func (r *HayvanSatisYerleriRepo) Delete(id uint64) error {
	var data entity.HayvanSatisYerleri
	var err error
	err = r.db.Debug().Where("id = ? and durum <> 0", id).Delete(&data).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//GetByID get data
func (r *HayvanSatisYerleriRepo) GetByID(id uint64) (*entity.HayvanSatisYerleri, error) {
	var data entity.HayvanSatisYerleri
	var err error
	err = r.db.Debug().Where("id = ? and durum <> 0", id).Take(&data).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return &data, nil
}

//GetAll all data
func (r *HayvanSatisYerleriRepo) GetAll() ([]entity.HayvanSatisYerleri, error) {
	var datas []entity.HayvanSatisYerleri
	var err error
	err = r.db.Debug().Where("durum <> 0").Order("created_at desc").Find(&datas).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return datas, nil
}

//GetAllP pagination all data
func (r *HayvanSatisYerleriRepo) GetAllP(perPage int, offset int) ([]entity.HayvanSatisYerleri, error) {
	var datas []entity.HayvanSatisYerleri

	var err error
	err = r.db.Debug().Where("durum <> 0").Limit(perPage).Offset(offset).Order("created_at desc").Find(&datas).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	fmt.Println(datas)
	return datas, nil
}

//ListPostData upate data
func (r *HayvanSatisYerleriRepo) ListDataTable(c *gin.Context) (dto.HayvanSatisYerleriDataTableResult, error) {
	var total, filtered int
	var err error
	var data []dto.HayvanSatisYerleriDataTable
	query := r.db.Table(entity.HayvanSatisYerleriTableName)
	query = query.Where("id <> 1")
	query = query.Offset(stnchelper.QueryOffset(c))
	query = query.Limit(stnchelper.QueryLimit(c))
	query = query.Order(r.queryOrder(c))
	query = query.Scopes(r.searchScope(c), stnchelper.DateTimeScope(c))
	err = query.Find(&data).Error
	query = query.Offset(0)
	query.Table(entity.HayvanSatisYerleriTableName).Count(&filtered)
	// Total data count
	r.db.Table(entity.HayvanSatisYerleriTableName).Count(&total)

	result := dto.HayvanSatisYerleriDataTableResult{
		Total:    total,
		Filtered: filtered,
		Data:     data,
	}
	return result, err
}

func (r *HayvanSatisYerleriRepo) queryOrder(c *gin.Context) string {
	columnMap := map[string]string{
		"0": "id",
		"1": "yer_adi",
		"2": "ilgili_kisi",
		"3": "adresi",
		"4": "telefon",
	}

	column := c.DefaultQuery("order[0][column]", "0")
	dir := c.DefaultQuery("order[0][dir]", "desc")
	orderString := columnMap[column] + " " + dir

	return orderString
}

func (r *HayvanSatisYerleriRepo) searchScope(c *gin.Context) func(DB *gorm.DB) *gorm.DB {
	return func(DB *gorm.DB) *gorm.DB {
		query := DB
		search := c.QueryMap("search")
		fmt.Println(search)
		dbdriver := os.Getenv("DB_DRIVER")

		if dbdriver == "mysql" {
			if search["value"] != "" {
				query = query.Where(" yer_adi LIKE ? ", "%"+search["value"]+"%")
				query = query.Or("ilgili_kisi LIKE ? ", "%"+search["value"]+"%")
				query = query.Or("adresi LIKE ? ", "%"+search["value"]+"%")
				query = query.Or("telefon LIKE ? ", "%"+search["value"]+"%")
			}
		} else if dbdriver == "postgres" {
			if search["value"] != "" {
				query = query.Where(" yer_adi ILIKE ? ", "%"+search["value"]+"%")
				query = query.Or("ilgili_kisi ILIKE ? ", "%"+search["value"]+"%")
				query = query.Or("adresi ILIKE ? ", "%"+search["value"]+"%")
				query = query.Or("telefon ILIKE ? ", "%"+search["value"]+"%")
			}
		}
		return query
	}
}

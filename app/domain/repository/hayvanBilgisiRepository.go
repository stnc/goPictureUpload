package repository

import (
	"errors"
	"fmt"
	"math"
	"os"
	"stncCms/app/domain/dto"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stnccollection"
	"stncCms/app/domain/helpers/stnchelper"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var hayvanBilgisiTableName string = "hayvan_bilgisi"

//HayvanBilgisiRepo struct
type HayvanBilgisiRepo struct {
	db *gorm.DB
}

//HayvanBilgisiRepositoryInit initial
func HayvanBilgisiRepositoryInit(db *gorm.DB) *HayvanBilgisiRepo {
	return &HayvanBilgisiRepo{db}
}

//HayvanBilgisiRepo implements the repository.HayvanBilgisiRepository interface
// var _ interfaces.HayvanBilgisiAppInterface = &HayvanBilgisiRepo{}

//Save data
func (r *HayvanBilgisiRepo) Save(data *entity.HayvanBilgisi) (*entity.HayvanBilgisi, map[string]string) {
	dbErr := map[string]string{}

	var err error
	err = r.db.Debug().Omit("HayvanSatisYerleri").Create(&data).Error
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
func (r *HayvanBilgisiRepo) Update(data *entity.HayvanBilgisi) (*entity.HayvanBilgisi, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Omit("HayvanSatisYerleri").Save(&data).Error
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

//Delete data
func (r *HayvanBilgisiRepo) Delete(id uint64) error {
	var data entity.HayvanBilgisi
	var err error
	err = r.db.Debug().Where("id = ?", id).Delete(&data).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//GetByID get data
func (r *HayvanBilgisiRepo) GetByID(id uint64) (*entity.HayvanBilgisi, error) {
	var data entity.HayvanBilgisi
	var err error
	err = r.db.Debug().Where("id = ?", id).Take(&data).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return &data, nil
}

//TODO: buradaki değişken isimleri hiç anlaşılır değil düzeltilmesi lazım
//GetByIDRelated get data /admin/hayvanBilgisi/hayvanListeAjax/ uses
func (r *HayvanBilgisiRepo) GetByIDRelated(id uint64) (*dto.HayvanBilgisi, error) {

	var data dto.HayvanBilgisi
	var err error

	var satisFiyati1, satisFiyati2, satisFiyati3, agirlikfloat float64

	appOptions := OptionRepositoryInit(r.db)
	fiyatStr1 := appOptions.GetOption("satis_birim_fiyati_1")
	fiyatStr2 := appOptions.GetOption("satis_birim_fiyati_2")
	fiyatStr3 := appOptions.GetOption("satis_birim_fiyati_3")
	hisseAdeti := appOptions.GetOption("hisse_adeti")

	err = r.db.Debug().Where("id = ?", id).Preload("HayvanSatisYerleri").Find(&data).Error

	agirlikfloat = float64(data.Agirlik)

	satisFiyati1, _ = stnccollection.StringToFloat64(fiyatStr1)
	satisFiyati2, _ = stnccollection.StringToFloat64(fiyatStr2)
	satisFiyati3, _ = stnccollection.StringToFloat64(fiyatStr3)
	hisseAdetim, _ := stnccollection.StringToFloat64(hisseAdeti)

	data.Fiyat1 = satisFiyati1 * agirlikfloat
	data.SatisFiyati1 = satisFiyati1

	data.Fiyat2 = satisFiyati2 * agirlikfloat
	data.SatisFiyati2 = satisFiyati2

	data.Fiyat3 = satisFiyati3 * agirlikfloat
	data.SatisFiyati3 = satisFiyati3
	dusen1 := (satisFiyati1 * agirlikfloat) / hisseAdetim
	dusen2 := (satisFiyati2 * agirlikfloat) / hisseAdetim
	dusen3 := (satisFiyati3 * agirlikfloat) / hisseAdetim

	data.KisiBasiDusen1 = stnccollection.SayiYuvarla(math.Ceil(stnccollection.ToFixedDecimal(dusen1, 2)))
	data.KisiBasiDusen2 = stnccollection.SayiYuvarla(math.Ceil(stnccollection.ToFixedDecimal(dusen2, 2)))
	data.KisiBasiDusen3 = stnccollection.SayiYuvarla(math.Ceil(stnccollection.ToFixedDecimal(dusen3, 2)))

	if err != nil {
		return nil, errors.New("database error, please try again")
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return &data, nil
}

// https://www.mindbowser.com/golang-go-with-gorm/
//  err := Db.Model(&places).Association("town").Find(&places.Town).Error

// https://github.com/stnc-go/gorm_example/blob/master/1-to-1-relationship.go

//GetAllFindDurum all data
func (r *HayvanBilgisiRepo) GetAllFindDurum(durum int) ([]entity.HayvanBilgisi, error) {
	var datas []entity.HayvanBilgisi
	var err error
	err = r.db.Debug().Where("durum = ?", durum).Order("created_at desc").Find(&datas).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return datas, nil
}

//GetAllP pagination all data
func (r *HayvanBilgisiRepo) GetAllP(perPage int, offset int) ([]entity.HayvanBilgisi, error) {
	var datas []entity.HayvanBilgisi

	var err error
	err = r.db.Debug().Limit(perPage).Offset(offset).Order("created_at desc").Find(&datas).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	fmt.Println(datas)
	return datas, nil
}

//Count fat
func (r *HayvanBilgisiRepo) Count(dataTotalCount *int64) {
	var data entity.HayvanBilgisi
	var count int64
	r.db.Debug().Model(data).Count(&count)
	*dataTotalCount = count
}

//KupeNoCount fat
func (r *HayvanBilgisiRepo) KupeNoCount(kupeNo int, dataTotalCount *int) {
	var data entity.HayvanBilgisi
	var count int
	r.db.Debug().Model(data).Where("kupe_no = ?", kupeNo).Count(&count)
	*dataTotalCount = count
}

//UpdateSingleStatus upate data
func (r *HayvanBilgisiRepo) UpdateSingleStatus(id uint64, status int) {
	r.db.Debug().Table(hayvanBilgisiTableName).Where("id = ?", id).Update("durum", status)
}

//GetKupeNo upate data
func (r *HayvanBilgisiRepo) GetKupeNo(id uint64) (result string) {
	row := r.db.Debug().Table(hayvanBilgisiTableName).Select("kupe_no").Where("id = ? ", id).Order("id desc").Row()
	row.Scan(&result)
	return result
}

//GetAll all data
func (r *HayvanBilgisiRepo) GetAll() ([]entity.HayvanBilgisi, error) {
	var hayvanBilgisi []entity.HayvanBilgisi
	var err error
	err = r.db.Debug().Order("created_at desc").Find(&hayvanBilgisi).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return hayvanBilgisi, nil
}

//ListPostData upate data
func (r *HayvanBilgisiRepo) ListDataTable(c *gin.Context) (entity.HayvanBilgisiDataResult, error) {
	var total, filtered int
	var err error
	var data []entity.HayvanBilgisi
	query := r.db.Table(hayvanBilgisiTableName)
	query = query.Offset(stnchelper.QueryOffset(c))
	query = query.Limit(stnchelper.QueryLimit(c))
	query = query.Order(r.queryOrder(c))
	query = query.Scopes(r.searchScope(c), stnchelper.DateTimeScope(c))
	err = query.Find(&data).Error
	// if gorm.IsRecordNotFoundError(err) {
	// 	return nil, errors.New("data not found")
	// }
	query = query.Offset(0)
	query.Table(hayvanBilgisiTableName).Count(&filtered)
	// Total data count
	r.db.Table(hayvanBilgisiTableName).Count(&total)

	result := entity.HayvanBilgisiDataResult{
		Total:    total,
		Filtered: filtered,
		Data:     data,
	}
	return result, err
}

func (r *HayvanBilgisiRepo) queryOrder(c *gin.Context) string {
	columnMap := map[string]string{
		"0": "id",
		"1": "kupe_no",
		"2": "agirlik",
		"3": "alis_fiyati",

		// "5": "created_at",
	}

	column := c.DefaultQuery("order[0][column]", "0")
	dir := c.DefaultQuery("order[0][dir]", "desc")
	orderString := columnMap[column] + " " + dir

	return orderString
}

func (r *HayvanBilgisiRepo) searchScope(c *gin.Context) func(DB *gorm.DB) *gorm.DB {
	return func(DB *gorm.DB) *gorm.DB {
		query := DB
		search := c.QueryMap("search")
		fmt.Println(search)
		dbdriver := os.Getenv("DB_DRIVER")

		if dbdriver == "mysql" {
			if search["value"] != "" {
				query = query.Where("kupe_no LIKE ? ", "%"+search["value"]+"%")
				query = query.Or("agirlik LIKE ? ", "%"+search["value"]+"%")
				query = query.Or("alis_fiyati LIKE ? ", "%"+search["value"]+"%")

			}
		} else if dbdriver == "postgres" {
			if search["value"] != "" {
				query = query.Where(" CAST (kupe_no AS TEXT) ILIKE ? ", "%"+search["value"]+"%")
				query = query.Or("CAST (agirlik AS TEXT) ILIKE ? ", "%"+search["value"]+"%")
				query = query.Or("CAST (alis_fiyati AS TEXT) ILIKE ? ", "%"+search["value"]+"%")
			}
		}
		return query
	}
}

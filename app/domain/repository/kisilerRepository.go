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

//KisilerRepo struct
type KisilerRepo struct {
	db *gorm.DB
}

//KisilerRepositoryInit initial
func KisilerRepositoryInit(db *gorm.DB) *KisilerRepo {
	return &KisilerRepo{db}
}

//PostRepo implements the repository.PostRepository interface
// var _ interfaces.KisilerAppInterface = &KisilerRepo{}

//Save data
func (r *KisilerRepo) Save(Kisiler *entity.Kisiler) (*entity.Kisiler, map[string]string) {
	dbErr := map[string]string{}

	err := r.db.Debug().Create(&Kisiler).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "post title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return Kisiler, nil
}

//Update upate data
func (r *KisilerRepo) Update(Kisiler *entity.Kisiler) (*entity.Kisiler, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&Kisiler).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return Kisiler, nil
}

//GetByID get data
func (r *KisilerRepo) GetByID(id uint64) (*entity.Kisiler, error) {
	var Kisiler entity.Kisiler
	err := r.db.Debug().Where("id = ?", id).Take(&Kisiler).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &Kisiler, nil
}

//GetByID get data
func (r *KisilerRepo) GetByIDRel(id uint64) (*entity.Kisiler, error) {
	var kisiler entity.Kisiler
	// err := r.db.Debug().Where("id = ?", id).Take(&Kisiler).Error
	err := r.db.Debug().Where("id = ?", id).Preload("Kurbanlar", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}).Take(&kisiler).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &kisiler, nil
}

//GetByReferans referans data //TODO: tam olarak verecek belli değil kullanımda değil
func (r *KisilerRepo) GetByReferansID(id uint64) (*entity.Kisiler, error) {
	var Kisiler entity.Kisiler
	err := r.db.Debug().Where("referans_kisi1 = ?", id).Take(&Kisiler).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &Kisiler, nil
}

//GetAll all data
func (r *KisilerRepo) GetAll() ([]entity.Kisiler, error) {
	var Kisiler []entity.Kisiler
	err := r.db.Debug().Order("created_at desc").Find(&Kisiler).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return Kisiler, nil
}

//Delete delete data
func (r *KisilerRepo) Delete(id uint64) error {
	var Kisiler entity.Kisiler
	// err := r.db.Debug().Where("id = ?", id).Delete(&Kisiler).Error//soft delete
	err := r.db.Debug().Unscoped().Where("id = ?", id).Delete(&Kisiler).Error //#hard delete

	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//Search   data
func (r *KisilerRepo) Search(search string) (map[string]string, []entity.Kisiler, error) {
	var data []entity.Kisiler
	var err error
	var count int64
	dbdriver := os.Getenv("DB_DRIVER")

	if dbdriver == "mysql" {
		err = r.db.Debug().Where("ad_soyad LIKE ?", "%"+search+"%").Find(&data).Error
		r.db.Debug().Where("ad_soyad LIKE ?", "%"+search+"%").Model(&data).Count(&count)
	} else if dbdriver == "postgres" {
		err = r.db.Debug().Where("ad_soyad ILIKE ?", "%"+search+"%").Find(&data).Error
		r.db.Debug().Where("ad_soyad ILIKE ?", "%"+search+"%").Model(&data).Count(&count)
	}

	if err != nil {
		fmt.Println("err nil")
		return map[string]string{"status": "error"}, nil, err
	}

	if gorm.IsRecordNotFoundError(err) {
		fmt.Println("null nil")
		return map[string]string{"status": "error"}, nil, errors.New("post not found")
	}

	if count == 0 {
		fmt.Println("data nil")
		mapD := map[string]string{"status": "not found", "html": "<span></span>"}
		//mapB, _ := json.Marshal(mapD)
		//fmt.Println(string(mapB))
		return mapD, nil, err
	} else {
		return map[string]string{"status": "ok"}, data, nil
	}

}

//Count fat
func (r *KisilerRepo) Count(postTotalCount *int64) {
	var post entity.Kurban
	var count int64
	r.db.Debug().Model(post).Count(&count)
	*postTotalCount = count
}

//GetAllP pagination all data
func (r *KisilerRepo) GetAllP(postsPerPage int, offset int) ([]entity.Kisiler, error) {
	var posts []entity.Kisiler
	var err error
	err = r.db.Debug().Limit(postsPerPage).Offset(offset).Order("id asc").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return posts, nil
}

//ListPostData upate data
func (r *KisilerRepo) ListDataTable(c *gin.Context) (dto.KisilerDataTableResult, error) {
	var total, filtered int
	var err error
	var data []dto.Kisiler
	query := r.db.Table(dto.KisilerTableName)
	query = query.Where("id <> 1")
	query = query.Offset(stnchelper.QueryOffset(c))
	query = query.Limit(stnchelper.QueryLimit(c))
	query = query.Order(r.queryOrder(c))
	query = query.Scopes(r.searchScope(c), stnchelper.DateTimeScope(c))
	err = query.Find(&data).Error
	query = query.Offset(0)
	query.Table(dto.KisilerTableName).Count(&filtered)
	// Total data count
	r.db.Table(dto.KisilerTableName).Count(&total)

	result := dto.KisilerDataTableResult{
		Total:    total,
		Filtered: filtered,
		Data:     data,
	}
	return result, err
}

func (r *KisilerRepo) queryOrder(c *gin.Context) string {
	columnMap := map[string]string{
		"0": "id",
		"1": "ad_soyad",
		"2": "telefon",
		"3": "adres",
	}

	column := c.DefaultQuery("order[0][column]", "0")
	dir := c.DefaultQuery("order[0][dir]", "desc")
	orderString := columnMap[column] + " " + dir

	return orderString
}

func (r *KisilerRepo) searchScope(c *gin.Context) func(DB *gorm.DB) *gorm.DB {
	return func(DB *gorm.DB) *gorm.DB {
		query := DB
		search := c.QueryMap("search")
		fmt.Println(search)
		dbdriver := os.Getenv("DB_DRIVER")

		if dbdriver == "mysql" {
			if search["value"] != "" {
				query = query.Where(" ad_soyad LIKE ? ", "%"+search["value"]+"%")
				query = query.Or("telefon LIKE ? ", "%"+search["value"]+"%")
				query = query.Or("adres LIKE ? ", "%"+search["value"]+"%")

			}
		} else if dbdriver == "postgres" {
			if search["value"] != "" {
				query = query.Where(" ad_soyad ILIKE ? ", "%"+search["value"]+"%")
				query = query.Or("telefon ILIKE ? ", "%"+search["value"]+"%")
				query = query.Or("adres ILIKE ? ", "%"+search["value"]+"%")
			}
		}
		return query
	}
}

//HasKisiKurban kişi kurban listesinde var mı
func (r *KisilerRepo) HasKisiKurban(id uint64, postTotalCount *int) {
	var kurbanData entity.Kurban
	var count int
	r.db.Debug().Model(kurbanData).Where("kisi_id = ?", id).Count(&count)
	*postTotalCount = count
}

//HasKisiReferans kişi referans listesinde var mı?
func (r *KisilerRepo) HasKisiReferans(id uint64, postTotalCount *int) {
	var kisiData entity.Kisiler
	var count int
	r.db.Debug().Model(kisiData).Where("referans_kisi1 = ?", id).Count(&count)
	*postTotalCount = count
}

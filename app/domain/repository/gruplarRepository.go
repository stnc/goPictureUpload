package repository

import (
	"errors"
	"fmt"
	"stncCms/app/domain/dto"
	"stncCms/app/domain/entity"
	"strings"

	"github.com/jinzhu/gorm"
)

//GruplarRepo struct
type GruplarRepo struct {
	db *gorm.DB
}

//GruplarRepositoryInit initial
func GruplarRepositoryInit(db *gorm.DB) *GruplarRepo {
	return &GruplarRepo{db}
}

//GruplarRepo implements the repository.KurbanRepository interface
// var _ interfaces.PostAppInterface = &GruplarRepo{}

//Save data
func (r *GruplarRepo) Save(post *entity.Gruplar) (*entity.Gruplar, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&post).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "post title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return post, nil
}

//Update upate data
func (r *GruplarRepo) Update(post *entity.Gruplar) (*entity.Gruplar, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&post).Error
	//db.Table("libraries").Where("id = ?", id).Update(postData)

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
	return post, nil
}

//Delete data
func (r *GruplarRepo) Delete(id uint64) error {
	var post entity.Gruplar
	err := r.db.Debug().Where("id = ?", id).Delete(&post).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//GetByID get data
func (r *GruplarRepo) GetByID(id uint64) (*entity.Gruplar, error) {
	var post entity.Gruplar
	err := r.db.Debug().Where("id = ? and durum <> 1  ", id).Preload("Kurban").Take(&post).Error
	fmt.Printf("%+v\n", post)
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &post, nil
}

//GetByIDAllRelations get //TODO:  data hayvanı atanmış olnaları getirir
func (r *GruplarRepo) GetByIDAllRelations(id uint64) (*entity.Gruplar, error) {
	var gruplarList entity.Gruplar
	err := r.db.Debug().Where("id = ? and  durum <> 1 and hayvan_bilgisi_id <> 0 ", id).Preload("Kurban").Take(&gruplarList).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &gruplarList, nil
}

//GetByIDAllRelations get //TODO:  data hayvanı atanmamış olanları getirir
func (r *GruplarRepo) GetByIDAllRelationsHayvanOlmayanlar(id uint64) (*entity.Gruplar, error) {
	var gruplarList entity.Gruplar
	// var err error
	err := r.db.Debug().Where("id = ? and  durum <> 1 and hayvan_bilgisi_id = 0 ", id).Preload("Kurban").Take(&gruplarList).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &gruplarList, nil
}

//GetByAllRelations listeleme
func (r *GruplarRepo) GetByAllRelations() ([]dto.GruplarExcelandIndex, error) {
	var gruplarList []dto.GruplarExcelandIndex

	// err = r.db.Debug().Where("  durum <> 1  ").Preload("Kurban").Find(&gruplarList).Error //entity.gruplara gore
	err := r.db.Debug().Table("gruplar").Where("  durum <> 1  ").Order("kesim_sira_no asc").Find(&gruplarList).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return gruplarList, nil
}

//TODO: burayı iptal et kullanan yok sanırım
//GetAllFindDurum get data
func (r *GruplarRepo) GetAllFindDurum(durum int) ([]entity.Gruplar, error) {
	var gruplarList []entity.Gruplar
	err := r.db.Debug().Where("  durum = ?  ", durum).Find(&gruplarList).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return gruplarList, nil
}

//GetAllFindDurum all data
func (r *GruplarRepo) GetAllFindDurumAndAgirlikTipi(durum int, agirlikTipi int) ([]entity.Gruplar, error) {
	var datas []entity.Gruplar
	err := r.db.Debug().Where("durum = ? and agirlik_tipi= ?", durum, agirlikTipi).Preload("Kurban").Order("id asc").Find(&datas).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return datas, nil
}

//GetAll all data
func (r *GruplarRepo) GetAll() ([]entity.Gruplar, error) {
	var rows []entity.Gruplar
	err := r.db.Debug().Where("  durum <> 1  ").Order("kesim_sira_no asc").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return rows, nil
}

//GetAllP pagination all data
func (r *GruplarRepo) GetAllP(postsPerPage int, offset int) ([]entity.Gruplar, error) {
	var posts []entity.Gruplar
	// var err error
	err := r.db.Debug().Where(" durum <> 1  ").Limit(postsPerPage).Offset(offset).Order("created_at desc").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return posts, nil
}

//KurbanFiyati kurban id ye gore odenen miktar toplamı
func (r *GruplarRepo) KurbanFiyati(kurbanID uint64) float64 {
	var result float64
	row := r.db.Debug().Table(entity.KurbanTableName).Select("kurban_fiyati").Where("id = ? AND borc_durum =1 ", kurbanID).Row()
	row.Scan(&result)
	return result
}

//ToplamOdemeler kurban id ye gore odenen miktar toplamı
func (r *GruplarRepo) ToplamOdemeler(GrupID uint64) float64 {
	var result float64
	row := r.db.Debug().Table(entity.KurbanTableName).Select("sum(bakiye) as total").Where("grup_id = ? AND borc_durum <> 3 AND borc_durum <> 6 AND borc_durum <> 7 AND durum <> 8 AND durum <> 7  AND durum <> 8", GrupID).Row()
	row.Scan(&result)
	return result
}

//ToplamOdemeler kurban id ye gore odenen miktar toplamı
func (r *GruplarRepo) KalanBorclar(GrupID uint64) float64 {
	var result float64
	row := r.db.Debug().Table(entity.KurbanTableName).Select("sum(alacak) as total").Where("grup_id = ? AND borc_durum <> 3 ", GrupID).Row()
	row.Scan(&result)
	return result
}

//ToplamOdemeler kurban id ye gore odenen miktar toplamı
func (r *GruplarRepo) KasaBorcu(GrupID uint64) float64 {
	var result float64
	row := r.db.Debug().Table(entity.KurbanTableName).Select("sum(borc) as total").Where("grup_id = ? AND borc_durum = 3  AND durum <> 6 AND durum <> 7 AND durum <> 8 AND durum <> 7  AND durum <> 8", GrupID).Row()
	row.Scan(&result)
	return result
}

//SatisFiyatTuru kurban id ye gore odenen miktar toplamı
func (r *GruplarRepo) SatisFiyatTuru(hayvanBilgisiID uint64) int {
	var result int
	row := r.db.Debug().Table(entity.GruplarTableName).Select("satis_fiyat_turu").Where("hayvan_bilgisi_id = ? ", hayvanBilgisiID).Row()
	row.Scan(&result)
	return result
}

//Count fat
func (r *GruplarRepo) Count(postTotalCount *int64) {
	var post entity.Gruplar
	var count int64
	r.db.Debug().Model(post).Where(" durum <> 1 ").Count(&count)
	*postTotalCount = count
}

//SetGrupLideri upate data
func (r *GruplarRepo) SetGrupID(grupID uint64, kurbanID uint64, grupLideri int) {
	r.db.Debug().Table(entity.KurbanTableName).Where(" id = ?", kurbanID).Update(entity.Kurban{GrupID: grupID, GrupLideri: grupLideri})
}

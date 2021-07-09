package repository

import (
	"errors"
	"fmt"
	"os"
	"stncCms/app/domain/dto"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stnccollection"
	"stncCms/app/domain/helpers/stnchelper"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var KurbanBorcDurumIlkEklenenFiyat string = stnccollection.IntToString(entity.KurbanBorcDurumIlkEklenenFiyat)

//KurbanRepo struct
type KurbanRepo struct {
	db *gorm.DB
}

//KurbanRepositoryInit initial
func KurbanRepositoryInit(db *gorm.DB) *KurbanRepo {
	return &KurbanRepo{db}
}

//KurbanRepo implements the repository.GenelKurbanRepository interface
// var _ interfaces.PostAppInterface = &KurbanRepo{}

//Save data
func (r *KurbanRepo) Save(post *entity.Kurban) (*entity.Kurban, map[string]string) {
	dbErr := map[string]string{}
	//The images are uploaded to digital ocean spaces. So we need to prepend the url. This might not be your use case, if you are not uploading image to Digital Ocean.

	var err error
	err = r.db.Debug().Create(&post).Error
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
func (r *KurbanRepo) Update(post *entity.Kurban) (*entity.Kurban, map[string]string) {
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

//UpdatePriceList add/new data
func (r *KurbanRepo) UpdatePriceList(post *dto.KurbanUpdateRead) (*dto.KurbanUpdateRead, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Table("kurbanlar").Save(&post).Error
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
func (r *KurbanRepo) Delete(id uint64) error {
	var post entity.Kurban
	var err error
	err = r.db.Debug().Where("id = ?", id).Delete(&post).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

// https://www.mindbowser.com/golang-go-with-gorm/
//  err := Db.Model(&places).Association("town").Find(&places.Town).Error

// https://github.com/stnc-go/gorm_example/blob/master/1-to-1-relationship.go

//GetByID get data
func (r *KurbanRepo) GetByID(id uint64) (*entity.Kurban, error) {
	var kurbandata entity.Kurban
	var err error
	// err = r.db.Debug().Where("id = ?", id).Preload("Odemeler", "borc_durum <> "+KurbanBorcDurumIlkEklenenFiyat).Find(&post).Error
	err = r.db.Debug().Where("id = ?", id).Preload("Odemeler", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}).Find(&kurbandata).Error

	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &kurbandata, nil
}

//GetAll all data
func (r *KurbanRepo) GetAll() ([]entity.Kurban, error) {
	var kurbandata []entity.Kurban
	var err error
	err = r.db.Debug().Order("created_at desc").Find(&kurbandata).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return kurbandata, nil
}

//GetAll all data
func (r *KurbanRepo) GetByGrupID(grupID uint64) ([]entity.Kurban, error) {
	var data []entity.Kurban
	var err error
	err = r.db.Debug().Where("grup_id = ?", grupID).Order("created_at desc").Find(&data).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return data, nil
}

//Count fat
func (r *KurbanRepo) Count(postTotalCount *int64) {
	var kurbandata entity.Kurban
	var count int64
	r.db.Debug().Model(kurbandata).Count(&count)
	*postTotalCount = count
}

//GetAllP pagination all data
func (r *KurbanRepo) GetAllP(postsPerPage int, offset int) ([]entity.Kurban, error) {
	var kurbandata []entity.Kurban
	var err error
	err = r.db.Debug().Limit(postsPerPage).Offset(offset).Order("id asc").Find(&kurbandata).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return kurbandata, nil
}

//TODO: burada daha fazla and gelmeli mi?
//OdemelerToplami kurban id ye gore odenen miktar toplamı
func (r *KurbanRepo) OdemelerToplami(kurbanID uint64) float64 {
	var result float64
	row := r.db.Debug().Table(entity.KurbanTableName).Select("bakiye").Where("id = ?", kurbanID).Row()
	row.Scan(&result)
	return result
}

//KalanUcret kurbanın kalan ücret bilgisi
func (r *KurbanRepo) KalanUcret(kurbanID uint64) float64 {
	var result float64
	row := r.db.Debug().Table(entity.KurbanTableName).Select("alacak").Where("id = ?", kurbanID).Row()
	row.Scan(&result)
	return result
}

//KurbanFiyati kurban id ye gore odenen miktar toplamı
func (r *KurbanRepo) KurbanFiyati(kurbanID uint64) float64 {
	var result float64
	row := r.db.Debug().Table(entity.KurbanTableName).Select("kurban_fiyati").Where("id = ?", kurbanID).Row()
	row.Scan(&result)
	return result
}

//KasaBorcu kurban id ye gore odenen miktar toplamı
func (r *KurbanRepo) KasaBorcu(kurbanID uint64) float64 {
	var result float64
	row := r.db.Debug().Table(entity.KurbanTableName).Select("borc").Where("id = ?", kurbanID).Row()
	row.Scan(&result)
	return result
}

//GetKurbanDurum kurban id ye gore odenen miktar toplamı
func (r *KurbanRepo) GetKurbanDurum(kurbanID uint64) (durum int) {
	row := r.db.Debug().Table(entity.KurbanTableName).Select("durum").Where("id = ? ", kurbanID).Row()
	row.Scan(&durum)
	return durum
}

//GetKurbanHayvanID hisseli kurbana atanmış hayvanın id değerini verir
func (r *KurbanRepo) GetKurbanHayvanID(grupID uint64) (hayvanID int) {
	row := r.db.Debug().Table(entity.GruplarTableName).Select("hayvan_bilgisi_id").Where("id = ? ", grupID).Row()
	row.Scan(&hayvanID)
	return hayvanID
}

//GetKurbanGrupID hisseli kurbana atanmış hayvanın id değerini verir
func (r *KurbanRepo) GetKurbanGrupID(kurbanID uint64) (GrupID int) {
	row := r.db.Debug().Table(entity.KurbanTableName).Select("grup_id").Where("id = ? ", kurbanID).Row()
	row.Scan(&GrupID)
	return GrupID
}

//GetKurbanTuru kurban id ye gore odenen miktar toplamı
func (r *KurbanRepo) GetKurbanTuru(kurbanID uint64) (kurbanTuru int) {
	row := r.db.Debug().Table(entity.KurbanTableName).Select("kurban_turu").Where("id = ? ", kurbanID).Row()
	row.Scan(&kurbanTuru)
	return kurbanTuru
}

//GetGrupLideri kurban id ye gore odenen miktar toplamı
func (r *KurbanRepo) GetGrupLideri(kurbanID uint64) (value int) {
	row := r.db.Debug().Table(entity.KurbanTableName).Select("grup_lideri").Where("id = ? ", kurbanID).Row()
	row.Scan(&value)
	return value
}

//GetKurbanKisiVarmi kurban id ye gore odenen miktar toplamı
func (r *KurbanRepo) GetKurbanKisiVarmi(kurbanID uint64, kisiID *int) {
	var kisi int
	row := r.db.Debug().Table(entity.KurbanTableName).Select("kisi_id").Where("id = ? ", kurbanID).Row()
	row.Scan(&kisi)
	*kisiID = kisi
}

//Count fat
func (r *KurbanRepo) GetByGrupIDCount(grupID uint64, postTotalCount *int) {
	var kurban entity.Kurban
	var count int
	r.db.Debug().Model(kurban).Where("grup_id = ? AND durum <> 3", grupID).Count(&count)
	*postTotalCount = count
}

//SetKurbanBakiyeUpdate upate data
func (r *KurbanRepo) SetKurbanBakiyeUpdate(id uint64, price float64) {
	r.db.Debug().Table(entity.KurbanTableName).Where("id = ?", id).Update("bakiye", price)
}

//SetKurbanKalanUcretiUpdate upate data
func (r *KurbanRepo) SetKurbanKalanUcretiUpdate(id uint64, price float64) {
	r.db.Debug().Table(entity.KurbanTableName).Where("id = ?", id).Update("alacak", price)
}

//SetKurbanKalanUcretiUpdate upate data
func (r *KurbanRepo) SetKurbanKasaBorcuUpdate(id uint64, price float64) {
	r.db.Debug().Table(entity.KurbanTableName).Where("id = ?", id).Update("borc", price)
}

//SetKurbanFiyatiUpdate upate data
func (r *KurbanRepo) SetKurbanFiyatiUpdate(id uint64, price float64) {
	r.db.Debug().Table(entity.KurbanTableName).Where("id = ?", id).Update("kurban_fiyati", price)
}

//KurbanBorcDurumUpdate upate data
func (r *KurbanRepo) SetKurbanBorcDurumUpdate(id uint64, status int) {
	r.db.Debug().Table(entity.KurbanTableName).Where("id = ?", id).Update("borc_durum", status)
}

//SetKurbanDurumUpdate upate data
func (r *KurbanRepo) SetKurbanDurumUpdate(id uint64, status int) {
	r.db.Debug().Table(entity.KurbanTableName).Where("id = ?", id).Update("durum", status)
}

//SetKurbanDurumUpdate upate data
func (r *KurbanRepo) SetKurbanAgirligi(id uint64, agirlik int) {
	r.db.Debug().Table(entity.KurbanTableName).Where("id = ?", id).Update("agirlik", agirlik)
}

//SetGrupLideri upate data
func (r *KurbanRepo) SetGrupLideri(kurbanID uint64, value int) {
	// r.db.Debug().Table(entity.KurbanTableName).Where("grup_id = ?", grupID).Update("grup_lideri", 0)
	r.db.Debug().Table(entity.KurbanTableName).Where(" id = ?", kurbanID).Update("grup_lideri", value)
	// UPDATE kurbanlar SET grup_lideri=15 WHERE grup_id=10 AND id=8
}

//SetGrupLideri upate data
func (r *KurbanRepo) SetVekaletDurumu(kurbanID uint64, value int) {
	r.db.Debug().Table(entity.KurbanTableName).Where(" id = ?", kurbanID).Update("vekalet_durumu", value)
}

//ListPostData upate data
func (r *KurbanRepo) ListDataTable(c *gin.Context) (dto.KurbanBilgisiDataTableResult, error) {
	var total, filtered int
	var err error
	var data []dto.KurbanTable
	// SELECT  kisiler.ad_soyad , kisiler.telefon, kurbanlar.alacak as  ToplamBorc, kurbanlar.borc as KalanBorc, kurbanlar.bakiye as ToplamOdeme ,
	// kurbanlar.id as kurbanid, kisiler.id as kisiID
	// FROM   kurbanlar  INNER JOIN kisiler  ON kurbanlar.kisi_id <> 1

	query := r.db.Table(entity.KurbanTableName)
	// query = query.Select("kisiler.ad_soyad as KisiAdSoyad,kisiler.telefon")
	query = query.Select("kurbanlar.*, kisiler.ad_soyad as kisi_ad_soyad,kisiler.telefon")
	query = query.Joins(" join kisiler on kurbanlar.kisi_id <> 1  and kurbanlar.kisi_id = kisiler.id ")
	// query = query.Joins(" join kisiler on kurbanlar.kisi_id <> 1  ")
	query = query.Offset(stnchelper.QueryOffset(c))
	query = query.Limit(stnchelper.QueryLimit(c))
	//query = query.Order(r.queryOrder(c))
	query = query.Scopes(r.searchScope(c), stnchelper.DateTimeScope(c))
	err = query.Find(&data).Error

	query = query.Offset(0)

	query.Table(entity.KurbanTableName).Count(&filtered)
	// Total data count
	// r.db.Table(entity.KurbanTableName).Count(&total)
	query.Table(entity.KurbanTableName).Count(&total)

	result := dto.KurbanBilgisiDataTableResult{
		Total:    total,
		Filtered: filtered,
		Data:     data,
	}
	return result, err
}

func (r *KurbanRepo) queryOrder(c *gin.Context) string {
	columnMap := map[string]string{
		"0": "kisi_ad_soyad",
		"1": "telefon",
		"2": "bakiye",
		"3": "alacak",
		"4": "borc",
	}

	column := c.DefaultQuery("order[0][column]", "0")
	dir := c.DefaultQuery("order[0][dir]", "desc")
	orderString := columnMap[column] + " " + dir

	return orderString
}

func (r *KurbanRepo) searchScope(c *gin.Context) func(DB *gorm.DB) *gorm.DB {
	return func(DB *gorm.DB) *gorm.DB {
		query := DB
		search := c.QueryMap("search")
		fmt.Println(search)
		dbdriver := os.Getenv("DB_DRIVER")
		// kisiler.ad_soyad LIKE CONCAT('%çürük ali%')
		if dbdriver == "mysql" {
			if search["value"] != "" {
				query = query.Where("kisiler.ad_soyad LIKE (", "%"+search["value"]+"%)")
				query = query.Or("kisiler.telefon  LIKE ? ", "%"+search["value"]+"%")
				query = query.Or("bakiye LIKE ? ", "%"+search["value"]+"%")
				query = query.Or("alacak LIKE ? ", "%"+search["value"]+"%")
				query = query.Or("borc LIKE ? ", "%"+search["value"]+"%")
			}
		} else if dbdriver == "postgres" {
			if search["value"] != "" {
				query = query.Where("kisiler.ad_soyad ILIKE  ? ", "%"+search["value"]+"%")
				query = query.Or("kisiler.telefon  ILIKE ? ", "%"+search["value"]+"%")
				query = query.Or("CAST (bakiye AS TEXT)  ILIKE  ? ", "%"+search["value"]+"%")
				query = query.Or("CAST (alacak AS TEXT)  ILIKE  ? ", "%"+search["value"]+"%")
				query = query.Or("CAST (borc AS TEXT)  ILIKE  ? ", "%"+search["value"]+"%")

			}
		}
		return query
	}
}

//GetAllKurbanAndKisiler
func (r *KurbanRepo) GetAllKurbanAndKisiler(grupId int) ([]dto.KurbanListForGrouplar, error) {
	var kurbandata []dto.KurbanListForGrouplar
	err := r.db.Debug().Table(entity.KurbanTableName).Select(`kurbanlar.id as id ,kurbanlar.user_id,kurbanlar.kisi_id,kurbanlar.aciklama,
	kurbanlar.vekalet_durumu,kurbanlar.agirlik,kurbanlar.grup_lideri,kurbanlar.slug,
	kurbanlar.durum,kurbanlar.borc_durum,kurbanlar.kurban_fiyati,kurbanlar.borc,kurbanlar.alacak,kurbanlar.bakiye,kurbanlar.hayvan_cinsi,
	kisiler.ad_soyad AS kisi_ad_soyad,kisiler.telefon AS kisi_telefon,kisiler.referans_kisi1 AS ref_kisi_id`).Joins("JOIN  kisiler ON kurbanlar.kisi_id=kisiler.id ").Where("grup_id=?", grupId).Order("kurbanlar.grup_lideri asc").Find(&kurbandata).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return kurbandata, nil
}

//GetOdemeRelation
func (r *KurbanRepo) GetKurbanOpenInfo(slug string) (*dto.KurbanOpenInfoData, error) {
	var kurbanOpenInfoData dto.KurbanOpenInfoData
	err := r.db.Debug().Table(entity.KurbanTableName).Select(`
	kurbanlar.grup_id AS kurban_grup_id,kurbanlar.id AS KurbanID,kurbanlar.durum AS kurban_durum,kurbanlar.borc_durum AS kurban_borc_durum,
	kurbanlar.borc AS kurban_borc,kurbanlar.alacak AS kurban_alacak,kurbanlar.bakiye AS kurban_bakiye,kurbanlar.kurban_fiyati,
	gruplar.grup_adi,gruplar.hissedar_adet,gruplar.kesim_sira_no AS kesim_no, gruplar.toplam_kurban_fiyati, gruplar.hayvan_bilgisi_id,
	hayvan_bilgisi.agirlik as hayvan_agirlik,
	kisiler.ad_soyad,kisiler.telefon,kisiler.id as kisi_id`).Joins("JOIN gruplar ON gruplar.id=kurbanlar.grup_id").Joins("JOIN kisiler ON kisiler.id=kurbanlar.kisi_id").Joins("JOIN hayvan_bilgisi ON hayvan_bilgisi.id=gruplar.hayvan_bilgisi_id").Where("kurbanlar.slug=?", slug).Order("kurbanlar.id asc").Find(&kurbanOpenInfoData).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &kurbanOpenInfoData, nil
}

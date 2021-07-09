package repository

import (
	"errors"
	"stncCms/app/domain/dto"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stnccollection"
	"strings"

	"github.com/jinzhu/gorm"
)

var odemelerBorcDurumKasaBorcluDurumda string = stnccollection.IntToString(entity.OdemelerBorcDurumKasaBorcluDurumda)
var ilkEklenenFiyatDegilse string = stnccollection.IntToString(entity.OdemelerBorcDurumIlkEklenenFiyat)

//OdemelerRepo struct
type OdemelerRepo struct {
	db *gorm.DB
}

//OdemelerRepositoryInit initial
func OdemelerRepositoryInit(db *gorm.DB) *OdemelerRepo {
	return &OdemelerRepo{db}
}

//Save data
func (r *OdemelerRepo) Save(odemeler *entity.Odemeler) (*entity.Odemeler, map[string]string) {
	dbErr := map[string]string{}
	var err error
	err = r.db.Debug().Create(&odemeler).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "odemeler title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return odemeler, nil
}

// var dataOdemeler = []entity.Odemeler{} //TODO: bunu nasıl çalıştırablirim
func (r *OdemelerRepo) SaveBatch(odemeler []entity.Odemeler) ([]entity.Odemeler, map[string]string) {
	dbErr := map[string]string{}
	var err error
	err = r.db.Debug().Create(&odemeler).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "odemeler title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return odemeler, nil
}

//Update upate data
func (r *OdemelerRepo) Update(odemeler *entity.Odemeler) (*entity.Odemeler, map[string]string) {
	dbErr := map[string]string{}
	var err error
	err = r.db.Debug().Save(&odemeler).Error
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
	return odemeler, nil
}

//Update upate data
func (r *OdemelerRepo) SetUpdateHesaplar(odemeler *dto.Odemeler) (*dto.Odemeler, map[string]string) {
	dbErr := map[string]string{}
	var err error
	err = r.db.Debug().Save(&odemeler).Error
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
	return odemeler, nil
}

//Delete data
func (r *OdemelerRepo) Delete(id uint64) error {
	var odemeler entity.Odemeler
	var err error
	err = r.db.Debug().Where("id = ?", id).Delete(&odemeler).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//GetByID get data
func (r *OdemelerRepo) GetByID(id uint64) (*entity.Odemeler, error) {
	var odemeler entity.Odemeler
	var err error
	err = r.db.Debug().Where("id = ?", id).Take(&odemeler).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("odemeler not found")
	}
	return &odemeler, nil
}

//GetAll all data
func (r *OdemelerRepo) GetAll() ([]entity.Odemeler, error) {
	var odemelers []entity.Odemeler
	var err error
	err = r.db.Debug().Order("created_at desc").Find(&odemelers).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("odemeler not found")
	}
	return odemelers, nil
}

//TODO: get değeri ekle başına
//LastPrice get data last price
func (r *OdemelerRepo) LastPrice(kurbanID uint64) (*dto.OdemelerSonFiyat, error) {
	var data dto.OdemelerSonFiyat
	var err error
	err = r.db.Debug().Table(dto.OdemelerTableName).Where("kurban_id = ?", kurbanID).Last(&data).Error

	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return &data, nil
}

//TODO: get değeri ekle başına
//KurbanOdenenMiktarForKurbanID kurban id ye gore odenen miktar toplamı
func (r *OdemelerRepo) OdemelerToplami(kurbanID uint64) float64 {
	var result float64
	row := r.db.Debug().Table(dto.OdemelerTableName).Select("sum(verilen_ucret) as total").Where("kurban_id = ? AND borc_durum <> "+odemelerBorcDurumKasaBorcluDurumda+" AND borc_durum <> "+ilkEklenenFiyatDegilse+"", kurbanID).Row()
	row.Scan(&result)
	return result
}

//TODO: burası bekli silineablir decimal falan hesaplamalrı içintutuldu
//KurbanOdenenMiktarForKurbanID kurban id ye gore odenen miktar toplamı
func (r *OdemelerRepo) OdemelerTopl2amiOLDeskiSistemeGoreYenideBakiyeEklendi(kurbanID uint64) float64 {

	/* 	var data dto.OdemelerSonFiyat
	var err error
	err = r.db.Debug().Table(dto.OdemelerTableName).Select("sum(verilen_ucret) as total").Where("kurban_id = ? AND durum <> 3", kurbanID).Last(&data).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return &data, nil
	*/
	//http://learningprogramming.net/golang/gorm/sum-in-gorm/
	var result float64

	row := r.db.Debug().Table(dto.OdemelerTableName).Select("sum(verilen_ucret) as total").Where("kurban_id = ? AND borc_durum <> "+odemelerBorcDurumKasaBorcluDurumda+" AND borc_durum <> "+ilkEklenenFiyatDegilse+"", kurbanID).Row()
	row.Scan(&result)
	/*

	   type Decimal struct {
	   	value *big.Int
	   	exp   int32
	   }
	*/
	// price, err := decimal.NewFromString(result)
	// if err != nil {
	// 	panic(err)
	// }

	return result

}

//TODO: get değeri ekle başına
//KurbanSonKalanUcret kurban id ye gore odenen miktar toplamı
func (r *OdemelerRepo) KurbanSonKalanUcret(kurbanID uint64) (result float64) {
	row := r.db.Debug().Table(dto.OdemelerTableName).Select("alacak").Where("kurban_id = ?  AND borc_durum <> "+odemelerBorcDurumKasaBorcluDurumda+" ", kurbanID).Order("id desc").Row()
	row.Scan(&result)
	return result
}

//TODO: get değeri ekle başına
//KasaBorcu kasa borcunu verir ** UPDATE `teknopark`.`genel_kurban` SET `kurban_fiyati`='600' ,durum=1 WHERE  `id`=10;
func (r *OdemelerRepo) KasaBorcu(kurbanID uint64) float64 {
	var result float64
	row := r.db.Debug().Table(dto.OdemelerTableName).Select("alacak").Where("kurban_id = ?  AND borc_durum ="+odemelerBorcDurumKasaBorcluDurumda+"", kurbanID).Order("id desc").Row()
	row.Scan(&result)
	return result
}

//GetOdemeRelation
func (r *OdemelerRepo) GetOdemeRelation(odemeID uint64) (*dto.OdemeMakbuzu, error) {
	var odemeData dto.OdemeMakbuzu
	err := r.db.Debug().Table(entity.OdemelerTableName).Select(`odemeler.id AS odeme_id,odemeler.aciklama,odemeler.makbuz,odemeler.durum AS odeme_durum,odemeler.borc_durum AS odemeler_borc_durum,
	odemeler.kurban_fiyati AS odemeler_kurban_fiyati,
	odemeler.verilen_ucret AS odemeler_verilen_ucret,odemeler.verildigi_tarih AS odemeler_verildigi_tarih,
	kurbanlar.grup_id AS kurban_grup_id,kurbanlar.id AS KurbanID,kurbanlar.durum AS kurban_durum,kurbanlar.borc_durum AS kurban_borc_durum,
	kurbanlar.borc AS kurban_borc,kurbanlar.alacak AS kurban_alacak,kurbanlar.bakiye AS kurban_bakiye,
	gruplar.grup_adi,gruplar.hissedar_adet,gruplar.kesim_sira_no as kesim_no,
	kisiler.ad_soyad,kisiler.telefon,kisiler.id as kisi_id`).Joins("JOIN kurbanlar ON odemeler.kurban_id =kurbanlar.id ").Joins("JOIN gruplar ON gruplar.id=kurbanlar.grup_id").Joins("JOIN kisiler ON kisiler.id=kurbanlar.kisi_id").Where("odemeler.id=?", odemeID).Order("odemeler.id asc").Find(&odemeData).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &odemeData, nil
}

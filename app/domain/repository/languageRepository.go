package repository

import (
	"errors"
	"stncCms/app/domain/entity"
	"strings"

	"github.com/jinzhu/gorm"
)

//LanguageRepo struct
type LanguageRepo struct {
	db *gorm.DB
}

//LanguageRepositoryInit initial
func LanguageRepositoryInit(db *gorm.DB) *LanguageRepo {
	return &LanguageRepo{db}
}

//languageRepo implements the repository.languageRepository interface
// var _ interfaces.languageAppInterface = &languageRepo{}

//Save data
func (r *LanguageRepo) Save(language *entity.Languages) (*entity.Languages, map[string]string) {
	dbErr := map[string]string{}
	var err error
	err = r.db.Debug().Create(&language).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "language title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return language, nil
}

//Update upate data
func (r *LanguageRepo) Update(language *entity.Languages) (*entity.Languages, map[string]string) {
	dbErr := map[string]string{}

	err := r.db.Debug().Save(&language).Error
	//db.Table("libraries").Where("id = ?", id).Update(languageData)

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
	return language, nil
}

//Delete data
func (r *LanguageRepo) Delete(id uint64) error {
	var language entity.Languages
	var err error
	err = r.db.Debug().Where("id = ?", id).Delete(&language).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//GetByID get data
func (r *LanguageRepo) GetByID(id uint64) (*entity.Languages, error) {
	var language entity.Languages
	var err error
	err = r.db.Debug().Where("id = ?", id).Take(&language).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("language not found")
	}
	return &language, nil
}

//GetAll all data
func (r *LanguageRepo) GetAll() ([]entity.Languages, error) {
	var languages []entity.Languages
	var err error
	err = r.db.Debug().Order("created_at desc").Find(&languages).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("language not found")
	}
	return languages, nil
}

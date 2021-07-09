package repository

import (
	"errors"
	"stncCms/app/domain/entity"
	"strings"

	"github.com/jinzhu/gorm"
)

//WebArchiveRepo struct
type WebArchiveRepo struct {
	db *gorm.DB
}

//WebArchiveRepositoryInit initial
func WebArchiveRepositoryInit(db *gorm.DB) *WebArchiveRepo {
	return &WebArchiveRepo{db}
}

//WebArchiveRepo implements the repository.WebArchiveRepository interface
// var _ interfaces.PostAppInterface = &WebArchiveRepo{}

//Save data
func (r *WebArchiveRepo) Save(data *entity.WebArchive) (*entity.WebArchive, map[string]string) {
	dbErr := map[string]string{}
	var err error
	err = r.db.Debug().Create(&data).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "webarc title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return data, nil
}

//Update upate data
func (r *WebArchiveRepo) Update(data *entity.WebArchive) (*entity.WebArchive, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&data).Error
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
	return data, nil
}

//Count fat
func (r *WebArchiveRepo) Count(totalCount *int64) {
	var data entity.WebArchive
	var count int64
	r.db.Debug().Model(data).Count(&count)
	*totalCount = count
}

//Delete data
func (r *WebArchiveRepo) Delete(id uint64) error {
	var data entity.WebArchive
	var err error
	err = r.db.Debug().Where("id = ?", id).Unscoped().Delete(&data).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//GetByID get data
func (r *WebArchiveRepo) GetByID(id uint64) (*entity.WebArchive, error) {
	var data entity.WebArchive
	var err error
	// err = r.db.Debug().Where("id = ?", id).Take(&data).Error
	err = r.db.Debug().Where("id = ?", id).Preload("WebArchiveLinks").Find(&data).Error

	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("webarc not found")
	}
	return &data, nil
}

//GetAll all data
func (r *WebArchiveRepo) GetAll() ([]entity.WebArchive, error) {
	var data []entity.WebArchive
	var err error
	err = r.db.Debug().Order("created_at desc").Find(&data).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("webarc not found")
	}
	return data, nil
}

//GetAllP pagination all data
func (r *WebArchiveRepo) GetAllP(perPage int, offset int) ([]entity.WebArchive, error) {
	var data []entity.WebArchive
	var err error
	err = r.db.Debug().Limit(perPage).Offset(offset).Order("created_at desc").Find(&data).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("webarc not found")
	}
	return data, nil
}

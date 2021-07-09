package repository

import (
	"errors"
	"stncCms/app/domain/entity"
	"strings"

	"github.com/jinzhu/gorm"
)

//WebArchiveLinksRepo struct
type WebArchiveLinksRepo struct {
	db *gorm.DB
}

//WebArchiveLinksRepositoryInit initial
func WebArchiveLinksRepositoryInit(db *gorm.DB) *WebArchiveLinksRepo {
	return &WebArchiveLinksRepo{db}
}

//WebArchiveLinksRepo implements the repository.WebArchiveRepository interface
// var _ interfaces.PostAppInterface = &WebArchiveLinksRepo{}

//Save data
func (r *WebArchiveLinksRepo) Save(data *entity.WebArchiveLinks) (*entity.WebArchiveLinks, map[string]string) {
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
func (r *WebArchiveLinksRepo) Update(data *entity.WebArchiveLinks) (*entity.WebArchiveLinks, map[string]string) {
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
func (r *WebArchiveLinksRepo) Count(totalCount *int64) {
	var data entity.WebArchiveLinks
	var count int64
	r.db.Debug().Model(data).Count(&count)
	*totalCount = count
}

//Delete data
func (r *WebArchiveLinksRepo) Delete(id uint64) error {
	var data entity.WebArchiveLinks
	var err error
	err = r.db.Debug().Where("id = ?", id).Delete(&data).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//GetByID get data
func (r *WebArchiveLinksRepo) GetByID(id uint64) (*entity.WebArchiveLinks, error) {
	var data entity.WebArchiveLinks
	var err error
	err = r.db.Debug().Where("id = ?", id).Take(&data).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("webarc not found")
	}
	return &data, nil
}

//GetAll all data
func (r *WebArchiveLinksRepo) GetAll() ([]entity.WebArchiveLinks, error) {
	var data []entity.WebArchiveLinks
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
func (r *WebArchiveLinksRepo) GetAllP(perPage int, offset int) ([]entity.WebArchiveLinks, error) {
	var data []entity.WebArchiveLinks
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

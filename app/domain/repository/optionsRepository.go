package repository

import (
	"errors"
	"stncCms/app/domain/entity"
	"strings"

	"github.com/jinzhu/gorm"
)

var optionTableName string = "options"

//OptionRepositoryRepo struct
type OptionRepositoryRepo struct {
	db *gorm.DB
}

//OptionRepositoryInit initial
func OptionRepositoryInit(db *gorm.DB) *OptionRepositoryRepo {
	return &OptionRepositoryRepo{db}
}

//AddOption data (save)
func (r *OptionRepositoryRepo) AddOption(data *entity.Options) (*entity.Options, map[string]string) {
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

//GetOptionID get data
func (r *OptionRepositoryRepo) GetOptionID(name string) (returnValue int) {
	row := r.db.Debug().Table(optionTableName).Select("option_id").Where("option_name = ?", name).Row()
	row.Scan(&returnValue)
	return returnValue
}

//GetOption get data
func (r *OptionRepositoryRepo) GetOption(name string) string {
	var result string
	row := r.db.Debug().Table(optionTableName).Select("option_value").Where("option_name = ?", name).Row()
	row.Scan(&result)

	return result
}

//SetOption upate data
func (r *OptionRepositoryRepo) SetOption(name, value string) {
	id := r.GetOptionID(name)
	r.db.Debug().Table(optionTableName).Where("option_id = ?", id).Update("option_value", value)
	// var err error
	// dbErr := map[string]string{}
	// if err != nil {
	// 	//since our title is unique
	// 	if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
	// 		dbErr["unique_title"] = "title already taken"
	// 		return nil, dbErr
	// 	}
	// 	//any other db error
	// 	dbErr["db_error"] = "database error"
	// 	return nil, dbErr
	// }
	// return data, nil
}

//DeleteOptionID data
func (r *OptionRepositoryRepo) DeleteOptionID(id uint64) error {
	var data entity.Options
	var err error
	err = r.db.Debug().Where("id = ?", id).Unscoped().Delete(&data).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//DeleteOption data
func (r *OptionRepositoryRepo) DeleteOption(value string) error {
	var data entity.Options
	var err error
	err = r.db.Debug().Where("option_name = ?", value).Unscoped().Delete(&data).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//GetAll all data
func (r *OptionRepositoryRepo) GetAll() ([]entity.Options, error) {
	var data []entity.Options
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

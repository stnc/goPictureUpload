package repository

import (
	"errors"
	"stncCms/app/domain/entity"
	"strings"

	"github.com/jinzhu/gorm"
)

//MediaRepo struct
type MediaRepo struct {
	db *gorm.DB
}

//MediaRepositoryInit initial
func MediaRepositoryInit(db *gorm.DB) *MediaRepo {
	return &MediaRepo{db}
}

//MediaRepo implements the repository.MediaRepository interface
// var _ interfaces.MediaAppInterface = &MediaRepo{}

//Save data
func (r *MediaRepo) Save(Media *entity.Media) (*entity.Media, map[string]string) {
	dbErr := map[string]string{}
	var err error
	err = r.db.Debug().Create(&Media).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "Media title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return Media, nil
}

//Update upate data
func (r *MediaRepo) Update(Media *entity.Media) (*entity.Media, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&Media).Error
	//db.Table("libraries").Where("id = ?", id).Update(MediaData)

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
	return Media, nil
}

//Count fat
func (r *MediaRepo) Count(modulID int, contentID int, MediaTotalCount *int) {
	var Media entity.Media
	var count int
	r.db.Debug().Model(Media).Where("modul_id = ? AND content_id = ?", modulID, contentID).Count(&count)
	*MediaTotalCount = count
}

//Delete data
func (r *MediaRepo) Delete(id uint64) error {

	var err error
	// var Media entity.Media                                      //#soft delete
	// err = r.db.Debug().Where("id = ?", id).Delete(&Media).Error //#soft delete
	err = r.db.Debug().Unscoped().Where("id = ?", id).Delete(entity.Media{}).Error //#hard delete

	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//GetByID get data
func (r *MediaRepo) GetByID(id uint64) (*entity.Media, error) {
	var Media entity.Media
	var err error
	err = r.db.Debug().Where("id = ?", id).Take(&Media).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("Media not found")
	}
	return &Media, nil
}

//GetAll all data
func (r *MediaRepo) GetAll(modulID int, contentID int) ([]entity.Media, error) {
	var Medias []entity.Media
	var err error
	err = r.db.Debug().Where("modul_id = ? AND content_id = ?", modulID, contentID).Order("created_at desc").Find(&Medias).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("Media not found")
	}
	return Medias, nil
}

//GetAllforModul all data
func (r *MediaRepo) GetAllforModul(modulID int, contentID int) ([]entity.Media, error) {
	var Medias []entity.Media
	var err error
	err = r.db.Debug().Where("modul_id = ? AND content_id = ?", modulID, contentID).Order("created_at desc").Find(&Medias).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("Media not found")
	}
	return Medias, nil
}

//GetAllP pagination all data
func (r *MediaRepo) GetAllP(MediasPerPage int, offset int) ([]entity.Media, error) {
	var Medias []entity.Media
	var err error
	err = r.db.Debug().Limit(MediasPerPage).Offset(offset).Order("created_at desc").Find(&Medias).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("Media not found")
	}
	return Medias, nil
}

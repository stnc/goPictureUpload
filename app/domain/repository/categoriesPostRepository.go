package repository

import (
	"errors"
	"stncCms/app/domain/entity"
	"strings"

	"github.com/jinzhu/gorm"

)

//CatPostRepo struct
type CatPostRepo struct {
	db *gorm.DB
}

//CatPostRepositoryInit initial
func CatPostRepositoryInit(db *gorm.DB) *CatPostRepo {
	return &CatPostRepo{db}
}

//PostRepo implements the repository.PostRepository interface
// var _ interfaces.CatPostAppInterface = &CatPostRepo{}

//Save data
func (r *CatPostRepo) Save(cat *entity.CategoryPosts) (*entity.CategoryPosts, map[string]string) {
	dbErr := map[string]string{}

	err := r.db.Debug().Create(&cat).Error
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
	return cat, nil
}

//GetByID get data
func (r *CatPostRepo) GetByID(id uint64) (*entity.CategoryPosts, error) {
	var cat entity.CategoryPosts
	err := r.db.Debug().Where("id = ?", id).Take(&cat).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	// if gorm.IsRecordNotFoundError(err) {
	// 	return nil, errors.New("post not found")

	// }
	return &cat, nil
}

//GetAllforCatID get data
func (r *CatPostRepo) GetAllforCatID(catid uint64) ([]entity.CategoryPosts, error) {
	var cat []entity.CategoryPosts
	err := r.db.Debug().Limit(100).Where("category_id = ?", catid).Order("created_at desc").Find(&cat).Error
	if err != nil {
		return nil, err
	}

	// if err.IsRecordNotFoundError(err) {
	// 	return nil, errors.New("post not found")
	// }
	return cat, nil
}

//GetAllforPostID all data
func (r *CatPostRepo) GetAllforPostID(postid uint64) ([]entity.CategoryPosts, error) {
	var cat []entity.CategoryPosts
	err := r.db.Debug().Limit(100).Where("post_id = ?", postid).Order("created_at desc").Find(&cat).Error
	if err != nil {
		return nil, err
	}
	// if gorm.IsRecordNotFoundError(err) {
	// 	return nil, errors.New("post not found")
	// }
	return cat, nil
}

//GetAll all data
func (r *CatPostRepo) GetAll() ([]entity.CategoryPosts, error) {
	var cat []entity.CategoryPosts
	err := r.db.Debug().Limit(100).Order("created_at desc").Find(&cat).Error
	if err != nil {
		return nil, err
	}
	// if gorm.IsRecordNotFoundError(err) {
	// 	return nil, errors.New("post not found")
	// }
	return cat, nil
}

//Update upate data
func (r *CatPostRepo) Update(cat *entity.CategoryPosts) (*entity.CategoryPosts, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&cat).Error

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
	return cat, nil
}

//Delete delete data
func (r *CatPostRepo) Delete(id uint64) error {
	var cat entity.CategoryPosts
	err := r.db.Debug().Where("id = ?", id).Unscoped().Delete(&cat).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//DeleteForPostID delete data
func (r *CatPostRepo) DeleteForPostID(postID uint64) error {
	var cat entity.CategoryPosts
	err := r.db.Debug().Where("post_id = ?", postID).Unscoped().Delete(&cat).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//DeleteForCatID delete data
func (r *CatPostRepo) DeleteForCatID(CatID uint64) error {
	var cat entity.CategoryPosts
	err := r.db.Debug().Where("category_id = ?", CatID).Unscoped().Delete(&cat).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

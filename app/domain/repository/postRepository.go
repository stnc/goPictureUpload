package repository

import (
	"errors"
	"os"
	"stncCms/app/domain/entity"
	"strings"

	"github.com/jinzhu/gorm"
)

//PostRepo struct
type PostRepo struct {
	db *gorm.DB
}

//PostRepositoryInit initial
func PostRepositoryInit(db *gorm.DB) *PostRepo {
	return &PostRepo{db}
}

//PostRepo implements the repository.PostRepository interface
// var _ interfaces.PostAppInterface = &PostRepo{}

//Save data
func (r *PostRepo) Save(post *entity.Post) (*entity.Post, map[string]string) {
	dbErr := map[string]string{}
	//The images are uploaded to digital ocean spaces. So we need to prepend the url. This might not be your use case, if you are not uploading image to Digital Ocean.
	post.Picture = os.Getenv("DO_SPACES_URL") + post.Picture
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
func (r *PostRepo) Update(post *entity.Post) (*entity.Post, map[string]string) {
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

//Count fat
func (r *PostRepo) Count(postTotalCount *int64) {
	var post entity.Post
	var count int64
	r.db.Debug().Model(post).Count(&count)
	*postTotalCount = count
}

//Delete data
func (r *PostRepo) Delete(id uint64) error {
	var post entity.Post
	var err error
	err = r.db.Debug().Where("id = ?", id).Delete(&post).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

//GetByID get data
func (r *PostRepo) GetByID(id uint64) (*entity.Post, error) {
	var post entity.Post
	var err error
	err = r.db.Debug().Where("id = ?", id).Take(&post).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &post, nil
}

//GetAll all data
func (r *PostRepo) GetAll() ([]entity.Post, error) {
	var posts []entity.Post
	var err error
	err = r.db.Debug().Order("created_at desc").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return posts, nil
}

//GetAllP pagination all data
func (r *PostRepo) GetAllP(postsPerPage int, offset int) ([]entity.Post, error) {
	var posts []entity.Post
	var err error
	err = r.db.Debug().Limit(postsPerPage).Offset(offset).Order("created_at desc").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return posts, nil
}

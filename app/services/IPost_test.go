package services


/*
import (
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stnchelper"
	"stncCms/app/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

//IF YOU HAVE TIME, YOU CAN TEST ALL THE METHODS FAILURES

type fakePostRepo struct{}

var (
	savePostRepo   func(*entity.Post) (*entity.Post, map[string]string)
	getPostRepo    func(uint64) (*entity.Post, error)
	getAllPostRepo func() ([]entity.Post, error)
	updatePostRepo func(*entity.Post) (*entity.Post, map[string]string)
	deletePostRepo func(uint64) error
)

func (f *fakePostRepo) Save(post *entity.Post) (*entity.Post, map[string]string) {
	return savePostRepo(post)
}
func (f *fakePostRepo) GetByID(postId uint64) (*entity.Post, error) {
	return getPostRepo(postId)
}
func (f *fakePostRepo) GetAll() ([]entity.Post, error) {
	return getAllPostRepo()
}
func (f *fakePostRepo) Update(post *entity.Post) (*entity.Post, map[string]string) {
	return updatePostRepo(post)
}
func (f *fakePostRepo) Delete(postId uint64) error {
	return deletePostRepo(postId)
}

//var fakePost repository.PostRepository = &fakePostRepo{} //this is where the real implementation is swap with our fake implementation
var postAppFake services.PostAppInterface = &fakePostRepo{} //this is where the real implementation is swap with our fake implementation

func TestSavePost_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	savePostRepo = func(user *entity.Post) (*entity.Post, map[string]string) {
		return &entity.Post{
			ID:          1,
			PostTitle:   "post title",
			PostContent: "post description",
			PostExcerpt: "post exceprt",
			PostSlug:    stnchelper.Slugify("post title", 15),
			PostStatus:  1,
			PostType:    1,
			UserID:      1,
		}, nil
	}
	post := &entity.Post{
		ID:          1,
		PostTitle:   "post title",
		PostContent: "post description",
		PostExcerpt: "post exceprt",
		PostSlug:    stnchelper.Slugify("post title", 15),
		PostStatus:  1,
		PostType:    1,
		UserID:      1,
	}
	f, err := postAppFake.Save(post)
	assert.Nil(t, err)
	assert.EqualValues(t, f.PostTitle, "post title")
	assert.EqualValues(t, f.PostContent, "post description")
	assert.EqualValues(t, f.PostExcerpt, "post exceprt")
	assert.EqualValues(t, f.PostSlug, "post-title")
	assert.EqualValues(t, f.PostStatus, 1)
	assert.EqualValues(t, f.PostType, 1)
	assert.EqualValues(t, f.UserID, 1)
}

func TestGetPost_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getPostRepo = func(postID uint64) (*entity.Post, error) {
		return &entity.Post{
			ID:          1,
			PostTitle:   "post title",
			PostContent: "post description",
			UserID:      1,
		}, nil
	}
	postID := uint64(1)
	f, err := postAppFake.GetByID(postID)
	assert.Nil(t, err)
	assert.EqualValues(t, f.PostTitle, "post title")
	assert.EqualValues(t, f.PostContent, "post description")
	assert.EqualValues(t, f.UserID, 1)
}

func TestAllPost_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getAllPostRepo = func() ([]entity.Post, error) {
		return []entity.Post{
			{
				ID:          1,
				PostTitle:   "post title first",
				PostContent: "post description first",
				UserID:      1,
			},
			{
				ID:          2,
				PostTitle:   "post title second",
				PostContent: "post description second",
				UserID:      1,
			},
		}, nil
	}
	f, err := postAppFake.GetAll()
	assert.Nil(t, err)
	assert.EqualValues(t, len(f), 2)
}

func TestUpdatePost_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	updatePostRepo = func(user *entity.Post) (*entity.Post, map[string]string) {
		return &entity.Post{
			ID:          1,
			PostTitle:   "post title update",
			PostContent: "post description update",
			UserID:      1,
		}, nil
	}
	post := &entity.Post{
		ID:          1,
		PostTitle:   "post title update",
		PostContent: "post description update",
		UserID:      1,
	}
	f, err := postAppFake.Update(post)
	assert.Nil(t, err)
	assert.EqualValues(t, f.PostTitle, "post title update")
	assert.EqualValues(t, f.PostContent, "post description update")
	assert.EqualValues(t, f.UserID, 1)
}

func TestDeletePost_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	deletePostRepo = func(postID uint64) error {
		return nil
	}
	postID := uint64(1)
	err := postAppFake.Delete(postID)
	assert.Nil(t, err)
}
*/

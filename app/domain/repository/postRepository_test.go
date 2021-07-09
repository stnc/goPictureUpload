package repository

import (
	"stncCms/app/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSavePost_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var post = entity.Post{}
	post.PostTitle = "post title"
	post.PostContent = "post description"
	post.UserID = 1 //stncsession.GetUserID2(c)

	repo := PostRepositoryInit(conn)

	f, saveErr := repo.Save(&post)
	assert.Nil(t, saveErr)
	assert.EqualValues(t, f.PostTitle, "post title")
	assert.EqualValues(t, f.PostContent, "post description")
	assert.EqualValues(t, f.UserID, 1)
}

//Failure can be due to duplicate email, etc
//Here, we will attempt saving a post that is already saved
func TestSavePost_Failure(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the post
	_, err = seedPost(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var post = entity.Post{}
	post.PostTitle = "post title"
	post.PostContent = "post desc"
	post.UserID = 1

	repo := PostRepositoryInit(conn)
	f, saveErr := repo.Save(&post)

	dbMsg := map[string]string{
		"unique_title": "post title already taken",
	}
	assert.Nil(t, f)
	assert.EqualValues(t, dbMsg, saveErr)
}

func TestGetPostID_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	post, err := seedPost(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := PostRepositoryInit(conn)

	f, saveErr := repo.GetByID(post.ID)

	assert.Nil(t, saveErr)
	assert.EqualValues(t, f.PostTitle, post.PostTitle)
	assert.EqualValues(t, f.PostContent, post.PostContent)
	assert.EqualValues(t, f.UserID, post.UserID)
}

func TestGetAllPost_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	_, err = seedPosts(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := PostRepositoryInit(conn)
	posts, getErr := repo.GetAll()

	assert.Nil(t, getErr)
	assert.EqualValues(t, len(posts), 2)
}

func TestUpdatePost_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	post, err := seedPost(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//updating
	post.PostTitle = "post title update"
	post.PostContent = "post description update"

	repo := PostRepositoryInit(conn)
	f, updateErr := repo.Update(post)

	assert.Nil(t, updateErr)
	assert.EqualValues(t, f.ID, 1)
	assert.EqualValues(t, f.PostTitle, "post title update")
	assert.EqualValues(t, f.PostContent, "post description update")
	assert.EqualValues(t, f.UserID, 1)
}

//Duplicate title error
func TestUpdatePost_Failure(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	posts, err := seedPosts(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var secondPost entity.Post

	//get the second post title
	for _, v := range posts {
		if v.ID == 1 {
			continue
		}
		secondPost = v
	}
	secondPost.PostTitle = "first post" //this title belongs to the first post already, so the second post cannot use it
	secondPost.PostContent = "New description"

	repo := PostRepositoryInit(conn)
	f, updateErr := repo.Update(&secondPost)

	dbMsg := map[string]string{
		"unique_title": "title already taken",
	}
	assert.NotNil(t, updateErr)
	assert.Nil(t, f)
	assert.EqualValues(t, dbMsg, updateErr)
}

func TestDeletePost_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	post, err := seedPost(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := PostRepositoryInit(conn)

	deleteErr := repo.Delete(post.ID)

	assert.Nil(t, deleteErr)
}

package controller

import (
	"fmt"
	"net/http"
	"os"
	"stncCms/app/domain/entity"
	"stncCms/app/infrastructure/auth"
	"stncCms/app/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//Post constructor
type Post struct {
	postApp services.PostAppInterface
	userApp services.UserAppInterface
	tk      auth.TokenInterface
	rd      auth.AuthInterface
}

//InitPost post controller  constructor
func InitPost(fApp services.PostAppInterface, uApp services.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Post {
	return &Post{
		postApp: fApp,
		userApp: uApp,
		rd:      rd,
		tk:      tk,
	}
}

//SavePost save method
func (fo *Post) SavePost(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userID, err := fo.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var savePostError = make(map[string]string)
	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}
	//We initialize a new post for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyPost := entity.Post{}
	emptyPost.PostTitle = title
	emptyPost.PostContent = description
	savePostError = emptyPost.Validate()
	if len(savePostError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, savePostError)
		return
	}

	//check if the user exist
	_, err = fo.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	var post = entity.Post{}
	post.UserID = userID
	post.PostTitle = title
	post.PostContent = description
	// post.Image = uploadedFile
	post.Picture = ""
	savedPost, saveErr := fo.postApp.Save(&post)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, savedPost)
}

//UpdatePost update method
func (fo *Post) UpdatePost(c *gin.Context) {
	//Check if the user is authenticated first
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	//lookup the metadata in redis:
	userID, err := fo.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var updatePostError = make(map[string]string)

	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	//Since it is a multipart form data we sent, we will do a manual check on each item
	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
	}
	//We initialize a new post for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyPost := entity.Post{}
	emptyPost.PostTitle = title
	emptyPost.PostContent = description
	updatePostError = emptyPost.Validate()
	if len(updatePostError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, updatePostError)
		return
	}
	user, err := fo.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	//check if the post exist:
	post, err := fo.postApp.GetByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	//if the user id doesnt match with the one we have, dont update. This is the case where an authenticated user tries to update someone else post using postman, curl, etc
	if user.ID != post.UserID {
		c.JSON(http.StatusUnauthorized, "you are not the owner of this post")
		return
	}
	//Since this is an update request,  a new image may or may not be given.
	// If not image is given, an error occurs. We know this that is why we ignored the error and instead check if the file is nil.
	// if not nil, we process the file by calling the "UploadFile" method.
	// if nil, we used the old one whose path is saved in the database
	file, _ := c.FormFile("image")
	if file != nil {
		// post.Picture, err = fo.fileUpload.UploadFile(file)
		//since i am using Digital Ocean(DO) Spaces to save image, i am appending my DO url here. You can comment this line since you may be using Digital Ocean Spaces.
		post.Picture = os.Getenv("DO_SPACES_URL") + post.Picture
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"upload_err": err.Error(),
			})
			return
		}
	}
	//we dont need to update user's id
	post.PostTitle = title
	post.PostContent = description
	post.UpdatedAt = time.Now()
	updatedPost, dbUpdateErr := fo.postApp.Update(post)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}
	c.JSON(http.StatusOK, updatedPost)
}

//GetAllPost all
func (fo *Post) GetAllPost(c *gin.Context) {
	allpost, err := fo.postApp.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allpost)
}

//GetPostAndCreator crate
func (fo *Post) GetPostAndCreator(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	post, err := fo.postApp.GetByID(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	user, err := fo.userApp.GetUser(post.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	postAndUser := map[string]interface{}{
		"post":    post,
		"creator": user.PublicUser(),
	}
	c.JSON(http.StatusOK, postAndUser)
}

//DeletePost delete event
func (fo *Post) DeletePost(c *gin.Context) {
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	_, err = fo.userApp.GetUser(metadata.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = fo.postApp.Delete(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "post deleted")
}

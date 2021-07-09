package controller

import (
	"fmt"
	"net/http"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stncsession"
	"stncCms/app/infrastructure/security"
	"stncCms/app/services"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//Login constructor
type Login struct {
	userApp services.UserAppInterface
}

//InitLogin login controller constructor
func InitLogin(uApp services.UserAppInterface) *Login {
	return &Login{
		userApp: uApp,
	}
}

//Login func implement
func (l *Login) Login(c *gin.Context) {
	flashMsg := stncsession.GetFlashMessage(c)
	viewData := pongo2.Context{
		"paginator": paginator,
		"title":     "Giriş",
		"flashMsg":  flashMsg,
		"csrf":      csrf.GetToken(c),
	}

	c.HTML(
		http.StatusOK,
		"admin/login/login.html",
		viewData,
	)
}

//sifre?p=mu1tluerF9E&name=mutluer33@gmail.com
//sifre?p=ku1rtulusTF9&name=kurtulusiplikci@yahoo.com.tr
//sifre?p=su1salGCY&name=susalyurdu38@gmail.com
//sifre?p=ha1mzaGX5&name=hamzakamil@gmail.com

func (l *Login) SifreVer(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	hak := c.Query("p")
	name := c.Query("name")
	sifre := security.PassVer(hak)
	var post entity.User
	post.Email = name
	post.Password = sifre
	l.userApp.SaveUser(&post)
	c.JSON(http.StatusOK, sifre)

}

//Login func implement
func (l *Login) LoginPost(c *gin.Context) {
	var user = entity.User{}
	flashMsg := stncsession.GetFlashMessage(c)
	var savePostError = make(map[string]string)

	email := c.PostForm("Email")
	pass := c.PostForm("Password") //"111111-6" //sha1 hali cb5e6834e30cf762b38387db44c936daac667559
	user.Email = email
	user.Password = pass

	validateUser := user.Validate("login")
	if len(validateUser) > 0 {
		//	c.JSON(http.StatusUnprocessableEntity, validateUser)
		//stncsession.SetFlashMessage(validateUser, c)
		savePostError = validateUser

	} else {
		userData, userErr := l.userApp.GetUserByEmailAndPassword2(email, pass)
		if userErr != nil {
			savePostError = userErr
			stncsession.SetFlashMessage("Kullanıcı adı veya şifre hatalıdır", "success", c)
			c.Redirect(http.StatusMovedPermanently, "/login")
			return
		} else {
			stncsession.SetStoreUserID(c, userData.ID)
			stncsession.SetSession("UserName", userData.FirstName, c)
			fmt.Println(userData)
			//	c.SetCookie("username", "selmnn", 3600, "", "", false, true)
			c.Redirect(http.StatusMovedPermanently, "/")
			return
		}

	}

	viewData := pongo2.Context{
		"paginator": paginator,
		"err":       savePostError,
		"title":     "Giriş",
		//"posts":     userData,
		"flashMsg": flashMsg,
		"email":    email,
		"password": pass,
		"csrf":     csrf.GetToken(c),
	}

	c.HTML(
		http.StatusOK,
		"admin/login/login.html",
		viewData,
	)
}

//Login func implement
func (l *Login) LoginAPI(c *gin.Context) {
	var user = entity.User{}
	email := "selmantunc@gmail.com"
	pass := "111111-6"
	user.Email = email
	user.Password = pass
	hashPassword := security.Hash(pass)
	fmt.Println("selman: " + string(hashPassword))
	//validate request:
	//var user *entity.User
	validateUser := user.Validate("login")
	if len(validateUser) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateUser)
		return
	}
	userData, userErr := l.userApp.GetUserByEmailAndPassword2(email, pass)
	if userErr != nil {
		c.JSON(http.StatusInternalServerError, userErr)
		return
	} else {
		stncsession.SetStoreUserID(c, userData.ID)
		stncsession.SetSession("UserName", userData.FirstName, c)
		//	c.SetCookie("username", "selmnn", 3600, "", "", false, true)
	}
	fmt.Println(userData)
	c.Redirect(http.StatusMovedPermanently, "/")
	c.JSON(http.StatusOK, userData)
}

//Logout güvenli çıkış
func (au *Login) Logout(c *gin.Context) {
	stncsession.ClearUserID(c)
	c.Redirect(http.StatusTemporaryRedirect, "/login")

	//c.JSON(http.StatusOK, u)
}

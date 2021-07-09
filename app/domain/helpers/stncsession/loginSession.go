package stncsession

import (
	"fmt"
	"net/http"
	"stncCms/app/domain/helpers/stnccollection"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// IsLoggedIn checks if the current user is logged in
func IsLoggedIn(c *gin.Context) bool {
	session := sessions.Default(c)
	return session.Get("user2ID") != nil
}

//IsLoggedInRedirect redirect site
func IsLoggedInRedirect(c *gin.Context) {
	//#stnc-notes
	if l := IsLoggedIn(c); !l {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

//GetUserID get userID data
func GetUserID(c *gin.Context) (userID string) {
	_userID := GetSession("user2ID", c)
	if _userID == "" {
		userID = "null"
	} else {
		userID = _userID
		fmt.Println("data " + userID)
	}
	return userID
}

//GetUserID get userID data
func GetUserID2(c *gin.Context) uint64 {
	//  GetSession("userID", c)
	session := sessions.Default(c)
	userID := session.Get("user2ID")
	str := fmt.Sprintf("%v", userID) //string donur
	str2 := stnccollection.StringtoUint64(str)
	return str2
}

// SetStoreUserID stores the userId for teh current user
func SetStoreUserID(c *gin.Context, userID uint64) {
	session := sessions.Default(c)
	session.Set("user2ID", userID)
	session.Save()
}

// ClearUserID clears the userId for the current suer //ClearUserIDFromCookie
func ClearUserID(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("user2ID", nil)
	session.Save()
}

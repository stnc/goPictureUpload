package stncsession

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SetSession session set eder
func SetSession(key string, value string, c *gin.Context) {
	session := sessions.Default(c)
	session.Set(key, value)
	session.Save()
}

//GetSession session get eder *--- çalışmıyor
func GetSession(key string, c *gin.Context) string {
	session := sessions.Default(c)
	if session != nil {
		return session.Get(key).(string)
	} else {
		return ""
	}
}

/*
//setCookie
func setCookie(c echo.Context, userName string) {

	if encoded, err := sifreleEncode(userName); err == nil {
		cookie := new(http.Cookie)
		cookie.Name = "login"
		cookie.Value = encoded
		cookie.Path = "/"
		cookie.Expires = time.Now().Add(24 * time.Hour)
		c.SetCookie(cookie)

	}

}

// for logout
func clearSession(c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "login"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Expires = time.Unix(0, 0)
	c.SetCookie(cookie)

}


*/

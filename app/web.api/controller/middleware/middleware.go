package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"stncCms/app/infrastructure/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// func csrfMiddleware(c *gin.Context) {
// 	csrf.Options{
// 		Secret: "StncWeb",
// 		ErrorFunc: func(c *gin.Context) {
// 			c.String(400, "CSRF token mismatch")
// 			c.Abort()
// 		},
// 	}
// }

//cors ozelleştirme https://stackoverflow.com/questions/29418478/go-gin-framework-cors
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		//selman ekledi
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("cache-Control", "no-cache, must-revalidate")
		c.Writer.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}

}

//Avoid a large file from loading into memory
//If the file size is greater than 8MB dont allow it to even load into memory and waste our time.
func MaxSizeAllowed(n int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, n)
		buff, errRead := c.GetRawData()
		if errRead != nil {
			//c.JSON(http.StatusRequestEntityTooLarge,"too large")
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"status":     http.StatusRequestEntityTooLarge,
				"upload_err": "too large: upload an image less than 8MB",
			})
			c.Abort()
			return
		}
		buf := bytes.NewBuffer(buff)
		c.Request.Body = ioutil.NopCloser(buf)
	}
}

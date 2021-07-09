package stnchelper

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func QueryOffset(c *gin.Context) int {
	offset := c.Query("start")
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0
	}
	return offsetInt
}

func QueryLimit(c *gin.Context) int {
	limit := c.Query("length")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 25
	}
	return limitInt
}

func DateTimeScope(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		begin := c.DefaultQuery("begin", "")
		end := c.DefaultQuery("end", "")

		if begin == "" && end == "" {
			return db
		}

		var tBegin, tEnd time.Time
		layout := "02.01.2006 15:04"

		if begin != "" {
			tBegin, _ = time.Parse(layout, begin)
		} else {
			t := time.Now()
			tBegin = t.AddDate(-20, 0, 0)
		}

		if end != "" {
			tEnd, _ = time.Parse(layout, end)
		} else {
			tEnd = time.Now()
		}

		return db.Where("created_at BETWEEN ? AND ?", tBegin, tEnd)
	}
}

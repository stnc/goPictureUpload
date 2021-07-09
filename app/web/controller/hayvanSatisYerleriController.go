package controller

import (
	"fmt"
	"log"
	"net/http"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stncsession"
	"stncCms/app/services"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//HayvanSatisYerleri constructor
type HayvanSatisYerleri struct {
	HayvanSatisYerleriApp services.HayvanSatisYerleriAppInterface
}

const viewPathHayvanSatisYerleri = "admin/hayvanSatisYerleri/"

//InitHayvanSatisYerleri post controller constructor
func InitHayvanSatisYerleri(kurbanbayramiApp services.HayvanSatisYerleriAppInterface) *HayvanSatisYerleri {
	return &HayvanSatisYerleri{
		HayvanSatisYerleriApp: kurbanbayramiApp,
	}
}

//Index list
func (access *HayvanSatisYerleri) Index(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	viewData := pongo2.Context{
		"paginator": paginator,
		"flashMsg":  flashMsg,
		"title":     "Hayvan Satış Yerleri",
	}

	c.HTML(
		http.StatusOK,
		viewPathHayvanSatisYerleri+"index.html",
		viewData,
	)
}

//Index list
func (access *HayvanSatisYerleri) ListDataTable(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	if result, err := access.HayvanSatisYerleriApp.ListDataTable(c); err == nil {
		c.JSON(200, result)
	} else {
		c.AbortWithStatus(404)
		log.Println(err)
		return
	}
}

//Create all list
func (access *HayvanSatisYerleri) Create(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	viewData := pongo2.Context{
		"title":    "Hayvan Satış Yeri Ekleme",
		"csrf":     csrf.GetToken(c),
		"flashMsg": flashMsg,
	}
	c.HTML(
		http.StatusOK,
		viewPathHayvanSatisYerleri+"create.html",
		viewData,
	)
}

//Store save method
func (access *HayvanSatisYerleri) Store(c *gin.Context) {
	//color.Blue("Prints %s in blue.", "text")
	stncsession.IsLoggedInRedirect(c)

	var kurbanbayrami, _, _ = kurbanBayramiModel(c)
	var saveError = make(map[string]string)

	saveError = kurbanbayrami.Validate()
	fmt.Println(saveError)
	if len(saveError) == 0 {

		saveData, saveErr := access.HayvanSatisYerleriApp.Save(&kurbanbayrami)
		if saveErr != nil {
			saveError = saveErr
		}
		lastID := strconv.FormatUint(uint64(saveData.ID), 10)
		stncsession.SetFlashMessage("Kayıt başarı ile eklendi", "success", c)
		c.Redirect(http.StatusMovedPermanently, "/admin/hayvanSatisYerleri/edit/"+lastID)
	} else {
		stncsession.SetFlashMessage("Zorunlu alanları lütfen doldurunuz", "danger", c)
	}
	flashMsg := stncsession.GetFlashMessage(c)
	viewData := pongo2.Context{
		"title":    "Hayvan Satış Yeri Ekleme",
		"csrf":     csrf.GetToken(c),
		"err":      saveError,
		"post":     kurbanbayrami,
		"flashMsg": flashMsg,
	}

	c.HTML(
		http.StatusOK,
		viewPathHayvanSatisYerleri+"create.html",
		viewData,
	)

}

//Edit genel kurban düzenleme işler
func (access *HayvanSatisYerleri) Edit(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	if kID, err := strconv.ParseUint(c.Param("kID"), 10, 64); err == nil {
		if posts, err := access.HayvanSatisYerleriApp.GetByID(kID); err == nil {
			viewData := pongo2.Context{
				"title":    "Hayvan Satış Yeri Düzenleme",
				"post":     posts,
				"csrf":     csrf.GetToken(c),
				"flashMsg": flashMsg,
			}
			c.HTML(
				http.StatusOK,
				viewPathHayvanSatisYerleri+"edit.html",
				viewData,
			)

		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

//Update data
func (access *HayvanSatisYerleri) Update(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	var kurban, _, id = kurbanBayramiModel(c)
	var saveError = make(map[string]string)

	saveError = kurban.Validate()

	if len(saveError) == 0 {

		_, saveErr := access.HayvanSatisYerleriApp.Update(&kurban)
		if saveErr != nil {
			saveError = saveErr
		}
		stncsession.SetFlashMessage("Kayıt başarı ile düzenlendi", "success", c)
		c.Redirect(http.StatusMovedPermanently, "/admin/hayvanSatisYerleri/edit/"+id)
		return
	}
	viewData := pongo2.Context{
		"title": "Hayvan Satış Yeri Düzenleme",
		"err":   saveError,
		"csrf":  csrf.GetToken(c),
		"post":  kurban,
	}
	c.HTML(
		http.StatusOK,
		viewPathHayvanSatisYerleri+"edit.html",
		viewData,
	)
}

/***  DATA MODEL   ***/
func kurbanBayramiModel(c *gin.Context) (data entity.HayvanSatisYerleri, idD uint64, idStr string) {
	id := c.PostForm("ID")
	idInt, _ := strconv.Atoi(id)
	var idN uint64
	idN = uint64(idInt)
	data.ID = idN
	data.Durum = 1
	data.UserID = stncsession.GetUserID2(c)
	data.YerAdi = c.PostForm("YerAdi")
	data.Adresi = c.PostForm("Adresi")
	data.IlgiliKisi = c.PostForm("IlgiliKisi")
	data.Telefon = c.PostForm("Telefon")

	return data, idN, id
}

package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stnccollection"
	"stncCms/app/domain/helpers/stncdatetime"
	"stncCms/app/domain/helpers/stncsession"
	"stncCms/app/services"

	"strconv"

	"github.com/astaxie/beego/utils/pagination"
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//Kisiler constructor
type Kisiler struct {
	KisilerApp services.KisilerAppInterface
}

const viewPathGKisiler = "admin/kisiler/"

//InitGKisiler post controller constructor
func InitGKisiler(gkApp services.KisilerAppInterface) *Kisiler {
	return &Kisiler{
		KisilerApp: gkApp,
	}
}

//Index list
func (access *Kisiler) Index(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	viewData := pongo2.Context{
		"paginator": paginator,
		"flashMsg":  flashMsg,
		"title":     "Kişi ",
	}

	c.HTML(
		http.StatusOK,
		viewPathGKisiler+"index.html",
		viewData,
	)
}

//Index list
func (access *Kisiler) ListDataTable(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)

	if result, err := access.KisilerApp.ListDataTable(c); err == nil {
		c.JSON(200, result)
	} else {
		c.AbortWithStatus(404)
		log.Println(err)
		return
	}

}

//Index list
func (access *Kisiler) IndexV1(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	var tarih stncdatetime.Inow
	var total int64
	access.KisilerApp.Count(&total)
	postsPerPage := 5
	paginator := pagination.NewPaginator(c.Request, postsPerPage, total)
	offset := paginator.Offset()
	posts, _ := access.KisilerApp.GetAllP(postsPerPage, offset)

	viewData := pongo2.Context{
		"paginator": paginator,
		"title":     "Kişi Ekleme",
		"posts":     posts,
		"flashMsg":  flashMsg,
		"csrf":      csrf.GetToken(c),
		"tarih":     tarih,
	}

	c.HTML(
		http.StatusOK,
		viewPathGKisiler+"indexList.html",
		viewData,
	)
}

//Create all list f
func (access *Kisiler) Create(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	kisiler, _ := access.KisilerApp.GetAll()
	viewData := pongo2.Context{
		"title":    "Kişi  Ekleme",
		"csrf":     csrf.GetToken(c),
		"flashMsg": flashMsg,
		"kisiler":  kisiler,
	}
	c.HTML(
		http.StatusOK,
		viewPathGKisiler+"create.html",
		viewData,
	)
}

//Store save method
func (access *Kisiler) Store(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	var Kisiler, _, _ = gKisilerModel("create", c)
	var savePostError = make(map[string]string)
	kisiler, _ := access.KisilerApp.GetAll()
	savePostError = Kisiler.Validate()

	if len(savePostError) == 0 {

		// fmt.Println(Kisiler)
		// return
		saveData, saveErr := access.KisilerApp.Save(&Kisiler)
		if saveErr != nil {
			savePostError = saveErr
		}
		lastID := strconv.FormatUint(uint64(saveData.ID), 10)
		stncsession.SetFlashMessage("Kayıt başarı ile eklendi", "success", c)
		c.Redirect(http.StatusMovedPermanently, "/admin/kisiler/edit/"+lastID)
		return
	} else {
		stncsession.SetFlashMessage("Zorunlu alanları lütfen doldurunuz", "danger", c)
	}

	flashMsg := stncsession.GetFlashMessage(c)

	viewData := pongo2.Context{
		"title":    "Kisi Ekleme",
		"csrf":     csrf.GetToken(c),
		"err":      savePostError,
		"post":     Kisiler,
		"kisiler":  kisiler,
		"flashMsg": flashMsg,
	}
	c.HTML(
		http.StatusOK,
		viewPathGKisiler+"create.html",
		viewData,
	)
}

//Edit genel Kisiler düzenleme işler
//Edit edit data
func (access *Kisiler) Edit(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	var tarih stncdatetime.Inow
	action := c.DefaultQuery("action", "ekle")
	if ID, err := strconv.ParseUint(c.Param("ID"), 10, 64); err == nil {
		if posts, err := access.KisilerApp.GetByIDRel(ID); err == nil {
			kisiler, _ := access.KisilerApp.GetAll()

			referansData, _ := access.KisilerApp.GetByID(posts.ReferansKisi1)

			viewData := pongo2.Context{
				"title":        "Kişi ekleme",
				"post":         posts,
				"csrf":         csrf.GetToken(c),
				"flashMsg":     flashMsg,
				"tarih":        tarih,
				"kisiler":      kisiler,
				"referansData": referansData,
				"actionHref":   action,
			}
			c.HTML(
				http.StatusOK,
				viewPathGKisiler+"edit.html",
				viewData,
			)
			empJSON, err := json.MarshalIndent(posts, "", "  ")
			if err != nil {
				log.Fatalf(err.Error())
			}
			fmt.Printf("MarshalIndent funnction output\n %s\n", string(empJSON))
			return
		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

//Update data
func (access *Kisiler) Update(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	KisilerUpdate, id, _ := gKisilerModel("edit", c)

	var savePostError = make(map[string]string)

	savePostError = KisilerUpdate.Validate()
	kisiler, _ := access.KisilerApp.GetAll()
	if len(savePostError) == 0 {

		_, saveErr := access.KisilerApp.Update(&KisilerUpdate)
		if saveErr != nil {
			savePostError = saveErr
		}

		stncsession.SetFlashMessage("Kayıt başarı ile düzenlendi", "success", c)

		c.Redirect(http.StatusMovedPermanently, "/admin/kisiler/edit/"+id)
		return
	} else {
		stncsession.SetFlashMessage("Zorunlu alanları lütfen doldurunuz", "danger", c)
	}

	flashMsg := stncsession.GetFlashMessage(c)

	viewData := pongo2.Context{
		"title":    "Kisiler Düzenleme",
		"err":      savePostError,
		"flashMsg": flashMsg,
		"csrf":     csrf.GetToken(c),
		"post":     KisilerUpdate,
		"kisiler":  kisiler,
	}

	c.HTML(
		http.StatusOK,
		viewPathGKisiler+"edit.html",
		viewData,
	)
}

//SearchAjax ajax search
func (access *Kisiler) KisiAraAjax(c *gin.Context) {
	q := c.PostForm("q")
	fmt.Println(q)
	stncsession.IsLoggedInRedirect(c)

	returnData, posts, _ := access.KisilerApp.Search(q)
	fmt.Println(posts)
	if returnData["status"] != "ok" {
		//c.String(http.StatusOK, returnData)
		c.JSON(http.StatusOK, returnData)
	} else {
		c.JSON(http.StatusOK, posts)
	}
}

//referansEkleAjax save buraası referans kişi eklerken kullanılıyor
func (access *Kisiler) KisiEkleAjax(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	// var kurban, _, _ = gkurbanModel("referans", c)
	var kisiler = entity.Kisiler{}

	kisiler.UserID = stncsession.GetUserID2(c)
	adSoyad := c.PostForm("AdSoyad")

	kisiler.AdSoyad = adSoyad

	kisiler.Telefon = c.PostForm("Telefon")
	kisiler.Adres = c.PostForm("Adres")

	var savePostError = make(map[string]string)

	savePostError = kisiler.Validate()

	if len(savePostError) == 0 {
		saveData, saveErr := access.KisilerApp.Save(&kisiler)
		if saveErr != nil {
			savePostError = saveErr
		}

		lastID := strconv.FormatUint(uint64(saveData.ID), 10)
		viewData := pongo2.Context{
			"title":    "Kişi Ekleme",
			"csrf":     csrf.GetToken(c),
			"lastID":   lastID,
			"viewID":   c.PostForm("viewID"),
			"username": saveData.AdSoyad,
			"tel":      saveData.Telefon,
			"err":      savePostError,
			"status":   "ok",
			"msg":      "Kayıt Başarı ile Eklendi",
		}
		c.JSON(http.StatusOK, viewData)
		return
	} else {
		viewData := pongo2.Context{
			"title":  "Kişi Ekleme",
			"csrf":   csrf.GetToken(c),
			"status": "error",
			"err":    savePostError,
		}
		c.JSON(http.StatusOK, viewData)
		return
	}

}

/***  POST MODEL   ***/
func gKisilerModel(formType string, c *gin.Context) (data entity.Kisiler, idString string, err error) {

	id := c.PostForm("ID")

	idInt, _ := strconv.Atoi(id)

	var idN uint64

	idN = uint64(idInt)

	data.ID = idN
	data.UserID = stncsession.GetUserID2(c)
	adSoyad := c.PostForm("AdSoyad")
	data.AdSoyad = adSoyad
	data.Telefon = c.PostForm("Telefon")
	data.Email = c.PostForm("Email")
	data.Adres = c.PostForm("Adres")
	data.Aciklama = c.PostForm("Aciklama")
	data.ReferansKisi1 = stnccollection.StringtoUint64(c.PostForm("ReferansKisi1"))
	// data.ReferansKisi2 = stnccollection.StringtoUint64(c.PostForm("ReferansKisi2"))
	return data, id, nil
}

//Delete data
func (access *Kisiler) Delete(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)

	if postID, err := strconv.ParseUint(c.Param("ID"), 10, 64); err == nil {
		var toplamKisi int
		var toplamRefKisi int
		access.KisilerApp.HasKisiKurban(postID, &toplamKisi)
		access.KisilerApp.HasKisiReferans(postID, &toplamRefKisi)
		fmt.Println("toplamKisi")
		fmt.Println(toplamKisi)

		fmt.Println("toplamRefKisi")
		fmt.Println(toplamRefKisi)

		if toplamRefKisi > 0 {
			stncsession.SetFlashMessage("Kişi Referans Olarak Kayıtlıdır,Silemezsiniz", "danger", c)
			c.Redirect(http.StatusMovedPermanently, "/admin/kisiler")
			return
		}
		if toplamKisi > 0 {
			stncsession.SetFlashMessage("Kişi Kurbana Kayıtlıdır,Silemezsiniz", "danger", c)
			c.Redirect(http.StatusMovedPermanently, "/admin/kisiler")
			return
		} else {
			access.KisilerApp.Delete(postID)
			stncsession.SetFlashMessage("Silindi", "success", c)
			c.Redirect(http.StatusMovedPermanently, "/admin/kisiler")
			return
		}

	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}

}

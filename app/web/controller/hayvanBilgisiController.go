package controller

import (
	"fmt"
	"log"
	"net/http"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stnc2upload"
	"stncCms/app/domain/helpers/stnccollection"
	"stncCms/app/domain/helpers/stncsession"
	"stncCms/app/domain/helpers/stncupload"
	"stncCms/app/domain/repository"
	"stncCms/app/services"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

const viewPathHayvanBilgisi = "admin/hayvanBilgisi/"

//localhayvanSatis constructor
type localHayvanSatis struct {
	HayvanSatisApp services.HayvanSatisYerleriAppInterface
	OptionsApp     services.OptionsAppInterface
}

//localHayvanSatisYerleri post controller constructor
func localHayvanSatisYerleri(hayvanSatisApp services.HayvanSatisYerleriAppInterface, optionsApp services.OptionsAppInterface) *localHayvanSatis {
	return &localHayvanSatis{
		HayvanSatisApp: hayvanSatisApp,
		OptionsApp:     optionsApp,
	}
}

//HayvanBilgisi constructor
type HayvanBilgisi struct {
	HayvanBilgisiApp services.HayvanBilgisiAppInterface
	MediaApp         services.MediaAppInterface
}

//InitHayvanBilgisi post controller constructor
func InitHayvanBilgisi(hayvanBilgisiApp services.HayvanBilgisiAppInterface, mediaApp services.MediaAppInterface) *HayvanBilgisi {
	return &HayvanBilgisi{
		HayvanBilgisiApp: hayvanBilgisiApp,
		MediaApp:         mediaApp,
	}
}

//UploadConfig upload media
func (access *HayvanBilgisi) UploadConfig() map[string]string {
	returnData := make(map[string]string)

	maxFiles := 20
	uploadSize := 10
	modulID := 3
	bigImageWidth := 1200
	bigImageHeight := 768
	thumbImageWidth := 100
	thumbImageHeight := 100
	uploadFile := "hayvanBilgisi"

	uploadPath := "public/upl/" + uploadFile + "/"
	uploadFsPath := "upload/" + uploadFile + "/"

	returnData["uploadPath"] = uploadPath
	returnData["modulName"] = "hayvanBilgisi"
	returnData["modulID"] = stnccollection.IntToString(modulID)

	returnData["uploadSymbol"] = uploadFsPath
	returnData["maxFiles"] = stnccollection.IntToString(maxFiles)

	returnData["bigImageWidth"] = stnccollection.IntToString(bigImageWidth)
	returnData["bigImageHeight"] = stnccollection.IntToString(bigImageHeight)

	returnData["thumbImageWidth"] = stnccollection.IntToString(thumbImageWidth)
	returnData["thumbImageHeight"] = stnccollection.IntToString(thumbImageHeight)

	returnData["uploadSize"] = stnccollection.IntToString(uploadSize)

	returnData["uploadSize"] = stnccollection.IntToString(uploadSize)
	returnData["fileType"] = "video/mp4,image/jpeg,image/jpg,image/gif,image/png,video/webm" //application/pdf
	returnData["text"] = ""
	return returnData
}

/*******datatable *****/
func (access *HayvanBilgisi) Index(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	viewData := pongo2.Context{
		"paginator": paginator,
		"flashMsg":  flashMsg,
		"title":     "Hayvan Bilgisi",
	}

	c.HTML(
		http.StatusOK,
		viewPathHayvanBilgisi+"index.html",
		viewData,
	)
}

//Index list
func (access *HayvanBilgisi) ListDataTable(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)

	if result, err := access.HayvanBilgisiApp.ListDataTable(c); err == nil {
		c.JSON(200, result)
	} else {
		c.AbortWithStatus(404)
		log.Println(err)
		return
	}

}

//Create all list
func (access *HayvanBilgisi) Create(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	db := repository.DB
	services, err1 := repository.RepositoriesInit(db)
	if err1 != nil {
		panic(err1)
	}

	hayvanSatisYerler := localHayvanSatisYerleri(services.HayvanSatisYerleri, services.Options)
	hayvanSatisyerleri, _ := hayvanSatisYerler.HayvanSatisApp.GetAll()

	viewData := pongo2.Context{
		"title":              "Hayvan Bilgisi Ekleme",
		"csrf":               csrf.GetToken(c),
		"alisFiyati1":        hayvanSatisYerler.OptionsApp.GetOption("alis_birim_fiyati_1"),
		"alisFiyati2":        hayvanSatisYerler.OptionsApp.GetOption("alis_birim_fiyati_2"),
		"alisFiyati3":        hayvanSatisYerler.OptionsApp.GetOption("alis_birim_fiyati_3"),
		"hayvanSatisyerleri": hayvanSatisyerleri,
		"flashMsg":           flashMsg,
	}
	c.HTML(
		http.StatusOK,
		viewPathHayvanBilgisi+"create.html",
		viewData,
	)
}

//Store save method
func (access *HayvanBilgisi) Store(c *gin.Context) {
	//color.Blue("Prints %s in blue.", "text")
	stncsession.IsLoggedInRedirect(c)
	var hayvanBilgisiData, _, _ = hayvanBilgisiModel(c)
	var saveError = make(map[string]string)

	saveError = hayvanBilgisiData.Validate()
	fmt.Println(saveError)

	db := repository.DB
	services, err1 := repository.RepositoriesInit(db)
	if err1 != nil {
		panic(err1)
	}

	yerler := localHayvanSatisYerleri(services.HayvanSatisYerleri, services.Options)
	hayvanSatisyerleri, _ := yerler.HayvanSatisApp.GetAll()

	sendFileName := "resim"
	filenameForm, _ := c.FormFile(sendFileName)
	filename, uploadError := stnc2upload.NewFileUpload().UploadFile(filenameForm, c.PostForm("Resim2"))

	if filename == "false" {
		saveError[sendFileName+"_error"] = uploadError
		saveError[sendFileName+"_valid"] = "is-invalid"
	}

	kupeNo := stnccollection.StringToint(c.PostForm("KupeNo"))
	var kupeNoCount int
	access.HayvanBilgisiApp.KupeNoCount(kupeNo, &kupeNoCount)

	if kupeNoCount > 0 {
		saveError["KupeNo_error"] = "küpe numarası daha önce girilmiş, lütfen başka bir numara giriniz"
		saveError["KupeNo_valid"] = "is-invalid"
		saveError["KupeNo"] = "küpe numarası daha önce girilmiş, lütfen başka bir numara giriniz"
		hayvanBilgisiData.KupeNo = uint64(kupeNo)
	}

	if len(saveError) == 0 {
		hayvanBilgisiData.Resim = filename
		hayvanBilgisiData.KupeNo = uint64(kupeNo)
		saveData, saveErr := access.HayvanBilgisiApp.Save(&hayvanBilgisiData)

		if saveErr != nil {
			saveError = saveErr
		}

		lastID := strconv.FormatUint(uint64(saveData.ID), 10)

		stncsession.SetFlashMessage("Kayıt başarı ile eklendi", "success", c)

		if c.PostForm("kaydet") == "kaydet" {
			c.Redirect(http.StatusMovedPermanently, "/admin/hayvanBilgisi/edit/"+lastID)
		} else {
			c.Redirect(http.StatusMovedPermanently, "/admin/hayvanBilgisi/create")
		}
		return
	} else {
		stncsession.SetFlashMessage("Zorunlu alanları lütfen doldurunuz", "danger", c)

	}
	flashMsg := stncsession.GetFlashMessage(c)
	viewData := pongo2.Context{
		"title":              "Hayvan Bilgisi Ekleme",
		"csrf":               csrf.GetToken(c),
		"err":                saveError,
		"hayvanSatisyerleri": hayvanSatisyerleri,
		"post":               hayvanBilgisiData,
		"alisFiyati1":        yerler.OptionsApp.GetOption("alis_birim_fiyati_1"),
		"alisFiyati2":        yerler.OptionsApp.GetOption("alis_birim_fiyati_2"),
		"alisFiyati3":        yerler.OptionsApp.GetOption("alis_birim_fiyati_3"),
		"flashMsg":           flashMsg,
	}
	c.HTML(
		http.StatusOK,
		viewPathHayvanBilgisi+"create.html",
		viewData,
	)
}

//Edit genel kurban düzenleme işler
func (access *HayvanBilgisi) Edit(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	db := repository.DB
	services, err1 := repository.RepositoriesInit(db)
	if err1 != nil {
		panic(err1)
	}

	hayvanSatisYerler := localHayvanSatisYerleri(services.HayvanSatisYerleri, services.Options)
	hayvanSatisyerleri, _ := hayvanSatisYerler.HayvanSatisApp.GetAll()

	if kID, err := strconv.ParseUint(c.Param("kID"), 10, 64); err == nil {
		if posts, err := access.HayvanBilgisiApp.GetByID(kID); err == nil {
			accData := access.UploadConfig()
			mediaModulID := stnccollection.StringToint(accData["modulID"])
			mediaData, _ := access.MediaApp.GetAll(mediaModulID, int(kID))

			viewData := pongo2.Context{
				"title":              "Hayvan Bilgisi duzenleme",
				"post":               posts,
				"hayvanSatisyerleri": hayvanSatisyerleri,
				"alisFiyati1":        hayvanSatisYerler.OptionsApp.GetOption("alis_birim_fiyati_1"),
				"alisFiyati2":        hayvanSatisYerler.OptionsApp.GetOption("alis_birim_fiyati_2"),
				"alisFiyati3":        hayvanSatisYerler.OptionsApp.GetOption("alis_birim_fiyati_3"),
				"csrf":               csrf.GetToken(c),
				"flashMsg":           flashMsg,
				"medias":             mediaData,
				"fileConfig":         accData,
				"ID":                 kID,
			}
			c.HTML(
				http.StatusOK,
				viewPathHayvanBilgisi+"edit.html",
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
func (access *HayvanBilgisi) Update(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	var kurban, _, id = hayvanBilgisiModel(c)
	var saveError = make(map[string]string)

	saveError = kurban.Validate()

	sendFileName := "resim"
	filenameForm, _ := c.FormFile(sendFileName)

	db := repository.DB
	services, err1 := repository.RepositoriesInit(db)
	if err1 != nil {
		panic(err1)
	}

	yerler := localHayvanSatisYerleri(services.HayvanSatisYerleri, services.Options)
	hayvanSatisyerleri, _ := yerler.HayvanSatisApp.GetAll()

	filename, uploadError := stnc2upload.NewFileUpload().UploadFile(filenameForm, c.PostForm("Resim2"))
	if filename == "false" {
		saveError[sendFileName+"_error"] = uploadError
		saveError[sendFileName+"_valid"] = "is-invalid"
	}

	kupeNo := stnccollection.StringToint(c.PostForm("KupeNo"))
	var kupeNoCount int
	access.HayvanBilgisiApp.KupeNoCount(kupeNo, &kupeNoCount)

	if kupeNoCount > 1 {
		saveError["KupeNo_error"] = "küpe numarası daha önce girilmiş, lütfen başka bir numara giriniz"
		saveError["KupeNo_valid"] = "is-invalid"
		saveError["KupeNo"] = "küpe numarası daha önce girilmiş, lütfen başka bir numara giriniz"
		kurban.KupeNo = uint64(kupeNo)
	}

	if len(saveError) == 0 {
		kurban.Resim = filename
		kurban.KupeNo = uint64(kupeNo)
		_, saveErr := access.HayvanBilgisiApp.Update(&kurban)
		if saveErr != nil {
			saveError = saveErr
		}
		stncsession.SetFlashMessage("Kayıt başarı ile düzenlendi", "success", c)
		if c.PostForm("kaydet") == "kaydet" {
			c.Redirect(http.StatusMovedPermanently, "/admin/hayvanBilgisi/edit/"+id)
		} else {
			c.Redirect(http.StatusMovedPermanently, "/admin/hayvanBilgisi/create")
		}
		return
	}
	viewData := pongo2.Context{
		"title":              "Ödeme Düzenleme",
		"err":                saveError,
		"csrf":               csrf.GetToken(c),
		"post":               kurban,
		"flashMsg":           flashMsg,
		"alisFiyati1":        yerler.OptionsApp.GetOption("alis_birim_fiyati_1"),
		"alisFiyati2":        yerler.OptionsApp.GetOption("alis_birim_fiyati_2"),
		"alisFiyati3":        yerler.OptionsApp.GetOption("alis_birim_fiyati_3"),
		"hayvanSatisyerleri": hayvanSatisyerleri,
	}
	c.HTML(
		http.StatusOK,
		viewPathHayvanBilgisi+"edit.html",
		viewData,
	)
}

//HayvanListeAjax genel kurban düzenleme işler
func (access *HayvanBilgisi) HayvanListeAjax(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)

	if hID, err := strconv.ParseUint(c.Param("hID"), 10, 64); err == nil {
		if jsonData, err := access.HayvanBilgisiApp.GetByIDRelated(hID); err == nil {

			viewData := pongo2.Context{
				"jsonData": jsonData,
				"csrf":     csrf.GetToken(c),
			}
			c.JSON(http.StatusOK, viewData)
			return

		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

/***  DATA MODEL   ***/
func hayvanBilgisiModel(c *gin.Context) (data entity.HayvanBilgisi, idD uint64, idStr string) {
	id := c.PostForm("ID")
	idInt, _ := strconv.Atoi(id)
	var idN uint64
	idN = uint64(idInt)
	data.ID = idN
	data.UserID = stncsession.GetUserID2(c)
	data.HayvanCinsi = c.PostForm("HayvanCinsi")
	data.HayvanSatisYerleriID = stnccollection.StringtoUint64(c.PostForm("HayvanSatisYerleriID"))

	data.Agirlik = stnccollection.StringToint(c.PostForm("Agirlik"))
	data.AlisFiyatTuru = stnccollection.StringToint(c.PostForm("AlisFiyatTuru"))
	data.AlisFiyati, _ = stnccollection.StringToFloat64(c.PostForm("AlisFiyati"))
	data.Durum = 1
	data.HayvanSatisYerleri.YerAdi = "_"
	data.HayvanSatisYerleri.Adresi = "_"
	data.HayvanSatisYerleri.Durum = 1
	return data, idN, id
}

//Upload
func (access *HayvanBilgisi) Upload(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	IDint := stnccollection.StringToint(c.Param("ID"))
	filenameForm, _ := c.FormFile("file")
	accData := access.UploadConfig()
	uploadSize := stnccollection.StringToint(accData["uploadSize"])
	fileType := accData["fileType"]
	fileTypes := strings.Split(fileType, ",")
	maxFiles := stnccollection.StringToint(accData["maxFiles"])
	uploadPath := accData["uploadPath"]
	upl := stncupload.FileUpload{
		UploadPath: uploadPath,
		UploadSize: uploadSize,
		MaxFiles:   maxFiles,
		Types:      fileTypes,
	}

	var up stncupload.UploadFileInterface = upl

	up.InitUploader(c, IDint, filenameForm, accData)

}

//MediaDelete Delete file data
func (access *HayvanBilgisi) MediaDelete(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	if ID, err := strconv.ParseUint(c.Param("ID"), 10, 64); err == nil {
		accData := access.UploadConfig()
		uploadPath := accData["uploadPath"]
		stncupload.NewFileUpload().MediaDelete(c, ID, uploadPath)
	}
}

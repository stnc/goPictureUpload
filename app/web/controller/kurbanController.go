package controller

import (
	"fmt"
	"log"
	"net/http"
	"stncCms/app/domain/dto"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stnccollection"
	"stncCms/app/domain/helpers/stncdatetime"
	"stncCms/app/domain/helpers/stnchelper"
	"stncCms/app/domain/helpers/stncsession"
	"stncCms/app/domain/helpers/stncupload"
	"stncCms/app/domain/repository"
	"stncCms/app/services"
	"strings"

	"strconv"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//Kurban constructor
type Kurban struct {
	kurbanApp services.KurbanAppInterface
	kisiApp   services.KisilerAppInterface
	MediaApp  services.MediaAppInterface
}

const viewPathGkurban = "admin/kurbanlar/"

//InitGkurban post controller constructor
func InitGkurban(gkApp services.KurbanAppInterface, kisiApp services.KisilerAppInterface, mediaApp services.MediaAppInterface) *Kurban {
	return &Kurban{
		kurbanApp: gkApp,
		kisiApp:   kisiApp,
		MediaApp:  mediaApp,
	}
}
func (access *Kurban) UploadConfig(modulID int, text string) map[string]string {
	returnData := make(map[string]string)

	maxFiles := 5
	uploadSize := 5
	// modulID := 3
	bigImageWidth := 600
	bigImageHeight := 600
	thumbImageWidth := 150
	thumbImageHeight := 150
	var uploadFile string = "kurbanlar"

	if modulID == 3 {
		uploadFile = "hayvanBilgisi"
	}

	uploadPath := "public/upl/" + uploadFile + "/"
	uploadFsPath := "upload/" + uploadFile + "/"

	returnData["uploadPath"] = uploadPath
	returnData["modulName"] = "kurban"
	returnData["modulID"] = stnccollection.IntToString(modulID)
	returnData["uploadSymbol"] = uploadFsPath
	returnData["maxFiles"] = stnccollection.IntToString(maxFiles)

	returnData["bigImageWidth"] = stnccollection.IntToString(bigImageWidth)
	returnData["bigImageHeight"] = stnccollection.IntToString(bigImageHeight)

	returnData["thumbImageWidth"] = stnccollection.IntToString(thumbImageWidth)
	returnData["thumbImageHeight"] = stnccollection.IntToString(thumbImageHeight)

	returnData["uploadSize"] = stnccollection.IntToString(uploadSize)
	returnData["fileType"] = "video/mp4,image/jpeg,image/jpg,image/gif,image/png,video/webm" //application/pdf
	returnData["text"] = text
	return returnData
}

//Index list
func (access *Kurban) Index(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	viewData := pongo2.Context{
		"paginator": paginator,
		"flashMsg":  flashMsg,
		"title":     "Kurban",
	}

	c.HTML(
		http.StatusOK,
		viewPathGkurban+"index.html",
		viewData,
	)

} //Index list
func (access *Kurban) KurbanBilgi(c *gin.Context) {
	var tarih stncdatetime.Inow
	slug := c.Param("slug")

	if kurbanData, err := access.kurbanApp.GetKurbanOpenInfo(slug); err == nil {

		var kisilerList = []dto.KurbanListForGrouplar{}
		kisilerList, _ = access.kurbanApp.GetAllKurbanAndKisiler(int(kurbanData.KurbanGrupId))

		// stnchelper.TestModeCode()
		for no, kisi := range kisilerList {
			if kisi.KisiID != 1 {
				ilkKarakterler := stnchelper.LetterPrefix(kisi.KisiAdSoyad)
				fmt.Println(ilkKarakterler)
				sonKarakterler := stnchelper.LetterSuffix(kisi.KisiAdSoyad)
				fmt.Println(sonKarakterler)
				ortaKarakterler := stnchelper.LetterMiddle(kisi.KisiAdSoyad)
				fmt.Println(ortaKarakterler)
				ortaYizdizla := stnchelper.LetterMiddleStars(ortaKarakterler)
				fmt.Println(ortaYizdizla)
				yeniKarakter := ilkKarakterler + ortaYizdizla + sonKarakterler
				fmt.Println(yeniKarakter)
				kisilerList[no].KisiAdSoyad = yeniKarakter
			} else {
				kisilerList[no].KisiAdSoyad = "BOŞ"
			}
		}

		kisiBasiKilo := kurbanData.HayvanAgirlik / uint64(kurbanData.HissedarAdet)

		mediaData, _ := access.MediaApp.GetAll(3, int(kurbanData.HayvanBilgisiId))

		viewData := pongo2.Context{
			"title":        "Bilgi",
			"durum":        true,
			"post":         kurbanData,
			"gruplar":      kisilerList,
			"kisiBasiKilo": kisiBasiKilo,
			"tarih":        tarih,
			"media":        mediaData,
			"slug":         slug,
		}

		c.HTML(
			http.StatusOK,
			"admin/KurbanBilgi.html",
			viewData,
		)
	} else {
		viewData := pongo2.Context{
			"title": "Bilgi",
			"durum": false,
		}

		c.HTML(
			http.StatusOK,
			"admin/KurbanBilgi.html",
			viewData,
		)
	}
}

//Index list
func (access *Kurban) ListDataTable(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	if result, err := access.kurbanApp.ListDataTable(c); err == nil {
		c.JSON(200, result)
	} else {
		c.AbortWithStatus(404)
		log.Println(err)
		return
	}
}

//Create all list f
func (access *Kurban) Create(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	viewData := pongo2.Context{
		"title":    "Kurban  Ekleme",
		"csrf":     csrf.GetToken(c),
		"flashMsg": flashMsg,
	}
	c.HTML(
		http.StatusOK,
		viewPathGkurban+"create.html",
		viewData,
	)
}

//Store save method
func (access *Kurban) Store(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	var kurban, _, _ = gkurbanModel("create", c)
	var savePostError = make(map[string]string)

	savePostError = kurban.Validate()

	if len(savePostError) == 0 {

		// fmt.Println(kurban)
		// return
		saveData, saveErr := access.kurbanApp.Save(&kurban)
		if saveErr != nil {
			savePostError = saveErr
		}

		if c.PostForm("kaydet") == "gruplaraDon" {
			c.Redirect(http.StatusMovedPermanently, "/admin/gruplar/")
			return
		} else {
			lastID := strconv.FormatUint(uint64(saveData.ID), 10)
			stncsession.SetFlashMessage("Kayıt başarı ile eklendi", "success", c)
			c.Redirect(http.StatusMovedPermanently, "/admin/kurban/edit/"+lastID)
		}

		return
	} else {
		stncsession.SetFlashMessage("Zorunlu alanları lütfen doldurunuz", "danger", c)
	}
	flashMsg := stncsession.GetFlashMessage(c)
	kisi1 := stnccollection.StringtoUint64(c.PostForm("Kisi1"))
	referansKisi1 := stnccollection.StringtoUint64(c.PostForm("ReferansKisi1"))
	// referansKisi2 := stnccollection.StringtoUint64(c.PostForm("ReferansKisi2"))

	kisiBilgileri, _ := access.kisiApp.GetByID(kisi1)
	referansKisiData1, _ := access.kisiApp.GetByID(referansKisi1)
	// referansKisiData2, _ := access.kisiApp.GetByIDReferans(referansKisi2)

	viewData := pongo2.Context{
		"title":         "Kurban Ekleme",
		"csrf":          csrf.GetToken(c),
		"err":           savePostError,
		"post":          kurban,
		"flashMsg":      flashMsg,
		"referansData1": referansKisiData1,
		"kisiBilgileri": kisiBilgileri,
	}
	c.HTML(
		http.StatusOK,
		viewPathGkurban+"create.html",
		viewData,
	)
}

//Edit genel kurban düzenleme işlero
func (access *Kurban) Edit(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	var ucretOdemesiDurumu bool = false
	// var kurbanSonKalanBorc float64
	var tarih stncdatetime.Inow
	action := c.DefaultQuery("action", "ekle")

	flashMsg := stncsession.GetFlashMessage(c)
	var referansKisi1 *entity.Kisiler
	var kisiBilgileri *entity.Kisiler
	if kurbanID, err := strconv.ParseUint(c.Param("kurbanID"), 10, 64); err == nil {
		if kurbanData, err := access.kurbanApp.GetByID(kurbanID); err == nil {

			if kurbanData.KisiID != 0 {
				kisiBilgileri, _ = access.kisiApp.GetByID(kurbanData.KisiID)

				if kisiBilgileri != nil && kisiBilgileri.ReferansKisi1 != 0 {
					referansKisi1, _ = access.kisiApp.GetByID(kisiBilgileri.ReferansKisi1)
				}
			}

			// referansKisi2, _ := access.kisiApp.GetByIDReferans(kurbanData.ReferansKisi2)

			kurbanTuru := kurbanData.KurbanTuru
			var modulId int = 1
			var text string = ""
			var mediaData []entity.Media
			if kurbanTuru == 12 {
				modulId = 3
				hayvanID := access.kurbanApp.GetKurbanHayvanID(kurbanData.GrupID)
				if hayvanID != 0 {
					text = `<a href="/admin/hayvanBilgisi/edit/` + stnccollection.IntToString(hayvanID) + `">Resim düzenlemek yada eklemek tıklayınız</a>`
				} else {
					text = `Gruba hayvan atanmadığı için resim düzenleyemezsiniz`
				}
				mediaData, _ = access.MediaApp.GetAllforModul(modulId, hayvanID)
			} else {
				mediaData, _ = access.MediaApp.GetAll(modulId, int(kurbanID))
			}

			accData := access.UploadConfig(modulId, text)

			// 	//selman tam liste #liste # fmtlist #stncfmt #fmtstnc
			// fmt.Printf("%+v\n", kurbanData)

			if kurbanData.BorcDurum == entity.KurbanBorcDurumKasaBorcluDurumda {
				fmt.Println("Kurba Borc DurumKasaBorcluDurumda")

				//burada kasa borçlu demek onun işlmeleri olacak
				//onun için ayrı bir sayfa include edilebilir bunun orneği evfada var
				kasaBorcu := access.kurbanApp.KasaBorcu(kurbanID)
				kurbanOdenenMiktar := access.kurbanApp.OdemelerToplami(kurbanID)

				viewData := pongo2.Context{
					"title":         "Kurban duzenleme",
					"durum":         "Kasa Borçlu",
					"post":          kurbanData,
					"referansData1": referansKisi1,
					// "referansData2":      referansKisi2,
					"kisiBilgileri":      kisiBilgileri,
					"tarih":              tarih,
					"ucretOdemesiDurumu": ucretOdemesiDurumu,
					"kurbanOdenenMiktar": kurbanOdenenMiktar,
					"kasaBorcu":          kasaBorcu,
					"hayvanAtanmasi":     true,
					"kurbanDurum":        kurbanData.Durum,
					"csrf":               csrf.GetToken(c),
					"flashMsg":           flashMsg,
					"medias":             mediaData,
					"fileConfig":         accData,
					"actionHref":         action,
					"ID":                 kurbanID,
				}
				c.HTML(
					http.StatusOK,
					viewPathGkurban+"edit.html",
					viewData,
				)
			} else if kurbanData.BorcDurum == entity.KurbanBorcDurumHesapKapandi {
				fmt.Println("Kurban Borc Durum Hesap Kapandi")
				// return
				//burada kasa borçlu demek onun işlmeleri olacak
				//onun için ayrı bir sayfa include edilebilir bunun orneği evfada var
				kasaBorcu := access.kurbanApp.KasaBorcu(kurbanID)
				kurbanOdenenMiktar := access.kurbanApp.OdemelerToplami(kurbanID)
				viewData := pongo2.Context{
					"title":         "Kurban duzenleme",
					"durum":         "kasa kapandı",
					"post":          kurbanData,
					"referansData1": referansKisi1,
					// "referansData2":      referansKisi2,
					"kisiBilgileri":      kisiBilgileri,
					"tarih":              tarih,
					"ucretOdemesiDurumu": true,
					"kurbanOdenenMiktar": kurbanOdenenMiktar,
					"kasaBorcu":          kasaBorcu,
					"hayvanAtanmasi":     true,
					"kurbanDurum":        kurbanData.Durum,
					"csrf":               csrf.GetToken(c),
					"flashMsg":           flashMsg,
					"fileConfig":         accData,
					"medias":             mediaData,
					"actionHref":         action,
					"ID":                 kurbanID,
				}
				c.HTML(
					http.StatusOK,
					viewPathGkurban+"edit.html",
					viewData,
				)
			} else if kurbanData.BorcDurum == entity.KurbanBorcDurumTaksitOdemesi || kurbanData.BorcDurum == entity.KurbanBorcDurumIlkEklenenFiyat || kurbanData.BorcDurum == entity.KurbanBorcDurumFiyatManuelDegistirildi {
				fmt.Println("üçlü data ")
				kasaBorcu := access.kurbanApp.KasaBorcu(kurbanID)
				kurbanOdenenMiktar := access.kurbanApp.OdemelerToplami(kurbanID)
				kurbanSonKalanBorc := access.kurbanApp.KalanUcret(kurbanID)

				viewData := pongo2.Context{
					"title":         "Kurban duzenleme",
					"durum":         "false",
					"post":          kurbanData,
					"referansData1": referansKisi1,
					// "referansData2":      referansKisi2,
					"kisiBilgileri":      kisiBilgileri,
					"tarih":              tarih,
					"ucretOdemesiDurumu": false,
					"hayvanAtanmasi":     false,
					"kurbanSonKalanBorc": kurbanSonKalanBorc,
					"kurbanOdenenMiktar": kurbanOdenenMiktar,
					"kasaBorcu":          kasaBorcu,
					"kurbanDurum":        kurbanData.Durum,
					"csrf":               csrf.GetToken(c),
					"flashMsg":           flashMsg,
					"fileConfig":         accData,
					"medias":             mediaData,
					"actionHref":         action,
					"ID":                 kurbanID,
				}
				c.HTML(
					http.StatusOK,
					viewPathGkurban+"edit.html",
					viewData,
				)
			}

		} else {
			c.Redirect(http.StatusMovedPermanently, "/404")
		}

	} else {
		c.Redirect(http.StatusMovedPermanently, "/404")
	}
}

//Update data
func (access *Kurban) Update(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	kurban, id, _ := gkurbanModel("edit", c)

	var savePostError = make(map[string]string)
	action := c.DefaultQuery("action", "ekle")

	savePostError = kurban.Validate()

	if len(savePostError) == 0 {
		_, saveErr := access.kurbanApp.Update(&kurban)
		if saveErr != nil {
			savePostError = saveErr
		}

		if c.PostForm("kaydet") == "gruplaraDon" {
			c.Redirect(http.StatusMovedPermanently, "/admin/gruplar/")
			return
		} else {
			stncsession.SetFlashMessage("Kayıt başarı ile eklendi", "success", c)
			c.Redirect(http.StatusMovedPermanently, "/admin/kurban/edit/"+id)
		}

	} else {
		stncsession.SetFlashMessage("Zorunlu alanları lütfen doldurunuz", "danger", c)
	}

	flashMsg := stncsession.GetFlashMessage(c)

	idInt1, _ := strconv.Atoi(id)
	idInt := uint64(idInt1)
	kurbanDurum := access.kurbanApp.GetKurbanDurum(idInt)

	kurbanTuru := access.kurbanApp.GetKurbanTuru(idInt)

	var modulId int = 1
	var mediaData []entity.Media
	if kurbanTuru == 12 {
		modulId = 3
		// grupID := access.kurbanApp.GetKurbanGrupID(idInt)
		// hayvanID := access.kurbanApp.GetKurbanHayvanID(uint64(grupID))
		// mediaData, _ = access.MediaApp.GetAllforModul(modulId, hayvanID)
	}
	// else {
	// 	mediaData, _ = access.MediaApp.GetAll(modulId, int(idInt))
	// }

	accData := access.UploadConfig(modulId, "")

	kisi1 := stnccollection.StringtoUint64(c.PostForm("Kisi1"))
	referansKisi1 := stnccollection.StringtoUint64(c.PostForm("ReferansKisi1"))
	// referansKisi2 := stnccollection.StringtoUint64(c.PostForm("ReferansKisi2"))

	kisiBilgileri, _ := access.kisiApp.GetByID(kisi1)
	referansKisiData1, _ := access.kisiApp.GetByID(referansKisi1)
	// referansKisiData2, _ := access.kisiApp.GetByIDReferans(referansKisi2)

	viewData := pongo2.Context{
		"title":         "Kurban Düzenleme",
		"err":           savePostError,
		"medias":        mediaData,
		"csrf":          csrf.GetToken(c),
		"post":          kurban,
		"kisiBilgileri": kisiBilgileri,
		"referansData1": referansKisiData1,
		// "referansData2": referansKisiData2,
		"flashMsg":    flashMsg,
		"actionHref":  action,
		"kurbanDurum": kurbanDurum,
		"fileConfig":  accData,
	}

	c.HTML(
		http.StatusOK,
		viewPathGkurban+"edit.html",
		viewData,
	)
}

/**---API ModALBOX KISIMLAR --*/

//ReferansCreateModalBox create modalbox
func (access *Kurban) ReferansCreateModalBox(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	viewID := c.Param("viewID")
	adSoyad := c.Query("metin")
	viewData := pongo2.Context{
		"title":   "Kişi Ekleme",
		"viewID":  viewID,
		"adSoyad": adSoyad,
		"csrf":    csrf.GetToken(c),
	}
	c.HTML(
		http.StatusOK,
		viewPathGkurban+"referansCreateModalBox.html",
		viewData,
	)
}

//OdemeEkleCreateModalBox takistler
func (access *Kurban) OdemeEkleCreateModalBox(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	kurbanID := c.DefaultQuery("kurbanID", "0")
	viewData := pongo2.Context{
		"title":    "Kişi Ekleme",
		"csrf":     csrf.GetToken(c),
		"kurbanID": kurbanID,
	}
	c.HTML(
		http.StatusOK,
		viewPathGkurban+"odemeEkleCreateModalBox.html",
		viewData,
	)
}

//GrupLideriAta takistler
func (access *Kurban) GrupLideriAta(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	grupID := stnccollection.StringtoUint64(c.DefaultQuery("grupID", "0"))
	kurbanID := stnccollection.StringtoUint64(c.Param("kurbanID")) // TODO: bunun param olarak gelmesi veri güvenliği sorunu olablir mi #gfsecurity
	kurbanlarGrub, _ := access.kurbanApp.GetByGrupID(grupID)
	for no, _ := range kurbanlarGrub {
		kurbanlarGrubkurbanID := kurbanlarGrub[no].ID
		a := no
		a++
		access.kurbanApp.SetGrupLideri(kurbanlarGrubkurbanID, a)
	}
	access.kurbanApp.SetGrupLideri(kurbanID, 0)
	stncsession.SetFlashMessage("Grup Lideri Değiştirildi", "success", c)
	c.Redirect(http.StatusMovedPermanently, "/admin/gruplar")
}

//GrupLideriAta takistler
func (access *Kurban) VekaletDurumu(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	status := stnccollection.StringToint(c.DefaultQuery("status", "0"))
	kurbanID := stnccollection.StringtoUint64(c.Param("kurbanID"))
	access.kurbanApp.SetVekaletDurumu(kurbanID, status)
	stncsession.SetFlashMessage("Vekalet Durumu Değiştirildi", "success", c)
	c.Redirect(http.StatusMovedPermanently, "/admin/gruplar")
}

/***  POST MODEL   ***/
func gkurbanModel(formType string, c *gin.Context) (dataKurban entity.Kurban, idString string, err error) {

	id := c.PostForm("ID")

	idInt, _ := strconv.Atoi(id)

	var idN uint64
	var kisi1veri uint64

	idN = uint64(idInt)

	t := time.Now()

	dataKurban.ID = idN
	dataKurban.UserID = stncsession.GetUserID2(c)

	vekaletDurumu, _ := strconv.Atoi(c.PostForm("VekaletDurumu"))
	kurbanTuru, _ := strconv.Atoi(c.PostForm("KurbanTuru"))

	// referansKisi1 := stnccollection.StringtoUint64(c.PostForm("ReferansKisi1"))
	// referansKisi2 := stnccollection.StringtoUint64(c.PostForm("ReferansKisi2"))
	fmt.Println("carmı")
	fmt.Println(stnccollection.StringtoUint64(c.PostForm("Kisi1")))
	kisi1veri = stnccollection.StringtoUint64(c.PostForm("Kisi1"))
	if kisi1veri == 0 {
		kisi1veri = 1
	}
	dataKurban.UserID = stncsession.GetUserID2(c)
	dataKurban.KisiID = kisi1veri
	// dataKurban.ReferansKisi1 = referansKisi1
	// dataKurban.ReferansKisi2 = referansKisi2
	dataKurban.Agirlik, _ = strconv.Atoi(c.PostForm("Agirlik"))

	dataKurban.VekaletDurumu = vekaletDurumu
	dataKurban.KurbanTuru = kurbanTuru
	HayvanCinsi, _ := strconv.Atoi(c.PostForm("HayvanCinsi"))
	dataKurban.HayvanCinsi = HayvanCinsi

	slug := stnchelper.RandSlugV1(2) + stnccollection.IntToString(t.Second()) + stnchelper.RandSlugV1(2)

	dataKurban.Slug = slug

	// dataKurban.BorcDurum = entity.KurbanBorcDurumIlkEklenenFiyat
	// dataKurban.Durum = entity.KurbanDurumKurbanEklendiKurbanBayraminaAitDegil
	dataKurban.Aciklama = c.PostForm("Aciklama")

	if formType == "create" {
		var kurbanFiyati float64
		dataKurban.GrupID = stnccollection.StringtoUint64(c.PostForm("KurbanTuru"))
		kurbanFiyati, _ = stnccollection.StringToFloat64(c.PostForm("KurbanFiyati"))
		dataKurban.KurbanFiyati = kurbanFiyati
		dataKurban.Alacak = kurbanFiyati

		dataKurban.BorcDurum = entity.KurbanBorcDurumIlkEklenenFiyat
		dataKurban.Odemeler = []dto.Odemeler{{Aciklama: "İlk Eklenen Fiyat", KurbanFiyati: kurbanFiyati, VerilenUcret: 0, Alacak: kurbanFiyati, VerildigiTarih: time.Now(), BorcDurum: entity.OdemelerBorcDurumIlkEklenenFiyat}}
	}
	//bu kısımda kurban fiyatındaki değişimleri hesaplama yapar
	if formType == "edit" {
		db := repository.DB
		services, err1 := repository.RepositoriesInit(db)
		if err1 != nil {
			panic(err1)
		}

		Kurban := InitGkurban(services.Kurban, services.Kisiler, services.Media)

		// appOdeme := InitOdemeler(services.Kodemeler)
		// kurbanSonKalanBorc = appOdeme.OdemelerApp.KurbanSonKalanUcret(idN)
		odemelerToplami := Kurban.kurbanApp.OdemelerToplami(idN)
		kalanUcret := Kurban.kurbanApp.KalanUcret(idN)
		kurbanDurum := Kurban.kurbanApp.GetKurbanDurum(idN)

		fmt.Println("odemelerToplami")
		fmt.Println(odemelerToplami)

		fmt.Println("kalanUcret")
		fmt.Println(kalanUcret)

		kurbanFiyati, _ := strconv.ParseFloat(c.PostForm("KurbanFiyati"), 64)
		kurbanFiyatiOLD := Kurban.kurbanApp.KurbanFiyati(idN)

		fmt.Println("kurbanFiyati")
		fmt.Println(kurbanFiyati)
		// dataKurban.KurbanFiyati, _ = stnccollection.StringToFloat64(c.PostForm("KurbanFiyati"))
		dataKurban.GrupID = stnccollection.StringtoUint64(c.PostForm("GrupID"))

		kurbanGrupLideri := Kurban.kurbanApp.GetGrupLideri(idN)

		dataKurban.GrupLideri = kurbanGrupLideri

		if kurbanFiyati != kurbanFiyatiOLD {
			//kasa borçlu durumda
			if kurbanFiyati > kalanUcret {
				fmt.Println("kasa borçlu")
				dataKurban.Bakiye = odemelerToplami
				kasaBorcu := odemelerToplami - kurbanFiyati
				dataKurban.Borc = kasaBorcu
				dataKurban.KurbanFiyati = kurbanFiyati

				dataKurban.BorcDurum = entity.KurbanBorcDurumKasaBorcluDurumda
				dataKurban.Durum = entity.KurbanDurumGrupOlusmusKesimlikHayvaniVar
				// TODO: user id eklenecek

				//Kasa BORÇLU:  Fiyat Düşürüldü  KF>KÜ kurbanFiyati > kalanUcret -fiyat elle değiştirildi"
				dataKurban.Odemeler = []dto.Odemeler{{Aciklama: "Kasa BORÇLU:  Fiyat Düşürüldü KF>KÜ -fiyat elle değiştirildi", Bakiye: 0, VerilenUcret: 0, KurbanFiyati: kurbanFiyati, Alacak: 0, Borc: kasaBorcu, VerildigiTarih: time.Now(), Durum: entity.OdemelerDurumGrupOlusmusKesimlikHayvaniVar, BorcDurum: entity.OdemelerBorcDurumKasaBorcluDurumda}}
			}
			//TODO: yukarıdakinden ne farkı var
			if kurbanFiyati < kalanUcret {
				fmt.Println("girilen fiyat küçük ")

				dataKurban.Bakiye = odemelerToplami
				dataKurban.Borc = odemelerToplami - kurbanFiyati
				dataKurban.Alacak = 0
				dataKurban.KurbanFiyati = kurbanFiyati

				dataKurban.BorcDurum = entity.KurbanBorcDurumKasaBorcluDurumda
				dataKurban.Durum = entity.KurbanDurumGrupOlusmusKesimlikHayvaniVar

				//	dataKurban.Odemeler = []dto.Odemeler{{Aciklama: "Kasa BORÇLU: Fiyat Düşürüldü  kurbanFiyati < kalanUcret -fiyat elle değiştirildi", Bakiye: 0, VerilenUcret: kurbanFiyati, KurbanFiyati: kurbanFiyati, Alacak: 0, Borc: kalanUcret, VerildigiTarih: time.Now(), Durum: entity.OdemelerDurumGrupOlusmusKesimlikHayvaniVar, BorcDurum: entity.OdemelerBorcDurumKasaBorcluDurumda}} //TODO: user id eklenecek

				dataKurban.Odemeler = []dto.Odemeler{{Aciklama: "Kasa BORÇLU: Fiyat Düşürüldü  KF<KÜ -fiyat elle değiştirildi", Bakiye: 0, VerilenUcret: 0, KurbanFiyati: kurbanFiyati, Alacak: 0, Borc: kalanUcret, VerildigiTarih: time.Now(), Durum: entity.OdemelerDurumGrupOlusmusKesimlikHayvaniVar, BorcDurum: entity.OdemelerBorcDurumKasaBorcluDurumda}} //TODO: user id eklenecek
			}
			if kurbanFiyati > odemelerToplami {
				//TODO: odemesi 0 ise ödemeler logdurumunda sorun oluyor
				fmt.Println("girilen fiyat(degisen) büyük ")
				farkUcret := kurbanFiyati - odemelerToplami
				dataKurban.KurbanFiyati = kurbanFiyati
				dataKurban.Alacak = farkUcret
				dataKurban.Bakiye = odemelerToplami

				dataKurban.BorcDurum = entity.KurbanBorcDurumTaksitOdemesi
				dataKurban.Durum = entity.KurbanDurumGrupOlusmusKesimlikHayvaniVar
				var status int = entity.OdemelerBorcDurumFiyatManuelDegistirildiDusukFiyat
				var aciklama string = "Kurban Fiyatı Düşürüldü -fiyat elle değiştirildi"
				if kurbanFiyati > kurbanFiyatiOLD {
					status = entity.OdemelerBorcDurumFiyatManuelDegistirildiYuksekFiyat
					aciklama = "Kurban Fiyati Artırıldı -fiyat elle değiştirildi"
				}
				dataKurban.Odemeler = []dto.Odemeler{{Aciklama: aciklama, VerilenUcret: 0, KurbanFiyati: kurbanFiyati, Bakiye: odemelerToplami, Alacak: farkUcret, VerildigiTarih: time.Now(), Durum: entity.OdemelerDurumGrupOlusmusKesimlikHayvaniVar, BorcDurum: status}} //TODO: user id eklenecek
			}
			if kurbanFiyati == odemelerToplami {
				//TODO: odemesi 0 ise ödemeler logdurumunda sorun oluyor
				fmt.Println("Ödeme tamamlandı ")
				farkUcret := kurbanFiyati - odemelerToplami
				dataKurban.KurbanFiyati = kurbanFiyati
				dataKurban.Alacak = 0
				dataKurban.Bakiye = odemelerToplami
				dataKurban.BorcDurum = entity.KurbanBorcDurumHesapKapandi
				dataKurban.Durum = entity.KurbanDurumGrupOlusmusKesimlikHayvaniVar
				dataKurban.Odemeler = []dto.Odemeler{{Aciklama: "Ödeme tamamlandı -fiyat elle değiştirildi", VerilenUcret: 0, Bakiye: odemelerToplami, KurbanFiyati: kurbanFiyati, Alacak: farkUcret, VerildigiTarih: time.Now(), Durum: entity.OdemelerDurumHesapKapandi, BorcDurum: entity.OdemelerBorcDurumFiyatManuelDegistirildiEsitFiyat}} //TODO: user id eklenecek
			}
		}
		if kurbanFiyati == kurbanFiyatiOLD {
			fmt.Println("eski ve yni fiyat eşit ")
			fmt.Println("kurbanFiyati")
			fmt.Println(kurbanFiyati)
			if kurbanFiyati == 0 {
				kurbanFiyati = -1
			}
			fmt.Println(kurbanFiyati)
			dataKurban.KurbanFiyati = kurbanFiyati
			dataKurban.Alacak = kurbanFiyati
			dataKurban.Bakiye = odemelerToplami
			dataKurban.BorcDurum = entity.KurbanBorcDurumIlkEklenenFiyat
			dataKurban.Durum = entity.KurbanDurumGrupOlusmusKesimlikHayvaniVar
		}
		if kurbanDurum == entity.KurbanDurumGrupOlusmusHayvanYok {
			fmt.Println("hayvan yok ")
			fmt.Println(kurbanFiyati)
			if kurbanFiyati == 0 {
				kurbanFiyati = -1
			}

			dataKurban.KurbanFiyati = kurbanFiyati
			dataKurban.Borc = 0
			dataKurban.Alacak = 0
			dataKurban.Bakiye = odemelerToplami
			dataKurban.BorcDurum = entity.KurbanBorcDurumIlkEklenenFiyat
			dataKurban.Durum = entity.KurbanDurumGrupOlusmusHayvanYok
		}
	}
	return dataKurban, id, nil
}

//Upload
func (access *Kurban) Upload(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	IDint := stnccollection.StringToint(c.Param("ID"))
	filenameForm, _ := c.FormFile("file")
	kurbanTuru := access.kurbanApp.GetKurbanTuru(uint64(IDint))
	var modulId int = 1
	if kurbanTuru == 12 {
		modulId = 3
	}
	accData := access.UploadConfig(modulId, "")
	uploadSize := stnccollection.StringToint(accData["uploadSize"])
	fileType := accData["fileType"]
	fileTypes := strings.Split(fileType, ",")
	maxFiles := stnccollection.StringToint(accData["maxFiles"])
	uploadPath := accData["uploadPath"]
	upl := stncupload.FileUpload{
		UploadPath: uploadPath,
		UploadSize: uploadSize,
		MaxFiles:   maxFiles,
		// Types:      []string{"video/mp4", "image/jpeg", "image/jpg", "image/gif", "image/png", "video/webm", "application/pdf"},
		// Types: []string{},
		Types: fileTypes,
	}

	var up stncupload.UploadFileInterface = upl

	up.InitUploader(c, IDint, filenameForm, accData)

}

//MediaDelete Delete file data
func (access *Kurban) MediaDelete(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	if ID, err := strconv.ParseUint(c.Param("ID"), 10, 64); err == nil {

		kurbanTuru := access.kurbanApp.GetKurbanTuru(uint64(ID))
		var modulId int = 1
		if kurbanTuru == 12 {
			modulId = 3
		}
		accData := access.UploadConfig(modulId, "")
		uploadPath := accData["uploadPath"]
		stncupload.NewFileUpload().MediaDelete(c, ID, uploadPath)
	}
}

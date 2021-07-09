package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stncdatetime"
	"stncCms/app/domain/helpers/stncsession"
	"stncCms/app/domain/repository"
	"stncCms/app/services"
	"strconv"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//Odemeler constructor
type Odemeler struct {
	OdemelerApp services.OdemelerAppInterface
}

const viewPathOdemeler = "admin/kurbanlar/"

//InitOdemeler post controller constructor
func InitOdemeler(odemeApp services.OdemelerAppInterface) *Odemeler {
	return &Odemeler{
		OdemelerApp: odemeApp,
	}
}

/*
func InlineRepo() *Odemeler {
	return &Odemeler{
		OdemelerApp: repository.OdemelerRepositoryInit(db),
	}
}
*/

//Create all list
func (access *Odemeler) Create(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)

	viewData := pongo2.Context{
		"title": "Kurban  Ekleme",
		"csrf":  csrf.GetToken(c),
	}
	c.HTML(
		http.StatusOK,
		viewPathOdemeler+"create.html",
		viewData,
	)
}

//Store save method
func (access *Odemeler) Store(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)

	var KurbanID uint64

	KurbanID, _ = strconv.ParseUint(c.PostForm("kurbanID"), 10, 64)

	var kalanKurbanFiyati float64
	if posts, err := access.OdemelerApp.LastPrice(KurbanID); err == nil {
		kalanKurbanFiyati = posts.Alacak
		fmt.Printf("%+v\n", posts)
	}

	if odeme, _, errorR := odemelerModel(kalanKurbanFiyati, c); errorR == nil {
		var savePostError = make(map[string]string)
		savePostError = odeme.Validate()

		if len(savePostError) == 0 {
			saveData, saveErr := access.OdemelerApp.Save(&odeme)
			if saveErr != nil {
				savePostError = saveErr
			}
			lastID := strconv.FormatUint(uint64(saveData.ID), 10)
			stncsession.SetFlashMessage("Kayıt başarı ile eklendi", "success", c)

			formtype := c.PostForm("formtype")
			if formtype == "inline" {
				kurbanID := c.PostForm("kurbanID")
				c.Redirect(http.StatusMovedPermanently, "/admin/kurban/edit/"+kurbanID)
			} else {
				c.Redirect(http.StatusMovedPermanently, "/admin/kurban/edit/"+lastID)
			}

			return
		}

		viewData := pongo2.Context{
			"title": "Kurban Ekleme",
			"csrf":  csrf.GetToken(c),
			"err":   savePostError,
			"post":  odeme,
		}
		c.HTML(
			http.StatusOK,
			viewPathOdemeler+"create.html",
			viewData,
		)
	}
	//  else {
	// 	if errorR.Error() == "Ücret Büyük" {
	// 		c.JSON(http.StatusOK, "büyük data")
	// 		return
	// 	}
	// }

}

//referansEkleAjax save method
func (access *Odemeler) OdemeEkleAjax(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)

	var KurbanID uint64

	KurbanID, _ = strconv.ParseUint(c.PostForm("kurbanID"), 10, 64)

	var kalanKurbanFiyati float64
	if posts, err := access.OdemelerApp.LastPrice(KurbanID); err == nil {
		kalanKurbanFiyati = posts.Alacak
		fmt.Printf("%+v\n", posts)
	}

	db := repository.DB
	appKurban := repository.KurbanRepositoryInit(db)
	var kisiID int
	appKurban.GetKurbanKisiVarmi(KurbanID, &kisiID)
	if kisiID == 1 {
		// sahte veri girişi yani kişi atanmamış kurbana ödeme yapmaya çalışıyor  TODO: bunun loglanması lazım
		viewData := pongo2.Context{
			"title":  "Kurban Ekleme",
			"csrf":   csrf.GetToken(c),
			"status": "err",
			"err":    "fk", // sahte veri girişi TODO: bunun loglanması lazım
			"errMsg": "beklenmeyen bir hata oluştu",
		}
		c.JSON(http.StatusOK, viewData)
		return
	} else {
		if odeme, _, errorR := odemelerModel(kalanKurbanFiyati, c); errorR == nil {
			var savePostError = make(map[string]string)
			savePostError = odeme.Validate()
			fmt.Printf("%+v\n", odeme)
			if len(savePostError) == 0 {
				_, saveErr := access.OdemelerApp.Save(&odeme)

				if saveErr != nil {
					savePostError = saveErr
				}
				viewData := pongo2.Context{
					"title":  "Kurban Ekleme",
					"csrf":   csrf.GetToken(c),
					"err":    savePostError,
					"status": "ok",
					"path":   "/admin/kurban/edit/" + c.PostForm("kurbanID"),
					"id":     c.PostForm("kurbanID"),
					"post":   odeme,
				}
				c.JSON(http.StatusOK, viewData)
			}
		} else {
			if errorR.Error() == "Ücret Büyük" {
				// fmt.Println("kalanKurbanFiyati")
				// fmt.Println(kalanKurbanFiyati)

				viewData := pongo2.Context{
					"title":  "Kurban Ekleme",
					"csrf":   csrf.GetToken(c),
					"status": "err",
					"err":    "büyük data",
					"errMsg": kalanKurbanFiyati,
				}
				c.JSON(http.StatusOK, viewData)
				return
			}
		}
	}
}

//Edit genel kurban düzenleme işler
func (access *Odemeler) Edit(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	if gID, err := strconv.ParseUint(c.Param("gID"), 10, 64); err == nil {
		if posts, err := access.OdemelerApp.GetByID(gID); err == nil {

			viewData := pongo2.Context{
				"title": "Ödeme duzenleme",
				"post":  posts,

				"csrf":     csrf.GetToken(c),
				"flashMsg": flashMsg,
			}
			c.HTML(
				http.StatusOK,
				viewPathOdemeler+"edit.html",
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
func (access *Odemeler) Update(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	//TODO: burada sıfır olması update edit de sıkıntı çıkarabilir

	kurban, id, _ := odemelerModel(0, c)

	// if odemeERROR != nil {
	// 	return "odemeERROR"
	// }

	var savePostError = make(map[string]string)

	savePostError = kurban.Validate()

	if len(savePostError) == 0 {

		_, saveErr := access.OdemelerApp.Update(&kurban)
		if saveErr != nil {
			savePostError = saveErr
		}
		stncsession.SetFlashMessage("Kayıt başarı ile düzenlendi", "success", c)
		c.Redirect(http.StatusMovedPermanently, "/admin/kurban/edit/"+id)
		return
	}
	viewData := pongo2.Context{
		"title": "Ödeme Düzenleme",
		"err":   savePostError,
		"csrf":  csrf.GetToken(c),
		"post":  kurban,
	}
	c.HTML(
		http.StatusOK,
		viewPathOdemeler+"edit.html",
		viewData,
	)
}

/***  DATA MODEL   ***/

func odemelerModel(kurbanSonKalanUcret float64, c *gin.Context) (odemelerData entity.Odemeler, idStr string, err error) {
	id := c.PostForm("ID")
	idInt, _ := strconv.Atoi(id)
	var idN uint64
	idN = uint64(idInt)
	odemelerData.ID = idN

	var kurbanID uint64
	kurbanID, _ = strconv.ParseUint(c.PostForm("kurbanID"), 10, 64)

	db := repository.DB
	appKurban := repository.KurbanRepositoryInit(db)
	appOdeme := repository.OdemelerRepositoryInit(db)

	// KurbanFiyati := appKurban.KurbanlarApp.KurbanFiyati(KurbanID)
	// kurbanOdenenMiktar := appOdeme.OdemelerApp.OdemelerToplami(KurbanID)

	var verilenUcret float64
	var err2 error

	if verilenUcret, err2 = strconv.ParseFloat(c.PostForm("VerilenUcret"), 64); err2 == nil {
		odemelerData.VerilenUcret = verilenUcret
	}
	aciklama := c.PostForm("Aciklama")

	kurbanDurum := appKurban.GetKurbanDurum(kurbanID)
	kurbanFiyati := appKurban.KurbanFiyati(kurbanID)
	fmt.Printf("%+v\n", kurbanDurum)

	//TODO: kesimi tamamnlanış durumu en önemli kriter

	var odemelerToplami float64
	odemelerToplami = appOdeme.OdemelerToplami(kurbanID) //TODO: hayvasız gruplarda bunu test et sorun olcak mı

	odemelerToplami = odemelerToplami + verilenUcret

	//Grup Olusmus Hayvan Yok ise
	if kurbanDurum == entity.KurbanDurumGrupOlusmusHayvanYok {

		odemelerData.Alacak = 0
		odemelerData.KurbanFiyati = kurbanFiyati
		odemelerData.Borc = 0
		odemelerData.Bakiye = odemelerToplami
		odemelerData.Aciklama = "Kurban atanmadı, kapora eklendi: " + c.PostForm("VerilenUcret") + " TL"
		odemelerData.Durum = entity.OdemelerDurumGrupOlusmusHayvanYok
		odemelerData.BorcDurum = entity.OdemelerBorcDurumKaporaOdemesiHayvanBos
		appKurban.SetKurbanBakiyeUpdate(kurbanID, odemelerToplami)
		appKurban.SetKurbanBorcDurumUpdate(kurbanID, entity.KurbanBorcDurumTaksitOdemesi)
		appKurban.SetKurbanKalanUcretiUpdate(kurbanID, 0)
	}

	//Grup Olusmus Kesimlik Hayvani Var
	if kurbanDurum == entity.KurbanDurumGrupOlusmusKesimlikHayvaniVar {
		if verilenUcret > kurbanSonKalanUcret {
			fmt.Println("buyukkkk")
			return odemelerData, "0", errors.New("Ücret Büyük")
		}
		kalanUcret := kurbanSonKalanUcret - verilenUcret
		//belki ilerde sorun olursa acarız
		if odemelerToplami < kurbanSonKalanUcret {
			fmt.Println("girer")
			odemelerData.Alacak = kalanUcret
			odemelerData.Borc = 0
			odemelerData.Bakiye = odemelerToplami
			odemelerData.Aciklama = "Taksit Eklendi: " + c.PostForm("VerilenUcret") + " ₺"
			odemelerData.Durum = entity.OdemelerDurumGrupOlusmusKesimlikHayvaniVar
			odemelerData.BorcDurum = entity.OdemelerBorcDurumTaksitOdemesi
			appKurban.SetKurbanBakiyeUpdate(kurbanID, odemelerToplami)
			appKurban.SetKurbanBorcDurumUpdate(kurbanID, entity.KurbanBorcDurumTaksitOdemesi)
			appKurban.SetKurbanKalanUcretiUpdate(kurbanID, kalanUcret)
		}

		if kalanUcret == odemelerToplami {
			fmt.Println("hesap kapandı")
			odemelerData.Alacak = 0
			odemelerData.Borc = 0
			odemelerData.Bakiye = odemelerToplami
			odemelerData.Aciklama = "Taksit Eklendi: " + c.PostForm("VerilenUcret") + " ₺ / Hesap eşit"
			odemelerData.Durum = entity.OdemelerDurumGrupOlusmusKesimlikHayvaniVar
			odemelerData.BorcDurum = entity.OdemelerBorcDurumHesapKapandi
			appKurban.SetKurbanBakiyeUpdate(kurbanID, odemelerToplami)
			appKurban.SetKurbanBorcDurumUpdate(kurbanID, entity.KurbanBorcDurumHesapKapandi)
			appKurban.SetKurbanKasaBorcuUpdate(kurbanID, 0)
			appKurban.SetKurbanKalanUcretiUpdate(kurbanID, 0)
		}
		if kalanUcret == 0 {
			fmt.Println("hesap kapandı")
			odemelerData.Alacak = 0
			odemelerData.Borc = 0
			odemelerData.Bakiye = odemelerToplami
			odemelerData.Aciklama = "Taksit Eklendi: " + c.PostForm("VerilenUcret") + " ₺ / Hesap Kapandı"
			odemelerData.Durum = entity.OdemelerDurumGrupOlusmusKesimlikHayvaniVar
			odemelerData.BorcDurum = entity.OdemelerBorcDurumHesapKapandi
			odemelerData.Bakiye = odemelerToplami
			appKurban.SetKurbanBakiyeUpdate(kurbanID, odemelerToplami)
			appKurban.SetKurbanBorcDurumUpdate(kurbanID, entity.KurbanBorcDurumHesapKapandi)
			appKurban.SetKurbanKalanUcretiUpdate(kurbanID, 0)
		} else {
			odemelerData.Alacak = kalanUcret
			odemelerData.Borc = 0
			odemelerData.Bakiye = odemelerToplami
			odemelerData.Aciklama = "Taksit Eklendi: " + c.PostForm("VerilenUcret") + " ₺"
			odemelerData.Durum = entity.OdemelerDurumGrupOlusmusKesimlikHayvaniVar
			odemelerData.BorcDurum = entity.OdemelerBorcDurumTaksitOdemesi
			appKurban.SetKurbanBakiyeUpdate(kurbanID, odemelerToplami)
			appKurban.SetKurbanBorcDurumUpdate(kurbanID, entity.KurbanBorcDurumTaksitOdemesi)
			appKurban.SetKurbanKalanUcretiUpdate(kurbanID, kalanUcret)
		}
	}

	// kurban eklenmiş ama kurban bayramına ait değil yani direk kurban girişi
	if kurbanDurum == entity.KurbanDurumKurbanEklendiKurbanBayraminaAitDegil {
		if verilenUcret > kurbanSonKalanUcret {
			fmt.Println("buyukkkk")
			return odemelerData, "0", errors.New("Ücret Büyük")
		}
		kalanUcret := kurbanSonKalanUcret - verilenUcret

		if kalanUcret == 0 {
			//fmt.Println("hesap kapandı")
			odemelerData.Alacak = 0
			odemelerData.Borc = 0
			odemelerData.Bakiye = odemelerToplami
			odemelerData.Aciklama = "Taksit Eklendi: " + c.PostForm("VerilenUcret") + " ₺ / Hesap Kapandı"
			odemelerData.Durum = entity.KurbanDurumKurbanEklendiKurbanBayraminaAitDegil
			odemelerData.BorcDurum = entity.OdemelerBorcDurumHesapKapandi
			appKurban.SetKurbanBakiyeUpdate(kurbanID, odemelerToplami)
			appKurban.SetKurbanBorcDurumUpdate(kurbanID, entity.KurbanBorcDurumHesapKapandi)
			appKurban.SetKurbanKalanUcretiUpdate(kurbanID, 0)
		} else {
			odemelerData.Alacak = kalanUcret
			odemelerData.Borc = 0
			odemelerData.Bakiye = odemelerToplami
			odemelerData.Aciklama = "Taksit Eklendi: " + c.PostForm("VerilenUcret") + " ₺"
			odemelerData.Durum = entity.KurbanDurumKurbanEklendiKurbanBayraminaAitDegil
			odemelerData.BorcDurum = entity.OdemelerBorcDurumTaksitOdemesi
			appKurban.SetKurbanBakiyeUpdate(kurbanID, odemelerToplami)
			appKurban.SetKurbanBorcDurumUpdate(kurbanID, entity.KurbanBorcDurumTaksitOdemesi)
			appKurban.SetKurbanKalanUcretiUpdate(kurbanID, kalanUcret)
		}
	}

	if aciklama != "" {
		odemelerData.Aciklama = aciklama
	}
	odemelerData.KurbanID = kurbanID
	odemelerData.Makbuz = c.PostForm("Makbuz")
	odemelerData.UserID = stncsession.GetUserID2(c)
	odemelerData.VerildigiTarih = time.Now()
	odemelerData.KurbanFiyati = kurbanFiyati
	return odemelerData, id, nil
}

/*

	//buraya erişim gkurbancontoller da index içinde var
	//toplamOdeme := access.gKurbanApp.OdemelerToplami(gID)
	// verilenUcret, _ := strconv.ParseFloat(c.PostForm("VerilenUcret"), 64)//gereksiz usstte var
	var hesaplaVerilenUcret float64
	if kurbanDurum != entity.KurbanDurumGrupOlusmusHayvanYok {
		if kurbanDurum != entity.KurbanDurumKurbanKesimiTamamlanmis {
			if verilenUcret > kurbanSonKalanUcret {
				fmt.Println("buyukkkk")
				return data, "0", errors.New("Ücret Büyük")
			}
		}
		if kurbanDurum != entity.KurbanDurumKurbanKesimiTamamlanmis {
			hesaplaVerilenUcret = kurbanSonKalanUcret - verilenUcret
			if hesaplaVerilenUcret == 0 {
				data.KalanUcret = -1
				data.Durum = entity.OdemelerDurumHesapKapandi
				data.BorcDurum = entity.OdemelerBorcDurumHesapKapandi
				appKurban.KurbanlarApp.SetKurbanBorcDurumUpdate(KurbanID, entity.KurbanBorcDurumHesapKapandi)
			} else {
				data.KalanUcret = hesaplaVerilenUcret
			}
		}
		//TODO: bu ayrı yerde depency yada observer pattern olmalı
		appKurban.KurbanlarApp.SetKurbanKalanUcretiUpdate(KurbanID, hesaplaVerilenUcret)

		odenenToplamMiktar := KurbanFiyati - hesaplaVerilenUcret

		fmt.Println("kurbanUcreti")
		fmt.Println(kurbanUcreti)

		fmt.Println("odenenToplamMiktar")
		fmt.Println(odenenToplamMiktar)

		//TODO: ustteki todo ile aynı olmalı
		//tek bir fonksiyon olmalı aslında gruplarda değiştir , kurban fiyatı değişikliği durumlarda da bunlar kullanılabilir olmalı
		appKurban.KurbanlarApp.SetKurbanBakiyeUpdate(KurbanID, odenenToplamMiktar)
	} else {
		fmt.Println("verilenUcret")
		fmt.Println(verilenUcret)
		data.KalanUcret = verilenUcret
		appKurban.KurbanlarApp.SetKurbanBakiyeUpdate(KurbanID, verilenUcret)
		appKurban.KurbanlarApp.SetKurbanKalanUcretiUpdate(KurbanID, verilenUcret)
	}

	data.Aciklama = c.PostForm("Aciklama")
	data.Makbuz = c.PostForm("Makbuz")
	data.KurbanID = KurbanID

	data.KasaBorcu = -1

	data.UserID = stncsession.GetUserID2(c)
	data.BorcDurum = entity.OdemelerBorcDurumTaksitOdemesi
	data.VerildigiTarih = time.Now()
	return data, id, nil

*/
func (access *Odemeler) Makbuz(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	var tarih stncdatetime.Inow
	//TODO: odeme durumunun 2,4,6 olmasına bakılacak #secure
	flashMsg := stncsession.GetFlashMessage(c)
	if gID, err := strconv.ParseUint(c.Param("gID"), 10, 64); err == nil {
		if makbuz, err := access.OdemelerApp.GetOdemeRelation(gID); err == nil {

			empJSON, err := json.MarshalIndent(makbuz, "", "  ")
			if err != nil {
				log.Fatalf(err.Error())
			}
			fmt.Printf("MarshalIndent funnction output\n %s\n", string(empJSON))
			currentTime := time.Now()
			verildigitarih := tarih.OnlyDate(currentTime.String())
			viewData := pongo2.Context{
				"title":    "Makbuz Kesme",
				"makbuz":   makbuz,
				"csrf":     csrf.GetToken(c),
				"flashMsg": flashMsg,
				"tarih":    verildigitarih,
			}
			c.HTML(
				http.StatusOK,
				viewPathOdemeler+"makbuzKes.html",
				viewData,
			)
		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

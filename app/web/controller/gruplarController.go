package controller

import (
	"fmt"
	"math"
	"net/http"
	"stncCms/app/domain/dto"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/msgedit"
	"stncCms/app/domain/helpers/stnccollection"
	"stncCms/app/domain/helpers/stnchelper"
	"stncCms/app/domain/helpers/stncsession"
	"stncCms/app/domain/repository"
	"stncCms/app/services"
	"time"

	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/leekchan/accounting"
	csrf "github.com/utrack/gin-csrf"
)

//Gruplar constructor
type Gruplar struct {
	HayvanBilgisiApp services.HayvanBilgisiAppInterface
	KurbanApp        services.KurbanAppInterface
	GruplarApp       services.GruplarAppInterface
	OptionsApp       services.OptionsAppInterface
	OdemelerApp      services.OdemelerAppInterface
	kisiApp          services.KisilerAppInterface
}

const viewPathGruplar = "admin/gruplar/"

//InitGruplar post controller constructor
func InitGruplar(hayvanBilgisiApp services.HayvanBilgisiAppInterface, kurbanApp services.KurbanAppInterface, gruplarbayramiApp services.GruplarAppInterface, optionsApp services.OptionsAppInterface, odemelerApp services.OdemelerAppInterface, kisiApp services.KisilerAppInterface) *Gruplar {
	return &Gruplar{
		HayvanBilgisiApp: hayvanBilgisiApp,
		KurbanApp:        kurbanApp,
		GruplarApp:       gruplarbayramiApp,
		OptionsApp:       optionsApp,
		OdemelerApp:      odemelerApp,
		kisiApp:          kisiApp,
	}
}

func GrupDataGenerate(access *Gruplar) []dto.GruplarExcelandIndex {
	// db := repository.DB
	// HayvanBilgisiApp := repository.HayvanBilgisiRepositoryInit(db)
	// kisiApp := repository.KisilerRepositoryInit(db)
	// GruplarApp := repository.GruplarRepositoryInit(db)

	var veriler []dto.GruplarExcelandIndex
	veriler, _ = access.GruplarApp.GetByAllRelations()
	for num, v := range veriler {
		fmt.Println(access.HayvanBilgisiApp.GetKupeNo(v.HayvanBilgisiID))
		veriler[num].KupeNo = access.HayvanBilgisiApp.GetKupeNo(v.HayvanBilgisiID)
		fmt.Println("fiyat")
		fmt.Printf("%+v\n", v.ToplamKurbanFiyati)
		if v.ToplamKurbanFiyati == 1 {
			veriler[num].ToplamKurbanFiyati = 0
		} else {
			veriler[num].ToplamKurbanFiyati = v.ToplamKurbanFiyati
		}
		if v.AgirlikTipi == 1 {
			veriler[num].GrupIsoTopeName = "dusuk"
			veriler[num].GrupIsoTopeTRname = "Düşük"
			veriler[num].GrupIsoTopeAlert = "danger"
		} else if v.AgirlikTipi == 2 {
			veriler[num].GrupIsoTopeName = "orta"
			veriler[num].GrupIsoTopeTRname = "Orta"
			veriler[num].GrupIsoTopeAlert = "warning"
		} else if v.AgirlikTipi == 3 {
			veriler[num].GrupIsoTopeName = "yuksek"
			veriler[num].GrupIsoTopeTRname = "Yüksek"
			veriler[num].GrupIsoTopeAlert = "success"
		}

		veriler[num].ToplamOdemeler = access.GruplarApp.ToplamOdemeler(v.ID)
		veriler[num].KalanBorclar = access.GruplarApp.KalanBorclar(v.ID)
		veriler[num].KasaBorcu = access.GruplarApp.KasaBorcu(v.ID)

		var kisilerList = []dto.KurbanListForGrouplar{}
		kisilerList, _ = access.KurbanApp.GetAllKurbanAndKisiler(int(v.ID))
		// fmt.Printf("%+v\n", kisilerList)
		for no, kisi := range kisilerList {
			if kisi.RefKisiID != 0 {
				referansKisi, _ := access.kisiApp.GetByID(kisi.RefKisiID)
				kisilerList[no].ReferansID = referansKisi.ID
				kisilerList[no].ReferansAdSoyad = referansKisi.AdSoyad
				kisilerList[no].ReferansTelefon = referansKisi.Telefon
			}
		}
		veriler[num].KurbanKisiList = kisilerList
		// var post = entity.Post{}
		// var kurbandataKisi *entity.Kisiler
		kisiBasiDusenFiyatAraHesaplamaKusuratli := v.ToplamKurbanFiyati / float64(veriler[num].HissedarAdet)
		kisiBasiDusenFiyat := stnccollection.SayiYuvarla(math.Ceil(stnccollection.ToFixedDecimal(kisiBasiDusenFiyatAraHesaplamaKusuratli, 2)))
		veriler[num].KisiBasiDusenHisseFiyati = kisiBasiDusenFiyat
	}
	//#json formatter #stncjson
	// empJSON, err := json.MarshalIndent(veriler, "", "  ")
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	// fmt.Printf("MarshalIndent funnction output\n %s\n", string(empJSON))
	return veriler
}

//Index all list
func (access *Gruplar) Index(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	// whatsAppMsg := access.OptionsApp.GetOption("whatsAppMsg")
	var msg msgedit.Msg

	veriler := GrupDataGenerate(access)
	viewData := pongo2.Context{
		"paginator":     paginator,
		"title":         "Grup listesi",
		"grupBilgileri": veriler,
		"flashMsg":      flashMsg,
		"msg":           msg,
	}

	c.HTML(
		http.StatusOK,
		viewPathGruplar+"index.html",
		viewData,
	)
}

//Index all list
func (access *Gruplar) IndexAta(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	// fmt.Printf("%+v\n", veriler)
	veriler := GrupDataGenerate(access)
	viewData := pongo2.Context{
		"paginator":     paginator,
		"title":         "Grup listesi",
		"grupBilgileri": veriler,
		"flashMsg":      flashMsg,
		"csrf":          csrf.GetToken(c),
	}

	c.HTML(
		http.StatusOK,
		viewPathGruplar+"indexAta.html",
		viewData,
	)
}

//yerDegistir  yer değişikliği yapar
func (access *Gruplar) YerDegistir(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	targetGrupID := stnccollection.StringtoUint64(c.PostForm("targetGrupID"))
	targetKurbanID := stnccollection.StringtoUint64(c.PostForm("targetKurbanID"))
	targetGruplideri := stnccollection.StringToint(c.PostForm("targetGruplideri"))
	sourceGrupID := stnccollection.StringtoUint64(c.PostForm("sourceGrupID"))
	sourceKurbanID := stnccollection.StringtoUint64(c.PostForm("sourceKurbanID"))
	sourceGruplideri := stnccollection.StringToint(c.PostForm("sourceGruplideri"))

	access.GruplarApp.SetGrupID(sourceGrupID, targetKurbanID, sourceGruplideri)
	access.GruplarApp.SetGrupID(targetGrupID, sourceKurbanID, targetGruplideri)

	stncsession.SetFlashMessage("Kayıt başarı ile eklendi", "success", c)
	c.Redirect(http.StatusMovedPermanently, "/admin/gruplar")

}

//Create all list
func (access *Gruplar) Create(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	hayvanlar, _ := access.HayvanBilgisiApp.GetAllFindDurum(entity.HayvanBilgisiDurumHayvanBosta)
	hisseAdeti := access.OptionsApp.GetOption("hisse_adeti")
	kurbanYili := access.OptionsApp.GetOption("kurban_yili")
	otomatikSira2021 := access.OptionsApp.GetOption("otomatik_sira_buyukbas_2021")

	viewData := pongo2.Context{
		"title":            "Gruplar  Ekleme",
		"csrf":             csrf.GetToken(c),
		"hayvanlar":        hayvanlar,
		"hisseAdeti":       hisseAdeti,
		"kurbanYili":       kurbanYili,
		"otomatikSira2021": otomatikSira2021,
		"flashMsg":         flashMsg,
	}
	c.HTML(
		http.StatusOK,
		viewPathGruplar+"create.html",
		viewData,
	)

	//var satisFiyati1, agirlik float64
	// fiyatStr1 := access.OptionsApp.GetOption("satis_birim_fiyati_1")
	// satisFiyati1, _ = stnccollection.StringToFloat64(fiyatStr1)
	// fiyat1 := satisFiyati1 * agirlik
	// kisiBasiDusenFiyat1 := (satisFiyati1 * agirlik) / hisseAdeti
}

//Store save method
func (access *Gruplar) Store(c *gin.Context) {
	//color.Blue("Prints %s in blue.", "text")
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	var Gruplarbayrami, _, _ = gruplarBayramilerModel("createHayvanAtanmis", c)
	var savePostError = make(map[string]string)
	savePostError = Gruplarbayrami.Validate()

	hayvanlar, _ := access.HayvanBilgisiApp.GetAllFindDurum(entity.HayvanBilgisiDurumHayvanBosta)
	hisseAdeti := access.OptionsApp.GetOption("hisse_adeti")
	kurbanYili := access.OptionsApp.GetOption("kurban_yili")
	otomatikSira2021 := access.OptionsApp.GetOption("otomatik_sira_buyukbas_2021")

	if len(savePostError) == 0 {

		saveData, saveErr := access.GruplarApp.Save(&Gruplarbayrami)

		//TODO: burası bir func olabilir
		otomatikSira2021int := stnccollection.StringToint(otomatikSira2021) + 1
		otomatikSira2021str := stnccollection.IntToString(otomatikSira2021int)
		access.OptionsApp.SetOption("otomatik_sira_buyukbas_2021", otomatikSira2021str)

		hayvanBilgisiID := stnccollection.StringtoUint64(c.PostForm("HayvanBilgisiID"))
		access.HayvanBilgisiApp.UpdateSingleStatus(hayvanBilgisiID, 2)

		if saveErr != nil {
			savePostError = saveErr
		}

		lastID := strconv.FormatUint(uint64(saveData.ID), 10)
		stncsession.SetFlashMessage("Kayıt başarı ile eklendi", "success", c)

		if c.PostForm("kaydet") == "kaydet" {
			c.Redirect(http.StatusMovedPermanently, "/admin/gruplar/edit/"+lastID)
			return
		} else {
			c.Redirect(http.StatusMovedPermanently, "/admin/gruplar/create")
			return
		}
	}

	viewData := pongo2.Context{
		"title":            "Gruplar Ekleme",
		"csrf":             csrf.GetToken(c),
		"err":              savePostError,
		"post":             Gruplarbayrami,
		"hayvanlar":        hayvanlar,
		"hisseAdeti":       hisseAdeti,
		"kurbanYili":       kurbanYili,
		"flashMsg":         flashMsg,
		"otomatikSira2021": otomatikSira2021,
	}

	c.HTML(
		http.StatusOK,
		viewPathGruplar+"create.html",
		viewData,
	)

}

//Create all list
func (access *Gruplar) CreateEmpty(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	hayvanlar, _ := access.HayvanBilgisiApp.GetAllFindDurum(entity.HayvanBilgisiDurumHayvanBosta)
	hisseAdeti := access.OptionsApp.GetOption("hisse_adeti")
	kurbanYili := access.OptionsApp.GetOption("kurban_yili")
	otomatikSira2021 := access.OptionsApp.GetOption("otomatik_sira_buyukbas_2021")

	viewData := pongo2.Context{
		"title":            "Gruplar  Ekleme",
		"csrf":             csrf.GetToken(c),
		"hayvanlar":        hayvanlar,
		"hisseAdeti":       hisseAdeti,
		"kurbanYili":       kurbanYili,
		"otomatikSira2021": otomatikSira2021,
		"flashMsg":         flashMsg,
	}
	c.HTML(
		http.StatusOK,
		viewPathGruplar+"createEmpty.html",
		viewData,
	)
}

//StoreEmpty save method
func (access *Gruplar) StoreEmpty(c *gin.Context) {
	//color.Blue("Prints %s in blue.", "text")
	stncsession.IsLoggedInRedirect(c)

	var Gruplarbayrami, _, _ = gruplarBayramilerModel("empty", c)
	var savePostError = make(map[string]string)
	savePostError = Gruplarbayrami.Validate()
	fmt.Println(savePostError)
	hisseAdeti := access.OptionsApp.GetOption("hisse_adeti")
	kurbanYili := access.OptionsApp.GetOption("kurban_yili")
	otomatikSira2021 := access.OptionsApp.GetOption("otomatik_sira_buyukbas_2021")

	if len(savePostError) == 0 {
		saveData, saveErr := access.GruplarApp.Save(&Gruplarbayrami)
		//TODO: burası bir func olabilir
		otomatikSira2021int := stnccollection.StringToint(otomatikSira2021) + 1
		otomatikSira2021str := stnccollection.IntToString(otomatikSira2021int)
		access.OptionsApp.SetOption("otomatik_sira_buyukbas_2021", otomatikSira2021str)
		if saveErr != nil {
			savePostError = saveErr
		}
		lastID := stnccollection.Uint64toString(saveData.ID)
		stncsession.SetFlashMessage("Kayıt başarı ile eklendi", "success", c)

		if c.PostForm("kaydet") == "kaydet" {
			c.Redirect(http.StatusMovedPermanently, "/admin/gruplar/editEmpty/"+lastID)
		} else {
			c.Redirect(http.StatusMovedPermanently, "/admin/gruplar/createEmpty")
		}

		return
	} else {
		stncsession.SetFlashMessage("Zorunlu alanları lütfen doldurunuz", "danger", c)
	}
	flashMsg := stncsession.GetFlashMessage(c)
	viewData := pongo2.Context{
		"title":            "Gruplar Ekleme",
		"csrf":             csrf.GetToken(c),
		"err":              savePostError,
		"post":             Gruplarbayrami,
		"hisseAdeti":       hisseAdeti,
		"kurbanYili":       kurbanYili,
		"flashMsg":         flashMsg,
		"otomatikSira2021": otomatikSira2021,
	}

	c.HTML(
		http.StatusOK,
		viewPathGruplar+"createEmpty.html",
		viewData,
	)

}

//Edit genel Gruplar düzenleme işler TODO: kullanılmıyor v2 de yapılacak
func (access *Gruplar) Edit(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	if kID, err := strconv.ParseUint(c.Param("kID"), 10, 64); err == nil {
		if posts, err := access.GruplarApp.GetByID(kID); err == nil {

			hayvanlar, _ := access.HayvanBilgisiApp.GetAllFindDurum(entity.HayvanBilgisiDurumHayvanBosta)
			hisseAdeti := access.OptionsApp.GetOption("hisse_adeti")
			kurbanYili := access.OptionsApp.GetOption("kurban_yili")
			//gruplardaki index in verileri gelecek
			// kisiBilgileri, _ := access.kisiApp.GetByIDReferans(posts.Kurban[0].KisiID)
			viewData := pongo2.Context{
				"title":      "Gruplar  Düzenleme",
				"csrf":       csrf.GetToken(c),
				"post":       posts,
				"hayvanlar":  hayvanlar,
				"hisseAdeti": hisseAdeti,
				"kurbanYili": kurbanYili,
				"flashMsg":   flashMsg,
			}

			c.HTML(
				http.StatusOK,
				viewPathGruplar+"edit.html",
				viewData,
			)

		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

//Edit genel Gruplar düzenleme işler
func (access *Gruplar) EditEmpty(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	if kID, err := strconv.ParseUint(c.Param("kID"), 10, 64); err == nil {
		if posts, err := access.GruplarApp.GetByID(kID); err == nil {

			hayvanlar, _ := access.HayvanBilgisiApp.GetAllFindDurum(entity.HayvanBilgisiDurumHayvanBosta)
			hisseAdeti := access.OptionsApp.GetOption("hisse_adeti")
			kurbanYili := access.OptionsApp.GetOption("kurban_yili")

			viewData := pongo2.Context{
				"title":      "Gruplar  Düzenleme",
				"csrf":       csrf.GetToken(c),
				"post":       posts,
				"hayvanlar":  hayvanlar,
				"hisseAdeti": hisseAdeti,
				"kurbanYili": kurbanYili,
				"flashMsg":   flashMsg,
			}

			c.HTML(
				http.StatusOK,
				viewPathGruplar+"editEmpty.html",
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
func (access *Gruplar) Update(c *gin.Context) {

	stncsession.IsLoggedInRedirect(c)
	var Gruplar, _, id = gruplarBayramilerModel("edit", c)
	// stncsession.SetFlashMessage("Kayıt Düzenleme Kapatılmıştır, düzenleme geçersizdir", "warning", c)
	// c.Redirect(http.StatusMovedPermanently, "/admin/gruplar/edit/"+id)
	// return //TODO: v2 de açık olacak

	var savePostError = make(map[string]string)

	savePostError = Gruplar.Validate()

	if len(savePostError) == 0 {

		_, saveErr := access.GruplarApp.Update(&Gruplar)
		//TODO: burası bir func olabilir
		otomatikSira2021 := access.OptionsApp.GetOption("otomatik_sira_buyukbas_2021")
		otomatikSira2021int := stnccollection.StringToint(otomatikSira2021) + 1
		otomatikSira2021str := stnccollection.IntToString(otomatikSira2021int)
		access.OptionsApp.SetOption("otomatik_sira_buyukbas_2021", otomatikSira2021str)

		if saveErr != nil {
			savePostError = saveErr
		}
		stncsession.SetFlashMessage("Kayıt başarı ile düzenlendi", "success", c)

		//TODO: bu sender olayı daha mantıklı bişey olabilir
		if c.PostForm("sender") == "empty" {
			c.Redirect(http.StatusMovedPermanently, "/admin/gruplar/editEmpty/"+id)
		} else {
			c.Redirect(http.StatusMovedPermanently, "/admin/gruplar/edit/"+id)
		}

		return
	}
	viewData := pongo2.Context{
		"title": "Ödeme Düzenleme",
		"err":   savePostError,
		"csrf":  csrf.GetToken(c),
		"post":  Gruplar,
	}

	c.HTML(
		http.StatusOK,
		viewPathGruplar+"edit.html",
		viewData,
	)
}

//Degistir degistirme
func (access *Gruplar) Degistir(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	gruplar, _ := access.GruplarApp.GetAll()

	viewData := pongo2.Context{
		"title":    "Gruplar Değiştirme",
		"csrf":     csrf.GetToken(c),
		"gruplar":  gruplar,
		"flashMsg": flashMsg,
	}
	c.HTML(
		http.StatusOK,
		viewPathGruplar+"degistir.html",
		viewData,
	)
}

//HayvanAtamasiYap
func (access *Gruplar) HayvanAtamasiYap(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	hayvanlar, _ := access.HayvanBilgisiApp.GetAllFindDurum(entity.HayvanBilgisiDurumHayvanBosta)
	//	gruplar, _ := access.GruplarApp.GetAllFindDurum(entity.GruplarDurumGrupOlusmusHayvanYok)

	//#selman #dd
	// fmt.Printf("%+v\n", hayvanlar)
	viewData := pongo2.Context{
		"title":     "Hayvan ataması",
		"csrf":      csrf.GetToken(c),
		"hayvanlar": hayvanlar,
		// "gruplar":   gruplar,
		"flashMsg": flashMsg,
	}
	c.HTML(
		http.StatusOK,
		viewPathGruplar+"hayvanAta.html",
		viewData,
	)
}

func (access *Gruplar) HayvanAtamasiStore(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)

	var dataGruplar entity.Gruplar
	// agirlikID := stnccollection.StringToint(c.Param("agirlikID"))
	hayvanID := stnccollection.StringtoUint64(c.PostForm("HayvanBilgisiID"))

	grupBilgisiID := stnccollection.StringtoUint64(c.PostForm("GrupBilgisiIDSource"))
	fmt.Println("grupBilgisiIDSource")
	fmt.Println(grupBilgisiID)

	//hayvan bilgisi
	hayvanBilgisiData, _ := access.HayvanBilgisiApp.GetByID(hayvanID)
	//grup bilgisi
	grupBilgisiM, _ := access.GruplarApp.GetByID(grupBilgisiID)
	// empJSON, err := json.MarshalIndent(grupBilgisiM, "", "  ")
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	// fmt.Printf("MarshalIndent funnction output\n %s\n", string(empJSON))
	// fmt.Println(dataGruplar.HissedarAdet)
	// return
	satisFiyatTuru := stnccollection.StringToint(c.PostForm("SatisFiyatTuru"))

	fmt.Println("satisFiyatTuru")
	fmt.Println(satisFiyatTuru)
	fmt.Println("grupBilgisiM.Agirlik")
	fmt.Println(grupBilgisiM.Agirlik)

	kisiBasiDusenFiyat, kurbanToplamFiyati, kurbanBirimSatisFiyati := satisFiyatiHesaplayici(satisFiyatTuru, float64(hayvanBilgisiData.Agirlik), float64(grupBilgisiM.HissedarAdet))

	fmt.Println("kisiBasiDusenFiyat")
	fmt.Println(kisiBasiDusenFiyat)

	fmt.Println("kurbanToplamFiyati")
	fmt.Println(kurbanToplamFiyati)

	fmt.Println("kurbanBirimSatisFiyati")
	fmt.Println(kurbanBirimSatisFiyati)

	access.HayvanBilgisiApp.UpdateSingleStatus(hayvanID, 2) //TODO: aç bunu unutma sakın

	dataGruplar.ID = grupBilgisiID
	dataGruplar.UserID = stncsession.GetUserID2(c)
	dataGruplar.HayvanBilgisiID = hayvanBilgisiData.ID
	dataGruplar.GrupAdi = grupBilgisiM.GrupAdi
	dataGruplar.KesimSiraNo = grupBilgisiM.KesimSiraNo
	dataGruplar.HissedarAdet = grupBilgisiM.HissedarAdet
	dataGruplar.Aciklama = "Gruba Kurbanlık hayvan eklendi"
	dataGruplar.SatisFiyatTuru = satisFiyatTuru
	dataGruplar.SatisFiyati = kurbanBirimSatisFiyati
	dataGruplar.Siralama = grupBilgisiM.Siralama
	dataGruplar.AgirlikTipi = grupBilgisiM.AgirlikTipi
	dataGruplar.Durum = entity.GruplarDurumGrupOlusmusKesimlikHayvaniVar
	// dataGruplar.GrupLideri = grupBilgisiM.GrupLideri
	dataGruplar.Slug = grupBilgisiM.Slug
	dataGruplar.ToplamKurbanFiyati = kurbanToplamFiyati
	dataGruplar.KurbanBayramiYili = grupBilgisiM.KurbanBayramiYili
	hayvanAgirligi := hayvanBilgisiData.Agirlik
	dataGruplar.Agirlik = hayvanAgirligi
	access.GruplarApp.Update(&dataGruplar)

	///****** KURBAN FİYATLARI HRESAPLANIYOR ******
	kurbanlarData, _ := access.KurbanApp.GetByGrupID(grupBilgisiID)

	var hisseAdeti int
	access.KurbanApp.GetByGrupIDCount(grupBilgisiID, &hisseAdeti)

	kisiBasiDusenAgirlik := hayvanAgirligi / hisseAdeti
	// fmt.Printf("%+v\n", kurbanlarData)

	fmt.Printf("%+v\n", kurbanlarData[0].ID)

	var odemelerToplami float64

	for no, _ := range kurbanlarData {
		kurbanID := kurbanlarData[no].ID
		fmt.Println("------")
		fmt.Println(kurbanID)
		fmt.Println("numara")

		fmt.Println("------")

		odemelerToplami = access.OdemelerApp.OdemelerToplami(kurbanID)

		fmt.Println("odemelerToplami")
		fmt.Println(odemelerToplami)

		fmt.Println("kisiBasiDusenFiyat")
		fmt.Println(kisiBasiDusenFiyat)

		fmt.Println("ara toplam ")
		fmt.Println(odemelerToplami - kisiBasiDusenFiyat)

		if odemelerToplami > kisiBasiDusenFiyat {
			fmt.Println("kasa borçlu olacak")
			odemeEkle := entity.Odemeler{}
			odemeEkle.KurbanID = kurbanID
			odemeEkle.Alacak = 0
			odemeEkle.VerilenUcret = 0
			odemeEkle.KurbanFiyati = kisiBasiDusenFiyat
			kasaBorcu := odemelerToplami - kisiBasiDusenFiyat
			odemeEkle.Borc = kasaBorcu
			odemeEkle.Aciklama = "Kurbanlık hayvan eklendi / kasa borçlandı"
			odemeEkle.Durum = entity.OdemelerDurumGrupOlusmusKesimlikHayvaniVar
			odemeEkle.BorcDurum = entity.OdemelerBorcDurumKasaBorcluDurumda
			odemeEkle.UserID = stncsession.GetUserID2(c)
			odemeEkle.VerildigiTarih = time.Now()
			access.KurbanApp.SetKurbanKasaBorcuUpdate(kurbanID, kasaBorcu)
			access.KurbanApp.SetKurbanBakiyeUpdate(kurbanID, 0)
			// access.KurbanApp.SetKurbanBakiyeUpdate(kurbanID, odemelerToplami)
			access.KurbanApp.SetKurbanKalanUcretiUpdate(kurbanID, 0)
			access.KurbanApp.SetKurbanFiyatiUpdate(kurbanID, kisiBasiDusenFiyat)
			access.KurbanApp.SetKurbanBorcDurumUpdate(kurbanID, entity.KurbanBorcDurumKasaBorcluDurumda)
			access.KurbanApp.SetKurbanDurumUpdate(kurbanID, entity.KurbanBorcDurumHesapKapandi)
			access.KurbanApp.SetKurbanFiyatiUpdate(kurbanID, kisiBasiDusenFiyat)
			access.KurbanApp.SetKurbanAgirligi(kurbanID, kisiBasiDusenAgirlik)
			access.OdemelerApp.Save(&odemeEkle)
		}

		if odemelerToplami < kisiBasiDusenFiyat {
			fmt.Println("normal borçlu durum")
			odemeEkle := entity.Odemeler{}
			odemeEkle.KurbanID = kurbanID
			odemeEkle.KurbanID = kurbanID
			kalanUcret := kisiBasiDusenFiyat - odemelerToplami
			odemeEkle.Alacak = kalanUcret
			// verilenUcret := kisiBasiDusenFiyat - kalanUcret
			// odemeEkle.VerilenUcret = verilenUcret
			odemeEkle.VerilenUcret = 0
			odemeEkle.Bakiye = odemelerToplami
			odemeEkle.Borc = 0
			odemeEkle.KurbanFiyati = kisiBasiDusenFiyat
			odemeEkle.Aciklama = "Kurbanlık hayvan eklendi NOT: Ara Fark Otomatik Hesaplandı"
			odemeEkle.Durum = entity.OdemelerDurumGrupOlusmusKesimlikHayvaniVar
			odemeEkle.BorcDurum = entity.OdemelerBorcDurumTaksitOdemesi
			odemeEkle.UserID = stncsession.GetUserID2(c)
			odemeEkle.VerildigiTarih = time.Now()
			access.KurbanApp.SetKurbanKasaBorcuUpdate(kurbanID, 0)
			access.KurbanApp.SetKurbanBakiyeUpdate(kurbanID, odemelerToplami)
			access.KurbanApp.SetKurbanKalanUcretiUpdate(kurbanID, kalanUcret)
			access.KurbanApp.SetKurbanFiyatiUpdate(kurbanID, kisiBasiDusenFiyat)
			access.KurbanApp.SetKurbanBorcDurumUpdate(kurbanID, entity.KurbanBorcDurumTaksitOdemesi)
			access.KurbanApp.SetKurbanDurumUpdate(kurbanID, entity.KurbanDurumGrupOlusmusKesimlikHayvaniVar)
			access.KurbanApp.SetKurbanAgirligi(kurbanID, kisiBasiDusenAgirlik)
			access.OdemelerApp.Save(&odemeEkle)
		}

		if odemelerToplami == kisiBasiDusenFiyat {
			fmt.Println("hesap kapandı")
			odemeEkle := entity.Odemeler{}
			odemeEkle.Alacak = 0
			odemeEkle.Borc = 0
			odemeEkle.Aciklama = "Taksit Eklendi: Hesap Kapandı"
			odemeEkle.KurbanFiyati = kisiBasiDusenFiyat
			odemeEkle.Durum = entity.OdemelerDurumGrupOlusmusKesimlikHayvaniVar
			odemeEkle.BorcDurum = entity.OdemelerBorcDurumHesapKapandi
			odemeEkle.UserID = stncsession.GetUserID2(c)
			odemeEkle.VerildigiTarih = time.Now()
			access.KurbanApp.SetKurbanBakiyeUpdate(kurbanID, odemelerToplami)
			access.KurbanApp.SetKurbanKasaBorcuUpdate(kurbanID, 0)
			access.KurbanApp.SetKurbanKalanUcretiUpdate(kurbanID, 0)
			access.KurbanApp.SetKurbanFiyatiUpdate(kurbanID, kisiBasiDusenFiyat)
			access.KurbanApp.SetKurbanBorcDurumUpdate(kurbanID, entity.KurbanBorcDurumHesapKapandi)
			access.KurbanApp.SetKurbanDurumUpdate(kurbanID, entity.KurbanBorcDurumHesapKapandi)
			access.KurbanApp.SetKurbanAgirligi(kurbanID, kisiBasiDusenAgirlik)
			access.OdemelerApp.Save(&odemeEkle)
		}

	}
	// access.OdemelerApp.Save1(dataOdemeler)
	stncsession.SetFlashMessage("Kayıt başarı ile yapıldı", "success", c)
	c.Redirect(http.StatusMovedPermanently, "/admin/gruplar/hayvanata")
	return

}

func satisFiyatiHesaplayici(satisFiyatTuru int, agirlikfloat float64, hisseAdeti float64) (kisiBasiDusenFiyat float64, kurbanToplamFiyati float64, kurbanBirimSatisFiyati float64) {

	var satisFiyati1, satisFiyati2, satisFiyati3 float64

	db := repository.DB
	services, err1 := repository.RepositoriesInit(db)
	if err1 != nil {
		panic(err1)
	}

	options := InitOptions(services.Options)
	//TODO: burası da dışardan verilebilir aslında en son buna bi bakalım
	fiyatStr1 := options.OptionsApp.GetOption("satis_birim_fiyati_1")
	fiyatStr2 := options.OptionsApp.GetOption("satis_birim_fiyati_2")
	fiyatStr3 := options.OptionsApp.GetOption("satis_birim_fiyati_3")

	satisFiyati1, _ = stnccollection.StringToFloat64(fiyatStr1)
	satisFiyati2, _ = stnccollection.StringToFloat64(fiyatStr2)
	satisFiyati3, _ = stnccollection.StringToFloat64(fiyatStr3)

	toplamFiyataGorefiyat1 := satisFiyati1 * agirlikfloat
	toplamFiyataGorefiyat2 := satisFiyati2 * agirlikfloat
	toplamFiyataGorefiyat3 := satisFiyati3 * agirlikfloat

	fmt.Println("toplamFiyataGorefiyat1")
	fmt.Println(toplamFiyataGorefiyat1)

	fmt.Println("toplamFiyataGorefiyat2")
	fmt.Println(toplamFiyataGorefiyat2)

	fmt.Println("toplamFiyataGorefiyat3")
	fmt.Println(toplamFiyataGorefiyat3)

	// data.SatisFiyati3 = satisFiyati3
	kisiBasiDusenFiyatAraHesaplamaKusuratli1 := (satisFiyati1 * agirlikfloat) / hisseAdeti
	kisiBasiDusenFiyatAraHesaplamaKusuratli2 := (satisFiyati2 * agirlikfloat) / hisseAdeti
	kisiBasiDusenFiyatAraHesaplamaKusuratli3 := (satisFiyati3 * agirlikfloat) / hisseAdeti
	fmt.Println("kisiBasiDusenFiyatAraHesaplamaKusuratli1")
	fmt.Println(kisiBasiDusenFiyatAraHesaplamaKusuratli1)

	kisiBasiDusenFiyat1 := stnccollection.SayiYuvarla(math.Ceil(stnccollection.ToFixedDecimal(kisiBasiDusenFiyatAraHesaplamaKusuratli1, 2)))
	kisiBasiDusenFiyat2 := stnccollection.SayiYuvarla(math.Ceil(stnccollection.ToFixedDecimal(kisiBasiDusenFiyatAraHesaplamaKusuratli2, 2)))
	kisiBasiDusenFiyat3 := stnccollection.SayiYuvarla(math.Ceil(stnccollection.ToFixedDecimal(kisiBasiDusenFiyatAraHesaplamaKusuratli3, 2)))

	if satisFiyatTuru == 1 {
		fmt.Println("1 girer")
		kisiBasiDusenFiyat = kisiBasiDusenFiyat1
		kurbanToplamFiyati = toplamFiyataGorefiyat1
		kurbanBirimSatisFiyati = satisFiyati1
	} else if satisFiyatTuru == 2 {
		fmt.Println("2 girer")
		kisiBasiDusenFiyat = kisiBasiDusenFiyat2
		kurbanToplamFiyati = toplamFiyataGorefiyat2
		kurbanBirimSatisFiyati = satisFiyati2
	} else if satisFiyatTuru == 3 {
		fmt.Println("3 girer")
		kisiBasiDusenFiyat = kisiBasiDusenFiyat3
		kurbanToplamFiyati = toplamFiyataGorefiyat3
		kurbanBirimSatisFiyati = satisFiyati3
	}
	return kisiBasiDusenFiyat, kurbanToplamFiyati, kurbanBirimSatisFiyati
}

//HayvanListeAjax genel kurban düzenleme işler
func (access *Gruplar) GruplarListeAjaxAgirlikTuru(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)

	agirlikID := stnccollection.StringToint(c.Param("agirlikID"))
	jsonData, _ := access.GruplarApp.GetAllFindDurumAndAgirlikTipi(entity.GruplarDurumGrupOlusmusHayvanYok, agirlikID)

	viewData := pongo2.Context{
		"jsonData": jsonData,
		"csrf":     csrf.GetToken(c),
	}
	c.JSON(http.StatusOK, viewData)
	return
}

//TODO: önemli boş ve boş grup değieşcekse kesim sıraları değişsin

//ıkı grup olacak ikisinin de karşılık lı olması mı gerekir
//bir grup daki kişi sayısı diğerinden fazla olursa ???
//gruplar map de mi olacak
//dongu nasıl işleyecek
//en son map de mi veriler aktarılacak
//------dongusel hata ********
//okunacak kurban degerini okumadan once ıkısınde de aynı sayıda mı kurban var

func (access *Gruplar) DegistirStore(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	var degeriOkunacakGrupID uint64 = stnccollection.StringtoUint64(c.PostForm("HayvanBilgisiIDSource"))
	var hedefGrupID uint64 = stnccollection.StringtoUint64(c.PostForm("HayvanBilgisiIDTarget"))

	degeriOkunacakGrupStr := c.PostForm("HayvanBilgisiIDSource")
	hedefGrupStr := c.PostForm("HayvanBilgisiIDTarget")

	var total int
	access.getByGrupIDCountInline(degeriOkunacakGrupID, &total)

	//nasıl olamnlı

	//TODO: kesim sıra no durumu ne olacak gruplardaki bilgi

	// once hedef grup da hayvanın kurban fiyatını elinde tut, sonra kilosunu, sonra id değerini

	//kurbanların borcu var mı
	//burası atanacak grup
	var yeniKurbanFiyati float64
	grupA, _ := access.GruplarApp.GetByID(degeriOkunacakGrupID)
	fmt.Println("---GRUP DURUMLARI-----")
	fmt.Printf("%+v\n", grupA)
	gaToplamKurbanFiyati := grupA.ToplamKurbanFiyati
	yeniKurbanFiyati = gaToplamKurbanFiyati / float64(total)
	//TODO: sayıyı yuvarla olayı depency ve clean code olmalı unit test eklenecek buna
	yeniKurbanFiyati = stnccollection.SayiYuvarla(math.Ceil(stnccollection.ToFixedDecimal(yeniKurbanFiyati, 2)))
	gaAgirlik := grupA.Agirlik / total
	//gaID := grupA.ID

	// burası karşı tarafın kurbanın grubunun değişeceği kurban alanları

	// oncelikle karasa borcu varmı onun için borç durumuna bak
	// varsa
	// kasa borcu <=  yeni fiyattan == yeni fiyat + kasa borcu bu da kalan borcu verecek
	// yeni fiyat < kasa borcundan == yeni fiyat - kasa borcu burad kalan borcu verecek

	fmt.Println("total")
	fmt.Println(total)

	var updateKurbanlar = dto.KurbanUpdateRead{}
	var updateOdemeler = dto.Odemeler{}

	updateOdemeler.Makbuz = ""
	updateOdemeler.Aciklama = "Gruplar arasında yer değişimi oldu; " + degeriOkunacakGrupStr + " ile " + hedefGrupStr + " yer değiştirdi"

	updateKurbanlar.Agirlik = gaAgirlik
	kurbanlarHedef1, _ := access.KurbanApp.GetByGrupID(hedefGrupID)
	kurbanlarHedef2, _ := access.KurbanApp.GetByGrupID(degeriOkunacakGrupID)

	fmt.Printf("%+v\n", kurbanlarHedef1)
	fmt.Println("------------------")
	fmt.Printf("%+v\n", kurbanlarHedef2)
	fmt.Println("------------------")

	var okunanDeger int = 1

	//grupB.Durum == entity.GruplarDurumGrupOlusmusKesimlikHayvaniVar  boyle durumda ne olur TODO: ???
	//burada grupun boş olmasına (hayvan atanmamış) iki grupunda boş olmasına göre işlem yapıyoruz, göre bir değişim yapıyoruz
	if grupA.Durum == entity.GruplarDurumGrupOlusmusKesimlikHayvaniVar {
		for no, _ := range kurbanlarHedef1 {
			// fmt.Println(veri)

			odenenToplamUcret := kurbanlarHedef1[no].Bakiye
			if kurbanlarHedef1[no].BorcDurum == entity.KurbanBorcDurumKasaBorcluDurumda {

			}

			fmt.Println("kurbanlarHedef[no].KurbanFiyati")
			fmt.Println(kurbanlarHedef1[no].KurbanFiyati)
			fmt.Println("yeniKurbanFiyati")
			fmt.Println(yeniKurbanFiyati)
			fmt.Println("odenenToplamUcret")
			fmt.Println(odenenToplamUcret)
			fmt.Println("ıd")
			fmt.Println(kurbanlarHedef1[no].ID)

			kurbanID := kurbanlarHedef1[no].ID

			// odenenToplamUcret := kurbanlarHedef[no].Bakiye
			// yeniBorc := odenenToplamUcret - gfg
			updateKurbanlar.KurbanFiyati = 0
			updateKurbanlar.KasaBorcu = 0
			updateKurbanlar.KalanUcret = 0
			updateKurbanlar.BorcDurum = kurbanlarHedef1[no].BorcDurum
			updateKurbanlar.ID = kurbanID
			updateOdemeler.BorcDurum = entity.OdemelerBorcDurumIkiGrupYerDegistirdi
			updateOdemeler.KurbanID = kurbanID
			updateOdemeler.UserID = stncsession.GetUserID2(c)
			updateOdemeler.Alacak = 0
			updateOdemeler.Borc = 0
			access.KurbanApp.UpdatePriceList(&updateKurbanlar)
			access.OdemelerApp.SetUpdateHesaplar(&updateOdemeler)
			okunanDeger++
			fmt.Println("---------deger---------")
			fmt.Println(okunanDeger)
		}
	}

	fmt.Printf("%+v\n", kurbanlarHedef1)

	viewData := pongo2.Context{
		"title": "Gruplar Karşilaştır",
		"csrf":  csrf.GetToken(c),

		"hayvanlar": kurbanlarHedef1,
		"grup":      grupA,
		"flashMsg":  flashMsg,
	}

	c.HTML(
		http.StatusOK,
		viewPathGruplar+"degistir.html",
		viewData,
	)
}

//KarsilastirAPI
func (access *Gruplar) GrupBilgisiAPI(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	if kID, err := strconv.ParseUint(c.Param("kID"), 10, 64); err == nil {
		if jsonData, err := access.GruplarApp.GetByIDAllRelations(kID); err == nil {

			viewData := pongo2.Context{
				"title":    "Grup Düzenleme",
				"csrf":     csrf.GetToken(c),
				"jsonData": jsonData,
				"flashMsg": flashMsg,
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

//KarsilastirAPI
func (access *Gruplar) GrupBilgisiHayvanBosOlanlarAPI(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)

	if kID, err := strconv.ParseUint(c.Param("kID"), 10, 64); err == nil {
		if jsonData, err := access.GruplarApp.GetByIDAllRelationsHayvanOlmayanlar(kID); err == nil {

			viewData := pongo2.Context{
				"title":    "Grup Düzenleme",
				"csrf":     csrf.GetToken(c),
				"jsonData": jsonData,
				"flashMsg": flashMsg,
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

//odemelere kayıt gir , yuvarla, eşiitri durumlaında ne olur , bir tarafın yedeği alnıp atılacak ,notficaiton oluşması lazım
//Karsilastir karsılaştirma işi

//excel hücreleri içindir
func IntToLetters(number int32) (letters string) {
	number--
	if firstLetter := number / 26; firstLetter > 0 {
		letters += IntToLetters(firstLetter)
		letters += string('A' + number%26)
	} else {
		letters += string('A' + number)
	}

	return
}
func (access *Gruplar) getByGrupIDCountInline(id uint64, postTotalCount *int) {

	var total int
	access.KurbanApp.GetByGrupIDCount(id, &total)
	*postTotalCount = total
}

func tableHeader(cellNumber int) map[string]string {
	number := stnccollection.IntToString(cellNumber + 1)
	return map[string]string{"A" + number: "NO",
		"B" + number: "Adı Soyadı",
		"C" + number: "Telefon",
		"D" + number: "Toplam Ödeme",
		"E" + number: "Kalan Borç",
		"F" + number: "Referans",
		"G" + number: "Vekalet",
	}

}

//https://xuri.me/excelize/en/cell.html#SetCellStyle
//Logout güvenli çıkış
func (access *Gruplar) Excel(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	f := excelize.NewFile()
	//https://xuri.me/excelize/en/style.html#number_format
	//font  "underline":"single",
	//				"vertical": "",
	styleBOLD, err := f.NewStyle(`
		{


			"font":
			     {

				 "bold": true,
				 "italic": false,
				 "family": "Times New Roman",
				 "size": 11,
				 "color": "#000000"
				  },

		     "alignment":
					{
				    "horizontal": "center",
					"vertical": "center",

					"wrap_text": true
					}
		}

		`)
	fmt.Println(styleBOLD)

	var styleText1 int
	styleText1, _ = f.NewStyle(`
	{
		   "alignment":
				{
				"horizontal": "center",
				"vertical": "center",
				"wrap_text": true
		   }
	} `)

	styleDusuk, _ := f.NewStyle(`
		{

			"fill":
					{
					"type":"pattern",
					"color":["#e55353"],
			 		"pattern":1
				  },

			"font":
			     {

				 "bold": true,
				 "italic": false,
				 "family": "Times New Roman",
				 "size": 11,
				 "color": "#ffffff"
				  },

		     "alignment":
					{
				    "horizontal": "center",
				    "vertical": "center",
					"ident": 0,
					"justify_last_line": true,
					"reading_order": 1,
					"relative_indent": 0,
					"shrink_to_fit": true,

					"wrap_text": true
					},

		"border": [
						{
							"type": "left",
							"color": "d8dbe0",
							"style": 3
						},
						{
							"type": "top",
							"color": "d8dbe0",
							"style": 4
						},
						{
							"type": "bottom",
							"color": "d8dbe0",
							"style": 5
						},
						{
							"type": "right",
							"color": "d8dbe0",
							"style": 6
						}
						]

		}
		`)

	fmt.Println(styleDusuk)
	styleOrta, _ := f.NewStyle(`
		{

			"fill":
					{
					"type":"pattern",
					"color":["#f9b115"],
			 		"pattern":1
				  },

			"font":
			     {
				 "bold": true,
				 "italic": false,
				 "family": "Times New Roman",
				 "size": 11,
				 "color": "#ffffff"
				  },

		     "alignment":
					{
				    "horizontal": "center",
				    "vertical": "center",
					"ident": 0,
					"justify_last_line": true,
					"reading_order": 1,
					"relative_indent": 0,
					"shrink_to_fit": true,
					"wrap_text": true
					},

	        	"border": [
						{
							"type": "left",
							"color": "d8dbe0",
							"style": 3
						},
						{
							"type": "top",
							"color": "d8dbe0",
							"style": 4
						},
						{
							"type": "bottom",
							"color": "d8dbe0",
							"style": 5
						},
						{
							"type": "right",
							"color": "d8dbe0",
							"style": 6
						}
						]

		}
		`)
	fmt.Println(styleOrta)

	styleYuksek, err := f.NewStyle(`
		{

			"fill":
					{
					"type":"pattern",
					"color":["#2eb85c"],
			 		"pattern":1
				  },

			"font":
			     {

				 "bold": true,
				 "italic": false,
				 "family": "Times New Roman",
				 "size": 11,
				 "color": "#ffffff"
				  },

		     "alignment":
					{
				    "horizontal": "center",
				    "vertical": "center",
					"ident": 0,
					"justify_last_line": true,
					"reading_order": 1,
					"relative_indent": 0,
					"shrink_to_fit": true,

					"wrap_text": true
					},

		       "border": [
						{
							"type": "left",
							"color": "d8dbe0",
							"style": 3
						},
						{
							"type": "top",
							"color": "d8dbe0",
							"style": 4
						},
						{
							"type": "bottom",
							"color": "d8dbe0",
							"style": 5
						},
						{
							"type": "right",
							"color": "d8dbe0",
							"style": 6
						}
						]

		}
		`)

	fmt.Println(styleYuksek)
	if err != nil {
		fmt.Println(err)
	}
	var veriler []dto.GruplarExcelandIndex = GrupDataGenerate(access)
	ac := accounting.Accounting{Precision: 0, Thousand: ".", Decimal: ","}

	//burayı fonk yap sonra parametre ver b ve sıra no gibi
	var i int = 1

	for num, v := range veriler {
		i++
		f.MergeCell("Sheet1", "A"+stnccollection.IntToString(i), "G"+stnccollection.IntToString(i))

		if v.AgirlikTipi == 1 {
			f.SetCellStyle("Sheet1", "A"+stnccollection.IntToString(i), "G"+stnccollection.IntToString(i), styleDusuk)
		} else if v.AgirlikTipi == 2 {
			f.SetCellStyle("Sheet1", "A"+stnccollection.IntToString(i), "G"+stnccollection.IntToString(i), styleOrta)
		} else if v.AgirlikTipi == 3 {
			f.SetCellStyle("Sheet1", "A"+stnccollection.IntToString(i), "G"+stnccollection.IntToString(i), styleYuksek)

		}

		f.SetRowHeight("Sheet1", i, 50)

		kesimSiraNoText := "Kesim No: " + stnccollection.IntToString(v.KesimSiraNo) + " - "
		agirlikTipiText := "Ağırlık Tipi: " + veriler[num].GrupIsoTopeTRname + " - "
		grupAdiText := "Grup Adı: " + v.GrupAdi + " - "
		kurbanFiyatiText := "Kurban Fiyatı: " + ac.FormatMoney(v.ToplamKurbanFiyati) + " -  TL"
		kupeNoText := " Küpe No: " + v.KupeNo
		KisiBasiDusenHisseFiyatiText := "\n Kişi Başı Düşen Hisse Fiyati: " + ac.FormatMoney(v.KisiBasiDusenHisseFiyati) + " TL"
		toplamAgirlikText := " / Toplam Ağırlık: " + stnccollection.IntToString(v.Agirlik)
		kisibasiAgirlikText := " / Kişi Başı Düşen Tahmini Kilo: " + stnccollection.IntToString(v.Agirlik/7)

		f.SetCellValue("Sheet1", "A"+stnccollection.IntToString(i), kesimSiraNoText+agirlikTipiText+grupAdiText+kurbanFiyatiText+kupeNoText+KisiBasiDusenHisseFiyatiText+toplamAgirlikText+kisibasiAgirlikText)

		//if i%2 == 0 {
		//f.SetRowHeight("Sheet1", num, 20)
		// f.UnmergeCell("Sheet1", "A"+stnccollection.IntToString(num), "J"+stnccollection.IntToString(num))
		//		}

		categories := tableHeader(i)
		for k, v := range categories {
			f.SetCellValue("Sheet1", k, v)
			f.SetCellStyle("Sheet1", k, k, styleBOLD)

		}

		// fmt.Println(len(v.Kurban))
		a := 1
		for _, kurban := range v.KurbanKisiList {
			f.SetCellStyle("Sheet1", "A"+stnccollection.IntToString(i+2), "G"+stnccollection.IntToString(i+2), styleText1)

			f.SetCellValue("Sheet1", "A"+stnccollection.IntToString(i+2), a)

			f.SetCellValue("Sheet1", "B"+stnccollection.IntToString(i+2), kisiKontrol(kurban.KisiAdSoyad))

			f.SetCellValue("Sheet1", "C"+stnccollection.IntToString(i+2), kurban.KisiTelefon)

			f.SetCellValue("Sheet1", "D"+stnccollection.IntToString(i+2), toplamOdemeSonucu(kurban.BorcDurum, ac.FormatMoney(kurban.Borc), ac.FormatMoney(kurban.Alacak), ac.FormatMoney(kurban.Bakiye))) //Toplam Ödeme

			f.SetCellValue("Sheet1", "E"+stnccollection.IntToString(i+2), kalanBorcSonucu(kurban.BorcDurum, ac.FormatMoney(kurban.Borc), ac.FormatMoney(kurban.Alacak), ac.FormatMoney(kurban.Bakiye))) //Kalan Borç

			f.SetCellValue("Sheet1", "F"+stnccollection.IntToString(i+2), kurban.ReferansAdSoyad+" \n Tel:"+kurban.ReferansTelefon)

			f.SetCellValue("Sheet1", "G"+stnccollection.IntToString(i+2), VekaletSonucu(kurban.VekaletDurumu))

			// f.SetRowHeight("Sheet1", 1, 10) //row büyütür

			i++
			a++
		}
		i++
	}
	// f.SetCellValue("Sheet1", "A2", "Grup Adı: "+v.GrupAdi)

	f.SetColWidth("Sheet1", "A", "A", 7)
	f.SetColWidth("Sheet1", "B", "B", 12)
	f.SetColWidth("Sheet1", "C", "C", 12)
	f.SetColWidth("Sheet1", "D", "D", 12)
	f.SetColWidth("Sheet1", "E", "E", 12)
	f.SetColWidth("Sheet1", "F", "F", 16)
	f.SetColWidth("Sheet1", "G", "G", 13)

	// Create a new sheet.
	f.SetDocProps(&excelize.DocProperties{
		Category:      "krbn",
		ContentStatus: "Puplish",
		// Created:        "2019-06-04T22:00:10Z",
		Creator:        "TeamWork Excelize",
		Description:    "This file created by TeamWork",
		Identifier:     "xlsx",
		Keywords:       "Spreadsheet",
		LastModifiedBy: "Go TeamWork",
		// Modified:       "2019-06-04T22:00:10Z",
		Revision: "0",
		Subject:  "Kurban Listesi",
		Title:    "Kurban Listesi",
		Language: "tr-TR",
		Version:  "1.0.0",
	})
	// Save spreadsheet by the given path.
	if err := f.SaveAs("gruplarExcel.xlsx"); err != nil {
		fmt.Println(err)
	}
	c.File("gruplarExcel.xlsx")
}

func VekaletSonucu(vekaletStatus int) string {
	if vekaletStatus == 1 {
		return "Vekalet Alınmamış"
	} else {
		return "Vekalet Alındı"
	}
}

func toplamOdemeSonucu(borcDurum int, borc string, alacak string, bakiye string) string {
	var ret string = ""
	if borcDurum == 3 {
		ret = borc

		ret += " TL Kasamız Borçlu Durumda "

	} else if borcDurum != 3 {
		ret = bakiye
	}

	return ret
}

func kalanBorcSonucu(borcDurum int, borc string, alacak string, bakiye string) string {
	var ret string = alacak
	if borcDurum == 6 {
		ret += " Ödeme Tamamlandı "
	}

	return ret
}

func kisiKontrol(kisiAdSoyad string) string {
	if kisiAdSoyad == "Boş GRUP, Kişi kaydı için tıklayınız" {
		return "---"
	} else {
		return kisiAdSoyad
	}
}

//gruplarBayramilerModel model
func gruplarBayramilerModel(formType string, c *gin.Context) (grupData entity.Gruplar, idD uint64, idStr string) {
	id := c.PostForm("ID")
	idInt, _ := strconv.Atoi(id)
	var idN uint64
	var hissedarAdet int
	var slug string
	idN = uint64(idInt)

	t := time.Now()

	agirlik := stnccollection.StringToint(c.PostForm("Agirlik"))

	grupData.ID = idN
	grupData.Aciklama = c.PostForm("Aciklama")
	grupData.GrupAdi = c.PostForm("GrupAdi")
	grupData.ToplamKurbanFiyati, _ = stnccollection.StringToFloat64(c.PostForm("ToplamKurbanFiyati"))
	hissedarAdet, _ = strconv.Atoi(c.PostForm("HissedarAdet"))
	grupData.HissedarAdet = hissedarAdet
	grupData.KurbanBayramiYili, _ = strconv.Atoi(c.PostForm("KurbanBayramiYili"))
	grupData.KesimSiraNo, _ = strconv.Atoi(c.PostForm("KesimSiraNo"))
	grupData.Siralama, _ = strconv.Atoi(c.PostForm("Siralama"))
	fmt.Println("ağırlık ")
	fmt.Println(strconv.Atoi(c.PostForm("AgirlikTipi")))

	grupData.AgirlikTipi, _ = strconv.Atoi(c.PostForm("AgirlikTipi"))
	grupData.SatisFiyatTuru, _ = strconv.Atoi(c.PostForm("SatisFiyatTuru"))
	grupData.SatisFiyati, _ = stnccollection.StringToFloat64(c.PostForm("SatisFiyati"))
	grupData.HayvanBilgisiID = stnccollection.StringtoUint64(c.PostForm("HayvanBilgisiID"))
	grupData.Agirlik = agirlik

	//EMPTY
	if formType == "empty" {
		// db := repository.DB
		// services, err1 := repository.RepositoriesInit(db)
		// if err1 != nil {
		// 	panic(err1)
		// }
		// options := InitOptions(services.Options)
		// hisseAdetiOption := options.OptionsApp.GetOption("hisse_adeti")

		grupData.Durum = entity.GruplarDurumGrupOlusmusHayvanYok
		//TODO: user id eklenecek
		slug = stnchelper.RandSlugV1(2) + stnccollection.IntToString(t.Second()) + stnchelper.RandSlugV1(2)
		grupData.Slug = slug
		grupData.Kurban = []entity.Kurban{}
		for i := 0; i < hissedarAdet; i++ {
			slug = stnchelper.RandSlugV1(2) + stnccollection.IntToString(t.Second()) + stnchelper.RandSlugV1(2)
			a := i
			a++
			appendData := entity.Kurban{
				KisiID:        1,
				Aciklama:      " ",
				VekaletDurumu: entity.VekaletDurumuAlinmadi,
				KurbanTuru:    12,
				Agirlik:       0,
				KurbanFiyati:  0,
				HayvanCinsi:   5,
				Alacak:        0,
				Borc:          0,
				Slug:          slug,
				GrupLideri:    a,
				UserID:        stncsession.GetUserID2(c),
				Durum:         entity.KurbanDurumGrupOlusmusHayvanYok,
				BorcDurum:     entity.KurbanBorcDurumIlkEklenenFiyat,
				Odemeler: []dto.Odemeler{{
					Aciklama:       "Kurban Henüz Atanmamıştır",
					VerilenUcret:   0,
					Alacak:         0,
					Makbuz:         "----",
					UserID:         stncsession.GetUserID2(c),
					VerildigiTarih: time.Now(),
					Durum:          entity.OdemelerDurumGrupOlusmusHayvanYok,
					BorcDurum:      entity.OdemelerBorcDurumIlkEklenenFiyat,
				}},
			}

			grupData.Kurban = append(grupData.Kurban, appendData)
		}
	}
	//TODO: tr karakter createHayvanAtanmis
	if formType == "createHayvanAtanmis" {

		grupData.Durum = entity.GruplarDurumGrupOlusmusKesimlikHayvaniVar

		db := repository.DB
		services, err1 := repository.RepositoriesInit(db)
		if err1 != nil {
			panic(err1)
		}
		options := InitOptions(services.Options)

		hisseAdeti := options.OptionsApp.GetOption("hisse_adeti")

		toplamKurbanFiyati, _ := stnccollection.StringToFloat64(c.PostForm("ToplamKurbanFiyati"))

		toplamKurbanFiyatiFix := stnccollection.ToFixedDecimal(toplamKurbanFiyati, 1)

		var hisseFiyati float64

		hisseAdetiFloat, _ := stnccollection.StringToFloat64(hisseAdeti)

		hisseAdetiInt := stnccollection.StringToint(hisseAdeti)

		hisseFiyati = math.Round(toplamKurbanFiyatiFix / hisseAdetiFloat)

		hisseFiyati = stnccollection.SayiYuvarla(math.Ceil(stnccollection.ToFixedDecimal(hisseFiyati, 2)))

		kisiBasiAgirlik := agirlik / hisseAdetiInt

		grupData.Kurban = []entity.Kurban{}
		for i := 0; i < hisseAdetiInt; i++ {
			slug = stnchelper.RandSlugV1(2) + stnccollection.IntToString(t.Second()) + stnchelper.RandSlugV1(2)
			appendData := entity.Kurban{
				KisiID:        1,
				Aciklama:      " ",
				Slug:          slug,
				VekaletDurumu: entity.VekaletDurumuAlinmadi,
				KurbanTuru:    12,
				Durum:         entity.KurbanDurumGrupOlusmusKesimlikHayvaniVar,
				Agirlik:       kisiBasiAgirlik,
				KurbanFiyati:  hisseFiyati,
				HayvanCinsi:   5,
				Alacak:        hisseFiyati,
				UserID:        stncsession.GetUserID2(c),
				Odemeler: []dto.Odemeler{{
					Aciklama:       "İlk Eklenen Fiyat",
					VerilenUcret:   hisseFiyati,
					Alacak:         hisseFiyati,
					Bakiye:         0,
					UserID:         stncsession.GetUserID2(c),
					VerildigiTarih: time.Now(),
					BorcDurum:      entity.OdemelerBorcDurumIlkEklenenFiyat,
				}},
			}
			grupData.Kurban = append(grupData.Kurban, appendData)
		}
	}

	return grupData, idN, id
}

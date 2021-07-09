package controller

import (
	"net/http"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stncsession"
	"stncCms/app/domain/repository"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

const viewPathIndex = "admin/index/"

//Index all list f
func Index(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	viewData := pongo2.Context{
		"title": "Posts",
		"csrf":  csrf.GetToken(c),
	}
	c.HTML(
		http.StatusOK,
		viewPathIndex+"index.html",
		viewData,
	)
}

//OptionsDefault all list f
func OptionsDefault(c *gin.Context) {
	// stncsession.IsLoggedInRedirect(c)

	//buraya bir oprion otılacak bunlar giriş yaptıktan sonra veri varmı yok mu bakacak

	db := repository.DB

	option1 := entity.Options{OptionName: "siteUrl", OptionValue: "http://localhost:8888/"}
	db.Debug().Create(&option1)

	option2 := entity.Options{OptionName: "kurban_yili", OptionValue: "2021"}
	db.Debug().Create(&option2)

	option3 := entity.Options{OptionName: "hisse_adeti", OptionValue: "7"}
	db.Debug().Create(&option3)

	option4 := entity.Options{OptionName: "satis_birim_fiyati_1", OptionValue: "20"}
	db.Debug().Create(&option4)

	option5 := entity.Options{OptionName: "satis_birim_fiyati_2", OptionValue: "25"}
	db.Debug().Create(&option5)

	option6 := entity.Options{OptionName: "satis_birim_fiyati_3", OptionValue: "30"}
	db.Debug().Create(&option6)

	option7 := entity.Options{OptionName: "hayvan_dusuk_agirligi", OptionValue: "0-200"}
	db.Debug().Create(&option7)

	option78 := entity.Options{OptionName: "hayvan_orta_agirligi", OptionValue: "200-600"}
	db.Debug().Create(&option78)

	option786 := entity.Options{OptionName: "hayvan_yuksek_agirligi", OptionValue: "600-1500"}
	db.Debug().Create(&option786)

	option8 := entity.Options{OptionName: "alis_birim_fiyati_1", OptionValue: "10"}
	db.Debug().Create(&option8)

	option9 := entity.Options{OptionName: "alis_birim_fiyati_2", OptionValue: "15"}
	db.Debug().Create(&option9)

	option10 := entity.Options{OptionName: "alis_birim_fiyati_3", OptionValue: "20"}
	db.Debug().Create(&option10)

	option11 := entity.Options{OptionName: "otomatik_sira_buyukbas_2021", OptionValue: "1"}
	db.Debug().Create(&option11)

	option12 := entity.Options{OptionName: "otomatik_sira_kuyukbas_2021", OptionValue: "1"}
	db.Debug().Create(&option12)

	//mutluerF9E
	user := entity.User{FirstName: "Sel", LastName: "t", Email: "hk@gmail.com", Password: "4544bcb2ce39fe656c64f0860895bdaccff7b8c0"}
	db.Debug().Create(&user)

	hayvansatisyeri := entity.HayvanSatisYerleri{UserID: stncsession.GetUserID2(c), YerAdi: "Bu ilk kurulum için oluşturulmuştur, hayvanlar henüz belli değilse bunu seçiniz", Durum: 0, Slug: "bosVeri"}
	db.Debug().Create(&hayvansatisyeri)

	//TODO:KurbanBayramiYili optiondan gelecek

	//TODO: delete yapılcak
	adak := entity.Gruplar{KesimSiraNo: 0, Aciklama: "Adak Olarak Kesilecek Hayvanların Üst Grubudur, değişikliğe kapalıdır", HissedarAdet: 1, Durum: 1, KurbanBayramiYili: 2021, HayvanBilgisiID: 0, Slug: "adak"}
	db.Debug().Create(&adak)

	akika := entity.Gruplar{KesimSiraNo: 0, Aciklama: "Akika Olarak Kesilecek Hayvanların Üst Grubudur, değişikliğe kapalıdır", HissedarAdet: 1, Durum: 1, KurbanBayramiYili: 2021, HayvanBilgisiID: 0, Slug: "akika"}
	db.Debug().Create(&akika)

	sukur := entity.Gruplar{KesimSiraNo: 0, Aciklama: "Şükür Olarak Kesilecek Hayvanların Üst Grubudur, değişikliğe kapalıdır", HissedarAdet: 1, Durum: 1, KurbanBayramiYili: 2021, HayvanBilgisiID: 0, Slug: "sukur"}
	db.Debug().Create(&sukur)

	sahibNi := entity.Gruplar{KesimSiraNo: 0, Aciklama: "SAHİBİNİN NİYETİNE Olarak Kesilecek Hayvanların Üst Grubudur, değişikliğe kapalıdır", HissedarAdet: 1, Durum: 1, KurbanBayramiYili: 2021, HayvanBilgisiID: 0, Slug: "sahibi-niyetine"}
	db.Debug().Create(&sahibNi)

	bagis := entity.Gruplar{KesimSiraNo: 0, Aciklama: "BAĞIŞ Olarak Kesilecek Hayvanların Üst Grubudur, değişikliğe kapalıdır", HissedarAdet: 1, Durum: 1, KurbanBayramiYili: 2021, HayvanBilgisiID: 0, Slug: "hayir"}
	db.Debug().Create(&bagis)

	naf := entity.Gruplar{KesimSiraNo: 0, Aciklama: "NAFİLE Olarak Kesilecek Hayvanların Üst Grubudur, değişikliğe kapalıdır", HissedarAdet: 1, Durum: 1, KurbanBayramiYili: 2021, HayvanBilgisiID: 0, Slug: "nafile"}
	db.Debug().Create(&naf)

	sifa := entity.Gruplar{KesimSiraNo: 0, Aciklama: "Şifa Olarak Kesilecek Hayvanların Üst Grubudur, değişikliğe kapalıdır", HissedarAdet: 1, Durum: 1, KurbanBayramiYili: 2021, HayvanBilgisiID: 0, Slug: "sifa"}
	db.Debug().Create(&sifa)

	krubanBayrami := entity.Gruplar{KesimSiraNo: 0, Aciklama: "Kurban Bayramı Kesilecek Büyük Baş Hayvan Atanmamış Değer", HissedarAdet: 1, Durum: 1, KurbanBayramiYili: 2021, HayvanBilgisiID: 0, Slug: "kurbanBayramiBuyukbas"}
	db.Debug().Create(&krubanBayrami)

	krubanBayramiKucukbas := entity.Gruplar{KesimSiraNo: 0, Aciklama: "Kurban Bayramı Kesilecek Büyük Baş Hayvan Atanmamış Değer", HissedarAdet: 1, Durum: 1, KurbanBayramiYili: 2021, HayvanBilgisiID: 0, Slug: "kurbanBayramiKucukBas"}
	db.Debug().Create(&krubanBayramiKucukbas)

	ModulKurban := entity.Modules{ModulName: "kurban", Status: 1, UserID: 1}
	db.Debug().Create(&ModulKurban)

	Modulodemeler := entity.Modules{ModulName: "odemeler", Status: 1, UserID: 1}
	db.Debug().Create(&Modulodemeler)

	ModulHayvanBilgisi := entity.Modules{ModulName: "hayvanBilgisi", Status: 1, UserID: 1}
	db.Debug().Create(&ModulHayvanBilgisi)

	ModulDashborad := entity.Modules{ModulName: "dashboard", Status: 1, UserID: 1}
	db.Debug().Create(&ModulDashborad)

	ModulAyarlar := entity.Modules{ModulName: "ayarlar", Status: 1, UserID: 1}
	db.Debug().Create(&ModulAyarlar)

	ModulGruplar := entity.Modules{ModulName: "gruplar", Status: 1, UserID: 1}
	db.Debug().Create(&ModulGruplar)

	ModulhayvanSatisYerleri := entity.Modules{ModulName: "hayvanSatisYerleri", Status: 1, UserID: 1}
	db.Debug().Create(&ModulhayvanSatisYerleri)

	ModulKisiler := entity.Modules{ModulName: "kisiler", Status: 1, UserID: 1}
	db.Debug().Create(&ModulKisiler)

	Modulkullanici := entity.Modules{ModulName: "kullanici", Status: 1, UserID: 1}
	db.Debug().Create(&Modulkullanici)

	Kisiler1 := entity.Kisiler{AdSoyad: "Boş GRUP, Kişi kaydı için tıklayınız ", Aciklama: "Boş olarak oluşacak grupların buraya kaydedilmesi içindir"}
	db.Debug().Create(&Kisiler1)

	// db.Debug().Delete(&entity.Kisiler{}, 1)

	c.JSON(http.StatusOK, "yapıldı")
}

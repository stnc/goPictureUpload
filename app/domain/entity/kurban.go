package entity

import (
	"fmt"
	"html"
	"stncCms/app/domain/dto"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	tr_translations "gopkg.in/go-playground/validator.v9/translations/tr"
)

const (
	//Kurbanlar  DURUM BİLGİLERİ
	KurbanDurumKurbanEklendiKurbanBayraminaAitDegil = 1 // kurban eklenmiş ama kurban bayramına ait değil yani direk kurban girişi,
	KurbanDurumGrupOlusmusHayvanYok                 = 2 // grup oluşmuş ama kimse atanmamış , yani kesimlik inek verilmemiştir
	KurbanDurumGrupOlusmusKesimlikHayvaniVar        = 3 // grup atanmış yani bir kesimlik inek verilmiş
	KurbanDurumKurbanKesimiTamamlanmis              = 4 // kurban kesimi tamamlanmış
	KurbanDurumIkiGrupYerDegistirdi                 = 5 //  iki grup arasi degisim yapildi
	KurbanDurumHesapKapandi                         = 6 //  tüm borcunu ödemiş
	KurbanDurumHayvanKesildi                        = 7 //  hayvan kesildi
	KurbanDurumHayvanKesildiTeslimEdildi            = 8 //  hayvan kesildi teslim edildi
	KurbanDurumFiyatManuelDegistirildi              = 9 //  kurbanın düzenleme alanında fiyatı değiştirilmiş

	//Kurbanlar BoRC DURUM BİLGİLERİ

	KurbanBorcDurumIlkEklenenFiyat           = 1  // ilk eklenen fiyat değeri
	KurbanBorcDurumTaksitOdemesi             = 2  // taksit eklemiş
	KurbanBorcDurumKasaBorcluDurumda         = 3  //  kasa borçlu kalmışsa
	KurbanBorcDurumKaporaOdemesiHayvanBos    = 4  //  kapora odendi ama hayvan atanmamışdır
	KurbanBorcDurumIkiGrupYerDegistirdi      = 5  //  iki grup arasi degisim yapildi
	KurbanBorcDurumHesapKapandi              = 6  //  tüm borcunu ödemiş
	KurbanBorcDurumHayvanKesildi             = 7  //  hayvan kesildi
	KurbanBorcDurumHayvanKesildiTeslimEdildi = 8  //  hayvan kesildi teslim edildi
	KurbanBorcDurumFiyatManuelDegistirildi   = 9  //  kurbanın düzenleme alanında fiyatı değiştirilmiş
	KurbanBorcDurumTaksitHayvanAtandi        = 10 //  referans kişi eklenmiş

	VekaletDurumuAlinmadi = 1
	VekaletDurumuAlindi   = 2
)

//KurbanTableName table name
var KurbanTableName string = "kurbanlar"

//TODO: bunun dto kopyası olsun orada validate olsun
//Kurban strcut
type Kurban struct {
	ID            uint64         `gorm:"primary_key;auto_increment" json:"id"`
	UserID        uint64         `gorm:"not null;" json:"user_id"`
	GrupID        uint64         `gorm:"not null;DEFAULT:'0'" json:"grup_id"`
	KisiID        uint64         `gorm:"not null;DEFAULT:'0'" validate:"numeric,omitempty"  json:"kisi_id"`
	Aciklama      string         `gorm:"type:text;" validate:"omitempty,required"  json:"aciklama" `
	VekaletDurumu int            `gorm:"type:smallint ;NOT NULL;DEFAULT:'0'" validate:"required,omitempty"  json:"vekalet"`
	Agirlik       int            `gorm:"type:smallint ; NULL;"  json:"agirlik"`
	Slug          string         `gorm:"size:255 ;null;" json:"slug"`
	KurbanTuru    int            `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" validate:"required,omitempty"  json:"kurbanTuru"`
	Durum         int            `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" validate:"omitempty,required"  json:"durum"`
	BorcDurum     int            `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" validate:"required,omitempty"  json:"BorcDurum"`
	GrupLideri    int            `gorm:"type:smallint ;NOT NULL;DEFAULT:'0'" validate:"omitempty,required"  json:"grupLideri"`
	KurbanFiyati  float64        `gorm:"type:decimal(10,2); NOT NULL; DEFAULT:'0';" validate:"required,omitempty"  json:"kurbanFiyati"`
	Borc          float64        `gorm:"type:decimal(10,2); NOT NULL; DEFAULT:'0';" validate:"numeric"  json:"borc"`
	Alacak        float64        `gorm:"type:decimal(10,2); NOT NULL; DEFAULT:'0';" validate:"numeric,omitempty"  json:"alacak"`
	Bakiye        float64        `gorm:"type:decimal(10,2); NOT NULL; DEFAULT:'0';" validate:"numeric,omitempty"  json:"bakiye"`
	HayvanCinsi   int            `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" validate:"required" json:"hayvanCinsi"`
	Tarih         *time.Time     `json:"tarih"`
	CreatedAt     time.Time      ` json:"created_at"`
	UpdatedAt     time.Time      ` json:"updated_at"`
	DeletedAt     *time.Time     `json:"deleted_at"`
	Odemeler      []dto.Odemeler `gorm:"foreignKey:KurbanID;references:ID"`
}

//BeforeSave init
func (gk *Kurban) BeforeSave() {
	gk.Aciklama = html.EscapeString(strings.TrimSpace(gk.Aciklama))
}

//Prepare init
func (gk *Kurban) Prepare() {
	gk.CreatedAt = time.Now()
	gk.UpdatedAt = time.Now()
}

//TableName override
func (gk *Kurban) TableName() string {
	return KurbanTableName
}

//Validate fluent validation
func (gk *Kurban) Validate() map[string]string {
	var (
		validate *validator.Validate
		uni      *ut.UniversalTranslator
	)
	tr := en.New()
	uni = ut.New(tr, tr)
	trans, _ := uni.GetTranslator("tr")
	validate = validator.New()
	tr_translations.RegisterDefaultTranslations(validate, trans)
	errorLog := make(map[string]string)
	err := validate.Struct(gk)
	fmt.Println(err)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		fmt.Println(errs)
		for _, e := range errs {
			// can translate each error one at a time.
			lng := strings.Replace(e.Translate(trans), e.Field(), "Burası", 1)
			errorLog[e.Field()+"_error"] = e.Translate(trans)
			// errorLog[e.Field()] = e.Translate(trans)
			errorLog[e.Field()] = lng
			errorLog[e.Field()+"_valid"] = "is-invalid"
		}
	}
	return errorLog
}

/*
Kurban Türleri
11-12 hisseli büyükbaş
9- kurban bayramı küçükbaş
1- adak
2-akika
3-şükür
4-SAHİBİNİN NİYETİNE
5-hayır
6- nafile
7- şifa

*/

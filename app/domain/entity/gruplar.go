package entity

import (
	// "stncCms/app/domain/dto"
	"html"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	tr_translations "gopkg.in/go-playground/validator.v9/translations/tr"
)

var GruplarTableName string = "gruplar"

const (
	GruplarDurumKurbanBayramiDisiKesim        = 1 // kurban bayramına ait değil büyükbaş yada küçük baş gibi yani direk kurban girişi, ve aynı zamanda diğer kurban çeşitleri
	GruplarDurumGrupOlusmusHayvanYok          = 2 // grup oluşmuş ama kimse atanmamış , yani kesimlik inek verilmemiştir
	GruplarDurumGrupOlusmusKesimlikHayvaniVar = 3 // grup atanmış yani bir kesimlik inek verilmiş
	GruplarDurumKurbanKesimiTamamlanmis       = 4 // kurban kesimi tamamlanmış ?? bunu belirleyen bir durum observer falan olmalı
)

//Gruplar struct
type Gruplar struct {
	ID                 uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID             uint64    `gorm:"NOT NUL;" json:"userId"`
	HayvanBilgisiID    uint64    `gorm:"NOT NUL;DEFAULT:'0'" json:"hayvan_id"  validate:"numeric,omitempty"`
	GrupAdi            string    `gorm:"size:255; null;" json:"grup_adi"`
	KesimSiraNo        int       `gorm:"type:smallint ;NOT NULL;DEFAULT:'0'" validate:"required,max=32767"`
	HissedarAdet       int       `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" validate:"numeric,required"`
	Aciklama           string    `gorm:"type:text ;" json:"aciklama" validate:"omitempty"`
	SatisFiyatTuru     int       `gorm:"type:smallint ;NOT NULL;DEFAULT:'0'" validate:"omitempty,required"`        //ayalardaki satış fiyatlraında hangisi kullanımış
	SatisFiyati        float64   `gorm:"type:decimal(10,2); NOT NULL; DEFAULT:'0';" validate:"omitempty,numeric" ` //ayarlardaki kilo başı satış fiyatının hangisi seçilmiş
	Siralama           int       `gorm:"type:smallint ;NOT NULL;DEFAULT:'0'" validate:"required"`
	AgirlikTipi        int       `gorm:"type:smallint ;NOT NULL;DEFAULT:'0'" validate:"required" json:"kurban_fiyati_kilo_tipi"` // yüksek ağırlık , orta ağırlık , düşük ağırlık gibi
	Durum              int       `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" validate:"omitempty,required"`
	Slug               string    `gorm:"size:255; null;" json:"slug"`
	ToplamKurbanFiyati float64   `gorm:"type:decimal(10,2); NOT NULL; DEFAULT:'0';" validate:"required,numeric" `
	KurbanBayramiYili  int       `gorm:"type:smallint ;NOT NULL;" validate:"omitempty,required,numeric"` //otomaitk atılacak isitem tarafından ayarlradan ceksin
	Agirlik            int       `gorm:"type:smallint ; NULL;"`
	CreatedAt          time.Time ` json:"created_at"`
	UpdatedAt          time.Time ` json:"updated_at"`

	DeletedAt *time.Time `json:"deleted_at"`
	Kurban    []Kurban   `gorm:"foreignKey:GrupID;references:ID"`
}

//BeforeSave init
func (gk *Gruplar) BeforeSave() {
	gk.Aciklama = html.EscapeString(strings.TrimSpace(gk.Aciklama))
}

//Prepare init
func (gk *Gruplar) Prepare() {
	gk.CreatedAt = time.Now()
	gk.UpdatedAt = time.Now()
}

//TableName override
func (gk *Gruplar) TableName() string {
	return "gruplar"
}

//Validate fluent validation
func (gk *Gruplar) Validate() map[string]string {
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
	//fmt.Println(err)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		//fmt.Println(errs)
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

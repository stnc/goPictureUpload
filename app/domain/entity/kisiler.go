package entity

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	tr_translations "gopkg.in/go-playground/validator.v9/translations/tr"
)

//KisilerTableName table name
var KisilerTableName string = "kisiler"

//TODO: bunun dto kopyası olsun orada validate olsun
//Kisiler strcut
type Kisiler struct {
	ID            uint64     `gorm:"primary_key;auto_increment" json:"id"`
	UserID        uint64     `gorm:"not null;" json:"user_id"`
	AdSoyad       string     `gorm:"size:255;not null;" validate:"required,omitempty,min=2,max=50" json:"adSoyad" `
	Telefon       string     `gorm:"size:255;not null;"  validate:"required,omitempty" json:"tel"`
	Email         string     `gorm:"type:varchar(255);" validate:"omitempty,email"  json:"emailAdres"`
	Aciklama      string     `gorm:"type:text;" validate:"omitempty,required"  json:"aciklama" `
	ReferansKisi1 uint64     `gorm:"not null;" json:"referansKisi1"`
	Adres         string     `gorm:"type:text;" validate:"omitempty,required"  json:"adres" `
	Durum         int        `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" validate:"omitempty,required"`
	DovizCinsi    int        `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" validate:"omitempty,required"`
	CreatedAt     time.Time  ` json:"created_at"`
	UpdatedAt     time.Time  ` json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
	Kurbanlar     []Kurban   `gorm:"foreignKey:KisiID;references:ID"`
}

//BeforeSave init
func (gk *Kisiler) BeforeSave() {
	gk.AdSoyad = html.EscapeString(strings.TrimSpace(gk.AdSoyad))
	gk.Telefon = html.EscapeString(strings.TrimSpace(gk.Telefon))
	gk.Adres = html.EscapeString(strings.TrimSpace(gk.Adres))
	gk.Aciklama = html.EscapeString(strings.TrimSpace(gk.Aciklama))
}

//Prepare init
func (gk *Kisiler) Prepare() {
	gk.AdSoyad = html.EscapeString(strings.TrimSpace(gk.AdSoyad))
	gk.Adres = html.EscapeString(strings.TrimSpace(gk.Adres))
	gk.CreatedAt = time.Now()
	gk.UpdatedAt = time.Now()
}

//TableName override
func (gk *Kisiler) TableName() string {
	return KisilerTableName
}

//Validate fluent validation
func (gk *Kisiler) Validate() map[string]string {
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

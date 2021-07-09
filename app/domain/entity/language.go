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

//Languages strcut
type Languages struct {
	ID            uint64     `gorm:"primary_key;auto_increment" json:"id"`
	UserID        uint64     `gorm:"not null;" json:"userId"`
	PostID        uint64     `gorm:"not null;" json:"postId"`
	Language      string     `gorm:"type:char(20) ;not null;" validate:"required"  json:"lang"`
	RelationsType string     `gorm:"type:text ;" json:"relations_type" validate:"required"`
	Type          int        `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" validate:"required"`
	Status        int        `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" validate:"required"`
	CreatedAt     time.Time  ` json:"created_at"`
	UpdatedAt     time.Time  ` json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

//BeforeSave init
func (l *Languages) BeforeSave() {
	l.Language = html.EscapeString(strings.TrimSpace(l.Language))
}

//Prepare init
func (l *Languages) Prepare() {
	l.Language = html.EscapeString(strings.TrimSpace(l.Language))
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()
}

//Validate fluent validation
func (l *Languages) Validate() map[string]string {
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

	err := validate.Struct(l)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		fmt.Println(errs)
		for _, e := range errs {
			// can translate each error one at a time.
			lng := strings.Replace(e.Translate(trans), e.Field(), "Burası", 1)
			errorLog[e.Field()+"_error"] = e.Translate(trans)
			errorLog[e.Field()] = lng
			errorLog[e.Field()+"_valid"] = "is-invalid"
		}
	}
	return errorLog
}

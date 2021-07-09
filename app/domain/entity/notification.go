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

//notification notların oldugu alan
type Notification struct {
	ID         uint64     `gorm:"primary_key;auto_increment" json:"id"`
	UserID     uint64     `gorm:"not null;" json:"user_id"`
	SourceID   uint64     `gorm:"not null;" json:"source_id"`
	sourceType string     `gorm:"size:50 ; null;" validate:"required"`
	Type       int        `gorm:"smallint:6;DEFAULT:'0'" `
	Read       int        `gorm:"smallint:1;DEFAULT:'1'" `
	Trash      int        `gorm:"smallint:1;DEFAULT:'1'" `
	Content    string     `gorm:"type:text ;" validate:"required"  json:"content"`
	CreatedAt  time.Time  ` json:"created_at"`
	UpdatedAt  time.Time  ` json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

//BeforeSave init
func (f *Notification) BeforeSave() {
	f.Content = html.EscapeString(strings.TrimSpace(f.Content))
}

/*
func (post *Post) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
 }
*/

//Prepare init
func (f *Notification) Prepare() {
	f.Content = html.EscapeString(strings.TrimSpace(f.Content))
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
}

//Validate fluent validation
func (f *Notification) Validate() map[string]string {
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

	err := validate.Struct(f)
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

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

//Options seçenekelr
type Options struct {
	OptionID    int        `gorm:"primary_key;auto_increment"`
	OptionName  string     `gorm:"size:255 ;not null;" validate:"required"`
	OptionValue string     `gorm:"type:text ;"  validate:"omitempty,required"`
	Status      int        `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" validate:"required"`
	CreatedAt   time.Time  ` json:"created_at"`
	UpdatedAt   time.Time  ` json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

//BeforeSave init
func (f *Options) BeforeSave() {
	f.OptionName = html.EscapeString(strings.TrimSpace(f.OptionName))
	f.OptionValue = html.EscapeString(strings.TrimSpace(f.OptionValue))
}

/*
func (post *Post) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
 }
*/

//Prepare init
func (f *Options) Prepare() {
	f.OptionName = html.EscapeString(strings.TrimSpace(f.OptionName))
	f.OptionValue = html.EscapeString(strings.TrimSpace(f.OptionValue))
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
}

//Validate fluent validation
func (f *Options) Validate() map[string]string {
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

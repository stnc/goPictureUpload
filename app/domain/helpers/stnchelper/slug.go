package stnchelper

import (
	"fmt"
	"math/rand"
	"stncCms/app/domain/entity"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/gosimple/slug"
	"gopkg.in/go-playground/validator.v9"
	tr_translations "gopkg.in/go-playground/validator.v9/translations/tr"
)

/*
usage
	rand.Seed(time.Now().UnixNano())

	fmt.Println(RandSlugV1(5))
*/
//RandSlugV1 random slug
func RandSlugV1(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//Slugify slug genereate
func Slugify(title string, size int) string {
	slug.MaxLength = size
	return slug.MakeLang(title, "tr")
}

//GenericName for uplaod generic name
func GenericName(title string, size int) string {
	var name string
	name = Slugify(title, size)
	currentTime := time.Now()
	dateadd := currentTime.Format("15_04_05")
	return name + "_" + dateadd
}

//Validate fluent validation
func validateTESSST(f entity.Post) map[string]string {
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

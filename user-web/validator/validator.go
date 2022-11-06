package validator

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
)

var (
	pattern   = `^1(3\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\d|9[0-35-9])\d{8}$`
	mobileReg *regexp.Regexp
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 正则表达式
	if mobileReg == nil {
		mobileReg = regexp.MustCompile(pattern)
	}
	if ok := mobileReg.MatchString(mobile); !ok {
		return false
	}
	return true
}
func ValidateTrans(v *validator.Validate, trans ut.Translator, name, text string) {
	err := v.RegisterTranslation(name, trans, func(ut ut.Translator) error {
		return ut.Add("mobile", fmt.Sprintf("{0} %s", text), true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(name, fe.Field())
		return t
	})
	if err != nil {
		log.Println("can not register validate translate", "name", name)
	}
}

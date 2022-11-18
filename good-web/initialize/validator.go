package initialize

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"web-api/good-web/global"
)

// 如果有自定义 匹配格式
// 默认的错误 会是英文 在这里可以转化问自己定义的 text 错误信息 name 是我们的匹配规则类似 required
func InitValidator(locale string) {
	uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
	global.Trans, _ = uni.GetTranslator(locale)
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		switch locale {
		case "zh":
			_ = zh_translations.RegisterDefaultTranslations(v, global.Trans)
		case "en":
			_ = en_translations.RegisterDefaultTranslations(v, global.Trans)
		default:
			_ = zh_translations.RegisterDefaultTranslations(v, global.Trans)
			break
		}
	}
}

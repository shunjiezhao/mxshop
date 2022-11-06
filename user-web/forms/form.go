package forms

import (
	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
	"strings"
	"web-api/user-web/global"
)

type ValidError struct {
	//Key     string
	Message string
}

type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

func BindAndValid(c *gin.Context, v interface{}) (ValidErrors, bool) {
	var errs ValidErrors
	err := c.ShouldBind(v)
	if err != nil {
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			return errs, false
		}
		for _, value := range verrs.Translate(global.Trans) {
			errs = append(errs, &ValidError{
				//Key:     key,
				Message: value,
			})
		}

		return errs, false
	}

	return nil, true
}

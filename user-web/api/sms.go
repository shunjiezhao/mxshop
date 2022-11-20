package api

import (
	"context"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"net/http"
	"time"
	"web-api/user-web/forms"
	"web-api/user-web/global"
)

func SendMessage(c *gin.Context) {
	form := forms.SendSmsForm{}
	if err, ok := forms.BindAndValid(c, &form); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err,
		})
		return
	}
	switch {

	}
	code := rand.Intn(10000)
	err := _main(form.Mobile, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "验证码发送失败",
		})
		return
	}

	global.RedisClient.Set(context.Background(), form.Mobile, code, time.Minute*global.ServerConfig.RedisInfo.ExpireMin)
	zap.L().Info("set value", zap.String("mobile", form.Mobile), zap.Int("code", code))

	c.JSON(http.StatusOK, gin.H{
		"msg": "验证码成功",
	})
	return
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func _main(mobile string, code int) (_err error) {

	client, _err := CreateClient(tea.String("LTAI5tL58ops2uUL2xc1Jmmd"), tea.String("qYq7YupTsAniJxsuB6ggB54kgSZzL9"))
	if _err != nil {
		log.Println(_err)
		return _err
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("阿里云短信测试"),
		TemplateCode:  tea.String("SMS_154950909"),
		PhoneNumbers:  tea.String(mobile),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%d\"}", code)),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				log.Println(r)
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, _err = client.SendSmsWithOptions(sendSmsRequest, runtime)
		if _err != nil {
			log.Println(_err)
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
			log.Println(error.Message)
		}
		// 如有需要，请打印 error
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			log.Println(error.Message)
			return _err
		}
	}
	return _err
}

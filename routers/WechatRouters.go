package routers

import (
	"encoding/json"
	"fmt"
	"gin_spring/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type userCode struct {
	Code       string `json:"code"`
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
}

func WechatRouter(router *gin.Engine) {
	var code userCode
	APPID := "wxb5ffc058d0c82c22"
	SECRET := "27735d75d1fe375f83c99b993f541f28"
	//APPID := "wx451e415a7b9a69bf"
	//SECRET := "1f5dbde56e2e6edbe7af882e4d83e8fa"
	router.POST("/api/openid", func(context *gin.Context) {
		util.JsonBind(context, func() {
			wechatLoginUrl := "https://api.weixin.qq.com/sns/jscode2session?appid=" + APPID + "&secret=" + SECRET + "&js_code=" + code.Code + "&grant_type=authorization_code"
			request, err := http.NewRequest("GET", wechatLoginUrl, nil)
			if err != nil {
				fmt.Println(err.Error())
			}
			var response *http.Response
			response, err = http.DefaultClient.Do(request)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			//fmt.Println(string(body))
			_ = json.Unmarshal(body, &code)
			context.JSON(http.StatusOK, code)
		}, &code)
	})
}

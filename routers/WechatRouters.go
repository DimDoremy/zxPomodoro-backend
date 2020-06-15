package routers

import (
	"encoding/json"
	"fmt"
	"gin_spring/model"
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

// WechatRouter doc
// @Summary 请求微信登录授权
// @Description 用于服务器请求微信服务器进行用户登录授权使用
// @Tags wechat-router
// @Accept json
// @Produce json
// @Param handBook body userCode true "用户的登录包"
// @Success 200 {object} userCode "请求结果"
// @Failure 400 {object} model.MessageBind "请求没到达"
// @Router /openid [post]
func WechatRouter(router *gin.Engine) {
	var code userCode
	APPID := "wxb5ffc058d0c82c22"
	SECRET := "27735d75d1fe375f83c99b993f541f28"
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
			if response == nil {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "request failed!"})
				return
			}
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			//fmt.Println(string(body))
			_ = json.Unmarshal(body, &code)
			context.JSON(http.StatusOK, code)
		}, &code)
	})
}

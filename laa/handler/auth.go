package handler

import (
	"authorization.setruth.com/laa/model"
	"authorization.setruth.com/laa/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ActivateAuth(context *gin.Context) {
	authCode := context.GetHeader("AuthCode")
	if authCode == "" {
		context.JSON(http.StatusBadRequest, model.BaseRes[struct{}]{
			Msg:  "不存在的授权码",
			Data: nil,
		})
		return
	}
	authData, err := util.VerificationAuthCode(authCode)
	if err != nil {
		context.JSON(http.StatusBadRequest, model.BaseRes[struct{}]{
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	if authData.UniqueCode != model.UniqueCodeCache {
		context.JSON(http.StatusUnauthorized, model.BaseRes[struct{}]{
			Msg:  "授权码不是授权此设备的",
			Data: nil,
		})
		return
	}
	err = util.UpsertAuthCode(authCode)
	if err != nil {
		context.JSON(http.StatusInternalServerError, model.BaseRes[struct{}]{
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	model.AuthDetailCache = &model.AuthDetail{
		AuthCode:     authCode,
		EndTimestamp: authData.EndTimestamp,
		UniqueCode:   model.UniqueCodeCache,
	}
	context.JSON(http.StatusOK, model.BaseRes[struct{}]{
		Msg:  "授权成功",
		Data: nil,
	})
}

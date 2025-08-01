package main

import (
	"authorization.setruth.com/laa/handler"
	"authorization.setruth.com/laa/model"
	"authorization.setruth.com/laa/resource"
	"authorization.setruth.com/laa/task"
	"authorization.setruth.com/laa/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRoutes(router *gin.Engine) {
	//WebGUI
	router.StaticFS("/web", http.FS(resource.WebStatic))
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web/index.html")
	})
	router.Static("/static", "resource/web")

	//API
	rootPath := router.Group("/api")
	{
		rootPath.GET("/uniqueCode", func(context *gin.Context) {
			context.JSON(http.StatusOK, model.BaseRes[string]{
				Msg:  "获取成功",
				Data: &model.UniqueCodeCache,
			})
		})
		authPath := rootPath.Group("/auth")
		{
			authPath.POST("", handler.ActivateAuth)
			authPath.DELETE("", func(context *gin.Context) {
				util.ClearAuthCode()
				model.AuthDetailCache = nil
				context.JSON(http.StatusOK, model.BaseRes[struct{}]{
					Msg:  "清空授权成功",
					Data: nil,
				})
			})
			authPath.GET("", func(context *gin.Context) {
				if model.AuthDetailCache == nil {
					emptyCode := ""
					context.JSON(http.StatusOK, model.BaseRes[string]{
						Msg:  "获取成功",
						Data: &emptyCode,
					})
				} else {
					context.JSON(http.StatusOK, model.BaseRes[string]{
						Msg:  "获取成功",
						Data: &model.AuthDetailCache.AuthCode,
					})
				}
			})
		}
		exposePath := rootPath.Group("/status")
		{
			exposePath.GET("", func(context *gin.Context) {
				context.JSON(http.StatusOK, model.BaseRes[model.AuthStatus]{
					Msg:  "获取成功",
					Data: task.GetAuthStatus(),
				})
			})
			exposePath.GET("/subscribe", func(context *gin.Context) {
				context.Header("Content-Type", "text/event-stream;charset=utf-8")
				context.Header("Cache-Control", "no-cache")
				context.Header("Connection", "keep-alive")
				context.Status(http.StatusOK)
				task.Subscribe(context)
			})
		}
	}

}

package main

import (
	"authorization.setruth.com/ams/entity"
	"authorization.setruth.com/ams/model"
	"authorization.setruth.com/ams/util"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

var webContent embed.FS

func InitRoutes(context *gin.Engine, db *gorm.DB) {
	//WebGUI
	context.Static("/web", "resource/web")
	context.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web/index.html")
	})

	//API
	rootPath := context.Group("/api")
	{
		authListPath := rootPath.Group("/authList")
		{
			authListPath.POST("add", func(context *gin.Context) {
				var addAuthRecordDTO model.AddAuthRecordDTO
				if err := context.ShouldBindJSON(&addAuthRecordDTO); err != nil {
					context.JSON(http.StatusBadRequest, model.BaseRes[struct{}]{
						Msg:  "添加授权失败,数据结构有误",
						Data: nil,
					})
					return
				}
				var repeatCount int64
				db.Model(&entity.AuthRecord{}).Where("name = ? OR unique_code =?", addAuthRecordDTO.Name, addAuthRecordDTO.UniqueCode).Count(&repeatCount)
				if repeatCount > 0 {
					context.JSON(http.StatusConflict, model.BaseRes[struct{}]{
						Msg:  "记录的设备码或名称已存在，请勿重复添加。",
						Data: nil,
					})
					return
				}
				record := entity.AuthRecord{
					Name:         addAuthRecordDTO.Name,
					UniqueCode:   addAuthRecordDTO.UniqueCode,
					EndTimestamp: addAuthRecordDTO.EndTimestamp,
				}
				result := db.Create(&record)
				if result.Error != nil {
					log.Printf("添加授权失败:%s", result.Error)
					context.JSON(http.StatusInternalServerError, model.BaseRes[struct{}]{
						Msg:  "添加授权失败,请稍后重试",
						Data: nil,
					})
					return
				}
				fmt.Printf("添加授权成功：%d", record.ID)
				context.JSON(http.StatusOK, model.BaseRes[struct{}]{
					Msg:  "添加授权成功",
					Data: nil,
				})
			})
			authListPath.POST("updateEndTimestamp", func(context *gin.Context) {
				var updateEndTimestampDTO model.UpdateEndTimestampDTO
				if err := context.ShouldBindJSON(&updateEndTimestampDTO); err != nil {
					context.JSON(http.StatusBadRequest, model.BaseRes[struct{}]{
						Msg:  "更新失败，提交的更新数据结构有误",
						Data: nil,
					})
					return
				}
				db.Model(&entity.AuthRecord{}).Where("id = ?", updateEndTimestampDTO.ID).Update("end_timestamp", updateEndTimestampDTO.EndTimestamp)
				context.JSON(http.StatusOK, model.BaseRes[struct{}]{
					Msg:  "更新成功",
					Data: nil,
				})
				return
			})
			authListPath.DELETE(":id", func(context *gin.Context) {
				idStr := context.Param("id")
				id, err := strconv.Atoi(idStr)
				if err != nil {
					context.JSON(http.StatusBadRequest, model.BaseRes[string]{
						Msg:  "删除授权失败,ID格式有误",
						Data: nil,
					})
					return
				}
				db.Delete(&entity.AuthRecord{}, id)
				context.JSON(http.StatusOK, model.BaseRes[struct{}]{
					Msg:  "删除授权成功",
					Data: nil,
				})
				return
			})
			authListPath.GET("/all", func(context *gin.Context) {
				var records []entity.AuthRecord
				db.Find(&records)
				context.JSON(http.StatusOK, model.BaseRes[[]entity.AuthRecord]{
					Msg:  "获取授权列表成功",
					Data: &records,
				})
			})
			authListPath.GET("/authCode/:id", func(context *gin.Context) {
				idStr := context.Param("id")
				id, err := strconv.Atoi(idStr)
				if err != nil {
					context.JSON(http.StatusBadRequest, model.BaseRes[string]{
						Msg:  "请求出错，ID格式有误",
						Data: nil,
					})
					return
				}
				var record entity.AuthRecord
				err = db.Model(&entity.AuthRecord{}).Where("id = ?", id).First(&record).Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						context.JSON(http.StatusNotFound, model.BaseRes[struct{}]{
							Msg:  fmt.Sprintf("ID为 %d 的记录未找到。", id),
							Data: nil,
						})
						return
					} else {
						log.Printf("查询 ID 为 %d 的记录时发生数据库错误: %v", id, err)
						context.JSON(http.StatusInternalServerError, model.BaseRes[struct{}]{
							Msg:  "数据库查询失败，请稍后再试。",
							Data: nil,
						})
						return
					}
				}
				code, err := util.GenerateAuthorizationCode(record.UniqueCode, record.EndTimestamp)
				if err != nil {
					log.Printf("授权码创建失败:%v", err)
					context.JSON(http.StatusInternalServerError, model.BaseRes[struct{}]{
						Msg:  "授权码创建失败,请稍后重试",
						Data: nil,
					})
					return
				}
				context.JSON(http.StatusOK, model.BaseRes[string]{
					Msg:  "获取授权码成功",
					Data: &code,
				})
			})
		}
		rootPath.GET("generateKeys", func(context *gin.Context) {

			context.Header("Access-Control-Allow-Origin", "*")
			// 允许所有常用 HTTP 方法
			context.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
			// 允许所有常用请求头
			context.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, User-Agent, Authorization, X-Requested-With")
			// 允许浏览器访问的响应头（如果你的后端有自定义响应头，前端需要读取）
			context.Header("Access-Control-Expose-Headers", "Content-Length")
			// 是否允许发送 Cookie 或 HTTP 认证信息
			context.Header("Access-Control-Allow-Credentials", "true")
			// 预检请求（OPTIONS）的缓存时间，这里设置为 1 天 (86400 秒)
			context.Header("Access-Control-Max-Age", "86400")
			aesKey := make([]byte, 32)
			_, err := rand.Read(aesKey)
			if err != nil {
				fmt.Printf("生成 AES 密钥失败: %v\n", err)
				context.JSON(http.StatusInternalServerError, model.BaseRes[struct{}]{
					Msg:  "AES 密钥生成失败，请稍后重试",
					Data: nil,
				})
				return
			}
			aesKeyBase64 := base64.StdEncoding.EncodeToString(aesKey)
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				fmt.Printf("生成 RSA 密钥对失败: %v\n", err)
				context.JSON(http.StatusInternalServerError, model.BaseRes[struct{}]{
					Msg:  "RSA 密钥对生成失败，请稍后重试",
					Data: nil,
				})
				return
			}
			publicKey := &privateKey.PublicKey
			pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
			if err != nil {
				fmt.Printf("编码 RSA 公钥失败: %v\n", err)
				context.JSON(http.StatusInternalServerError, model.BaseRes[struct{}]{
					Msg:  "RSA 公钥编码失败，请稍后重试",
					Data: nil,
				})
				return
			}
			publicKeyBase64 := base64.StdEncoding.EncodeToString(pubBytes)
			privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
			if err != nil {
				fmt.Printf("编码 RSA 私钥失败: %v\n", err)
				context.JSON(http.StatusInternalServerError, model.BaseRes[struct{}]{
					Msg:  "RSA 私钥编码失败，请稍后重试",
					Data: nil,
				})
				return
			}
			privateKeyBase64 := base64.StdEncoding.EncodeToString(privBytes)
			context.JSON(http.StatusOK, model.BaseRes[model.Keys]{
				Msg: "获取密钥成功",
				Data: &model.Keys{
					AESKey:        aesKeyBase64,
					RSAPublicKey:  publicKeyBase64,
					RSAPrivateKey: privateKeyBase64,
				},
			})
		})
	}

}

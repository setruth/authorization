package main

import (
	"authorization.setruth.com/ams/entity"
	"authorization.setruth.com/ams/model"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	initRsaPrivateKey()
	db, err := InitDB()
	if err != nil {
		println("数据库初始化失败: %s\n", err)
		return
	}
	startServer(db)
}
func initRsaPrivateKey() {
	rsaPrivateKeyBase64 := os.Getenv(model.RsaPrivateEnvKey)
	if rsaPrivateKeyBase64 == "" {
		log.Fatalf("部署授权服务的机器的环境变量没有配置RSA私钥，请检查")
	}
	rsaPrivateKey, err := base64.StdEncoding.DecodeString(rsaPrivateKeyBase64)
	if err != nil {
		log.Fatalf("你的RSA私钥不是标准的BASE64编码内容无法解码: %v", err)
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(rsaPrivateKey)
	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		log.Fatal("进行断言发现并不是PublicKey类型")
	}
	model.UpdateRsaPrivateKey(rsaKey)
}
func startServer(db *gorm.DB) {
	port := "1023"
	context := gin.New()
	context.Use(gin.Logger())
	context.Use(gin.Recovery())
	InitRoutes(context, db)
	err := context.Run(":" + port)
	if err != nil {
		log.Printf("启动失败: %s\n", err)
	}
}
func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %s\n", err)
	}
	fmt.Println("数据库连接成功")
	err = db.AutoMigrate(&entity.AuthRecord{})
	if err != nil {
		return nil, fmt.Errorf("表更新失败: %s\n", err)
	}
	return db, nil
}

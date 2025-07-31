package main

import (
	"authorization.setruth.com/laa/model"
	"authorization.setruth.com/laa/task"
	util "authorization.setruth.com/laa/util"
	"bufio"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	uniqueCode, err := util.GetUniqueCode()
	if err != nil {
		log.Printf("获取设备的唯一码失败，无法启动服务: %s\n", err)
		log.Println("请按回车键退出...")
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		return
	}
	model.UniqueCodeCache = uniqueCode
	model.TaskWg.Add(1)
	go task.LaunchAuthStatusCheckTask()
	initAuthStatus(uniqueCode)
	startServer()
}
func initAuthStatus(uniqueCode string) {
	authCode, _ := util.ReadAuthCode()
	if authCode == "" {
		return
	}
	authData, err := util.VerificationAuthCode(authCode)
	if err != nil {
		util.ClearAuthCode()
		log.Println("存储的授权码不正确,已清空")
		return
	}
	if uniqueCode != authData.UniqueCode {
		util.ClearAuthCode()
		log.Println("授权码的唯一标识与机器不符,已清空")
		return
	}
	model.AuthDetailCache = &model.AuthDetail{
		AuthCode:     authCode,
		UniqueCode:   authData.UniqueCode,
		EndTimestamp: authData.EndTimestamp,
	}
}
func startServer() {
	port := "1022"
	context := gin.New()
	context.Use(gin.Logger())
	context.Use(gin.Recovery())
	InitRoutes(context)
	err := context.Run(":" + port)
	if err != nil {
		log.Printf("启动失败: %s\n", err)
	}
	close(model.TaskStopChan)
	model.TaskWg.Wait()
}

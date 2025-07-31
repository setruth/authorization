package task

import (
	"authorization.setruth.com/laa/model"
	"authorization.setruth.com/laa/state"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"log"
	"net/http"
	"time"
)

var (
	authStatusFlow = state.NewGoStateFlow(model.AuthStatus{
		Tag:          model.Unauthorized,
		EndTimestamp: model.EndTimestampNil,
	})
)

func LaunchAuthStatusCheckTask() {
	defer model.TaskWg.Done()
	log.Println("认证状态检查任务已启动...")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			currentUnixTime := time.Now().UnixMilli()
			if model.AuthDetailCache == nil {
				authStatusFlow.Set(model.AuthStatus{
					Tag:          model.Unauthorized,
					EndTimestamp: model.EndTimestampNil,
				})
			} else if model.AuthDetailCache.EndTimestamp > currentUnixTime {
				authStatusFlow.Set(model.AuthStatus{
					Tag:          model.Authorized,
					EndTimestamp: model.AuthDetailCache.EndTimestamp,
				})
			} else {
				authStatusFlow.Set(model.AuthStatus{
					Tag:          model.Expire,
					EndTimestamp: model.AuthDetailCache.EndTimestamp,
				})
			}
		case <-model.TaskStopChan:
			log.Println("服务关闭，正在退出任务。")
			return
		}
	}
}

func GetAuthStatus() *model.AuthStatus {
	authStatus := authStatusFlow.Get()
	return &authStatus
}

func Subscribe(context *gin.Context) {
	flusher, ok := context.Writer.(http.Flusher)
	if !ok {
		log.Printf("无法进行SSE的写入")
		return
	}
	subscribe, closeSubscribe := authStatusFlow.Subscribe()
	for {
		select {
		case <-context.Request.Context().Done():
			log.Printf("客户端关闭，取消状态监听和发送")
			closeSubscribe()
			return
		case newStatus := <-subscribe:
			jsonByte, err := json.Marshal(newStatus)
			if err != nil {
				log.Printf("错误的授权状态解析:%s\n", err)
			}
			sseMsg := fmt.Sprintf("data: %s\n\n", string(jsonByte))
			_, err = context.Writer.WriteString(sseMsg)
			if err != nil {
				log.Printf("sse订阅发送出错:%s\n", err)
				return
			}
		}
		flusher.Flush()
	}
}

package blive

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

var (
	logger = logrus.WithField("service", "blive")
)

func StartWebSocket(ctx context.Context, wg *sync.WaitGroup) {

	websocketHost := os.Getenv("WEBSOCKET_URL")

	logrus.Debugf("prepare to connect %v", websocketHost)
	con, _, err := websocket.DefaultDialer.Dial(websocketHost, nil)

	if err != nil {
		logger.Errorf("連線到 Websocket %s 時出現錯誤: %v", websocketHost, err)
		logger.Warnf("十秒後重試")
		<-time.After(time.Second * 10)
		StartWebSocket(ctx, wg)
		return
	}

	logger.Infof("連線到 Websocket %s 成功", websocketHost)

	con.SetCloseHandler(func(code int, text string) error {
		return con.WriteMessage(websocket.CloseMessage, nil)
	})

	wg.Add(1)
	go onReceiveMessage(ctx, con, wg)
}

func onReceiveMessage(ctx context.Context, conn *websocket.Conn, wg *sync.WaitGroup) {
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Errorf("關閉 Websocket 時出現錯誤: %v", err)
		} else {
			logger.Debugf("連接關閉成功。")
		}
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("正在關閉 Websocket...")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "stop"))
			if err != nil {
				logger.Errorf("發送 websocket 關閉訊息時出現錯誤: %v", err)
			}
			return
		default:
			_, message, err := conn.ReadMessage()
			// Error
			if err != nil {
				logger.Errorf("Websocket 嘗試讀取消息時出現錯誤: %v", err)
				go retryDelay(ctx, wg)
				return
			}
			go handleMessage(message)
		}
	}
}

func retryDelay(ctx context.Context, wg *sync.WaitGroup) {
	logger.Warnf("五秒後重連...")
	<-time.After(time.Second * 5)
	StartWebSocket(ctx, wg)
}

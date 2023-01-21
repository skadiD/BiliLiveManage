package room

import (
	"errors"
	"github.com/fexli/logger"
	"github.com/gorilla/websocket"
	"os"
	"os/signal"
	"time"
)

// 创建 websocket 链接
func (room *Client) Join() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	var err = errors.New("[ID: " + room.ID + "] 进入直播间失败")
	for _, chatHost := range room.Info.prefixHost() {
		room.Conn, _, err = websocket.DefaultDialer.Dial(chatHost.String(), nil)
		if err != nil {
			logger.RootLogger.Error(logger.WithContent("[ID: "+room.ID+"] 连接弹幕服务器失败", err))
			continue
		}
		logger.RootLogger.System(logger.WithContent("[ID: " + room.ID + "] 连接弹幕服务器成功"))
		logger.RootLogger.Debug(logger.WithContent("[ID: "+room.ID+"]", chatHost.String()))
		break
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			packType, message, readErr := room.Conn.ReadMessage()
			if readErr != nil {
				logger.RootLogger.Warning(logger.WithContent("[ID: "+room.ID+"] 直播间断开：", readErr))
				//retry++
				//logger.RootLogger.Notice(logger.WithContent("正在尝试重连...第", retry, "次"))
				//Init(sockets, host, roomID, retry)
				//TODO: 掉线重连
				return
			}
			switch packType {
			case websocket.PingMessage:
				logger.RootLogger.Debug(logger.WithContent("PING", message))
				return
			case websocket.TextMessage:
				// 从来没见到过
				logger.RootLogger.Debug(logger.WithContent("[ID: "+room.ID+"] 收到消息", message))
				return
			case websocket.BinaryMessage:
				// logger.RootLogger.Debug(logger.WithContent("[ID: "+room.ID+"] 原始消息", message))
				go room.execPacket(message)
			}
			//logger.RootLogger.Debug(logger.WithContent("[ID: "+room.ID+"] 收到消息类型:", packType))
			//go event.ParseCore(message)
		}
	}()
	room.SendChan <- room.FirstHandlePack() // 连接认证
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case data := <-room.SendChan:
			go func() {
				_err := room.Conn.WriteMessage(websocket.BinaryMessage, data)
				if _err != nil {
					logger.RootLogger.Warning(logger.WithContent("[ID: "+room.ID+"] 直播间发包异常", _err))
					return
				}
				logger.RootLogger.System(logger.WithContent("[ID: " + room.ID + "] 直播间发送数据"))
			}()
		case <-interrupt:
			logger.RootLogger.System(logger.WithContent("[ID: " + room.ID + "] 直播间正在退出"))
			_err := room.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if _err != nil {
				logger.RootLogger.Warning(logger.WithContent("[ID: "+room.ID+"] 直播间退出异常", _err))
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
				logger.RootLogger.Warning(logger.WithContent("[ID: " + room.ID + "] 直播间异常警告"))
			}
			return
		}
	}
}

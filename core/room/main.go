package room

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/fexli/logger"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/skadiD/BiliLiveManage/utils/networkUtils"
	"github.com/tidwall/gjson"
	"io"
	"time"
	"unicode"
)

// 初始化直播间状态
func getRoom(id string) *Init {
	var url = "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id="
	ret, err := networkUtils.Get(url+id+"&type=0", nil)
	var status Init
	if err != nil {
		logger.RootLogger.Warning(logger.WithContent("[ID: " + id + "] 初始化直播间信息失败"))
		logger.RootLogger.Debug(logger.WithContent(err))
		return nil
	}
	// fmt.Println(string(ret))
	json.Unmarshal(ret, &status)
	return &status
}

// 创建直播间 Client
func Create(id string) *Client {
	roomInfo := getRoom(id)
	if roomInfo.Code == 0 {
		//logger.RootLogger.Debug(logger.WithContent("[ID: " + id + "] token: " + roomInfo.Data.Token))
		//logger.RootLogger.Debug(logger.WithContent("[ID: " + id + "] host: " + roomInfo.Data.HostList[0].Host))
		// fmt.Println(roomInfo)
		client := Client{
			ID:       id,
			Info:     roomInfo,
			Auth:     nil,
			SendChan: make(chan []byte, 1e4),
		}
		// client.FirstHandlePack()
		client.Join()
		return &client
	}
	logger.RootLogger.Error(logger.WithContent("[ID: " + id + "] 获取直播间信息失败"))
	return nil
}

// 解析返回包
func (room *Client) execPacket(data []byte) {
	header := PacketHead{
		PacketLen: binary.BigEndian.Uint32(data[0:4]),
		HeaderLen: binary.BigEndian.Uint16(data[4:6]),
		MsgType:   binary.BigEndian.Uint16(data[6:8]),
		Operation: binary.BigEndian.Uint32(data[8:12]),
		Sequence:  binary.BigEndian.Uint32(data[12:16]),
	}
	switch header.Operation {
	case WS_OP_CONNECT_SUCCESS: // 连接成功
		logger.RootLogger.System(logger.WithContent("[ID: " + room.ID + "] 进入直播间成功"))
		go room.sendHeartbeat() // 发送心跳包
		//todo: 添加状态统计
	case WS_OP_MESSAGE: // 弹幕
		var payload []byte
		if header.MsgType == WS_BODY_PROTOCOL_VERSION_BROTLI { // BROTLI 压缩
			b := brotli.NewReader(bytes.NewReader(data[16:]))
			if _payload, err := io.ReadAll(b); err != nil { // JSON 流式解析
				logger.RootLogger.Error(logger.WithContent("[ID: "+room.ID+"] 解压弹幕数据失败(BROTLI)", err, data[16:]))
				return
			} else {
				payload = _payload
			}
		} else if header.MsgType == WS_BODY_PROTOCOL_VERSION_NORMAL {
			go room.execCMD(data[16:])
			return
		} else {
			logger.RootLogger.Error(logger.WithContent("[ID: "+room.ID+"] 解压弹幕数据失败 类型错误", header.MsgType, data[16:]))
			return
		}
		var offset = 0
		for offset < len(payload) {
			jsonLen := int(binary.BigEndian.Uint32(payload[offset : offset+4])) // 下一 JSON 包长度
			// fmt.Println("json 正文:", string(payload[offset+16:offset+jsonLen]))
			go room.execCMD(payload[offset+16 : offset+jsonLen])
			offset += jsonLen
		}
		//logger.RootLogger.Debug(logger.WithContent("[ID: " + room.ID + "] 直播间收到弹幕"))
	case WS_OP_HEARTBEAT_REPLY: // 心跳包 PONG
		// logger.RootLogger.Debug(logger.WithContent("[ID: " + room.ID + "] 收到心跳包返回"))
		popularity := binary.BigEndian.Uint16(data[18:])
		logger.RootLogger.Common(logger.WithContent("[ID: "+room.ID+"] 直播间当前人气:", popularity))
	}
}

// 发送心跳包
func (room *Client) sendHeartbeat() {
	for {
		err := room.Conn.WriteMessage(websocket.BinaryMessage, packetGen([]byte{},
			WS_OP_HEARTBEAT, WS_HEADER_DEFAULT_SEQUENCE))
		if err != nil {
			logger.RootLogger.Error(logger.WithContent("[ID: "+room.ID+"] 直播间心跳包维系失败", err))
			time.Sleep(5 * time.Second)
			continue
		}
		time.Sleep(20 * time.Second)
	}

}

// 处理数据包
func (room *Client) execCMD(data []byte) {
	parse := gjson.ParseBytes(data)
	cmd := parse.Get("cmd").String()

	switch cmd {
	case WS_CMD_WATCHED_CHANGE: // 【N人观看】数据包
		logger.RootLogger.Notice(logger.WithContent(fmt.Sprintf(
			"[Room: %s] 直播间 %d 人观看",
			room.ID, parse.Get("data.num").Int()),
		))
	case WS_CMD_LIKE_INFO_V3_CLICK: // 【点赞】数据包
		var packet struct {
			Cmd  string    `json:"cmd"`
			Data likeClick `json:"data"`
		}
		err := json.Unmarshal(data, &packet)
		if err != nil {
			logger.RootLogger.Warning(logger.WithContent("[ID: "+room.ID+"] 解析点赞数据包失败", err))
			return
		}
		logger.RootLogger.Notice(logger.WithContent(fmt.Sprintf(
			"[Room: %s] %s(uid: %d)为直播间点赞",
			room.ID, packet.Data.Uname, packet.Data.Uid),
		))
	case WS_CMD_LIKE_INFO_V3_UPDATE:
		//fmt.Println(string(data))
		logger.RootLogger.Notice(logger.WithContent(fmt.Sprintf(
			"[Room: %s] 直播间 %d 人点赞",
			room.ID, parse.Get("data.click_count").Int()),
		))
	case WS_CMD_ONLINE_RANK_V2: // 【在线排行榜】数据包
		var packet struct {
			Cmd  string `json:"cmd"`
			Data struct {
				List     []onlineRank `json:"list"`
				RankType string       `json:"rank_type"`
			} `json:"data"`
		}
		if err := json.Unmarshal(data, &packet); err != nil {
			logger.RootLogger.Warning(logger.WithContent("[ID: "+room.ID+"] 解析高能榜数据包失败", err))
			return
		}
		var listStr string
		for _, v := range packet.Data.List {
			var count int // 0
			for _, c := range v.Uname {
				if unicode.Is(unicode.Han, c) {
					count++
				}
			}
			switch count {
			case 0, 1, 2:
				listStr += fmt.Sprintf("\n#%2d (uid: % 10d)%-40s\t-- %5s", v.Rank, v.Uid, v.Uname, v.Score)
			case 3, 4, 5, 6, 7, 8, 9, 10:
				listStr += fmt.Sprintf("\n#%2d (uid: % 10d)%-32s\t-- %5s", v.Rank, v.Uid, v.Uname, v.Score)
			case 11, 12, 13, 14, 15, 16:
				listStr += fmt.Sprintf("\n#%2d (uid: % 10d)%-24s\t-- %5s", v.Rank, v.Uid, v.Uname, v.Score)
			}
		}
		logger.RootLogger.Common(logger.WithContent(fmt.Sprintf(
			"[Room: %s] 直播间高能榜刷新: %s",
			room.ID, listStr),
		))
		// logger.RootLogger.Debug(logger.WithContent("[Room: "+room.ID+"] 数据包命令:", cmd, string(data)))
	case WS_CMD_ONLINE_RANK_TOP3:
		// 【在线排行榜】数据包
		// {"cmd":"ONLINE_RANK_TOP3","data":{"dmscore":112,"list":[{"msg":"恭喜 \u003c%Mr-大先生%\u003e 成为高能用户","rank":2}]}}
		var packet struct {
			Cmd  string         `json:"cmd"`
			Data onlineRankTop3 `json:"data"`
		}
		if err := json.Unmarshal(data, &packet); err != nil {
			logger.RootLogger.Warning(logger.WithContent("[ID: "+room.ID+"] 解析高能榜数据包失败", err))
			return
		}
		var listStr string
		for _, v := range packet.Data.List {
			listStr += fmt.Sprintf("%s\t[榜%d]", v.Msg, v.Rank)
		}
		logger.RootLogger.Common(logger.WithContent(fmt.Sprintf(
			"[Room: %s] %s",
			room.ID, listStr),
		))
		// logger.RootLogger.System(logger.WithContent("[Room: "+room.ID+"] 数据包命令:", cmd))
		fmt.Println(string(data))
	case WS_CMD_INTERACT_WORD:
		// 进入直播间
		return
		var packet struct {
			Cmd  string       `json:"cmd"`
			Data interactRoom `json:"data"`
		}
		if err := json.Unmarshal(data, &packet); err != nil {
			logger.RootLogger.Warning(logger.WithContent("[ID: "+room.ID+"] 解析 [进入直播间] 数据包失败", err))
			return
		}
		logger.RootLogger.Debug(logger.WithContent(fmt.Sprintf(
			"[Room: %s] %s(uid: %d)进入直播间",
			room.ID, packet.Data.Uname, packet.Data.Uid),
		))
	case WS_CMD_DANMU_MSG:
		// 弹幕消息
		var packet struct {
			Cmd  string `json:"cmd"`
			Info []any  `json:"info"`
		}
		if err := json.Unmarshal(data, &packet); err != nil {
			logger.RootLogger.Warning(logger.WithContent("[ID: "+room.ID+"] 解析 [弹幕消息] 数据包失败", err))
			return
		}
		// todo: len check
		msg := packet.Info[1].(string)
		uid, uname := packet.Info[2].([]any)[0].(float64), packet.Info[2].([]any)[1].(string)
		logger.RootLogger.Notice(logger.WithContent(fmt.Sprintf(
			"[Room: %s] %s(uid: %d)发送弹幕: %s",
			room.ID, uname, int(uid), msg),
		))
	case WS_CMD_ROOM_BLOCK_MSG:
		// 禁言事件
		var packet muteEvent
		if err := json.Unmarshal(data, &packet); err != nil {
			logger.RootLogger.Warning(logger.WithContent("[ID: "+room.ID+"] 解析 [禁言事件] 数据包失败", err))
			return
		}
		logger.RootLogger.Warning(logger.WithContent(fmt.Sprintf(
			"[Room: %s] %s(uid: %s)被禁言",
			room.ID, packet.Uname, packet.Uid),
		))
	case WS_CMD_UNUSED_STOP_LIVE_ROOM_LIST:
		break
	case WS_CMD_NOTICE_MSG, WS_CMD_WIDGET_BANNER: // 垃圾广告不看
		break
	case WS_CMD_ENTRY_EFFECT:
		// 舰长进入直播间
		var packet struct {
			Cmd  string       `json:"cmd"`
			Data captainEntry `json:"data"`
		}
		if err := json.Unmarshal(data, &packet); err != nil {
			logger.RootLogger.Warning(logger.WithContent("[ID: "+room.ID+"] 解析 [舰长进入直播间] 数据包失败", err))
			return
		}
		logger.RootLogger.Notice(logger.WithContent(fmt.Sprintf(
			"[Room: %s] %s",
			room.ID, packet.Data.CopyWriting),
		))
	case WS_CMD_UNUSED_ONLINE_RANK_COUNT:
		logger.RootLogger.Common(logger.WithContent(fmt.Sprintf(
			"[Room: %s] 直播间高能榜人数总计：%d",
			room.ID, gjson.ParseBytes(data).Get("data.count").Int()),
		))
	case WS_CMD_ROOM_REAL_TIME_MESSAGE_UPDATE:
		// 直播间 人气值变化
		logger.RootLogger.Common(logger.WithContent(fmt.Sprintf(
			"[Room: %s] 直播间 fans 数：%d\tfans_club 数：%d",
			room.ID, gjson.ParseBytes(data).Get("data.fans").Int(), gjson.ParseBytes(data).Get("data.fans_club").Int()),
		))
	case WS_CMD_SEND_GIFT:
		var packet struct {
			Cmd  string   `json:"cmd"`
			Data sendGift `json:"data"`
		}
		if err := json.Unmarshal(data, &packet); err != nil {
			logger.RootLogger.Warning(logger.WithContent("[ID: "+room.ID+"] 解析 [赠送礼物] 数据包失败", err))
			return
		}
		logger.RootLogger.Help(logger.WithContent(fmt.Sprintf(
			"[Room: %s] %s(uid: %d)%s了 %d 个 [%s]",
			room.ID, packet.Data.Uname, packet.Data.Uid, packet.Data.Action, packet.Data.Num, packet.Data.GiftName),
		))
	// 礼物消息
	case WS_CMD_ANCHOR_LOT_START, WS_CMD_DANMU_AGGREGATION:
		// 天选时刻 - 抽奖弹幕		无视
	default:
		logger.RootLogger.Debug(logger.WithContent("[Room: "+room.ID+"] 数据包命令:", cmd, string(data)))
	}
}

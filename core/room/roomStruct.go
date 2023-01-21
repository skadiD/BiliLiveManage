package room

import "github.com/gorilla/websocket"

// 直播间 结构体
type Client struct {
	ID       string          // 房间ID
	Info     *Init           // 初始化信息
	Auth     *ConnAuth       // 授权认证信息
	Conn     *websocket.Conn // 持久链接
	SendChan chan []byte     // 发送数据管道
}
type Init struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Group            string  `json:"group"`
		BusinessId       int     `json:"business_id"`
		RefreshRowFactor float64 `json:"refresh_row_factor"`
		RefreshRate      int     `json:"refresh_rate"`
		MaxDelay         int     `json:"max_delay"`
		Token            string  `json:"token"`
		HostList         []struct {
			Host    string `json:"host"`
			Port    int    `json:"port"`
			WssPort int    `json:"wss_port"`
			WsPort  int    `json:"ws_port"`
		} `json:"host_list"`
	} `json:"data"`
}

type ConnAuth struct {
	UID      uint8  `json:"uid"`
	RoomID   uint32 `json:"roomid"`
	ProtoVer uint8  `json:"protover"`
	Platform string `json:"platform"`
	Type     uint8  `json:"type"`
	Key      string `json:"key"`
}

type RespPacket struct {
	Header *PacketHead
	Data   []byte
}
type PacketHead struct {
	PacketLen uint32
	HeaderLen uint16 // 固定值 16
	MsgType   uint16 // 1:原始JSON	2:zlib	3:brotli
	Operation uint32 // 操作码 5目测是弹幕
	Sequence  uint32 // 发送为 1	返回为 0
}

package room

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/goccy/go-json"
	"net/url"
	"strconv"
)

// 预处理弹幕服务器地址
func (data *Init) prefixHost() []url.URL {
	hostList := make([]url.URL, 0, 3)

	for _, host := range data.Data.HostList {
		u := url.URL{Scheme: "wss", Host: fmt.Sprintf("%s:%d", host.Host, host.WssPort), Path: "/sub"}
		hostList = append(hostList, u)
	}
	return hostList
}

// 发送首次握手包
func (room *Client) FirstHandlePack() []byte {
	roomID, _ := strconv.Atoi(room.ID)
	authData := ConnAuth{
		UID:      0,
		RoomID:   uint32(roomID),
		ProtoVer: WS_BODY_PROTOCOL_VERSION_BROTLI,
		Platform: "web", // 似乎是固定值
		Type:     2,     // 看起来也是固定值
		Key:      room.Info.Data.Token,
	} // 创建认证包
	authBytes, _ := json.Marshal(authData)
	headBytes := packetGen(authBytes, WS_OP_USER_AUTHENTICATION, WS_HEADER_DEFAULT_OPERATION)

	//logger.RootLogger.Debug(logger.WithContent("Auth包长度：", len(headBytes)))
	//hexStr := fmt.Sprintf("%08x", headBytes)
	//re, _ := regexp.Compile(".{4}")
	//fmt.Println(hexStr, re.FindAllString(hexStr, -1))
	//logger.RootLogger.Debug(logger.WithContent("Auth包内容：", re.FindAllString(hexStr, -1)))
	//fmt.Println(string(headBytes))
	return headBytes
}

// 数据包构造
func packetGen(data []byte, Operation, Sequence int) []byte {
	header := new(bytes.Buffer)

	binary.Write(header, binary.BigEndian, uint32(len(data)+WS_PACKAGE_HEADER_TOTAL_LENGTH))
	binary.Write(header, binary.BigEndian, uint16(WS_PACKAGE_HEADER_TOTAL_LENGTH))
	binary.Write(header, binary.BigEndian, uint16(WS_HEADER_DEFAULT_VERSION))
	binary.Write(header, binary.BigEndian, uint32(Operation))
	binary.Write(header, binary.BigEndian, uint32(Sequence))

	socketData := append(header.Bytes(), data...)
	return socketData
}

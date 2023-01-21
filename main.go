package main

import (
	"github.com/fexli/logger"
	"github.com/skadiD/BiliLiveManage/core/room"
)

func main() {
	logger.InitGlobLog("running.log", "BLM")
	room.Create("222521")
	//test()
}

//func main() {
//	logger.InitGlobLog("test.log")
//	for i := 1; i < 10000000; i++ {
//		go func(i int) {
//			//logger.RootLogger.Common(logger.WithContent("test content", i))
//			//if i%1000 == 0 {
//			//	logger.RootLogger.ClearLogInfo()
//			//}
//			fmt.Println("test ", i)
//		}(i)
//	}
//
//}

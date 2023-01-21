package room

import (
	"bytes"
	"encoding/binary"
	"github.com/fexli/logger"
)

type Int interface {
	~int64 | ~int32 | ~int16 | ~int
}

func itob[T Int](num T) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		logger.RootLogger.Warning(logger.WithContent(err))
	}
	return buffer.Bytes()
}

func decodeBrotli() []byte {
	return nil
}

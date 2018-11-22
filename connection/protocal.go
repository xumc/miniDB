/*
   protocol is a custom defined protocal to solve stick package problem.
*/
package connection

import "github.com/xumc/miniDB/utils"

const (
	ConstHeader         = "====protocol delimiter===="
	ConstHeaderLength   = 26
	ConstSaveDataLength = 4
)

// Packet implemtns pack logic
func Packet(message []byte) []byte {
	return append(append([]byte(ConstHeader), utils.IntToBytes(len(message))...), message...)
}

// Unpack implements unpack logic
func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstHeaderLength+ConstSaveDataLength {
			break
		}
		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := utils.BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstSaveDataLength])
			if length < i+ConstHeaderLength+ConstSaveDataLength+messageLength {
				break
			}
			data := buffer[i+ConstHeaderLength+ConstSaveDataLength : i+ConstHeaderLength+ConstSaveDataLength+messageLength]
			readerChannel <- data

			i += ConstHeaderLength + ConstSaveDataLength + messageLength - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

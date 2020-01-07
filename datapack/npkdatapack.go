package datapack

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"locust/log"
	"locust/message"
	"locust/utils"
)

/**
 * @Author: stydm
 * @Date: 2019-10-30 19:12
 */

// 定义常量 7D02  <=> 7E,  7D01  <=> 7D
const (
	npkBegin byte = 0xaa
	npkBeginChange byte = 0xbb
	npkEnd byte = 0xcc
	npkEndChange byte = 0xdd
	npk01 byte = 0x01
	npk02 byte = 0x02
)

//封包拆包类实例，暂时不需要成员
type NPKDataPack struct {}

//封包拆包实例初始化方法
func NewDataPack() *NPKDataPack {
	return &NPKDataPack{}
}

// 解码
func decodeTrans(d []byte)([]byte , bool){
	bytes := make([]byte,len(d),cap(d))
	copy(bytes,d)
	for i:= len(d)-1 ;i>=0 ;i-- {
		if d[i] == npkBeginChange {
			if i==len(d)-1 {
				return nil,false
			}
			if d[i+1]==npk01 {
				if (i+1)<len(d) {
					bytes = append(bytes[:i+1],bytes[i+2:]...)
				}else{
					bytes = append(bytes[:i+1])
				}
				bytes[i] = npkBegin
			}else if d[i+1]==npk02 {
				if (i+1)<len(d) {
					bytes = append(bytes[:i+1],bytes[i+2:]...)
				}else{
					bytes = append(bytes[:i+1])
				}
				bytes[i] = npkBeginChange
			}else{
				return nil,false
			}
		}
		if d[i] == npkEndChange {
			if i==len(d)-1 {
				return nil,false
			}
			if d[i+1]==npk01 {
				if (i+1)<len(d) {
					bytes = append(bytes[:i+1],bytes[i+2:]...)
				}else{
					bytes = append(bytes[:i+1])
				}
				bytes[i] = npkEnd
			}else if d[i+1]==npk02 {
				if (i+1)<len(d) {
					bytes = append(bytes[:i+1],bytes[i+2:]...)
				}else{
					bytes = append(bytes[:i+1])
				}
				bytes[i] = npkEndChange
			}else{
				return nil,false
			}
		}
	}

	return bytes,true
}


//获取包头长度方法
func(dp *NPKDataPack) GetHeadLen() uint32 {
	//Id uint32(4字节) +  DataLen uint32(4字节)
	return 8
}
//封包方法(压缩数据)
func(dp *NPKDataPack) Pack(msg message.IMessage)([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil ,err
	}

	return dataBuff.Bytes(), nil
}
//拆包方法(解压数据)
func(dp *NPKDataPack) Unpack(binaryData []byte)([]message.IMessage, int8) {
	returnMsg := make([]message.IMessage, 0)

	indexBegins,isOk := utils.SliceSearch(binaryData,npkBegin)
	if !isOk{
		return nil,-1
	}
	if len(indexBegins)<1 {
		return nil,-1
	}
	indexEnds,isOk := utils.SliceSearch(binaryData,npkEnd)

	if len(indexEnds)<1 {
		msg := &message.NPKMessage{}
		msg.Data = binaryData[indexBegins[len(indexBegins)-1]:]
		returnMsg = append(returnMsg, msg)
		return returnMsg,1
	}
	for i:=0;i<len(indexBegins);i++ {
		for j:=0;j<len(indexEnds);j++ {
			if indexBegins[i] < indexEnds[j] {
				tempData,b := decodeTrans(binaryData[indexBegins[i]+1:indexEnds[j]])
				if !b || len(tempData)<10 {
					break
				}
				if utils.CheckSumModBus(tempData[:len(tempData)-2]) == utils.BytesToUint16(tempData[len(tempData)-2:]) {
					msg := &message.NPKMessage{}
					msg.Data = tempData[:len(tempData)-2]
					msg.Id= uint32(tempData[0])
					if msg.Id!=255 {
						if (len(msg.Data)-13)==int(msg.Data[11]) {
							returnMsg = append(returnMsg, msg)
						}else{
							log.LogFile.WithFields(logrus.Fields{
								"operation": "长度不对",
								"data":   hex.EncodeToString(tempData),
							}).Info("A group of walrus emerges from the ocean")
						}
					}
					msg = nil
				}else{
					log.LogFile.WithFields(logrus.Fields{
						"operation": "校验不通过",
						"data":   hex.EncodeToString(tempData),
					}).Info("A group of walrus emerges from the ocean")
				}
				break
			}
		}
	}


	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return returnMsg, 0
}
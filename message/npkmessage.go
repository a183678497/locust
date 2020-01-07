package message

/**
 * @Author: stydm
 * @Date: 2019-10-30 19:09
 */

type NPKMessage struct {
	Id      uint32 //消息的ID
	Data    []byte //消息的内容
}

//获取消息ID
func (msg *NPKMessage) GetMsgId() uint32 {
	return msg.Id
}

//获取消息内容
func (msg *NPKMessage) GetData() []byte {
	return msg.Data
}


//设置消息ID
func (msg *NPKMessage) SetMsgId(msgId uint32) {
	msg.Id = msgId
}

//设置消息内容
func (msg *NPKMessage) SetData(data []byte) {
	msg.Data = data
}


//创建一个Message消息包
func NewMsgPackage(id uint32, data []byte) *NPKMessage {
	return &NPKMessage{
		Id:     id,
		Data:   data,
	}
}
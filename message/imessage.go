package message

/**
 * @Author: stydm
 * @Date: 2019-10-30 19:03
 */

/*
   将请求的一个消息封装到message中，定义抽象层接口
*/
type IMessage interface {
	GetMsgId() uint32    //获取消息ID
	GetData() []byte     //获取消息内容

	SetMsgId(uint32)     //设计消息ID
	SetData([]byte)      //设计消息内容
}
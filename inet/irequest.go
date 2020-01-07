package inet

/**
 * @Author: stydm
 * @Date: 2019-10-30 16:05
 */

/*
   IRequest 接口：
   实际上是把客户端请求的链接信息 和 请求的数据 包装到了 Request里
*/
type IRequest interface{
	GetConnection() IConnection    //获取请求连接信息
	GetData() []byte            //获取请求消息的数据
	GetMsgID() uint32           //获取请求的消息ID
}
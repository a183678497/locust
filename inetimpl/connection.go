package inetimpl

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"locust/common"
	"locust/datapack"
	"locust/inet"
	"locust/log"
	"locust/message"
	"net"
	"sync"
)

/**
 * @Author: stydm
 * @Date: 2019-10-30 14:43
 */

type Connection struct {
	//当前Conn属于哪个Server
	TcpServer    inet.IServer       //当前conn属于哪个server，在conn初始化的时候添加即可
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool
	//消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler inet.IMsgHandle
	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan        chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	// ================================
	//链接属性
	property     map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
	// ================================
}


//创建连接的方法
func NewConntion(server inet.IServer, conn *net.TCPConn, connID uint32, msgHandler inet.IMsgHandle) *Connection{
	//初始化Conn属性
	c := &Connection{
		TcpServer:server,
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		MsgHandler: msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, common.ConfigObject.MaxMsgChanLen),
		property:     make(map[string]interface{}), //对链接属性map初始化
	}

	//将新创建的Conn添加到链接管理中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

/*
   写消息Goroutine， 用户将数据发送给客户端
*/
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
			//针对有缓冲channel需要些的数据处理
		case data, ok:= <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				break
				fmt.Println("msgBuffChan is Closed")
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is  running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	//bufferData := make([]byte, 0)

	for  {
		dp := datapack.NewDataPack()
		//读取我们最大的数据到buf中
		buf := make([]byte, common.ConfigObject.MaxPacketSize)
		dataLen,err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			c.ExitBuffChan <- true
			continue
		}

		log.LogFile.WithFields(logrus.Fields{
			"operation": "接收数据",
			"data":   hex.EncodeToString(buf[:dataLen]),
		}).Info("A group of walrus emerges from the ocean")


		//bufferData = append(bufferData[:len(bufferData)],buf[:dataLen]...)
		//拆包，得到msgid 和 datalen 放在msg中
		msg , code := dp.Unpack(buf)
		if code == 0{
			for _,val := range msg {
				//得到当前客户端请求的Request数据
				req := Request{
					conn:c,
					msg:val,
				}

				if common.ConfigObject.WorkerPoolSize > 0 {
					//已经启动工作池机制，将消息交给Worker处理
					c.MsgHandler.SendMsgToTaskQueue(&req)
				} else {
					//从绑定好的消息和对应的处理方法中执行对应的Handle方法
					go c.MsgHandler.DoMsgHandler(&req)
				}
			}

		}

	}
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data封包，并且发送
	dp := datapack.NewDataPack()
	msg, err := dp.Pack(message.NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return  errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgChan <- msg   //将之前直接回写给conn.Write的方法 改为 发送给Channel 供Writer读取

	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	//将data封包，并且发送
	dp := datapack.NewDataPack()
	msg, err := dp.Pack(message.NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return  errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgBuffChan <- msg

	return nil
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {

	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	//==================
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
	//==================

	for {
		select {
		case <- c.ExitBuffChan:
			//得到退出消息，不再阻塞
			return
		}
	}
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//==================
	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)
	//==================

	// 关闭socket链接
	c.Conn.Close()
	//关闭Writer Goroutine
	c.ExitBuffChan <- true

	//将链接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)  //删除conn从ConnManager中

	//关闭该链接全部管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接ID
func (c *Connection) GetConnID() uint32{
	return c.ConnID
}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok  {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
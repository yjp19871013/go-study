package tcp_client

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	TcpStateReconnecting = iota
	TcpStateConnected
	TcpStateDisconnect
)

type TCPClient struct {
	ServerAddress string
	ServerPort    uint
	TimeoutSec    time.Duration

	conn              net.Conn
	reconnectChan     chan byte
	stopReconnectChan chan byte
	listenChan        chan int

	sync.Mutex
}

func CreateTCPClient(serverAddress string, serverPort uint, timeout time.Duration) *TCPClient {
	return &TCPClient{
		ServerAddress: serverAddress,
		ServerPort:    serverPort,
		TimeoutSec:    timeout,
	}
}

func DestroyTCPClient(client *TCPClient) {
	if client == nil {
		return
	}

	client.TimeoutSec = 0
	client.ServerPort = 0
	client.ServerAddress = ""

	client = nil
}

func (c *TCPClient) StartConnect(listenChan chan int) {
	c.Lock()
	defer c.Unlock()

	if c.isStartConnectCalled() {
		fmt.Println("StartConnect has called")
		return
	}

	// 创建需要使用的通道
	c.reconnectChan = make(chan byte, 1)
	c.stopReconnectChan = make(chan byte, 1)
	c.listenChan = listenChan

	// 启动重连协程
	go c.reconnect()

	// 触发一次重连
	c.triggerReconnect()
}

func (c *TCPClient) StopConnect() {
	c.Lock()
	defer c.Unlock()

	if !c.isStartConnectCalled() {
		fmt.Println("StartConnect has not called")
		return
	}

	c.triggerStopReconnect()
}

func (c *TCPClient) Read(readBuffer []byte) (int, error) {
	c.Lock()
	defer c.Unlock()

	if c.reconnectChan == nil {
		return 0, fmt.Errorf("You should call StartConnect first.")
	}

	if !c.isConnected() {
		return 0, fmt.Errorf("reconnect")
	}

	_ = c.conn.SetDeadline(time.Now().Add(time.Second * c.TimeoutSec))
	n, err := c.conn.Read(readBuffer)
	if err != nil {
		netError, ok := err.(net.Error)
		if ok && netError.Timeout() {
			return 0, err
		}

		c.triggerReconnect()
		return 0, err
	}

	return n, nil
}

func (c *TCPClient) Write(data []byte) error {
	c.Lock()
	defer c.Unlock()

	if c.reconnectChan == nil {
		return fmt.Errorf("You should call StartConnect First.")
	}

	if !c.isConnected() {
		return fmt.Errorf("reconnect")
	}

	_ = c.conn.SetDeadline(time.Now().Add(time.Second * c.TimeoutSec))
	_, err := c.conn.Write(data)
	if err != nil {
		netError, ok := err.(net.Error)
		if ok && netError.Timeout() {
			return err
		}

		c.triggerReconnect()
		return err
	}

	return nil
}

func (c *TCPClient) reconnect() {
	for {
		select {
		case <-c.reconnectChan:
			c.doReconnect()
		case <-c.stopReconnectChan:
			c.stopReconnect()
		}
	}
}

func (c *TCPClient) doReconnect() {
	c.Lock()
	defer c.Unlock()

	c.notify(TcpStateReconnecting)

	c.closeConn()
	err := c.createConn()
	if err != nil {
		fmt.Println(err)
		c.triggerReconnect()
		return
	}

	fmt.Println("重连成功")
	c.notify(TcpStateConnected)
}

func (c *TCPClient) stopReconnect() {
	c.Lock()
	defer c.Unlock()

	c.disconnect()
}

func (c *TCPClient) disconnect() {
	c.closeConn()
	c.reconnectChan = nil
	c.stopReconnectChan = nil

	c.notify(TcpStateDisconnect)
}

func (c *TCPClient) createConn() error {
	address := c.ServerAddress + ":" + strconv.Itoa(int(c.ServerPort))

	conn, err := net.DialTimeout("tcp", address, time.Second*c.TimeoutSec)
	if err != nil {
		return err
	}
	c.conn = conn

	return nil
}

func (c *TCPClient) closeConn() {
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
}

func (c *TCPClient) triggerReconnect() {
	if len(c.reconnectChan) < cap(c.reconnectChan) {
		c.reconnectChan <- 1
	}
}

func (c *TCPClient) triggerStopReconnect() {
	if len(c.stopReconnectChan) < cap(c.stopReconnectChan) {
		c.stopReconnectChan <- 1
	}
}

func (c *TCPClient) isConnected() bool {
	return c.conn != nil
}

func (c *TCPClient) isStartConnectCalled() bool {
	return c.reconnectChan != nil && c.stopReconnectChan != nil
}

func (c *TCPClient) notify(state int) {
	if c.listenChan == nil {
		return
	}
	go func(state int) {
		c.listenChan <- state
	}(state)
}

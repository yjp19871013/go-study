package tcp_server

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)


// accept回调函数
type OnAccept func(conn *ClientConn, readChan chan []byte, writeChan chan []byte)

// 客户端连接
type ClientConn struct {
	conn net.Conn
}

type TCPServer struct {
	// 服务器IP
	serverIP string

	// 服务器端口号
	serverPort int

	// 服务器地址
	address string

	// 接收缓冲区大小，单位字节
	recvBufferLen int

	// listen后返回的监听
	listener net.Listener

	// 保存当前的客户端连接
	clientConns []*ClientConn
	clientConnsMutex sync.Mutex

	// 保存accept回调
	onAccept OnAccept
}

func CreateTCPServer(serverIP string, serverPort int, recvBufferLen int) *TCPServer {
	server := new(TCPServer)

	server.recvBufferLen = recvBufferLen

	server.serverIP = serverIP
	server.serverPort = serverPort
	server.address = serverIP + ":" + strconv.Itoa(serverPort)

	return server
}

func DestroyServer(server *TCPServer) {
	if server == nil {
		fmt.Println("DestroyServer server nil")
		return
	}

	if server.listener != nil {
		fmt.Println("DestroyServer before server stop")
		return
	}

	server.address = ""
	server.serverPort = 0
	server.serverIP = ""

	server.recvBufferLen = 0

	server = nil
}

func (server *TCPServer) Start(onAccept OnAccept) error {
	if server.listener != nil {
		fmt.Println("Start listener nil")
		return nil
	}

	server.onAccept = onAccept

	listener, err := net.Listen("tcp", server.address)
	if err != nil {
		fmt.Println("Start", err)
		return err
	}

	server.listener = listener

	server.clientConnsMutex.Lock()
	server.clientConns = make([]*ClientConn, 0)
	server.clientConnsMutex.Unlock()

	go server.accept()

	return nil
}

func (server *TCPServer) Stop() {
	if server.listener == nil {
		return
	}

	_ = server.listener.Close()
	server.listener = nil

	server.clientConnsMutex.Lock()

	// 关闭所有未关闭的连接
	for _, clientConn := range server.clientConns {
		_ = clientConn.conn.Close()
		clientConn.conn = nil
	}

	server.clientConns = nil

	server.clientConnsMutex.Unlock()

	server.onAccept = nil
}

func (server *TCPServer) CloseClientConn(conn *ClientConn) {
	server.clientConnsMutex.Lock()
	defer server.clientConnsMutex.Unlock()

	if server.clientConns == nil {
		return
	}

	for i, clientConn := range server.clientConns {
		if clientConn == conn {
			_ = clientConn.conn.Close()
			clientConn.conn = nil
		}

		server.clientConns = append(server.clientConns[0:i], server.clientConns[i+1:]...)

		break
	}
}

func (server *TCPServer) accept() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		server.clientConnsMutex.Lock()
		clientConn := new(ClientConn)
		clientConn.conn = conn
		server.clientConns = append(server.clientConns, clientConn)
		server.clientConnsMutex.Unlock()

		go server.handleConn(clientConn)
	}
}

func (server *TCPServer) handleConn(conn *ClientConn) {
	readChan := make(chan []byte)
	writeChan := make(chan []byte)

	if server.onAccept != nil {
		go server.onAccept(conn, readChan, writeChan)
	}

	go server.readConn(conn, readChan)
	go server.writeConn(conn, writeChan)
}

func (server *TCPServer) readConn(clientConn *ClientConn, readChan chan<- []byte) {
	for {
		data := make([]byte, server.recvBufferLen)
		_, err := clientConn.conn.Read(data)
		if err != nil {
			fmt.Println(err)
			break
		}

		readChan <- data
	}

	server.CloseClientConn(clientConn)
}

func (server *TCPServer) writeConn(clientConn *ClientConn, writeChan <-chan []byte) {
	for {
		data := <-writeChan
		_, err := clientConn.conn.Write(data)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	server.CloseClientConn(clientConn)
}



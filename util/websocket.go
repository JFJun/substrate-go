package util

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Socket struct {
	Conn              *websocket.Conn
	WebsocketDialer   *websocket.Dialer
	Url               string
	ConnectionOptions ConnectionOptions
	RequestHeader     http.Header
	OnConnected       func(socket Socket)
	OnTextMessage     func(message string, socket Socket)
	OnBinaryMessage   func(data []byte, socket Socket)
	OnConnectError    func(err error, socket Socket)
	OnDisconnected    func(err error, socket Socket)
	OnPingReceived    func(data string, socket Socket)
	OnPongReceived    func(data string, socket Socket)
	IsConnected       bool
	sendMu            *sync.Mutex // Prevent "concurrent write to websocket connection"
	receiveMu         *sync.Mutex
	reConnectNum      int
}

type ConnectionOptions struct {
	UseCompression bool
	UseSSL         bool
	Proxy          func(*http.Request) (*url.URL, error)
	Subprotocols   []string
}

func NewWebsocket(url string) Socket {
	return Socket{
		Url:           url,
		RequestHeader: http.Header{},
		ConnectionOptions: ConnectionOptions{
			UseCompression: false,
			UseSSL:         false,
		},
		WebsocketDialer: &websocket.Dialer{},
		sendMu:          &sync.Mutex{},
		receiveMu:       &sync.Mutex{},
	}
}

func (socket *Socket) setConnectionOptions() {
	socket.WebsocketDialer.EnableCompression = socket.ConnectionOptions.UseCompression
	socket.WebsocketDialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: socket.ConnectionOptions.UseSSL}
	socket.WebsocketDialer.Proxy = socket.ConnectionOptions.Proxy
	socket.WebsocketDialer.Subprotocols = socket.ConnectionOptions.Subprotocols
}

func (socket *Socket) Connect() {
	var err error
	socket.setConnectionOptions()

	socket.Conn, _, err = socket.WebsocketDialer.Dial(socket.Url, socket.RequestHeader)

	if err != nil {
		log.Println("Error while connecting to server ", err)
		socket.IsConnected = false
		if socket.OnConnectError != nil {
			socket.OnConnectError(err, *socket)
		}
		return
	}

	log.Println("Connected to server")

	if socket.OnConnected != nil {
		socket.IsConnected = true
		socket.OnConnected(*socket)
	}

	defaultPingHandler := socket.Conn.PingHandler()
	socket.Conn.SetPingHandler(func(appData string) error {
		log.Println("Received PING from server")
		if socket.OnPingReceived != nil {
			socket.OnPingReceived(appData, *socket)
		}
		return defaultPingHandler(appData)
	})

	defaultPongHandler := socket.Conn.PongHandler()
	socket.Conn.SetPongHandler(func(appData string) error {
		log.Println("Received PONG from server")
		if socket.OnPongReceived != nil {
			socket.OnPongReceived(appData, *socket)
		}
		return defaultPongHandler(appData)
	})

	defaultCloseHandler := socket.Conn.CloseHandler()
	socket.Conn.SetCloseHandler(func(code int, text string) error {
		result := defaultCloseHandler(code, text)
		log.Println("Disconnected from server ", result)
		if socket.OnDisconnected != nil {
			socket.IsConnected = false
			socket.OnDisconnected(errors.New(text), *socket)
		}
		return result
	})

	go func() {
		for {
			socket.receiveMu.Lock()
			messageType, message, err := socket.Conn.ReadMessage()
			socket.receiveMu.Unlock()
			if err != nil {
				log.Println("read:", err)
			}
			switch messageType {
			case websocket.TextMessage:
				if socket.OnTextMessage != nil {
					socket.OnTextMessage(string(message), *socket)
				}
			case websocket.BinaryMessage:
				if socket.OnBinaryMessage != nil {
					socket.OnBinaryMessage(message, *socket)
				}
			}
		}
	}()
}

func (socket *Socket) SendText(message string) {
	err := socket.send(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("write:", err)
		return
	}
	log.Println("发送数据： ", message)
}

func (socket *Socket) SendBinary(data []byte) {
	err := socket.send(websocket.BinaryMessage, data)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func (socket *Socket) send(messageType int, data []byte) error {
	socket.sendMu.Lock()
	err := socket.Conn.WriteMessage(messageType, data)
	socket.sendMu.Unlock()
	return err
}

func (socket *Socket) Close() {
	err := socket.send(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
	}
	socket.Conn.Close()
	if socket.OnDisconnected != nil {
		socket.IsConnected = false
		socket.OnDisconnected(err, *socket)
	}
}

type JsonRpcResponse struct {
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Id      int         `json:"id"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}

func (socket *Socket) SendRequest(method string, params []interface{}) ([]byte, error) {
	err := socket.reConnect()
	if err != nil {
		return nil, err
	}
	reqData := map[string]interface{}{
		"id":      time.Now().Unix(),
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}
	dd, _ := json.Marshal(reqData)
	socket.sendMu.Lock()
	err = socket.Conn.WriteMessage(websocket.BinaryMessage, dd)
	socket.sendMu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("ws send req data error,Err=%v", err)
	}
	socket.receiveMu.Lock()
	_, respData, err := socket.Conn.ReadMessage()
	socket.receiveMu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("ws resp data error,Err=%v", err)
	}
	if len(respData) == 0 {
		return nil, errors.New("ws resp data is null")
	}
	var resp JsonRpcResponse
	err = json.Unmarshal(respData, &resp)
	if err != nil {
		return nil, fmt.Errorf("ws json unmarshal resp data error,err=%v", err)
	}
	if resp.Error != nil {
		errData, _ := json.Marshal(resp.Error)
		return nil, fmt.Errorf("ws resp error ,err=%s", string(errData))
	}
	if resp.Result == nil {
		return nil, errors.New("ws resp result is null")
	}
	switch resp.Result.(type) {
	case string:
		return []byte(resp.Result.(string)), nil
	default:
		res, _ := json.Marshal(resp.Result)
		return res, nil
	}
}

func (socket *Socket) reConnect() error {
	if socket.Conn == nil {
		var err error
		socket.setConnectionOptions()

		socket.Conn, _, err = socket.WebsocketDialer.Dial(socket.Url, socket.RequestHeader)

		if err != nil {
			log.Printf("error while connecting to server,err=%v,reConnect num is %d ", err, socket.reConnectNum)
			socket.reConnectNum++
			socket.IsConnected = false
			if socket.OnConnectError != nil {
				socket.OnConnectError(err, *socket)
			}
			if socket.reConnectNum == 3 {
				return err
			}
			socket.reConnect()
		}
		socket.reConnectNum = 0
	}
	return nil
}

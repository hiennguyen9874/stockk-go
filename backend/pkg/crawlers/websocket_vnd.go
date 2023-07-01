package crawlers

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type vndWebsocketCrawlers struct {
	connection  *websocket.Conn
	isConnected bool
	cfg         *config.Config
	logger      logger.Logger
	mu          sync.Mutex
}

func NewVNDWebsocketCrawlers(cfg *config.Config, logger logger.Logger) WebsocketCrawlers {
	return &vndWebsocketCrawlers{isConnected: false, cfg: cfg, logger: logger}
}

func (wsc *vndWebsocketCrawlers) Connect() error {
	if wsc.isConnected {
		return nil
	}

	if wsc.connection != nil {
		wsc.isConnected = false
		wsc.connection.Close()
	}

	u := url.URL{
		Scheme: "wss",
		Host:   "price-api-free.vndirect.com.vn",
		Path:   "/realtime",
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	wsc.connection = c
	wsc.isConnected = true
	return nil
}

func (wsc *vndWebsocketCrawlers) Close() error {
	return wsc.connection.Close()
}

func (wsc *vndWebsocketCrawlers) WriteMessage(messages []string) error {
	for _, message := range messages {
		wsc.mu.Lock()
		err := wsc.connection.WriteMessage(websocket.TextMessage, []byte(message))
		wsc.mu.Unlock()
		if err != nil {
			wsc.isConnected = false
			return err
		}
	}
	return nil
}

func (wsc *vndWebsocketCrawlers) ReadMessage() <-chan Message {
	type MessageRaw struct {
		messageStr string
		messageErr *error
	}

	messageCh := make(chan MessageRaw)

	go func() {
		defer close(messageCh)

		for {
			_, message, err := wsc.connection.ReadMessage()

			if err != nil {
				wsc.isConnected = false
				messageCh <- MessageRaw{
					messageErr: &err,
				}
			} else {
				messageCh <- MessageRaw{
					messageStr: string(message),
				}
			}
		}
	}()

	resultCh := make(chan Message)

	go func() {
		defer close(resultCh)

		for message := range messageCh {
			if message.messageErr != nil {
				resultCh <- Message{
					MessageErr: message.messageErr,
				}
			} else {
				decodeArr, decodeType, err := vndDecodeWebsocketMessage(message.messageStr)
				if err != nil {
					resultCh <- Message{
						MessageErr: &err,
					}
				} else {
					resultCh <- Message{
						MessageDict: &decodeArr,
						MessageType: &decodeType,
					}
				}
			}
		}
	}()

	return resultCh
}

func vndDecodeWebsocketMessage(message string) (map[string]string, string, error) {
	messageSplit := strings.Split(message, "|")
	messageType := messageSplit[0]
	messageData := strings.Join(messageSplit[1:], "|")

	switch strings.ToUpper(messageType) {
	case "D":
		messageDecodeSplit := strings.Split(messageData, ":")
		messageDecodeType := messageDecodeSplit[0]
		messageDecodeData := strings.Join(messageDecodeSplit[1:], ":")
		messageArray := VNDDecodeMessage(messageDecodeData)

		messageDict, err := VNDMessageArrayToDict(messageDecodeType, messageArray)
		return messageDict, messageDecodeType, err
	default:
		return nil, "", fmt.Errorf("not implemented error: %v", messageType)
	}
}

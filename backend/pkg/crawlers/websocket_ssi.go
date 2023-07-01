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

type ssiWebsocketCrawlers struct {
	connection  *websocket.Conn
	isConnected bool
	cfg         *config.Config
	logger      logger.Logger
	mu          sync.Mutex
}

func NewSSIWebsocketCrawlers(cfg *config.Config, logger logger.Logger) WebsocketCrawlers {
	return &ssiWebsocketCrawlers{isConnected: false, cfg: cfg, logger: logger}
}

func (wsc *ssiWebsocketCrawlers) Connect() error {
	if wsc.isConnected {
		return nil
	}

	if wsc.connection != nil {
		wsc.isConnected = false
		wsc.connection.Close()
	}

	u := url.URL{
		Scheme: "wss",
		Host:   "pricestream-iboard.ssi.com.vn",
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

func (wsc *ssiWebsocketCrawlers) Close() error {
	return wsc.connection.Close()
}

func (wsc *ssiWebsocketCrawlers) WriteMessage(messages []string) error {
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

func (wsc *ssiWebsocketCrawlers) ReadMessage() <-chan Message {
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
				continue
			}

			messageDict, messageType, err := ssiDecodeWebsocketMessage(message.messageStr)
			if err != nil {
				resultCh <- Message{
					MessageErr: &err,
				}
				continue
			}

			resultCh <- Message{
				MessageDict: &messageDict,
				MessageType: &messageType,
			}
		}
	}()

	return resultCh
}

func ssiDecodeWebsocketMessage(message string) (map[string]string, string, error) {
	if strings.HasPrefix(message, "I#") {
		messageDict, err := zipToDict([]string{
			"indexID",
			"indexValue",
			"totalQtty",
			"totalValue",
			"advances",
			"declines",
			"nochanges",
			"ceiling",
			"floor",
			"allQty",
			"allValue",
			"time",
			"vol",
			"timeTDW",
			"totalQttyPrevTDW",
			"chartOpen",
			"chartHigh",
			"chartLow",
			"change",
			"changePercent",
			"stockChange",
		}, strings.Split(strings.Split(message, "#")[1], "|"))
		return messageDict, "I", err
	}

	if strings.HasPrefix(message, "S#") {
		messageDict, err := zipToDict([]string{
			"stockNo",
			"stockSymbol",
			"best1Bid",
			"best1BidVol",
			"best2Bid",
			"best2BidVol",
			"best3Bid",
			"best3BidVol",
			"best4Bid",
			"best4BidVol",
			"best5Bid",
			"best5BidVol",
			"best6Bid",
			"best6BidVol",
			"best7Bid",
			"best7BidVol",
			"best8Bid",
			"best8BidVol",
			"best9Bid",
			"best9BidVol",
			"best10Bid",
			"best10BidVol",
			"best1Offer",
			"best1OfferVol",
			"best2Offer",
			"best2OfferVol",
			"best3Offer",
			"best3OfferVol",
			"best4Offer",
			"best4OfferVol",
			"best5Offer",
			"best5OfferVol",
			"best6Offer",
			"best6OfferVol",
			"best7Offer",
			"best7OfferVol",
			"best8Offer",
			"best8OfferVol",
			"best9Offer",
			"best9OfferVol",
			"best10Offer",
			"best10OfferVol",
			"matchedPrice",
			"matchedVolume",
			"highest",
			"exchange",
			"lowest",
			"avgPrice",
			"buyForeignQtty",
			"buyForeignValue",
			"sellForeignQtty",
			"sellForeignValue",
			"priceChange",
			"priceChangePercent",
			"nmTotalTradedQty",
			"nmTotalTradedValue",
			"currentBidQty",
			"currentOfferQty",
			"openInterest",
			"ceiling",
			"floor",
			"refPrice",
			"maturityDate",
			"session",
			"caStatus",
			"remainForeignQtty",
			"stockType",
			"tradingStatus",
			"lastTradingDate",
			"coveredWarrantType",
			"underlyingSymbol",
			"exercisePrice",
			"listedShare",
			"issuerName",
			"issueDate",
			"openPrice",
			"prevMatchedPrice",
			"lastMatchedPrice",
			"basis",
			"issuer",
			"underlyingPrice",
			"breakEven",
			"breakEvenMarketPrice",
			"exerciseRatio",
			"underlyingRef",
			"underlyingCeiling",
			"underlyingFloor",
			"oddSession",
			"oddMatchedPrice",
			"oddMatchedVolume",
			"oddPriceChange",
			"oddPriceChangePercent",
		}, strings.Split(strings.Split(message, "#")[1], "|"))
		return messageDict, "S", err
	}

	if strings.HasPrefix(message, "M#") {
		messageDict, err := zipToDict([]string{
			"stockNo",
			"price",
			"totalVol",
			"buyUpVol",
			"sellDownVol",
			"unknownVol",
			"weight",
			"changeType",
			"stockVol",
			"stockBUVol",
			"stockSDVol",
			"stockUnknownVol",
		}, strings.Split(strings.Split(message, "#")[1], "|"))
		return messageDict, "M", err
	}

	if strings.HasPrefix(message, "L#") {
		messageDict, err := zipToDict([]string{
			"stockNo",
			"price",
			"vol",
			"accumulatedVol",
			"time",
			"ref",
			"side",
			"priceChange",
			"priceChangePercent",
			"changeType",
		}, strings.Split(strings.Split(message, "#")[1], "|"))
		return messageDict, "L", err
	}

	return nil, "", fmt.Errorf("not support message")
}

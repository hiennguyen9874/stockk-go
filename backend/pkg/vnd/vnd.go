package vnd

import (
	"fmt"
	"strconv"
	"strings"
)

func zipToDict(keys []string, values []string) (map[string]string, error) {
	// if len(keys) != len(values) {
	// 	return nil, fmt.Errorf("keys and values hasn't same length: len key: %v, len value: %v", keys, values)
	// }

	length := len(keys)
	if length > len(values) {
		length = len(values)
	}

	pairs := make(map[string]string, length)
	for i := 0; i < length; i++ {
		pairs[keys[i]] = values[i]
	}
	return pairs, nil
}

func DecodeMessage(message string) []string {
	var decodeString strings.Builder

	for i, char := range message {
		decodeString.WriteRune(rune(int(char) + i%5))
	}

	return strings.Split(decodeString.String(), "|")
}

func ConvertToFloat(valString string) (float32, error) {
	if valString == "" {
		return 0, nil
	}

	valFloat64, err := strconv.ParseFloat(valString, 32)
	if err != nil {
		return 0, err
	}

	return float32(valFloat64), nil
}

func MessageArrayToDict(messageType string, messageArray []string) (map[string]string, error) {
	switch strings.ToUpper(messageType) {
	case "S":
		switch strings.ToUpper(messageArray[0]) {
		case "SFU":
			switch strings.ToUpper(messageArray[2]) {
			case "ST":
				return zipToDict([]string{
					"code",
					"stockType",
					"floorCode",
					"basicPrice",
					"floorPrice",
					"ceilingPrice",
					"bidPrice01",
					"bidPrice02",
					"bidPrice03",
					"bidPrice04",
					"bidPrice05",
					"bidPrice06",
					"bidPrice07",
					"bidPrice08",
					"bidPrice09",
					"bidPrice10",
					"bidQtty01",
					"bidQtty02",
					"bidQtty03",
					"bidQtty04",
					"bidQtty05",
					"bidQtty06",
					"bidQtty07",
					"bidQtty08",
					"bidQtty09",
					"bidQtty10",
					"offerPrice01",
					"offerPrice02",
					"offerPrice03",
					"offerPrice04",
					"offerPrice05",
					"offerPrice06",
					"offerPrice07",
					"offerPrice08",
					"offerPrice09",
					"offerPrice10",
					"offerQtty01",
					"offerQtty02",
					"offerQtty03",
					"offerQtty04",
					"offerQtty05",
					"offerQtty06",
					"offerQtty07",
					"offerQtty08",
					"offerQtty09",
					"offerQtty10",
					"totalBidQtty",
					"totalOfferQtty",
					"tradingSessionId",
					"buyForeignQtty",
					"sellForeignQtty",
					"highestPrice",
					"lowestPrice",
					"accumulatedVal",
					"accumulatedVol",
					"matchPrice",
					"matchQtty",
					"currentPrice",
					"currentQtty",
					"projectOpen",
					"totalRoom",
					"currentRoom",
				}, messageArray[1:])
			case "W":
				return zipToDict([]string{
					"code",
					"stockType",
					"floorCode",
					"basicPrice",
					"floorPrice",
					"ceilingPrice",
					"underlyingSymbol",
					"issuerName",
					"exercisePrice",
					"exerciseRatio",
					"bidPrice01",
					"bidPrice02",
					"bidPrice03",
					"bidQtty01",
					"bidQtty02",
					"bidQtty03",
					"offerPrice01",
					"offerPrice02",
					"offerPrice03",
					"offerQtty01",
					"offerQtty02",
					"offerQtty03",
					"totalBidQtty",
					"totalOfferQtty",
					"tradingSessionId",
					"buyForeignQtty",
					"sellForeignQtty",
					"highestPrice",
					"lowestPrice",
					"accumulatedVal",
					"accumulatedVol",
					"matchPrice",
					"matchQtty",
					"currentPrice",
					"currentQtty",
					"projectOpen",
					"totalRoom",
					"currentRoom"},
					messageArray[1:],
				)
			default:
				return zipToDict([]string{
					"code",
					"stockType",
					"floorCode",
					"basicPrice",
					"floorPrice",
					"ceilingPrice",
					"bidPrice01",
					"bidPrice02",
					"bidPrice03",
					"bidQtty01",
					"bidQtty02",
					"bidQtty03",
					"offerPrice01",
					"offerPrice02",
					"offerPrice03",
					"offerQtty01",
					"offerQtty02",
					"offerQtty03",
					"totalBidQtty",
					"totalOfferQtty",
					"tradingSessionId",
					"buyForeignQtty",
					"sellForeignQtty",
					"highestPrice",
					"lowestPrice",
					"accumulatedVal",
					"accumulatedVol",
					"matchPrice",
					"matchQtty",
					"currentPrice",
					"currentQtty",
					"projectOpen",
					"totalRoom",
					"currentRoom",
					"iNav",
				}, messageArray[1:])
			}
		case "SBA":
			switch strings.ToUpper(messageArray[2]) {
			case "ST":
				return zipToDict([]string{
					"code",
					"stockType",
					"bidPrice01",
					"bidPrice02",
					"bidPrice03",
					"bidPrice04",
					"bidPrice05",
					"bidPrice06",
					"bidPrice07",
					"bidPrice08",
					"bidPrice09",
					"bidPrice10",
					"bidQtty01",
					"bidQtty02",
					"bidQtty03",
					"bidQtty04",
					"bidQtty05",
					"bidQtty06",
					"bidQtty07",
					"bidQtty08",
					"bidQtty09",
					"bidQtty10",
					"offerPrice01",
					"offerPrice02",
					"offerPrice03",
					"offerPrice04",
					"offerPrice05",
					"offerPrice06",
					"offerPrice07",
					"offerPrice08",
					"offerPrice09",
					"offerPrice10",
					"offerQtty01",
					"offerQtty02",
					"offerQtty03",
					"offerQtty04",
					"offerQtty05",
					"offerQtty06",
					"offerQtty07",
					"offerQtty08",
					"offerQtty09",
					"offerQtty10",
					"totalBidQtty",
					"totalOfferQtty",
				}, messageArray[1:])
			default:
				return zipToDict([]string{
					"code",
					"stockType",
					"bidPrice01",
					"bidPrice02",
					"bidPrice03",
					"bidQtty01",
					"bidQtty02",
					"bidQtty03",
					"offerPrice01",
					"offerPrice02",
					"offerPrice03",
					"offerQtty01",
					"offerQtty02",
					"offerQtty03",
					"totalBidQtty",
					"totalOfferQtty",
				}, messageArray[1:])
			}
		case "SMA":
			return zipToDict([]string{
				"code",
				"stockType",
				"tradingSessionId",
				"buyForeignQtty",
				"sellForeignQtty",
				"highestPrice",
				"lowestPrice",
				"accumulatedVal",
				"accumulatedVol",
				"matchPrice",
				"matchQtty",
				"currentPrice",
				"currentQtty",
				"projectOpen",
				"totalRoom",
				"currentRoom",
				"iNav",
			}, messageArray[1:])
		case "SBS":
			switch strings.ToUpper(messageArray[2]) {
			case "W":
				return zipToDict([]string{
					"code",
					"stockType",
					"floorCode",
					"basicPrice",
					"floorPrice",
					"ceilingPrice",
					"underlyingSymbol",
					"issuerName",
					"exercisePrice",
					"exerciseRatio",
				}, messageArray[1:])
			default:
				return zipToDict([]string{
					"code",
					"stockType",
					"floorCode",
					"basicPrice",
					"floorPrice",
					"ceilingPrice",
				}, messageArray[1:])
			}
		default:
			return nil, fmt.Errorf("not support message array type: %s", messageArray[0])
		}
	case "D":
		switch strings.ToUpper(messageArray[0]) {
		case "DFU":
			return zipToDict([]string{
				"code",
				"time",
				"bidPrice01",
				"bidPrice02",
				"bidPrice03",
				"bidPrice04",
				"bidPrice05",
				"bidPrice06",
				"bidPrice07",
				"bidPrice08",
				"bidPrice09",
				"bidPrice10",
				"bidQtty01",
				"bidQtty02",
				"bidQtty03",
				"bidQtty04",
				"bidQtty05",
				"bidQtty06",
				"bidQtty07",
				"bidQtty08",
				"bidQtty09",
				"bidQtty10",
				"offerPrice01",
				"offerPrice02",
				"offerPrice03",
				"offerPrice04",
				"offerPrice05",
				"offerPrice06",
				"offerPrice07",
				"offerPrice08",
				"offerPrice09",
				"offerPrice10",
				"offerQtty01",
				"offerQtty02",
				"offerQtty03",
				"offerQtty04",
				"offerQtty05",
				"offerQtty06",
				"offerQtty07",
				"offerQtty08",
				"offerQtty09",
				"offerQtty10",
				"totalBidQtty",
				"totalOfferQtty",
				"tradingSessionId",
				"buyForeignQtty",
				"sellForeignQtty",
				"highestPrice",
				"lowestPrice",
				"accumulatedVal",
				"accumulatedVol",
				"matchPrice",
				"currentPrice",
				"matchQtty",
				"currentQtty",
				"floorCode",
				"stockType",
				"tradingDate",
				"lastTradingDate",
				"underlying",
				"basicPrice",
				"floorPrice",
				"ceilingPrice",
				"openInterest",
				"openPrice",
			}, messageArray[1:])
		case "DBA":
			return zipToDict([]string{
				"code",
				"bidPrice01",
				"bidPrice02",
				"bidPrice03",
				"bidPrice04",
				"bidPrice05",
				"bidPrice06",
				"bidPrice07",
				"bidPrice08",
				"bidPrice09",
				"bidPrice10",
				"bidQtty01",
				"bidQtty02",
				"bidQtty03",
				"bidQtty04",
				"bidQtty05",
				"bidQtty06",
				"bidQtty07",
				"bidQtty08",
				"bidQtty09",
				"bidQtty10",
				"offerPrice01",
				"offerPrice02",
				"offerPrice03",
				"offerPrice04",
				"offerPrice05",
				"offerPrice06",
				"offerPrice07",
				"offerPrice08",
				"offerPrice09",
				"offerPrice10",
				"offerQtty01",
				"offerQtty02",
				"offerQtty03",
				"offerQtty04",
				"offerQtty05",
				"offerQtty06",
				"offerQtty07",
				"offerQtty08",
				"offerQtty09",
				"offerQtty10",
				"totalBidQtty",
				"totalOfferQtty",
			}, messageArray[1:])
		case "DMA":
			return zipToDict([]string{
				"code",
				"time",
				"tradingSessionId",
				"buyForeignQtty",
				"sellForeignQtty",
				"highestPrice",
				"lowestPrice",
				"accumulatedVal",
				"accumulatedVol",
				"matchPrice",
				"currentPrice",
				"matchQtty",
				"currentQtty",
			}, messageArray[1:])
		case "DBS":
			return zipToDict([]string{
				"code",
				"floorCode",
				"stockType",
				"tradingDate",
				"lastTradingDate",
				"underlying",
				"basicPrice",
				"floorPrice",
				"ceilingPrice",
				"openInterest",
				"openPrice",
			}, messageArray[1:])
		default:
			return nil, fmt.Errorf("not support message array type: %s", messageArray[0])
		}
	case "MI":
		switch strings.ToUpper(messageArray[0]) {
		case "MI":
			return zipToDict([]string{
				"floorCode",
				"tradingTime",
				"status",
				"advance",
				"noChange",
				"decline",
				"marketIndex",
				"priorMarketIndex",
				"highestIndex",
				"lowestIndex",
				"totalShareTraded",
				"totalValueTraded",
				"totalNormalTradedQttyRd",
				"totalNormalTradedValueRd",
				"predictionIndex",
			}, messageArray[1:])
		default:
			return nil, fmt.Errorf("not support message array type: %s", messageArray[0])
		}
	default:
		return nil, fmt.Errorf("not support message type: %s", messageType)
	}
}

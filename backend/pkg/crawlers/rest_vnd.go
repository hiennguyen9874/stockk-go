package crawlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
)

func (cr *restCrawler) VNDMapExchange(exchange string) (string, error) {
	switch exchange {
	case "UPCOM":
		return "UPCOM", nil
	case "HNX":
		return "HNX", nil
	case "HOSE":
		return "HOSE", nil
	default:
		return "", fmt.Errorf("not support exchange: %v", exchange)
	}
}

func (cr *restCrawler) VNDMapResolutionToString(resolution Resolution) (string, error) {
	switch resolution {
	case R1:
		return "1", nil
	case R5:
		return "5", nil
	case R15:
		return "15", nil
	case R30:
		return "30", nil
	case R60:
		return "60", nil
	case RD:
		return "D", nil
	default:
		return "", fmt.Errorf("not support resolution: %v", resolution)
	}
}

func (cr *restCrawler) VNDCrawlStockSymbols(ctx context.Context) ([]Ticker, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api-finfo.vndirect.com.vn/v4/stocks?q=type:IFC,ETF,STOCK~status:LISTED&fields=code,companyName,companyNameEng,shortName,floor,industryName&size=3000", nil)
	if err != nil {
		return nil, httpErrors.ErrMakeRequest(err)
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "vi,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", "https://trade.vndirect.com.vn")
	req.Header.Set("Referer", "https://trade.vndirect.com.vn/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="111", "Not(A:Brand";v="8", "Chromium";v="111"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)

	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		cr.logger.Warn("crawl stock symbols using vnd fail")
		return nil, httpErrors.ErrCallRequest(err)
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		cr.logger.Warn("crawl stock symbols using vnd fail")
		return nil, httpErrors.ErrReadBodyRequest(err)
	}

	type VNDStocksData struct {
		Code           string `json:"code"`
		CompanyName    string `json:"companyName"`
		CompanyNameEng string `json:"companyNameEng"`
		Floor          string `json:"floor"`
		ShortName      string `json:"shortName"`
	}

	type VNDStocksResponse struct {
		CurrentPage   int             `json:"currentPage"`
		Data          []VNDStocksData `json:"data"`
		Size          int             `json:"size"`
		TotalElements int             `json:"totalElements"`
		TotalPage     int             `json:"totalPage"`
	}

	var response VNDStocksResponse
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}

	tickers := make([]Ticker, response.TotalElements)
	for i, ticker := range response.Data {
		exchange, err := cr.VNDMapExchange(ticker.Floor)
		if err != nil {
			return nil, err
		}

		tickers[i] = Ticker{
			Symbol:    ticker.Code,
			Exchange:  exchange,
			FullName:  ticker.CompanyName,
			ShortName: ticker.ShortName,
			Type:      "Stock",
		}
	}
	return tickers, nil
}

func (cr *restCrawler) VNDCrawlStockHistory(ctx context.Context, symbol string, resolution Resolution, from int64, to int64) ([]Bar, error) {
	strResolution, err := cr.VNDMapResolutionToString(resolution)
	if err != nil {
		cr.logger.Warn("crawl history using vnd fail, error when map resolution to string")
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://dchart-api.vndirect.com.vn/dchart/history?resolution=%v&symbol=%v&from=%v&to=%v", strResolution, symbol, from, to), nil)
	if err != nil {
		cr.logger.Warn("crawl history using vnd fail, error create request")
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "vi,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", "https://dchart.vndirect.com.vn")
	req.Header.Set("Referer", "https://dchart.vndirect.com.vn/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="111", "Not(A:Brand";v="8", "Chromium";v="111"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)

	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		cr.logger.Warn("crawl history using vnd fail, error when request")
		return nil, httpErrors.ErrCallRequest(err)
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		cr.logger.Warn("crawl history using vnd fail, error when read response")
		return nil, httpErrors.ErrReadBodyRequest(err)
	}

	type VNDHistoryData struct {
		C []float64 `json:"c"`
		H []float64 `json:"h"`
		L []float64 `json:"l"`
		O []float64 `json:"o"`
		T []int64   `json:"t"`
		V []int64   `json:"v"`
		S string    `json:"s"`
	}

	var response VNDHistoryData
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		cr.logger.Warn("crawl history using vnd fail, error when deserialize response")
		return nil, err
	}

	bars := make([]Bar, len(response.C))
	for i := 0; i < len(response.C); i++ {
		bars[i] = Bar{
			Time:   time.Unix(response.T[i], 0),
			Open:   response.O[i],
			High:   response.H[i],
			Low:    response.L[i],
			Close:  response.C[i],
			Volume: response.V[i],
		}
	}

	return bars, nil
}

func (cr *restCrawler) VNDTransformMessage(ctx context.Context, message string) ([]string, error) {
	var decodeString strings.Builder

	for i, char := range message {
		decodeString.WriteRune(rune(int(char) + i%5))
	}

	return strings.Split(decodeString.String(), "|")[1:], nil
}

func (cr *restCrawler) VNDCrawlStockSnapshot(ctx context.Context, symbols []string) ([]StockSnapshot, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://price-api.vndirect.com.vn/stocks/snapshot?code=%v", strings.Join(symbols, ",")), nil)
	if err != nil {
		log.Fatal(err)
	}

	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		cr.logger.Warn("crawl history using vnd fail")
		return nil, httpErrors.ErrCallRequest(err)
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		cr.logger.Warn("crawl history using vnd fail")
		return nil, httpErrors.ErrReadBodyRequest(err)
	}

	type VNDStockSnapshot []string

	var response VNDStockSnapshot
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}

	var snapshots []StockSnapshot

	for _, msg := range response {
		messageArray := VNDDecodeMessage(msg)

		messageDict, err := VNDMessageArrayToDict("S", messageArray)
		if err != nil {
			return nil, err
		}

		var snapshot StockSnapshot

		if value, ok := messageDict["accumulatedVal"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.AccumulatedVal = valueFloat
		}
		if value, ok := messageDict["accumulatedVol"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.AccumulatedVol = valueFloat
		}
		if value, ok := messageDict["basicPrice"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.BasicPrice = valueFloat
		}
		if value, ok := messageDict["buyForeignQtty"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.BuyForeignQtty = valueFloat
		}
		if value, ok := messageDict["ceilingPrice"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.CeilingPrice = valueFloat
		}
		if value, ok := messageDict["code"]; ok && value != "" {
			snapshot.Code = value
		}
		if value, ok := messageDict["currentRoom"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.CurrentRoom = valueFloat
		}
		if value, ok := messageDict["floorCode"]; ok && value != "" {
			snapshot.FloorCode = value
		}
		if value, ok := messageDict["floorPrice"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.FloorPrice = valueFloat
		}
		if value, ok := messageDict["highestPrice"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.HighestPrice = valueFloat
		}
		if value, ok := messageDict["lowestPrice"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.LowestPrice = valueFloat
		}
		if value, ok := messageDict["matchPrice"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.MatchPrice = valueFloat
		}
		if value, ok := messageDict["matchQtty"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.MatchQtty = valueFloat
		}
		if value, ok := messageDict["projectOpen"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.ProjectOpen = valueFloat
		}
		if value, ok := messageDict["sellForeignQtty"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.SellForeignQtty = valueFloat
		}
		if value, ok := messageDict["totalRoom"]; ok && value != "" {
			valueFloat, err := ConvertToFloat(value)
			if err != nil {
				return nil, err
			}
			snapshot.TotalRoom = valueFloat
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots, nil
}

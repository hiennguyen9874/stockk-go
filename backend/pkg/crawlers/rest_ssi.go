package crawlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
)

func (cr *restCrawler) SSIMapExchange(exchange string) (string, error) {
	switch exchange {
	case "UpcomIndex":
		return "UPCOM", nil
	case "HNXIndex":
		return "HNX", nil
	case "VNINDEX":
		return "HOSE", nil
	default:
		return "", fmt.Errorf("not support exchange: %v", exchange)
	}
}

func (cr *restCrawler) SSICrawlStockSymbols(ctx context.Context) ([]Ticker, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://fiin-core.ssi.com.vn/Master/GetListOrganization?language=vi", nil)
	if err != nil {
		cr.logger.Warn("crawl stock symbols using ssi fail")
		return nil, httpErrors.ErrMakeRequest(err)
	}

	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		cr.logger.Warn("crawl stock symbols using ssi fail")
		return nil, httpErrors.ErrCallRequest(err)
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		cr.logger.Warn("crawl stock symbols using ssi fail")
		return nil, httpErrors.ErrReadBodyRequest(err)
	}

	type SSIStocksData struct {
		OrganCode      string `json:"organCode"`
		Ticker         string `json:"ticker"`
		ComGroupCode   string `json:"comGroupCode"`
		IcbCode        string `json:"icbCode"`
		OrganTypeCode  string `json:"organTypeCode"`
		ComTypeCode    string `json:"comTypeCode"`
		OrganName      string `json:"organName"`
		OrganShortName string `json:"organShortName"`
	}

	type SSIStocksResponse struct {
		Page       int             `json:"page"`
		PageSize   int             `json:"pageSize"`
		TotalCount int             `json:"totalCount"`
		Items      []SSIStocksData `json:"items"`
		PackageId  *int            `json:"packageId"`
		Status     string          `json:"status"`
		Errors     *string         `json:"errors"`
	}

	var response SSIStocksResponse
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		cr.logger.Warn("crawl history using ssi fail, error when deserialize response")
		return nil, err
	}

	tickers := make([]Ticker, response.TotalCount)
	for i, ticker := range response.Items {
		exchange, err := cr.SSIMapExchange(ticker.ComGroupCode)
		if err != nil {
			return nil, err
		}

		tickers[i] = Ticker{
			Symbol:    ticker.Ticker,
			Exchange:  exchange,
			FullName:  ticker.OrganName,
			ShortName: ticker.OrganShortName,
			Type:      "Stock",
		}
	}
	return tickers, nil
}

func (cr *restCrawler) SSIMapResolutionToString(resolution Resolution) (string, error) {
	switch resolution {
	case R1:
		return "1", nil
	case RD:
		return "D", nil
	default:
		return "", fmt.Errorf("not support resolution: %v", resolution)
	}
}

func (cr *restCrawler) SSICrawlStockHistory(ctx context.Context, symbol string, resolution Resolution, from int64, to int64) ([]Bar, error) {
	strResolution, err := cr.SSIMapResolutionToString(resolution)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://iboard.ssi.com.vn/dchart/api/history?resolution=%v&symbol=%v&from=%v&to=%v", strResolution, symbol, from, to), nil)
	if err != nil {
		cr.logger.Warn("crawl history using ssi fail")
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "vi,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Not?A_Brand";v="8", "Chromium";v="108", "Google Chrome";v="108"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)

	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		cr.logger.Warn("crawl history using ssi fail")
		return nil, httpErrors.ErrCallRequest(err)
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		cr.logger.Warn("crawl history using ssi fail")
		return nil, httpErrors.ErrReadBodyRequest(err)
	}

	// cr.logger.Info(responseData)

	type SSIHistoryData struct {
		C []string `json:"c"`
		H []string `json:"h"`
		L []string `json:"l"`
		O []string `json:"o"`
		T []int64  `json:"t"`
		V []string `json:"v"`
		S string   `json:"s"`
	}

	var response SSIHistoryData
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}

	bars := make([]Bar, len(response.C))
	for i := 0; i < len(response.C); i++ {
		// cr.logger.Info(response.O[i])

		Open, err := strconv.ParseFloat(response.O[i], 64)
		if err != nil {
			return nil, err
		}
		High, err := strconv.ParseFloat(response.H[i], 64)
		if err != nil {
			return nil, err
		}
		Low, err := strconv.ParseFloat(response.L[i], 64)
		if err != nil {
			return nil, err
		}
		Close, err := strconv.ParseFloat(response.C[i], 64)
		if err != nil {
			return nil, err
		}
		Volume, err := strconv.ParseInt(response.V[i], 10, 64)
		if err != nil {
			return nil, err
		}

		bars[i] = Bar{
			Time:   time.Unix(response.T[i], 0),
			Open:   Open,
			High:   High,
			Low:    Low,
			Close:  Close,
			Volume: Volume,
		}
	}

	return bars, nil
}

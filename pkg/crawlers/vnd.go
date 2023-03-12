package crawlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
)

func (cr *crawler) VNDCrawlStockSymbols() ([]Ticker, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api-finfo.vndirect.com.vn/v4/stocks?q=type:IFC,ETF,STOCK~status:LISTED&fields=code,companyName,companyNameEng,shortName,floor,industryName&size=3000", nil)
	if err != nil {
		return nil, httpErrors.ErrMakeRequest(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, httpErrors.ErrCallRequest(err)
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
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
	json.Unmarshal(responseData, &response)

	tickers := make([]Ticker, response.TotalElements)
	for i, ticker := range response.Data {
		tickers[i] = Ticker{
			Symbol:    ticker.Code,
			Exchange:  ticker.Floor,
			FullName:  ticker.CompanyName,
			ShortName: ticker.ShortName,
			Type:      "Stock",
		}
	}
	return tickers, nil
}

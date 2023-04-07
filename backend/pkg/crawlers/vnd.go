package crawlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
)

func (cr *crawler) VNDMapExchange(exchange string) (string, error) {
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

func (cr *crawler) VNDMapResolutionToString(resolution Resolution) (string, error) {
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

func (cr *crawler) VNDCrawlStockSymbols(ctx context.Context) ([]Ticker, error) {
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

func (cr *crawler) VNDCrawlStockHistory(ctx context.Context, symbol string, resolution Resolution, from int64, to int64) ([]Bar, error) {
	strResolution, err := cr.VNDMapResolutionToString(resolution)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://dchart-api.vndirect.com.vn/dchart/history?resolution=%v&symbol=%v&from=%v&to=%v", strResolution, symbol, from, to), nil)
	if err != nil {
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
		cr.logger.Warn("crawl history using vnd fail")
		return nil, httpErrors.ErrCallRequest(err)
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		cr.logger.Warn("crawl history using vnd fail")
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

func (cr *crawler) VNDTransformMessage(ctx context.Context, message string) ([]string, error) {
	var decodeString strings.Builder

	for i, char := range message {
		decodeString.WriteRune(rune(int(char) + i%5))
	}

	return strings.Split(decodeString.String(), "|")[1:], nil
}

func (cr *crawler) VNDCrawlStockSnapshot(ctx context.Context, symbols []string) ([]StockSnapshot, error) {
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

	convertToFloat := func(valString string) (float32, error) {
		if valString == "" {
			return 0, nil
		}

		valFloat64, err := strconv.ParseFloat(valString, 32)
		if err != nil {
			return 0, err
		}

		return float32(valFloat64), nil
	}

	for _, msg := range response {
		snapshotArray, _ := cr.VNDTransformMessage(ctx, msg)

		var snapshot StockSnapshot

		if snapshotArray[1] == "S" && snapshotArray[2] == "10" {
			snapshot.Ticker = snapshotArray[0]

			refprice, err := convertToFloat(snapshotArray[3])
			if err != nil {
				return nil, err
			}
			snapshot.RefPrice = float32(refprice)
			ceilprice, err := convertToFloat(snapshotArray[5])
			if err != nil {
				return nil, err
			}
			snapshot.CeilPrice = float32(ceilprice)
			floorprice, err := convertToFloat(snapshotArray[4])
			if err != nil {
				return nil, err
			}
			snapshot.FloorPrice = float32(floorprice)
			tltvol, err := convertToFloat(snapshotArray[26])
			if err != nil {
				return nil, err
			}
			snapshot.TltVol = float32(tltvol)
			tltval, err := convertToFloat(snapshotArray[25])
			if err != nil {
				return nil, err
			}
			snapshot.TltVal = float32(tltval)
			priceb3, err := convertToFloat(snapshotArray[8])
			if err != nil {
				return nil, err
			}
			snapshot.PriceB3 = float32(priceb3)
			priceb2, err := convertToFloat(snapshotArray[7])
			if err != nil {
				return nil, err
			}
			snapshot.PriceB2 = float32(priceb2)
			priceb1, err := convertToFloat(snapshotArray[6])
			if err != nil {
				return nil, err
			}
			snapshot.PriceB1 = float32(priceb1)
			volb3, err := convertToFloat(snapshotArray[11])
			if err != nil {
				return nil, err
			}
			snapshot.VolB3 = float32(volb3)
			volb2, err := convertToFloat(snapshotArray[10])
			if err != nil {
				return nil, err
			}
			snapshot.VolB2 = float32(volb2)
			volb1, err := convertToFloat(snapshotArray[9])
			if err != nil {
				return nil, err
			}
			snapshot.VolB1 = float32(volb1)
			price, err := convertToFloat(snapshotArray[27])
			if err != nil {
				return nil, err
			}
			snapshot.Price = float32(price)
			vol, err := convertToFloat(snapshotArray[28])
			if err != nil {
				return nil, err
			}
			snapshot.Vol = float32(vol)
			prices3, err := convertToFloat(snapshotArray[14])
			if err != nil {
				return nil, err
			}
			snapshot.PriceS3 = float32(prices3)
			prices2, err := convertToFloat(snapshotArray[13])
			if err != nil {
				return nil, err
			}
			snapshot.PriceS2 = float32(prices2)
			prices1, err := convertToFloat(snapshotArray[12])
			if err != nil {
				return nil, err
			}
			snapshot.PriceS1 = float32(prices1)
			vols3, err := convertToFloat(snapshotArray[15])
			if err != nil {
				return nil, err
			}
			snapshot.VolS3 = float32(vols3)
			vols2, err := convertToFloat(snapshotArray[16])
			if err != nil {
				return nil, err
			}
			snapshot.VolS2 = float32(vols2)
			vols1, err := convertToFloat(snapshotArray[17])
			if err != nil {
				return nil, err
			}
			snapshot.VolS1 = float32(vols1)
			high, err := convertToFloat(snapshotArray[23])
			if err != nil {
				return nil, err
			}
			snapshot.High = float32(high)
			low, err := convertToFloat(snapshotArray[24])
			if err != nil {
				return nil, err
			}
			snapshot.Low = float32(low)
			buyforeign, err := convertToFloat(snapshotArray[21])
			if err != nil {
				return nil, err
			}
			snapshot.BuyForeign = float32(buyforeign)
			sellforeign, err := convertToFloat(snapshotArray[22])
			if err != nil {
				return nil, err
			}
			snapshot.SellForeign = float32(sellforeign)

			snapshots = append(snapshots, snapshot)

		} else if snapshotArray[1] == "ST" && (snapshotArray[2] == "02" || snapshotArray[2] == "03") {
			snapshot.Ticker = snapshotArray[0]

			refprice, err := convertToFloat(snapshotArray[3])
			if err != nil {
				return nil, err
			}
			snapshot.RefPrice = float32(refprice)
			ceilprice, err := convertToFloat(snapshotArray[5])
			if err != nil {
				return nil, err
			}
			snapshot.CeilPrice = float32(ceilprice)
			floorprice, err := convertToFloat(snapshotArray[4])
			if err != nil {
				return nil, err
			}
			snapshot.FloorPrice = float32(floorprice)
			tltvol, err := convertToFloat(snapshotArray[54])
			if err != nil {
				return nil, err
			}
			snapshot.TltVol = float32(tltvol)
			tltval, err := convertToFloat(snapshotArray[53])
			if err != nil {
				return nil, err
			}
			snapshot.TltVal = float32(tltval)
			priceb3, err := convertToFloat(snapshotArray[8])
			if err != nil {
				return nil, err
			}
			snapshot.PriceB3 = float32(priceb3)
			priceb2, err := convertToFloat(snapshotArray[7])
			if err != nil {
				return nil, err
			}
			snapshot.PriceB2 = float32(priceb2)
			priceb1, err := convertToFloat(snapshotArray[6])
			if err != nil {
				return nil, err
			}
			snapshot.PriceB1 = float32(priceb1)
			volb3, err := convertToFloat(snapshotArray[18])
			if err != nil {
				return nil, err
			}
			snapshot.VolB3 = float32(volb3)
			volb2, err := convertToFloat(snapshotArray[17])
			if err != nil {
				return nil, err
			}
			snapshot.VolB2 = float32(volb2)
			volb1, err := convertToFloat(snapshotArray[16])
			if err != nil {
				return nil, err
			}
			snapshot.VolB1 = float32(volb1)
			price, err := convertToFloat(snapshotArray[55])
			if err != nil {
				return nil, err
			}
			snapshot.Price = float32(price)
			vol, err := convertToFloat(snapshotArray[56])
			if err != nil {
				return nil, err
			}
			snapshot.Vol = float32(vol)
			prices3, err := convertToFloat(snapshotArray[28])
			if err != nil {
				return nil, err
			}
			snapshot.PriceS3 = float32(prices3)
			prices2, err := convertToFloat(snapshotArray[27])
			if err != nil {
				return nil, err
			}
			snapshot.PriceS2 = float32(prices2)
			prices1, err := convertToFloat(snapshotArray[26])
			if err != nil {
				return nil, err
			}
			snapshot.PriceS1 = float32(prices1)
			vols3, err := convertToFloat(snapshotArray[38])
			if err != nil {
				return nil, err
			}
			snapshot.VolS3 = float32(vols3)
			vols2, err := convertToFloat(snapshotArray[37])
			if err != nil {
				return nil, err
			}
			snapshot.VolS2 = float32(vols2)
			vols1, err := convertToFloat(snapshotArray[36])
			if err != nil {
				return nil, err
			}
			snapshot.VolS1 = float32(vols1)
			high, err := convertToFloat(snapshotArray[51])
			if err != nil {
				return nil, err
			}
			snapshot.High = float32(high)
			low, err := convertToFloat(snapshotArray[52])
			if err != nil {
				return nil, err
			}
			snapshot.Low = float32(low)
			buyforeign, err := convertToFloat(snapshotArray[49])
			if err != nil {
				return nil, err
			}
			snapshot.BuyForeign = float32(buyforeign)
			sellforeign, err := convertToFloat(snapshotArray[50])
			if err != nil {
				return nil, err
			}
			snapshot.SellForeign = float32(sellforeign)

			snapshots = append(snapshots, snapshot)
		} else {
			return nil, errors.New("not support message")
		}
	}

	return snapshots, nil
}

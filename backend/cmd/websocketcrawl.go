package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hiennguyen9874/stockk-go/config"
	stockSnapshotRepository "github.com/hiennguyen9874/stockk-go/internal/stockssnapshot/repository"
	stockSnapshotUseCase "github.com/hiennguyen9874/stockk-go/internal/stockssnapshot/usecase"
	tickerRepository "github.com/hiennguyen9874/stockk-go/internal/tickers/repository"
	tickerUseCase "github.com/hiennguyen9874/stockk-go/internal/tickers/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/db/postgres"
	"github.com/hiennguyen9874/stockk-go/pkg/db/redis"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
	"github.com/hiennguyen9874/stockk-go/pkg/sentry"
	"github.com/hiennguyen9874/stockk-go/pkg/vnd"
	"github.com/hiennguyen9874/stockk-go/pkg/websocketCrawlers"
	"github.com/spf13/cobra"
)

var websocketCrawlCmd = &cobra.Command{
	Use:   "websocketcrawl",
	Short: "websocket crawl",
	Long:  "websocket crawl",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()

		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()
		appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

		err := sentry.Init(cfg)
		if err != nil {
			appLogger.Fatalf("Sentry init: %s", err)
		}
		defer sentry.Flush()

		psqlDB, err := postgres.NewPsqlDB(cfg)
		if err != nil {
			appLogger.Fatalf("Postgresql init: %s", err)
		} else {
			appLogger.Infof("Postgres connected")
		}

		redisClient := redis.NewRedis(cfg)

		// Repository
		tickerPgRepo := tickerRepository.CreateTickerPgRepository(psqlDB)
		tickerRedisRepo := tickerRepository.CreateTickerRedisRepository(redisClient)
		stockSnapshotRedisRepo := stockSnapshotRepository.CreateStockSnapshotRedisRepository(redisClient)

		// UseCase
		tickerUC := tickerUseCase.CreateTickerUseCaseI(tickerPgRepo, tickerRedisRepo, cfg, appLogger)
		stockSnapshotUC := stockSnapshotUseCase.CreateTickerUseCaseI(tickerUC, stockSnapshotRedisRepo, cfg, appLogger)

		ctx := context.Background()

		websocketCrawlers := websocketCrawlers.NewWebsocketCrawlers(cfg, appLogger)

		err = websocketCrawlers.Connect()
		if err != nil {
			appLogger.Fatalf("Websocket connect fail: %s", err)
		}
		defer websocketCrawlers.Close()

		go func() {
			for {
				err := websocketCrawlers.Connect()
				if err != nil {
					appLogger.Warnf("can't not connect to websocket: %v", err)
					time.Sleep(10 * time.Second)
					continue
				}

				errCh := make(chan error)

				go func() {
					tickers, err := tickerUC.GetAllActive(ctx, true)
					if err != nil {
						errCh <- err
						return
					}

					var symbols []string
					for _, ticker := range tickers {
						symbols = append(symbols, ticker.Symbol)
					}

					err = websocketCrawlers.WriteMessage([]string{
						"a",
						fmt.Sprintf("s|S:%v", strings.Join(symbols, ",")),
						fmt.Sprintf("ss|S:code=%v", strings.Join(symbols, ",")),
					})
					if err != nil {
						errCh <- err
						return
					}

					for {
						err := websocketCrawlers.WriteMessage([]string{"1"})
						if err != nil {
							errCh <- err
							return
						}

						time.Sleep(10 * time.Second)
					}
				}()

				go func() {
					for message := range websocketCrawlers.ReadMessage() {
						if message.MessageErr != nil {
							errCh <- err
							return
						}

						if message.MessageType == nil || message.MessageDict == nil {
							continue
						}

						switch *message.MessageType {
						case "S":
							values := make(map[string]interface{})

							if code, ok := (*message.MessageDict)["code"]; ok && code != "" {
								if value, ok := (*message.MessageDict)["accumulatedVal"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["AccumulatedVal"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["accumulatedVol"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["AccumulatedVol"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["basicPrice"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["BasicPrice"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["buyForeignQtty"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["BuyForeignQtty"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["ceilingPrice"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["CeilingPrice"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["currentRoom"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["CurrentRoom"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["floorCode"]; ok && value != "" {
									values["FloorCode"] = value
								}
								if value, ok := (*message.MessageDict)["floorPrice"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["FloorPrice"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["highestPrice"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["HighestPrice"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["lowestPrice"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["LowestPrice"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["matchPrice"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["MatchPrice"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["matchQtty"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["MatchQtty"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["projectOpen"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["ProjectOpen"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["sellForeignQtty"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["SellForeignQtty"] = valueFloat
								}
								if value, ok := (*message.MessageDict)["totalRoom"]; ok && value != "" {
									valueFloat, err := vnd.ConvertToFloat(value)
									if err != nil {
										errCh <- err
										return
									}
									values["TotalRoom"] = valueFloat
								}

								err := stockSnapshotUC.UpdateStockSnapshotBySymbol(ctx, code, values)

								if err != nil {
									errCh <- err
								}
							}
						default:
							errCh <- fmt.Errorf("not support message type: %v", *message.MessageType)
							return
						}
					}
				}()

				err = <-errCh
				if err != nil {
					appLogger.Warnf("error: %v", err)
				}

				time.Sleep(10 * time.Second)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		sig := <-quit

		appLogger.Infof("Shutting down server... Reason: %s", sig)
	},
}

func init() {
	RootCmd.AddCommand(websocketCrawlCmd)
}

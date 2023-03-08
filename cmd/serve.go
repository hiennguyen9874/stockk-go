package cmd

import (
	"log"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/server"
	"github.com/hiennguyen9874/stockk-go/pkg/db/postgres"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()

		psqlDB, err := postgres.NewPsqlDB(cfg)
		if err != nil {
			log.Fatalf("Postgresql init: %s", err)
		} else {
			log.Println("Postgres connected")
		}

		server, err := server.NewServer(cfg, psqlDB)
		if err != nil {
			log.Fatal(err)
		}
		server.Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

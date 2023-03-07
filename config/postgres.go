package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetDNSConfig() string {
	masterName := viper.GetString("MASTER_DB_NAME")
	masterUser := viper.GetString("MASTER_DB_USER")
	masterPassword := viper.GetString("MASTER_DB_PASSWORD")
	masterHost := viper.GetString("MASTER_DB_HOST")
	masterPort := viper.GetString("MASTER_DB_PORT")
	masterSslMode := viper.GetString("MASTER_SSL_MODE")

	masterDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		masterHost, masterUser, masterPassword, masterName, masterPort, masterSslMode,
	)
	return masterDSN
}

# Go Restful API Boilerplate

An API Boilerplate written in Golang with chi-route and Gorm. Write restful API with fast development and developer friendly.

## Architecture

In this project use 3 layer architecture

-   Models
-   Repository
-   Usecase
-   Delivery

## Features

-   CRUD
-   Jwt, refresh token saved in redis
-   Cached user in redis
-   Email verification
-   Forget/reset password, send email

## Technical

-   `chi`: router and middleware
-   `viper`: configuration
-   `cobra`: CLI features
-   `gorm`: orm
-   `validator`: data validation
-   `jwt`: jwt authentication
-   `zap`: logger
-   `gomail`: email
-   `hermes`: generate email body
-   `air`: hot-reload

## Start Application

### Generate the Private and Public Keys

-   Generate the private and public keys: [travistidwell.com/jsencrypt/demo/](https://travistidwell.com/jsencrypt/demo/)
-   Copy the generated private key and visit this Base64 encoding website to convert it to base64
-   Copy the base64 encoded key and add it to the `config/config-local.yml` file as `jwt`
-   Similar for public key

### Stmp mail config

-   Create [mailtrap](https://mailtrap.io/) account
-   Create new inboxes
-   Update smtp config `config/config-local.yml` file as `smtpEmail`

### Run

-   `docker-compose up`

## TODO

-   React-Query

## Acknowledgements

-   https://github.com/c9s/bbgo
-   https://github.com/Gituser143/cryptgo
-   https://github.com/m1/go-finnhub
-   https://github.com/achannarasappa/ticker
-   https://betterprogramming.pub/build-a-real-time-crypto-ticker-with-go-and-influxdb-2-89e968c65b7e
-   https://ajiybanesij.medium.com/building-go-applications-using-influxdb-87b462fd9d70
-   https://github.com/wentao-yang/StockAnalysisApplication

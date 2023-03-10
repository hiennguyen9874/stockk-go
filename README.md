# TODO

- Swagger
- Cookie
- Session using redis
- Traefik
- Celery
- Config using .env
- Save refresh token in redis, refresh token will check token available in redis, logout will remove token in redis

## Boilerplate Structure

<pre>├── <font color="#3465A4"><b>config</b></font>
├── <font color="#3465A4"><b>controllers</b></font>
├── <font color="#3465A4"><b>helpers</b></font>
├── <font color="#3465A4"><b>infra</b></font>
│   ├── <font color="#3465A4">database</font>
│   └── <font color="#3465A4">logger</font>
├── <font color="#3465A4"><b>migrations</b></font>
├── <font color="#3465A4"><b>models</b></font>
├── <font color="#3465A4"><b>repository</b></font>
├── <font color="#3465A4"><b>routers</b></font>
│   ├── <font color="#3465A4">middlewares</font>
</pre>

## JWT

- Generate the Private and Public Keys (for both access token and refresh token)
  - Generate the private and public keys: [Link](https://travistidwell.com/jsencrypt/demo/)
  - Copy the generated private key and visit this Base64 encoding website to convert it to base64
  - Copy the base64 encoded key and add it to the app.env file as ACCESS_TOKEN_PRIVATE_KEY
  - Similar for public key

## Acknowledgements

- https://github.com/dhax/go-base
- https://github.com/akmamun/go-fication
- https://github.com/wpcodevo/golang-fiber-jwt
- https://github.com/wpcodevo/golang-fiber
- https://github.com/kienmatu/togo
- https://github.com/AleksK1NG/Go-Clean-Architecture-REST-API
- https://github.com/bxcodec/go-clean-arch
- https://codevoweb.com/golang-and-gorm-user-registration-email-verification/
- https://codevoweb.com/golang-gorm-postgresql-user-registration-with-refresh-tokens/

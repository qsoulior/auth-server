# auth-server
[![Go Reference](https://img.shields.io/badge/-reference-007d9c?style=flat-square&logo=go&logoColor=fafafa&labelColor=555555)](https://pkg.go.dev/github.com/qsoulior/auth-server)
![Go Version](https://img.shields.io/github/go-mod/go-version/qsoulior/auth-server?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/qsoulior/auth-server?style=flat-square)](https://goreportcard.com/report/github.com/qsoulior/auth-server)

`auth-server` is a microservice that provides authentication and authorization using access and refresh tokens.

## ▶️ Features
- __High security__ — two types of tokens, asymmetric signing, fingerprint verification, password hashing and much more
- __Easy installation and startup__ — use [Docker](https://github.com/qsoulior/auth-server#-docker) or [Docker Compose](https://github.com/qsoulior/auth-server#-docker-compose)
- __Flexible and clear configuration__ — environment variables or .env file (see [Configuration](https://github.com/qsoulior/auth-server#%EF%B8%8F-configuration))
- __Simple interaction__ — RESTful API to interact with other services (see [Endpoints](https://github.com/qsoulior/auth-server#%EF%B8%8F-endpoints))

## ▶️ Installation and Startup
In production, the database must be running and migrations must be applied. 

Migrations are located in the `migrations` directory.
### 🖥️ Locally
Create [configuration](https://github.com/qsoulior/auth-server#%EF%B8%8F-configuration) file and specify its path instead of `<config_file>` in the following commands.

`POSTGRES_URI` must be set to URI of running PostgreSQL database.
`KEY_PRIVATE` and `KEY_PUBLIC` must be set to generated keys paths.
```
go mod download
go build ./cmd/main.go
main -c <config_file>
```
### 🐳 Docker
Create [configuration](https://github.com/qsoulior/auth-server#%EF%B8%8F-configuration) file and specify its path instead of `<config_file>` in the following commands.

Private and public keys are generated using the `ECDSA` algorithm when the image is built. There is no effect of changing `KEY_PRIVATE`, `KEY_PUBLIC` or `AT_ALG`.

`POSTGRES_URI` must be set to URI of running PostgreSQL database. 
```
docker build -t auth-server .
docker run -p <host_port>:<app_port> --env-file <config_file> --name auth-web auth-server
```

You can copy generated keys to host in the following way:
```console
$ docker ps
CONTAINER ID   IMAGE         COMMAND    CREATED         STATUS         PORTS                  NAMES
80fb44dc1638   auth-server   "./main"   4 seconds ago   Up 2 seconds   0.0.0.0:3000->80/tcp   auth-web

$ docker cp auth-web:/keys ./
Successfully copied 3.58kB

$ ls ./keys
ecdsa  ecdsa.pub
```
### 🐙 Docker Compose
As when running using [Docker](https://github.com/qsoulior/auth-server#-docker), private and public keys are generated when the image is built.

#### For development
Compose uses `configs/docker.dev.env` to configure web server for development. It also runs database server and applies migrations.

There is no effect of changing `POSTGRES_URI`.
```
docker compose -f docker-compose.dev.yaml up --build
```
#### For production
Compose uses `configs/docker.prod.env` to configure web server for production.

`POSTGRES_URI` must be changed to URI of running PostgreSQL database.
```
docker compose -f docker-compose.prod.yaml up --build
```

## ▶️ Configuration
Сonfiguration is performed using environment variables or .env file.

List of environment variables:
| Variable        | Default     | Constraint          | Description                                                   |
| -               | -           | -                   | -                                                             |
| `APP_NAME`      | auth        |                     | Application name used in the "iss" JWT claim                  |
| `APP_ENV`       | development |                     | Application environment                                       |
| `KEY_PUBLIC`    |             |                     | Path to public key encoded in PEM format                      |
| `KEY_PRIVATE`   |             |                     | Path to private key encoded in PEM format                     |
| `HTTP_HOST`     | 0.0.0.0     |                     | Host for the server to listen on                              |
| `HTTP_PORT`     | 3000        |                     | Port for the server to listen on                              |
| `HTTP_ORIGINS`  | *           | Separated by comma  | List of origins a cross-domain request can be executed from   |
| `POSTGRES_URI`  |             | [PostgreSQL connection URI](https://www.postgresql.org/docs/current/libpq-connect.html#id-1.7.3.8.3.6) | Database connection string in URI format |
| `AT_ALG`        | HS256       | [RFC7518](https://datatracker.ietf.org/doc/html/rfc7518#section-3.1), [RFC8037](https://datatracker.ietf.org/doc/html/rfc8037#section-3.1) | Algorithm used to sign the JWT |
| `AT_AGE`        | 15          | 1 — 60              | Number of __minutes__ until the access token expires          |
| `RT_CAP`        | 10          | > 1                 | Max number of refresh tokens per user until overwriting       |
| `RT_AGE`        | 30          | 1 — 365             | Number of __days__ until the refresh token expires            |
| `BCRYPT_COST`   | 4           | 4 — 31              | Cost parameter of bcrypt algorithm used for password hashing  |

.env file example:
```dotenv
APP_NAME=auth
APP_ENV=development
KEY_PUBLIC=/secrets/ecdsa.pub
KEY_PRIVATE=/secrets/ecdsa
HTTP_HOST=0.0.0.0
HTTP_PORT=3000
HTTP_ORIGINS=https://*.example1.com,http://example2.com
POSTGRES_URI=postgres://postgres:test@localhost:5432/postgres?search_path=auth
AT_ALG=ES384
AT_AGE=15
RT_CAP=10
RT_AGE=30
BCRYPT_COST=10
```

## ▶️ Tokens
### 🔐 Access token
Access token is a JSON Web Token (JWT) signed using one of the algorithms: `HMAC SHA`, `RSA`, `ECDSA` or `EdDSA`. Token contains a payload with two custom claims: `fingerprint` and `roles`.

Payload example:
```json
{
  "fingerprint": "fb57b63a63bb4923031a191fa0abd37db24d8c56c6ba33d26ca34529a505eeab",
  "roles": ["admin"],
  "iss": "auth",
  "sub": "522198cc-42d9-4b47-b20e-1def58dc2709",
  "exp": 1687173288
}
```
Access token is created by `auth-server` and used by other microservices to authorize requests.
Recipient microservice must parse token, check token issuer (`iss`) and expiration date (`exp`), compare `fingerprint` from payload with user fingerprint and optionally check `roles`. 

Subject claim (`sub`) contains user ID. 

Fingerprint is created from HTTP headers `Sec-CH-UA`, `User-Agent`, `Accept-Language`, `Upgrade-Insecure-Requests` and hashed using `SHA-256` in the following way:
```cpp
SHA256(Sec-CH-UA + ":" + User-Agent + ":" + Accept-Language + ":" + Upgrade-Insecure-Requests)
```
This repository also contains package `jwt` that provides [Parser](https://github.com/qsoulior/auth-server/blob/master/pkg/jwt/parser.go#L11) interface to parse access token:
```go
import (
	"fmt"
 	
	"github.com/qsoulior/auth-server/pkg/jwt"
)

parser, err := jwt.NewParser(jwt.Params{issuer, alg, publicKey})
if err != nil {
	return err
}
claims, err := parser.Parse(token)
if err != nil {
	return err
}
fmt.Println(claims.Subject) // "522198cc-42d9-4b47-b20e-1def58dc2709"
```

### ♻️ Refresh token
Refresh token is stored in database and used to refresh the access token.

Refresh token entity:
```go
type RefreshToken struct {
	ID          uuid.UUID `json:"id"`
	ExpiresAt   time.Time `json:"expires_at"`
	Fingerprint []byte    `json:"fingerprint"`
	Session     bool      `json:"session"`
	UserID      uuid.UUID `json:"-"`
}
```
This token is issued by the server upon successful authentication and is refreshed along with refresh of the access token. Client receives a cookie in response:
```http
Set-Cookie: refresh_token=da5067f7-0235-4ca2-ab38-a650e44d7bbc; Path=/v1/token; Expires=Sat, 22 Jul 2023 16:35:36 GMT; HttpOnly; Secure; SameSite=None
```

## ▶️ Endpoints
### 💁 Create user
`POST /user`

Request:
```json
{
  "name": "test",
  "password": "Test123$"
}
```
Response:
```
201 Created
```

### 💁 Get user
`GET /user`

Request:
```http
Authorization: Bearer <access_token>
```
Response:
```
200 OK
```
```json
{
  "id": "522198cc-42d9-4b47-b20e-1def58dc2709",
  "username": "test"
}
```

### 💁 Update user password
`PUT /user/password`

Request:
```http
Authorization: Bearer <access_token>
```
```json
{
  "current_password": "Ttest123$",
  "new_password": "Ttest123$"
}
```
Response:
```
204 No Content
```

### 💁 Delete user
`DELETE /user`

Request:
```http
Authorization: Bearer <access_token>
```
```json
{
  "password": "Ttest123$"
}
```
Response:
```
204 No Content
```

### 🔑 Create token
`POST /token`

Request:
```json
{
  "name": "test",
  "password": "Test123$",
  "session": false
}
```
Response:
```
201 Created
```
```http
Set-Cookie: refresh_token=da5067f7-0235-4ca2-ab38-a650e44d7bbc; Path=/v1/token; Expires=Sat, 22 Jul 2023 16:35:36 GMT; HttpOnly; Secure; SameSite=None
```
```json
{
  "access_token": "<access_token>"
}
```

### 🔑 Refresh token
`POST /token/refresh`

Request:
```http
Cookie: refresh_token=d337672c-d6e9-4058-b838-a634bbc5bddc
```
Response:
```
201 Created
```
```http
Set-Cookie: refresh_token=da5067f7-0235-4ca2-ab38-a650e44d7bbc; Path=/v1/token; Expires=Sat, 22 Jul 2023 16:35:36 GMT; HttpOnly; Secure; SameSite=None
```
```json
{
  "access_token": "<access_token>"
}
```

### 🔑 Revoke token
`POST /token/revoke`

Request:
```http
Cookie: refresh_token=d337672c-d6e9-4058-b838-a634bbc5bddc
```
Response:
```
204 No Content
```

### 🔑 Revoke all tokens
`POST /token/revoke-all`

Request:
```http
Cookie: refresh_token=d337672c-d6e9-4058-b838-a634bbc5bddc
```
Response:
```
204 No Content
```

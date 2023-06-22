# auth-server
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/qsoulior/auth-server?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/qsoulior/auth-server?style=flat-square)](https://goreportcard.com/report/github.com/qsoulior/auth-server)

`auth-server` is a microservice that provides authentication and authorization using access and refresh tokens.

## ▶️ Tokens

### 🔓 Access token
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

### ♻️ Refresh token
Refresh token is stored in database and used to refresh the access token.

Refresh token entity:
```go
type RefreshToken struct {
	ID          uuid.UUID `json:"id"`
	ExpiresAt   time.Time `json:"expires_at"`
	Fingerprint []byte    `json:"fingerprint"`
	UserID      uuid.UUID `json:"-"`
}
```
This token is issued by the server upon successful authentication and is refreshed along with refresh of the access token. Client receives a cookie in response:
```http
Set-Cookie: refresh_token=d337672c-d6e9-4058-b838-a634bbc5bddc; Expires=Wed, 19 Jul 2023 14:04:07 GMT; HttpOnly; Secure; SameSite=Strict
```

## ▶️ Installation and Running
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

`POSTGRES_URI` must be set to URI of running PostgreSQL database.
```
docker build -t auth-server .
docker run --env-file <config_file> auth-server
```
### 🐙 Docker Compose
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
`configs/dev.env`
```dotenv
APP_NAME=auth
APP_ENV=development
KEY_PUBLIC=/secrets/ecdsa.pub
KEY_PRIVATE=/secrets/ecdsa
HTTP_HOST=0.0.0.0
HTTP_PORT=3000
HTTP_ORIGINS=https://*.example1.com,http://example2.com
POSTGRES_URI=postgres://postgres:test@localhost:5432/postgres?search_path=auth
AT_ALG=ES256
AT_AGE=15
RT_CAP=10
RT_AGE=30
BCRYPT_COST=4
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
  "password": "Test123$"
}
```
Response:
```
201 Created
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
Cookie: refresh_token=d337672c-d6e9-4058-b838-a634bbc5bddc; Expires=Wed, 19 Jul 2023 14:04:07 GMT; HttpOnly; Secure; SameSite=Strict
```
Response:
```
201 Created
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
Cookie: refresh_token=d337672c-d6e9-4058-b838-a634bbc5bddc; Expires=Wed, 19 Jul 2023 14:04:07 GMT; HttpOnly; Secure; SameSite=Strict
```
Response:
```
204 No Content
```

### 🔑 Revoke all tokens
`POST /token/revoke-all`

Request:
```http
Cookie: refresh_token=d337672c-d6e9-4058-b838-a634bbc5bddc; Expires=Wed, 19 Jul 2023 14:04:07 GMT; HttpOnly; Secure; SameSite=Strict
```
Response:
```
204 No Content
```

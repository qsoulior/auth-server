# auth-server
`auth-server` is a microservice that provides authentication and authorization using access and refresh tokens.

## Tokens

### Access token
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

### Refresh token
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

## Run
> TODO

## Configuration
`configs/dev.env`
```properties
APP_NAME=auth
APP_ENV=development
KEY_PUBLIC_PATH=/secrets/ecdsa.pub
KEY_PRIVATE_PATH=/secrets/ecdsa
HTTP_PORT=3000
POSTGRES_URI=postgres://postgres:test123@localhost:5432/postgres?search_path=auth
AT_ALG=ES256
AT_AGE=15
RT_CAP=10
RT_AGE=30
BCRYPT_COST=4
```

## Endpoints
### Create user
`POST /user`

Request:
```json
{
  "name": "test",
  "password": "Test123$"
}
```
Response:
```json
201 Created
```

### Get user
`GET /user`

Request:
```http
Authorization: Bearer <access_token>
```
Response:
```json
200 OK
{
  "id": "522198cc-42d9-4b47-b20e-1def58dc2709",
  "username": "test"
}
```

### Update user password
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
```json
204 No Content
```

### Delete user
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
```json
204 No Content
```

### Create token
`POST /token`

Request:
```json
{
  "name": "test",
  "password": "Test123$"
}
```
Response:
```json
201 Created
{
  "access_token": "<access_token>"
}
```

### Refresh token
`POST /token/refresh`

Request:
```http
Cookie: refresh_token=d337672c-d6e9-4058-b838-a634bbc5bddc; Expires=Wed, 19 Jul 2023 14:04:07 GMT; HttpOnly; Secure; SameSite=Strict
```
Response:
```json
201 Created
{
  "access_token": "<access_token>"
}
```

### Revoke token
`POST /token/revoke`

Request:
```http
Cookie: refresh_token=d337672c-d6e9-4058-b838-a634bbc5bddc; Expires=Wed, 19 Jul 2023 14:04:07 GMT; HttpOnly; Secure; SameSite=Strict
```
Response:
```json
204 No Content
```

### Revoke all tokens
`POST /token/revoke-all`

Request:
```http
Cookie: refresh_token=d337672c-d6e9-4058-b838-a634bbc5bddc; Expires=Wed, 19 Jul 2023 14:04:07 GMT; HttpOnly; Secure; SameSite=Strict
```
Response:
```json
204 No Content
```
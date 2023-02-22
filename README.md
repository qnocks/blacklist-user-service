# Blacklist user service

## Usage

Create `.env` file in root directory:
```
PORT=8080
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_USER=root
POSTGRES_PASS=root
POSTGRES_DB=blacklistdb
POSTGRES_SSL=disable
PASSWORD_SALT=dasdjasdjasjj213
TOKEN_TTL=1  
TOKEN_SECRET=my-secret-key
```
Consider following:
* `POSTGRES_HOST` has to be the same with docker-compose service name for your postgres
* `TOKEN_TTL` value is expiration time for JWT token in **minutes**
* please do use `PASSWORD_SALT=dasdjasdjasjj213` as this is used for the pre-inserted password and don't do this in the production environment!

Run `make docker` to build docker image of backend\
Use `docker-compose up` to spin up backend with database

To see swagger docs go to `https://localhost:<your_port>/swagger/index.html`

To authorize request use login endpoint with following credentials:
```json
{
  "username": "admin",
  "password": "admin"
}
```

Returned token use to add to header as `Bearer <generated-token`

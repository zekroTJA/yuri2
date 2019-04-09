# API DOCS

## Index

- [REST API](#rest-api)

- [Web Socket API](#web-socket-api)

---

# REST API

The Yuri REST API is generally just used for authentication and for getting the local sounds index.

## Authentication

Because Yuri is using the Discord OAuth 2 flow, you need to manually generate an API token by authentication with the Yuri Discord App. Therefore, open the `/token` endpoint using a browser-like application to press `Authorize` in the Discord authorization page. Then, you will be relocated back to the `/token/authorize`* endpoint, which will respond with a randomly generated token, your Discord user ID and the tokens lifetime on successful authentication.  

The token's lifetime will be exteded each time authorizing yourself with it.

This token must be passed via an Basic `Authorization` header using the `/api/localsounds` endpoint. Also, you must pass the token and your user ID to authenticate the web socket API initialization.

**) The `/token/authorize` endpoint is not dcumented because it should only be used as callback for the Discord OAuth 2 flow.*

## Endpoints

### Get Token

> GET /token


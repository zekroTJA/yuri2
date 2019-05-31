# REST API

The Yuri REST API is generally just used for authentication and for getting data like lokal sounds list and the sound playing logs and stats for guilds.

## Index

- [Authentication](#authentication)

- [Parameters](#parameters)

- [Response Data](#response-data)

- [Rate Limiting](#rate-limiting)

- [Endpoints](#endpoints)
  - [Get Token](#get-token)
  - [Get Local Sounds](#get-local-sounds)
  - [Get Play Log](#get-play-log)
  - [Get Play Stats](#get-play-stats)
  - [Get Favorites](#get-favorites)
  - [Set Favorite](#set-favorite)
  - [Unset Favorite](#unset-favorite)

- [Admin Endpoints](#admin-endpoints)
  - [Get System Stats](#get-system-stats)
  - [Restart](#restart)
  - [Refetch](#refetch)

---

## Authentication

Because Yuri is using the Discord OAuth 2 flow, you need to manually generate an API token by authentication with the Yuri Discord App. Therefore, open the `/token` endpoint using a browser-like application to press `Authorize` in the Discord authorization page. Then, you will be relocated back to the `/token/authorize`* endpoint, which will respond with a randomly generated token, your Discord user ID and the tokens lifetime on successful authentication.  

The token's lifetime will be extended each time when you are authorizing yourself with it.

There are **two options** to autenticate against the REST API:

### I) Basic Authorization Header

Using the REST API, you need to pass your Discord User ID **and** your API token as base64 encoded value as Basic `Authorization` header **on every request**.

Example:

1. Your Discord User ID:  
`221905671296253953`  

2. Your API Token:  
`gDURWm1gkLEjmFcjKs1CzWkUIkIDJQ486iheIfcr728jb6MxG2RUaoLnTdCILxLJ`

3. Both assembled together:  
`221905671296253953:gDURWm1gkLEjmFcjKs1CzWkUIkIDJQ486iheIfcr728jb6MxG2RUaoLnTdCILxLJ`

4. Bas64 encoded, assembled authorization value:  
`MjIxOTA1NjcxMjk2MjUzOTUzOmdEVVJXbTFna0xFam1GY2pLczFDeldrVUlrSURKUTQ4NmloZUlmY3I3MjhqYjZNeEcyUlVhb0xuVGRDSUx4TEo=`

5. Resulting Basic Authorization Header:  
`Authorization: Basic MjIxOTA1NjcxMjk2MjUzOTUzOmdEVVJXbTF...`

### II) Authorization Cookies

If you have visited the `/login` endpoint and authorized the Discord API App, two cookies are set to authenticate against the REST and WS API.

```
< HTTP/1.1 307 Temporary Redirect
< Location: /
< Set-Cookie: token=UGVOEHNgu6P7iteHPTWrL08FvthQEmokZKdbY2jkQ6sxw7Y720vxvGOdJxxXRxDQ; Path=/
< Set-Cookie: userid=221905671296253953; Path=/
< Date: Thu, 11 Apr 2019 07:52:11 GMT
< Content-Length: 0
```

These cookeis will be automatically detected and checked for authentication if you sent them with your request to REST API endpoints.

## Parameters

Parameters with `default` formated names are **required** and parameters with *`italic`* formated names are ***optional***.  
The type *(int, string, ...)* and the passing method *(URL Query, JSON Body, Resource Path, ...)* are described in the parameter tables.

## Response Data

Generally, the API will never omit response data keys if they are unset or not existent. They will always be defined as `null` or as the default value o the specific data type like `""` for strings or `0` for integers.  
Those data properties are marked with *`italic`* font style.

## Rate Limiting

The TREST API is rate limited based on a per-user limiter globally over all endpoints.

Every **2000 Milliseconds, 1 token** is regenerated to a total burst number of **5 tokens**.

On each request, the response contains inforamtion about the users rate limiting status in the following three response headers:

```
< X-Ratelimit-Limit: 5
< X-Ratelimit-Remaining: 4
< X-Ratelimit-Reset: 0
```

If the rate limit has been exceed, an error response like following will be returned:

```
< HTTP/1.1 429 Too Many Requests
< Content-Type: application/json
< X-Ratelimit-Limit: 5
< X-Ratelimit-Remaining: 0
< X-Ratelimit-Reset: 1585
< Date: Mon, 15 Apr 2019 14:00:43 GMT
< Content-Length: 72
```
```json
{
  "error": {
    "code": 429,
    "message": "rate limit exceed"
  }
}
```

The value of the `X-Ratelimit-Reset` header indicates the time in milliseconds which the user has to wait until another request can be executed.

---

## Endpoints

### Get Token

> GET /token

Because this endpoint redirects to Discord's OAuth2 App authentication, this endpoint needs to be requested from a browser-like application which is capable of rendering HTML/CSS and executing JavaScript.

After authorization, Discord will redirect to the passed callback URI, which will be `/token/authorize`. The response of this will have content type `application/json`.

#### Response

```
> HTTP/1.1 307 Temporary Redirect
> Location: https://discordapp.com/api/oauth2/authorize?client_id=529947...
```

The resulting response of the callback to `/token/authorize` will have follwoing response format:

```
< HTTP/1.1 200 OK
< Content-Length: 166
< Content-Type: application/json
< Date: Wed, 10 Apr 2019 06:32:40 GMT
```
```json
{
  "token": "gDURWm1gkLEjmFcjKs1CzWkUIkIDJQ486iheIfcr728jb6MxG2RUaoLnTdCILxLJ",
  "user_id": "221905671296253953",
  "expires": "2019-04-17T08:29:21.2972027+02:00"
}
```

### Get Local Sounds

> GET /api/localsounds

#### Parameters

| Name | Passed by | Type | Description |
| -----|-----------|------|-------------|
| *`sort`* | URL Query | `string` | Wether sort results by:<br/>- `NAME`<br/>- `DATE` |
| *`from`* | URL Query | `int` | Start index of result entries. |
| *`limit`* | URL Query | `int` | Ammount of result entries. |



#### Response

```
< HTTP/1.1 200 OK
< Content-Length: 289
< Content-Type: application/json
< Date: Wed, 10 Apr 2019 06:32:40 GMT
```
```json
{
  "n": 3,
  "results": [
    {
      "name": "zugriff.pm3",
      "last_modified": "2019-03-13T09:45:21Z"
    },
    {
      "name": "oof.ogg",
      "last_modified": "2019-02-11T19:32:12Z"
    },
    {
      "name": "danke.wav",
      "last_modified": "2018-10-23T14:12:56Z"
    }
  ]
}
```

### Get Play Log

> GET /api/logs/:GUILDID

#### Parameters

| Name | Passed by | Type | Description |
| -----|-----------|------|-------------|
| `GUILDID` | Resource Path | `string` | The ID of the Discord guild. |
| *`from`* | URL Query | `int` | Start index of result entries. |
| *`limit`* | URL Query | `int` | Ammount of result entries. The default limit is `1000`. |


#### Response

```
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Wed, 10 Apr 2019 11:19:03 GMT
< Content-Length: 672
```
```json
{
  "n": 3,
  "results": [
    {
      "time": "2019-04-06T18:48:02Z",
      "user_id": "221905671296253953",
      "user_tag": "zekro#9131",
      "guild_id": "526196711962705925",
      "source": "local",
      "sound": "derbergruft"
    },
    {
      "time": "2019-04-06T18:47:53Z",
      "user_id": "221905671296253953",
      "user_tag": "zekro#9131",
      "guild_id": "526196711962705925",
      "source": "local",
      "sound": "derbergruft"
    },
    {
      "time": "2019-04-06T18:47:45Z",
      "user_id": "221905671296253953",
      "user_tag": "zekro#9131",
      "guild_id": "526196711962705925",
      "source": "local",
      "sound": "nice"
    }
  ]
}
```

### Get Play Stats

> GET /api/stats/:GUILDID

#### Parameters

| Name | Passed by | Type | Description |
| -----|-----------|------|-------------|
| `GUILDID` | Resource Path | `string` | The ID of the Discord guild. |
| *`limit`* | URL Query | `int` | Ammount of result entries. The default limit is `1000`. |


#### Response

```
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Wed, 10 Apr 2019 11:33:30 GMT
< Content-Length: 201
```
```json
{
  "n": 3,
  "results": [
    {
      "sound": "nice",
      "count": 98
    },
    {
      "sound": "derbergruft",
      "count": 23
    },
    {
      "sound": "danke",
      "count": 11
    }
  ]
}
```

### Get Favorites

> GET /api/favorites

#### Parameters

*No parameters passed.*

#### Response

```
< HTTP/1.1 200 OK
< Date: Sat, 11 May 2019 10:44:47 GMT
< Content-Type: application/json
< Content-Length: 104
```
```json
{
  "n": 4,
  "results": [
    "ojamoin",
    "jamoin",
    "ichsagenein",
    "echtjetzt"
  ]
}
```

### Set Favorite

> POST /api/favorite/:SOUND

#### Parameters

| Name | Passed by | Type | Description |
| -----|-----------|------|-------------|
| `SOUND` | Resource Path | `string` | The name of the sound. |

#### Response

```
< HTTP/1.1 201 Created
< Date: Sat, 11 May 2019 10:46:53 GMT
```

### Unset Favorite

> DELETE /api/favorite/:SOUND

#### Parameters

| Name | Passed by | Type | Description |
| -----|-----------|------|-------------|
| `SOUND` | Resource Path | `string` | The name of the sound. |

#### Response

```
< HTTP/1.1 200 Created
< Date: Sat, 11 May 2019 10:46:53 GMT
```

### Get Fast Trigger

> GET /api/settings/fasttrigger

#### Parameters

*No parameters passed.*

#### Response

```
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 31 May 2019 14:32:10 GMT
< Content-Length: 40
```
```json
{
  "ident": "nice",
  "random": false
}
```

### SET Fast Trigger

> POST /api/settings/fasttrigger

#### Parameters

| Name | Passed by | Type | Description |
| -----|-----------|------|-------------|
| `random` | JSON Body | `boolean` | Set fast trigger to random or not. |
| `ident` | JSON Body | `string` | Ident of sound to set fast trigger to if not random. |

#### Response

```
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 31 May 2019 14:34:28 GMT
< Content-Length: 40
```
```json
{
  "ident": "danke",
  "random": false
}
```

## Admin Endpoints

### Get System Stats

> GET /api/admin/stats

#### Parameters

*No parameters passed.*

#### Response

```
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 05 May 2019 01:22:14 GMT
< Content-Length: 377
```
```json
{
  "guilds": [
    {
      "name": "5w4gg3rn4ut_5t4t10n",
      "id": "526196711962705925"
    }
  ],
  "voice_connections": [],
  "system": {
    "arch": "amd64",
    "os": "windows",
    "go_version": "go1.11.2",
    "cpu_used_cores": 6,
    "go_routines": 15,
    "heap_use_b": 2768896,
    "stack_use_b": 327680,
    "uptime_seconds": 13123.9373705
  }
}
```

### Get Sound Stats

> GET /api/admin/soundstats

***Note:** `log_len` is the total ammount of recorded played sounds on all guilds and not for one single or the current recognized guild.*

#### Parameters

*No parameters passed.*

#### Response

```
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 05 May 2019 01:22:14 GMT
< Content-Length: 63
```
```json
{
  "sounds_len": 3,
  "log_len": 144,
  "size_b": 486432
}
```

### Restart

> POST /api/admin/restart

#### Parameters

*No parameters passed.*

#### Response

```
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 05 May 2019 01:22:14 GMT
```

### Refetch

> POST /api/admin/refetch

#### Parameters

*No parameters passed.*

#### Response

```
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 05 May 2019 01:22:14 GMT
```
# WS API

The web socket API is used for controlling players and receiving player events.

WS events will always be fired wether the command was executed from a Discord client or a client connected to the WS API.

## Index

- [Command Structure](#command-structure)

- [Event Structure](#event-structure)

- [Commands](#commands)
  - [INIT](#init)
  - [JOIN](#join)
  - [LEAVE](#leave)
  - [PLAY](#play)
  - [RANDOM](#random)
  - [STOP](#stop)
  - [VOLUME](#volume)

- [Events](#events)
  - [HELLO](#hello)
  - [ERROR](#error)
  - [PLAYING](#playing)
  - [END](#end)
  - [PLAY_ERROR](#play_error)
  - [STUCK](#stuck)
  - [VOLUME_CHANGED](#volume_changed)
  - [JOINED](#joined)
  - [LEFT](#left)

---

## Command Structure

Commands are send by the client in form of a JSON text message with following strucrure:

| Field | Type | Description |
|-------|------|-------------|
| `name` | `string` | Name of the command which must be all uppercase. |
| `data` | `any?` | The data playload which can be any type of JSON object. |

## Event Structure

Events are received by the client in form of a JSON text message with following structure:

| Field | Type | Description |
|-------|------|-------------|
| `name` | `string` | Name of the event which will be all uppercase. |
| `data` | `any?` | The data playload which can be any type of JSON object. |

## Source Types

Source types specify the type of source of a sound.

| Type | Description |
|------|-------------|
| `0` | local sound file |
| `1` | YouTube video ID |
| `2` | HTTP link |

---

## Commands

### INIT

This command must be send on start of each new web socket connection to authorize and identify this connection.

```json
{
    "name": "INIT",
    "data": {
        "user_id": "221905671296253953",
        "token": "gDURWm1gkLEjmFcjKs1CzWkUIkIDJQ486iheIfcr728jb6MxG2RUaoLnTdCILxLJ"
    }
}
```

#### Data

The delivered data consists of following JSON object:

| Field | Type | Description |
|-------|------|-------------|
| `user_id` | `string` | Your Discord user ID. |
| `token` | `string` | Your API token. |


### JOIN

Join the current voice channel you are connected to.

```json
{
    "name": "JOIN"
}
```


### LEAVE

Leave the voice channel on your guild where you are also in.

```json
{
    "name": "LEAVE"
}
```


### PLAY

Play a local sound or a resource from an online resoucre like YouTube or via HTTP link.

```json
{
    "name": "PLAY",
    "data": {
        "ident": "yay",
        "source": 0
    }
}
```

#### Data

The delivered data consists of following JSON object:

| Field | Type | Description |
|-------|------|-------------|
| `ident` | `string` | When playing a local sound file, this will be the name of the sound filw without file extension. Else, this is the specific resource identifier *(YouTube Video ID, HTTP link...)*. |
| `source` | `int` | The source type. |


### RANDOM

This is shorthand for `PLAY` command which plays a random picked, local sound.

```json
{
    "name": "RANDOM"
}
```


### STOP

Stop a currently playing sound.

```json
{
    "name": "STOP"
}
```


### VOLUME

Set the volume of a player in a voice channel. The set volume will be saved and also applied to future players created and used on the guild.

```json
{
    "name": "VOLUME",
    "data": 150
}
```

#### Data

The delivered data specifies the `volume` as `integer` value in a range of `[0, 1000]`. This value represents the volume in `%`.

---

## Events

### HELLO

This event is fired when you successfully authenticated and initialized the WS connection with the [`INIT`](#init) command. You should wait for this event until sending further commands.

```json
{
    "name": "HELLO"
}
```

### ERROR

This event is fired every time a command could not be executed or something other unexpected exceptions occure on the server side.

```json
{
  "name": "ERROR",
  "data": {
    "code": 0,
    "type": "bad command args",
    "message": "ident must be a valid string value"
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `code` | `int` | The integer code of the error type. |
| `type` | `string` | The description of the error type. |
| `message` | `string` | Further information about the error. |


### PLAYING

Is fired when a track starts playing.

```json
{
  "name": "PLAYING",
  "data": {
     "ident": "danke",
     "source": 0,
     "guild_id": "526196711962705925",
     "channel_id": "549871583364382771",
     "user_id": "221905671296253953",
     "user_tag": "zekro#9131"
  }
}
```

#### Data

The event data is a JSON object with following properties:

| Field | Type | Description |
|-------|------|-------------|
| `ident` | `string` | When playing a local sound file, this will be the name of the sound filw without file extension. Else, this is the specific resource identifier *(YouTube Video ID, HTTP link...)*. |
| `source` | `int` | The source type. |
| `guild_id` | `string` | The guild the track was played on. |
| `user_id` | `string` | The users ID who played the track. |
| `user_tag` | `string` | The users Tag who played the track. |


### END

Is fired when a track stops playing.

```json
{
  "name": "END",
  "data": {
    "ident": "danke",
    "source": 0,
    "guild_id": "526196711962705925",
    "channel_id": "549871583364382771",
    "user_id": "221905671296253953",
    "user_tag": "zekro#9131"
  }
}
```

#### Data

The event data is a JSON object with following properties:

| Field | Type | Description |
|-------|------|-------------|
| `ident` | `string` | When playing a local sound file, this will be the name of the sound filw without file extension. Else, this is the specific resource identifier *(YouTube Video ID, HTTP link...)*. |
| `source` | `int` | The source type. |
| `guild_id` | `string` | The guild the track was played on. |
| `user_id` | `string` | The users ID who played the track. |
| `user_tag` | `string` | The users Tag who played the track. |


### PLAY_ERROR

Is fired when an unexpected exception occures while trying to play a sound.

```json
{
  "name": "PLAY_ERROR",
  "data": {
    "reason": "some error description here",
    "track": {
      "ident": "danke",
      "source": 0,
      "guild_id": "526196711962705925",
      "channel_id": "549871583364382771",
      "user_id": "221905671296253953",
      "user_tag": "zekro#9131"
    }
  }
}
```

#### Data

The event data is a JSON object with following properties:

| Field | Type | Description |
|-------|------|-------------|
| `reason` | `string` | Further information about th exception |
| `track.ident` | `string` | When playing a local sound file, this will be the name of the sound filw without file extension. Else, this is the specific resource identifier *(YouTube Video ID, HTTP link...)*. |
| `track.source` | `int` | The source type. |
| `track.guild_id` | `string` | The guild the track was played on. |
| `track.user_id` | `string` | The users ID who played the track. |
| `track.user_tag` | `string` | The users Tag who played the track. |


### STUCK

Actually, I have no Idea when this event is getting fired or what the values below mean, because [I was not able to find any documentation about this](https://github.com/Frederikam/Lavalink/blob/f88938819976c7973d38d5dabff777cd4faa5fcd/LavalinkServer/src/main/java/lavalink/server/player/EventEmitter.java#L86).

```json
{
  "name": "PLAY_ERROR",
  "data": {
    "threshold": 1337,
    "track": {
      "ident": "danke",
      "source": 0,
      "guild_id": "526196711962705925",
      "channel_id": "549871583364382771",
      "user_id": "221905671296253953",
      "user_tag": "zekro#9131"
    }
  }
}
```

#### Data

The event data is a JSON object with following properties:

| Field | Type | Description |
|-------|------|-------------|
| `threshold` | `ind` | No idea what this says but be lucky if you know... |
| `track.ident` | `string` | When playing a local sound file, this will be the name of the sound filw without file extension. Else, this is the specific resource identifier *(YouTube Video ID, HTTP link...)*. |
| `track.source` | `int` | The source type. |
| `track.guild_id` | `string` | The guild the track was played on. |
| `track.user_id` | `string` | The users ID who played the track. |
| `track.user_tag` | `string` | The users Tag who played the track. |


### VOLUME_CHANGED

Is fired when the volume for the guild was changed.

```json
{
  "name": "VOLUME_CHANGED",
  "data": {
    "guild_id": "526196711962705925",
    "vol": 150
  }
}
```

#### Data

The event data is a JSON object with following properties:

| Field | Type | Description |
|-------|------|-------------|
| `guild_id` | `string` | The ID of the guild where the volume was changed. |
| `vol` | `int` | The new volume value in `%`. |


### JOINED

Is fired when the bot joined a voice channel.

```json
{
  "name": "JOINED",
  "data": {
    "guild_id": "526196711962705925",
    "channel_id": "549871583364382771"
  }
}
```

#### Data

The event data is a JSON object with following properties:

| Field | Type | Description |
|-------|------|-------------|
| `guild_id` | `string` | The ID of the guild where the bot joined. |
| `channel_id` | `string` | The ID of the channel where the bot joined into. |


### LEFT

Is fired when the bot joined a voice channel.

```json
{
  "name": "LEFT",
  "data": {
    "guild_id": "526196711962705925",
    "channel_id": "549871583364382771"
  }
}
```

#### Data

The event data is a JSON object with following properties:

| Field | Type | Description |
|-------|------|-------------|
| `guild_id` | `string` | The ID of the guild where the bot has quitted. |
| `channel_id` | `string` | The ID of the channel which the bot has quitted. |
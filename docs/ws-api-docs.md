# WS API

The web socket API is used for controlling players and receiving player events.

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

![](https://img.shields.io/badge/-fully%20implemented-green.svg)

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

![](https://img.shields.io/badge/-not%20implemented%20yet-red.svg)

Join the current voice channel you are connected to.

```json
{
    "name": "JOIN"
}
```


### LEAVE

![](https://img.shields.io/badge/-not%20implemented%20yet-red.svg)

Leave the voice channel on your guild where you are also in.

```json
{
    "name": "LEAVE"
}
```


### PLAY

![](https://img.shields.io/badge/-fully%20implemented-green.svg)

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

![](https://img.shields.io/badge/-not%20implemented%20yet-red.svg)

This is shorthand for `PLAY` command which plays a random picked, local sound.

```json
{
    "name": "RANDOM"
}
```

### VOLUME

![](https://img.shields.io/badge/-not%20implemented%20yet-red.svg)

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

### PLAYING

![](https://img.shields.io/badge/-fully%20implemented-green.svg)

Is fired when a track starts playing.

```json
{
    "name": "PLAYING",
    "data": {
        "ident": "danke",
        "source": 0,
        "guild_id": "526196711962705925",
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
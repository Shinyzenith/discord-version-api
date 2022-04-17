# Discord-version-API

## Dependencies:

1. `git`
1. `go`
1. `make`

## Building:

```bash
git clone https://github.com/shinyzenith/discord-version-api;cd discord-version-api/
make
```

## Usage:

Run the binary then use the `/version/:release_channel` endpoint to get the data.

The API is served on port 3000 by default. Use the `-p` flag to alter this behavior.

```bash
./discord-version-api
curl http://localhost:3000/version/canary -H "X-API-Key: <yourapikey>"
curl http://localhost:3000/buildid/stable -H "X-API-Key: <yourapikey>"
```

Rename `env.sample` to `.env` and add your API key.
The above method of authentication is really cringe but I need something useable right now so I'll switch to redis when I have time.

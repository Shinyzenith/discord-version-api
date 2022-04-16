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

The API is served on port 3000 by default.

```bash
./discord-version-api
curl http://localhost:3000/version/canary
curl http://localhost:3000/buildid/stable
```

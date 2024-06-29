# ESC Discord Bot

This is the official "Engineering Student Committee, Chulalongkorn University" Discord bot.

## Features

- [x] Welcome message
- [x] Webhook syncing users data from NocoDB
- [x] User authentication via Bot Direct Message
- [x] Role assignment via Bot Direct Message
- [ ] Link and PATCH user Discord ID to userss data in NocoDB

## Getting Started

### Download Modules

```bash
go mod download
```

### Configuration

Create a `.env` file in the root directory, example given in [`.env.example`](./.env.example).

For further configuration, [`config-local.yml`](./config/config-local.yml) is provided.

If webhook is in use, do port forwarding of port `8080` so that NocoDB can trigger the webhook.

### Run the Bot

```bash
go run ./cmd/bot/main.go

# or

air # for hot reload: https://github.com/cosmtrek/air
```

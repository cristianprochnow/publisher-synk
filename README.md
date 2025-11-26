# Publisher
Service to publish posts from system.

# Getting Started

First step is to create a `.env` file in project root and change example values to your config. You can use `example.env` file from `_setup` folder as template.

And then, run `docker compose up -d` into project root to start project.

## Tests

The easy way to run tests is just run `docker compose up -d` command to start project with variables. So, enter in `synk_publisher` with `docker exec` and run `go test ./tests -v`.

## Certificates

This app must run in HTTPS to authentication works properly. So, to install it, just setup `[mkcert](https://github.com/FiloSottile/mkcert)` into your machine and then run command below into root directory of this project.

```
mkcert -key-file ./.cert/key.pem -cert-file ./.cert/cert.pem localhost synk_publisher
```

## Network

You can use a custom network for this services, using then `synk_network` you must create before run it. So, to create on just run command below once during initial setup.

```
docker network create synk_network
```

# Setup integrations

## Discord

**Benefícios:**

- O Truque: Não crie um Bot. Use um Webhook.
- Por que é fácil: Você não precisa de autenticação OAuth, nem de tokens complexos. O Discord te dá uma URL única; qualquer JSON que você enviar para lá vira uma mensagem.

**Como fazer:**

1. Crie um servidor seu no Discord.
2. Vá nas configurações de um canal de texto -> Integrações -> Webhooks.
3. Crie um novo Webhook e copie a URL do Webhook.

**Exemplo:**

```go
type DiscordResponse struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
}

func PostToDiscord(content string) (string, error) {
	webhookURL := "YOUR_WEBHOOK_URL_HERE" + "?wait=true"

	payload := map[string]string{"content": content}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("discord returned status: %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)

	var result DiscordResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", err
	}

	fmt.Printf("Message posted successfully! ID: %s\n", result.ID)

	return result.ID, nil
}
```

## Telegram

**Benefícios:**

- Por que é fácil: Você fala com um bot chamado @BotFather, ele te dá um token e pronto. A API é uma URL simples.

**Como fazer:**

1. Abra o Telegram e procure por `@BotFather`.
2. Envie `/newbot`, dê um nome e um username.
3. Ele te dará um token (ex: `123456:ABC-DEF...`).
4. Você precisará saber o `chat_id` para onde enviar (mande uma mensagem para seu bot e acesse `https://api.telegram.org/bot<SEU_TOKEN>/getUpdates` para descobrir seu ID).

**Exemplo:**

```go
func PostToTelegram(content string) error {
    botToken := "SEU_TOKEN"
    chatID := "SEU_CHAT_ID"

    apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", botToken, chatID, url.QueryEscape(content))

    resp, err := http.Get(apiURL)
    if err != nil { return err }
    defer resp.Body.Close()

    return nil
}
```

# Routes

## Get info about app

> `GET` /about

### Response

```json
{
	"ok": true,
	"error": "",
	"info": {
		"server_port": "8080",
		"app_port": "8083",
		"db_working": true
	},
	"list": null
}
```

## Send message through Discord

> `GET` /discord/publish

### Request

> `webhook_url` can be get in the text chat settings on a server from Discord. So in that settings access Integrations > Webhooks > New Webhook.

```json
{
	"message": "showwwwwwwwwwwwww",
	"webhook_url": "https://discord.com/api/webhooks/123456789/asadsadfasdwefef323112ewefdwed"
}
```

### Response

```json
{
	"resource": {
		"ok": true,
		"error": ""
	},
	"post": {
		"id": "123456789",
		"channel_id": "123456789",
		"webhook_id": "123456789"
	},
	"raw": "{\"type\":0,\"content\":\"showwwwwwwwwwwwww\",\"mentions\":[],\"mention_roles\":[],\"attachments\":[],\"embeds\":[],\"timestamp\":\"2025-11-26T00:23:48.134000+00:00\",\"edited_timestamp\":null,\"flags\":0,\"components\":[],\"id\":\"123456789\",\"channel_id\":\"123456789\",\"author\":{\"id\":\"123456789\",\"username\":\"Captain Hook\",\"avatar\":null,\"discriminator\":\"0000\",\"public_flags\":0,\"flags\":0,\"bot\":true,\"global_name\":null,\"clan\":null,\"primary_guild\":null},\"pinned\":false,\"mention_everyone\":false,\"tts\":false,\"webhook_id\":\"123456789\"}\n"
}
```

## Send message through Telegram

> `GET` /telegram/publish

### Request

> `bot_token` and `chat_id` can be got following instructions above for Setup Integrations of Telegram.

```json
{
	"bot_token": "123456789:123456789-123456789-123456789",
	"chat_id": "123456789",
	"message": "showwwwwwww"
}
```

### Response

```json
{
	"resource": {
		"ok": true,
		"error": ""
	},
	"post": {
		"message_id": "4"
	},
	"raw": "{\"ok\":true,\"result\":{\"message_id\":4,\"from\":{\"id\":123456789,\"is_bot\":true,\"first_name\":\"Show\",\"username\":\"showbot\"},\"chat\":{\"id\":123456789,\"first_name\":\"Cristian\",\"last_name\":\"Prochnow\",\"type\":\"private\"},\"date\":1764178981,\"text\":\"showwwwwwww\"}}"
}
```
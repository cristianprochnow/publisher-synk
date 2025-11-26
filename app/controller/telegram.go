package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"synk/gateway/app/model"
)

type Telegram struct {
	model *model.Telegram
}

type HandleTelegramPublishResponse struct {
	Resource ResponseHeader                    `json:"resource"`
	Post     HandleTelegramPublishDataResponse `json:"post"`
	Raw      any                               `json:"raw"`
}

type HandleTelegramPublishDataResponse struct {
	Id        string `json:"id"`
	ChannelId string `json:"channel_id"`
	WebhookId string `json:"webhook_id"`
}

type HandleTelegramPublishRequest struct {
	Message    string `json:"message"`
	WebhookUrl string `json:"webhook_url"`
}

func NewTelegram(db *sql.DB) *Telegram {
	telegram := Telegram{
		model: model.NewTelegram(db),
	}

	return &telegram
}

func (d *Telegram) HandlePublish(w http.ResponseWriter, r *http.Request) {
	SetJsonContentType(w)

	response := HandleTelegramPublishResponse{
		Resource: ResponseHeader{
			Ok: true,
		},
		Post: HandleTelegramPublishDataResponse{},
	}

	bodyContent, bodyErr := io.ReadAll(r.Body)

	if bodyErr != nil {
		response.Resource.Ok = false
		response.Resource.Error = "error on read message body"

		WriteErrorResponse(w, response, "/telegram/publish", response.Resource.Error, http.StatusBadRequest)

		return
	}

	var post HandleTelegramPublishRequest

	jsonErr := json.Unmarshal(bodyContent, &post)

	if jsonErr != nil {
		response.Resource.Ok = false
		response.Resource.Error = "some fields can be in invalid format"

		WriteErrorResponse(w, response, "/telegram/publish", response.Resource.Error, http.StatusBadRequest)

		return
	}

	post.Message = strings.TrimSpace(post.Message)
	post.WebhookUrl = strings.TrimSpace(post.WebhookUrl)

	if post.Message == "" || post.WebhookUrl == "" {
		response.Resource.Ok = false
		response.Resource.Error = "field `message` and `webhook_url` can not be empty"

		WriteErrorResponse(w, response, "/telegram/publish", response.Resource.Error, http.StatusBadRequest)

		return
	}

	post.WebhookUrl += "?wait=true" // To return response body

	payload := map[string]string{"content": post.Message}
	jsonPayload, jsonPayloadErr := json.Marshal(payload)

	if jsonPayloadErr != nil {
		response.Resource.Ok = false
		response.Resource.Error = "some fields can be in invalid format on sending message"

		WriteErrorResponse(w, response, "/telegram/publish", response.Resource.Error, http.StatusBadRequest)

		return
	}

	respMessage, errMessage := http.Post(post.WebhookUrl, "application/json", bytes.NewBuffer(jsonPayload))
	if errMessage != nil {
		response.Resource.Ok = false
		response.Resource.Error = errMessage.Error()

		WriteErrorResponse(w, response, "/telegram/publish", response.Resource.Error, http.StatusBadRequest)

		return
	}

	defer respMessage.Body.Close()

	bodyBytes, _ := io.ReadAll(respMessage.Body)

	response.Raw = string(bodyBytes)

	json.Unmarshal(bodyBytes, &response.Post)

	WriteSuccessResponse(w, response)
}

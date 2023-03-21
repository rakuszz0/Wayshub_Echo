package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	channeldto "wayshub/dto/channel"
	dto "wayshub/dto/result"
	"wayshub/models"
	"wayshub/pkg/bcrypt"
	"wayshub/repositories"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// var os.Getenv("PATH_FILE") = "http://localhost:8080/uploads/"

type handlerChannel struct {
	ChannelRepository repositories.ChannelRepository
}

func HandlerChannel(ChannelRepository repositories.ChannelRepository) *handlerChannel {
	return &handlerChannel{ChannelRepository}
}

func (h *handlerChannel) FindChannels(c echo.Context) error {
	channels, err := h.ChannelRepository.FindChannels()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	for i, p := range channels {
		channels[i].Cover = os.Getenv("PATH_FILE") + p.Cover
	}

	for i, p := range channels {
		channels[i].Photo = os.Getenv("PATH_FILE") + p.Photo
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: "success", Data: channels})
}

func (h *handlerChannel) GetChannel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	channel, err := h.ChannelRepository.GetChannel(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	channel.Cover = os.Getenv("PATH_FILE") + channel.Cover
	channel.Photo = os.Getenv("PATH_FILE") + channel.Photo

	for i, p := range channel.Video {
		channel.Video[i].Thumbnail = os.Getenv("PATH_FILE") + p.Thumbnail
	}

	for i, p := range channel.Video {
		channel.Video[i].Video = os.Getenv("PATH_FILE") + p.Video
	}

	for i, p := range channel.Subscription {
		channel.Subscription[i].OtherPhoto = os.Getenv("PATH_FILE") + p.OtherPhoto
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: "success", Data: convertResponse(channel)})
}

func (h *handlerChannel) UpdateChannel(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(c.Param("id"))

	channelInfo := c.Get("channelInfo").(jwt.MapClaims)
	channelID := int(channelInfo["id"].(float64))

	if channelID != id {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Code: http.StatusUnauthorized, Message: "Can't update channel!"})
	}

	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	dataCover := c.Get("dataCover")
	fileCover := dataCover.(string)

	dataPhoto := c.Get("dataPhoto")
	filePhoto := dataPhoto.(string)

	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	resp1, err := cld.Upload.Upload(ctx, fileCover, uploader.UploadParams{Folder: "WaysHub"})
	if err != nil {
		fmt.Println(err.Error())
	}

	resp2, err := cld.Upload.Upload(ctx, filePhoto, uploader.UploadParams{Folder: "WaysHub"})
	if err != nil {
		fmt.Println(err.Error())
	}

	var request channeldto.UpdateChannelRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	password, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	channel, err := h.ChannelRepository.GetChannel(int(id))
	if err != nil {
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		return c.JSON(http.StatusBadRequest, response)
	}
	if request.Email != "" {
		channel.Email = request.Email
	}

	if request.Password != "" {
		channel.Password = password
	}

	if request.ChannelName != "" {
		channel.ChannelName = request.ChannelName
	}

	if request.Cover != "false" {
		channel.Cover = resp1.SecureURL
	}

	if request.Photo != "false" {
		channel.Photo = resp2.SecureURL
	}

	if request.Description != "" {
		channel.Description = request.Description
	}

	data, err := h.ChannelRepository.UpdateChannel(channel)
	if err != nil {
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		return c.JSON(http.StatusInternalServerError, response)
	}

	channel.Cover = os.Getenv("PATH_FILE") + channel.Cover
	channel.Photo = os.Getenv("PATH_FILE") + channel.Photo

	response := dto.SuccessResult{Code: http.StatusOK, Data: data}
	return c.JSON(http.StatusOK, response)

}

func (h *handlerChannel) DeleteChannel(c echo.Context) error {
	channelInfo := c.Get("channelInfo").(jwt.MapClaims)
	channelID := int(channelInfo["id"].(float64))

	channel, err := h.ChannelRepository.GetChannel(channelID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.ChannelRepository.DeleteChannel(channel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: "success", Data: deleteResponse(data)})
}

func (h *handlerChannel) PlusSubscriber(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	channel, err := h.ChannelRepository.GetChannel(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	channel.Subscriber = channel.Subscriber + 1

	data, err := h.ChannelRepository.PlusSubscriber(channel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: "success", Data: data})
}

func (h *handlerChannel) MinusSubscriber(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	channel, err := h.ChannelRepository.GetChannel(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	channel.Subscriber = channel.Subscriber - 1

	data, err := h.ChannelRepository.MinusSubscriber(channel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: "success", Data: data})
}

func convertResponse(u models.Channel) channeldto.ChannelResponse {
	return channeldto.ChannelResponse{
		ID:           u.ID,
		Email:        u.Email,
		ChannelName:  u.ChannelName,
		Description:  u.Description,
		Cover:        u.Cover,
		Photo:        u.Photo,
		Video:        u.Video,
		Subscription: u.Subscription,
		Subscriber:   u.Subscriber,
	}
}

func updateResponse(u models.Channel) channeldto.ChannelResponse {
	return channeldto.ChannelResponse{
		ID:          u.ID,
		Email:       u.Email,
		ChannelName: u.ChannelName,
		Description: u.Description,
		Cover:       u.Cover,
		Photo:       u.Photo,
	}
}

func deleteResponse(u models.Channel) channeldto.DeleteResponse {
	return channeldto.DeleteResponse{
		ID: u.ID,
	}
}

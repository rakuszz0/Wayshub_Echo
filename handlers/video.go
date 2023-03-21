package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	dto "wayshub/dto/result"
	videodto "wayshub/dto/video"
	"wayshub/models"
	"wayshub/repositories"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerVideo struct {
	VideoRepository repositories.VideoRepository
}

func HandlerVideo(VideoRepository repositories.VideoRepository) *handlerVideo {
	return &handlerVideo{VideoRepository}
}

func (h *handlerVideo) FindVideos(c echo.Context) error {
	videos, err := h.VideoRepository.FindVideos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	for i, p := range videos {
		// videos[i].Thumbnail = path_file + p.Thumbnail
		thumbnailPath := os.Getenv("PATH_FILE") + p.Thumbnail
		videos[i].Thumbnail = thumbnailPath
	}

	for i, p := range videos {
		// videos[i].Thumbnail = path_file + p.Thumbnail
		videoPath := os.Getenv("PATH_FILE") + p.Video
		videos[i].Video = videoPath
	}

	response := dto.SuccessResult{Status: "success", Data: videos}
	return c.JSON(http.StatusOK, response)
}

func (h *handlerVideo) GetVideo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var video models.Video
	video, err := h.VideoRepository.GetVideo(id)
	if err != nil {
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		return c.JSON(http.StatusBadRequest, response)
	}

	video.Thumbnail = os.Getenv("PATH_FILE") + video.Thumbnail
	video.Video = os.Getenv("PATH_FILE") + video.Video

	response := dto.SuccessResult{Status: "success", Data: subscribe}
	return c.JSON(http.StatusOK, response)
}

func (h *handlerVideo) CreateVideo(c echo.Context) error {
	channelInfo := c.Get("channelInfo").(jwt.MapClaims)
	channelID := int(channelInfo["id"].(float64))

	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	dataThumbnail := c.Get("dataThumbnail")
	fileThumbnail := dataThumbnail.(string)

	dataVideo := c.Get("dataVideo")
	fileVideo := dataVideo.(string)

	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	resp1, err := cld.Upload.Upload(ctx, fileThumbnail, uploader.UploadParams{Folder: "wayshub"})
	if err != nil {
		fmt.Println(err.Error())
	}

	resp2, err := cld.Upload.Upload(ctx, fileVideo, uploader.UploadParams{Folder: "wayshub"})
	if err != nil {
		fmt.Println(err.Error())
	}

	request := videodto.CreateVideoRequest{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
	}

	validation := validator.New()
	err = validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	video := models.Video{
		Title:       request.Title,
		Thumbnail:   resp1.SecureURL,
		Description: request.Description,
		Video:       resp2.SecureURL,
		CreatedAt:   time.Now(),
		ChannelID:   channelID,
	}

	video.FormatTime = video.CreatedAt.Format("2 Jan 2006")

	video, err = h.VideoRepository.CreateVideo(video)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	video, _ = h.VideoRepository.GetVideo(video.ID)

	response := dto.SuccessResult{Status: "success", Data: video}
	return c.JSON(http.StatusOK, response)
}

func (h *handlerVideo) UpdateVideo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	channelInfo := c.Get("channelInfo").(jwt.MapClaims)
	channelID := int(channelInfo["id"].(float64))

	if channelID != id {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Code: http.StatusUnauthorized, Message: "Can't update channel!"})
	}

	dataThumbnail := c.Get("dataThumbnail")
	fileThumbnail := dataThumbnail.(string)

	dataVideo := c.Get("dataVideo")
	fileVideo := dataVideo.(string)

	request := videodto.UpdateVideoRequest{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
	}

	video, err := h.VideoRepository.GetVideo(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if request.Title != "" {
		video.Title = request.Title
	}

	if request.Description != "" {
		video.Description = request.Description
	}

	if fileThumbnail != "false" {
		video.Thumbnail = fileThumbnail
	}

	if fileVideo != "false" {
		video.Video = fileVideo
	}

	data, err := h.VideoRepository.UpdateVideo(video)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	video.Thumbnail = os.Getenv("PATH_FILE") + video.Thumbnail
	video.Video = os.Getenv("PATH_FILE") + video.Video

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: "success", Data: data})
}

func (h *handlerVideo) DeleteVideo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	video, err := h.VideoRepository.GetVideo(int(id))

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.VideoRepository.DeleteVideo(video)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: "success", Data: DeleteVideoResponse(data)})
}

func (h *handlerVideo) FindVideosByChannelId(c echo.Context) error {
	channelInfo := c.Get("channelInfo").(jwt.MapClaims)
	channelID := int(channelInfo["id"].(float64))

	videos, err := h.VideoRepository.FindVideosByChannelId(channelID)
	if err != nil {
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		return c.JSON(http.StatusInternalServerError, response)
	}

	videos.Thumbnail = os.Getenv("PATH_FILE") + videos.Thumbnail
	videos.Video = os.Getenv("PATH_FILE") + videos.Video

	response := dto.SuccessResult{Status: "success", Data: videos}
	return c.JSON(http.StatusOK, response)
}

func (h *handlerVideo) FindMyVideos(c echo.Context) error {
	channelInfo := c.Get("channelInfo").(jwt.MapClaims)
	channelID := int(channelInfo["id"].(float64))

	videos, err := h.VideoRepository.FindMyVideos(channelID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	for i, p := range videos {
		videos[i].Thumbnail = os.Getenv("PATH_FILE") + p.Thumbnail
	}

	for i, p := range videos {
		videos[i].Video = os.Getenv("PATH_FILE") + p.Video
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: "success", Data: videos})
}

func (h *handlerVideo) UpdateViews(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	video, err := h.VideoRepository.GetVideo(int(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	video.ViewCount = video.ViewCount + 1

	data, err := h.VideoRepository.UpdateViews(video)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: "success", Data: data})
}

func DeleteVideoResponse(u models.Video) videodto.DeleteResponse {
	return videodto.DeleteResponse{
		ID: u.ID,
	}
}

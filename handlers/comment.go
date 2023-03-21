package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	commentdto "wayshub/dto/comment"
	dto "wayshub/dto/result"
	"wayshub/models"
	"wayshub/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerComment struct {
	CommentRepository respositories.CommentRepository
}

func HandlerComment(CommentRepository repositories.CommentRepository) *handlerComment {
	return &handlerComment{CommentRepository}
}

func (h *handlerComment) FindComments(c echo.Context) error {
	comments, err := h.CommentRepository.FindComments()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: comment})
}

func (h *handlerComment) GetComment(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	comment, err := h.CommentRepository.GetComment(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: comment})
}

func (h *handlerComment) CreateComment(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	channelInfo := c.Get("channelInfo").(jwt.MapClaims)
	channelID := int(channelInfo["id"].(float64))

	request := commentdto.CreateCommentRequest{
		Comment: c.FormValue("comment"),
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	comment := models.Comments{
		Comment:   request.Comment,
		ChannelID: channelID,
		VideoID:   id,
	}

	comment, err = h.CommentRepository.CreateComment(comment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	comment, _ = h.CommentRepository.GetComment(comment.ID)

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: "success", Data: comment})
}

func (h *handlerComment) UpdateComment(c echo.Context) error {
	idStr := c.Param("id")
	idInt, _ := strconv.Atoi(idStr)

	channelInfo := c.Get("channelInfo").(jwt.MapClaims)
	channelID := int(channelInfo["id"].(float64))

	if channelID != idInt {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Code: http.StatusUnauthorized, Message: "Can't update channel!"})
	}

	request := commentdto.CreateCommentRequest{
		Comment: c.FormValue("comment"),
	}

	comment, err := h.CommentRepository.GetComment(int(idInt))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if request.Comment != "" {
		comment.Comment = request.Comment
	}

	data, err := h.CommentRepository.UpdateComment(comment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: "success", Data: data})
}

func (h *handlerComment) DeleteComment(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	channelInfo := c.Get("channelInfo").(jwt.MapClaims)
	channelID := int(channelInfo["id"].(float64))

	comment, err := h.CommentRepository.GetComment(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if channelID != comment.ChannelID {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Code: http.StatusUnauthorized, Message: "Please Login First!"})
	}

	data, err := h.CommentRepository.DeleteComment(comment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: "success", Data: DeleteCommentResponse(data)})
}

func DeleteCommentResponse(u models.Comments) commentdto.DeleteResponse {
	return commentdto.DeleteResponse{
		ID: u.ID,
	}
}

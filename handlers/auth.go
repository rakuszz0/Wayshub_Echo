package handlers

import (
	"time"
	"wayhub/repositories"
	authdto "wayshub/dto/auth"
	dto "wayshub/dto/result"
	"wayshub/models"
)

type handlersAuth struct {
	AuthRepository repositories.AuthRepository
}

func HandlerAuth(AuthRepository repositories.AuthRepository) *handlerAuth {
	return &handlerAuth{AuthRepository}
}

func (h *handlerAuth) Register(c echo.Context) error {
	request := new(authdto.RegisterRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	validation := validator.New()
	err := validation.struct(request)
	if err := JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message:err.Error()})

	password, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code : http.StatusBadRequest, Message: err.Error()})
	}

	channel := models.Channel{
		Email:       request.Email,
		Password:    password,
		ChannelName: request.ChannelName,
		Description: request.Description,
	}

	data, err := h.AuthRepository.Register(channel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.ErrorResult{Code : StatusOK,Data: data})
}

func (h *handlerAuth) Login(c echo.Context) error {
    request := new(authdto.LoginRequest)
    if err := c.Bind(request); err != nil {
        return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
    }

    channel := models.Channel{
        Email:    request.Email,
        Password: request.Password,
    }

    channel, err := h.AuthRepository.Login(channel.Email)
    if err != nil {
        return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: "Wrong email or password!"})
    }

    isValid := bcrypt.CheckPasswordHash(request.Password, channel.Password)
    if !isValid {
        return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: "Wrong email or password!"})
    }

    claims := jwt.MapClaims{}
    claims["id"] = channel.ID
    claims["photo"] = channel.Photo
    claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

    token, errGenerateToken := jwtToken.GenerateToken(&claims)
    if errGenerateToken != nil {
        log.Println(errGenerateToken)
        fmt.Println("Unauthorize")
        return c.NoContent(http.StatusUnauthorized)
    }

    loginResponse := authdto.LoginResponse{
        ID:    channel.ID,
        Email: channel.Email,
        Photo: channel.Photo,
        Token: token,
    }

    channel.Photo = os.Getenv("PATH_FILE") + channel.Photo

    response := dto.SuccessResult{Status: "success", Data: loginResponse}
    return c.JSON(http.StatusOK, response)
}

func (h *handlerAuth) CheckAuth(c echo.Context) error {
	channelInfo := c.Get("channelInfo").(jwt.MapClaims)
	channelId := int(channelInfo["id"].(float64))

	channel, err := h.AuthRepository.Getchannel(channelId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	CheckAuthResponse := authdto.CheckAuthResponse{
		ID:    channel.ID,
		Email: channel.Email,
		Photo: channel.Photo,
	}

	channel.Photo = os.Getenv("PATH_FILE") + channel.Photo

	response := dto.SuccessResult{Status: "success", Data: CheckAuthResponse}
	return c.JSON(http.StatusOK, response)
}

func registerResponse(u models.Channel) authdto.RegisterResponse {
	return authdto.RegisterResponse{
		Email: u.Email,
	}
}
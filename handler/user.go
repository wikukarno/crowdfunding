package handler

import (
	"fmt"
	"net/http"

	"backend-crowdfunding/auth"
	"backend-crowdfunding/helper"
	"backend-crowdfunding/storage"
	"backend-crowdfunding/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService user.Service
	authService auth.Service
	uploader    storage.Uploader
}

func NewUserHandler(userService user.Service, authService auth.Service, uploader storage.Uploader) *UserHandler {
	return &UserHandler{userService, authService, uploader}
}

// @Summary Register a new user
// @Tags    users
// @Param   payload body user.RegisterUserInput true "Registration details"
// @Success 200 {object} helper.Response
// @Router  /users [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		c.Error(err)
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

// @Summary Authenticate and receive a JWT
// @Tags    users
// @Param   payload body user.LoginInput true "Login credentials"
// @Success 200 {object} helper.Response
// @Router  /sessions [post]
func (h *UserHandler) Login(c *gin.Context) {
	var input user.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, err := h.userService.Login(input)
	if err != nil {
		// Keep the client message generic so we don't reveal whether the email
		// exists or what the password check returned, but log the real cause.
		c.Error(err)
		errorMessage := gin.H{"errors": "invalid email or password"}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)
	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, token)
	response := helper.APIResponse("Successfully logged in", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

// @Summary Check whether an email is still available
// @Tags    users
// @Param   payload body user.CheckEmailInput true "Email to check"
// @Success 200 {object} helper.Response
// @Router  /email_checkers [post]
func (h *UserHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}
		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server Error"}
		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	metaMessage := "Email has been registered"
	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	data := gin.H{"is_available": isEmailAvailable}
	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

// @Summary  Upload the current user's avatar
// @Tags     users
// @Param    avatar formData file true "Avatar image"
// @Success  200 {object} helper.Response
// @Security BearerAuth
// @Router   /avatars [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	opened, err := file.Open()
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	defer opened.Close()

	key := objectKey(fmt.Sprintf("avatars/%s", userID), file.Filename)
	url, err := h.uploader.Upload(c.Request.Context(), key, opened, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		c.Error(err)
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if _, err := h.userService.SaveAvatar(userID, url); err != nil {
		c.Error(err)
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfully uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

// @Summary  Fetch the current authenticated user
// @Tags     users
// @Success  200 {object} helper.Response
// @Security BearerAuth
// @Router   /users/fetch [get]
func (h *UserHandler) FetchUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	formatter := user.FormatUser(currentUser, "")
	response := helper.APIResponse("Successfully fetch user data", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

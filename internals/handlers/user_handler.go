package handlers

import (
	"IFEST/helpers"
	"IFEST/internals/config"
	"IFEST/internals/core/domain"
	"IFEST/internals/services"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type UserHandler struct {
	userService services.IUserService
	validate    *validator.Validate
}

func NewUserHandler(userService services.IUserService) UserHandler {
	return UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

func (uh *UserHandler) Create(c *fiber.Ctx) error {
	var userRequest domain.UserRequest
	if err := c.BodyParser(&userRequest); err != nil {
		return helpers.HttpBadRequest(c, err.Error())
	}

	if err := uh.validate.Struct(userRequest); err != nil {
		var errors []string
		for _, errs := range err.(validator.ValidationErrors) {
			errors = append(errors, helpers.FormatValidationError(errs))
		}
		return helpers.HttpBadRequest(c, "failed to create user", errors)
	}

	user, err := uh.userService.Create(&userRequest, false)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to create user", err)
	}

	return helpers.HttpSuccess(c, "new user created", 201, user)
}

func (uh *UserHandler) Login(c *fiber.Ctx) error {
	var user domain.UserLogin
	if err := c.BodyParser(&user); err != nil {
		return helpers.HttpBadRequest(c, err.Error())
	}
	if err := uh.validate.Struct(user); err != nil {
		var errors []string
		for _, errs := range err.(validator.ValidationErrors) {
			errors = append(errors, helpers.FormatValidationError(errs))
			return helpers.HttpBadRequest(c, "failed to login user", errors)
		}
	}

	userData, err := uh.userService.GetByEmail(user.Email)
	if err != nil {
		return helpers.HttpNotFound(c, "invalid email / password")
	}

	if userData.IsGoogleAuth == true {
		return helpers.HttpBadRequest(c, "login using google Oauth instead")
	}

	if err := helpers.CompareHashAndPassword(userData.Password, user.Password); err != nil {
		return helpers.HttpBadRequest(c, "invalid email / password")
	}
	uuidStr := userData.ID.String()
	token, err := helpers.JwtToken(uuidStr)
	return helpers.HttpSuccess(c, "Login succes", 200, map[string]string{"token": token})
}

func (uh *UserHandler) Profile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	userData, err := uh.userService.GetByID(userID)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to get user", err)
	}

	return helpers.HttpSuccess(c, "succes to get user", 200, userData)
}

func (uh *UserHandler) GoogleLogin(c *fiber.Ctx) error {
	OAuth := config.OAuthConfig()

	state, err := helpers.GenerateState(16)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to get state", err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(5 * time.Minute),
		Secure:   false,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	url := OAuth.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func (uh *UserHandler) GoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	cookieState := c.Cookies("oauth_state")

	if state == "" || cookieState == "" {
		return helpers.HttpsInternalServerError(c, "failed to get state", errors.New("missing state value"))
	}

	if state != cookieState {
		return helpers.HttpUnauthorized(c, "invalid oauth state")
	}

	c.ClearCookie("oauth_state")

	code := c.Query("code")
	if code == "" {
		return helpers.HttpsInternalServerError(c, "missing code value", errors.New("missing code value"))
	}

	token, err := config.OAuthConfig().Exchange(c.Context(), code)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to exchange token", err)
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to get user info", err)
	}

	usr := domain.UserGoogleInfo

	if err := json.NewDecoder(response.Body).Decode(&usr); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to parse user info")
	}

	user, err := uh.userService.GetByEmail(usr.Email)
	if err != nil {
		user, err = uh.userService.Create(&domain.UserRequest{
			Name:     usr.Name,
			Email:    usr.Email,
			Password: usr.ID,
		}, true)
	}

	uuidStr := user.ID.String()
	jwt, err := helpers.JwtToken(uuidStr)
	return helpers.HttpSuccess(c, "Login succes", 200, map[string]string{"token": jwt})
}

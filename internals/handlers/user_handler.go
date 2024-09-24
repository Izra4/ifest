package handlers

import (
	"IFEST/helpers"
	"IFEST/internals/core/domain"
	"IFEST/internals/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

	user, err := uh.userService.Create(&userRequest)
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

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

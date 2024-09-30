package handlers

import (
	"IFEST/helpers"
	"IFEST/internals/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserDocHandler struct {
	userDocService services.IUserDocService
	docService     services.IDocsService
	userService    services.IUserService
}

func NewUserDocHandler(userDocService services.IUserDocService) UserDocHandler {
	return UserDocHandler{
		userDocService: userDocService,
	}
}

func (udh *UserDocHandler) Create(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	if userIDStr == "" {
		return helpers.HttpUnauthorized(c, "unauthorized")
	}

	email := c.FormValue("email")
	docsIDStr := c.Params("id")

	docsID, err := uuid.Parse(docsIDStr)

	userAdded, err := udh.userService.GetByEmail(email)
	if err != nil {
		return helpers.HttpNotFound(c, "user not found")
	}

	access, err := udh.userDocService.Create(userAdded.ID, docsID)

	return helpers.HttpSuccess(c, "success to give data", 201, access)
}

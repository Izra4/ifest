package handlers

import (
	"IFEST/helpers"
	"IFEST/helpers/email"
	"IFEST/internals/config"
	"IFEST/internals/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"time"
)

type UserDocHandler struct {
	userDocService services.IUserDocService
	userService    services.IUserService
	docService     services.IDocsService
}

func NewUserDocHandler(userDocService services.IUserDocService, userService services.IUserService, docService services.IDocsService) UserDocHandler {
	return UserDocHandler{
		userDocService: userDocService,
		userService:    userService,
		docService:     docService,
	}
}

func (udh *UserDocHandler) Create(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	if userIDStr == "" {
		return helpers.HttpUnauthorized(c, "unauthorized")
	}

	user, err := udh.userService.GetByID(userIDStr)
	if err != nil {
		return helpers.HttpBadRequest(c, "user not found", err)
	}

	emailInput := c.FormValue("email")
	docsIDStr := c.Params("id")

	docsID, err := uuid.Parse(docsIDStr)

	userAdded, err := udh.userService.GetByEmail(emailInput)
	if err != nil {
		return helpers.HttpNotFound(c, "user not found")
	}

	access, err := udh.userDocService.Create(userAdded.ID, docsID, emailInput, user.Name)

	return helpers.HttpSuccess(c, "success to create", 201, access)
}

func (udh *UserDocHandler) Download(c *fiber.Ctx) error {
	//userID := c.Locals("userID").(string)
	//if userID == "" {
	//	return helpers.HttpUnauthorized(c, "unauthorized")
	//}

	token := c.Query("token")
	if token == "" {
		return helpers.HttpBadRequest(c, "token is required", nil)
	}

	data, err := udh.userDocService.FindByToken(token)
	if err != nil {
		return helpers.HttpBadRequest(c, "invalid token / expired token", err)
	}

	log.Println(time.Now().UTC(), " || ", data.Expired_at)

	if time.Now().UTC().After(data.Expired_at) {
		if err := udh.userDocService.DeleteAccessByToken(token); err != nil {
			return helpers.HttpsInternalServerError(c, "failed to delete access", err)
		}
		return helpers.HttpUnauthorized(c, "token expired")
	}

	docs, err := udh.docService.FindByID(data.DocID.String())
	if err != nil {
		return helpers.HttpNotFound(c, "document not found")
	}

	encryptedFile, err := config.SupabaseClient().DownloadFile("docs", docs.DocumentName)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to download file from cloud storage", err)
	}

	decryptedData, err := helpers.Decrypt(encryptedFile)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to decrypt file", err)
	}

	c.Set("Content-Disposition", "attachment; filename="+docs.DocumentName)
	c.Set("Content-Type", "application/octet-stream")
	return c.Send(decryptedData)
}

func (udh *UserDocHandler) TestEmail(c *fiber.Ctx) error {
	email.SendDownloadLink("akunuplay7@gmail.com", "aman", "www.youtube.com")
	return helpers.HttpSuccess(c, "success", 200, nil)
}

package handlers

import (
	"IFEST/helpers"
	"IFEST/internals/config"
	"IFEST/internals/core/domain"
	"IFEST/internals/services"
	"bytes"
	"encoding/base64"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	storage_go "github.com/supabase-community/storage-go"
	"io"
	"log"
	"strings"
)

type DocHandler struct {
	docService  services.IDocsService
	userService services.IUserService
	validator   *validator.Validate
}

func NewDocHandler(docService services.IDocsService, userService services.IUserService) DocHandler {
	return DocHandler{
		docService:  docService,
		userService: userService,
		validator:   validator.New(),
	}
}

func (dh *DocHandler) Upload(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)

	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		return helpers.HttpBadRequest(c, "failed to upload file", nil)
	}
	number := c.FormValue("number")
	fileType := c.FormValue("type")

	if number == "" || fileType == "" {
		return helpers.HttpBadRequest(c, "fill the form", nil)
	}

	src, err := file.Open()
	if err != nil {
		return helpers.HttpBadRequest(c, "failed to open the file", nil)
	}

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to read the file", err)
	}

	encryptedFile, err := helpers.Encrypt(fileBytes)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to encrypt file", err)
	}

	content := file.Header.Get("Content-Type")

	fileName := helpers.GenerateRandomString(10)
	_, err = config.SupabaseClient().UploadFile("docs", fileName, bytes.NewReader(encryptedFile), storage_go.FileOptions{ContentType: &content})
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to upload file to cloud storage", err)
	}
	userID, _ := uuid.Parse(userIDStr)

	encryptedNumberBytes, err := helpers.Encrypt([]byte(number))
	encryptedNumber := base64.StdEncoding.EncodeToString(encryptedNumberBytes)

	docs := domain.DocsUpload{
		UserID: userID,
		Name:   fileName,
		Number: encryptedNumber,
		Type:   fileType,
	}
	data, err := dh.docService.Upload(docs)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to upload file", err)
	}

	return helpers.HttpSuccess(c, "file uploaded", 201, data)
}

func (dh *DocHandler) GetAll(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	log.Println(userID)
	result, err := dh.docService.FindByUserID(userID)
	if err != nil {
		return helpers.HttpNotFound(c, "docs not found")
	}

	return helpers.HttpSuccess(c, "success to get data", 200, result)
}

func (dh *DocHandler) GetByID(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	if userIDStr == "" {
		return helpers.HttpUnauthorized(c, "unauthorized")
	}
	docsID := c.Params("id")
	log.Println(docsID)
	docs, err := dh.docService.FindByID(docsID)
	if err != nil {
		return helpers.HttpNotFound(c, "docs not found")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to parse user id", err)
	}

	if docs.UserID != userID {
		log.Println(docs.UserID, " || ", userID)
		return helpers.HttpUnauthorized(c, "unauthorized")
	}

	emails := strings.Split(docs.AccessEmails, ", ")
	var filteredEmails []string
	if emails[0] != "" {
		for _, data := range emails {
			filteredEmails = append(filteredEmails, data)
		}
	}

	fixed := domain.DocumentAccessInfo{
		DocumentID:     docs.DocumentID,
		DocumentName:   docs.DocumentName,
		DocumentType:   docs.DocumentType,
		DocumentStatus: docs.DocumentStatus,
		AccessCount:    docs.AccessCount,
		FixedEmails:    filteredEmails,
	}

	return helpers.HttpSuccess(c, "success to get data", 200, fixed)
}

func (dh *DocHandler) Update(c *fiber.Ctx) error {

	userIDStr := c.Locals("userID").(string)

	user, err := dh.userService.GetByID(userIDStr)
	if err != nil {
		return helpers.HttpUnauthorized(c, "unauthorized")
	}

	if user.Role != "admin" {
		return helpers.HttpUnauthorized(c, "unauthorized: not an admin")
	}
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to parse id", err)
	}

	var req domain.UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.HttpBadRequest(c, "failed to bind request", err)
	}

	if err := dh.docService.UpdateStatus(id, req.Status); err != nil {
		return helpers.HttpsInternalServerError(c, "failed to update status", err)
	}

	return helpers.HttpSuccess(c, "succes to update", 200, nil)
}

func (dh *DocHandler) GetUnverifiedDocs(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	if userIDStr == "" {
		return helpers.HttpUnauthorized(c, "unauthorized")
	}

	user, err := dh.userService.GetByID(userIDStr)
	if err != nil {
		return helpers.HttpNotFound(c, "user not found")
	}

	if user.Role != "admin" {
		return helpers.HttpUnauthorized(c, "unauthorized: not an admin")
	}

	docs, err := dh.docService.GetAllDocsByStatus(0)
	if err != nil {
		return helpers.HttpNotFound(c, "docs not found")
	}

	var unverifiedDocs []domain.UnverifiedDocs
	for _, list := range docs {

		decodedNumber, err := base64.StdEncoding.DecodeString(list.Number)
		if err != nil {
			return helpers.HttpsInternalServerError(c, "Failed to decode", err)
		}

		user, err := dh.userService.GetByID(list.UserID.String())
		if err != nil {
			return helpers.HttpsInternalServerError(c, "failed to get user", err)
		}

		decryptedNumber, err := helpers.Decrypt(decodedNumber)

		unverifiedDocs = append(unverifiedDocs, domain.UnverifiedDocs{
			ID:     list.ID.String(),
			Name:   user.Name,
			Type:   list.Type,
			Number: string(decryptedNumber),
			Date:   list.CreatedAt,
		})
	}

	return helpers.HttpSuccess(c, "success to get data", 200, unverifiedDocs)
}

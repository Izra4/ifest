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
)

type DocHandler struct {
	docService services.IDocsService
	validator  *validator.Validate
}

func NewDocHandler(docService services.IDocsService) DocHandler {
	return DocHandler{
		docService: docService,
		validator:  validator.New(),
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

func (dh *DocHandler) Download(c *fiber.Ctx) error {
	docsID := c.Params("id")
	log.Println(docsID)
	docs, err := dh.docService.FindByID(docsID)
	if err != nil {
		return helpers.HttpNotFound(c, "docs not found")
	}

	encryptedFile, err := config.SupabaseClient().DownloadFile("docs", docs.Name)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to download file from cloud storage", err)
	}

	decryptedData, err := helpers.Decrypt(encryptedFile)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to decrypt file", err)
	}

	c.Set("Content-Disposition", "attachment; filename="+docs.Name)
	c.Set("Content-Type", "application/octet-stream")
	return c.Send(decryptedData)
}

func (dh *DocHandler) GetAll(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	result, err := dh.docService.FindByUserID(userID)
	if err != nil {
		return helpers.HttpNotFound(c, "docs not found")
	}

	return helpers.HttpSuccess(c, "success to get data", 200, result)
}

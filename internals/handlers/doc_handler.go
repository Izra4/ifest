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
	"strconv"
	"strings"
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
	userID := c.Locals("userID").(string)
	if userID == "" {
		return helpers.HttpUnauthorized(c, "unauthorized")
	}

	docIDStr := c.Params("id")
	if docIDStr == "" {
		return helpers.HttpBadRequest(c, "document id is required", nil)
	}

	doc, err := dh.docService.FindByID(docIDStr)
	if err != nil {
		return helpers.HttpNotFound(c, "document not found")
	}

	var data domain.DocsUpdateRequest
	if err := c.BodyParser(&data); err != nil {
		return helpers.HttpBadRequest(c, "failed to parse request body", err)
	}

	if err := dh.validator.Struct(data); err != nil {
		var errors []string
		for _, errs := range err.(validator.ValidationErrors) {
			errors = append(errors, helpers.FormatValidationError(errs))
		}
		return helpers.HttpBadRequest(c, "failed to binding request", errors)
	}
	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to parse document id", err)
	}

	if data.Name == "" {
		data.Name = doc.DocumentName
	}
	if data.Type == "" {
		data.Number = doc.DocumentType
	}

	parsedDocStatus := strconv.Itoa(doc.DocumentStatus)

	if data.Status == "" {
		data.Status = parsedDocStatus
	}

	if data.Number == "" {
		data.Number = doc.DocumentNumber
	} else {
		encrypt, err := helpers.Encrypt([]byte(data.Number))
		if err != nil {
			return helpers.HttpsInternalServerError(c, "failed to encrypt document number", err)
		}
		data.Number = base64.StdEncoding.EncodeToString(encrypt)
	}

	dataInput := domain.DocsUpdateRequest{
		ID:     doc.DocumentID,
		Name:   data.Name,
		Type:   data.Type,
		Status: data.Status,
		Number: data.Number,
	}

	err = dh.docService.Update(docID, dataInput)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to update document", err)
	}

	return helpers.HttpSuccess(c, "document updated successfully", 200, nil)
}

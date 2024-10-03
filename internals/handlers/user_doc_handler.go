package handlers

import (
	"IFEST/helpers"
	"IFEST/helpers/email"
	"IFEST/internals/blockchain"
	"IFEST/internals/config"
	"IFEST/internals/core/domain"
	"IFEST/internals/services"
	"encoding/base64"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

type UserDocHandler struct {
	userDocService  services.IUserDocService
	userService     services.IUserService
	docService      services.IDocsService
	blockchain      *blockchain.Blockchain
	mutex           sync.Mutex
	processedTokens map[string]time.Time
	tokenMutex      sync.Mutex
	tokenTTL        time.Duration
	validate        *validator.Validate
}

func NewUserDocHandler(userDocService services.IUserDocService, userService services.IUserService,
	docService services.IDocsService, bc *blockchain.Blockchain) UserDocHandler {
	return UserDocHandler{
		userDocService:  userDocService,
		userService:     userService,
		docService:      docService,
		blockchain:      bc,
		processedTokens: make(map[string]time.Time),
		tokenTTL:        10 * time.Second,
		validate:        validator.New(),
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

	udh.tokenMutex.Lock()
	lastProcessed, exists := udh.processedTokens[token]
	if !exists || time.Since(lastProcessed) > udh.tokenTTL {
		udh.processedTokens[token] = time.Now()

		for tok, ts := range udh.processedTokens {
			if time.Since(ts) > udh.tokenTTL {
				delete(udh.processedTokens, tok)
			}
		}
		udh.tokenMutex.Unlock()
		udh.mutex.Lock()
		err = udh.blockchain.AddBlock(blockchain.Transaction{
			OwnerID:    docs.UserID.String(),
			AccessorID: data.UserID.String(),
			DocID:      data.DocID.String(),
			AccessTime: time.Now().UTC(),
		})
		udh.mutex.Unlock()
		if err != nil {
			log.Println("failed to add block to blockchain:", err)
		}
	} else {
		udh.tokenMutex.Unlock()
		log.Println("duplicate access token")
	}

	c.Set("Content-Disposition", "attachment; filename="+docs.DocumentName)
	c.Set("Content-Type", "application/octet-stream")
	return c.Send(decryptedData)
}

func (udh *UserDocHandler) GetHistoryByUserID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	if userID == "" {
		return helpers.HttpBadRequest(c, "userID is required", nil)
	}

	var histories []domain.AccessHistory
	historyList := udh.blockchain.GetHistoryByUserID(userID)

	for _, history := range historyList {

		accessor, err := udh.userService.GetByID(history.AccessorID)
		if err != nil {
			return helpers.HttpNotFound(c, "user not found")
		}

		docs, err := udh.docService.FindByID(history.DocID)
		if err != nil {
			return helpers.HttpNotFound(c, "user not found")
		}
		decodedNumber, err := base64.StdEncoding.DecodeString(docs.DocumentNumber)
		if err != nil {
			return helpers.HttpsInternalServerError(c, "Failed to decode", err)
		}
		decryptedNumber, err := helpers.Decrypt(decodedNumber)
		if err != nil {
			return helpers.HttpsInternalServerError(c, "failed to decrypt document number", err)
		}

		histories = append(histories, domain.AccessHistory{
			AcessorID:    history.AccessorID,
			DocID:        history.DocID,
			AccessorName: accessor.Name,
			Type:         docs.DocumentType,
			Number:       string(decryptedNumber),
			AccessTime:   history.AccessTime,
		})
	}

	return helpers.HttpSuccess(c, "history retrieved successfully", 200, histories)
}

func (udh *UserDocHandler) DeleteAccess(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	if userID == "" {
		return helpers.HttpUnauthorized(c, "unauthorized")
	}

	var deleteRequest domain.AcessDeleteRequest
	if err := c.BodyParser(&deleteRequest); err != nil {
		return helpers.HttpBadRequest(c, err.Error())
	}

	if err := udh.validate.Struct(deleteRequest); err != nil {
		var errors []string
		for _, errs := range err.(validator.ValidationErrors) {
			errors = append(errors, helpers.FormatValidationError(errs))
		}
		log.Println(deleteRequest.DocID)
		log.Println(deleteRequest.AcessorID)
		return helpers.HttpBadRequest(c, "failed to binding request", errors)
	}

	parsedAcessorID, err := uuid.Parse(deleteRequest.AcessorID)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "Failed to delete access", err)
	}

	parsedDocID, err := uuid.Parse(deleteRequest.DocID)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "Failed to delete access", err)
	}

	err = udh.userDocService.DeleteAccessByUserID(parsedAcessorID, parsedDocID)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "Failed to delete access", err)
	}
	return helpers.HttpSuccess(c, "success to delete the access", 200, nil)
}

func (udh *UserDocHandler) TestEmail(c *fiber.Ctx) error {
	email.SendDownloadLink("akunuplay7@gmail.com", "aman", "www.youtube.com")
	return helpers.HttpSuccess(c, "success", 200, nil)
}

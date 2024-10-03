package handlers

import (
	"IFEST/helpers"
	"IFEST/internals/core/domain"
	"IFEST/internals/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
)

type ReportHandler struct {
	reportService services.IReportsService
	validator     *validator.Validate
}

func NewReportHandler(reportService services.IReportsService) ReportHandler {
	return ReportHandler{
		reportService: reportService,
		validator:     validator.New(),
	}
}

func (rh *ReportHandler) CreateReport(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	if userIDStr == "" {
		return helpers.HttpUnauthorized(c, "unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to parse user ID", err)
	}

	var req domain.ReportCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.HttpBadRequest(c, "failed to parse request body", err)
	}

	if err := rh.validator.Struct(&req); err != nil {
		var errors []string
		for _, errs := range err.(validator.ValidationErrors) {
			errors = append(errors, helpers.FormatValidationError(errs))
		}
		return helpers.HttpBadRequest(c, "validation failed", errors)
	}

	req.UserID = userID

	report, err := rh.reportService.CreateReport(req)
	if err != nil {
		log.Println("Error creating report:", err)
		return helpers.HttpsInternalServerError(c, "failed to create report", err)
	}

	return helpers.HttpSuccess(c, "report created successfully", 201, report)
}

func (rh *ReportHandler) GetReportByID(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	if userIDStr == "" {
		return helpers.HttpUnauthorized(c, "unauthorized")
	}

	reportIDStr := c.Params("id")
	if reportIDStr == "" {
		return helpers.HttpBadRequest(c, "report ID is required", nil)
	}

	reportID, err := uuid.Parse(reportIDStr)
	if err != nil {
		return helpers.HttpBadRequest(c, "invalid report ID", err)
	}

	report, err := rh.reportService.GetReportByID(reportID)
	if err != nil {
		log.Println("Error fetching report:", err)
		return helpers.HttpNotFound(c, "report not found")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to parse user ID", err)
	}

	if report.UserID != userID {
		return helpers.HttpUnauthorized(c, "unauthorized to access this report")
	}

	return helpers.HttpSuccess(c, "report fetched successfully", 200, report)
}

func (rh *ReportHandler) GetReportsByUserID(c *fiber.Ctx) error {
	reports, err := rh.reportService.GetReports()
	if err != nil {
		log.Println("Error fetching reports:", err)
		return helpers.HttpNotFound(c, "reports not found")
	}

	return helpers.HttpSuccess(c, "reports fetched successfully", 200, reports)
}

func (rh *ReportHandler) UpdateReport(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	if userIDStr == "" {
		return helpers.HttpUnauthorized(c, "unauthorized")
	}

	reportIDStr := c.Params("id")
	if reportIDStr == "" {
		return helpers.HttpBadRequest(c, "report ID is required", nil)
	}

	reportID, err := uuid.Parse(reportIDStr)
	if err != nil {
		return helpers.HttpBadRequest(c, "invalid report ID", err)
	}

	report, err := rh.reportService.GetReportByID(reportID)
	if err != nil {
		return helpers.HttpNotFound(c, "report not found")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to parse user ID", err)
	}

	if report.UserID != userID {
		return helpers.HttpUnauthorized(c, "unauthorized to update this report")
	}

	var req domain.ReportUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.HttpBadRequest(c, "failed to parse request body", err)
	}

	if err := rh.validator.Struct(&req); err != nil {
		var errors []string
		for _, errs := range err.(validator.ValidationErrors) {
			errors = append(errors, helpers.FormatValidationError(errs))
		}
		return helpers.HttpBadRequest(c, "validation failed", errors)
	}

	updatedReport, err := rh.reportService.UpdateReport(reportID, req)
	if err != nil {
		log.Println("Error updating report:", err)
		return helpers.HttpsInternalServerError(c, "failed to update report", err)
	}

	return helpers.HttpSuccess(c, "report updated successfully", 200, updatedReport)
}

func (rh *ReportHandler) DeleteReport(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	if userIDStr == "" {
		return helpers.HttpUnauthorized(c, "unauthorized")
	}

	reportIDStr := c.Params("id")
	if reportIDStr == "" {
		return helpers.HttpBadRequest(c, "report ID is required", nil)
	}

	reportID, err := uuid.Parse(reportIDStr)
	if err != nil {
		return helpers.HttpBadRequest(c, "invalid report ID", err)
	}

	report, err := rh.reportService.GetReportByID(reportID)
	if err != nil {
		return helpers.HttpNotFound(c, "report not found")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return helpers.HttpsInternalServerError(c, "failed to parse user ID", err)
	}

	if report.UserID != userID {
		return helpers.HttpUnauthorized(c, "unauthorized to delete this report")
	}

	err = rh.reportService.DeleteReport(reportID)
	if err != nil {
		log.Println("Error deleting report:", err)
		return helpers.HttpsInternalServerError(c, "failed to delete report", err)
	}

	return helpers.HttpSuccess(c, "report deleted successfully", 200, nil)
}

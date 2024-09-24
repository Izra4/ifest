package helpers

import "github.com/gofiber/fiber/v2"

func HttpsInternalServerError(c *fiber.Ctx, msg string, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  "failed",
		"message": msg,
		"error":   err.Error(),
	})
}

func HttpSuccess(c *fiber.Ctx, msg string, code int, data any) error {
	return c.Status(code).JSON(fiber.Map{
		"status":  "success",
		"message": msg,
		"data":    data,
	})
}

func HttpBadRequest(c *fiber.Ctx, msg string, data ...any) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status":  "failed",
		"message": msg,
		"data":    data,
	})
}

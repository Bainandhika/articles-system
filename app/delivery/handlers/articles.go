package handlers

import (
	"articles-system/app/logging"
	"articles-system/app/services"
	"articles-system/lib/models"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ArticlesHandler struct {
	validator *validator.Validate
	service   services.ArticlesService
}

func NewArticlesHandler(validator *validator.Validate, service services.ArticlesService) *ArticlesHandler {
	return &ArticlesHandler{
		validator: validator,
		service:   service,
	}
}

func (h *ArticlesHandler) Create(c *fiber.Ctx) error {
	var payload models.AddArticle
	if err := c.BodyParser(&payload); err != nil {
		response := models.Response{
			Code:    fiber.StatusBadRequest,
			Message: fiber.ErrBadRequest.Message,
			Data:    err.Error(),
		}

		logging.Error.Println(logDetail(payload, response, err.Error()))
		return c.Status(response.Code).JSON(response)
	}

	if err := h.validator.Struct(payload); err != nil {
		response := models.Response{
			Code:    fiber.StatusBadRequest,
			Message: "Validation failed",
		}

		logging.Error.Println(logDetail(payload, response, err.Error()))
		return c.Status(response.Code).JSON(response)
	}

	if err := h.service.Create(payload); err != nil {
		response := models.Response{
			Code:    fiber.StatusInternalServerError,
			Message: fiber.ErrInternalServerError.Message,
		}

		logging.Error.Println(logDetail(payload, response, err.Error()))
		return c.Status(response.Code).JSON(response)
	}

	response := models.Response{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    "Article created successfully",
	}

	logging.Info.Println(logDetail(payload, response, response.Data.(string)))
	return c.JSON(response)
}

func (h *ArticlesHandler) GetArticles(c *fiber.Ctx) error {
	queryParams := models.GetArticles{
		Query:  c.Query("query"),
		Author: c.Query("author"),
	}

	articles, err := h.service.GetArticles(queryParams)
	if err != nil {
		response := models.Response{
			Code:    fiber.StatusInternalServerError,
			Message: fiber.ErrInternalServerError.Message,
		}

		logging.Error.Println(logDetail(queryParams, response, err.Error()))
		return c.Status(response.Code).JSON(response)
	}

	response := models.Response{
		Code:    fiber.StatusOK,
		Message: "Article retrieved successfully",
		Data:    articles,
	}

	logging.Info.Println(logDetail(queryParams, response, response.Message))
	return c.JSON(response)
}

func logDetail(request, response any, message string) string {
	detail := struct {
		Request  any    `json:"request"`
		Response any    `json:"response"`
		Message  string `json:"message"`
	}{
		Request:  request,
		Response: response,
		Message:  message,
	}

	detailBytes, _ := json.Marshal(detail)
	return string(detailBytes)
}

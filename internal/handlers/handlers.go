package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/usawyer/testZeroAgency/internal/service"
	"github.com/usawyer/testZeroAgency/models"
)

type Handler struct {
	s *service.Service
}

var limitPerPage = 5

func New(s *service.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitRoutes(app *fiber.App) {
	app.Put("/posts/:Id", h.EditPost)
	app.Post("/posts", h.CreatePost)
	app.Get("/posts", h.GetPosts)
}

func (h *Handler) CreatePost(c *fiber.Ctx) error {
	news := models.News{}

	err := c.BodyParser(&news)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Success": false, "Error": "invalid data input", "Message": err})
	}

	err = h.s.CreatePost(news)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Success": false, "Error": "news was not created", "Message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Success": true, "Id": news.Id})
}

func (h *Handler) EditPost(c *fiber.Ctx) error {
	news := models.News{}

	id, err := c.ParamsInt("Id", -1)
	if err != nil || id < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Success": false, "Error": "invalid id input"})
	}

	err = c.BodyParser(&news)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Success": false, "Error": "invalid data input"})
	}

	err = h.s.EditPost(id, news)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Success": false, "Error": "news wasn't update", "Message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Success": true, "Id": id})
}

func (h *Handler) GetPosts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	size := c.QueryInt("size", limitPerPage)

	searchParams := models.SearchParams{
		Offset: (page - 1) * limitPerPage,
		Limit:  size,
	}

	news, err := h.s.GetPosts(searchParams)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Success": false, "News": nil, "Message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Success": true, "News": news})
}

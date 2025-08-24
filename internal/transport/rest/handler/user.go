package handler

import (
	"WebProject/internal/core"
	"context"
	"net/http"
	"time"
)
import "github.com/gofiber/fiber/v2"

type UserService interface {
	GetAll(ctx context.Context) ([]*core.User, error)
	GetById(ctx context.Context, id string) (*core.User, error)
	CreateUser(ctx context.Context, user *core.User) (*core.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (handler *UserHandler) InitRoutes(app *fiber.App) {
	app.Get("/users", handler.GetAll)
	app.Get("/users/:userId", handler.GetById)
	app.Post("/users", handler.CreateUser)
}

// GetAll
// @Summary Get all users
// @Tags users
// @Description Returns the list of the users
// @Produce json
// @Status 200
// @Router /users [get]
func (handler *UserHandler) GetAll(ctx *fiber.Ctx) error {
	users, err := handler.service.GetAll(ctx.UserContext())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(
		fiber.Map{
			"users": users,
		})
}

func (handler *UserHandler) GetById(ctx *fiber.Ctx) error {
	ctxTimeOut, cancel := context.WithTimeout(ctx.UserContext(), 2*time.Millisecond)
	defer cancel()

	userChannel := make(chan *core.User)
	var err error
	go func(channel <-chan *core.User) {
		var user *core.User
		user, err = handler.service.GetById(ctxTimeOut, ctx.Params("userId"))
		userChannel <- user
	}(userChannel)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	select {
	case user := <-userChannel:
		{
			return ctx.Status(http.StatusOK).JSON(
				fiber.Map{
					"user": user,
				})
		}
	case <-ctxTimeOut.Done():
		{
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": ctxTimeOut.Err().Error(),
			})
		}
	}

}

func (handler *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	user := &core.User{}
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	savedUser, err := handler.service.CreateUser(ctx.UserContext(), user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"user": savedUser,
	})
}

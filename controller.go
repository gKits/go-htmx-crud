package crud

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	AddRoutes(fiber.Router)
}

/**
Utilities
**/

func HXError(c *fiber.Ctx, err error) error {
	c.Set("HX-Retarget", "#errors")
	c.Set("HX-Reswap", "afterbegin")
	return c.Render("_error", err)
}

func HXMiddleware(c *fiber.Ctx) error {
	if strings.ToLower(c.Get("HX-Request", "false")) != "true" {
		return HXError(c, nil)
	}

	return c.Next()
}

func GetID(c *fiber.Ctx) (uint, error) {
	paramId := c.Params("id")
	if paramId == "" {
		return 0, errors.New("id missing from request header")
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		return 0, fmt.Errorf("invalid id '%s'", paramId)
	}

	return uint(id), nil
}

/**
User Controller
**/

type UserController struct {
	repo Repository
}

func NewUserController(repo Repository) UserController {
	return UserController{repo}
}

func (ctrl *UserController) AddRoutes(g fiber.Router) {
	g.Get("/users", ctrl.getUsersHandler)
	g.Get("/user/:id", ctrl.getUserHandler)
	g.Post("/user/", ctrl.createUserHandler)
	g.Put("/user/:id", ctrl.updateUserHandler)
	g.Delete("/user/:id", ctrl.deleteUserHandler)
	g.Get("/users/avg", ctrl.avg)
}

func (ctrl *UserController) getUsersHandler(c *fiber.Ctx) error {
	var users []User
	var err error

	search := c.Query("search")
	if search != "" {
		users, err = ctrl.repo.Search(c.Context(), search)
	} else {
		users, err = ctrl.repo.All(c.Context())
	}
	if err != nil {
		return HXError(c, err)
	}

	c.Set("HX-Trigger", "user_table_change")

	return c.Render("usersTable", users)
}

func (ctrl *UserController) getUserHandler(c *fiber.Ctx) error {
	id, err := GetID(c)
	if err != nil {
		return HXError(c, err)
	}

	user, err := ctrl.repo.GetByID(c.Context(), id)
	if err != nil {
		return HXError(c, err)
	}

	return c.Render("userForm", user)
}

func (ctrl *UserController) createUserHandler(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return HXError(c, err)
	}

	created, err := ctrl.repo.Create(c.Context(), user)
	if err != nil {
		return HXError(c, err)
	}

	c.Set("HX-Trigger", "user_table_change")

	return c.Render("userRow", created)
}

func (ctrl *UserController) updateUserHandler(c *fiber.Ctx) error {
	id, err := GetID(c)
	if err != nil {
		return HXError(c, err)
	}

	var update User
	if err := c.BodyParser(&update); err != nil {
		return HXError(c, err)
	}

	updated, err := ctrl.repo.Update(c.Context(), id, update)
	if err != nil {
		return HXError(c, err)
	}

	c.Set("HX-Trigger", "user_table_change")

	return c.Render("userRow", updated)
}

func (ctrl *UserController) deleteUserHandler(c *fiber.Ctx) error {
	id, err := GetID(c)
	if err != nil {
		return HXError(c, err)
	}

	if err := ctrl.repo.Delete(c.Context(), id); err != nil {
		return HXError(c, err)
	}

	c.Set("HX-Trigger", "user_table_change")

	return c.SendStatus(fiber.StatusOK)
}

func (ctrl *UserController) avg(c *fiber.Ctx) error {
	avgAge, err := ctrl.repo.AvgAge(c.Context())
	if err != nil {
		return HXError(c, err)
	}

	avgHeight, err := ctrl.repo.AvgHeight(c.Context())
	if err != nil {
		return HXError(c, err)
	}

	return c.Render("avg", fiber.Map{"Age": avgAge, "Height": avgHeight})
}

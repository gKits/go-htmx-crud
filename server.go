package crud

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type Server struct {
	addr  string
	ctrls map[string]Controller
}

func NewServer(addr string) Server {
	return Server{
		addr:  addr,
		ctrls: map[string]Controller{},
	}
}

func (s *Server) Run() error {
	app := fiber.New(fiber.Config{
		Views: html.New("views", ".html"),
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("_index", nil)
	})

	app.Get("/empty", func(c *fiber.Ctx) error {
		return c.SendString("")
	})

	for pref, ctrl := range s.ctrls {
		g := app.Group(pref)
		ctrl.AddRoutes(g)
	}

	if err := app.Listen(s.addr); err != nil {
		return err
	}
	return nil
}

func (s *Server) Attach(pref string, ctrl Controller) {
	s.ctrls[pref] = ctrl
}

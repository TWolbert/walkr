package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var HTML []byte = getHTML()
var dir = getDir()

func ServeClient(c *fiber.Ctx) error {
	p := c.Path()
	if len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}

	clean := filepath.Clean(p)

	if clean == ".." || strings.HasPrefix(clean, ".. "+string(os.PathSeparator)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You little snoop",
		})
	}

	fsPath := filepath.Join(dir, "client", "dist", clean)

	if info, err := os.Stat(fsPath); err == nil && !info.IsDir() {
		return c.SendFile(fsPath)
	}

	return c.Type("html").SendString(string(HTML))
}

func ServeHTML(c *fiber.Ctx) error {
	return c.Type("html").SendString(string(HTML))
}

func getHTML() []byte {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return nil
	}

	htmlPath := filepath.Join((dir), "client/dist/index.html")
	html, err := os.ReadFile(htmlPath)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return html
}

func getDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return ""
	}

	return dir
}

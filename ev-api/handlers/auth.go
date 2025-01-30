package handlers

import (
	"context"
	"ev-api/config"
	"ev-api/models"
	"ev-api/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	query := `
	SELECT u.id, u.name, u.phone, u.email, u.password, u.customer_id, r.type 
	FROM tbl_user u
	JOIN tbl_role r ON u.role_id = r.id
	WHERE u.email = $1
	`

	err := config.DB.QueryRow(context.Background(), query, req.Email).Scan(
		&user.ID, &user.Name, &user.Phone, &user.Email, &user.Password, &user.CustomerID, &user.RoleType,
	)
	if err != nil {
		log.Println("Error querying user:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := utils.GenerateJWT(user.ID, user.Name, user.Email, user.RoleType, user.CustomerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": token})
}

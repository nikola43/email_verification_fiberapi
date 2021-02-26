package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sethvargo/go-password/password"
	"log"
	"os"
	"strings"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
	u "github.com/nikola43/rfpjforex_email_verification_api/utils"
)

type VerificationData struct {
	gorm.Model
	Email string `gorm:"type:varchar(128)" json:"email"`
	Code  string `gorm:"type:varchar(128)" json:"code"`
}

var DBGorm *gorm.DB

func main() {
	app := fiber.New()

	InitializeDbCorrection(
		GetEnvVariable("MYSQL_USER"),
		GetEnvVariable("MYSQL_PASSWORD"),
		GetEnvVariable("MYSQL_DATABASE"))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://rfpjforex.com",
		AllowHeaders:  "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"message": "hi",
		})
	})

	app.Get("/verify/:code", func(c *fiber.Ctx) error {
		code := c.Params("code")
		verificationData := &VerificationData{}

		// comprobamos la longitud del email
		if len(code) == 0 {
			return c.JSON(fiber.Map{"error": "empty code"})
		}

		// buscamos si existe algun registro con el email recibido
		DBGorm.Where("code = ?", code).First(&verificationData)

		if len(verificationData.Code) > 0 {
			verificationData.Code = code
			DBGorm.Where("code = ?", code).Unscoped().Delete(&verificationData)

			// veirificación OK !!
			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"success": true,
			})
		}

		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"success": false,
		})
	})

	app.Get("/auth/:email", func(c *fiber.Ctx) error {
		email := c.Params("email")
		verificationData := &VerificationData{}

		// comprobamos la longitud del email
		if len(email) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
			})
		}

		code, err := password.Generate(6, 0, 0, true, false)
		if err != nil {
			return err
		}
		code = strings.ToUpper(code)

		// buscamos si existe algun registro con el email recibido
		DBGorm.Where("email = ?", email).First(&verificationData)

		//fmt.Println(result)
		fmt.Println("verificationData")
		fmt.Println(verificationData)

		if len(verificationData.Email) > 0 {
			verificationData.Code = code
			verificationData.Email = email
			DBGorm.Save(&verificationData)
		} else {
			verificationData.Email = email
			verificationData.Code = code
			DBGorm.Create(&verificationData)
		}

		emailManager := u.Info{Code: code}
		emailManager.SendMailRecovery(verificationData.Email)

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
		})
	})

	app.Listen(":8080")
}

func InitializeDbCorrection(user, password, database_name string) {
	var err error
	connectionString := fmt.Sprintf(
		"%s:%s@/%s?parseTime=true",
		user,
		password,
		database_name,
	)
	DBGorm, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	//DBGorm.AutoMigrate(&VerificationData{})
}

// use godot package to load/read the .env file and
// return the value of the key
func GetEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

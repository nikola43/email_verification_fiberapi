package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	u "github.com/nikola43/rfpjforex_email_verification_api/utils"
	"github.com/sethvargo/go-password/password"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type VerificationData struct {
	gorm.Model
	Email string `gorm:"type:varchar(128)" json:"email"`
	Code  string `gorm:"type:varchar(128)" json:"code"`
	Ip    string `gorm:"type:varchar(24)" json:"ip"`
}

var DBGorm *gorm.DB

func main() {
	app := fiber.New()

	InitializeDbCorrection(
		GetEnvVariable("MYSQL_USER"),
		GetEnvVariable("MYSQL_PASSWORD"),
		GetEnvVariable("MYSQL_DATABASE"))

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"message": "hi",
		})
	})

	app.Get("/verify/:code", func(c *fiber.Ctx) error {
		code := c.Params("code")
		verificationData := &VerificationData{}
		ip := c.Get("X-Real-Ip")
		verificationData.Ip = ip
		verificationData.Code = code

		// comprobamos la longitud del email
		if len(code) == 0 {
			return c.JSON(fiber.Map{"error": "empty code"})
		}

		// buscamos si existe algun registro con el email recibido
		DBGorm.Where("code = ? AND ip = ? ", code, ip).First(&verificationData).Limit(1)

		if len(verificationData.Code) > 0 {
			today := time.Now()
			expDate := verificationData.CreatedAt.Add(15 * time.Minute)
			if today.Before(expDate) {
				// veirificación OK !!
				return c.Status(fiber.StatusOK).JSON(&fiber.Map{
					"success": true,
				})
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"success": false,
		})
	})

	app.Get("/auth/en/:email", func(c *fiber.Ctx) error {
		email := c.Params("email")
		verificationData := &VerificationData{}
		ip := c.Get("X-Real-Ip")

		// comprobamos la longitud del email
		if len(email) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"error": "empty email",
			})
		}
		emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		if !emailRegexp.MatchString(email) {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"error": "invalid email",
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

		DBGorm.Unscoped().Delete(&verificationData)

		verificationData.Code = code
		verificationData.Email = email
		verificationData.Ip = ip
		DBGorm.Create(&verificationData)

		emailManager := u.Info{Code: code}
		emailManager.SendMailRecoveryEn(verificationData.Email)

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
		})
	})

	app.Get("/auth/es/:email", func(c *fiber.Ctx) error {
		email := c.Params("email")
		verificationData := &VerificationData{}
		ip := c.Get("X-Real-Ip")

		// comprobamos la longitud del email
		if len(email) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"error": "empty email",
			})
		}
		emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		if !emailRegexp.MatchString(email) {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"error": "invalid email",
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

		DBGorm.Unscoped().Delete(&verificationData)

		verificationData.Code = code
		verificationData.Email = email
		verificationData.Ip = ip
		DBGorm.Create(&verificationData)

		emailManager := u.Info{Code: code, Ip: ip}
		emailManager.SendMailRecoveryEs(verificationData.Email)

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
		})
	})

	app.Listen(":8081")
}

func InitializeDbCorrection(user, password, database_name string) {
	var err error
	connectionString := fmt.Sprintf(
		"%s:%s@/%s?parseTime=true",
		user,
		password,
		database_name,
	)
	DBGorm, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
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

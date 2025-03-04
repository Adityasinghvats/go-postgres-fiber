package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/adix/books-fiber-postgres/models"
	"github.com/adix/books-fiber-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}
	//convert json from http request into golang struct Book
	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request failed"})
		return err
	}

	//asking gorm to create a db with same structure as my struct
	err = r.DB.Create(&book).Error
	if err != nil {
		//similar to creating a new Json reponse in express with json message and status
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "error creating db book"})
		return err
	}

	context.Status(http.StatusCreated).JSON(
		&fiber.Map{"message": "book has been added"})
	//everythin goes well then we will return nil for return type error
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"messgae": "id cannot be empty"})
		return nil
	}
	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"messgae": "could not delete book"})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"messgae": "book deleted successfully"})

	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "couldn't get the books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Books found",
		"data":    bookModels,
	})

	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"messgae": "id cannot be empty"})
		return nil
	}
	fmt.Println("the id is", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "couldn't get the book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Books found",
		"data":    bookModel,
	})

	return nil
}

// struct method signature
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	//r.*Book are all repository methods

	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_books/:id", r.DeleteBook)
	api.Get("/books", r.GetBooks)
	api.Get("/get_book/:id", r.GetBookByID)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("Coludnot load database", err)
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()

	//struct method
	r.SetupRoutes(app)
	app.Listen(":8080")
}

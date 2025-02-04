package controllers

import (
	"context"
	"net/http"
	"quotesapi/configs"
	"quotesapi/models"
	"quotesapi/responses"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var quoteCollection *mongo.Collection = configs.GetCollection(configs.DB, "quotes")
var validate = validator.New()

func CreateQuote(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var quote models.Quote
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&quote); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.QuoteResponse{Status: http.StatusBadRequest, Message: "BodyParser error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&quote); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.QuoteResponse{Status: http.StatusBadRequest, Message: "Validaton error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newQuote := models.Quote{
		Id:        primitive.NewObjectID(),
		Name:      quote.Name,
		BirthDate: quote.BirthDate,
		Email:     quote.Email,
		HomeSize:  quote.HomeSize,
		CarYear:   quote.CarYear,
		CarModel:  quote.CarModel,
		Status:    "new",
	}

	result, err := quoteCollection.InsertOne(ctx, newQuote)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.QuoteResponse{Status: http.StatusInternalServerError, Message: "Insert error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.QuoteResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAQuote(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	quoteId := c.Params("quoteId")
	var quote models.Quote
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(quoteId)

	err := quoteCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&quote)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.QuoteResponse{Status: http.StatusInternalServerError, Message: quoteId, Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.QuoteResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": quote}})
}

func GetAllQuotes(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var quotes []models.Quote
	defer cancel()
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetLimit(100)
	results, err := quoteCollection.Find(ctx, bson.M{}, opts)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.QuoteResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleQuote models.Quote
		if err = results.Decode(&singleQuote); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.QuoteResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		quotes = append(quotes, singleQuote)
	}

	return c.Status(http.StatusOK).JSON(
		responses.QuoteResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": quotes}},
	)
}

func GetMyQuotes(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	name := c.Params("name")
	var quotes []models.Quote
	defer cancel()

	results, err := quoteCollection.Find(ctx, bson.M{"name": name})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.QuoteResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleQuote models.Quote
		if err = results.Decode(&singleQuote); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.QuoteResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
		quotes = append(quotes, singleQuote)
	}

	return c.Status(http.StatusOK).JSON(
		responses.QuoteResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": quotes}},
	)
}
func coalesce(a string, b string) string {
	if len(a) > 0 {
		return a
	}
	return b
}
func EditAQuote(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	quoteId := c.Params("quoteId")
	var quote models.Quote
	var currentquote models.Quote
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(quoteId)

	//validate the request body
	if err := c.BodyParser(&quote); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.QuoteResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	// if validationErr := validate.Struct(&quote); validationErr != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(responses.QuoteResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	// }

	//get current quote details
	err := quoteCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&currentquote)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.QuoteResponse{Status: http.StatusInternalServerError, Message: quoteId, Data: &fiber.Map{"data": err.Error()}})
	}

	// update fields or use existing
	update := bson.M{"name": coalesce(quote.Name, currentquote.Name), "birthdate": coalesce(quote.BirthDate, currentquote.BirthDate), "email": coalesce(quote.Email, currentquote.Email), "homesize": coalesce(quote.HomeSize, currentquote.Email), "caryear": coalesce(quote.CarYear, currentquote.CarYear), "carmodel": coalesce(quote.CarModel, currentquote.CarModel), "status": coalesce(quote.Status, currentquote.Status)}

	result, err := quoteCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.QuoteResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated quote details
	var updatedQuote models.Quote
	if result.MatchedCount == 1 {
		err := quoteCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedQuote)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.QuoteResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.QuoteResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedQuote}})
}

func DeleteAQuote(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	quoteId := c.Params("quoteId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(quoteId)

	result, err := quoteCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.QuoteResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.QuoteResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Quote with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.QuoteResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Quote successfully deleted!"}},
	)
}

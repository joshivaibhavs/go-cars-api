package main

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
)

func registerRoutes() {

	app := fiber.New()

	app.Get("/cars", func(c *fiber.Ctx) error {
		query := bson.D{{}}
		cursor, err := mg.Db.Collection(CollectionName).Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var cars []Car = make([]Car, 0)

		if err := cursor.All(c.Context(), &cars); err != nil {
			return c.Status(500).SendString(err.Error())

		}

		return c.JSON(cars)
	})

	app.Get("/cars/:id", func(c *fiber.Ctx) error {
		params := c.Params("id")

		_id, err := primitive.ObjectIDFromHex(params)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		filter := bson.D{{Key: "_id", Value: _id}}

		var result Car

		if err := mg.Db.Collection(CollectionName).FindOne(c.Context(), filter).Decode(&result); err != nil {
			return c.Status(500).SendString("Something went wrong.")
		}

		return c.Status(fiber.StatusOK).JSON(result)
	})

	app.Post("/cars", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection(CollectionName)

		car := new(Car)
		if err := c.BodyParser(car); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		insertionResult, err := collection.InsertOne(c.Context(), car)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
		createdRecord := collection.FindOne(c.Context(), filter)

		createdCar := &Car{}
		createdRecord.Decode(createdCar)

		return c.Status(201).JSON(createdCar)
	})

	app.Put("/cars/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		carID, err := primitive.ObjectIDFromHex(idParam)

		if err != nil {
			return c.SendStatus(400)
		}

		car := new(Car)
		if err := c.BodyParser(car); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: carID}}
		update := bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "make", Value: car.Make},
					{Key: "model", Value: car.Model},
				},
			},
		}
		err = mg.Db.Collection(CollectionName).FindOneAndUpdate(c.Context(), query, update).Err()

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.SendStatus(404)
			}
			return c.SendStatus(500)
		}

		car.ID = idParam
		return c.Status(200).JSON(car)
	})

	app.Delete("/cars/:id", func(c *fiber.Ctx) error {
		carID, err := primitive.ObjectIDFromHex(
			c.Params("id"),
		)

		if err != nil {
			return c.SendStatus(400)
		}

		query := bson.D{{Key: "_id", Value: carID}}
		result, err := mg.Db.Collection(CollectionName).DeleteOne(c.Context(), &query)

		if err != nil {
			return c.SendStatus(500)
		}

		if result.DeletedCount < 1 {
			return c.SendStatus(404)
		}

		return c.SendStatus(204)
	})

	log.Fatal(app.Listen(":8080"))
}

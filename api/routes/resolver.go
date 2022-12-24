package routes

import (
	"github.com/Hariharan148/Url-Shortener-Go-Redis/api/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)


func ResolveUrl(c *fiber.Ctx)error{
	url := c.Params("url")

	r := database.Client(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error":"short url not found in database"})
	} else if err != nil{
		return c.Status(fiber.StatusInternalServerError ).JSON(fiber.Map{"error":"can't connect to db"})
	}

	dbInr := database.Client(1)
	defer dbInr.Close()

	_ = dbInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)
}
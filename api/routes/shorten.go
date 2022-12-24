package routes 

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/go-redis/redis/v8"
	"strconv"
	"google/uuid"
	"github.com/asaskevich/govalidator"
)


type Request struct {
	URL			string		  `json:"url"`
	CustomShort	string  	  `json:"custom_short"`
	Expiry		time.Duration `json:"expiry"`
}

type Response struct{
	URL				string		  `json:"url"`
	CustomShort		string  	  `json:"custom_short"`
	Expiry			time.Duration `json:"expiry"`
	XRateLimit  	int			  `json:"x_rate_limit`
	XRateLimitRest	time.Duration `json:"x_rate_limit_reset"`
}


func ShortenUrl(c *fiber.Ctx){
	var body *Request
	
	if err := c.BodyParser(&body); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"cannot parse the body"})
	}


	//rate limiting

	rd1 := database.Client(1)
	defer rd1.Close()

	val, err := rd1.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil{
		_ = rd1.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		valInt := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := rd1.TTL(database.Ctx. c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailabe).JSON(fiber.Map{
				"error": "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	//valid url 

	if !govalidator.IsURL(body.URL)	{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid url"})
	}

	//domain error

	if helpers.DomainError(body.URL){
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error":"Dont use our domain url! "})
	}

	// https enforcement

	body.URL = helpers.Enforehttp(body.URL)

	var id string

	if body.CustomShort == ""{
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort 
	}

	r := database.Client(0)
	defer r.Close()

	if body.Expiry == 0{
		body.Expiry = 24
	}

	val, _:= r.Get(database.Ctx, id).Result()
	if val != ""{
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error":"short url already in use"})
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry *3600 *time.Second).Err()

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"cant connect to database"})
	}



	resp := Response{
		URL: body.URL,
		CustomShort: "",
		Expiry: body.Expiry,
		XRateLimit: 10,
		XRateLimitRest: 30,

	}

	rd1.Decr(database.Ctx, c.IP())

	limit, _ := rd1.Get(database.Ctx, c.IP()).Result()
	resp.XRateLimit:= strconv.Atoi(limit)


	limitReset, _ : rd1.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitRest = limitReset / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id


	return c.Status(fiber.StatusOK).JSON(resp)
}
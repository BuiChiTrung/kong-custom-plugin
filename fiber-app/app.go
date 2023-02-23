package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/qinains/fastergoding"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var redisCtx = context.Background()

func GetCacheKeyHandler(c *fiber.Ctx) error {
	response := GetCacheKeyResponse{Data: nil}

	cacheKey := c.Params("key")
	val, err := rdb.Get(redisCtx, cacheKey).Result()

	if err == redis.Nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "cache key not exist"
	} else if err != nil {
		response.Code = fiber.StatusInternalServerError
		response.Message = "internal server err"
	} else {
		response.Code = fiber.StatusOK
		response.Message = "success"
		response.Data = &GetCacheKeyResponseData{Value: val}
	}
	return c.Status(response.Code).JSON(response)
}

func DelCacheKeyHandler(c *fiber.Ctx) error {
	response := DelCacheKeyResponse{}

	cacheKey := c.Params("key")
	_, err := rdb.Get(redisCtx, cacheKey).Result()

	if err == redis.Nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "cache key not exist"
	} else if err != nil {
		response.Code = fiber.StatusInternalServerError
		response.Message = "internal server err"
	} else {
		if err := rdb.Del(redisCtx, cacheKey).Err(); err != nil {
			response.Code = fiber.StatusInternalServerError
			response.Message = "internal server err"
		}
		response.Code = fiber.StatusOK
		response.Message = "success"
	}
	return c.Status(response.Code).JSON(response)
}

func UpsertCacheKeyHandler(c *fiber.Ctx) error {
	reqBody := UpsertCacheKeyRequest{}
	response := UpsertCacheKeyResponse{}

	if err := c.BodyParser(&reqBody); err != nil {
		return err
	}

	_, err := rdb.Get(redisCtx, reqBody.CacheKey).Result()
	if err == redis.Nil || err == nil {
		if err := rdb.Set(redisCtx, reqBody.CacheKey, reqBody.Value, 0).Err(); err != nil {
			response.Code = fiber.StatusInternalServerError
			response.Message = "internal server err"
		} else {
			response.Code = fiber.StatusOK
			response.Message = "success upsert key"
		}
	} else {
		response.Code = fiber.StatusInternalServerError
		response.Message = "internal server err"
	}

	return c.Status(response.Code).JSON(response)
}

func main() {
	fastergoding.Run() // hot reload
	app := fiber.New()

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// TODO: trung.bc - key => cache-key
	app.Get("/proxy-cache/:key", GetCacheKeyHandler)
	app.Delete("/proxy-cache/:key", DelCacheKeyHandler)
	app.Post("/proxy-cache", UpsertCacheKeyHandler)

	app.Listen(":3000")
}

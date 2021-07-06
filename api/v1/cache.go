package v1

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/khan1507017/redis-app/database/rds"
	"github.com/khan1507017/redis-app/dto"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type CacheControllerInf interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Get(c echo.Context) error
}

type CacheControllerInstance struct {
}

func CacheController() CacheControllerInf {
	return new(CacheControllerInstance)
}

func (e CacheControllerInstance) Create(c echo.Context) error {
	var input dto.RedisObject
	if err := c.Bind(&input); err != nil {
		log.Println("input error: ", err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	rdsErr := rds.GetRedisMaster().Set(context.Background(), input.Key, input.Value, 0).Err()
	if rdsErr != nil {
		log.Println("database error: ", rdsErr.Error())
		return c.JSON(http.StatusInternalServerError, rdsErr.Error())
	}
	return c.JSON(http.StatusOK, "successfully added")
}

func (e CacheControllerInstance) Update(c echo.Context) error {

	var input dto.RedisObject
	if err := c.Bind(&input); err != nil {
		log.Println("input error: ", err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	_, rdsErr := rds.GetRedisMaster().Get(context.Background(), input.Key).Result()
	if rdsErr == redis.Nil {
		log.Println("key not found")
		return c.JSON(http.StatusBadRequest, errors.New("key not found"))
	} else if rdsErr != nil {
		log.Println("database error: ", rdsErr.Error())
		return c.JSON(http.StatusInternalServerError, rdsErr.Error())
	} else {
		err := rds.GetRedisMaster().Set(context.Background(), input.Key, input.Value, 0)
		if err != nil {
			log.Println("key update failed")
			return c.JSON(http.StatusInternalServerError, errors.New("key update failed: "+err.String()))
		}
	}
	return c.JSON(http.StatusOK, nil)
}

func (e CacheControllerInstance) Delete(c echo.Context) error {
	params := c.QueryParams()
	if params.Get("key") == "" {
		log.Println("key not found in params")
		return c.JSON(http.StatusBadRequest, "key missing in query param")
	}
	_, rdsErr := rds.GetRedisMaster().Get(context.Background(), params.Get("key")).Result()
	if rdsErr == redis.Nil {
		log.Println("key not found")
		return c.JSON(http.StatusBadRequest, "key not found")
	} else {
		err := rds.GetRedisMaster().Del(context.Background(), params.Get("key")).Err()
		if err != nil {
			log.Println("key deletion failed: " + err.Error())
			return c.JSON(http.StatusBadRequest, "key update failed: "+err.Error())
		}
	}

	return c.JSON(http.StatusOK, "successfully deleted")
}

func (e CacheControllerInstance) Get(c echo.Context) error {
	params := c.QueryParams()
	if params.Get("key") == "" {
		log.Println("key not found in params")
		return c.JSON(http.StatusBadRequest, "key missing in query param")
	}
	value, rdsErr := rds.GetRedisMaster().Get(context.Background(), params.Get("key")).Result()
	if rdsErr == redis.Nil {
		log.Println("key not found")
		return c.JSON(http.StatusBadRequest, "key not found")
	}

	respObj := dto.RedisObject{Key: params.Get("key"), Value: value}
	return c.JSON(http.StatusOK, respObj)
}

package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

// ServerPort and RedisPort are embedded into actual code, as they are part of the codebase
//If you run Redis in a different port, please change the default value
//If you want to use a different port for the application, change the server port. Please also make changes to docker file (expose the same port)
const ServerPort = "4040"
const RedisPort = "6379"

//envVars
var RedisPassword string
var RedisMasterEndpoint string
var RedisSlaveEndpoints [50]string
var RedisSlaveCount int

//temp variable
var boolVal bool
var slaveCountTemp string
var RunMode string

func InitEnvironmentVariables() error {

	//checking runMode
	RunMode = os.Getenv("RUN_MODE")
	if RunMode == "" {
		RunMode = DEVELOP
	}
	var err error
	log.Println("RUN MODE:", RunMode)

	//loading envArs from .env file is runMode != PRODUCTION
	if RunMode != PRODUCTION {
		err = godotenv.Load()
		if err != nil {
			log.Println("ERROR: ", err.Error())
			return err
		}
	}

	//DB CREDENTIALS + Cluster Endpoint + Others
	RedisPassword, boolVal = os.LookupEnv("REDIS_PASSWORD")
	if boolVal == false {
		return errors.New("REDIS_PASSWORD not found in envVars")
	}
	RedisMasterEndpoint, boolVal = os.LookupEnv("MASTER_ENDPOINT")
	if boolVal == false {
		return errors.New("MASTER_ENDPOINT not found in envVars")
	}
	err = initSlaveEndpoints()
	if err != nil {
		return err
	}
	fmt.Println("environment vars loaded")
	return nil
}
func initSlaveEndpoints() error {
	slaveCountTemp, boolVal = os.LookupEnv("SLAVE_COUNT")
	if boolVal == true {
		var err error
		RedisSlaveCount, err = strconv.Atoi(slaveCountTemp)
		if err != nil {
			return err
		}
		if RedisSlaveCount < 0 || RedisSlaveCount > 50 {
			return errors.New("invalid slave number: " + slaveCountTemp)
		}
	} else {
		RedisSlaveCount = 0
		return nil
	}
	for i := 0; i < RedisSlaveCount; i++ {
		RedisSlaveEndpoints[i], boolVal = os.LookupEnv("SLAVE_ENDPOINT_" + strconv.Itoa(i))
		if boolVal == false {
			return errors.New("SLAVE_ENDPOINT_" + strconv.Itoa(i) + " not found in envVars")
		}
	}
	return nil
}

package container

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	db "github.com/syedomair/ex-paygate/lib/db"
	"github.com/syedomair/ex-paygate/lib/tools/logger"
)

// Const
const (
	ServiceNameEnvVar = "SERVICE_NAME"
	LogLevelEnvVar    = "LOG_LEVEL"
	DatabaseURLEnvVar = "DATABASE_URL"
	PortEnvVar        = "PORT"
	SigningKeyEnvVar  = "SIGNINGKEY"
)

// Container interface
type Container interface {
	Logger() logger.Logger
	Db() *gorm.DB
	ServiceName() string
	Port() string
	SigningKey() string
}

type container struct {
	logger               logger.Logger
	db                   *gorm.DB
	environmentVariables map[string]string
}

func (c *container) Db() *gorm.DB {
	db, err := db.New(c.getRequiredEnvVar(DatabaseURLEnvVar))
	db.LogMode(true)
	if err != nil {
		c.Logger().Critical("", "Could not connect to DB: %v", err)
	} else {
		c.Logger().Info("", "Connected to DB")
	}
	return db
}

func (c *container) Logger() logger.Logger {
	if c.logger == nil {
		c.logger = logger.New(c.getRequiredEnvVar(LogLevelEnvVar), c.ServiceName()+"#", os.Stdout)
	}
	return c.logger
}

func (c *container) SigningKey() string {
	return c.getRequiredEnvVar(SigningKeyEnvVar)
}

func (c *container) ServiceName() string {
	return c.getRequiredEnvVar(ServiceNameEnvVar)
}

func (c *container) Port() string {
	return c.getRequiredEnvVar(PortEnvVar)
}

func (c *container) getRequiredEnvVar(key string) string {
	value, ok := c.environmentVariables[key]
	if !ok {
		panic(fmt.Errorf("missing mandatory env var:: %q", key))
	}
	return value
}

// New Public
func New(envVars map[string]string) Container {
	return &container{
		environmentVariables: envVars,
	}
}

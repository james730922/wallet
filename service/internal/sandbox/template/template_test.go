package template

import (
	"context"
	"os"
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	"go.uber.org/dig"

	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/thirdparty/db"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	redisCli "github.com/james730922/wallet/service/internal/thirdparty/redis"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

type templateTestApp struct {
	dig.In

	Template ITemplate
}

type exampleTestSuite struct {
	suite.Suite
	app      *templateTestApp
	db       *gorm.DB
	ctx      context.Context
	redisCli *redis.Client
}

func (s *exampleTestSuite) SetupSuite() {
	// Init config
	tools.InitRootFolder("../../../..")

	// Init logger for sandbox.
	logger.TestMock()

	//set unit test local cache
	cache.NewCache()

	// DB connection to local mysql.
	s.db = db.NewGORMMock()
	// ctx
	s.ctx = context.Background()
	// Redis connection to local Redis cluster.
	s.redisCli = redisCli.NewRedisMock()

	// Init Your Provider
	binder := dig.New()
	s.Require().Nil(binder.Provide(NewTemplate))
	s.Require().Nil(binder.Invoke(func(app templateTestApp) {
		s.app = &app
	}))
}

func (s *exampleTestSuite) SetupTest() {
}

func (s *exampleTestSuite) TearDownTest() {
}

func (s *exampleTestSuite) TearDownSuite() {
}

func (s *exampleTestSuite) TestCase1() {
}

func TestTemplate(t *testing.T) {
	if os.Getenv("RUN_DB_INTEGRATION_TESTS") != "1" {
		t.Skip("set RUN_DB_INTEGRATION_TESTS=1 to run MySQL/Redis integration tests")
	}
	suite.Run(t, &exampleTestSuite{})
}

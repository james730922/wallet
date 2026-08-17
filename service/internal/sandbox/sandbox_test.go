package sandbox

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"

	"github.com/james730922/wallet/service/internal/thirdparty/cache"
	"github.com/james730922/wallet/service/internal/thirdparty/db"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	redisCli "github.com/james730922/wallet/service/internal/thirdparty/redis"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

type sandboxTestSuite struct {
	suite.Suite
	db       *gorm.DB
	ctx      context.Context
	redisCli *redis.Client
}

func (s *sandboxTestSuite) SetupSuite() {
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
}

func (s *sandboxTestSuite) TestLocalMySQL() {
	type testModel struct {
		ID        uint64     `gorm:"column:id;primary_key;NOT NULL"`
		Name      string     `gorm:"column:name;type:varchar(20)"`
		Value     string     `gorm:"column:value;type:varchar(50)"`
		CreatedAt *time.Time `gorm:"column:created_at;type:DATETIME(6);default:CURRENT_TIMESTAMP(6);NOT NULL"`
		UpdatedAt *time.Time `gorm:"column:updated_at;type:DATETIME(6);default:CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)"`
	}

	// Drop & create table
	if s.db.HasTable(&testModel{}) {
		s.db.DropTable(&testModel{})
	}
	s.db.CreateTable(&testModel{})

	// Insert & select
	data1 := testModel{
		Name:  "kobe",
		Value: "8",
	}
	s.Require().Nil(s.db.Create(&data1).Error)

	var res1 testModel
	s.Require().Nil(s.db.
		Model(&testModel{}).Where("name = ?", data1.Name).First(&res1).Error)
	s.Require().Equal(data1.Value, res1.Value)

	s.db.DropTable(&testModel{})
}

func (s *sandboxTestSuite) SetupTest() {
}

func (s *sandboxTestSuite) TearDownTest() {
}

func (s *sandboxTestSuite) TearDownSuite() {
}

func TestSandbox(t *testing.T) {
	if os.Getenv("RUN_DB_INTEGRATION_TESTS") != "1" {
		t.Skip("set RUN_DB_INTEGRATION_TESTS=1 to run MySQL/Redis integration tests")
	}
	suite.Run(t, &sandboxTestSuite{})
}

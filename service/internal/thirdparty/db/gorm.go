package db

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/james730922/wallet/service/internal/models"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/conf"
)

const (
	GormQueryOption = "gorm:query_option"
)

const (
	SyntaxForUpdate = "FOR UPDATE"
)

var (
	GormSetSelectForUpdate = func() (string, string) { return GormQueryOption, SyntaxForUpdate }
)

func NewGORM() *gorm.DB {
	connInfo := conf.DB().GetRawSQLUrl()

	logger.ApLog().Info("init primary gorm")
	db, err := gorm.Open("mysql", connInfo)
	if err != nil {
		logger.ApLog().Panic(err)
	}
	db.DB().SetMaxIdleConns(2)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(conf.DB().GetGormLogMode())

	return db
}

func NewGORMSlave() *gorm.DB {
	connInfo := conf.DB().GetSQLSlaveUrl()

	logger.ApLog().Info("init read gorm")
	db, err := gorm.Open("mysql", connInfo)
	if err != nil {
		logger.ApLog().Panic(err)
	}
	db.DB().SetMaxIdleConns(2)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(conf.DB().GetGormLogMode())

	return db
}

func NewGORMMock() *gorm.DB {
	connInfo := conf.DB().GetMockSQLUrl()
	logger.ApLog().Debug("init mock gorm")
	db, err := gorm.Open("mysql", connInfo)
	if err != nil {
		logger.ApLog().Panic(err)
	}
	db.DB().SetMaxIdleConns(2)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(conf.DB().GetGormLogMode())

	return db
}

// ParseWhere 1.parse start_at_xxx, end_at_xxx, in_xxx, like_xxx  2.for others, using = or IN operator
func ParseWhere(whereCond map[string]interface{}) func(dc *gorm.DB) *gorm.DB {
	if len(whereCond) == 0 {
		return func(dc *gorm.DB) *gorm.DB { return dc }
	}
	newMap := map[string]interface{}{}
	for key, val := range whereCond {
		if strings.HasPrefix(key, "start_at_") {
			// start_at_xxx -> where(xxx >= ?, val(timestamp))
			mapKey := strings.Replace(key, "start_at_", "", 1) + " >="

			switch val.(type) {
			case int64:
				newMap[mapKey] = time.Unix(val.(int64), 0).UTC()
			case *int64:
				newMap[mapKey] = time.Unix(*val.(*int64), 0).UTC()
			default:
				newMap[mapKey] = val
			}
		} else if strings.HasPrefix(key, "end_at_") {
			// end_at_xxx -> where(xxx < ?, val(timestamp))
			mapKey := strings.Replace(key, "end_at_", "", 1) + " <"

			switch val.(type) {
			case int64:
				newMap[mapKey] = time.Unix(val.(int64), 0).UTC()
			case *int64:
				newMap[mapKey] = time.Unix(*val.(*int64), 0).UTC()
			default:
				newMap[mapKey] = val
			}
		} else if strings.HasPrefix(key, "in_") {
			// in_xxx -> where(xxx IN (?), ...val)
			mapKey := strings.Replace(key, "in_", "", 1) + " IN"
			newMap[mapKey] = val
		} else if strings.HasPrefix(key, "null_") {
			// null_xxx -> where(xxx NULL (?), ...val)
			mapKey := strings.Replace(key, "null_", "", 1) + " NULL"
			newMap[mapKey] = val
		} else if strings.HasPrefix(key, "like_") {
			// like_xxx -> where(xxx LIKE ?, %val%)
			mapKey := strings.Replace(key, "like_", "", 1) + " LIKE"

			switch reflect.ValueOf(val).Kind() {
			case reflect.Ptr:
				newMap[mapKey] = fmt.Sprintf("%%%v%%", reflect.ValueOf(val).Elem())
			default:
				newMap[mapKey] = fmt.Sprintf("%%%v%%", val)
			}
		} else if strings.HasPrefix(key, "start_num_") {
			// start_num_xxx -> where(xxx >= ?, val(num))
			mapKey := strings.Replace(key, "start_num_", "", 1) + " >="

			switch val.(type) {
			case int:
				newMap[mapKey] = val.(int)
			case *int:
				newMap[mapKey] = *val.(*int)
			case int64:
				newMap[mapKey] = val.(int64)
			case *int64:
				newMap[mapKey] = *val.(*int64)
			case float64:
				newMap[mapKey] = val.(float64)
			case *float64:
				newMap[mapKey] = *val.(*int64)
			default:
				newMap[mapKey] = val
			}
		} else if strings.HasPrefix(key, "end_num_") {
			// end_num_xxx -> where(xxx < ?, val(end_num))
			mapKey := strings.Replace(key, "end_num_", "", 1) + " <"

			switch val.(type) {
			case int:
				newMap[mapKey] = val.(int)
			case *int:
				newMap[mapKey] = *val.(*int)
			case int64:
				newMap[mapKey] = val.(int64)
			case *int64:
				newMap[mapKey] = *val.(*int64)
			case float64:
				newMap[mapKey] = val.(float64)
			case *float64:
				newMap[mapKey] = *val.(*int64)
			default:
				newMap[mapKey] = val
			}
		} else {
			newMap[key] = val
		}
	}

	return func(dc *gorm.DB) *gorm.DB {
		return buildWhere(dc, newMap)
	}
}

func buildWhere(db *gorm.DB, wheres map[string]interface{}) *gorm.DB {
	for key, value := range wheres {
		// check key is valid
		keys := strings.Split(key, " ")
		if len(keys) > 2 {
			continue
		}

		refVal := reflect.ValueOf(value)
		if refVal.Kind() == reflect.Ptr && !refVal.IsNil() {
			refVal = refVal.Elem()
			value = refVal.Interface()
		}

		// set sqlKey, ex: lock -> `lock`, member.id -> `member`.`id`
		sqlKey := keys[0]

		sqlKeys := strings.Split(sqlKey, ".")
		switch len(sqlKeys) {
		case 1:
			sqlKey = "`" + sqlKey + "`"
		case 2:
			sqlKey = "`" + sqlKeys[0] + "`.`" + sqlKeys[1] + "`"
		default:
			continue
		}

		// set operator
		operator := ""
		switch len(keys) {
		case 1:
			switch refVal.Kind() {
			case reflect.Slice:
				operator = "IN"
			default:
				operator = "="
			}
		case 2:
			operator = strings.ToUpper(keys[1])
		}

		switch operator {
		case "=", ">", "<", ">=", "<=", "!=", "<>", "<=>", "LIKE":
			db = db.Where(sqlKey+" "+operator+" ?", value)
		case "IN":
			db = db.Where(sqlKey+" "+operator+" (?)", value)
		case "NULL":
			db = db.Where(sqlKey + " is " + operator)
		default:
			continue
		}
	}

	return db
}

func ParsePaging(paging *models.Paging) func(dc *gorm.DB) *gorm.DB {
	return func(dc *gorm.DB) *gorm.DB {
		if paging != nil && paging.Size != 0 {
			dc = dc.Limit(paging.GetSize()).Offset(paging.GetOffset())
		}

		return dc
	}
}

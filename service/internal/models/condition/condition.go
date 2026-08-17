package condition

import (
	"reflect"
	"strings"
	"sync"

	"github.com/james730922/wallet/service/internal/models"
)

var (
	valPaging        = reflect.ValueOf(models.Paging{})
	pagingTypeString = valPaging.Type().String()
)

func NewQuery(cond interface{}) *Query {
	where := newCond(structToMap(cond))

	var paging models.Paging
	if val, ok := where.Result()[pagingTypeString]; ok {
		where.Delete(pagingTypeString)
		if tmp, ok := val.(models.Paging); ok {
			paging = tmp
		}
	}

	return &Query{
		where:  where,
		paging: &paging,
	}
}

type IQuery interface {
	Where() map[string]interface{}
	Paging() *models.Paging
}

type Query struct {
	where  *cond
	paging *models.Paging
}

func (q *Query) WhereCond() *cond {
	return q.where
}

func (q *Query) Where() map[string]interface{} {
	return q.where.Result()
}

func (q *Query) Paging() *models.Paging {
	return q.paging
}

func NewUpdate(update interface{}, where interface{}) *Update {
	return &Update{
		update: newCond(structToMap(update)),
		where:  newCond(structToMap(where)),
	}
}

type IUpdate interface {
	Update() map[string]interface{}
	Where() map[string]interface{}
}

type Update struct {
	update *cond
	where  *cond
}

func (q *Update) UpdateCond() *cond {
	return q.update
}

func (q *Update) Update() map[string]interface{} {
	return q.update.Result()
}

func (q *Update) WhereCond() *cond {
	return q.where
}

func (q *Update) Where() map[string]interface{} {
	return q.where.Result()
}

func newCond(m _map) *cond {
	if m == nil {
		m = make(map[string]interface{})
	}

	return &cond{
		_map: m,
	}
}

type _map map[string]interface{}

type cond struct {
	mx sync.RWMutex
	_map
}

func (c *cond) Set(key string, val interface{}) *cond {
	if !c.isNil(val) {
		c.mx.Lock()
		c._map[key] = val
		c.mx.Unlock()
	}
	return c
}

func (c *cond) NullSet(key string, val interface{}) *cond {
	c.mx.Lock()
	c._map[key] = val
	c.mx.Unlock()
	return c
}

func (c *cond) Delete(key string) *cond {
	c.mx.Lock()
	delete(c._map, key)
	c.mx.Unlock()
	return c
}

func (c *cond) Result() map[string]interface{} {
	return c._map
}

func (c *cond) isNil(v interface{}) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		return rv.IsNil()
	}
	return false
}

func structToMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if item == nil {
		return res
	}

	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		fieldValue := reflectValue.Field(i)

		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		field := fieldValue.Interface()
		if tag != "" && tag != "-" {
			tags := strings.Split(tag, ",")
			tag = strings.TrimSpace(tags[0])
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else {
				res[tag] = field
			}
		} else if fieldValue.Type() == valPaging.Type() {
			res[pagingTypeString] = field
		}
	}

	return res
}

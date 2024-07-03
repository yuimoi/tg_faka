package functions

import (
	"errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"reflect"
	"strings"
	_type "tg_go_faka/internal/utils/type"
	"time"
	"unicode"
)

type StructToMapData struct {
	Mode          StructMarshalMode
	Keys          []string
	IgnoreNilFlag bool
}
type StructMarshalMode int

const (
	StructToMapIncludeMode StructMarshalMode = iota
	StructToMapExcludeMode
)

func StructToMap(v interface{}, opts ...StructToMapData) map[string]interface{} {
	// 设置默认值
	var data StructToMapData
	if len(opts) > 0 {
		data = opts[0]
	} else {
		data = StructToMapData{
			Mode:          StructToMapExcludeMode,
			Keys:          make([]string, 0),
			IgnoreNilFlag: false,
		}
	}
	mode := data.Mode
	keys := data.Keys
	ignoreNilFlag := data.IgnoreNilFlag

	resultMap := make(map[string]interface{})
	vValue := reflect.Indirect(reflect.ValueOf(v)) // Automatically handles pointers

	for i := 0; i < vValue.NumField(); i++ {
		field := vValue.Field(i)
		typeField := vValue.Type().Field(i)
		jsonTag := typeField.Tag.Get("json")
		tagParts := strings.Split(jsonTag, ",")
		jsonKey := tagParts[0]

		if ignoreNilFlag {
			//跳过值为nil,判断前要判断是否为指针
			//1.用于请求参数，为EditProduct handler在用，同一个edit，用于切换开关和修改内容
			//2.用于返回参数，有这个跳过nil会导致值为nil的kv缺失，如果这里出问题，希望给上order_id:null这种返回值，则函数加一个参数切换
			if field.Kind() == reflect.Ptr && field.IsNil() {
				continue
			}
		}

		//跳过jsongtag为空
		if jsonTag == "" {
			continue
		}
		//排除模式、包括模式
		if mode == StructToMapIncludeMode {
			if !SliceContainString(keys, jsonKey) {
				continue
			}
		} else {
			if SliceContainString(keys, jsonKey) {
				continue
			}
		}

		// 跳过含有omitempty的,但保留输入include为最优先
		if len(tagParts) >= 2 && SliceContainString(tagParts[1:], "omitempty") {
			if !(mode == StructToMapIncludeMode && SliceContainString(keys, jsonKey)) {
				continue
			}
		}

		// 重置uuid默认值,这里是用于返回请求有关外键的参数,有些外键为null,但是默认值是00000
		if field.Type() == reflect.TypeOf(uuid.UUID{}) && field.Interface() == uuid.Nil {
			resultMap[jsonKey] = nil

		} else if field.Type() == reflect.TypeOf(decimal.Decimal{}) {
			resultMap[jsonKey] = field.Interface()
		} else {
			resultMap[jsonKey] = field.Interface()
		}
	}

	return resultMap
}
func SliceContainString(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func ParseUUID(idString string) (uuid.UUID, error) {
	myUUID, err := uuid.Parse(idString)
	if err != nil {
		return myUUID, errors.New("无法解析UUID")
	}
	return myUUID, nil
}
func IsWhitespace(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
func ApplyPaginationQueryData(query *gorm.DB, paginationQueryData _type.PaginationQueryDataStruct) *gorm.DB {
	if paginationQueryData.Limit <= 0 || paginationQueryData.Page <= 0 {
		paginationQueryData.Limit = 1
		paginationQueryData.Page = 1
	}

	query = query.Limit(paginationQueryData.Limit)
	offset := (paginationQueryData.Page - 1) * paginationQueryData.Limit
	query = query.Offset(offset)

	return query
}
func TimestampToDatetime(timestamp int64) string {
	// 将时间戳转换为 time.Time 对象
	t := time.Unix(timestamp, 0)
	formattedTime := t.Format("2006-01-02 15:04:05")
	return formattedTime
}

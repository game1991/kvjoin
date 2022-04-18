package kvjoin

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func (j *join) parseURLValues(values url.Values) {
	for key := range values {
		j.data[key] = values.Get(key)
	}
}

func (j *join) parseURLString() (err error) {
	str, ok := j.src.(string)
	if !ok {
		return nil
	}
	index := strings.Index(str, "?")
	values, err := url.ParseQuery(str[index+1:])
	if err != nil {
		return err
	}
	j.parseURLValues(values)
	return
}

func (j *join) parseStruct(rv reflect.Value) (err error) {
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("unsupported type :%s", rv.Type().Name())
	}
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		key := rt.Field(i).Name
		if j.options.StructTag != "" {
			if j.isSkippedField(rt.Field(i).Tag) {
				continue
			}
			key = rt.Field(i).Tag.Get(j.options.StructTag)
		}
		value := rv.Field(i)
		if err = j.parseValue(key, value); err != nil {
			return
		}
	}
	return
}

// 是否是忽略的字段
func (j *join) isSkippedField(tag reflect.StructTag) bool {
	v, ok := tag.Lookup(j.options.StructTag)
	return ok && v == "-"
}

func (j *join) parseMap(rv reflect.Value) (err error) {
	if rv.Kind() != reflect.Map {
		return fmt.Errorf("unsupported type :%s", rv.Type().Name())
	}
	for _, key := range rv.MapKeys() {
		kv := ""
		switch key.Kind() {
		case reflect.Bool:
			kv = strconv.FormatBool(key.Bool())
		case reflect.String:
			kv = key.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			kv = strconv.FormatInt(key.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			kv = strconv.FormatUint(key.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			kv = strconv.FormatFloat(key.Float(), 'f', -1, 64)
		default:
			return fmt.Errorf("unsupported key type :%s", rv.Type().Name())
		}
		value := rv.MapIndex(key)
		if err = j.parseValue(kv, value); err != nil {
			return
		}
	}
	return nil
}

func (j *join) parseValue(key string, rv reflect.Value) (err error) {
	if !rv.CanInterface() {
		return
	}
	if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		//if j.options.Unwrap {
		//	err = j.parseSlice(key, rv)
		//}

	case reflect.Struct:
		if j.options.Unwrap {
			err = j.parseStruct(rv)
		}
	case reflect.Map:
		if j.options.Unwrap {
			err = j.parseMap(rv)
		}
	case reflect.Ptr, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return fmt.Errorf("unsupported type :%s", rv.Type().Name())
	default:
		j.data[key] = rv.Interface()
	}

	return
}

func (j *join) parseSlice(key string, rv reflect.Value) (err error) {
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return
	}
	// for i := 0; i < rv.Len(); i++ {
	// 	itemV := rv.Index(i)
	// 	switch itemV.Kind() {
	// 	case reflect.Bool:

	// 	}

	// }
	j.data[key] = rv.Interface()
	return
}

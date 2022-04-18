package kvjoin

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strings"
)

func (j *join) joinInASCII() ([]string, error) {
	keys := make([]string, 0, len(j.data))
	for key := range j.data {
		if j.exceptKeys(key) {
			continue
		}
		keys = append(keys, key)
	}
	switch j.options.Order {
	case ASCII:
		sort.Strings(keys)
	case ASCIIDesc:
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] > keys[j]
		})
	}
	list := make([]string, 0, len(keys))
	for _, key := range keys {
		if value, ignore := j.getValue(key); !ignore {
			list = append(list, value)
		}
	}
	return list, nil
}

func (j *join) exceptKeys(key string) bool {
	if len(j.options.ExceptKeys) > 0 {
		exceptKeys := make(map[string]int)
		for _, except := range j.options.ExceptKeys {
			exceptKeys[except] = 1
		}

		if _, ok := exceptKeys[key]; ok {
			return true
		}
	}

	// 如果标签中包含omitempty并且对应的值为空值，则进行排除
	temp := j.data[key]
	isZero := reflect.ValueOf(temp).IsZero()
	// 读取key中是否存在omitempty
	split := strings.Split(key, ",")
	for _, item := range split {
		if item == "omitempty" && isZero {
			return true
		}
	}
	return false
}

func (j *join) joinInDefined() ([]string, error) {
	list := make([]string, 0, len(j.data))
	for _, key := range j.options.DefinedOrders {
		if _, ok := j.data[key]; !ok {
			continue
		}
		if value, ignore := j.getValue(key); !ignore {
			list = append(list, value)
		}
	}
	return list, nil
}

func (j *join) getValue(key string) (value string, ignore bool) {
	temp := j.data[key]
	rv := reflect.ValueOf(temp)
	if j.options.IgnoreEmpty && rv.IsZero() {
		return "", true
	}
	// 如果tag对应的存在多字段以，分割时取第一个为key
	index := strings.Index(key, ",")
	if index > 0 {
		key = key[:index]
	}

	value = fmt.Sprintf("%v", temp)
	switch j.options.URLCoding {
	case None:
	case Encoding:
		value = url.QueryEscape(value)
	case Decoding:
		value, _ = url.QueryUnescape(value)
	}
	if j.options.IgnoreKey {
		return value, false
	}
	return key + j.options.KVSep + value, false
}

package kvjoin

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

/*
	1、实现生成签名的键值对kv规则格式的字符串，例如k1=v1&k2=v2&k3=v3&.....kn=vn，键排序默认采用ASCII码由小到大排序，如果键对应的存在空值，默认舍弃。
	2、本方法Join支持自定义对象接口实现，调用NewJoiner方法注入接口实现
	3、由于签名规则是键值对形式出现，故本方法不支持嵌套复杂结构对象，因为嵌套解析出来的对象key值可能存在重名现象等其他问题，故不做嵌套深度解析。结构体字段对应属性应为基础类型，例如string、int、bool
	4、Join方法传参src源数据类型只支持url地址字符串、结构体和map字典类型，其他类型暂不支持。

   注意: Join函数的入参，如果出现同名字段并且option.UnWrap字段为true时，会解析结构体或者map中的字段，同时，如果整体字段存在重复会被覆盖，最终的urlValue对相同字段的解析中只会存在一个kv值
*/

// Joiner ... 实现生成签名格式kv形式的拼接字符串
type Joiner interface {
	Join(src interface{}, opts ...Optional) (dst string, err error)
}

var defaultJoiner Joiner = &join{}

type join struct {
	src     interface{}
	options *Option
	data    map[string]interface{} // src parsed from src
}

// Join ... 实现生成签名格式kv形式的拼接字符串
func Join(src interface{}, opts ...Optional) (dst string, err error) {
	return defaultJoiner.Join(src, opts...)
}

// 如果传入参数为空，则给默认值
func defaultOption() *Option {
	op := &Option{
		Sep:           "&",
		KVSep:         "=",
		IgnoreKey:     false,
		IgnoreEmpty:   true,
		Order:         ASCII,
		DefinedOrders: nil,
		StructTag:     "json",
		URLCoding:     0,
	}

	// if option.URLCoding != 0 {
	// 	op.URLCoding = option.URLCoding
	// }
	// if option.Unwrap {
	// 	op.Unwrap = option.Unwrap
	// }
	return op

}

// Optional 可选项
type Optional func(option *Option)

// WithSep 每组kv连接的分隔符，例如k1=v1&k2=v2
func WithSep(sep string) Optional {
	return func(option *Option) {
		option.Sep = sep
	}
}

// WithKVSep k和v之间的分隔符选项，例如k=v
func WithKVSep(kvSep string) Optional {
	return func(option *Option) {
		option.KVSep = kvSep
	}
}

// WithIgnoreKey 忽略key拼接，保留vaule
func WithIgnoreKey(ignoreKey bool) Optional {
	return func(option *Option) {
		option.IgnoreKey = ignoreKey
	}
}

// WithExceptKeys 忽略的key对象
func WithExceptKeys(exceptKeys []string) Optional {
	return func(option *Option) {

	}
}

// WithIgnoreEmpty 忽略空值对象
func WithIgnoreEmpty(ignoreEmpty bool) Optional {
	return func(option *Option) {
		option.IgnoreEmpty = ignoreEmpty
	}
}

// WithOrder 自定义排序
func WithOrder(orderStyle orderStyle, definedOrders []string) Optional {
	return func(option *Option) {
		if len(definedOrders) > 0 && orderStyle == Defined {
			option.DefinedOrders = definedOrders
		}
		option.Order = orderStyle
	}
}

// WithStructTag 自定义tag
func WithStructTag(structTag string) Optional {
	return func(option *Option) {
		if structTag != "" {
			option.StructTag = structTag
		}
	}
}

// WithURLCoding url编码
func WithURLCoding(urlCoding urlCoding) Optional {
	return func(option *Option) {
		option.URLCoding = urlCoding
	}
}

// WithUnwrap 深度解析对象(不建议使用，未完善)
func WithUnwrap(unwrap bool) Optional {
	return func(option *Option) {
		option.Unwrap = unwrap
	}
}

func (j *join) initJoiner(src interface{}, opts ...Optional) {
	options := defaultOption()
	for _, opt := range opts {
		opt(options)
	}

	j.src = src
	j.options = options
	j.data = map[string]interface{}{}
}

func (j *join) Join(src interface{}, opts ...Optional) (dst string, err error) {
	// 初始化写入源
	j.initJoiner(src, opts...)
	// 解析源数据
	values, ok := j.src.(url.Values)
	if ok {
		j.parseURLValues(values)
	} else {
		rv := reflect.ValueOf(j.src)
		if rv.Kind() == reflect.Ptr {
			rv = reflect.Indirect(rv)
		}
		switch rv.Kind() {
		case reflect.String:
			err = j.parseURLString()
		case reflect.Struct:
			err = j.parseStruct(rv)
		case reflect.Map:
			err = j.parseMap(rv)
		default:
			return "", fmt.Errorf("unsupported type :%s", rv.Type().Name())
		}
	}
	if err != nil {
		return "", err
	}
	var list []string
	switch j.options.Order {
	case ASCII, ASCIIDesc:
		list, err = j.joinInASCII()
	case Defined:
		if len(j.options.DefinedOrders) == 0 {
			return "", errors.New("need 'DefinedOrders' in Defied order mode")
		}
		list, err = j.joinInDefined()
	default:
		return "", errors.New("unsupported order")
	}
	return strings.Join(list, j.options.Sep), nil

}

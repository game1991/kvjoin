package kvjoin

type orderStyle uint

// 排序方式
const (
	ASCII orderStyle = iota
	ASCIIDesc
	Defined
)

type urlCoding uint

// url地址是否
const (
	None urlCoding = iota
	Encoding
	Decoding
)

// Option 可选项
type Option struct {
	Sep           string     // the separator between keys (or contains their values)
	KVSep         string     // the separator between key & its value
	IgnoreKey     bool       // whether ignore key , if yes, the key will be ignored, but the value will reserve
	IgnoreEmpty   bool       // whether ignore empty value , if yes, the key & its value will be ignored
	ExceptKeys    []string   // the keys & their values will be ignored
	Order         orderStyle // the join order
	DefinedOrders []string   // the keys order, using with Order == Defined
	StructTag     string     // struct tag, using when src struct type, if not set, will use struct filed name, only support export fields
	URLCoding     urlCoding  // the value format, using when format value
	Unwrap        bool       // whether unwrap the internal map or struct (由于签名kv形式，如果对象存在嵌套解析出来字段名会存在重名现象或者其他问题，故此功能默认不解析)--> 启用该功能，不考虑重名问题，这个由使用者自行决定，如果有重名则会覆盖
}

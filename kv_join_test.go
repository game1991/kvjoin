package kvjoin

import (
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
)

var randomStringSeed = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	//A2Z 大写的A-Z
	A2Z = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//A2z 小写的A-z
	A2z = "abcdefghijklmnopqrstuvwxyz"
	//Zero2Nine 0-9
	Zero2Nine = "0123456789"
	//A2F 大写A-F
	A2F = "ABCDEF"
	//A2f 小写a-f
	A2f = "abcdef"
	//A2Z2Num 大写的A-Z和0-9
	A2Z2Num = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	//A2z2Num 小写的a-z和0-9
	A2z2Num = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type callback struct {
	Bool    bool                   `json:"b,omitempty"`
	APPID   string                 `json:"-"`
	APPName string                 `json:"app_name,omitempty"`
	APPNo   string                 `json:"app_no"`
	Count   int                    `json:"count,omitempty"`
	Score   []int                  `json:"score,omitempty"`
	Player  []*foo                 `json:"player"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type foo struct {
	Animal []*animal
	Count  int
}

type animal struct {
	Name  string
	Bingo bool
}

func TestJoin(t *testing.T) {
	type args struct {
		src  interface{}
		opts []Optional
	}
	tests := []struct {
		name    string
		args    args
		wantDst string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "测试json对象标签",
			args: args{
				src: &callback{
					Bool:    true,
					APPID:   UUID32(),
					APPName: RandomString(A2Z, 8),
					APPNo:   UUID32(),
					Count:   0,
					Score:   []int{},
					Player:  []*foo{},
				},
				opts: []Optional{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDst, err := Join(tt.args.src, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Join() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			bts, _ := json.MarshalIndent(tt.args.src, "", "\t")
			t.Log(string(bts))
			t.Log(gotDst)

			// if gotDst != tt.wantDst {
			// 	t.Errorf("Join() = %v, want %v", gotDst, tt.wantDst)
			// }
		})
	}
}

func Test_join_Join(t *testing.T) {
	type fields struct {
		src     interface{}
		options *Option
		data    map[string]interface{}
	}
	type args struct {
		src  interface{}
		opts []Optional
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantDst string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &join{
				src:     tt.fields.src,
				options: tt.fields.options,
				data:    tt.fields.data,
			}
			gotDst, err := j.Join(tt.args.src, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("join.Join() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDst != tt.wantDst {
				t.Errorf("join.Join() = %v, want %v", gotDst, tt.wantDst)
			}
		})
	}
}

// RandomString 生成随机字符串(伪随机)(在高并发下使用，rand是并发不安全的)
func RandomString(source string, length int) string {
	runes := []rune(source)
	size := len(runes)
	result := make([]rune, length)
	for i := range result {
		result[i] = runes[randomStringSeed.Intn(size)]
	}
	return string(result)
}

func UUID32() string {
	u := uuid.New()
	var buf = make([]byte, 32)
	hex.Encode(buf, u[:4])
	hex.Encode(buf[8:12], u[4:6])
	hex.Encode(buf[12:16], u[6:8])
	hex.Encode(buf[16:20], u[8:10])
	hex.Encode(buf[20:], u[10:])
	return string(buf[:])
}

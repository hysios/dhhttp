package dhhttp

import (
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/hysios/utils/convert"
)

func TestDecoder_parseBrackets(t *testing.T) {
	type fields struct {
		rd io.Reader
	}
	type args struct {
		val string
	}

	test := func(val string, key string, brackets []string, ok bool) {
		var tt = struct {
			name       string
			fields     fields
			args       args
			wantKey    string
			wantIndexs []string
			wantOk     bool
		}{
			fields:     fields{rd: nil},
			args:       args{val: val},
			wantKey:    key,
			wantIndexs: brackets,
			wantOk:     ok,
		}

		t.Run(tt.name, func(t *testing.T) {
			dec := &Decoder{
				rd: tt.fields.rd,
			}
			gotKey, gotIndexs, gotOk := dec.parseBrackets(tt.args.val)
			if gotKey != tt.wantKey {
				t.Errorf("Decoder.parseBrackets() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if !reflect.DeepEqual(gotIndexs, tt.wantIndexs) {
				t.Errorf("Decoder.parseBrackets() gotIndexs = %v, want %v", gotIndexs, tt.wantIndexs)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Decoder.parseBrackets() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
	test("Regions[0][0]", "Regions", []string{"0", "0"}, true)
	test("Regions[0][2]", "Regions", []string{"0", "2"}, true)
	test("Regions[0][3]", "Regions", []string{"0", "3"}, true)
	test("Regions[1][0]", "Regions", []string{"1", "0"}, true)
	test(`Regions[1]["test"]`, "Regions", []string{"1", "test"}, true)
}

func TestDecoder_adjustSlice(t *testing.T) {
	type args struct {
		v   interface{}
		idx int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{v: &[]int{1, 2, 3}, idx: 1},
			want: 3,
		},
		{
			args: args{v: &[]int{1, 2}, idx: 3},
			want: 4,
		},
		{
			args: args{v: &[]int{1, 2, 3}, idx: 10},
			want: 11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := &Decoder{}
			if got := dec.adjustSlice(reflect.ValueOf(tt.args.v), tt.args.idx); !reflect.DeepEqual(got.Len(), tt.want) {
				t.Errorf("Decoder.adjustSlice() = %v, len %v want %v", got, got.Len(), tt.want)
			}
		})
	}
}

func TestDecoder_adjustMap(t *testing.T) {

	type args struct {
		v   interface{}
		key string
	}
	tests := []struct {
		name   string
		args   args
		want   reflect.Value
		canSet bool
	}{
		{
			args:   args{v: &map[string]interface{}{"hello": "world"}, key: "bob"},
			canSet: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := &Decoder{}
			if got := dec.adjustMap(reflect.ValueOf(tt.args.v), tt.args.key); !got.Elem().CanSet() {
				t.Errorf("Decoder.adjustMap() = %v, want %v", got, tt.want)
			}
			t.Logf("%#v", tt.args.v)
		})
	}
}

func TestDecoder_parseTag(t *testing.T) {
	type fields struct {
	}
	type args struct {
		t reflect.StructTag
	}
	tests := []struct {
		name  string
		args  args
		want  *Tag
		skip  bool
		want1 bool
	}{
		{
			args:  args{t: reflect.StructTag(`dahua:"password,omitempty"`)},
			want:  &Tag{Name: "password", OmitEmpty: true},
			want1: true,
		},
		{
			args:  args{t: reflect.StructTag(`dahua:"-"`)},
			skip:  true,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := &Decoder{}
			got, skip, got1 := dec.parseTag(tt.args.t)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decoder.parseTag() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Decoder.parseTag() got1 = %v, want %v", got1, tt.want1)
			}
			if skip != tt.skip {
				t.Errorf("Decoder.parseTag() skip = %v, want %v", skip, tt.skip)
			}
		})
	}
}

func TestDecoder_lookupField(t *testing.T) {
	type args struct {
		name string
		v    interface{}
		fn   func(st reflect.StructField, fv reflect.Value, ok bool) error
	}

	test := func(val interface{}, fieldName string, found bool, wantErr bool) {
		tt := struct {
			name    string
			args    args
			wantErr bool
		}{
			args: args{v: val, name: fieldName, fn: func(st reflect.StructField, fv reflect.Value, ok bool) error {
				t.Logf("fv %v, canset %v, ok %v", fv, fv.CanSet(), ok)
				if ok != found {
					t.Errorf("Decoder.lookupField() cb ok = %v, wanFound %v", ok, found)
				}
				return nil
			}},
			wantErr: wantErr,
		}
		t.Run(tt.name, func(t *testing.T) {
			dec := &Decoder{}
			if err := dec.lookupField(tt.args.name, reflect.ValueOf(tt.args.v), tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.lookupField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	test(&struct{ Name string }{}, "name", true, false)
	test(&struct {
		Name string
		Age  int
	}{}, "age", true, false)

	test(&struct {
		Name string
	}{}, "notfound", false, true)

	test(&struct {
		Name     string
		UserName int `dahua:"user_name"`
	}{}, "user_name", true, false)
}

func TestDecoder_searchBySelect(t *testing.T) {

	type args struct {
		v        interface{}
		selector string
		set      FieldSetter
	}

	test := func(strc interface{}, selector string, found bool, val interface{}) {
		tt := struct {
			name    string
			args    args
			wantErr bool
		}{
			args: args{
				v:        strc,
				selector: selector,
				set: func(v reflect.Value, typ reflect.Type, ok bool) error {
					if v.CanSet() {
						switch typ.Kind() {
						case reflect.String:
							s, _ := convert.String(val)
							v.SetString(s)
						case reflect.Bool:
							bo, _ := convert.Bool(val)
							v.SetBool(bo)
						case reflect.Int:
							i, _ := convert.Int(val)
							v.SetInt(int64(i))
						case reflect.Float64:
							f, _ := convert.Float(val)
							v.SetFloat(f)
						}
					}
					return nil
				},
			},
		}
		t.Run(tt.name, func(t *testing.T) {
			dec := &Decoder{}
			if err := dec.searchBySelect(reflect.ValueOf(tt.args.v), tt.args.selector, tt.args.set); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.searchBySelect() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("v %+v", tt.args.v)
		})
	}

	test(&struct {
		Name string
	}{}, "name", true, "hello")

	test(&struct {
		User struct {
			Name string
			Age  int
		}
	}{}, "user.name", true, "hello")

	test(&struct {
		User struct {
			Name    string
			Age     int
			Admin   bool
			Friends []struct {
				Name   string
				UserID int
			}
		}
	}{}, "user.friends[0].name", true, "bob")

	test(&struct {
		User struct {
			Name    string
			Age     int
			Admin   bool
			Friends []struct {
				Name   string
				UserID int
			}
		}
	}{}, "user.friends[3].UserID", true, 13)

	test(&struct {
		User struct {
			Name    string
			Age     int
			Admin   bool
			Friends []struct {
				Name    string
				UserID  int
				Profile struct {
					Money   float64
					Expires time.Time
				}
			}
		}
	}{}, "user.friends[3].profile.money", true, 18.8)

	test(&struct {
		User struct {
			Name   string
			Age    int
			Admin  bool
			Matrix [][]string
		}
	}{}, "user.matrix[3][3]", true, "hello")
}

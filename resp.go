package dhhttp

import (
	"bufio"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/hysios/utils/convert"
)

type Decoder struct {
	rd io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{rd: r}
}

func (dec *Decoder) Decode(val interface{}) error {
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Ptr {
		return errors.New("decode argument val can't be assigment")
	}

	v = reflect.Indirect(v)
	var s = bufio.NewScanner(dec.rd)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		line := s.Text()
		if err := dec.decodeLine(line, v); err != nil {
			continue
			// return err
		}
	}

	if err := s.Err(); err != nil {
		return err
	}
	return nil
}

type FieldSetter func(val reflect.Value, typ reflect.Type, index int, ok bool) error

func (dec *Decoder) decodeLine(s string, v reflect.Value) error {
	ss := strings.SplitN(s, "=", 2)
	if len(ss) < 2 {
		return errors.New("invalid expression format")
	}

	selector, value := ss[0], ss[1]

	return dec.searchBySelect(v, selector, func(val reflect.Value, typ reflect.Type, index int, ok bool) error {
		if !ok {
			return nil
		}

		switch typ.Kind() {
		case reflect.Int:
			i, _ := convert.Int(value)
			val.SetInt(int64(i))
		case reflect.Float64:
			f, _ := convert.Float(value)
			val.SetFloat(f)
		case reflect.Bool:
			b, _ := convert.Bool(value)
			val.SetBool(b)
		case reflect.String:
			s, _ := convert.String(value)
			val.SetString(s)
		case reflect.Slice, reflect.Array:
			e := val.Index(index)
			switch e.Kind() {
			case reflect.Int:
				i, _ := convert.Int(value)
				if e.CanSet() {
					e.SetInt(int64(i))
				}
			case reflect.Float64:
				f, _ := convert.Float(value)
				if e.CanSet() {
					e.SetFloat(f)
				}
			case reflect.Bool:
				b, _ := convert.Bool(value)
				if e.CanSet() {
					e.SetBool(b)
				}
			case reflect.String:
				s, _ := convert.String(value)
				if e.CanSet() {
					e.SetString(s)
				}
			default:
				return errors.New("invalid type")
			}
		default:
			return ErrNonimplement
		}
		return nil
	})
}

func (dec *Decoder) lookupField(name string, v reflect.Value, fn func(st reflect.StructField, fv reflect.Value, ok bool) error) error {
	v = reflect.Indirect(v)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		var fieldName string
		if tag, skip, ok := dec.parseTag(ft.Tag); ok {
			fieldName = tag.Name
		} else if skip {
			continue
		} else {
			fieldName = ft.Name
		}

		if strings.EqualFold(fieldName, name) {
			if err := fn(ft, v.Field(i), true); err != nil {
				return err
			}
			return nil
		}
	}

	if err := fn(reflect.StructField{}, reflect.Value{}, false); err != nil {
		return err
	}
	return errors.New("not found field")
}

func (dec *Decoder) searchBySelect(v reflect.Value, selector string, set FieldSetter) error {
	var (
		ss = strings.Split(selector, ".")
		l  = len(ss) - 1
		p  = reflect.Indirect(v)
	)

	for i, s := range ss {
		// t := p.Type()
		var child reflect.Value
		if key, brackets, ok := dec.parseBrackets(s); ok {
			if err := dec.lookupField(key, p, func(st reflect.StructField, fv reflect.Value, ok bool) error {
				if !ok {
					return nil
				}
				var (
					bl = len(brackets) - 1
					// bv reflect.Value
				)
				for j, bra := range brackets {
					if idx, ok := convert.Int(bra); ok { // Slice
						fv.Set(dec.adjustSlice(fv.Addr(), idx))

						if j == bl {
							if i == l {
								set(fv, fv.Type(), idx, true)
								return nil
							}
							child = fv.Index(idx)
						} else {
							fv = fv.Index(idx)
						}
					} else if key, ok := convert.String(bra); ok { // Map
						fv.SetMapIndex(reflect.ValueOf(key), dec.adjustMap(fv.Addr(), key))
						if j == bl {
							if i == l {
								set(fv, fv.Type(), 0, true)
								return nil
							}
						}

						child = fv.MapIndex(reflect.ValueOf(key))
					} else {
						return nil
					}
				}
				p = fv
				return nil
			}); err != nil {
				return err
			}

		} else {
			if err := dec.lookupField(s, p, func(st reflect.StructField, fv reflect.Value, ok bool) error {
				if !ok {
					return nil
				}

				if i < l {
					switch st.Type.Kind() {
					case reflect.Struct:
					default:
						return errors.New("scalar type can't out multi dot sections of selector")
					}

					child = fv
				} else {
					set(fv, st.Type, 0, true)
				}

				return nil
			}); err != nil {
				return err
			}
		}
		p = child
	}

	return nil
}

func (dec *Decoder) parseBrackets(val string) (key string, indexs []string, ok bool) {
	var (
		s     scanner.Scanner
		start bool
		c     int
	)
	s.Init(strings.NewReader(val))

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		// fmt.Printf("%s: %s\n", s.Position, s.TokenText())
		switch tok {
		case scanner.Ident:
			if !start {
				key = s.TokenText()
			} else {
				indexs = append(indexs, s.TokenText())
			}
		case '[':
			if start {
				ok = false
				return
			}
			start = true
		case ']':
			if !start {
				ok = false
				return
			}
			c++
			start = false
		case scanner.String:
			raw, _ := strconv.Unquote(s.TokenText())
			indexs = append(indexs, raw)
		case scanner.Int:
			indexs = append(indexs, s.TokenText())
		case ',':
			return
		default:
			ok = false
			return
		}
	}

	ok = c > 0
	return
}

func (dec *Decoder) adjustSlice(v reflect.Value, idx int) reflect.Value {
	if v.Kind() != reflect.Ptr {
		panic(errors.New("alice must can be set"))
	}

	e := reflect.Indirect(v)

	if !(e.Kind() == reflect.Array || e.Kind() == reflect.Slice) {
		panic(errors.New("must slice or array value"))
	}

	if idx < e.Len() {
		return v.Elem()
	}

	if v.IsValid() {
		dst := reflect.MakeSlice(e.Type(), idx+1, idx+1)

		reflect.Copy(dst, v.Elem())

		v.Elem().Set(dst)
	}

	return v.Elem()
}

func (dec *Decoder) adjustMap(v reflect.Value, key string) reflect.Value {
	if v.Kind() != reflect.Ptr {
		panic(errors.New("alice must can be set"))
	}

	v = reflect.Indirect(v)

	if v.Kind() != reflect.Map {
		panic(errors.New("must map"))
	}

	t := v.Type()
	if v.IsValid() {
		if !(t.Elem().Kind() == reflect.Ptr || t.Elem().Kind() == reflect.Interface) {
			panic(errors.New("map value type must ptr or interface"))
		}

		ev := reflect.New(t.Elem())
		v.SetMapIndex(reflect.ValueOf(key), ev)
		return ev
	}

	return reflect.Value{}
}

type Tag struct {
	Name      string
	OmitEmpty bool
	Value     string
	Type      reflect.Type
}

const (
	TagName = "dahua"
)

func pure(s string) string {
	val, _ := strconv.Unquote(s)
	return val
}

func (dec *Decoder) parseTag(t reflect.StructTag) (tag *Tag, skip bool, ok bool) {
	var (
		value = t.Get(TagName)
	)
	tag = &Tag{}

	if len(value) == 0 {
		return nil, false, false
	}

	ss := strings.SplitN(value, ",", -1)
	if len(ss) == 1 { // only name
		if len(ss[0]) == 1 && ss[0][0] == '-' {
			skip = true
			return nil, skip, true
		}
		tag.Name = ss[0]
		ok = true
		return
	} else if len(ss) == 2 { // name or type
		tag.Name = ss[0]
		if ss[1] == "omitempty" {
			tag.OmitEmpty = true
		} else {
			if typ, ok := builtinTypes[ss[1]]; ok {
				tag.Type = typ
			} else {
				panic(errors.New("must register type"))
			}
		}

		ok = true
		return
	}

	ok = false
	return
}

var builtinTypes = map[string]reflect.Type{
	"int":    reflect.TypeOf(new(int)).Elem(),
	"string": reflect.TypeOf(new(string)).Elem(),
	"bool":   reflect.TypeOf(new(bool)).Elem(),
	"float":  reflect.TypeOf(new(float64)).Elem(),
	"byte":   reflect.TypeOf(new(byte)).Elem(),
	"rune":   reflect.TypeOf(new(rune)).Elem(),
}

func Registry(typ string, val interface{}) {
	if _, ok := builtinTypes[typ]; ok {
		return
	}

	builtinTypes[typ] = reflect.TypeOf(val)
}

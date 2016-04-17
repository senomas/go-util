package util

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
)

const (
	// MaxFieldLen const
	MaxFieldLen = 100
	// MaxSliceLen const
	MaxSliceLen = 100
	tabSpace    = "  "
)

// Dump value
func Dump(v interface{}) string {
	var bb bytes.Buffer
	fdump(&bb, 0, 0, "", reflect.ValueOf(v))
	return bb.String()
}

// Fdump value
// func Fdump(w io.Writer, v interface{}) {
// 	fdump(w, 0, 0, "", reflect.ValueOf(v), "", "")
// }

func fdump(w io.Writer, rx int, rd int, tab string, value reflect.Value) {
	rx++
	if rx > 100 {
		fmt.Fprint(w, "/* too many recursive */")
	}
	if rd > 10 {
		fmt.Fprint(w, "/* too deep */")
	}
	if value.Kind() == reflect.Interface && !value.IsNil() {
		elm := value.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			value = elm
		}
	}
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() == reflect.String {
		fmt.Fprintf(w, "\"%s\" /* Type:%s */", value, value.Type().String())
	} else if kindIn(value.Kind(), reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64) {
		fmt.Fprintf(w, "%+v /* Type:%s */", value, value.Type().String())
	} else if kindIn(value.Kind(), reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64) {
		fmt.Fprintf(w, "%+v /* Type:%s */", value, value.Type().String())
	} else if value.Kind() == reflect.Struct {
		ntab := tab + tabSpace
		fmt.Fprintf(w, "{")
		il := value.NumField()
		for i, j := 0, 0; i < il && i < MaxFieldLen; i++ {
			tf := value.Type().Field(i)
			tag := strings.Split(value.Type().Field(i).Tag.Get("dump"), ",")
			if !InStringSlice("ignore", tag...) {
				if j > 0 {
					fmt.Fprintf(w, ",\n%s\"%s\": ", ntab, tf.Name)
				} else {
					fmt.Fprintf(w, "\n%s\"%s\": ", ntab, tf.Name)
				}
				vf := value.Field(i)
				fdump(w, rx, rd+1, ntab, vf)
				j++
			}
		}
		fmt.Fprintf(w, "\n%s}", tab)
	} else if value.Kind() == reflect.Slice {
		ntab := tab + tabSpace
		if value.Type().String() == "[]uint8" {
			il := value.Len()
			fmt.Fprintf(w, "[ /* Type:%s Length:%v */\n%s      00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F\n%s0000  ", value.Type().String(), il, ntab, ntab)
			for i, j := 0, 0; i < il; i++ {
				if j == 16 {
					fmt.Fprintf(w, "\n%s%04X  ", ntab, i)
					j = 1
				} else {
					j++
				}
				fmt.Fprintf(w, "%02X ", value.Index(i).Interface())
			}
			fmt.Fprintf(w, "\n%s]%s", tab, "")
		} else {
			fmt.Fprintf(w, "[ /* Type:%s Length:%v */\n%s", value.Type().String(), value.Len(), ntab)
			il := value.Len()
			for i := 0; i < il && i < MaxSliceLen; i++ {
				vi := value.Index(i)
				if i > 0 {
					fmt.Fprintf(w, ",\n%s/* Index:%v */ ", ntab, i)
				} else {
					fmt.Fprint(w, "/* Index:0 */ ")
				}
				fdump(w, rx, rd+1, ntab, vi)
			}
			if il > MaxSliceLen {
				fmt.Fprintf(w, "\n%s/* too many items */", ntab)
			}
			fmt.Fprintf(w, "\n%s]%s", tab, "")
		}
	} else if value.Kind() == reflect.Map {
		ntab := tab + tabSpace
		fmt.Fprintf(w, "{ /* Type:%s */\n%s", value.Type().String(), ntab)
		keys := value.MapKeys()
		for k, kv := range keys {
			if k > 0 {
				fmt.Fprintf(w, ",\n%s\"%v\": ", ntab, kv)
			} else {
				fmt.Fprintf(w, "\"%v\": ", kv)
			}
			fdump(w, rx, rd+1, ntab, value.MapIndex(kv))
		}
		fmt.Fprintf(w, "\n%s}%s", tab, "")
	} else {
		fmt.Fprintf(w, "undefined /* Type:%s Value:[%v] */", value.Type().String(), value)
	}
}

func kindIn(value reflect.Kind, list ...reflect.Kind) bool {
	for _, v := range list {
		if value == v {
			return true
		}
	}
	return false
}

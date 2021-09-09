package config

import (
	"io/fs"
	"io/ioutil"
	"reflect"
	"strconv"

	env "github.com/enorith/environment"
	"github.com/enorith/supports/reflection"
	"gopkg.in/yaml.v3"
)

var EnvPrefix = ""

func Unmarshal(file string, out interface{}) error {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	return UnmarshalBytes(data, out)
}

func UnmarshalFS(fsys fs.FS, filename string, out interface{}) error {
	data, err := fs.ReadFile(fsys, filename)
	if err != nil {
		return err
	}
	return UnmarshalBytes(data, out)
}

func UnmarshalBytes(data []byte, out interface{}) error {

	if err := yaml.Unmarshal(data, out); err != nil {
		return err
	}

	UnmarshalEnv(out)
	return nil
}

func UnmarshalNode(node yaml.Node, out interface{}) error {

	if err := (&node).Decode(out); err != nil {
		return err
	}

	UnmarshalEnv(out)
	return nil
}

func UnmarshalEnv(config interface{}) {
	v := reflection.StructValue(config)
	t := reflection.StructType(config)
	if t.Kind() == reflect.Struct {
		decodeEnvStruct(t, v)
	}
}
func decodeEnvStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		ft := sf.Type
		fv := v.Field(i)
		if ft.Kind() == reflect.Struct {
			decodeEnvStruct(ft, fv)
		} else if ft.Kind() == reflect.Map {
			//
			for _, k := range fv.MapKeys() {
				mv := fv.MapIndex(k)
				if mv.Kind() == reflect.Struct {
					decodeEnvStruct(mv.Type(), mv)
				}

				if mv.Kind() == reflect.Ptr {
					mvv := reflect.Indirect(mv)
					decodeEnvStruct(mvv.Type(), mvv)
				}
			}
		} else {
			if key := sf.Tag.Get("env"); key != "" {
				decodeEnv(ft, fv, key, true) // use env fisrt
			}

			if def := sf.Tag.Get("default"); def != "" {
				applyDefault(ft, fv, def)
			}
		}
	}
}

func decodeEnv(ft reflect.Type, fv reflect.Value, key string, prioritize bool) {
	key = EnvPrefix + key
	if env.GetString(key) == "" {
		// return if env not set
		return
	}
	if !fv.CanAddr() {
		return
	}

	if fv.IsZero() || prioritize {
		switch ft.Kind() {
		case reflect.String:
			fv.SetString(env.GetString(key))
		case reflect.Int, reflect.Int32, reflect.Int64:
			fv.SetInt(env.GetInt64(key))
		case reflect.Bool:
			fv.SetBool(env.GetBoolean(key))
		case reflect.Float32, reflect.Float64:
			fv.SetFloat(env.GetFloat64(key))
		case reflect.Slice:
			if ft.Elem().Kind() == reflect.Uint8 {
				fv.SetBytes([]byte(env.GetString(key)))
			}
		}
	}

}

func applyDefault(ft reflect.Type, fv reflect.Value, def string) {
	if fv.IsZero() && fv.CanAddr() {
		switch ft.Kind() {
		case reflect.String:
			fv.SetString(def)
		case reflect.Int, reflect.Int32, reflect.Int64:
			i, _ := strconv.ParseInt(def, 10, 64)
			fv.SetInt(i)
		case reflect.Bool:
			b, _ := strconv.ParseBool(def)
			fv.SetBool(b)
		case reflect.Float32, reflect.Float64:
			f, _ := strconv.ParseFloat(def, 64)
			fv.SetFloat(f)
		}
	}
}

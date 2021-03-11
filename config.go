package config

import (
	"io/fs"
	"io/ioutil"
	"reflect"
	"strconv"

	env "github.com/enorith/environment"
	"github.com/enorith/supports/reflection"
	"gopkg.in/yaml.v2"
)

type Config interface {
	GetValue(key string) (interface{}, bool)
	GetInt(key string) (int, bool)
	GetString(key string) (string, bool)
	GetBool(key string) (bool, bool)
}

type SimpleConfig struct {
	config map[string]interface{}
}

func (c *SimpleConfig) Load(file string) (*SimpleConfig, error) {
	data, err := ioutil.ReadFile(file)

	if e := yaml.Unmarshal(data, &c.config); e != nil {
		return nil, e
	}

	return c, err
}

func (c *SimpleConfig) GetValue(key string) (interface{}, bool) {
	if data, ok := c.config[key]; ok {
		return data, true
	}

	return nil, false
}

func (c *SimpleConfig) GetInt(key string) (int, bool) {
	if data, ok := c.GetValue(key); ok {
		if result, o := data.(int); o {
			return result, true
		}
		return 0, false
	}
	return 0, false
}

func (c *SimpleConfig) GetString(key string) (string, bool) {
	if data, ok := c.GetValue(key); ok {
		if result, o := data.(string); o {
			return result, true
		}
		return "", false
	}

	return "", false
}

func (c *SimpleConfig) GetBool(key string) (bool, bool) {
	if data, ok := c.GetValue(key); ok {
		if result, o := data.(bool); o {
			return result, true
		}
		return false, false
	}
	return false, false
}

func Load(file string) (*SimpleConfig, error) {
	c := &SimpleConfig{}

	return c.Load(file)
}

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

func UnmarshalEnv(config interface{}) {
	v := reflection.StructValue(config)
	t := reflection.StructType(config)
	decodeEnvStruct(t, v)
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
	if env.GetString(key) == "" {
		// return if env not set
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
		}
	}

}

func applyDefault(ft reflect.Type, fv reflect.Value, def string) {
	if fv.IsZero() {
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

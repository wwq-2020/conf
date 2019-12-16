package conf

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/wwq1988/errors"

	"github.com/BurntSushi/toml"
	"github.com/hashicorp/consul/api"
)

// ReadFile 读取文件
func ReadFile(file string, obj interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Trace(err)
	}
	return Unmarshal(data, obj)
}

// Unmarshal 反序列化
func Unmarshal(data []byte, obj interface{}) error {
	if err := UnmarshalExt(data, obj, toml.Unmarshal); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// UnmarshalExt 反序列化带反序列化方法
func UnmarshalExt(data []byte, obj interface{}, unmarshaler func([]byte, interface{}) error) error {
	data, err := Render(data)
	if err != nil {
		return errors.Trace(err)
	}
	return unmarshaler(data, obj)
}

// Render 渲染
func Render(data []byte) ([]byte, error) {
	tmpl, err := template.New("conf").Funcs(template.FuncMap{
		"kv": getKV,
	}).Parse(string(data))
	if err != nil {
		return nil, errors.Trace(err)
	}
	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, ""); err != nil {
		return nil, errors.Trace(err)
	}
	return buf.Bytes(), nil
}

func getKV(key string) string {
	cfg := api.DefaultConfig()
	c, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	val, _, err := c.KV().Get(key, nil)
	if err != nil {
		panic(err)
	}
	if val == nil {
		panic(key + " not found")
	}
	return string(val.Value)
}

package conf

import "testing"

type conf struct {
	Hello string `toml:"hello"`
}

func TestReadFile(t *testing.T) {
	conf := &conf{}
	if err := ReadFile("conf.toml", conf); err != nil {
		t.Fatalf("ReadFile expected:nil, got:%+v", err)
	}
}

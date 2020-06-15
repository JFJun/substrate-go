package source

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

type TypeStruct struct {
	Type        string     `json:"type"`
	TypeString  string     `json:"type_string"`
	TypeMapping [][]string `json:"type_mapping,omitempty"`
	ValueList   []string   `json:"value_list,omitempty"`
}

func ReadJson(coinName string)string{
	cc, err := ioutil.ReadFile(fmt.Sprintf("%s.json", coinName))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(cc)
}

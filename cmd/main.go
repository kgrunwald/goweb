package main

import (
	"fmt"

	"github.com/kgrunwald/goweb/xml"
)

type NS1 struct {
	XMLNameSpace string `xml:"ns1=https://ns1.com"`
}

type NS2 struct {
	XMLNameSpace string `xml:"ns2=https://ns2.com"`
}

type Test struct {
	NS1 `xml:"-"`
	XMLName string `xml:"Test"`
	Key     string `xml:"key"`
	Ready   bool `xml:",attr"`
	Nested Nested
}

type Nested struct {
	NS2 `xml:"-"`
	Value string
}

func main() {
	// controller.Register()
	// goweb.Start()

	nested := Nested{Value: "NestedValue"}
	test := Test{Key: "Hi", Ready: false, Nested: nested}
	enc := xml.Encoder{}
	res, _ := enc.Marshal(test)
	fmt.Println(string(res))
}

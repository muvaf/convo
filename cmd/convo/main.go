/*
Copyright 02 November 2019 Muvaffak Onus.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"reflect"
)

type A struct {
	FieldA string `json:"fieldA,omitempty"`
	FieldB *string `json:"fieldB,omitempty"`
}

type B struct {
	FieldA string `json:"fieldA,omitempty"`
	FieldB *string `json:"fieldB,omitempty"`
}

func main() {
	//f := jen.NewFile("main")
	//f.Func().Id("main").Params().Block(
	//	jen.Qual("fmt", "Println").Call(jen.Lit("Hello, world")),
	//)
	//fmt.Printf("%#v", f)
	aMap := map[string]interface{}{}
	bMap := map[string]interface{}{}

	aRef := reflect.ValueOf(A{FieldA: "val1"})
	for i := 0; i < aRef.NumField(); i++ {
		fmt.Printf("fieldname: %s, type: %s, tag: %s, value: %s\n", aRef.Type().Field(i).Name, aRef.Type().Field(i).Type, aRef.Type().Field(i).Tag, aRef.Field(i).Interface())
		aMap[aRef.Type().Field(i).Name] = aRef.Field(i).Interface()
	}

	bRef := reflect.ValueOf(A{FieldA: "val1"})
	for i := 0; i < bRef.NumField(); i++ {
		fmt.Printf("fieldname: %s, type: %s, tag: %s, value: %s\n", bRef.Type().Field(i).Name, bRef.Type().Field(i).Type, bRef.Type().Field(i).Tag, bRef.Field(i).Interface())
		bMap[bRef.Type().Field(i).Name] = bRef.Field(i).Interface()
	}
	var statementList []jen.Code
	for fieldName := range aMap {
		s := jen.Id("target").Dot(fieldName).Op("=").Id("source").Dot(fieldName)
		statementList = append(statementList, s)
	}

	f := jen.NewFile("main")
	f.Func().Id("convert").Params(
		jen.Id("target").Op("*").Id("B"),
		jen.Id("source").Op("*").Id("A"),
		).Block(statementList...)
	fmt.Printf("%#v", f)
}
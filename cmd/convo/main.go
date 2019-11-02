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
	"strings"
)

type A struct {
	FieldA string
	FieldB *string
	Fieldc *string
	FieldD string
}

type B struct {
	FieldA string
	Fieldb *string
	FieldC string
	FieldD *string
}

func main() {
	aRef := reflect.ValueOf(A{})
	aMap := map[string]reflect.StructField{}
	for i := 0; i < aRef.NumField(); i++ {
		aMap[strings.ToLower(aRef.Type().Field(i).Name)] = aRef.Type().Field(i)
	}

	bRef := reflect.ValueOf(B{})
	fieldMap := map[string]reflect.StructField{}
	for i := 0; i < bRef.NumField(); i++ {
		fieldMap[strings.ToLower(bRef.Type().Field(i).Name)] = bRef.Type().Field(i)
	}
	var statementList []jen.Code
	statementList = append(statementList, jen.Id("r").Op(":=").Op("&").Id("B").Values())
	for name, field := range fieldMap {
		// string -> string
		// *string -> *string
		if aMap[name].Type == field.Type {
			statementList = append(statementList, jen.Id("r").Dot(field.Name).Op("=").Id("a").Dot(aMap[name].Name))
		}
		// *string -> string
		if aMap[name].Type.Kind() == reflect.Ptr && field.Type.Kind() != reflect.Ptr {
			s := jen.If(jen.Id("a").Dot(aMap[name].Name).Op("!=").Nil()).Block(
				jen.Id("r").Dot(name).Op("=").Op("*").Id("a").Dot(aMap[name].Name),
				)
			statementList = append(statementList, s)
		}
		// string -> *string
		if aMap[name].Type.Kind() != reflect.Ptr && field.Type.Kind() == reflect.Ptr {
			statementList = append(statementList, jen.Id("r").Dot(field.Name).Op("=").Op("&").Id("a").Dot(aMap[name].Name))
		}
	}
	statementList = append(statementList, jen.Return(jen.Id("r")))

	f := jen.NewFile("main")
	f.Func().Params(jen.Id("a").Id("*A")).Id("ConvertToB").Params().Op("*").Id("B").Block(statementList...)
	fmt.Printf("%#v", f)
}
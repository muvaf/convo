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
	FieldA string
	FieldB *string
	FieldC *string
}

type B struct {
	FieldA string
	FieldB *string
}

func main() {
	bRef := reflect.ValueOf(B{})
	fieldList := make([]string, bRef.NumField())
	for i := 0; i < bRef.NumField(); i++ {
		fieldList[i] = bRef.Type().Field(i).Name
	}
	var statementList []jen.Code
	statementList = append(statementList, jen.Id("r").Op(":=").Op("&").Id("B").Values())
	for _, fieldName := range fieldList {
		s := jen.Id("r").Dot(fieldName).Op("=").Id("a").Dot(fieldName)
		statementList = append(statementList, s)
	}
	statementList = append(statementList, jen.Return(jen.Id("r")))

	f := jen.NewFile("main")
	f.Func().Params(jen.Id("a").Id("*A")).Id("ConvertToB").Params().Op("*").Id("B").Block(statementList...)
	fmt.Printf("%#v", f)
}
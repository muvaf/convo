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
package convo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
)

// BasicConversion converts the classes whose matching fields' are assignable
// to each other no matter whether they are Ptr or not.
// Does not work with anonymous structs.
func BasicConversion(a interface{}, b interface{}) *jen.Statement {
	// a and b has to be non-pointer.
	aRef := reflect.ValueOf(a)
	bRef := reflect.ValueOf(b)
	if aRef.Type().Name() == "" || bRef.Type().Name() == "" {
		panic(fmt.Errorf("convo does not work with anonymous structs"))
	}
	aMap := map[string]reflect.StructField{}
	for i := 0; i < aRef.NumField(); i++ {
		aMap[strings.ToLower(aRef.Type().Field(i).Name)] = aRef.Type().Field(i)
	}
	aReceiverName := string(strings.ToLower(aRef.Type().Name())[0])

	fieldList := make([]reflect.StructField, bRef.NumField())
	for i := 0; i < bRef.NumField(); i++ {
		fieldList[i] = bRef.Type().Field(i)
	}
	statementList := []jen.Code{
		jen.Id("r").Op(":=").Op("&").Id(bRef.Type().Name()).Values(),
	}
	for _, field := range fieldList {
		equatedName := strings.ToLower(field.Name)
		// string -> string
		// *string -> *string
		if aMap[equatedName].Type == field.Type {
			statementList = append(statementList, jen.Id("r").Dot(field.Name).Op("=").Id(aReceiverName).Dot(aMap[equatedName].Name))
		}
		// *string -> string
		if aMap[equatedName].Type.Kind() == reflect.Ptr && field.Type.Kind() != reflect.Ptr {
			s := jen.If(jen.Id(aReceiverName).Dot(aMap[equatedName].Name).Op("!=").Nil()).Block(
				jen.Id("r").Dot(field.Name).Op("=").Op("*").Id(aReceiverName).Dot(aMap[equatedName].Name),
			)
			statementList = append(statementList, s)
		}
		// string -> *string
		if aMap[equatedName].Type.Kind() != reflect.Ptr && field.Type.Kind() == reflect.Ptr {
			statementList = append(statementList, jen.Id("r").Dot(field.Name).Op("=").Op("&").Id(aReceiverName).Dot(aMap[equatedName].Name))
		}
	}
	statementList = append(statementList, jen.Return(jen.Id("r")))

	return jen.Func().
		Params(jen.Id(aReceiverName).Op("*").Id(aRef.Type().Name())).
		Id(fmt.Sprintf("Get%v", bRef.Type().Name())).
		Params().Op("*").Id(bRef.Type().Name()).
		Block(statementList...)
}

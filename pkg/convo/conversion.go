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

// Convert converts the classes whose matching fields' are assignable
// to each other no matter whether they are Ptr or not.
// Does not work with anonymous structs.
func Convert(aRef reflect.Type, bRef reflect.Type) *jen.Statement {
	aReceiverName := string(strings.ToLower(aRef.Name())[0])
	statementList := append(structConversion(aRef, aReceiverName, "r", bRef), jen.Return(jen.Id("r")))
	return jen.Func().
		Params(jen.Id(aReceiverName).Op("*").Id(aRef.Name())).
		Id(fmt.Sprintf("Get%v", bRef.Name())).
		Params().Id(bRef.Name()).
		Block(statementList...)
}

func structConversion(aRef reflect.Type, sourcePath, resultPath string, bType reflect.Type) []jen.Code {
	aMap := map[string]reflect.StructField{}
	for i := 0; i < aRef.NumField(); i++ {
		// TODO: inline fields?
		aMap[strings.ToLower(aRef.Field(i).Name)] = aRef.Field(i)
	}
	// TODO: ptr typed field
	statementList := []jen.Code{jen.Id(resultPath).Op("=").Id(bType.Name()).Values()}
	if !strings.Contains(resultPath, ".") {
		statementList = []jen.Code{jen.Id(resultPath).Op(":=").Id(bType.Name()).Values()}
	}
	fieldList := make([]reflect.StructField, bType.NumField())
	for i := 0; i < bType.NumField(); i++ {
		fieldList[i] = bType.Field(i)
	}
	for _, field := range fieldList {
		statementList = append(statementList, fieldConversion(aMap, sourcePath, resultPath, field)...)
	}
	return statementList
}

func fieldConversion(aMap map[string]reflect.StructField, aSourcePath, resultPath string, field reflect.StructField) []jen.Code {
	var result []jen.Code
	switch field.Type.Kind() {
	case reflect.Struct:
		// TODO: name of the field might be different in result path.
		result = append(result, structConversion(field.Type, fmt.Sprintf("%s.%s", aSourcePath, field.Name), fmt.Sprintf("%s.%s", resultPath, field.Name), field.Type)...)
	default:
		equatedName := strings.ToLower(field.Name)
		// string -> string
		// *string -> *string
		if aMap[equatedName].Type == field.Type {
			result = append(result, jen.Id(resultPath).Dot(field.Name).Op("=").Id(aSourcePath).Dot(aMap[equatedName].Name))
		}
		// *string -> string
		if aMap[equatedName].Type.Kind() == reflect.Ptr && field.Type.Kind() != reflect.Ptr {
			result = append(result, jen.If(jen.Id(aSourcePath).Dot(aMap[equatedName].Name).Op("!=").Nil()).Block(
				jen.Id(resultPath).Dot(field.Name).Op("=").Op("*").Id(aSourcePath).Dot(aMap[equatedName].Name),
			))
		}
		// string -> *string
		if aMap[equatedName].Type.Kind() != reflect.Ptr && field.Type.Kind() == reflect.Ptr {
			result = append(result, jen.Id(resultPath).Dot(field.Name).Op("=").Op("&").Id(aSourcePath).Dot(aMap[equatedName].Name))
		}
	}
	return result
}

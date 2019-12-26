/*
Copyright 03 November 2019 Muvaffak Onus.

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
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type A struct {
	FieldA string
	FieldB *string
	FieldC *string
	FieldD string
	FieldE SubType
}

type B struct {
	FieldA string
	FieldB *string
	FieldC string
	FieldD *string
	FieldE SubType
}

type SubType struct {
	SubFieldA string
}

// TODO(muvaf): improve test cases to be smaller and more contained
func TestBasicConversion(t *testing.T) {

	type args struct {
		from reflect.Type
		to   reflect.Type
	}
	type want struct {
		result string
	}

	cases := map[string]struct {
		args
		want
	}{
		"PointerValueAndFieldNameCaseDifference": {
			args: args{
				from: reflect.TypeOf(A{}),
				to:   reflect.TypeOf(B{}),
			},
			want: want{
				result: `func (a *A) GetB() B {
	r := B{}
	r.FieldA = a.FieldA
	r.FieldB = a.FieldB
	if a.FieldC != nil {
		r.FieldC = *a.FieldC
	}
	r.FieldD = &a.FieldD
	r.FieldE = SubType{}
	r.FieldE.SubFieldA = a.FieldE.SubFieldA
	return r
}`,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result := Convert(tc.args.from, tc.args.to)
			if diff := cmp.Diff(tc.want.result, result.GoString()); diff != "" {
				t.Errorf("Convert() -want, +got:\n%s", diff)
			}
		})
	}
}

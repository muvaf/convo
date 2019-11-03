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
	"testing"

	"github.com/google/go-cmp/cmp"
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

// TODO(muvaf): improve test cases to be smaller and more contained
func TestBasicConversion(t *testing.T) {

	type args struct {
		from interface{}
		to   interface{}
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
				from: A{},
				to:   B{},
			},
			want: want{
				result: `func (a *A) GetB() *B {
	r := &B{}
	r.FieldA = a.FieldA
	r.Fieldb = a.FieldB
	if a.Fieldc != nil {
		r.FieldC = *a.Fieldc
	}
	r.FieldD = &a.FieldD
	return r
}`,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result := BasicConversion(tc.args.from, tc.args.to)
			if diff := cmp.Diff(tc.want.result, result.GoString()); diff != "" {
				t.Errorf("BasicConversion() -want, +got:\n%s", diff)
			}
		})
	}
}

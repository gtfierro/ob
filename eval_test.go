package ob

import (
	"reflect"
	"testing"
)

func TestEvalArray(t *testing.T) {
	for _, test := range []struct {
		op     ArrayOperator
		data   interface{}
		result interface{}
	}{
		{
			ArrayOperator{index: 0, all: false, slice: false},
			[]int{1, 2, 3, 4},
			1,
		},
		{
			ArrayOperator{index: 3, all: false, slice: false},
			[]int{1, 2, 3, 4},
			4,
		},
		{
			ArrayOperator{index: 5, all: false, slice: false},
			[]int{1, 2, 3, 4},
			4,
		},
		{
			ArrayOperator{index: -1, all: false, slice: false},
			[]int{1, 2, 3, 4},
			4,
		},
		{
			ArrayOperator{index: -2, all: false, slice: false},
			[]int{1, 2, 3, 4},
			3,
		},
		{
			ArrayOperator{all: true, slice: false},
			[]int{1, 2, 3, 4},
			[]int{1, 2, 3, 4},
		},
		{
			ArrayOperator{all: false, slice: true, slice_start: 0, slice_end: 4},
			[]int{1, 2, 3, 4},
			[]int{1, 2, 3, 4},
		},
		{
			ArrayOperator{all: false, slice: true, slice_start: 1, slice_end: 4},
			[]int{1, 2, 3, 4},
			[]int{2, 3, 4},
		},
		{
			ArrayOperator{all: false, slice: true, slice_start: 1, slice_end: 40},
			[]int{1, 2, 3, 4},
			[]int{2, 3, 4},
		},
		{
			ArrayOperator{all: false, slice: true, slice_start: 1, slice_end: 2},
			[]int{1, 2, 3, 4},
			[]int{2},
		},
	} {
		res := test.op.Eval(test.data)
		if !reflect.DeepEqual(res, test.result) {
			t.Errorf("Operator %+v on %+v gave %+v but wanted %+v", test.op, test.data, res, test.result)
		}
	}
}

func BenchmarkArrayIndex(b *testing.B) {
	op := ArrayOperator{index: 0, all: false, slice: false}
	data := []uint32{1, 2, 3, 4}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		op.Eval(data)
	}
}

func BenchmarkArraySlice(b *testing.B) {
	op := ArrayOperator{slice_start: 0, slice_end: 4, all: false, slice: true}
	data := []uint32{1, 2, 3, 4}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		op.Eval(data)
	}
}

func BenchmarkArrayAll(b *testing.B) {
	op := ArrayOperator{all: true, slice: false}
	data := []uint32{1, 2, 3, 4}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		op.Eval(data)
	}
}

func TestEvalMap(t *testing.T) {
	for _, test := range []struct {
		op     ObjectOperator
		data   interface{}
		result interface{}
	}{
		{
			ObjectOperator{"key1"},
			map[string]interface{}{"key1": "val1"},
			"val1",
		},
		{
			ObjectOperator{"key1"},
			map[string]interface{}{"key1": 12345},
			12345,
		},
		{
			ObjectOperator{"key1"},
			map[string]interface{}{"key1": []string{"a", "b"}},
			[]string{"a", "b"},
		},
	} {
		res := test.op.Eval(test.data)
		if !reflect.DeepEqual(res, test.result) {
			t.Errorf("Operator %+v on %+v gave %+v but wanted %+v", test.op, test.data, res, test.result)
		}
	}
}

func TestEvalStruct(t *testing.T) {
	for _, test := range []struct {
		op     ObjectOperator
		data   interface{}
		result interface{}
	}{
		{
			ObjectOperator{"Key1"},
			struct{ Key1 string }{Key1: "val1"},
			"val1",
		},
		{
			ObjectOperator{"Key2"},
			struct{ Key2 map[string]string }{Key2: map[string]string{"a": "b"}},
			map[string]string{"a": "b"},
		},
	} {
		res := test.op.Eval(test.data)
		if !reflect.DeepEqual(res, test.result) {
			t.Errorf("Operator %+v on %+v gave %+v but wanted %+v", test.op, test.data, res, test.result)
		}
	}
}

func BenchmarkObjectMap(b *testing.B) {
	op := ObjectOperator{"key1"}
	data := map[string]interface{}{"key1": "val1"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		op.Eval(data)
	}
}
func BenchmarkObjectStruct(b *testing.B) {
	op := ObjectOperator{"Key1"}
	data := struct{ Key1 string }{Key1: "val1"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		op.Eval(data)
	}
}

func TestEvalOperatorChain(t *testing.T) {
	for _, test := range []struct {
		expr   string
		data   interface{}
		result interface{}
	}{
		{
			"[0]",
			[]int{1, 2, 3, 4},
			1,
		},
		{
			"[3]",
			[]int{1, 2, 3, 4},
			4,
		},
		{
			"[5]",
			[]int{1, 2, 3, 4},
			4,
		},
		{
			"[-1]",
			[]int{1, 2, 3, 4},
			4,
		},
		{
			"[-3]",
			[]int{1, 2, 3, 4},
			2,
		},
		{
			"[:]",
			[]int{1, 2, 3, 4},
			[]int{1, 2, 3, 4},
		},
		{
			"[0:4]",
			[]int{1, 2, 3, 4},
			[]int{1, 2, 3, 4},
		},
		{
			"[1:4]",
			[]int{1, 2, 3, 4},
			[]int{2, 3, 4},
		},
		{
			"[1:40]",
			[]int{1, 2, 3, 4},
			[]int{2, 3, 4},
		},
		{
			"[1:2]",
			[]int{1, 2, 3, 4},
			[]int{2},
		},
		{
			"key1",
			map[string]interface{}{"key1": "val1"},
			"val1",
		},
		{
			"key1",
			map[string]interface{}{"key1": 12345},
			12345,
		},
		{
			"key1",
			map[string]interface{}{"key1": []string{"a", "b"}},
			[]string{"a", "b"},
		},
		{
			"key1",
			map[string]interface{}{"key1": "val1"},
			"val1",
		},
		{
			"key1",
			map[string]interface{}{"key1": 12345},
			12345,
		},
		{
			"key1[:]",
			map[string]interface{}{"key1": []string{"a", "b"}},
			[]string{"a", "b"},
		},
		{
			"key1[0]",
			map[string]interface{}{"key1": []string{"a", "b"}},
			"a",
		},
		{
			"[0].key1[1]",
			[]map[string]interface{}{map[string]interface{}{"key1": []string{"a", "b"}}},
			"b",
		},
		{
			"[:].key1",
			[]map[string]interface{}{map[string]interface{}{"key1": 123}, map[string]interface{}{"key1": 456}},
			[]interface{}{123, 456},
		},
		{
			"enter[:].key1",
			map[string]interface{}{"enter": []map[string]interface{}{map[string]interface{}{"key1": 123}, map[string]interface{}{"key1": 456}}},
			[]interface{}{123, 456},
		},
	} {
		res := Eval(Parse(test.expr), test.data)
		if !reflect.DeepEqual(res, test.result) {
			t.Errorf("Expr %+v on %+v gave %+v but wanted %+v", test.expr, test.data, res, test.result)
		}
	}
}

func BenchmarkParseEvalKey(b *testing.B) {
	expr := "key1"
	data := map[string]interface{}{"key1": "val1"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Eval(Parse(expr), data)
	}
}

func BenchmarkParseEvalIndex(b *testing.B) {
	expr := "[5]"
	data := []int{1, 2, 3, 4}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Eval(Parse(expr), data)
	}
}

func BenchmarkParseEvalComplex(b *testing.B) {
	expr := "[0].key1[1]"
	data := []map[string]interface{}{map[string]interface{}{"key1": []string{"a", "b"}}}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Eval(Parse(expr), data)
	}
}

func BenchmarkEvalComplex(b *testing.B) {
	expr := "[0].key1[1]"
	data := []map[string]interface{}{map[string]interface{}{"key1": []string{"a", "b"}}}
	ops := Parse(expr)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Eval(ops, data)
	}
}

func BenchmarkEvalKey(b *testing.B) {
	expr := "key1"
	data := map[string]interface{}{"key1": "val1"}
	ops := Parse(expr)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Eval(ops, data)
	}
}

func BenchmarkEvalIndex(b *testing.B) {
	expr := "[5]"
	data := []int{1, 2, 3, 4}
	ops := Parse(expr)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Eval(ops, data)
	}
}

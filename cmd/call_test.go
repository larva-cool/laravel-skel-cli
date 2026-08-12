package cmd

import (
	"testing"

	"laravel-skel-cli/internal/apidefs"
)

func TestParseValue(t *testing.T) {
	cases := []struct {
		name string
		in   string
		typ  string
		want any
		err  bool
	}{
		{name: "string", in: "abc", typ: "string", want: "abc"},
		{name: "integer", in: "42", typ: "integer", want: int64(42)},
		{name: "number", in: "3.14", typ: "number", want: 3.14},
		{name: "boolean true", in: "1", typ: "boolean", want: true},
		{name: "boolean false", in: "false", typ: "boolean", want: false},
		{name: "array", in: "1,2,3", typ: "array", want: []any{"1", "2", "3"}},
		{name: "integer invalid", in: "abc", typ: "integer", err: true},
		{name: "boolean invalid", in: "maybe", typ: "boolean", err: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseValue(c.in, c.typ)
			if c.err {
				if err == nil {
					t.Fatalf("期望报错，但得到 %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseValue(%q, %q) 报错: %v", c.in, c.typ, err)
			}
			if !equal(got, c.want) {
				t.Fatalf("parseValue(%q, %q) = %#v, 期望 %#v", c.in, c.typ, got, c.want)
			}
		})
	}
}

func TestIsRawBody(t *testing.T) {
	raw := &apidefs.Endpoint{BodyParams: []apidefs.Param{{Name: "body"}}}
	if !isRawBody(raw) {
		t.Fatal("含 body 原始参数应返回 true")
	}
	obj := &apidefs.Endpoint{BodyParams: []apidefs.Param{{Name: "prompt"}}}
	if isRawBody(obj) {
		t.Fatal("含普通字段应返回 false")
	}
	empty := &apidefs.Endpoint{}
	if isRawBody(empty) {
		t.Fatal("无请求体应返回 false")
	}
}

func equal(a, b any) bool {
	switch av := a.(type) {
	case []any:
		bv, ok := b.([]any)
		if !ok || len(av) != len(bv) {
			return false
		}
		for i := range av {
			if av[i] != bv[i] {
				return false
			}
		}
		return true
	default:
		return av == b
	}
}

package converter

import (
	"context"
	"reflect"
	"testing"

	rstypes "github.com/drewstone/go2rs/pkg/types"
)

func expect(t *testing.T, actual rstypes.Type, expected string) {
	t.Helper()
	if actual.String() != expected {
		t.Fatalf("expected \"%s\" to equal \"%s\"", expected, actual.String())
	}
}

func typ(i interface{}) reflect.Type {
	return reflect.TypeOf(i)
}

type User struct{ Name string }
type Nested struct{ Owner User }

func (*Nested) Method(arg string) {}

func TestPrimitives(t *testing.T) {
	c := NewConverter()
	expect(t, c.Convert(typ("")), "String")
	expect(t, c.Convert(typ(138)), "i32")
	expect(t, c.Convert(typ(int64(138))), "i64")

}

func TestStructs(t *testing.T) {
	c := NewConverter()
	expect(t, c.Convert(typ(User{})), "User")
	expect(t, c.Convert(typ(&User{})), "Option<User>")
	expect(t, c.Convert(typ(Nested{})), "Nested")
}

func TestFuncs(t *testing.T) {
	c := NewConverter()
	// params
	expect(t, c.Convert(typ(func() {})), "fn()")
	expect(t, c.Convert(typ(func(string, int, bool) {})), "fn(s: String, n: i32, b: bool)")
	expect(t, c.Convert(typ(func(struct{}) {})), "fn(empty: ())")
	expect(t, c.Convert(typ(func(struct{ ID string }) {})), "fn(data: struct { ID: String })")
	expect(t, c.Convert(typ(func(User) {})), "fn(user: User)")
	expect(t, c.Convert(typ(func(*User) {})), "fn(user: Option<User>)")
	expect(t, c.Convert(typ(func(Nested) {})), "fn(nested: Nested)")

	// returns
	expect(t, c.Convert(typ(func() string { return "foo" })), "fn() -> String")
	expect(t, c.Convert(typ(func() User { return User{} })), "fn() -> User")
	expect(t, c.Convert(typ(func() *User { return nil })), "fn() -> Option<User>")
	expect(t, c.Convert(typ(func() error { return nil })), "fn() -> Result<(), String>")
	expect(t, c.Convert(typ(func() (string, error) { return "", nil })), "fn() -> Result<String, String>")
	expect(t, c.Convert(typ(func() (string, string, error) { return "", "", nil })),
		"fn() -> Result<(String, String), String>")
	expect(t, c.Convert(typ(func() chan string { return nil })), "fn() -> Vec<String>")
	expect(t, c.Convert(typ(func(context.Context, string) {})), "fn(context: Context, s: String)")

	// methods
	n := Nested{}
	m, _ := typ(&n).MethodByName("Method")
	c.ConfigureFunc = func(t reflect.Type) FuncConfig {
		return FuncConfig{IsMethod: true, MethodName: m.Name}
	}
	expect(t, c.Convert(m.Type), "fn Method(&self, s: String)")

	// configurations
	c.ConfigureFunc = func(t reflect.Type) FuncConfig {
		return FuncConfig{ParamNames: []string{"foo", "bar", "baz"}}
	}
	expect(t, c.Convert(typ(func(string, int, bool) {})),
		"fn(foo: String, bar: i32, baz: bool)")

	c.ConfigureFunc = func(t reflect.Type) FuncConfig {
		return FuncConfig{IsAsync: true}
	}
	expect(t, c.Convert(typ(func() string { return "" })), "async fn() -> String")

	c.ConfigureFunc = func(t reflect.Type) FuncConfig {
		return FuncConfig{NoIgnoreContext: true}
	}
	expect(t, c.Convert(typ(func(ctx context.Context) {})), "fn(ctx: Context)")
}

func TestSlices(t *testing.T) {
	c := NewConverter()
	expect(t, c.Convert(typ([]string{})), "Vec<String>")
	expect(t, c.Convert(typ([]*User{})), "Vec<Option<User>>")
}

func TestArrays(t *testing.T) {
	c := NewConverter()
	expect(t, c.Convert(typ([2]string{"first", "second"})), "[String; 2]")
}

func TestMaps(t *testing.T) {
	c := NewConverter()
	expect(t, c.Convert(typ(map[string]int{})), "HashMap<String, i32>")
	expect(t, c.Convert(typ(map[string]User{})), "HashMap<String, User>")
}

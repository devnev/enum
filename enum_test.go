package enum_test

import (
	"testing"

	"github.com/devnev/enum"
	"github.com/matryer/is"
)

type Color enum.Member[string]

var (
	Red    = Color{"red"}
	Green  = Color{"green"}
	Blue   = Color{"blue"}
	Colors = enum.New(Red, Green, Blue)
)

func TestMember_Value(t *testing.T) {
	is := is.New(t)
	is.Equal(Red.Value, "red")
	is.Equal(Green.Value, "green")
	is.Equal(Blue.Value, "blue")
	is.Equal(enum.Member[string]{"blue"}.Value, "blue")
	is.Equal(enum.Member[int]{14}.Value, 14)
}

func TestEnum_Parse(t *testing.T) {
	is := is.New(t)
	parsed, ok := Colors.Parse("red")
	is.True(ok)
	is.Equal(parsed, Red)
	parsed, ok = Colors.Parse("purple")
	is.True(!ok)
	is.Equal(parsed, Color{})
}

func TestEnum_Empty(t *testing.T) {
	is := is.New(t)
	is.True(!Colors.Empty())
	is.True(enum.New[int, enum.Member[int]]().Empty())
}

func TestEnum_Len(t *testing.T) {
	is := is.New(t)
	is.Equal(Colors.Len(), 3)
	is.Equal(enum.New[int, enum.Member[int]]().Len(), 0)
}

func TestEnum_Contains(t *testing.T) {
	is := is.New(t)
	is.True(Colors.Contains(Red))
	is.True(Colors.Contains(Green))
	is.True(Colors.Contains(Blue))
	blue := Color{"blue"}
	is.True(Colors.Contains(blue))
	purple := Color{"purple"}
	is.True(!Colors.Contains(purple))
}

func TestEnum_Members(t *testing.T) {
	is := is.New(t)
	exp := []Color{Red, Green, Blue}
	is.Equal(Colors.Members(), exp)
}

func TestEnum_Values(t *testing.T) {
	is := is.New(t)
	exp := []string{"red", "green", "blue"}
	is.Equal(Colors.Values(), exp)
}

func TestEnum_Value(t *testing.T) {
	is := is.New(t)
	is.Equal(Colors.Value(Red), "red")
}

func TestEnum_Index(t *testing.T) {
	is := is.New(t)
	is.Equal(Colors.Index(Red), 0)
	is.Equal(Colors.Index(Green), 1)
	is.Equal(Colors.Index(Blue), 2)
}

func TestEnum_Index_Panic(t *testing.T) {
	is := is.New(t)
	defer func() {
		r := recover()
		is.Equal(r, "the given Member does not belong to this Enum")
	}()
	Colors.Index(Color{"purple"})
}

func TestBuilder(t *testing.T) {
	is := is.New(t)
	type Country enum.Member[string]
	var (
		b         = enum.NewBuilder[string, Country]()
		NL        = b.Add(Country{"Netherlands"})
		FR        = b.Add(Country{"France"})
		BE        = b.Add(Country{"Belgium"})
		Countries = b.Enum()
	)
	is.Equal(Countries.Members(), []Country{NL, FR, BE})
}

package util

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestDecompose(t *testing.T) {
	assert := assert.New(t)

	myObj := testType{}
	myObj.ID = 123
	myObj.Name = "Test Object"
	myObj.NotTagged = "Not Tagged"
	myObj.Tagged = "Is Tagged"
	myObj.SubTypes = append([]subType{}, subType{1, "One"})
	myObj.SubTypes = append(myObj.SubTypes, subType{2, "Two"})
	myObj.SubTypes = append(myObj.SubTypes, subType{3, "Three"})
	myObj.SubTypes = append(myObj.SubTypes, subType{4, "Four"})

	decomposed := Reflection.Decompose(myObj)

	_, hasKey := decomposed["ID"]
	assert.True(hasKey)

	_, hasKey = decomposed["Name"]
	assert.True(hasKey)

	_, hasKey = decomposed["NotTagged"]
	assert.True(hasKey)

	_, hasKey = decomposed["Tagged"]
	assert.True(hasKey)

	children, hasKey := decomposed["SubTypes"]
	assert.True(hasKey)
	assert.Len(children, 4)
}

func TestDecomposeStrings(t *testing.T) {
	assert := assert.New(t)

	input := mapStringsTest{
		Float64:  3.14,
		String:   "foo",
		Base64:   []byte("this is base64"),
		Bytes:    []byte("this is bytes"),
		CSV:      []string{"foo", "bar", "baz"},
		Duration: 10 * time.Second,
		Sub: mapStringsTestSubObject{
			Foo: "yes this is foo",
		},
	}

	output := Reflection.DecomposeStrings(input)
	assert.NotEmpty(output)

	assert.Equal("3.14", output["Float64"])
	assert.Equal("foo", output["String"])
	assert.NotEmpty(output["Base64"])
	assert.NotEqual("this is base64", output["Base64"])
	assert.Equal("[116 104 105 115 32 105 115 32 98 121 116 101 115]", output["Bytes"])
	assert.Equal("10s", output["Duration"])
	assert.Equal("yes this is foo", output["Foo"])
}

func TestDecomposeStringsWithTag(t *testing.T) {
	assert := assert.New(t)

	input := mapStringsTest{
		Float64:  3.14,
		String:   "foo",
		Base64:   []byte("this is base64"),
		Bytes:    []byte("this is bytes"),
		CSV:      []string{"foo", "bar", "baz"},
		Duration: 10 * time.Second,
		Sub: mapStringsTestSubObject{
			Foo: "yes this is foo",
		},
	}

	output := Reflection.DecomposeStrings(input, "secret")
	assert.NotEmpty(output)

	assert.Equal("3.14", output["float64"])
	assert.Equal("foo", output["string"])
	assert.NotEmpty(output["base64Field"])
	assert.NotEqual("this is base64", output["base64Field"])
	assert.Equal("this is bytes", output["bytesField"])
	assert.Equal("10s", output["duration"])
	assert.Equal("yes this is foo", output["foo"])
}

func TestPatchObject(t *testing.T) {
	assert := assert.New(t)

	myObj := testType{}
	myObj.ID = 123
	myObj.Name = "Test Object"
	myObj.NotTagged = "Not Tagged"
	myObj.Tagged = "Is Tagged"
	myObj.SubTypes = append([]subType{}, subType{1, "One"})
	myObj.SubTypes = append(myObj.SubTypes, subType{2, "Two"})
	myObj.SubTypes = append(myObj.SubTypes, subType{3, "Three"})
	myObj.SubTypes = append(myObj.SubTypes, subType{4, "Four"})

	patchData := make(map[string]interface{})
	patchData["Tagged"] = "Is Not Tagged"

	err := Reflection.Patch(&myObj, patchData)
	assert.Nil(err)
	assert.Equal("Is Not Tagged", myObj.Tagged)
}

func testCachedObject(obj interface{}) func() interface{} {
	return func() interface{} {
		return obj
	}
}

func TestReflectTypeInterface(t *testing.T) {
	assert := assert.New(t)

	proto := testCachedObject(testObject{ID: 1, Name: "Test"})

	assert.NotNil(proto())

	objType := Reflection.Type(proto())
	assert.NotNil(objType)
}

func TestReflectValueInterface(t *testing.T) {
	assert := assert.New(t)

	proto := testCachedObject(&testObject{ID: 1, Name: "Test"})

	assert.NotNil(proto())

	objValue := Reflection.Value(proto())
	assert.NotNil(objValue)
	assert.True(objValue.CanSet())
}

func TestPatchStrings(t *testing.T) {
	assert := assert.New(t)

	var mule mapStringsTest

	// ----
	// bool
	// ----

	boolValid := map[string]string{
		"bool": "true",
	}
	boolInvalid := map[string]string{
		"bool": "nottrue",
	}
	assert.Nil(Reflection.PatchStrings("secret", boolValid, &mule))
	assert.Equal(true, mule.Bool)
	assert.Nil(Reflection.PatchStrings("secret", boolInvalid, &mule))
	assert.Equal(false, mule.Bool)

	// -------
	// float32
	// -------

	float32Valid := map[string]string{
		"float32": "3.14",
	}
	float32Invalid := map[string]string{
		"float32": "random",
	}
	assert.Nil(Reflection.PatchStrings("secret", float32Valid, &mule))
	assert.Equal(3.14, mule.Float32)
	assert.NotNil(Reflection.PatchStrings("secret", float32Invalid, &mule))

	// -------
	// float64
	// -------

	float64Valid := map[string]string{
		"float64": "6.28",
	}
	float64Invalid := map[string]string{
		"float64": "random",
	}
	assert.Nil(Reflection.PatchStrings("secret", float64Valid, &mule))
	assert.Equal(6.28, mule.Float64)
	assert.NotNil(Reflection.PatchStrings("secret", float64Invalid, &mule))

	// -------
	// int8
	// -------

	int8Valid := map[string]string{
		"int8": "3",
	}
	int8Invalid := map[string]string{
		"int8": "random",
	}
	assert.Nil(Reflection.PatchStrings("secret", int8Valid, &mule))
	assert.Equal(3, mule.Int8)
	assert.NotNil(Reflection.PatchStrings("secret", int8Invalid, &mule))
	assert.Equal(3, mule.Int8)

	// -------
	// int16
	// -------

	int16Valid := map[string]string{
		"int16": "4",
	}
	int16Invalid := map[string]string{
		"int16": "random",
	}
	assert.Nil(Reflection.PatchStrings("secret", int16Valid, &mule))
	assert.Equal(4, mule.Int16)
	assert.NotNil(Reflection.PatchStrings("secret", int16Invalid, &mule))
	assert.Equal(4, mule.Int16)

	// -------
	// int32
	// -------

	int32Valid := map[string]string{
		"int32": "5",
	}
	int32Invalid := map[string]string{
		"int32": "random",
	}
	assert.Nil(Reflection.PatchStrings("secret", int32Valid, &mule))
	assert.Equal(5, mule.Int32)
	assert.NotNil(Reflection.PatchStrings("secret", int32Invalid, &mule))
	assert.Equal(5, mule.Int32)

	// -------
	// int
	// -------

	intValid := map[string]string{
		"int": "6",
	}
	intInvalid := map[string]string{
		"int": "random",
	}
	assert.Nil(Reflection.PatchStrings("secret", intValid, &mule))
	assert.Equal(6, mule.Int)
	assert.NotNil(Reflection.PatchStrings("secret", intInvalid, &mule))
	assert.Equal(6, mule.Int)

	// -------
	// int64
	// -------

	int64Valid := map[string]string{
		"int64": "7",
	}
	int64Invalid := map[string]string{
		"int64": "random",
	}
	assert.Nil(Reflection.PatchStrings("secret", int64Valid, &mule))
	assert.Equal(7, mule.Int64)
	assert.NotNil(Reflection.PatchStrings("secret", int64Invalid, &mule))
	assert.Equal(7, mule.Int64)

	// -------
	// uint8
	// -------

	uint8Valid := map[string]string{
		"uint8": "8",
	}
	uint8Invalid := map[string]string{
		"uint8": "random",
	}
	assert.Nil(Reflection.PatchStrings("secret", uint8Valid, &mule))
	assert.Equal(8, mule.Uint8)
	assert.NotNil(Reflection.PatchStrings("secret", uint8Invalid, &mule))
	assert.Equal(8, mule.Uint8)

	// -------
	// uint16
	// -------

	uint16Valid := map[string]string{
		"uint16": "9",
	}
	uint16Invalid := map[string]string{
		"uint16": "random",
	}
	assert.Nil(Reflection.PatchStrings("secret", uint16Valid, &mule))
	assert.Equal(9, mule.Uint16)
	assert.NotNil(Reflection.PatchStrings("secret", uint16Invalid, &mule))
	assert.Equal(9, mule.Uint16)

	// -------
	// uint32
	// -------

	uint32Valid := map[string]string{
		"uint32": "10",
	}
	uint32Invalid := map[string]string{
		"uint32": "random",
	}
	assert.Nil(Reflection.PatchStrings("secret", uint32Valid, &mule))
	assert.Equal(10, mule.Uint32)
	assert.NotNil(Reflection.PatchStrings("secret", uint32Invalid, &mule))
	assert.Equal(10, mule.Uint32)

	// -------
	// uint64
	// -------

	uint64Valid := map[string]string{
		"uint64": "11",
	}
	uint64Invalid := map[string]string{
		"uint64": "random",
	}
	assert.Nil(Reflection.PatchStrings("secret", uint64Valid, &mule))
	assert.Equal(11, mule.Uint64)
	assert.NotNil(Reflection.PatchStrings("secret", uint64Invalid, &mule))
	assert.Equal(11, mule.Uint64)

	// -------
	// string
	// -------

	stringValid := map[string]string{
		"string": "foo",
	}
	assert.Nil(Reflection.PatchStrings("secret", stringValid, &mule))
	assert.Equal("foo", mule.String)

	// -------
	// duration
	// -------

	durationValid := map[string]string{
		"duration": "10s",
	}
	durationInvalid := map[string]string{
		"duration": "random",
	}
	assert.Nil(Reflection.PatchStrings("secret", durationValid, &mule))
	assert.Equal(10*time.Second, mule.Duration)
	assert.NotNil(Reflection.PatchStrings("secret", durationInvalid, &mule))
	assert.Equal(10*time.Second, mule.Duration)

	// -------
	// csv
	// -------

	csvValid := map[string]string{
		"csvField": "foo,bar,baz",
	}
	assert.Nil(Reflection.PatchStrings("secret", csvValid, &mule))
	assert.Len(mule.CSV, 3)
	assert.Equal([]string{"foo", "bar", "baz"}, mule.CSV)

	// -------
	// base64
	// -------

	base64Valid := map[string]string{
		"base64Field": base64.StdEncoding.EncodeToString([]byte("this is only a test")),
	}
	base64Invalid := map[string]string{
		"base64Field": "thisisnonsense",
	}
	assert.Nil(Reflection.PatchStrings("secret", base64Valid, &mule))
	assert.Equal("this is only a test", string(mule.Base64))
	assert.NotNil(Reflection.PatchStrings("secret", base64Invalid, &mule))
	assert.Equal("this is only a test", string(mule.Base64))

	// -------
	// bytes
	// -------

	bytesValid := map[string]string{
		"bytesField": "this is bytes",
	}
	assert.Nil(Reflection.PatchStrings("secret", bytesValid, &mule))
	assert.Equal("this is bytes", string(mule.Bytes))

	// -------
	// child objects
	// -------

	childValid := map[string]string{
		"foo": "this is foo",
	}
	assert.Nil(Reflection.PatchStrings("secret", childValid, &mule))
	assert.Equal("this is foo", string(mule.Sub.Foo))
}

type subType struct {
	ID   int
	Name string
}

type testObject struct {
	ID   int
	Name string
}

type testType struct {
	ID        int
	Name      string
	NotTagged string
	Tagged    string
	SubTypes  []subType
}

type mapStringsTest struct {
	Bool     bool          `secret:"bool"`
	Float32  float32       `secret:"float32"`
	Float64  float64       `secret:"float64"`
	Int8     int8          `secret:"int8"`
	Int16    int16         `secret:"int16"`
	Int32    int32         `secret:"int32"`
	Int      int           `secret:"int"`
	Int64    int64         `secret:"int64"`
	Uint8    uint8         `secret:"uint8"`
	Uint16   uint16        `secret:"uint16"`
	Uint32   uint32        `secret:"uint32"`
	Uint64   uint32        `secret:"uint64"`
	String   string        `secret:"string"`
	Duration time.Duration `secret:"duration"`

	CSV    []string `secret:"csvField,csv"`
	Base64 []byte   `secret:"base64Field,base64"`
	Bytes  []byte   `secret:"bytesField,bytes"`

	Sub mapStringsTestSubObject

	unexported string `secret:"unexported"`
}

type mapStringsTestSubObject struct {
	Foo string `secret:"foo"`
}

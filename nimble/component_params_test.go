package nimble

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestComponentParams_parse_empty_uri(t *testing.T) {
	params, err := ParseParams("")
	assert.EqualError(t, err, "uri can not be empty")
	assert.Nil(t, params)
}

func TestComponentParams_parse_empty_schema_uri(t *testing.T) {
	params, err := ParseParams("abc")
	assert.EqualError(t, err, "scheme can not be empty")
	assert.Nil(t, params)
}

func TestComponentParams_parse_empty_name_uri(t *testing.T) {
	params, err := ParseParams("abc://")
	assert.EqualError(t, err, "name can not be empty")
	assert.Nil(t, params)
}

func TestComponentParams_parse_uri(t *testing.T) {
	params, err := ParseParams("abc://test")
	assert.Nil(t, err)
	assert.Equal(t, params.Uri, "abc://test")
	assert.Equal(t, params.Key, "abc")
	assert.Equal(t, params.Name, "test")
}

func TestComponentParams_parse_uri_with_query(t *testing.T) {
	params, err := ParseParams("abc://test?test=abc")
	assert.Nil(t, err)
	assert.Equal(t, params.Uri, "abc://test?test=abc")
	assert.Equal(t, params.Key, "abc")
	assert.Equal(t, params.Name, "test")
}

func TestComponentParams_get_id(t *testing.T) {
	params, err := ParseParams("abc://test")
	assert.Nil(t, err)
	assert.Equal(t, params.GetId(), "abc://test?id=1")
}

func TestComponentParams_get_id_with_custom(t *testing.T) {
	params, err := ParseParams("abc://test?routeId=testing")
	assert.Nil(t, err)
	assert.Equal(t, params.GetId(), "abc://test?id=testing")
}

func TestComponentParams_get_id_with_other_queries(t *testing.T) {
	params, err := ParseParams("abc://test?test=abc")
	assert.Nil(t, err)
	assert.Equal(t, params.GetId(), "abc://test?id=1")
}

func TestComponentParams_string(t *testing.T) {
	params, err := ParseParams("abc://test?test=abc")
	assert.Nil(t, err)
	val, err := params.String("test")
	assert.Nil(t, err)
	assert.Equal(t, val, "abc")
}

func TestComponentParams_string_no_val(t *testing.T) {
	params, err := ParseParams("abc://test")
	assert.Nil(t, err)
	val, err := params.String("test")
	assert.EqualError(t, err, "key not found")
	assert.Equal(t, val, "")
}

func TestComponentParams_string_def(t *testing.T) {
	params, err := ParseParams("abc://test?test=abc")
	assert.Nil(t, err)
	val := params.StringDef("test", "abc")
	assert.Equal(t, val, "abc")
}

func TestComponentParams_string_def_no_val(t *testing.T) {
	params, err := ParseParams("abc://test")
	assert.Nil(t, err)
	val := params.StringDef("test", "abc")
	assert.Equal(t, val, "abc")
}

func TestComponentParams_int(t *testing.T) {
	params, err := ParseParams("abc://test?test=123")
	assert.Nil(t, err)
	val, err := params.Int("test")
	assert.Nil(t, err)
	assert.Equal(t, val, 123)
}

func TestComponentParams_int_no_val(t *testing.T) {
	params, err := ParseParams("abc://test")
	assert.Nil(t, err)
	val, err := params.Int("test")
	assert.EqualError(t, err, "key not found")
	assert.Equal(t, val, 0)
}

func TestComponentParams_int_invalid_val(t *testing.T) {
	params, err := ParseParams("abc://test?test=abc")
	assert.Nil(t, err)
	val, err := params.Int("test")
	assert.EqualError(t, err, "strconv.Atoi: parsing \"abc\": invalid syntax")
	assert.Equal(t, val, 0)
}

func TestComponentParams_int_def(t *testing.T) {
	params, err := ParseParams("abc://test?test=123")
	assert.Nil(t, err)
	val := params.IntDef("test", 123)
	assert.Equal(t, val, 123)
}

func TestComponentParams_int_def_no_val(t *testing.T) {
	params, err := ParseParams("abc://test")
	assert.Nil(t, err)
	val := params.IntDef("test", 123)
	assert.Equal(t, val, 123)
}

func TestComponentParams_int_def_invalid_val(t *testing.T) {
	params, err := ParseParams("abc://test?test=abc")
	assert.Nil(t, err)
	val := params.IntDef("test", 123)
	assert.Equal(t, val, 123)
}

func TestComponentParams_duration(t *testing.T) {
	params, err := ParseParams("abc://test?test=1s")
	assert.Nil(t, err)
	val, err := params.Duration("test")
	assert.Nil(t, err)
	assert.Equal(t, val, time.Second)
}

func TestComponentParams_duration_no_val(t *testing.T) {
	params, err := ParseParams("abc://test")
	assert.Nil(t, err)
	val, err := params.Duration("test")
	assert.EqualError(t, err, "key not found")
	assert.Equal(t, val, time.Nanosecond)
}

func TestComponentParams_duration_invalid_val(t *testing.T) {
	params, err := ParseParams("abc://test?test=abc")
	assert.Nil(t, err)
	val, err := params.Duration("test")
	assert.EqualError(t, err, "time: invalid duration \"abc\"")
	assert.Equal(t, val, time.Nanosecond)
}

func TestComponentParams_duration_def(t *testing.T) {
	params, err := ParseParams("abc://test?test=1s")
	assert.Nil(t, err)
	val := params.DurationDef("test", time.Second)
	assert.Equal(t, val, time.Second)
}

func TestComponentParams_duration_def_no_val(t *testing.T) {
	params, err := ParseParams("abc://test")
	assert.Nil(t, err)
	val := params.DurationDef("test", time.Second)
	assert.Equal(t, val, time.Second)
}

func TestComponentParams_duration_def_invalid_val(t *testing.T) {
	params, err := ParseParams("abc://test?test=abc")
	assert.Nil(t, err)
	val := params.DurationDef("test", time.Second)
	assert.Equal(t, val, time.Second)
}

func TestComponentParams_bool(t *testing.T) {
	params, err := ParseParams("abc://test?test=true")
	assert.Nil(t, err)
	val, err := params.Bool("test")
	assert.Nil(t, err)
	assert.Equal(t, val, true)
}

func TestComponentParams_bool_no_val(t *testing.T) {
	params, err := ParseParams("abc://test")
	assert.Nil(t, err)
	val, err := params.Bool("test")
	assert.EqualError(t, err, "key not found")
	assert.Equal(t, val, false)
}

func TestComponentParams_bool_invalid_val(t *testing.T) {
	params, err := ParseParams("abc://test?test=abc")
	assert.Nil(t, err)
	val, err := params.Bool("test")
	assert.EqualError(t, err, "strconv.ParseBool: parsing \"abc\": invalid syntax")
	assert.Equal(t, val, false)
}

func TestComponentParams_bool_def(t *testing.T) {
	params, err := ParseParams("abc://test?test=true")
	assert.Nil(t, err)
	val := params.BoolDef("test", true)
	assert.Equal(t, val, true)
}

func TestComponentParams_bool_def_no_val(t *testing.T) {
	params, err := ParseParams("abc://test")
	assert.Nil(t, err)
	val := params.BoolDef("test", true)
	assert.Equal(t, val, true)
}

func TestComponentParams_bool_def_invalid_val(t *testing.T) {
	params, err := ParseParams("abc://test?test=abc")
	assert.Nil(t, err)
	val := params.BoolDef("test", true)
	assert.Equal(t, val, true)
}

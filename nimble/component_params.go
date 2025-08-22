package nimble

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type ComponentParams struct {
	Uri    string
	Key    string
	Name   string
	Values url.Values
}

func ParseParams(uri string) (*ComponentParams, error) {
	if uri == "" {
		return nil, errors.New("uri can not be empty")
	}

	parse, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if parse.Scheme == "" {
		return nil, errors.New("scheme can not be empty")
	}

	name := parse.Opaque + parse.Host + parse.Path
	if name == "" {
		return nil, errors.New("name can not be empty")
	}

	query, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return nil, err
	}

	return &ComponentParams{
		Uri:    uri,
		Key:    parse.Scheme,
		Name:   name,
		Values: query,
	}, nil
}

func (m *ComponentParams) GetId() string {
	id := m.StringDef("routeId", "1")
	return fmt.Sprintf("%s://%s?id=%s", m.Key, m.Name, id)
}

func (m *ComponentParams) String(key string) (string, error) {
	if m.Values.Has(key) {
		return m.Values.Get(key), nil
	}
	return "", errors.New("key not found")
}

func (m *ComponentParams) StringDef(key string, defVal string) string {
	if m.Values.Has(key) {
		return m.Values.Get(key)
	}
	return defVal
}

func (m *ComponentParams) Int(key string) (int, error) {
	val, err := m.String(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (m *ComponentParams) IntDef(key string, defVal int) int {
	val, err := m.String(key)
	if err != nil {
		return defVal
	}
	out, err := strconv.Atoi(val)
	if err != nil {
		return defVal
	}
	return out
}

func (m *ComponentParams) Duration(key string) (time.Duration, error) {
	val, err := m.String(key)
	if err != nil {
		return time.Nanosecond, err
	}

	dur, err := time.ParseDuration(val)
	if err != nil {
		return time.Nanosecond, err
	}

	return dur, nil
}

func (m *ComponentParams) DurationDef(key string, defVal time.Duration) time.Duration {
	val, err := m.String(key)
	if err != nil {
		return defVal
	}

	dur, err := time.ParseDuration(val)
	if err != nil {
		return defVal
	}

	return dur
}

func (m *ComponentParams) Bool(key string) (bool, error) {
	val, err := m.String(key)
	if err != nil {
		return false, err
	}

	parseBool, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}

	return parseBool, nil
}

func (m *ComponentParams) BoolDef(key string, defVal bool) bool {
	val := m.StringDef(key, strconv.FormatBool(defVal))
	parseBool, err := strconv.ParseBool(val)
	if err != nil {
		return defVal
	}
	return parseBool
}

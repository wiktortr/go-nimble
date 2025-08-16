package nimble

import (
	"errors"
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
	parse, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	query, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return nil, err
	}

	return &ComponentParams{
		Uri:    uri,
		Key:    parse.Scheme,
		Name:   parse.Opaque,
		Values: query,
	}, nil
}

func (m *ComponentParams) String(key string) (string, error) {
	if m.Values.Has(key) {
		return m.Values.Get(key), nil
	}
	return "", errors.New("key not found")
}

func (m *ComponentParams) IntDef(key string, defVal int) (int, error) {
	val, err := m.String(key)
	if err != nil {
		return defVal, nil
	}
	return strconv.Atoi(val)
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

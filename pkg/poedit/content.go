package poedit

import "errors"

var ErrNoExtensionMetaFound = errors.New("no extension meta found for specified format")

type ContentMeta struct {
	Extension, Type string
}

var ContentMetaMap = map[string]ContentMeta{
	"pot":             {Extension: "pot", Type: "text/plain; charset=utf-8"},
	"po":              {Extension: "po", Type: "text/plain; charset=utf-8"},
	"mo":              {Extension: "mo", Type: "application/octet-stream"},
	"xls":             {Extension: "xls", Type: "application/vnd.ms-excel"},
	"xlsx":            {Extension: "xlsx", Type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
	"csv":             {Extension: "csv", Type: "text/csv; charset=utf-8"},
	"ini":             {Extension: "ini", Type: "text/plain; charset=utf-8"},
	"resw":            {Extension: "resw", Type: "application/xml"},
	"resx":            {Extension: "resx", Type: "application/xml"},
	"android_strings": {Extension: "xml", Type: "application/xml"},
	"apple_strings":   {Extension: "strings", Type: "text/plain; charset=utf-8"},"xliff":           {Extension: "xliff", Type: "application/xml"},
	"properties":      {Extension: "properties", Type: "text/plain; charset=utf-8"},
	"key_value_json":  {Extension: "json", Type: "application/json"},
	"json":            {Extension: "json", Type: "application/json"},
	"yml":             {Extension: "yml", Type: "text/plain; charset=utf-8"},
	"xmb":             {Extension: "xmb", Type: "application/xml"},
	"xtb":             {Extension: "xtb", Type: "application/xml"},
	"arb":             {Extension: "arb", Type: "application/json"},
}

func GetContentMeta(format string) (*ContentMeta, error) {
	m, ok := ContentMetaMap[format]
	if !ok {
		return nil, ErrNoExtensionMetaFound
	}

	return &m, nil
}

func GetContentMetaMap() map[string]ContentMeta {
	return ContentMetaMap
}

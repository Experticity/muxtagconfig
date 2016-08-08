package muxtagconfig

import (
	"net/http"
	"net/url"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/Experticity/tagconfig"
)

// ParseMuxRequestToStruct is a convenience method so to wrap the operations of creating a MuxURLGetter and running
// tagconfig.Process
func ParseMuxRequestToStruct(req *http.Request, container interface{}) error {
	return tagconfig.Process(&MuxURLGetter{req}, container)
}

// MuxURLGetter Implements the TagValueGetter to parse struct tagconfig into a struct instance
type MuxURLGetter struct {
	*http.Request
}

func (mg *MuxURLGetter) TagName() string {
	return "mux.url"
}

// Get will be called from tagconfig.Process for any fields with the tag mux.url and consult secondary fields
// mux.param, mux.form and mux.path to retrieve the appropriate value for the struct field
func (mg *MuxURLGetter) Get(key string, f reflect.StructField) (v string) {
	if mg.Request == nil {
		return ""
	}

	// To avoid shadowing, pre-declaring
	var ok bool
	if f.Tag.Get("mux.param") != "" {
		// QueryParams
		v, ok = tryURLValues(key, mg.Request.URL.Query())

		if ok {
			return
		}
	}

	if f.Tag.Get("mux.form") != "" {
		// PostForm Values
		v, ok = tryURLValues(key, mg.Request.PostForm)

		if ok {
			return
		}

		// Form Values
		v, ok = tryURLValues(key, mg.Request.Form)
		if ok {
			return
		}
	}

	if f.Tag.Get("mux.path") != "" {
		vs := mux.Vars(mg.Request)
		return vs[key]
	}

	return ""
}

func tryURLValues(key string, vs url.Values) (v string, present bool) {
	v = vs.Get(key)
	present = v != ""

	return
}

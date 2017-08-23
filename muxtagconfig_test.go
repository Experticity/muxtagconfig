package muxtagconfig

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Experticity/tagconfig"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type HackerReferences struct {
	Handle        string `mux.url:"handle" mux.form:"true"`
	SuperComputer string `mux.url:"super.computer" mux.form:"true"`
	Rival         string `mux.url:"rival" mux.param:"true"`
	Display       string `mux.url:"type" mux.path:"true"`
}

type SaltyMeats struct {
	Bacon []string `mux.url:"bacon" mux.param:"true"`
}

func mockRequest(path string) *http.Request {
	req, _ := http.NewRequest("GET", path, nil)

	req.Form = url.Values{}
	req.PostForm = url.Values{}

	return req
}

// This is just to test the wrapper.  The rest of the testing is on what it wraps.
func TestParseMuxRequestToStruct(t *testing.T) {
	req := mockRequest("/root/period")

	handle := "lord.nikon"
	req.Form.Add("handle", handle)

	hr := &HackerReferences{}

	err := ParseMuxRequestToStruct(req, hr)
	assert.NoError(t, err)

	assert.Equal(t, handle, hr.Handle)
}

func TestMuxURLFormGet(t *testing.T) {
	req := mockRequest("/root/period")

	handle := "lord.nikon"
	req.Form.Add("handle", handle)

	hr := &HackerReferences{}
	mug := &MuxURLGetter{req}

	err := tagconfig.Process(mug, hr)
	assert.NoError(t, err)

	assert.Equal(t, handle, hr.Handle)
}

func TestMuxURLPostFormGet(t *testing.T) {
	req := mockRequest("/root/period")

	handle := "lord.nikon"
	req.PostForm.Add("handle", handle)

	hr := &HackerReferences{}
	mug := &MuxURLGetter{req}

	err := tagconfig.Process(mug, hr)
	assert.NoError(t, err)

	assert.Equal(t, handle, hr.Handle)
}

func TestMuxURLParamGet(t *testing.T) {
	rival := "acid.burn"
	req := mockRequest("/root/period/?rival=" + rival)

	hr := &HackerReferences{}
	mug := &MuxURLGetter{req}

	err := tagconfig.Process(mug, hr)
	assert.NoError(t, err)

	assert.Equal(t, rival, hr.Rival)
}

func TestMuxMultipleCSVURLParamToSliceGet(t *testing.T) {
	bacon := "slice"
	moreBacon := "bits"

	req := mockRequest("/salty/meats/?bacon=" + bacon + "," + moreBacon)

	hr := &SaltyMeats{}
	mug := &MuxURLGetter{req}

	err := tagconfig.Process(mug, hr)
	assert.NoError(t, err)

	assert.Equal(t, []string{bacon, moreBacon}, hr.Bacon)
}

func TestMuxMultipleURLParamToSliceGet(t *testing.T) {
	bacon := "slice"
	moreBacon := "bits"

	req := mockRequest("/salty/meats/?bacon=" + bacon + "&bacon=" + moreBacon)

	hr := &SaltyMeats{}
	mug := &MuxURLGetter{Request: req}

	err := tagconfig.Process(mug, hr)
	assert.NoError(t, err)

	assert.Equal(t, []string{bacon, moreBacon}, hr.Bacon)
}

func TestMuxURLPathVariableGet(t *testing.T) {
	r := mux.NewRouter()

	r.HandleFunc("/display/{type}", handler(t))

	s := httptest.NewServer(r)
	defer s.Close()

	client := &http.Client{}
	displayType := "activeMartix"
	res, err := client.Do(mockRequest(s.URL + "/display/" + displayType))

	assert.NoError(t, err)
	defer res.Body.Close()

	assert.NoError(t, err)

	var hrRes HackerReferences
	err = json.NewDecoder(res.Body).Decode(&hrRes)
	assert.NoError(t, err)

	assert.Equal(t, displayType, hrRes.Display)
}

func TestMuxURLGetterAll(t *testing.T) {
	// Expectations
	displayType := "activeMartix"
	rival := "crash.override"
	handle := "the.plague"
	scomp := "gibson"

	fv := url.Values{"handle": {handle}}
	pfv := url.Values{"super.computer": {scomp}}

	r := mux.NewRouter()
	r.HandleFunc("/display/{type}/", handlerWithForm(t, fv, pfv))
	s := httptest.NewServer(r)
	defer s.Close()

	u := s.URL + "/display/" + displayType + "/?rival=" + rival
	req := mockRequest(u)

	client := &http.Client{}
	res, err := client.Do(req)

	assert.NoError(t, err)

	var hrRes HackerReferences
	err = json.NewDecoder(res.Body).Decode(&hrRes)

	assert.NoError(t, err)

	assert.Equal(t, rival, hrRes.Rival)
	assert.Equal(t, displayType, hrRes.Display)
}

func handlerWithForm(t *testing.T, formValues, postFormValues url.Values) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		req.Form = formValues
		req.PostForm = postFormValues

		hr := &HackerReferences{}
		mug := &MuxURLGetter{req}

		err := tagconfig.Process(mug, hr)
		assert.NoError(t, err)

		err = json.NewEncoder(w).Encode(hr)
		assert.NoError(t, err)
	}
}

func handler(t *testing.T) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		hr := &HackerReferences{}
		mug := &MuxURLGetter{req}

		err := tagconfig.Process(mug, hr)
		assert.NoError(t, err)

		err = json.NewEncoder(w).Encode(hr)
		assert.NoError(t, err)
	}
}

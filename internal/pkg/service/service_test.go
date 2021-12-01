package service

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	tData      *Data
	tSaver     *testSaver
	tExtractor *testExtractor
	tEcho      *echo.Echo
	tReq       *http.Request
	tRec       *httptest.ResponseRecorder
)

func initTest(t *testing.T) {
	tSaver = &testSaver{name: "test.wav"}
	tExtractor = &testExtractor{res: "olia"}
	tData = newTestData(tSaver, tExtractor)
	tEcho = initRoutes(tData)
	tReq = newTestRequest("in.epub")
	tRec = httptest.NewRecorder()
}

func TestLive(t *testing.T) {
	initTest(t)
	req := httptest.NewRequest(http.MethodGet, "/live", nil)

	e := initRoutes(tData)
	e.ServeHTTP(tRec, req)
	assert.Equal(t, http.StatusOK, tRec.Code)
	assert.Equal(t, `{"service":"OK"}`, tRec.Body.String())
}

func TestExtract(t *testing.T) {
	initTest(t)

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusOK, tRec.Code)
	assert.Equal(t, `{"text":"test"}`+"\n", tRec.Body.String())
}

func TestExtract_FailData(t *testing.T) {
	initTest(t)
	req := newTestRequest("")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	tEcho.ServeHTTP(tRec, req)

	assert.Equal(t, http.StatusBadRequest, tRec.Code)
}

func TestExtract_FailType(t *testing.T) {
	initTest(t)
	tReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusBadRequest, tRec.Code)
}

func TestExtract_FailExt(t *testing.T) {
	initTest(t)
	req := newTestRequest("file.wav")

	tEcho.ServeHTTP(tRec, req)

	assert.Equal(t, http.StatusBadRequest, tRec.Code)
}

func TestExtract_FailSaver(t *testing.T) {
	initTest(t)

	tSaver.err = errors.New("olia")

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusInternalServerError, tRec.Code)
}

func TestExtract_FailConvert(t *testing.T) {
	initTest(t)

	tExtractor.err = errors.New("olia")

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusInternalServerError, tRec.Code)
}

func TestExtract_FailRead(t *testing.T) {
	initTest(t)

	tData.readFunc = func(string) ([]byte, error) { return nil, errors.New("olia") }

	tEcho.ServeHTTP(tRec, tReq)

	assert.Equal(t, http.StatusInternalServerError, tRec.Code)
}

func newTestRequest(file string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if file != "" {
		part, _ := writer.CreateFormFile("file", file)
		_, _ = io.Copy(part, strings.NewReader("body"))
	}
	writer.Close()
	req := httptest.NewRequest("POST", "/text", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

type testSaver struct {
	name string
	err  error
	data bytes.Buffer
}

func (s *testSaver) Save(name string, reader io.Reader) (string, error) {
	io.Copy(&s.data, reader)
	return s.name, s.err
}

type testExtractor struct {
	err     error
	nameIn  string
	nameOut string
	res     string
}

func (s *testExtractor) Extract(nameIn, nameOut string) error {
	s.nameIn = nameIn
	s.nameOut = nameOut
	return s.err
}

func newTestData(s FileSaver, e Extractor) *Data {
	return &Data{Saver: s, Extractor: e, readFunc: func(string) ([]byte, error) { return []byte("test"), nil }}
}

func Test_validateExt(t *testing.T) {
	type args struct {
		ext string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "txt", args: args{ext: ".txt"}, want: false},
		{name: "epub", args: args{ext: ".epub"}, want: true},
		{name: "mobi", args: args{ext: ".mobi"}, want: true},
		{name: "docx", args: args{ext: ".docx"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateExt(tt.args.ext); got != tt.want {
				t.Errorf("validateExt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNewFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "ext", args: args{file: "olia.docx"}, want: "olia_out.txt"},
		{name: "path", args: args{file: "/olia/olia.docx"}, want: "/olia/olia_out.txt"},
		{name: "several ext", args: args{file: "/olia/olia.docx.epub"}, want: "/olia/olia.docx_out.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNewFile(tt.args.file); got != tt.want {
				t.Errorf("getNewFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

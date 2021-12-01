package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type (
	// FileSaver saves the file with the provided name
	FileSaver interface {
		Save(name string, reader io.Reader) (string, error)
	}

	// Extractor extracts text
	Extractor interface {
		Extract(nameIn string, nameOut string) error
	}

	//Data is service operation data
	Data struct {
		Port int

		Saver     FileSaver
		Extractor Extractor

		readFunc func(string) ([]byte, error)
	}
)

//StartWebServer starts the HTTP service and listens for the convert requests
func StartWebServer(data *Data) error {
	goapp.Log.Infof("starting HTTP audio convert service at %d", data.Port)
	portStr := strconv.Itoa(data.Port)
	data.readFunc = ioutil.ReadFile
	e := initRoutes(data)

	e.Server.Addr = ":" + portStr
	e.Server.ReadHeaderTimeout = 5 * time.Second
	e.Server.ReadTimeout = 45 * time.Second
	e.Server.WriteTimeout = 45 * time.Second

	w := goapp.Log.Writer()
	defer w.Close()
	l := log.New(w, "", 0)
	gracehttp.SetLogger(l)

	return gracehttp.Serve(e.Server)
}

var promMdlw *prometheus.Prometheus

func init() {
	promMdlw = prometheus.NewPrometheus("get_text", nil)
}

func initRoutes(data *Data) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	promMdlw.Use(e)

	e.POST("/text", convert(data))
	e.GET("/live", live(data))

	goapp.Log.Info("Routes:")
	for _, r := range e.Routes() {
		goapp.Log.Infof("  %s %s", r.Method, r.Path)
	}
	return e
}

type output struct {
	Text string `json:"text"`
}

func convert(data *Data) func(echo.Context) error {
	return func(c echo.Context) error {
		defer goapp.Estimate("extract method")()

		form, err := c.MultipartForm()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "no multipart form data")
		}
		defer cleanFiles(form)

		files, ok := form.File["file"]
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "no file")
		}
		if len(files) > 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "multiple files")
		}

		file := files[0]
		ext := filepath.Ext(file.Filename)
		ext = strings.ToLower(ext)
		if !validateExt(ext) {
			return echo.NewHTTPError(http.StatusBadRequest, "wrong file type: "+ext)
		}

		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "can't read file")
		}
		defer src.Close()

		id := uuid.New().String()
		fileName := id + ext

		est := goapp.Estimate("Saving")
		fileNameIn, err := data.Saver.Save(fileName, src)
		if err != nil {
			goapp.Log.Error(err)
			return errors.Wrap(err, "can not save file")
		}
		defer deleteFile(fileNameIn)
		fileNameOut := getNewFile(fileNameIn)
		defer deleteFile(fileNameOut)

		est()

		est = goapp.Estimate("Extract")
		err = data.Extractor.Extract(fileNameIn, fileNameOut)
		if err != nil {
			goapp.Log.Error(err)
			return errors.Wrap(err, "can not extract txt")
		}
		est()

		est = goapp.Estimate("Read")
		fd, err := data.readFunc(fileNameOut)
		if err != nil {
			goapp.Log.Error(err)
			return errors.Wrap(err, "can not read file")
		}
		est()

		res := &output{}
		res.Text = string(fd)

		return c.JSON(http.StatusOK, res)
	}
}

func cleanFiles(f *multipart.Form) {
	if f != nil {
		f.RemoveAll()
	}
}

func deleteFile(file string) {
	os.RemoveAll(file)
}

func validateExt(ext string) bool {
	return ext == ".epub" || ext == ".mobi" || ext == ".docx"
}

func live(data *Data) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, []byte(`{"service":"OK"}`))
	}
}

func getNewFile(file string) string {
	f := filepath.Base(file)
	f = strings.TrimSuffix(f, filepath.Ext(f))
	d := filepath.Dir(file)
	return filepath.Join(d, fmt.Sprintf("%s_out.txt", f))
}

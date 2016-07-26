package img_test

import (
	"errors"
	"github.com/dooman87/kolibri/test"
	"github.com/dooman87/transformimgs/img"
	"net/http"
	"net/http/httptest"
	"testing"
)

type resizerMock struct{}

func (r *resizerMock) Resize(data []byte, size string) ([]byte, error) {
	if string(data) == "321" && size == "300x200" {
		return []byte("123"), nil
	}
	return nil, errors.New("resize_error")
}

func (r *resizerMock) FitToSize(data []byte, size string) ([]byte, error) {
	if string(data) == "321" && size == "300x200" {
		return []byte("123"), nil
	}
	return nil, errors.New("resize_error")
}

func (r *resizerMock) Optimise(data []byte) ([]byte, error) {
	if string(data) == "321" {
		return []byte("123"), nil
	}
	return nil, errors.New("resize_error")
}

type readerMock struct{}

func (r *readerMock) Read(url string) ([]byte, error) {
	if url == "http://site.com/img.png" {
		return []byte("321"), nil
	}
	return nil, errors.New("read_error")
}

func TestResize(t *testing.T) {
	service := &img.Service{
		Reader:    &readerMock{},
		Processor: &resizerMock{},
	}
	test.Service = service.ResizeUrl

	testCases := []test.TestCase{
		{"http://localhost/img?url=http://site.com/img.png&size=300x200", http.StatusOK, "Success",
			func(w *httptest.ResponseRecorder, t *testing.T) {
				if w.Header().Get("Cache-Control") != "public, max-age=86400" {
					t.Errorf("Expected to get Cache-Control header")
				}
				if w.Header().Get("Content-Length") != "3" {
					t.Errorf("Expected to get Content-Length header equal to 3 but got [%s]", w.Header().Get("Content-Length"))
				}
			}},
		{"http://localhost/img?size=300x200", http.StatusBadRequest, "Param url is required", nil},
		{"http://localhost/img?url=http://site.com/img.png", http.StatusBadRequest, "Param size is required", nil},
		{"http://localhost/img?url=NO_SUCH_IMAGE&size=300x200", http.StatusInternalServerError, "Read error", nil},
		{"http://localhost/img?url=http://site.com/img.png&size=BADSIZE", http.StatusInternalServerError, "Resize error", nil},
	}

	test.RunRequests(testCases, t)
}

func TestFitToSize(t *testing.T) {
	service := &img.Service{
		Reader:    &readerMock{},
		Processor: &resizerMock{},
	}
	test.Service = service.FitToSizeUrl

	testCases := []test.TestCase{
		{"http://localhost/fit?url=http://site.com/img.png&size=300x200", http.StatusOK, "Success",
			func(w *httptest.ResponseRecorder, t *testing.T) {
				if w.Header().Get("Cache-Control") != "public, max-age=86400" {
					t.Errorf("Expected to get Cache-Control header")
				}
				if w.Header().Get("Content-Length") != "3" {
					t.Errorf("Expected to get Content-Length header equal to 3 but got [%s]", w.Header().Get("Content-Length"))
				}
			}},
		{"http://localhost/fit?size=300x200", http.StatusBadRequest, "Param url is required", nil},
		{"http://localhost/fit?url=http://site.com/img.png", http.StatusBadRequest, "Param size is required", nil},
		{"http://localhost/fit?url=NO_SUCH_IMAGE&size=300x200", http.StatusInternalServerError, "Read error", nil},
		{"http://localhost/fit?url=http://site.com/img.png&size=BADSIZE", http.StatusBadRequest, "size param should be in format WxH", nil},
		{"http://localhost/fit?url=http://site.com/img.png&size=300", http.StatusBadRequest, "size param should be in format WxH", nil},
	}

	test.RunRequests(testCases, t)
}

func TestOptimise(t *testing.T) {
	service := &img.Service{
		Reader:    &readerMock{},
		Processor: &resizerMock{},
	}
	test.Service = service.OptimiseUrl

	testCases := []test.TestCase{
		{"http://localhost/img?url=http://site.com/img.png", http.StatusOK, "Success",
			func(w *httptest.ResponseRecorder, t *testing.T) {
				if w.Header().Get("Cache-Control") != "public, max-age=86400" {
					t.Errorf("Expected to get Cache-Control header")
				}
				if w.Header().Get("Content-Length") != "3" {
					t.Errorf("Expected to get Content-Length header equal to 3 but got [%s]", w.Header().Get("Content-Length"))
				}
			}},
		{"http://localhost/img", http.StatusBadRequest, "Param url is required", nil},
		{"http://localhost/fit?url=NO_SUCH_IMAGE", http.StatusInternalServerError, "Read error", nil},
	}

	test.RunRequests(testCases, t)
}

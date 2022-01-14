package aliyundrive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/chyroc/gorequests"
	"github.com/sirupsen/logrus"

	"github.com/chyroc/go-aliyundrive/internal/helper_tool"
)

type RawRequestReq struct {
	Scope  string
	API    string
	Method string
	URL    string
	Body   interface{}
	IsFile bool

	headers map[string]string
}

type Response struct {
	Method        string      // request method
	URL           string      // request url
	StatusCode    int         // http response status code
	Header        http.Header // http response header
	ContentLength int64       // http response content length
	Body          []byte
}

func (r *AliyunDrive) RawRequest(ctx context.Context, req *RawRequestReq, resp interface{}) (*Response, error) {
	return r.rawRequest(ctx, req, resp)
}

func (r *AliyunDrive) rawRequest(ctx context.Context, req *RawRequestReq, resp interface{}) (response *Response, err error) {
	r.log(ctx, LogLevelInfo, "%s#%s call api", req.Scope, req.API)

	req.headers, err = r.prepareHeaders(ctx, req)
	if err != nil {
		return nil, err
	}

	response, err = r.doRequest(ctx, req, resp)
	if err != nil {
		r.log(ctx, LogLevelError, "%s#%s %s %s failed, status_code: %d, error: %s", req.Scope, req.API, req.Method, req.URL, response.StatusCode, err)
		return response, err
	}
	msg := getCodeMsg(resp)
	if msg != "" {
		r.log(ctx, LogLevelError, "%s#%s %s %s failed, status_code: %d, msg: %s", req.Scope, req.API, req.Method, req.URL, response.StatusCode, msg)
		return response, fmt.Errorf(msg)
	}

	r.log(ctx, LogLevelDebug, "%s#%s success, status_code: %d, response: %s", req.Scope, req.API, response.StatusCode, "TODO")

	return response, nil
}

func (r *AliyunDrive) doRequest(ctx context.Context, requestParam *RawRequestReq, realResponse interface{}) (*Response, error) {
	response := new(Response)
	realReq, err := parseRequestParam(requestParam)
	if err != nil {
		return response, err
	}

	response.Method = realReq.Method
	response.URL = realReq.URL
	response.Header = map[string][]string{}

	if r.logLevel <= LogLevelTrace {
		r.log(ctx, LogLevelTrace, "request %s#%s, %s %s, header=%s, body=%s", requestParam.Scope, requestParam.API, realReq.Method, realReq.URL, jsonString(realReq.Headers), string(realReq.RawBody))
	}

	req := r.session.New(realReq.Method, realReq.URL).
		WithHeaders(realReq.Headers).
		WithBody(realReq.Body).WithLogger(gorequests.NewDiscardLogger())
	resp, err := req.Response()
	if err != nil {
		return response, err
	}

	_, media, _ := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	respFilename := ""
	if media != nil {
		respFilename = media["filename"]
	}

	response.StatusCode = resp.StatusCode
	response.Header = resp.Header
	response.ContentLength = resp.ContentLength

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	response.Body = bs

	if r.logLevel <= LogLevelTrace {
		r.log(ctx, LogLevelTrace, "response %s#%s, %s %s, body=%s", requestParam.Scope, requestParam.API, realReq.Method, realReq.URL, string(bs))
	}

	if realResponse != nil {
		if resp != nil && resp.StatusCode == http.StatusOK {
			isSpecResp := false
			if setter, ok := realResponse.(readerSetter); ok {
				isSpecResp = true
				setter.SetReader(bytes.NewReader(bs))
			}
			if setter, ok := realResponse.(filenameSetter); ok {
				isSpecResp = true
				setter.SetFilename(respFilename)
			}
			if isSpecResp {
				return response, nil
			}
		}

		if resp.StatusCode != http.StatusNoContent {
			if err = json.Unmarshal(bs, realResponse); err != nil {
				return response, fmt.Errorf("invalid json: %s", bs)
			}
		}
	}

	return response, nil
}

func (r *AliyunDrive) prepareHeaders(ctx context.Context, req *RawRequestReq) (map[string]string, error) {
	headers := map[string]string{
		"User-Agent":      userAgent,
		"Referer":         "https://www.aliyundrive.com/",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Host":            "api.aliyundrive.com",
	}
	if req.Method == http.MethodPost {
		headers["Content-Type"] = "application/json; charset=utf-8"
	}
	if token, _ := r.store.Get(ctx, ""); token != nil && token.AccessToken != "" {
		headers["Authorization"] = "Bearer " + token.AccessToken
	}

	return headers, nil
}

func parseRequestParam(req *RawRequestReq) (*realRequestParam, error) {
	uri := req.URL
	var body io.Reader
	var rawBody []byte
	headers := req.headers
	if headers == nil {
		headers = map[string]string{}
	}

	if req.Body != nil {
		if reader, ok := req.Body.(io.Reader); ok {
			body = reader
			rawBody = []byte("<io.Reader>")
		} else {
			vv := reflect.ValueOf(req.Body)
			vt := reflect.TypeOf(req.Body)

			if vt.Kind() == reflect.Ptr {
				vv = vv.Elem()
				vt = vt.Elem()
			}

			if vt.Kind() == reflect.Map {
				bs, err := json.Marshal(req.Body)
				if err != nil {
					return nil, err
				}
				rawBody = bs
				body = bytes.NewReader(bs)
			} else {
				q := url.Values{}
				isNeedQuery := false
				isNeedBody := false
				isNeedFormURLEncoded := false
				filedata := map[string]string{}
				var reader io.Reader
				fileKey := ""
				formURLEncoded := url.Values{}

				for i := 0; i < vt.NumField(); i++ {
					fieldVV := vv.Field(i)
					fieldVT := vt.Field(i)

					if fieldVV.Kind() == reflect.Ptr && fieldVV.IsNil() {
						continue
					}
					if fieldVV.Kind() == reflect.Slice && fieldVV.Len() == 0 {
						continue
					}

					// path, query, json, form-
					if path := fieldVT.Tag.Get("path"); path != "" {
						if strings.Contains(uri, ":"+path) {
							uri = strings.ReplaceAll(uri, ":"+path, helper_tool.ReflectToString(fieldVV))
						} else {
							uri = strings.ReplaceAll(uri, "{"+path+"}", helper_tool.ReflectToString(fieldVV))
						}
						continue
					} else if query := fieldVT.Tag.Get("query"); query != "" {
						isNeedQuery = true
						for _, v := range helper_tool.ReflectToQueryString(fieldVV) {
							q.Add(query, v)
						}
						continue
					} else if j := fieldVT.Tag.Get("json"); j != "" {
						if strings.HasSuffix(j, ",omitempty") {
							j = j[:len(j)-10]
						}
						if req.IsFile {
							fileKey = j
							if r, ok := fieldVV.Interface().(io.Reader); ok {
								reader = r
							} else {
								filedata[j] = helper_tool.ReflectToString(fieldVV)
							}
						} else {
							isNeedBody = true
						}
						continue
					} else if j := fieldVT.Tag.Get("form-url-encoded"); j != "" {
						if strings.HasSuffix(j, ",omitempty") {
							j = j[:len(j)-10]
						}
						isNeedFormURLEncoded = true
						for _, v := range helper_tool.ReflectToQueryString(fieldVV) {
							formURLEncoded.Add(j, v)
						}
						continue
					}
				}

				if isNeedBody {
					bs, err := json.Marshal(req.Body)
					if err != nil {
						return nil, err
					}
					rawBody = bs
					body = bytes.NewBuffer(bs)
				}

				if isNeedFormURLEncoded {
					s := formURLEncoded.Encode()
					rawBody = []byte(s)
					body = bytes.NewBuffer(rawBody)
					headers["Content-Type"] = "application/x-www-form-urlencoded"
				}

				if req.IsFile {
					contentType, bod, err := newFileUploadRequest(filedata, fileKey, reader)
					if err != nil {
						return nil, err
					}
					headers["Content-Type"] = contentType
					body = bod
					rawBody = []byte("<FILE>")
				}

				if isNeedQuery {
					uri = uri + "?" + q.Encode()
				}
			}
		}
	}

	return &realRequestParam{
		Method:  strings.ToUpper(req.Method),
		URL:     uri,
		Body:    body,
		RawBody: rawBody,
		Headers: headers,
	}, nil
}

type realRequestParam struct {
	Method      string
	URL         string
	Body        io.Reader
	ContentType string
	RawBody     []byte
	Headers     map[string]string
}

func newFileUploadRequest(params map[string]string, filekey string, reader io.Reader) (string, io.Reader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(filekey, "file.file")
	if err != nil {
		return "", nil, err
	}
	if reader != nil {
		if _, err = io.Copy(part, reader); err != nil {
			return "", nil, err
		}
	}
	for key, val := range params {
		if err = writer.WriteField(key, val); err != nil {
			return "", nil, err
		}
	}
	if err = writer.Close(); err != nil {
		return "", nil, err
	}

	return writer.FormDataContentType(), body, nil
}

type readerSetter interface {
	SetReader(file io.Reader)
}

type filenameSetter interface {
	SetFilename(filename string)
}

var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"

type logrusLogger struct{}

func (r *logrusLogger) Info(ctx context.Context, format string, v ...interface{}) {
	logrus.Infof(format, v...)
}

func (r *logrusLogger) Error(ctx context.Context, format string, v ...interface{}) {
	logrus.Errorf(format, v...)
}

func jsonString(v interface{}) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}

func getCodeMsg(v interface{}) (msg string) {
	if v == nil {
		return ""
	}
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	if vv.Kind() != reflect.Struct {
		return ""
	}

	codeMsg := vv.FieldByName("Message")
	if codeMsg.IsValid() {
		if codeMsg.Kind() == reflect.String {
			msg = codeMsg.String()
		}
	}
	return msg
}

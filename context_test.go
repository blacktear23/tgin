package tgin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ReadBodyString(resp *http.Response) string {
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func getRequest(path string, body string) *http.Request {
	buf := bytes.NewBufferString(body)
	return httptest.NewRequest("GET", path, buf)
}

func postRequest(path string, body string) *http.Request {
	buf := bytes.NewBufferString(body)
	return httptest.NewRequest("POST", path, buf)
}

func createTestContext(r *http.Request) *Context {
	w := httptest.NewRecorder()
	return newContext(w, r)
}

func TestBindJSON(t *testing.T) {
	req := postRequest("/", "{\"a\": 1}")
	ctx := createTestContext(req)
	data := H{}
	ctx.BindJSON(&data)
	assertEqual(t, fmt.Sprintf("%v", H{"a": 1}), fmt.Sprintf("%v", data), "Bind json result not correct")
}

func TestGetQuery(t *testing.T) {
	req := getRequest("/test?a=1&b=2", "")
	ctx := createTestContext(req)
	val, have := ctx.GetQuery("a")
	assertTrue(t, have, "Query not have key a")
	assertEqual(t, "1", val, "Query a value not correct")
}

func TestGetQueryArray(t *testing.T) {
	req := getRequest("/test?a=1&a=2", "")
	ctx := createTestContext(req)
	val, have := ctx.GetQueryArray("a")
	assertTrue(t, have, "Query not have key: a")
	assertEqual(t, []string{"1", "2"}, val, "Query a value not correct")
}

func TestGetHeader(t *testing.T) {
	req := getRequest("/", "")
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	ctx := createTestContext(req)
	assertEqual(t, "1.2.3.4", ctx.GetHeader("X-Forwarded-For"), "Header X-Forwarded-For not correct")
	assertEqual(t, "", ctx.GetHeader("Not-Exists"))
}

func TestPostForm(t *testing.T) {
	req := postRequest("/", "a=1&a=2&b=3")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx := createTestContext(req)
	vals := ctx.PostFormArray("a")
	assertEqual(t, []string{"1", "2"}, vals, "Form value a not correct")
	val := ctx.PostForm("b")
	assertEqual(t, "3", val, "Form value b not correct")
}

func TestPostFormWithWrongContentType(t *testing.T) {
	req := postRequest("/", "a=1&a=2&b=3")
	req.Header.Add("Content-Type", "application/xml")
	ctx := createTestContext(req)
	vals := ctx.PostFormArray("a")
	assertEqual(t, []string{}, vals, "Form value a not correct")
	val := ctx.PostForm("b")
	assertEqual(t, "", val, "Form value b not correct")
}

func TestFormFile(t *testing.T) {
	body := bytes.NewBuffer(nil)
	mpw := multipart.NewWriter(body)
	fw, _ := mpw.CreateFormFile("file", "upload.txt")
	fw.Write([]byte("This is upload.txt content"))
	mpw.Close()
	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Add("Content-Type", mpw.FormDataContentType())
	ctx := createTestContext(req)
	fh, err := ctx.FormFile("file")
	assertNil(t, err)
	assertEqual(t, "upload.txt", fh.Filename, "Filename not correct")
	fp, err := fh.Open()
	assertNil(t, err)
	fcontent, _ := ioutil.ReadAll(fp)
	assertEqual(t, []byte("This is upload.txt content"), fcontent)
}

func TestOutputJson(t *testing.T) {
	req := getRequest("/", "")
	ctx := createTestContext(req)
	ctx.JSON(200, H{"key": "value"})
	resp := ctx.Writer.(*httptest.ResponseRecorder).Result()
	assertEqual(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"), "Content-Type not correct")
	assertEqual(t, 200, resp.StatusCode, "Status code is not correct")
	body := ReadBodyString(resp)
	assertEqual(t, "{\"key\":\"value\"}\n", body, "Body not correct")
}

func TestOutputIndentedJson(t *testing.T) {
	req := getRequest("/", "")
	ctx := createTestContext(req)
	ctx.IndentedJSON(200, H{"key": "value"})
	resp := ctx.Writer.(*httptest.ResponseRecorder).Result()
	assertEqual(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"), "Content-Type not correct")
	assertEqual(t, 200, resp.StatusCode, "Status code is not correct")
	body := ReadBodyString(resp)
	assertEqual(t, "{\n    \"key\": \"value\"\n}\n", body, "Body not correct")
}

func TestOutputText(t *testing.T) {
	req := getRequest("/", "")
	ctx := createTestContext(req)
	ctx.Text(200, "body")
	resp := ctx.Writer.(*httptest.ResponseRecorder).Result()
	assertEqual(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), "Content-Type not correct")
	assertEqual(t, 200, resp.StatusCode, "Status Code not correct")
	body := ReadBodyString(resp)
	assertEqual(t, "body", body, "Body not correct")
}

func TestOutputString(t *testing.T) {
	req := getRequest("/", "")
	ctx := createTestContext(req)
	ctx.String(200, "hello %s", "world")
	resp := ctx.Writer.(*httptest.ResponseRecorder).Result()
	assertEqual(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), "Content-Type not correct")
	assertEqual(t, 200, resp.StatusCode, "Status Code not correct")
	body := ReadBodyString(resp)
	assertEqual(t, "hello world", body, "Body not correct")
}

func TestOutputRedirect(t *testing.T) {
	req := getRequest("/", "")
	ctx := createTestContext(req)
	ctx.Redirect(302, "/login")
	resp := ctx.Writer.(*httptest.ResponseRecorder).Result()
	assertEqual(t, 302, resp.StatusCode, "Status code not correct")
	assertEqual(t, "/login", resp.Header.Get("Location"), "Location header not correct")
}

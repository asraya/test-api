package helpers

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

var headers = map[string]string{
	"Content-Type": "application/json",
	"Accept":       "application/json",
}
var client = &http.Client{}

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Exist   []*net.NS
	Data    interface{} `json:"data"`
}

type EmptyObj struct{}

func BuildResponse(status bool, message string, data interface{}) Response {
	res := Response{
		Status:  status,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
	return res
}

func BuildResponseOverview(status bool, message string, data interface{}, exist []*net.NS) Response {
	res := Response{
		Status:  status,
		Message: message,
		Errors:  nil,
		Exist:   exist,
		Data:    data,
	}
	return res
}

func doRequest(method string, url string, body []byte) ([]byte, error) {

	payload := strings.NewReader(string(body))

	req, err := http.NewRequest(method, url, payload)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	fmt.Println("LOG ERR", err)

	if err != nil {
		return nil, err
	}

	fmt.Println("LOG REQ", req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if isStatusError(resp.StatusCode) {
		return nil, fmt.Errorf("Status error: %v %v", resp.StatusCode, string(respBody))
	}

	defer resp.Body.Close()
	return respBody, err
}
func Get(url string, body []byte) ([]byte, error) {
	return doRequest("GET", url, body)
}
func isStatusError(statusCode int) bool {
	return statusCode >= http.StatusBadRequest
}

func BuildErrorResponse(message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	res := Response{
		Status:  false,
		Message: message,
		Errors:  splittedError,
		Data:    data,
	}
	return res
}

func BuildErrorResponseCekNs(message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	res := Response{
		Status: true,
		Data:   splittedError,
	}
	return res
}

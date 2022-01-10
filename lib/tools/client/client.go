package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/syedomair/ex-paygate-lib/lib/models"
)

// APICall Public
type APICall struct {
	Method,
	PathParam,
	ID,
	Path,
	Body,
	ResponseVar,
	ResponseID string
}

// Public
const (
	Result  = "result"
	Data    = "data"
	Success = "success"
	Failure = "failure"
)

// LoadAPIRequest Public
func LoadAPIRequest(apiRequests []*models.APIRequest, recordIDs map[string]string, urlStr string,
	deptFile *os.File, userFile *os.File, customFile *os.File, itemFile *os.File) error {

	for _, apiRequest := range apiRequests {
		var jsonResult, jsonData []byte
		var err error

		if apiRequest.Path == "workflows" {
			time.Sleep(2 * time.Second)
			//time.Sleep(500 * time.Millisecond)
		}
		if apiRequest.Path == "upload-files" {
			jsonResult, jsonData, err = CallAPIUploadUserDeptFiles(urlStr+apiRequest.Path, recordIDs["token"], recordIDs["api_key"], deptFile, userFile)
		} else if apiRequest.Path == "upload-custom-files" {
			jsonResult, jsonData, err = CallAPIUploadCustomFile(urlStr+apiRequest.Path, recordIDs["token"], recordIDs["api_key"], customFile)
		} else if apiRequest.Path == "upload-item-files" {
			jsonResult, jsonData, err = CallAPIUploadItemFile(urlStr+apiRequest.Path, recordIDs["token"], recordIDs["api_key"], itemFile)
		} else if apiRequest.TitleVar == "" {
			jsonResult, jsonData, _, err = CallAPI("POST", urlStr+apiRequest.Path, ReplaceStringVar(apiRequest.Body, recordIDs), recordIDs["api_key"], recordIDs["token"])
		} else if apiRequest.TitleVar == "PATCH" {
			jsonResult, jsonData, _, err = CallAPI("PATCH", ReplaceStringVar(urlStr+apiRequest.FindPath, recordIDs), ReplaceStringVar(apiRequest.Body, recordIDs), recordIDs["api_key"], recordIDs["token"])
		} else if apiRequest.TitleVar == "GET" {
			jsonResult, jsonData, _, err = CallAPI("GET", ReplaceStringVar(urlStr+apiRequest.Path, recordIDs), apiRequest.Body, recordIDs["api_key"], recordIDs["token"])
		} else if apiRequest.TitleVar == "DELETE" {
			jsonResult, jsonData, _, err = CallAPI("DELETE", ReplaceStringVar(urlStr+apiRequest.Path, recordIDs), apiRequest.Body, recordIDs["api_key"], recordIDs["token"])
		} else {
			jsonResult, jsonData, _, err = CallAPI("GET", urlStr+apiRequest.FindPath+"?title="+url.QueryEscape(apiRequest.TitleVar), apiRequest.Body, recordIDs["api_key"], recordIDs["token"])
		}
		if err != nil {
			fmt.Printf("Error in API call err:%v", err)
		}
		if PrintResult(jsonResult, jsonData, ReplaceStringVar(apiRequest.Body, recordIDs)) == Success {
			recordIDs[apiRequest.IDVar] = ExtractVariable(jsonData, apiRequest.RespVar)
		} else {
			return errors.New("error in API load")
		}
		//time.Sleep(500 * time.Millisecond)
		fmt.Println("---------------------------------------------------------------------------")
	}
	return nil
}

// CallAPIUploadUserDeptFiles Public
func CallAPIUploadUserDeptFiles(url, token, apiKey string, deptFile *os.File, userFile *os.File) (result, data []byte, err error) {

	// New multipart writer.
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)
	fw, _ := writer.CreateFormField("user_header")
	_, _ = io.Copy(fw, strings.NewReader("1"))
	fw, _ = writer.CreateFormField("dept_header")
	_, _ = io.Copy(fw, strings.NewReader("1"))
	fw, _ = writer.CreateFormFile("user_file", "UserData.csv")
	_, _ = io.Copy(fw, userFile)
	fw, _ = writer.CreateFormFile("dept_file", "DeptData.csv")
	_, _ = io.Copy(fw, deptFile)

	// Close multipart writer.
	writer.Close()
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBuf.Bytes()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Token", token)
	req.Header.Set("ApiKey", apiKey)
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed with response code: %d err:%v", resp.StatusCode, err.Error())
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var bodyInterface map[string]interface{}
	json.Unmarshal(body, &bodyInterface)
	jsonResult, _ := json.Marshal(bodyInterface[Result])
	jsonData, _ := json.Marshal(bodyInterface[Data])
	return jsonResult, jsonData, nil
}

// CallAPIUploadCustomFile Public
func CallAPIUploadCustomFile(url, token, apiKey string, customFile *os.File) (result, data []byte, err error) {

	// New multipart writer.
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)
	fw, _ := writer.CreateFormField("custom_header")
	_, _ = io.Copy(fw, strings.NewReader("1"))
	fw, _ = writer.CreateFormFile("custom_file", "CustomData.csv")
	_, _ = io.Copy(fw, customFile)

	// Close multipart writer.
	writer.Close()
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBuf.Bytes()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Token", token)
	req.Header.Set("ApiKey", apiKey)
	client := &http.Client{Timeout: time.Second * 10}
	resp, _ := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed with response code: %d", resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var bodyInterface map[string]interface{}
	json.Unmarshal(body, &bodyInterface)
	jsonResult, _ := json.Marshal(bodyInterface[Result])
	jsonData, _ := json.Marshal(bodyInterface[Data])
	return jsonResult, jsonData, nil
}

// CallAPIUploadItemFile Public
func CallAPIUploadItemFile(url, token, apiKey string, itemFile *os.File) (result, data []byte, err error) {

	// New multipart writer.
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)
	fw, _ := writer.CreateFormField("item_header")
	_, _ = io.Copy(fw, strings.NewReader("1"))
	fw, _ = writer.CreateFormFile("item_file", "ItemData.csv")
	_, _ = io.Copy(fw, itemFile)

	// Close multipart writer.
	writer.Close()
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBuf.Bytes()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Token", token)
	req.Header.Set("ApiKey", apiKey)
	client := &http.Client{Timeout: time.Second * 10}
	resp, _ := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed with response code: %d", resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var bodyInterface map[string]interface{}
	json.Unmarshal(body, &bodyInterface)
	jsonResult, _ := json.Marshal(bodyInterface[Result])
	jsonData, _ := json.Marshal(bodyInterface[Data])
	return jsonResult, jsonData, nil
}

// CallAPI Public
func CallAPI(method, url, requestBody, apiKey, token string) (result, data []byte, httpResponseCode int, err error) {
	fmt.Println(method + " " + url)
	req, err := http.NewRequest(method, url, strings.NewReader(requestBody))
	if err != nil {
		return nil, nil, 0, err
	}
	req.Header.Set("ApiKey", apiKey)
	req.Header.Set("Token", token)
	req.Header.Set("FrontendURL", "xyz.com")
	//iStr := strconv.Itoa(i)
	//req.Header.Set("Testnum", iStr)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("request failed with response code: %d", resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var bodyInterface map[string]interface{}
	json.Unmarshal(body, &bodyInterface)
	jsonResult, _ := json.Marshal(bodyInterface[Result])
	jsonData, _ := json.Marshal(bodyInterface[Data])

	return jsonResult, jsonData, resp.StatusCode, nil
}

// PrintResult Public
func PrintResult(jsonResult, jsonData []byte, requestBody string) string {
	rtnString := Failure
	if strings.Compare(strings.Trim(string(jsonResult), "\""), Success) == 0 {
		fmt.Println("\033[32m" + " PASS" + "\033[39m")
		rtnString = Success
	} else {
		fmt.Println("\033[31m" + " FAIL" + "\033[39m")
		//fmt.Println(strconv.Itoa(i) + " POST " + string(url) + " " + requestBody)
		fmt.Println(string(requestBody))
		fmt.Println(string(jsonData))
		fmt.Println(string(jsonResult))
		rtnString = Failure
	}
	return rtnString
}

// PrintResultError Public
func PrintResultError(jsonResult, jsonData []byte, requestBody string, httpCode int, errCode string, successBln bool) string {
	rtnString := Failure
	if successBln {
		fmt.Println("\033[32m" + " PASS" + "\033[39m")
		rtnString = Success
	} else {
		fmt.Println("\033[31m" + " FAIL" + "\033[39m")
		//fmt.Println(strconv.Itoa(i) + " POST " + string(url) + " " + requestBody)
		fmt.Println(string(requestBody))
		fmt.Println(string(jsonData))
		fmt.Println(string(jsonResult))
		fmt.Println(httpCode)
		fmt.Println(errCode)
		rtnString = Failure
	}
	return rtnString
}

// ExtractVariable Public
func ExtractVariable(jsonData []byte, responseVar string) string {
	var varInterface map[string]interface{}
	json.Unmarshal(jsonData, &varInterface)
	jsonResponseVar, _ := json.Marshal(varInterface[responseVar])
	return string(bytes.Trim(jsonResponseVar, `"`))
}

// GetUserHomeDir Public
func GetUserHomeDir() string {
	user, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
	}
	return user.HomeDir + "/"
}

// ReplaceStringVar Public
func ReplaceStringVar(body string, recordIDs map[string]string) string {
	i := strings.Index(body, "#")
	if i > -1 {
		bodyPrefix := body[:i]
		bodySuffix := body[1+i:]
		//body = bodyPrefix + bodySuffix
		//fmt.Println("bodyPrefix:", bodyPrefix)
		//fmt.Println("bodySuffix:", bodySuffix)
		j := strings.Index(bodySuffix, "#")
		bodySuffixPrefix := bodySuffix[:j]
		bodySuffixSuffix := bodySuffix[1+j:]
		//fmt.Println("bodySuffixPrefix:", bodySuffixPrefix)
		//fmt.Println("bodySuffixSuffix:", bodySuffixSuffix)
		idValue := recordIDs[bodySuffixPrefix]
		return ReplaceStringVar(bodyPrefix+idValue+bodySuffixSuffix, recordIDs)

	}
	return body
}

// UniqueNetworkName Public
func UniqueNetworkName() string {
	currentTime := time.Now()
	uniqueString := strconv.FormatInt(currentTime.UnixNano(), 10)
	return "NETWORK" + uniqueString
}

// UniqueUserName Public
func UniqueUserName() string {
	currentTime := time.Now()
	uniqueString := strconv.FormatInt(currentTime.UnixNano(), 10)
	return "USER-" + uniqueString
}

// LoadJSON Public
func LoadJSON(fileName string, apiRequests []*models.APIRequest) []*models.APIRequest {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	jsonFile, err := os.Open(path + "/" + fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	apiRequestsFromJSONFile := []*models.APIRequest{}
	err = json.Unmarshal(byteValue, &apiRequestsFromJSONFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	apiRequests = append(apiRequests, apiRequestsFromJSONFile...)
	return apiRequests
}

// LoadJSONTest Public
func LoadJSONTest(fileName string, apiRequests []*models.APIRequestTest) []*models.APIRequestTest {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	jsonFile, err := os.Open(path + "/" + fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	apiRequestsFromJSONFile := []*models.APIRequestTest{}
	err = json.Unmarshal(byteValue, &apiRequestsFromJSONFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	apiRequests = append(apiRequests, apiRequestsFromJSONFile...)
	return apiRequests
}

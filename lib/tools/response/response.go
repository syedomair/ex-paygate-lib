package response

import (
	"encoding/json"
	"net/http"

	"github.com/syedomair/ex-paygate/lib/tools/logger"
)

const (
	Success = "success"
	Failure = "failure"
)

// ErrorResponseHelper func
func ErrorResponseHelper(requestID, methodName string, logger logger.Logger,
	w http.ResponseWriter, errorCode string, errorMessage string, httpStatus int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(httpStatus)
	w.Write(errorResponse(errorCode, errorMessage))
	logger.Error(requestID, "M:%v EC:%v EM:%v ", methodName, errorCode, errorMessage)
}

// SuccessResponseHelper func
func SuccessResponseHelper(w http.ResponseWriter, class interface{}, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(successResponse(class))
	//c.Logger.Log("METHOD", "SuccessResponseHelper", "END", c.successResponse(class))
}

// SuccessResponseList func
func SuccessResponseList(w http.ResponseWriter, class interface{}, offset string, limit string, count string) {
	tempResponse := make(map[string]interface{})
	tempResponse["count"] = count
	tempResponse["list"] = class

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(successResponse(tempResponse))
	//c.Logger.Log("METHOD", "SuccessResponseList", "END", c.successResponse(tempResponse))
}

func errorResponse(errorCode string, message string) []byte {
	class := map[string]string{"error_code": errorCode, "message": message}
	return commonResponse(class, Failure)
}

func successResponse(class interface{}) []byte {
	return commonResponse(class, Success)
}

func commonResponse(class interface{}, result string) []byte {
	response := make(map[string]interface{})
	response["result"] = result
	response["data"] = class
	jsonByte, _ := json.Marshal(response)
	return jsonByte
}

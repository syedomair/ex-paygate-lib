package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/syedomair/ex-pay-gateway/lib/models"
	"github.com/syedomair/ex-pay-gateway/lib/tools/logger"
)

const (
	errorCodePrefix = "5"
)

const (
	SUCCESS = "success"
	FAILURE = "failure"
	STRING  = "STRING"
)

// ContextKey Public
type ContextKey string

// ContextKeyRequestID Public
const ContextKeyRequestID ContextKey = "request_id"

// GetRequestID Public
func GetRequestID(r *http.Request) string {
	requestID, _ := r.Context().Value(ContextKeyRequestID).(string)
	return requestID
}

// ValidateInputParameters Public
func ValidateInputParameters(r *http.Request, requestID string, logger logger.Logger, paramConf map[string]models.ParamConf, pathParamConf map[string]string) (map[string]interface{}, string, error) {
	methodName := "ValidateInputParameters"
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		return nil, errorCodePrefix + "01", errors.New("Error reading the request body: " + err.Error())
	}
	var jsonMap map[string]interface{}
	decodedJSON := json.NewDecoder(bytes.NewReader(body))

	if err := decodedJSON.Decode(&jsonMap); err != nil {
		return nil, errorCodePrefix + "02", errors.New("Invalid JSON in request BODY:" + string(body) + "  ERROR:" + err.Error())
	}
	logger.Debug(requestID, "M:%v request:%v", methodName, r)
	logger.Debug(requestID, "M:%v jsonMap:%v", methodName, jsonMap)

	for k := range jsonMap {
		if _, ok := paramConf[k]; !ok {
			return nil, errorCodePrefix + "03", errors.New(k + " Invalid parameter in JSON")
		}
	}

	for k, v := range paramConf {
		if val, ok := jsonMap[k]; ok {
			if reflect.TypeOf(val).String() != "string" {
				return nil, errorCodePrefix + "04", errors.New(k + " must be a valid string parameter")
			}
			if !v.EmptyAllowed {
				if len(strings.TrimSpace(val.(string))) < 1 {
					return nil, errorCodePrefix + "05", errors.New(k + " cannot be blank")
				}
			}
			if len(strings.TrimSpace(val.(string))) < 1 && v.EmptyAllowed {
				break
			}
			switch fieldType := v.Type; fieldType {
			case STRING:
				if len(strings.TrimSpace(val.(string))) > 100 {
					return nil, errorCodePrefix + "06", errors.New(k + " allowed with max character of 100")
				}

			}

		} else {
			if v.Required {
				return nil, errorCodePrefix + "14", errors.New(k + " is a requird field")
			}
		}
	}
	return jsonMap, "", nil
}

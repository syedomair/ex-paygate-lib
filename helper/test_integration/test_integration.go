package main

import (
	"flag"
	"fmt"

	"github.com/syedomair/ex-paygate-lib/lib/tools/client"
)

func main() {
	serverName := flag.String("server_name", "localhost", "server_name")
	flag.Parse()

	approveURL := `http://` + *serverName + `:8321/v1/approve`
	//voidURL := `http://` + serverName + `:8322/v1/void`
	//captureURL := `http://` + serverName + `:8323/v1/capture`
	//refundURL := `http://` + serverName + `:8323/v1/refund`

	requestBody := `{"merchant_key":"KEY1", "cc_number":"4000000000000000", "cc_expiry":"12345", "currency":"USD", "amount":"88"}`
	jsonResult, jsonData, _, err := client.CallAPI("POST", approveURL, requestBody, "", "")
	if err != nil {
		fmt.Printf("Error in API call err:%v", err)
	}
	client.PrintResult(jsonResult, jsonData, requestBody)
	approveKey := client.ExtractVariable(jsonData, "approve_key")
	fmt.Println(approveKey)

}

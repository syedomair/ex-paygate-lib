package main

import (
	"flag"
	"fmt"

	"github.com/syedomair/ex-paygate-lib/lib/tools/client"
)

const (
	ValidMerchanKey       = "KEY1"
	InValidMerchanKey     = "KEYxyz"
	ValidCCNumber         = "4242424242424242"
	InValidCCNumber       = "0002424242424242"
	AuthFailedCCNumber    = "4000000000000119"
	MerchantKey           = "merchant_key"
	CcNumber              = "cc_number"
	CcCVV                 = "cc_cvv"
	CcMonth               = "cc_month"
	CcYear                = "cc_year"
	Currency              = "currency"
	Amount                = "amount"
	ApproveKey            = "approve_key"
	ApprovedAmountBalance = "approved_amount_balance"
)

func main() {
	serverName := flag.String("server_name", "localhost", "server_name")
	flag.Parse()

	approveURL := `http://` + *serverName + `:8321/v1/authorize`
	voidURL := `http://` + *serverName + `:8322/v1/void`
	captureURL := `http://` + *serverName + `:8323/v1/capture`
	refundURL := `http://` + *serverName + `:8324/v1/refund`

	// Invalid Merchant KEY
	requestBody := `{"` + MerchantKey + `":"` + InValidMerchanKey + `", "` + CcNumber +
		`":"` + ValidCCNumber + `", "` + CcCVV + `":"1234",  "` + CcMonth + `":"12", "` + CcYear + `":"2025","` +
		Currency + `":"USD", "` + Amount + `":"100"}`
	_, _, httpStatus, err := client.CallAPI("POST", approveURL, requestBody, "", "")
	if err != nil {
		fmt.Printf("Error in API call err:%v", err)
	}
	if httpStatus != 400 {
		fmt.Println("test broken: %v" + requestBody)
		return
	}
	// valid Merchant KEY
	requestBody = `{"` + MerchantKey + `":"` + ValidMerchanKey + `", "` + CcNumber +
		`":"` + ValidCCNumber + `", "` + CcCVV + `":"1234",  "` + CcMonth + `":"12", "` + CcYear + `":"2025","` +
		Currency + `":"USD", "` + Amount + `":"100"}`
	_, jsonData, httpStatus, err := client.CallAPI("POST", approveURL, requestBody, "", "")
	if err != nil {
		fmt.Printf("Error in API call err:%v", err)
	}
	if httpStatus != 200 {
		fmt.Println("test broken: %v" + requestBody)
		return
	}

	approveKey := client.ExtractVariable(jsonData, "approve_key")

	// Void approve
	requestBody = `{"` + ApproveKey + `":"` + approveKey + `" }`
	_, jsonData, httpStatus, err = client.CallAPI("POST", voidURL, requestBody, "", "")
	if err != nil {
		fmt.Printf("Error in API call err:%v", err)
	}
	if httpStatus != 200 {
		fmt.Println("test broken: %v" + requestBody)
		return
	}

	// Capture on Void approve
	requestBody = `{"` + ApproveKey + `":"` + approveKey + `", "` + Amount + `":"10" }`
	_, jsonData, httpStatus, err = client.CallAPI("POST", captureURL, requestBody, "", "")
	if err != nil {
		fmt.Printf("Error in API call err:%v", err)
	}

	if httpStatus != 400 {
		fmt.Println("test broken: %v" + requestBody)
		return
	}

	// valid Merchant KEY
	requestBody = `{"` + MerchantKey + `":"` + ValidMerchanKey + `", "` + CcNumber +
		`":"` + ValidCCNumber + `", "` + CcCVV + `":"1234",  "` + CcMonth + `":"12", "` + CcYear + `":"2025","` +
		Currency + `":"USD", "` + Amount + `":"100"}`
	_, jsonData, httpStatus, err = client.CallAPI("POST", approveURL, requestBody, "", "")
	if err != nil {
		fmt.Printf("Error in API call err:%v", err)
	}
	if httpStatus != 200 {
		fmt.Println("test broken: %v" + requestBody)
		return
	}

	approveKey = client.ExtractVariable(jsonData, "approve_key")

	for i := 0; i < 10; i++ {
		go func() {
			requestBody = `{"` + ApproveKey + `":"` + approveKey + `", "` + Amount + `":"10" }`
			_, jsonData, httpStatus, err = client.CallAPI("POST", captureURL, requestBody, "", "")
			if err != nil {
				fmt.Printf("Error in API call err:%v", err)
			}

			if httpStatus != 200 {
				fmt.Println("test broken: %v" + requestBody)
				return
			}
		}()
	}

	balance := client.ExtractVariable(jsonData, "approved_amount_balance")
	fmt.Println(balance)

	// Capture on approve
	requestBody = `{"` + ApproveKey + `":"` + approveKey + `", "` + Amount + `":"10" }`
	_, jsonData, httpStatus, err = client.CallAPI("POST", captureURL, requestBody, "", "")
	if err != nil {
		fmt.Printf("Error in API call err:%v", err)
	}

	if httpStatus != 200 {
		fmt.Println("test broken: %v" + requestBody)
		return
	}
	balance = client.ExtractVariable(jsonData, "approved_amount_balance")
	fmt.Println(balance)

	// Refund on approve
	requestBody = `{"` + ApproveKey + `":"` + approveKey + `", "` + Amount + `":"10" }`
	_, jsonData, httpStatus, err = client.CallAPI("POST", refundURL, requestBody, "", "")
	if err != nil {
		fmt.Printf("Error in API call err:%v", err)
	}

	if httpStatus != 200 {
		fmt.Println("test broken: %v" + requestBody)
		return
	}
	balance = client.ExtractVariable(jsonData, "approved_amount_balance")
	fmt.Println(balance)

}

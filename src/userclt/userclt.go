package main

import (
	"errors"
	"fmt"
	"os"
	"userdet"

	atmi "github.com/endurox-dev/endurox-go"
	"github.com/golang/protobuf/proto"
)

/*
#include <signal.h>
*/
import "C"

const (
	ProgSection = "userclt"
)

var MSomeConfigFlag string = ""
var MSomeOtherConfigFlag int = 0

//Run the listener
func apprun(ac *atmi.ATMICtx) error {

	//Do some work here

	return nil
}

//Init function
//@param ac	ATMI context
//@return error (if erro) or nil
func appinit(ac *atmi.ATMICtx) error {

	if err := ac.TpInit(); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

//Un-init & Terminate the application
//@param ac	ATMI Context
//@param restCode	Return code. atmi.FAIL (-1) or atmi.SUCCEED(0)
func unInit(ac *atmi.ATMICtx, retCode int) {

	ac.TpTerm()
	ac.FreeATMICtx()
	os.Exit(retCode)
}

//Cliet process main entry
func main() {

	ac, errA := atmi.NewATMICtx()

	if nil != errA {
		fmt.Fprintf(os.Stderr, "Failed to allocate cotnext %d:%s!\n",
			errA.Code(), errA.Message())
		os.Exit(atmi.FAIL)
	}

	if err := appinit(ac); nil != err {
		ac.TpLogError("Failed to init: %s", err)
		os.Exit(atmi.FAIL)
	}

	ac.TpLogWarn("Init complete, processing...")

	for i := 0; i < 20000; i++ {

		//Init data buffer
		user := userdet.Userdet{FirstName: "Jim", LastName: "Morrison", Age: 27}

		user.ProfileTags = int32(i/100 + 1)
		user.ProfileData = make([]string, user.ProfileTags)

		for j := 0; j < int(user.ProfileTags); j++ {
			user.ProfileData[j] = fmt.Sprintf("DATA STRING %d", j)
		}

		// Create BLOB
		data, err := proto.Marshal(&user)

		if err != nil {
			ac.TpLogError("Failed to marshal Userdet: %s", err.Error())
			os.Exit(atmi.FAIL)
		}

		// Create transport buffer, CARRAY type (BLOB for Endurox/X)
		transportBuf, errA := ac.NewCarray(data)

		if errA != nil {
			ac.TpLogError("Failed to allocate Carray buffer %s", errA.Message())
			os.Exit(atmi.FAIL)
		}

		// Call "USERADD" service
		_, errA = ac.TpCall("USERADD", transportBuf, 0)

		if errA != nil {
			ac.TpLogError("failed to call USERADD service: %s", errA.Message())
			os.Exit(atmi.FAIL)
		}

		// Parse back the protocol buffer
		var response userdet.Resultdet

		err = proto.Unmarshal(transportBuf.GetBytes(), &response)

		if err != nil {
			ac.TpLogError("Failed to unmarshal : %s", err.Error())
			os.Exit(atmi.FAIL)
		}

		if response.GetStatusCode() != 0 && response.GetStatusMessage() != "OK" {
			ac.TpLogError("INvalid value returned in response %d:%s",
				response.GetStatusCode(), response.GetStatusMessage())
			os.Exit(atmi.FAIL)
		}

	}

	if err := apprun(ac); nil != err {
		unInit(ac, atmi.FAIL)
	}

	unInit(ac, atmi.SUCCEED)
}

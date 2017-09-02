package main

import (
	"errors"
	"fmt"
	"os"
	"ubftab"

	atmi "github.com/endurox-dev/endurox-go"
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

		// Create UBF buffer
		ubf, errA := ac.NewUBF(56000)

		if errA != nil {
			ac.TpLogError("Failed to allocate Carray buffer %s", errA.Message())
			os.Exit(atmi.FAIL)
		}

		if errU := ubf.BAdd(ubftab.A_FIRSTNAME, "Jim"); nil != errU {
			ac.TpLogError("Failed to set A_FIRSTNAME: %s", errU.Message())
			os.Exit(atmi.FAIL)
		}

		if errU := ubf.BAdd(ubftab.A_LASTNAME, "Morrison"); nil != errU {
			ac.TpLogError("Failed to set A_LASTNAME: %s", errU.Message())
			os.Exit(atmi.FAIL)
		}

		if errU := ubf.BAdd(ubftab.A_AGE, 27); nil != errU {
			ac.TpLogError("Failed to set A_AGE: %s", errU.Message())
			os.Exit(atmi.FAIL)
		}

		profileTags := i/100 + 1

		if errU := ubf.BAdd(ubftab.A_PROFILETAGS, profileTags); nil != errU {
			ac.TpLogError("Failed to set A_PROFILETAGS: %s", errU.Message())
			os.Exit(atmi.FAIL)
		}

		for j := 0; j < profileTags; j++ {
			str := fmt.Sprintf("DATA STRING %d", j)
			if errU := ubf.BChg(ubftab.A_PROFILETAGS, j, str); nil != errU {
				ac.TpLogError("Failed to set A_PROFILETAGS %d: %s", j, errU.Message())
				os.Exit(atmi.FAIL)
			}
		}

		// Call "USERADDUBF" service
		_, errA = ac.TpCall("USERADDUBF", ubf, 0)

		if errA != nil {
			ac.TpLogError("failed to call USERADDUBF service: %s", errA.Message())
			os.Exit(atmi.FAIL)
		}

		statuscode, errU := ubf.BGetInt(ubftab.A_STATUSCODE, 0)

		if nil != errU {
			ac.TpLogError("Failed to get A_STATUSCODE: %s", errU.Message())
			os.Exit(atmi.FAIL)
		}

		if statuscode != 0 {
			ac.TpLogError("Invalid status code returned in response"+
				" %d", statuscode)
			os.Exit(atmi.FAIL)
		}

		statusmessage, errU := ubf.BGetString(ubftab.A_STATUSMESSAGE, 0)

		if nil != errU {
			ac.TpLogError("Failed to get A_STATUSMESSAGE: %s", errU.Message())
			os.Exit(atmi.FAIL)
		}

		if statusmessage != "OK" {
			ac.TpLogError("Invalid status message returned in response"+
				", status = %s", statusmessage)
			os.Exit(atmi.FAIL)
		}
	}

	if err := apprun(ac); nil != err {
		unInit(ac, atmi.FAIL)
	}

	unInit(ac, atmi.SUCCEED)
}

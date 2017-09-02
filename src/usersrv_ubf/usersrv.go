package main

import (
	"fmt"
	"os"
	"ubftab"

	atmi "github.com/endurox-dev/endurox-go"
)

const (
	SUCCEED = atmi.SUCCEED
	FAIL    = atmi.FAIL
)

var MCallNumber int = 0

//USERADDUBF service
//@param ac ATMI Context
//@param svc Service call information
func USERADDUBF(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {

	ret := SUCCEED

	//Return to the caller
	defer func() {
		if SUCCEED == ret {
			ac.TpReturn(atmi.TPSUCCESS, 0, &svc.Data, 0)
		} else {
			ac.TpReturn(atmi.TPFAIL, 0, &svc.Data, 0)
		}
		MCallNumber++
	}()

	//Get BUF Handler
	ubf, _ := ac.CastToUBF(&svc.Data)

	//Test the data received
	FirstName, errU := ubf.BGetString(ubftab.A_FIRSTNAME, 0)
	if nil != errU {
		ac.TpLogError("Failed to get A_FIRSTNAME: %s", errU.Message())
		ret = FAIL
		return
	}

	if FirstName != "Jim" {
		ac.TpLogError("Expected FirstName = \"Jim\", but got: %s",
			FirstName)
		ret = FAIL
		return
	}

	LastName, errU := ubf.BGetString(ubftab.A_LASTNAME, 0)
	if nil != errU {
		ac.TpLogError("Failed to get A_LASTNAME: %s", errU.Message())
		ret = FAIL
		return
	}

	if LastName != "Morrison" {
		ac.TpLogError("Expected LastName = \"Morrison\", but got: %s",
			LastName)
		ret = FAIL
		return
	}

	Age, errU := ubf.BGetInt(ubftab.A_AGE, 0)
	if nil != errU {
		ac.TpLogError("Failed to get A_AGE: %s", errU.Message())
		ret = FAIL
		return
	}

	if Age != 27 {
		ac.TpLogError("Expected Age = \"27\", but got: %d",
			Age)
		ret = FAIL
		return
	}

	ProfileTags, errU := ubf.BGetInt(ubftab.A_PROFILETAGS, 0)
	if nil != errU {
		ac.TpLogError("Failed to get A_PROFILETAGS: %s", errU.Message())
		ret = FAIL
		return
	}

	//Test data strings...
	for i := 0; i < ProfileTags; i++ {
		str := fmt.Sprintf("DATA STRING %d", i)

		bufval, errU := ubf.BGetString(ubftab.A_PROFILEDATA, i)
		if nil != errU {
			ac.TpLogError("Failed to get A_LASTNAME: %s", errU.Message())
			ret = FAIL
			return
		}

		if bufval != str {

			ac.TpLogError("ProfileData[%d] expected to be [%s] but got [%s]",
				i, str, bufval)
			ret = FAIL
			return
		}
	}

	//Now send back the status code
	//Reseize a bit buffer, to place new values
	size, errA := ubf.BSizeof()

	if errA != nil {
		ac.TpLogError("Failed to get buffer size: %s", errA.Message())
		ret = FAIL
		return
	}

	errA = ubf.TpRealloc(size + 1024)

	if nil != errA {
		ac.TpLogError("Failed to realloc buffer size: %s", errA.Message())
		ret = FAIL
		return
	}

	//Set the fields
	if errU := ubf.BAdd(ubftab.A_STATUSCODE, 0); nil != errU {
		ac.TpLogError("Failed to set A_STATUSCODE: %s", errU.Message())
		ret = FAIL
		return
	}

	if errU := ubf.BAdd(ubftab.A_STATUSMESSAGE, "OK"); nil != errU {
		ac.TpLogError("Failed to set A_STATUSMESSAGE: %s", errU.Message())
		ret = FAIL
		return
	}

	return
}

//Server init, called when process is booted
//@param ac ATMI Context
func Init(ac *atmi.ATMICtx) int {

	ac.TpLogWarn("Doing server init...")

	//Advertize service
	if err := ac.TpAdvertise("USERADDUBF", "USERADDUBF", USERADDUBF); err != nil {
		ac.TpLogError("Failed to Advertise: ATMI Error %d:[%s]\n",
			err.Code(), err.Message())
		return atmi.FAIL
	}

	return SUCCEED
}

//Server shutdown
//@param ac ATMI Context
func Uninit(ac *atmi.ATMICtx) {
	ac.TpLogWarn("Server is shutting down...")
}

//Executable main entry point
func main() {
	//Have some context
	ac, err := atmi.NewATMICtx()

	if nil != err {
		fmt.Fprintf(os.Stderr, "Failed to allocate new context: %s", err)
		os.Exit(atmi.FAIL)
	} else {
		//Run as server
		if err = ac.TpRun(Init, Uninit); nil != err {
			ac.TpLogError("Exit with failure")
			os.Exit(atmi.FAIL)
		} else {
			ac.TpLogInfo("Exit with success")
			os.Exit(atmi.SUCCEED)
		}
	}
}

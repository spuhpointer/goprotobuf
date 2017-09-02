package main

import (
	"fmt"
	"os"
	"userdet"

	atmi "github.com/endurox-dev/endurox-go"
	"github.com/golang/protobuf/proto"
)

const (
	SUCCEED = atmi.SUCCEED
	FAIL    = atmi.FAIL
)

var MCallNumber int = 0

//USERADD service
//@param ac ATMI Context
//@param svc Service call information
func USERADD(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {

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

	//Get CARRAY Handler
	requestBlob, _ := ac.CastToCarray(&svc.Data)

	var user userdet.Userdet

	err := proto.Unmarshal(requestBlob.GetBytes(), &user)

	if nil != err {
		ac.TpLogError("Failed to unmarshal User data: %s", err.Error())
		ret = FAIL
		return
	}

	//Test the data received

	if user.FirstName != "Jim" {
		ac.TpLogError("Expected FirstName = \"Jim\", but got: %s",
			user.FirstName)
		ret = FAIL
		return
	}

	if user.LastName != "Morrison" {
		ac.TpLogError("Expected LastName = \"Morrison\", but got: %s",
			user.LastName)
		ret = FAIL
		return
	}

	if user.Age != 27 {
		ac.TpLogError("Expected Age = \"27\", but got: %d",
			user.Age)
		ret = FAIL
		return
	}

	//Test data strings...
	for i := 0; i < int(user.ProfileTags); i++ {
		str := fmt.Sprintf("DATA STRING %d", i)
		if user.ProfileData[i] != str {

			ac.TpLogError("ProfileData[%d] expected to be [%s] but got [%s]",
				i, str, user.ProfileData[i])
			ret = FAIL
			return
		}
	}

	//Now send back the status code

	rsp := userdet.Resultdet{StatusCode: 0, StatusMessage: "OK"}

	data, err := proto.Marshal(&rsp)

	if err != nil {
		ac.TpLogError("Failed to marshal Resultdet: %s", err.Error())
		ret = FAIL
		return
	}

	//Will load the message in the same request object (as it goes back as rsp now)
	requestBlob.SetBytes(data)

	return
}

//Server init, called when process is booted
//@param ac ATMI Context
func Init(ac *atmi.ATMICtx) int {

	ac.TpLogWarn("Doing server init...")

	//Advertize service
	if err := ac.TpAdvertise("USERADD", "USERADD", USERADD); err != nil {
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

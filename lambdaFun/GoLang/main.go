package main

import (
	"audio_conference/lambdaFun/GoLang/alexa"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request alexa.Request) (alexa.Response, error) {

	if os.Getenv("GO_DEBUG_EN") == "1" {
		alexa.SetDebugGo(true)
	} else {
		log.Print("\r*DEBUG LOG OFF*\rEnvironment Variable GO_DEBUG_EN is either 0 or not set\r")
	}

	if os.Getenv("USERINFO_DEBUG_EN") == "1" {
		alexa.SetDebugUserInfo(true)
	} else {
		log.Print("\r*DEBUG LOG OFF*\rEnvironment Variable USERINFO_DEBUG_EN is either 0 or not set\r")
	}

	if os.Getenv("LOGTRACE_DEBUG_EN") == "1" {
		alexa.SetDebugLogTrace(true)
	} else {
		log.Print("\r*DEBUG LOG OFF*\rEnvironment Variable LOGTRACE_DEBUG_EN is either 0 or not set\r")
	}

	if os.Getenv("OPTIONS_DEBUG_EN") == "1" {
		alexa.SetDebugOptions(true)
	} else {
		log.Print("\r*DEBUG LOG OFF*\rEnvironment Variable OPTIONS_DEBUG_EN is either 0 or not set\r")
	}

	if alexa.GetDebugGo() {
		alexa.LogObject("Request", request)
	}

	return IntentDispatcher(request), nil

}

func main() {
	lambda.Start(Handler)
}

func HandleLaunchRequest(request alexa.Request) alexa.Response {

	var options map[string]interface{}
	options = make(map[string]interface{})

	//var LogTrace alexa.LogTrace
	var LogTrace = `LaunchRequest has been invoked`

	var speechText = "Welcome to the Audio Conference skill.  "
	speechText += "You can start a conference on a phone number or Be Anywhere device.  "
	speechText += "You can say, for example, "
	speechText += "ask audio conference to start a conference on "
	speechText += `<say-as interpret-as="telephone">2155551234</say-as>`
	speechText += ", or, "
	speechText += "ask audio conference to start a conference on My Cell. "
	options["speechText"] = speechText

	var cardContent = "Welcome to the Audio Conference skill.  "
	cardContent += "You can start a conference on a phone number or Be Anywhere device.  "
	cardContent += "You can say, for example, "
	cardContent += `'ask audio conference to start a conference on (215) 555-1234'`
	cardContent += ", or, "
	cardContent += `'ask audio conference to start a conference on My Cell'.`
	options["cardContent"] = cardContent

	options["imageObj"] = alexa.GetConferenceImg()

	var cardTitle = "Audio Conference"
	options["cardTitle"] = cardTitle

	options["endSession"] = true

	if alexa.GetDebugLogTrace() {
		alexa.LogObject("Trace", LogTrace)
	}

	if alexa.GetDebugUserInfo() {
		alexa.LogObject("User Info", alexa.UserInfo)
	}

	return alexa.BuildResponse(options)

}

func HandleStartConferenceIntent(request alexa.Request, AmazonID string) alexa.Response {

	var options map[string]interface{}
	options = make(map[string]interface{})

	var info = alexa.Info{}

	var LogTrace string

	LogTrace = "\rStartConference has been invoked"

	if alexa.UserInfoExists(AmazonID) {

		LogTrace += "\r     "
		LogTrace += "Current user info exists"

		info = alexa.GetUserInfo(AmazonID)

		slots := request.Body.Intent.Slots
		PNCur := slots["PN"].Value
		BeAnywhereCur := slots["BeAnywhere"].Value
		NumCur := slots["Num"].Value
		NumCheckCur := slots["NumCheck"].Value
		OrdinalCur := slots["Ordinal"].Value

		//Accuracy editing
		BeAnywhereCur = alexa.BeAnywhereHomomyn(BeAnywhereCur)

		//Accuracy editing
		if NumCur != "" {
			BeAnywhereCur = alexa.BeAnywhereNum(BeAnywhereCur, NumCur)
		} else if NumCheckCur != "" {
			BeAnywhereCur = alexa.BeAnywhereNumCheck(BeAnywhereCur, NumCheckCur)
		} else if OrdinalCur != "" {
			BeAnywhereCur = alexa.BeAnywhereOrdinal(BeAnywhereCur, OrdinalCur)
		}

		//Only one slot is filled
		if (PNCur == "" && BeAnywhereCur != "") || (PNCur != "" && BeAnywhereCur == "") {

			LogTrace += "\r          "
			LogTrace += "One slot has been filled"

			//A phone number is provided
			if PNCur != "" {

				LogTrace += "\r               "
				LogTrace += "A phone number has been provided"

				//The phone number passed verification
				if alexa.VerifyPN(PNCur) {

					LogTrace += "\r                    "
					LogTrace += "The phone number has passed verification"
					LogTrace += "\r                         "
					LogTrace += "Phone number: " + PNCur

					//PN_Pass
					options = OptionTemplates("PN_Pass", request, PNCur)

					info.StartIntent = true
					info.PN = PNCur

					//The phone number failed verification
				} else {

					LogTrace += "\r                    "
					LogTrace += "The phone number has failed verification"
					LogTrace += "\r                         "
					LogTrace += "Heard phone number: " + PNCur

					//PN_Fail
					options = OptionTemplates("PN_Fail", request, "")

				}

				//A Be Anywhere device was provided
			} else if BeAnywhereCur != "" {

				LogTrace += "\r               "
				LogTrace += "A Be Anywhere device has been provided"
				LogTrace += "\r                    "
				LogTrace += "Be Anywhere device: " + BeAnywhereCur

				//BeAnywhere
				options = OptionTemplates("BeAnywhere", request, BeAnywhereCur)

				info.StartIntent = true
				info.BeAnywhere = BeAnywhereCur

			}

			//Both slots are empty
		} else if PNCur == "" && BeAnywhereCur == "" {

			LogTrace += "\r          "
			LogTrace += "Neither slot has been filled"

			//Start_BothEmpty
			options = OptionTemplates("Start_BothEmpty", request, "")

			//Both slots are filled (RARE CASE)
		} else {

			LogTrace += "\r          "
			LogTrace += "Both slots have been filled"
			LogTrace += "\r               "
			LogTrace += "Heard phone number: " + PNCur
			LogTrace += "\r               "
			LogTrace += "Heard Be Anywhere device: " + BeAnywhereCur

			//Invalid
			options = OptionTemplates("Invalid", request, "")

		}

	} else {

		LogTrace += "\r     "
		LogTrace += "Current user info does not exist"

		options = OptionTemplates("Unknown", request, "")

		options["cardTitle"] = "ERROR: Audio Conference Start"

		info.StartIntent = false
		info.PN = ""
		info.BeAnywhere = ""

	}

	alexa.SetUserInfoObj(AmazonID, info)

	if alexa.GetDebugLogTrace() {
		alexa.LogObject("Trace", LogTrace)
	}

	if alexa.GetDebugUserInfo() {
		alexa.LogObject("Info", alexa.GetUserInfo(AmazonID))
	}

	return alexa.BuildResponse(options)

}

func HandleStartConferenceDeviceIntent(request alexa.Request, AmazonID string) alexa.Response {

	var options map[string]interface{}
	options = make(map[string]interface{})

	var info = alexa.Info{}

	var LogTrace = "\rStartConferenceDeviceIntent has been invoked"

	if alexa.UserInfoExists(AmazonID) {

		LogTrace += "\r     "
		LogTrace += "Current user info exists"

		info = alexa.GetUserInfo(AmazonID)

		//If session attributes exist and this intent was invoked by the intent StartConferenceIntent
		if request.Session.Attributes != nil && request.Session.Attributes["startConferenceIntent"] == true {

			LogTrace += "\r          "
			LogTrace += "In the current request, session attributes exist and this intent was invoked by StartConferenceIntent"

			slots := request.Body.Intent.Slots
			PNCur := slots["PN"].Value
			BeAnywhereCur := slots["BeAnywhere"].Value
			NumCur := slots["Num"].Value
			NumCheckCur := slots["NumCheck"].Value
			OrdinalCur := slots["Ordinal"].Value

			//Accuracy editing
			BeAnywhereCur = alexa.BeAnywhereHomomyn(BeAnywhereCur)

			//Accuracy editing
			if NumCur != "" {
				BeAnywhereCur = alexa.BeAnywhereNum(BeAnywhereCur, NumCur)
			} else if NumCheckCur != "" {
				BeAnywhereCur = alexa.BeAnywhereNumCheck(BeAnywhereCur, NumCheckCur)
			} else if OrdinalCur != "" {
				BeAnywhereCur = alexa.BeAnywhereOrdinal(BeAnywhereCur, OrdinalCur)
			}

			//Previously, a phone number was provided alongside the intent StartConferenceIntent and failed verification
			if _, ok := request.Session.Attributes["isPN"]; ok {

				LogTrace += "\r               "
				LogTrace += "In the previous request, a phone number was provided but failed verification"

				//If the slot is perceived as a phone number and passes verification
				if PNCur != "" && alexa.VerifyPN(PNCur) {

					LogTrace += "\r                    "
					LogTrace += "In the current request, the phone number has passed verification"
					LogTrace += "\r                         "
					LogTrace += "Phone number: " + PNCur

					//PN_Pass
					options = OptionTemplates("PN_Pass", request, PNCur)

					info.StartIntent = true
					info.PN = PNCur

					//If the slot is perceived as a Be Anywhere device OR the slot is perceived as a phone number but fails verification
				} else {

					LogTrace += "\r                    "
					LogTrace += "In the current request, the phone number has failed verificaiton"

					LogTrace += "\r                         "
					if PNCur != "" {
						LogTrace += "Heard phone number: " + PNCur
					} else {
						LogTrace += "Heard phone number: " + BeAnywhereCur
					}

					//PN_Fail
					options = OptionTemplates("PN_Fail", request, "")

					//IN CASE OF ERROR: DO NOT copy over Session Attributes in this block
					//Explicitly empty Session Attributes

				}

				//Previously, no slots were provided or an invalid request was made (e.g. both slots filled)
			} else {

				LogTrace += "\r               "
				LogTrace += "In the previous request, both slots were empty or an invalid request was made"

				//One of the slots is filled
				if (PNCur == "" && BeAnywhereCur != "") || (PNCur != "" && BeAnywhereCur == "") {

					LogTrace += "\r                    "
					LogTrace += "In the current request, one slot has been filled"

					//If a phone number is provided
					if PNCur != "" {

						LogTrace += "\r                         "
						LogTrace += "A phone number has been provided"

						//If the phone number passed verification
						if alexa.VerifyPN(PNCur) {

							LogTrace += "\r                              "
							LogTrace += "The phone number has passed verification"
							LogTrace += "\r                                   "
							LogTrace += "Phone number: " + PNCur

							//PN_Pass
							options = OptionTemplates("PN_Pass", request, PNCur)

							info.StartIntent = true
							info.PN = PNCur

							//If the phone number failed verification
						} else {

							LogTrace += "\r                              "
							LogTrace += "The phone number has failed verification"
							LogTrace += "\r                                   "
							LogTrace += "Heard phone number: " + PNCur

							//PN_Fail
							options = OptionTemplates("PN_Fail", request, "")

						}

						//If a Be Anywhere device is provided
					} else if BeAnywhereCur != "" {

						LogTrace += "\r                         "
						LogTrace += "A Be Anywhere device has been provided"
						LogTrace += "\r                              "
						LogTrace += "Be Anywhere device: " + BeAnywhereCur

						//BeAnywhere
						options = OptionTemplates("BeAnywhere", request, BeAnywhereCur)

						info.StartIntent = true
						info.BeAnywhere = BeAnywhereCur

					}

					//Both slots are either empty or filled (RARE CASE)
				} else {

					LogTrace += "\r                    "
					LogTrace += "In the current request, neither slot has been filled or both slots have been filled"
					LogTrace += "\r                         "
					LogTrace += "Heard phone number: " + PNCur
					LogTrace += "\r                         "
					LogTrace += "Heard Be Anywhere device: " + BeAnywhereCur

					//Invalid
					options = OptionTemplates("Invalid", request, "")

				}

			}

			//If there are no session attributes, this intent was NOT invoked by the intent StartConferenceIntent, or an invalid request has been made
		} else {

			LogTrace += "\r          "
			LogTrace += "In the current request, there are either no session attributes, this intent was not invoked by StartConferenceIntent, or an invalid request has been made"

			//Incorrect
			options = OptionTemplates("Incorrect", request, "")

		}

	} else {

		LogTrace += "\r     "
		LogTrace += "Current user info does not exist"

		options = OptionTemplates("Unknown", request, "")

		options["cardTitle"] = "ERROR: Audio Conference Start"

		info.StartIntent = false
		info.PN = ""
		info.BeAnywhere = ""

	}

	alexa.SetUserInfoObj(AmazonID, info)

	if alexa.GetDebugLogTrace() {
		alexa.LogObject("Trace", LogTrace)
	}

	if alexa.GetDebugUserInfo() {
		alexa.LogObject("Info", alexa.GetUserInfo(AmazonID))
	}

	return alexa.BuildResponse(options)

}

func HandleStopConferenceIntent(request alexa.Request, AmazonID string) alexa.Response {

	var options map[string]interface{}
	options = make(map[string]interface{})

	var info = alexa.Info{}

	//var LogTrace alexa.LogTrace
	var LogTrace = "StopConferenceIntent has been invoked"

	if alexa.UserInfoExists(AmazonID) {

		LogTrace += "\r     "
		LogTrace += "Current user info exists"

		info = alexa.GetUserInfo(AmazonID)

		if info.StartIntent {

			LogTrace += "\r          "
			LogTrace += "This intent was invoked after the StartIntent was invoked"

			slots := request.Body.Intent.Slots
			PNCur := slots["PN"].Value
			BeAnywhereCur := slots["BeAnywhere"].Value
			NumCur := slots["Num"].Value
			NumCheckCur := slots["NumCheck"].Value
			OrdinalCur := slots["Ordinal"].Value

			//Accuracy editing
			BeAnywhereCur = alexa.BeAnywhereHomomyn(BeAnywhereCur)

			//Accuracy editing
			if NumCur != "" {
				BeAnywhereCur = alexa.BeAnywhereNum(BeAnywhereCur, NumCur)
			} else if NumCheckCur != "" {
				BeAnywhereCur = alexa.BeAnywhereNumCheck(BeAnywhereCur, NumCheckCur)
			} else if OrdinalCur != "" {
				BeAnywhereCur = alexa.BeAnywhereOrdinal(BeAnywhereCur, OrdinalCur)
			}

			//A phone number was provided as a slot
			if PNCur != "" {

				LogTrace += "\r               "
				LogTrace += "A phone number has been provided as a slot"
				LogTrace += "\r               "
				LogTrace += "Phone number: " + PNCur

				var speechText = `Your conference was stopped on <say-as interpret-as="telephone">` + PNCur + "</say-as>. "
				options["speechText"] = speechText

				var cardContent = "Your conference was stopped on " + PNCur + ". "
				options["cardContent"] = cardContent

				options["imageObj"] = alexa.GetPhoneStopImg()

				var cardTitle = "Audio Conference Stop"
				options["cardTitle"] = cardTitle

				/* DO THIS ON MONDAY!!!
				if (alexa.VerifyPN(PNCur)) {



				} else {



				}
				*/

				//A Be Anywhere device has been provided as a slot
			} else if BeAnywhereCur != "" {

				LogTrace += "\r               "
				LogTrace += "A Be Anywhere device has been provided as a slot"
				LogTrace += "\r                    "
				LogTrace += "Be Anywhere: " + BeAnywhereCur

				var text = "Your conference was stopped on " + BeAnywhereCur + ". "
				options["speechText"] = text
				options["cardContent"] = text

				options["imageObj"] = alexa.GetPhoneStopImg()

				var cardTitle = "Audio Conference Stop"
				options["cardTitle"] = cardTitle

			} else if info.PN != "" {

				LogTrace += "\r               "
				LogTrace += "Previously, a phone number was provided to start the conference"
				LogTrace += "\r                    "
				LogTrace += "Phone number: " + PNCur

				var speechText = `Your conference was stopped on <say-as interpret-as="telephone">` + info.PN + "</say-as>. "
				options["speechText"] = speechText

				var cardContent = "Your conference was stopped on " + info.PN + ". "
				options["cardContent"] = cardContent

				options["imageObj"] = alexa.GetPhoneStopImg()

				var cardTitle = "Audio Conference Stop"
				options["cardTitle"] = cardTitle

			} else if info.BeAnywhere != "" {

				LogTrace += "\r               "
				LogTrace += "Previously, a Be Anywhere device was provided to start the conference"
				LogTrace += "\r                    "
				LogTrace += "Be Anywhere: " + info.BeAnywhere

				var text = "Your conference was stopped on " + info.BeAnywhere + ". "
				options["speechText"] = text
				options["cardContent"] = text

				options["imageObj"] = alexa.GetPhoneStopImg()

				var cardTitle = "Audio Conference Stop"
				options["cardTitle"] = cardTitle

				//Neither slot was filled
			} else {

				LogTrace += "\r               "
				LogTrace += "Neither slot has been filled"

				var speechText = "Invalid option.  "
				speechText += "To stop the conference, please provide a valid telephone number or BeAnywhere device. "
				options["speechText"] = speechText

				var cardContent = "Invalid option.  "
				cardContent += "To stop the conference, please provide a valid telephone number or BeAnywhere device. "
				options["cardContent"] = cardContent

				options["imageObj"] = alexa.GetPhoneErrorImg()

				options["cardTitle"] = "ERROR: Audio Conference Stop"

			}

		} else {

			LogTrace += "\r          "
			LogTrace += "This intent was not invoked after the StartIntent was invoked"

			var speechText = "Incorrect usage.  "
			speechText += "To stop a conference, a conference must first be started. "
			options["speechText"] = speechText

			var cardContent = "Incorrect usage.  "
			cardContent += "To stop a conference, a conference must first be started."
			options["cardContent"] = cardContent

			options["imageObj"] = alexa.GetPhoneErrorImg()

			options["cardTitle"] = "ERROR: Audio Conference Stop"

		}

	} else {

		LogTrace += "\r     "
		LogTrace += "Current user info doees not exist"

		options = OptionTemplates("Unknown", request, "")

		options["cardTitle"] = "ERROR: Audio Conference Stop"

		info.StartIntent = false
		info.PN = ""
		info.BeAnywhere = ""

	}

	options["endSession"] = true

	info.StartIntent = false
	info.PN = ""
	info.BeAnywhere = ""

	alexa.SetUserInfoObj(AmazonID, info)

	if alexa.GetDebugLogTrace() {
		alexa.LogObject("Trace", LogTrace)
	}

	if alexa.GetDebugUserInfo() {
		alexa.LogObject("Info", alexa.GetUserInfo(AmazonID))
	}

	return alexa.BuildResponse(options)

}

func IntentDispatcher(request alexa.Request) alexa.Response {

	var AmazonID = request.Session.User.UserID

	if !(alexa.UserInfoExists(AmazonID)) {
		alexa.SetUserInfo(AmazonID, false, "", "")
	}

	var response alexa.Response

	if request.Body.Type == "LaunchRequest" {
		response = HandleLaunchRequest(request)
	} else if request.Body.Type == "IntentRequest" {

		switch request.Body.Intent.Name {
		case "StartConferenceIntent":
			response = HandleStartConferenceIntent(request, AmazonID)
		case "StartConferenceDeviceIntent":
			response = HandleStartConferenceDeviceIntent(request, AmazonID)
		case "StopConferenceIntent":
			response = HandleStopConferenceIntent(request, AmazonID)
		}

	}

	return response

}

func OptionTemplates(name string, request alexa.Request, deviceCur string) map[string]interface{} {

	var options map[string]interface{}
	options = make(map[string]interface{})

	options["session"] = request.Session

	switch name {

	case "PN_Pass":

		var speechText = "Your conference was started on "
		speechText += `<say-as interpret-as="telephone">` + deviceCur + "</say-as>. "
		options["speechText"] = speechText

		var cardContent = "Your conference was started on " + alexa.FormatPN(deviceCur) + ". "
		options["cardContent"] = cardContent

		options["imageObj"] = alexa.GetPhoneStartImg()

		var cardTitle = "Audio Conference Start"
		options["cardTitle"] = cardTitle

		options["endSession"] = true

	case "PN_Fail":

		var speechText = "Invalid phone number.  "
		speechText += "Please provide a valid 10-digit phone number starting with the area code.  "
		speechText += "What phone number would you like to start your conference on? "
		options["speechText"] = speechText

		var repromptText = "For example, you could say "
		repromptText += `<say-as interpret-as="telephone">2155551234</say-as>. `
		repromptText += "What phone number would you like to start your conference on? "
		options["repromptText"] = repromptText

		var cardContent = "Invalid phone number.  "
		cardContent += "Please provide a valid 10-digit phone number starting with the area code.  "
		cardContent += "For example, you could say "
		cardContent += `'(215) 555-1234'.`
		options["cardContent"] = cardContent

		options["imageObj"] = alexa.GetPhoneErrorImg()

		var cardTitle = "ERROR: Audio Conference Start"
		options["cardTitle"] = cardTitle

		session := request.Session
		var attributes map[string]interface{}
		attributes = make(map[string]interface{})
		attributes["startConferenceIntent"] = true
		attributes["isPN"] = true
		session.Attributes = attributes
		options["session"] = session

		options["endSession"] = false

	case "BeAnywhere":

		var speechText = "Your conference was started on " + deviceCur + ". "
		options["speechText"] = speechText

		var cardContent = "Your conference was started on " + deviceCur + "."
		options["cardContent"] = cardContent

		options["imageObj"] = alexa.GetPhoneStartImg()

		var cardTitle = "Audio Conference Start"
		options["cardTitle"] = cardTitle

		options["endSession"] = true

	case "Start_BothEmpty":

		var speechText = "What device would you like to start your conference on? "
		options["speechText"] = speechText

		var repromptText = "You can say a phone number, such as "
		repromptText += `<say-as interpret-as="telephone">2155551234</say-as>`
		repromptText += ", or "
		repromptText += "say a Be Anywhere device, such as "
		repromptText += "My Cell. "
		options["repromptText"] = repromptText

		var cardContent = "You can say a phone number, such as "
		cardContent += `'(215) 555-1234'`
		cardContent += ", or "
		cardContent += "say a Be Anywhere device, such as "
		cardContent += `'My Cell'.`
		options["cardContent"] = cardContent

		options["imageObj"] = alexa.GetQuestionImg()
		options["endSession"] = false

		var cardTitle = "HELP: Audio Conference Start"
		options["cardTitle"] = cardTitle

		session := request.Session
		var attributes map[string]interface{}
		attributes = make(map[string]interface{})
		attributes["startConferenceIntent"] = true
		session.Attributes = attributes
		options["session"] = session

		options["endSession"] = false

	case "Invalid":

		var speechText = "Invalid request.  "
		speechText += "Please provide a valid phone number or Be Anywhere device.  "
		speechText += "What device would you like to start your conference on? "
		options["speechText"] = speechText

		var repromptText = "You can say a phone number, such as "
		repromptText += `<say-as interpret-as="telephone">2155551234</say-as>`
		repromptText += ", or "
		repromptText += "say a Be Anywhere device, such as "
		repromptText += "My Cell. "
		repromptText += "What device would you like to start your conference on? "
		options["repromptText"] = repromptText

		var cardContent = "Invalid request.  "
		cardContent += "Please provide a valid phone number or Be Anywhere device.  "
		cardContent = "You can say a phone number, such as "
		cardContent += `'(215) 555-1234'`
		cardContent += ", or "
		cardContent += "say a Be Anywhere device, such as "
		cardContent += `'My Cell'.`
		options["cardContent"] = cardContent

		options["imageObj"] = alexa.GetPhoneErrorImg()

		var cardTitle = "ERROR: Audio Conference Start"
		options["cardTitle"] = cardTitle

		session := request.Session
		var attributes map[string]interface{}
		attributes = make(map[string]interface{})
		attributes["startConferenceIntent"] = true
		session.Attributes = attributes
		options["session"] = session

		options["endSession"] = false

	case "Incorrect":

		var speechText = "Incorrect usage.  "
		speechText += "To start a conference, please provide a valid phone number or Be Anywhere device.  "
		speechText += "You can say, for example, "
		speechText += "ask audio conference to start a conference on "
		speechText += `<say-as interpret-as="telephone">2155551234</say-as>`
		speechText += ", or, "
		speechText += "ask audio conference to start a conference on My Cell. "
		options["speechText"] = speechText

		var cardContent = "Incorrect usage.  "
		cardContent += "To start a conference, please provide a valid phone number or Be Anywhere device.  "
		cardContent += "You can say, for example, "
		cardContent += `'ask audio conference to start a conference on (215) 555-1234'`
		cardContent += ", or, "
		cardContent += `'ask audio conference to start a conference on My Cell'.`
		options["cardContent"] = cardContent

		options["imageObj"] = alexa.GetPhoneErrorImg()

		var cardTitle = "ERROR: Audio Conference Start"
		options["cardTitle"] = cardTitle

		options["endSession"] = true

	case "Unknown":

		var speechText = "Unknown usage error.  "
		speechText += "The data associated with your Amazon ID for this skill was not able to be retrieved.  "
		speechText += "Your data has been refreshed.  "
		speechText += "Please start over and try again! "
		options["speechText"] = speechText

		var cardContext = "Unknown usage error.  "
		cardContext += "The data associated with your Amazon ID for this skill was not able to be retrieved.  "
		cardContext += "Your data has been refreshed.  "
		cardContext += "Please start over and try again!"
		options["cardContext"] = cardContext

		options["imageObj"] = alexa.GetPhoneErrorImg

		options["endSession"] = true

	}

	return options

}

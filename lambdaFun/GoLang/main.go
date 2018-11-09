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

	var speechText = "Welcome to the Audio Conference skill.  "
	speechText += "You can say, for example, "
	speechText += "ask audio conference to start a conference on "
	speechText += `<say-as interpret-as="telephone">2155551234</say-as>`
	speechText += ", or, "
	speechText += "ask audio conference to start a conference on My Cell. "
	options["speechText"] = speechText

	var cardContent = "Welcome to the Audio Conference skill.  "
	cardContent += "You can say, for example, "
	cardContent += `'ask audio conference to start a conference on (215) 555-1234'`
	cardContent += ", or, "
	cardContent += `'ask audio conference to start a conference on My Cell'.`
	options["cardContent"] = cardContent

	options["imageObj"] = alexa.GetConferenceImg()

	var cardTitle = "Audio Conference"
	options["cardTitle"] = cardTitle

	options["endSession"] = true

	if alexa.GetDebugOptions() {
		alexa.LogObject("Options (StartIntent)", options)
	}

	return alexa.BuildResponse(options)

}

func HandleStartIntent(request alexa.Request) alexa.Response {

	var options map[string]interface{}
	options = make(map[string]interface{})

	slots := request.Body.Intent.Slots
	PNCur := slots["PN"].Value
	BeAnywhereCur := slots["BeAnywhere"].Value

	//Only one slot is filled
	if (PNCur == "" && BeAnywhereCur != "") || (PNCur != "" && BeAnywhereCur == "") {

		alexa.LogObject("StartIntent invoked: One Slot provided", nil)

		//A phone number is provided
		if PNCur != "" {

			alexa.LogObject("Current Slot: PN", nil)

			//The phone number passed verification
			if alexa.VerifyPN(PNCur) {

				//PN_Pass
				options = OptionTemplates("PN_Pass", request, PNCur)

				alexa.LogObject("PN Passed Verification", nil)

				//The phone number failed verification
			} else {

				//PN_Fail
				options = OptionTemplates("PN_Fail", request, "")

				alexa.LogObject("PN Failed Verification", nil)

			}

			//A Be Anywhere device was provided
		} else if BeAnywhereCur != "" {

			//BeAnywhere
			options = OptionTemplates("BeAnywhere", request, BeAnywhereCur)

			alexa.LogObject("Current Slot: Be Anywhere", nil)

		}

		//Both slots are empty
	} else if PNCur == "" && BeAnywhereCur == "" {

		//Start_BothEmpty
		options = OptionTemplates("Start_BothEmpty", request, "")

		alexa.LogObject("StartIntent invoked: Both Slots Empty", nil)

		//Both slots are filled
	} else {

		//Invalid
		options = OptionTemplates("Invalid", request, "")

		alexa.LogObject("StartIntent invoked: Both Slots Filled", nil)

	}

	return alexa.BuildResponse(options)

}

func HandleStartDeviceIntent(request alexa.Request) alexa.Response {

	var options map[string]interface{}
	options = make(map[string]interface{})

	//If session attributes exist and this intent was invoked by the intent StartIntent
	if request.Session.Attributes != nil && request.Session.Attributes["startIntent"] == true {

		alexa.LogObject("StartDeviceIntent invoked: Session Attributes Exist, Invoked by StartIntent", nil)

		slots := request.Body.Intent.Slots
		PNCur := slots["PN"].Value
		BeAnywhereCur := slots["BeAnywhere"].Value

		//Previously, a phone number was provided alongside the intent StartIntent and failed verification
		if _, ok := request.Session.Attributes["isPN"]; ok {

			alexa.LogObject("Previous invocation: PN Provided & Failed Verification or Be Anywhere Provided instad of PN (PN Provided initially & Failed Verification)", nil)

			//If the slot is perceived as a phone number and passes verification
			if PNCur != "" && alexa.VerifyPN(PNCur) {

				//PN_Pass
				options = OptionTemplates("PN_Pass", request, PNCur)

				alexa.LogObject("Current Slot: PN & Passed Verification", nil)

				//If the slot is perceived as a BeAnywhere device OR the slot is perceived as a phone number but fails verification
			} else {

				//PN_Fail
				options = OptionTemplates("PN_Fail", request, "")

				alexa.LogObject("Current Slot: Be Anywhere / PN & Failed Verification", nil)

				//IN CASE OF ERROR: DO NOT copy over Session Attributes in this block
				//Explicitly empty Session Attributes

			}

			//Previously, no slots were provided or an invalid request was made (e.g. both slots filled)
		} else {

			alexa.LogObject("Previous invocation: No Slots / Invalid Request", nil)

			//One of the slots is filled
			if (PNCur == "" && BeAnywhereCur != "") || (PNCur != "" && BeAnywhereCur == "") {

				alexa.LogObject("Current Slot: One Slot", nil)

				//If a phone number is provided
				if PNCur != "" {

					alexa.LogObject("One Slot is PN", nil)

					//If the phone number passed verification
					if alexa.VerifyPN(PNCur) {

						//PN_Pass
						options = OptionTemplates("PN_Pass", request, PNCur)

						//If the phone number failed verification
					} else {

						//PN_Fail
						options = OptionTemplates("PN_Fail", request, "")

						alexa.LogObject("PN Failed Verification", nil)

					}

					//If a Be Anywhere device is provided
				} else if BeAnywhereCur != "" {

					//BeAnywhere
					options = OptionTemplates("BeAnywhere", request, BeAnywhereCur)

					alexa.LogObject("One Slot is Be Anywhere", nil)

				}

				//Both slots are either empty or filled
			} else {

				//Invalid
				options = OptionTemplates("Invalid", request, "")

				alexa.LogObject("Current Slot: Both Slots Provided / No Slots Provided", nil)

			}

		}

		//If there are no session attributes or this intent was NOT invoked by the intent StartIntent
	} else {

		//Incorrect
		options = OptionTemplates("Incorrect", request, "")

		alexa.LogObject("StartDeviceIntent invocation: No Session Attributes / Not Invoked by StartIntent", nil)

	}

	return alexa.BuildResponse(options)

}

func HandleStopIntent(request alexa.Request) alexa.Response {

	var options map[string]interface{}
	options = make(map[string]interface{})

	slots := request.Body.Intent.Slots
	PNCur := slots["PN"].Value
	BeAnywhereCur := slots["BeAnywhere"].Value

	if PNCur != "" {

		var speechText = `Your conference was stopped on <say-as interpret-as="telephone">` + PNCur + "</say-as>. "
		options["speechText"] = speechText

		var cardContent = "Your conference was stopped on " + PNCur + ". "
		options["cardContent"] = cardContent

	} else if BeAnywhereCur != "" {

		var text = "Your conference was stopped on " + BeAnywhereCur + ". "
		options["speechText"] = text
		options["cardContent"] = text

	} else {

		var text = "Your conference was stopped. "
		options["speechText"] = text
		options["cardContent"] = text

	}

	options["imageObj"] = alexa.GetPhoneStopImg()

	var cardTitle = "Audio Conference Stop"
	options["cardTitle"] = cardTitle

	options["endSession"] = true

	if alexa.GetDebugOptions() {
		alexa.LogObject("Options (StopIntent)", options)
	}

	return alexa.BuildResponse(options)

}

func IntentDispatcher(request alexa.Request) alexa.Response {

	var response alexa.Response

	if request.Body.Type == "LaunchRequest" {
		response = HandleLaunchRequest(request)
	} else if request.Body.Type == "IntentRequest" {

		switch request.Body.Intent.Name {
		case "StartIntent":
			response = HandleStartIntent(request)
		case "StartDeviceIntent":
			response = HandleStartDeviceIntent(request)
		case "StopIntent":
			response = HandleStopIntent(request)
		}

	}

	//alexa.ClearOptionsMap()

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
		attributes["startIntent"] = true
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
		attributes["startIntent"] = true
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
		attributes["startIntent"] = true
		session.Attributes = attributes
		options["session"] = session

		options["endSession"] = false

	case "Incorrect":

		var speechText = "Incorrect usage.  "
		speechText += "To start a conference, please provide a valid phone number or Be Anywhere device. "
		options["speechText"] = speechText

		var cardContent = "Incorrect usage.  "
		cardContent += "To start a conference, please provide a valid phone number or Be Anywhere device."
		options["cardContent"] = cardContent

		options["imageObj"] = alexa.GetPhoneErrorImg()

		var cardTitle = "ERROR: Audio Conference Start"
		options["cardTitle"] = cardTitle

		options["endSession"] = true

	}

	if alexa.GetDebugOptions() {
		alexa.LogObject("Options", options)
	}

	return options

}

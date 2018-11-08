package main

import (
	//"strings"

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
		log.Print("\r*DEBUG LOG OFF*\rEnvironment Variable OPTIONS_DEBUG_EN is either 0 or not set\r")
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

	var builder alexa.SSMLBuilder

	builder.Say("Welcome to the Audio Conference skill.  ")
	builder.Say("Using this skill, you can start an audio conference on a telephone number or a Be Anywhere device.  ")
	builder.Say("You can say, for example, ask audio conference to start a conference on ")
	builder.PN("2155551234")
	builder.Say(", or, ask audio conference to start a conference on My Cell. ")

	return alexa.NewSSMLResponse("LaunchRequest", builder.Build(), "", true, request.Session)

}

func HandleStartIntent(request alexa.Request) alexa.Response {

	//var options map[interface{}]interface{}
	//options = make(map[interface{}]interface{})

	var options map[string]interface{}
	options = make(map[string]interface{})

	//var builder alexa.SSMLBuilder

	slots := request.Body.Intent.Slots
	PNCur := slots["PN"].Value
	BeAnywhereCur := slots["BeAnywhere"].Value

	if PNCur != "" || BeAnywhereCur != "" {

		if PNCur != "" {

			var speechText = `Your conference was started on <say-as interpret-as="telephone">` + PNCur + "</say-as>. "
			options["speechText"] = speechText

			var cardContent = "Your conference was started on " + PNCur + ". "
			options["cardContent"] = cardContent

			/*
				builder.Say("Your conference was started on ")
				builder.PN(PNCur)
				builder.Say(". ")
			*/

		} else if BeAnywhereCur != "" {

			//builder.Say("Your conference was started on " + BeAnywhereCur + ". ")

			var text = "Your conference was started on " + BeAnywhereCur + ". "
			options["speechText"] = text
			options["cardContent"] = text

		}

		options["imageObj"] = alexa.GetPhoneStartImg()
		options["endSession"] = true

		//return alexa.NewSSMLResponse("StartIntent Device", builder.Build(), "", true, request.Session)

	} else {

		session := request.Session

		options["session"] = request.Session

		var builderReprompt alexa.SSMLBuilder

		var speechText = "What device would you like to start your conference on? "
		options["speechText"] = speechText

		var repromptText = `You can say a telephone number, such as <say-as interpret-as=\"telephone\">2155551234</say-as>, or say a Be Anywhere device, such as My Cell.`
		options["repromptText"] = repromptText

		var cardContent = "You can say a telephone number, such as 2155551234, or say a Be Anywhere device, such as My Cell."
		options["cardContent"] = cardContent

		options["imageObj"] = alexa.GetQuestionImg()
		options["endSession"] = false

		/*
			builder.Say("What device would you like to start your conference on? ")

			builderReprompt.Say("You can say a telephone number, such as ")
			builderReprompt.PN("2155551234")
			builderReprompt.Say(", or say a Be Anywhere device, such as My Cell. ")
		*/

		var attributes map[string]interface{}
		attributes = make(map[string]interface{})
		attributes["startIntent"] = true
		session.Attributes = attributes
		options["session"] = session

		options["endSession"] = false

		//return alexa.NewSSMLResponse("StartIntent NoDevices", builder.Build(), builderReprompt.Build(), false, session)

	}

	var cardTitle = "Audio Conference Start"
	options["cardTitle"] = cardTitle

	if alexa.GetDebugOptions() {
		alexa.LogObject("Options (StartIntent)", options)
	}

	return alexa.BuildResponse(options)

}

func HandleStartDeviceIntent(request alexa.Request) alexa.Response {

	var builder alexa.SSMLBuilder

	if request.Session.Attributes != nil && request.Session.Attributes["startIntent"] == true {

		slots := request.Body.Intent.Slots
		PNCur := slots["PN"].Value
		BeAnywhereCur := slots["BeAnywhere"].Value

		if PNCur != "" {

			builder.Say("Your conference was started on ")
			builder.PN(PNCur)
			builder.Say(". ")

		} else if BeAnywhereCur != "" {
			builder.Say("Your conference was started on " + BeAnywhereCur + ". ")
		}

	} else {
		builder.Say("Incorrect usage.  To start a conference, please provide a valid telephone number or Be Anywhere device. ")
	}

	return alexa.NewSSMLResponse("StartDeviceIntent Device", builder.Build(), "", true, request.Session)

}

func HandleStopIntent(request alexa.Request) alexa.Response {

	//var builder alexa.SSMLBuilder

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

		/*
			builder.Say("Your conference was stopped on ")
			builder.PN(PNCur)
			builder.Say(". ")
		*/

	} else if BeAnywhereCur != "" {

		var text = "Your conference was stopped on " + BeAnywhereCur + ". "
		options["speechText"] = text
		options["cardContent"] = text

		//builder.Say("Your conference was stopped on " + BeAnywhereCur + ". ")

	} else {

		var text = "Your conference was stopped. "
		options["speechText"] = text
		options["cardContent"] = text

		//builder.Say("Your conference was stopped. ")

	}

	options["imageObj"] = alexa.GetPhoneStopImg()

	var cardTitle = "Audio Conference Stop"
	options["cardTitle"] = cardTitle

	if alexa.GetDebugOptions() {
		alexa.LogObject("Options (StopIntent)", options)
	}

	//return alexa.NewSSMLResponse("StopIntent", builder.Build(), "", true, request.Session)

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

	return response

}

/*
package main

import (
	"fmt"
	//"encoding/json"
	"strings"
)

type Info struct {
	StartIntent bool
	PN string
	BeAnywhere string
}

type UserInfo struct {
	UserId string
	Info Info
}

type CardImg struct {
	Small string
	Large string
}

//<a href="https://www.iconfinder.com/icons/309047/conference_group_people_users_icon" target="_blank">"Conference, group, people, users icon"</a> by <a href="https://www.iconfinder.com/visualpharm" target="_blank">Ivan Boyko</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0" target="_blank">CC BY 3.0</a>
//"Conference, group, people, users icon" (https://www.iconfinder.com/icons/309047/conference_group_people_users_icon) by Ivan Boyko (https://www.iconfinder.com/visualpharm) is licensed under CC BY 3.0 (http://creativecommons.org/licenses/by/3.0)
var conferenceImg CardImg
conferenceImg.Small = "https://s3.amazonaws.com/audio-conference/images/conferenceSmall.png"
conferenceImg.Large = "https://s3.amazonaws.com/audio-conference/images/conferenceLarge.png"

//<a href="https://www.iconfinder.com/icons/3324959/outgoing_phone_icon" target="_blank">"Outgoing, phone icon"</a> by <a href="https://www.iconfinder.com/colebemis" target="_blank">Cole Bemis</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0" target="_blank">CC BY 3.0</a>
//"Outgoing, phone icon" (https://www.iconfinder.com/icons/3324959/outgoing_phone_icon) by Cole Bemis (https://www.iconfinder.com/colebemis) is licensed under CC BY 3.0 (http://creativecommons.org/licenses/by/3.0)
var phoneStartImg CardImg
phoneStartImg.Small = "https://s3.amazonaws.com/audio-conference/images/phoneStartSmall.png"
phoneStartImg.Large = "https://s3.amazonaws.com/audio-conference/images/phoneStartLarge.png"

//<a href="https://www.iconfinder.com/icons/3324961/missed_phone_icon" target="_blank">"Missed, phone icon"</a> by <a href="https://www.iconfinder.com/colebemis" target="_blank">Cole Bemis</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0" target="_blank">CC BY 3.0</a>
//"Missed, phone icon" (https://www.iconfinder.com/icons/3324961/missed_phone_icon) by Cole Bemis (https://www.iconfinder.com/colebemis) is licensed under CC BY 3.0 (http://creativecommons.org/licenses/by/3.0)
var phoneStopImg CardImg
phoneStopImg.Small = "https://s3.amazonaws.com/audio-conference/images/phoneStopSmall.png"
phoneStopImg.Large = "https://s3.amazonaws.com/audio-conference/images/phoneStopLarge.png"

//<a href="https://www.iconfinder.com/icons/3324960/off_phone_icon" target="_blank">"Off, phone icon"</a> by <a href="https://www.iconfinder.com/colebemis" target="_blank">Cole Bemis</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0" target="_blank">CC BY 3.0</a>
//"Off, phone icon" (https://www.iconfinder.com/icons/3324960/off_phone_icon) by Cole Bemis (https://www.iconfinder.com/colebemis) is licensed under CC BY 3.0 (http://creativecommons.org/licenses/by/3.0)
var phoneErrorImg CardImg
phoneErrorImg.Small = "https://s3.amazonaws.com/audio-conference/images/phoneErrorSmall.png"
phoneErrorImg.Large = "https://s3.amazonaws.com/audio-conference/images/phoneErrorLarge.png"

//<a href="https://www.iconfinder.com/icons/183285/help_mark_question_icon" target="_blank">"Help, mark, question icon"</a> by <a href="https://www.iconfinder.com/yanlu" target="_blank">Yannick Lung</a>
//"Help, mark, question icon" (https://www.iconfinder.com/icons/183285/help_mark_question_icon) by Yannick Lung (https://www.iconfinder.com/yanlu)
var questionImg CardImg
questionImg.Small = "https://s3.amazonaws.com/audio-conference/images/questionSmall.png"
questionImg.Large = "https://s3.amazonaws.com/audio-conference/images/questionLarge.png"


func main() {
	fmt.Println("Hello!")
}
*/

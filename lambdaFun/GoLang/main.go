package main

import (
	//AWS Lambda Context Object in Go
	"context"
	"fmt"
	"log"

	//Implements the Lambda programming model for Go
	"github.com/aws/aws-lambda-go/lambda" //go get -u github.com/aws/aws-lambda-go/lambda

	//Implements a deep pretty printer for Go data structures to aid in debugging
	"github.com/davecgh/go-spew/spew" //go get -u github.com/davecgh/go-spew/spew
)

//Input object
type AlexaRequest struct {
	//string "Version": skill version number
	Version string `json:"version"`
	//object "Request": current request data
	Request struct {
		//string "Type": current request type
		Type string `json:"type"`
		//string "Time": current request UTC time
		Time string `json:"timestamp"`
		//object "intent": current intent data
		Intent struct {
			//string "Name": current intent name
			Name string `json:"name"`
			//string "ConfirmationStatus": if the user needs to provide confirmation upon invocation of the current intent
			ConfirmationStatus string `json:"confirmationstatus"`
		} `json:"intent"`
	} `json:"request"`
}

//Output object
type AlexaResponse struct {
	//string "Version": skill version number
	Version string `json:"version"`
	//object "Response": current response data
	Response struct {
		//object "OutputSpeech": current output speech data
		OutputSpeech struct {
			//string "Type": current data type of the output speech
			Type string `json:"type"`
			//string "Text": current text of the output speech
			Text string `json:"text"`
		} `json:"outputSpeech"`
	} `json:"response"`
}

func CreateResponse() *AlexaResponse {
	var resp AlexaResponse
	resp.Version = "1.0"
	resp.Response.OutputSpeech.Type = "PlainText"
	resp.Response.OutputSpeech.Text = "Hello.  Please override this default output."
	return &resp
}

func (resp *AlexaResponse) Say(text string) {
	resp.Response.OutputSpeech.Text = text
}

/*
	name: HandleRequest()
	description: The Lambda handler signature, includes the code which will be executed
	parameters:
		ctx context.Context: Provides runtime information for your Lambda function invocation.
			ctx: The variable declared to leverage the information avaiable via AWS Lambda Context Object in Go
		i AlexaRequest: The input request
	return:
		AlexaResponse: The output response
		error: Returns standard error information
*/
func HandleRequest(ctx context.Context, i AlexaRequest) (AlexaResponse, error) {
	// Use Spew to output the request for debugging purposes:
	fmt.Println("---- Dumping Input Map: ----")
	spew.Dump(i)
	fmt.Println("---- Done. ----")

	// Example of accessing map value via index:
	log.Printf("Request type is ", i.Request.Intent.Name)

	// Create a response object
	resp := CreateResponse()

	// Customize the response for each Alexa Intent
	switch i.Request.Intent.Name {
	case "officetemp":
		resp.Say("The current temperature is 68 degrees.")
	case "hello":
		resp.Say("Hello there, Lambda appears to be working properly.")
	case "AMAZON.HelpIntent":
		resp.Say("This app is easy to use, just say: ask the office how warm it is")
	default:
		resp.Say("I'm sorry, the input does not look like something I understand.")
	}

	return *resp, nil
}

func main() {
	lambda.Start(HandleRequest)
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

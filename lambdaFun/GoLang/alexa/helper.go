package alexa

import (
	"encoding/json"
	"log"
	"strconv"
)

var debugGo = false
var debugUserInfo = false
var debugOptions = false

type CardImg struct {
	Small string
	Large string
}

//<a href="https://www.iconfinder.com/icons/309047/conference_group_people_users_icon" target="_blank">"Conference, group, people, users icon"</a> by <a href="https://www.iconfinder.com/visualpharm" target="_blank">Ivan Boyko</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0" target="_blank">CC BY 3.0</a>
//"Conference, group, people, users icon" (https://www.iconfinder.com/icons/309047/conference_group_people_users_icon) by Ivan Boyko (https://www.iconfinder.com/visualpharm) is licensed under CC BY 3.0 (http://creativecommons.org/licenses/by/3.0)
var conferenceImg CardImg = CardImg{
	Small: "https://s3.amazonaws.com/audio-conference/images/conferenceSmall.png",
	Large: "https://s3.amazonaws.com/audio-conference/images/conferenceLarge.png",
}

//<a href="https://www.iconfinder.com/icons/3324959/outgoing_phone_icon" target="_blank">"Outgoing, phone icon"</a> by <a href="https://www.iconfinder.com/colebemis" target="_blank">Cole Bemis</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0" target="_blank">CC BY 3.0</a>
//"Outgoing, pho/ne icon" (https://www.iconfinder.com/icons/3324959/outgoing_phone_icon) by Cole Bemis (https://www.iconfinder.com/colebemis) is licensed under CC BY 3.0 (http://creativecommons.org/licenses/by/3.0)
var phoneStartImg CardImg = CardImg{
	Small: "https://s3.amazonaws.com/audio-conference/images/phoneStartSmall.png",
	Large: "https://s3.amazonaws.com/audio-conference/images/phoneStartLarge.png",
}

//<a href="https://www.iconfinder.com/icons/3324961/missed_phone_icon" target="_blank">"Missed, phone icon"</a> by <a href="https://www.iconfinder.com/colebemis" target="_blank">Cole Bemis</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0" target="_blank">CC BY 3.0</a>
//"Missed, phone icon" (https://www.iconfinder.com/icons/3324961/missed_phone_icon) by Cole Bemis (https://www.iconfinder.com/colebemis) is licensed under CC BY 3.0 (http://creativecommons.org/licenses/by/3.0)
var phoneStopImg CardImg = CardImg{
	Small: "https://s3.amazonaws.com/audio-conference/images/phoneStopSmall.png",
	Large: "https://s3.amazonaws.com/audio-conference/images/phoneStopLarge.png",
}

//<a href="https://www.iconfinder.com/icons/3324960/off_phone_icon" target="_blank">"Off, phone icon"</a> by <a href="https://www.iconfinder.com/colebemis" target="_blank">Cole Bemis</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0" target="_blank">CC BY 3.0</a>
//"Off, phone icon" (https://www.iconfinder.com/icons/3324960/off_phone_icon) by Cole Bemis (https://www.iconfinder.com/colebemis) is licensed under CC BY 3.0 (http://creativecommons.org/licenses/by/3.0)
var phoneErrorImg CardImg = CardImg{
	Small: "https://s3.amazonaws.com/audio-conference/images/phoneErrorSmall.png",
	Large: "https://s3.amazonaws.com/audio-conference/images/phoneErrorLarge.png",
}

//<a href="https://www.iconfinder.com/icons/183285/help_mark_question_icon" target="_blank">"Help, mark, question icon"</a> by <a href="https://www.iconfinder.com/yanlu" target="_blank">Yannick Lung</a>
//"Help, mark, question icon" (https://www.iconfinder.com/icons/183285/help_mark_question_icon) by Yannick Lung (https://www.iconfinder.com/yanlu)
var questionImg CardImg = CardImg{
	Small: "https://s3.amazonaws.com/audio-conference/images/questionSmall.png",
	Large: "https://s3.amazonaws.com/audio-conference/images/questionLarge.png",
}

func GetDebugGo() bool {
	return debugGo
}

func SetDebugGo(dg bool) {
	debugGo = dg
}

func GetDebugUserInfo() bool {
	return debugUserInfo
}

func SetDebugUserInfo(dui bool) {
	debugUserInfo = dui
}

func GetDebugOptions() bool {
	return debugOptions
}

func SetDebugOptions(do bool) {
	debugOptions = do
}

func GetConferenceImg() CardImg {
	return conferenceImg
}

func GetPhoneStartImg() CardImg {
	return phoneStartImg
}

func GetPhoneStopImg() CardImg {
	return phoneStopImg
}

func GetPhoneErrorImg() CardImg {
	return phoneErrorImg
}

func GetQuestionImg() CardImg {
	return questionImg
}

func LogObject(identifier string, obj interface{}) {

	o, err := json.Marshal(obj)
	if err != nil {
		log.Print("\r" + identifier + ":\r" + "ERROR: could not convert object to JSON")
	} else {
		log.Print("\r" + identifier + ":\r" + string(o))
	}

}

func VerifyPN(PN string) bool {

	runes := []rune(PN)

	//Checking to see if the phone number is 10 digits and does not start with a 0
	if len(PN) == 10 && string(runes[0:1]) != "0" {

		//Checking to see if there are any non-numeric characters in the phone number
		if _, err := strconv.Atoi(PN); err == nil {
			return true
		}

	}

	return false

}

func FormatPN(PN string) string {

	var PNFormatted string

	if VerifyPN(PN) {

		runes := []rune(PN)

		var areaCode = string(runes[0:3])
		var middle = string(runes[3:6])
		var end = string(runes[6])

		PNFormatted = "(" + areaCode + ") " + middle + "-" + end

	} else {
		PNFormatted = "PHONE NUMBER FORMATTING ERROR"
	}

	return PNFormatted

}

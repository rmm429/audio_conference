'use strict';

var userInfo = {};

//<a href="https://www.iconfinder.com/icons/309047/conference_group_people_users_icon" target="_blank">"Conference, group, people, users icon"</a> by <a href="https://www.iconfinder.com/visualpharm" target="_blank">Ivan Boyko</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0" target="_blank">CC BY 3.0</a>
//"Conference, group, people, users icon" (https://www.iconfinder.com/icons/309047/conference_group_people_users_icon) by Ivan Boyko (https://www.iconfinder.com/visualpharm) is licensed under CC BY 3.0 (http://creativecommons.org/licenses/by/3.0)
var conferenceImg = {
    small: "https://s3.amazonaws.com/audio-conference/images/conferenceSmall.png",
    large: "https://s3.amazonaws.com/audio-conference/images/conferenceLarge.png"
};

//<a href="https://www.iconfinder.com/icons/3324959/outgoing_phone_icon" target="_blank">"Outgoing, phone icon"</a> by <a href="https://www.iconfinder.com/colebemis" target="_blank">Cole Bemis</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0" target="_blank">CC BY 3.0</a>
//"Outgoing, phone icon" (https://www.iconfinder.com/icons/3324959/outgoing_phone_icon) by Cole Bemis (https://www.iconfinder.com/colebemis) is licensed under CC BY 3.0 (http://creativecommons.org/licenses/by/3.0)
var phoneStartImg = {
    small: "https://s3.amazonaws.com/audio-conference/images/phoneStartSmall.png",
    large: "https://s3.amazonaws.com/audio-conference/images/phoneStartLarge.png"
};

//<a href="https://www.iconfinder.com/icons/3324961/missed_phone_icon" target="_blank">"Missed, phone icon"</a> by <a href="https://www.iconfinder.com/colebemis" target="_blank">Cole Bemis</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0" target="_blank">CC BY 3.0</a>
//"Missed, phone icon" (https://www.iconfinder.com/icons/3324961/missed_phone_icon) by Cole Bemis (https://www.iconfinder.com/colebemis) is licensed under CC BY 3.0 (http://creativecommons.org/licenses/by/3.0)
var phoneStopImg = {
    small: "https://s3.amazonaws.com/audio-conference/images/phoneStopSmall.png",
    large: "https://s3.amazonaws.com/audio-conference/images/phoneStopLarge.png"
};

//<a href="https://www.iconfinder.com/icons/3324960/off_phone_icon" target="_blank">"Off, phone icon"</a> by <a href="https://www.iconfinder.com/colebemis" target="_blank">Cole Bemis</a> is licensed under <a href="http://creativecommons.org/licenses/by/3.0" target="_blank">CC BY 3.0</a>
//"Off, phone icon" (https://www.iconfinder.com/icons/3324960/off_phone_icon) by Cole Bemis (https://www.iconfinder.com/colebemis) is licensed under CC BY 3.0 (http://creativecommons.org/licenses/by/3.0)
var phoneErrorImg = {
    small: "https://s3.amazonaws.com/audio-conference/images/phoneErrorSmall.png",
    large: "https://s3.amazonaws.com/audio-conference/images/phoneErrorLarge.png"
};

//<a href="https://www.iconfinder.com/icons/183285/help_mark_question_icon" target="_blank">"Help, mark, question icon"</a> by <a href="https://www.iconfinder.com/yanlu" target="_blank">Yannick Lung</a>
//"Help, mark, question icon" (https://www.iconfinder.com/icons/183285/help_mark_question_icon) by Yannick Lung (https://www.iconfinder.com/yanlu)
var questionImg = {
    small: "https://s3.amazonaws.com/audio-conference/images/questionSmall.png",
    large: "https://s3.amazonaws.com/audio-conference/images/questionLarge.png"
};

//event = input JSON
exports.handler = function(event,context) {

    try {

        //Outputting the input JSON to the console
		if(process.env.NODE_DEBUG_EN) {
			console.log("Request:\n" + JSON.stringify(event,null,2));
        }
        
        //specific objects of the event JSON
        var request = event.request;

        //Retrieving the Amazon ID of the current user of the skill
        var amazonId = event.context.System.user.userId;

        //Checking to see if previous data does not exists for the current user of the skill
        if(!(userInfo.hasOwnProperty(amazonId.valueOf()))) {

            //Creating an info JSON to store specific intent information
            var info = {
                startIntent: false,
                PN: "",
                BeAnywhere: ""
            };

            //Associating specific intent information to the current user of the skill
            userInfo[amazonId.valueOf()] = info;

        }
        
        if (request.type === "LaunchRequest") {

			handleLaunchRequest(context);

		} else if (request.type === "IntentRequest") {

            if (request.intent.name === "StartIntent"){

                handleStartIntent(request,context,amazonId);

            } else if (request.intent.name === "StopIntent") {
                
                handleStopIntent(request,context,amazonId);

            } else {
                throw("Unknown intent");
            }

        } else if (request.type === "SessionEndedRequest") {

		} else {
			throw("Unknown intent type");
		}

    } catch(e) {
        context.fail("Exception: " + e);
    }

};

function buildResponse(options) {

    //Outputting the response options to the console
	if(process.env.NODE_DEBUG_EN) {
		console.log("\nbuildResponse options:\n" + JSON.stringify(options,null,2));
    }
    
    options.speechText = addSpacing(options.speechText);

    //response = output JSON
    var response = {
        version: "1.0",
        response: {
            outputSpeech: {
                type: "SSML",
                ssml: "<speak>" + options.speechText + "</speak>"
            },
            shouldEndSession: options.endSession
        }
    };

    if (options.repromptText) {

        options.repromptText = addSpacing(options.repromptText);

        response.response.reprompt = {
            outputSpeech: {
                type: "SSML",
                ssml: "<speak>" + options.repromptText + "</speak>"
            }
        };
    }

    if (options.cardTitle) {
        response.response.card = {
            type: "Simple",
            title: options.cardTitle
        };

        if (options.imageUrl) {
            response.response.card.type = "Standard";
            response.response.card.text = options.cardContent;
            response.response.card.image = {
                smallImageUrl: options.imageUrl,
                largeImageUrl: options.imageUrl
            };
        } else if (options.imageObj) {
            response.response.card.type = "Standard";
            response.response.card.text = options.cardContent;
            response.response.card.image = {
                smallImageUrl: options.imageObj.small,
                largeImageUrl: options.imageObj.large
            };
        } else {
            response.response.card.content = options.cardContent;
        }
    }

    //Outputting the output JSON to the console
    if(process.env.NODE_DEBUG_EN) {
		console.log("\nResponse:\n" + JSON.stringify(response,null,2));
    }
    
    return response;

}

function handleLaunchRequest(context) {
    let options = {};
    options.speechText = "Welcome to the Audio Conference skill.Using this skill, you can start an audio conference on a telephone number or a Be Anywhere device.You can say, for example, ask audio conference to start a conference on <say-as interpret-as=\"telephone\">2155551234</say-as>, or, ask audio conference to start a conference on My Cell.";
    options.cardContent = "Welcome to the Audio Conference skill.  Using this skill, you can start an audio conference on a telephone number or a Be Anywhere device.  You can say, for example, ask audio conference to start a conference on 2155551234, or, ask audio conference to start a conference on My Cell.";
    options.cardTitle = "Audio Conference";
    options.imageObj = conferenceImg;
    options.endSession = true;

    //Outputting the userInfo JSON to the console
    if(process.env.USERINFO_DEBUG_EN) {
		console.log("\nuserInfo (LaunchRequest):\n" + JSON.stringify(userInfo,null,2));
    }

    //Outputting the Launch JSON to the console
    if(process.env.NODE_DEBUG_EN) {
		console.log("\nLaunch:\n" + JSON.stringify(options,null,2));
    }

    context.succeed(buildResponse(options));
}

function handleStartIntent(request,context,amazonId) {
    let options = {};

    //The info JSON that will store specific intent information
    var info = {};

    //Checking to see if previous data exists for the current user of the skill
    if(userInfo.hasOwnProperty(amazonId.valueOf())) {

        //Retrieving the info JSON mapped to the Amazon ID of the user of the skill
        info = userInfo[amazonId.valueOf()];

        //Checking to see if slots exist
        if (request.intent.slots.PN.value || request.intent.slots.BeAnywhere.value)
        {
            //Noting that we are coming from the intent StartIntent
            info.startIntent = true;

            //Checking to see which type of request is being made
            if (request.intent.slots.PN.value) {
                let PN = request.intent.slots.PN.value;
                options.speechText = `Your conference was started on <say-as interpret-as="telephone">${PN}</say-as>.`;
                options.cardContent = `Your conference was started on ${PN}.`;
                //Saving the phone number in the info JSON
                info.PN = PN;
            } else if (request.intent.slots.BeAnywhere.value) {
                let BeAnywhere = request.intent.slots.BeAnywhere.value;
                options.speechText = `Your conference was started on ${BeAnywhere}.`;
                options.cardContent = `Your conference was started on ${BeAnywhere}.`;
                //Saving the Be Anywhere device in the info JSON
                info.BeAnywhere = BeAnywhere;
            }

            options.endSession = true;
            options.imageObj = phoneStartImg;

        } else {

            options.speechText = "What device would you like to start your conference on?";
            options.repromptText = "You can say a telephone number, such as <say-as interpret-as=\"telephone\">2155551234</say-as>, or say a Be Anywhere device, such as My Cell.";
            options.cardContent = "You can say a telephone number, such as 2155551234, or say a Be Anywhere device, such as My Cell."
            options.imageObj = questionImg;
            options.endSession = false;

        }

        options.cardTitle = "Audio Conference Start";

    //If no previous data exists for the current user, throw an error and refresh the data
    } else {

        options.speechText = "Unknown usage error.The data associated with your Amazon ID for this skill was not able to be retrieved.Your data has been refreshed.Please start over and try again!";
        options.cardContent = "Unknown usage error.  The data associated with your Amazon ID for this skill was not able to be retrieved.  Your data has been refreshed.  Please start over and try again!";
        options.imageObj = phoneErrorImg;
        options.cardTitle = "ERROR: Audio Conference Start";

        //Creating an info JSON to store specific intent information
        info = {
            startIntent: false,
            PN: "",
            BeAnywhere: ""
        };

    }

    //Adding the info JSON to the userInfo JSON with the key set as the Amazon ID of the current user of the skill
    userInfo[amazonId.valueOf()] = info;

    //Outputting the userInfo JSON to the console
    if(process.env.USERINFO_DEBUG_EN) {
		console.log("\nuserInfo (StartIntent):\n" + JSON.stringify(userInfo,null,2));
    }

    //Outputting the StartIntent JSON to the console
    if(process.env.NODE_DEBUG_EN) {
		console.log("\nStartIntent:\n" + JSON.stringify(options,null,2));
    }

    context.succeed(buildResponse(options));
}

function handleStopIntent(request,context,amazonId) {
    let options = {};

    //The info JSON that will store specific intent information
    var info = {};

    //Checking to see if previous data exists for the current user of the skill
    if(userInfo.hasOwnProperty(amazonId.valueOf())) {

        //Retrieving the info JSON mapped to the Amazon ID of the user of the skill
        info = userInfo[amazonId.valueOf()];

        //Making sure we came from a the intent StartIntent
        if(info.startIntent) {
            //If a phone number was provided
            if (request.intent.slots.PN.value) {
                //stop conference on the phone number that was provided
                //stopConference(phone_number)
                options.speechText = `Your conference was stopped on <say-as interpret-as="telephone">${request.intent.slots.PN.value}</say-as>.`;
                options.cardContent = `Your conference was stopped on ${request.intent.slots.PN.value}.`;
                options.imageObj = phoneStopImg;
                options.cardTitle = "Audio Conference Stop";
            //If a Be Anywhere device was provided
            } else if (request.intent.slots.BeAnywhere.value) {
                //stop conference on the BeAnywhere device that was provided
                //stopConference(be_anywhere)
                options.speechText = `Your conference was stopped on ${request.intent.slots.BeAnywhere.value}.`;
                options.cardContent = `Your conference was stopped on ${request.intent.slots.BeAnywhere.value}.`;
                options.imageObj = phoneStopImg;
                options.cardTitle = "Audio Conference Stop";
            //If a phone number was used to start the conference
            } else if (info.PN != "") {
                //stop conference on the phone number that the conference was started on
                //stopConference(phone_number)
                options.speechText = `Your conference was stopped on <say-as interpret-as="telephone">${info.PN}</say-as>.`;
                options.cardContent = `Your conference was stopped on ${info.PN}.`;
                options.imageObj = phoneStopImg;
                options.cardTitle = "Audio Conference Stop";

            //If a Be Anywhere device was used to start the conference
            } else if (info.BeAnywhere != "") {
                //stop conference on the BeAnywhere device that the conference was started on
                //stopConference(be_anywhere)
                options.speechText = `Your conference was stopped on ${info.BeAnywhere}.`;
                options.cardContent = `Your conference was stopped on ${info.BeAnywhere}.`;
                options.imageObj = phoneStopImg;
                options.cardTitle = "Audio Conference Stop";

            //No cases were fulfilled, throw error
            } else {
                options.speechText = "Invalid option.To stop the conference, please provide a valid telephone number or BeAnywhere device.";
                options.cardContent = "Invalid option.  To stop the conference, please provide a valid telephone number or BeAnywhere device.";
                options.imageObj = phoneErrorImg;
                options.cardTitle = "ERROR: Audio Conference Stop";
            }

        //If we did not come from the intent StartIntent
        } else {
            options.speechText = "Incorrect usage.To stop a conference, a conference must first be started.";
            options.cardContent = "Incorrect usage.  To stop a conference, a conference must first be started.";
            options.imageObj = phoneErrorImg;
            options.cardTitle = "ERROR: Audio Conference Stop";
        }

    //If no previous data exists for the current user, throw an error and refresh the data
    } else {

        options.speechText = "Unknown usage error.The data associated with your Amazon ID for this skill was not able to be retrieved.Your data has been refreshed.Please start over and try again!";
        options.cardContent = "Unknown usage error.  The data associated with your Amazon ID for this skill was not able to be retrieved.  Your data has been refreshed.  Please start over and try again!";
        options.imageObj = phoneErrorImg;
        options.cardTitle = "ERROR: Audio Conference Stop";

    }

    options.endSession = true;

    //Emptying the info JSON for new usage of the skill
    info = {
        startIntent: false,
        PN: "",
        BeAnywhere: ""
    };

    //Adding the info JSON to the userInfo JSON with the key set as the Amazon ID of the current user of the skill
    userInfo[amazonId.valueOf()] = info;

    //Outputting the userInfo JSON to the console
    if(process.env.USERINFO_DEBUG_EN) {
		console.log("\nuserInfo (StopIntent):\n" + JSON.stringify(userInfo,null,2));
    }

    //Outputting the StopIntent JSON to the console
    if(process.env.NODE_DEBUG_EN) {
		console.log("\nStopIntent:\n" + JSON.stringify(options,null,2));
    }

    context.succeed(buildResponse(options));

}

function addSpacing(text) {

	/*
		Two spaces after period, question mark, and/or exclamation point, one space after final period, question mark, and/or exclamation point.
		Incoming string has no spaces after any period, question mark, and/or exclamation point (with the exception of a justifiable use of a period, e.g. "Mr. Jones").
	*/

	//Find period with no space after it, replace with period with two spaces after it
	var textSpace = text.replace(/\.(?=[^ ])/g, ".  ");
	//Find question mark with no space after it, replace with question mark with two spaces after it
	textSpace = textSpace.replace(/\?(?=[^ ])/g, "?  ");
	//Find exclamation point with no space after it, replace with exclamation point with two spaces after it
	textSpace = textSpace.replace(/\!(?=[^ ])/g, "!  ");

	//Add a final space onto the end of the string and return the string
	return textSpace + " ";

}
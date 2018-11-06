'use strict';

//event = input JSON
exports.handler = function(event,context) {

    try {

        //Outputting the input JSON to the console
		if(process.env.NODE_DEBUG_EN) {
			console.log("Request:\n" + JSON.stringify(event,null,2));
        }
        
        //specific objects of the event JSON
		var request = event.request;
        var session = event.session;
        
        //Session attributes: pass some information from one intent to another intent
		//If there are no session attributes, create an empty object for them and place it in the event JSON
		if(!event.session.attributes) {
			event.session.attributes = {};
        }
        
        if (request.type === "LaunchRequest") {

			handleLaunchRequest(context);

		} else if (request.type === "IntentRequest") {

            if (request.intent.name === "PNIntent") {

                handleStartPNIntent(request,context);

            } else if (request.intent.name === "BeAnywhereIntent") {

                handleStartBeAnywhereIntent(request,context);

            } else if (request.intent.name === "StopIntent") {
                
                handleStopIntent(request,context);

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

}

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
        } else {
            response.response.card.content = options.cardContent;
        }
    }

    if (options.session && options.session.attributes) {
        response.sessionAttributes = options.session.attributes;
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
    options.endSession = true;
    context.succeed(buildResponse(options));
}

function handleStartPNIntent(request,context) {
    let options = {};
    let PN = request.intent.slots.PN.value;
    options.speechText = `Your conference was started on <say-as interpret-as="telephone">${PN}</say-as>.`;
    options.endSession = true;
    context.succeed(buildResponse(options));
}

function handleStartBeAnywhereIntent(request,context) {
    let options = {};
    let BeAnywhere = request.intent.slots.BeAnywhere.value;
    options.speechText = `Your conference was started on ${BeAnywhere}.`;
    options.endSession = true;
    context.succeed(buildResponse(options));
}

function handleStopIntent(request,context) {
    let options = {};
    options.speechText = "Your conference was stopped.";
    options.endSession = true;
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
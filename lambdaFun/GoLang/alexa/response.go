package alexa

func NewSimpleResponse(title string, text string) Response {
	r := Response{
		Version: "1.0",
		Body: ResBody{
			OutputSpeech: &Payload{
				Type: "PlainText",
				Text: text,
			},
			Card: &Payload{
				Type:    "Simple",
				Title:   title,
				Content: text,
			},
			ShouldEndSession: true,
		},
	}
	return r
}

type Response struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Body              ResBody                `json:"response"`
}

type ResBody struct {
	OutputSpeech     *Payload     `json:"outputSpeech,omitempty"`
	Card             *Payload     `json:"card,omitempty"`
	Reprompt         *Reprompt    `json:"reprompt,omitempty"`
	Directives       []Directives `json:"directives,omitempty"`
	ShouldEndSession bool         `json:"shouldEndSession"`
}

type Reprompt struct {
	OutputSpeech Payload `json:"outputSpeech,omitempty"`
}

type Directives struct {
	Type          string         `json:"type,omitempty"`
	SlotToElicit  string         `json:"slotToElicit,omitempty"`
	UpdatedIntent *UpdatedIntent `json:"UpdatedIntent,omitempty"`
	PlayBehavior  string         `json:"playBehavior,omitempty"`
	AudioItems    struct {
		Stream struct {
			Token                string `json:"token,omitempty"`
			URL                  string `json:"url,omitempty"`
			OffsetInMilliseconds int    `json:"offsetInMilliseconds,omitempty"`
		} `json:"stream,omitempty"`
	} `json:"audioItem,omitempty"`
}

type UpdatedIntent struct {
	Name               string                 `json:"name,omitempty"`
	ConfirmationStatus string                 `json:"confirmationStatus,omitempty"`
	Slots              map[string]interface{} `json:"slots,omitempty"`
}

type Image struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

type Payload struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Text    string `json:"text,omitempty"`
	SSML    string `json:"ssml,omitempty"`
	Content string `json:"content,omitempty"`
	Image   Image  `json:"image,omitempty"`
}

func BuildResponse(options map[string]interface{}) Response {

	var response = Response{
		Version: "1.0",
		Body: ResBody{
			OutputSpeech: &Payload{
				Type: "SSML",
				SSML: "<speak>" + options["speechText"].(string) + "</speak>",
			},
			ShouldEndSession: options["endSession"].(bool),
		},
	}

	//If does not work, change repromptText to "var" and use options["repromptText"] instead of "repromptText"
	if repromptText, ok := options["repromptText"]; ok {

		response.Body.Reprompt = &Reprompt{
			OutputSpeech: Payload{
				Type: "SSML",
				SSML: "<speak>" + repromptText.(string) + "</speak>",
			},
		}

	}

	//If does not work, change cardTitle to "var" and use options["cardTitle"] instead of "cardTitle"
	if cardTitle, ok := options["cardTitle"]; ok {

		response.Body.Card = &Payload{
			Type:  "Simple",
			Title: cardTitle.(string),
		}

		//If does not work, change imageUrl to "var" and use options["imageUrl"] instead of "imageUrl"
		if imageUrl, ok := options["imageUrl"]; ok {

			response.Body.Card.Type = "Standard"
			response.Body.Card.Text = options["cardContent"].(string)
			response.Body.Card.Image = Image{
				SmallImageURL: imageUrl.(string),
				LargeImageURL: imageUrl.(string),
			}

			//If does not work, change imageObj to "var" and use options["imageObj"] instead of "imageObj"
		} else if imageObj, ok := options["imageObj"]; ok {

			response.Body.Card.Type = "Standard"
			response.Body.Card.Text = options["cardContent"].(string)
			response.Body.Card.Image = Image{
				SmallImageURL: (imageObj.(CardImg)).Small,
				LargeImageURL: (imageObj.(CardImg)).Large,
			}

		} else {
			response.Body.Card.Content = options["cardContent"].(string)
		}

	}

	return response

}

func NewSSMLResponse(title string, text string, reprompt string, endSession bool, session Session) Response {

	var r Response

	if reprompt == "" {

		if session.Attributes != nil {

			r = Response{
				Version:           "1.0",
				SessionAttributes: session.Attributes,
				Body: ResBody{
					OutputSpeech: &Payload{
						Type: "SSML",
						SSML: text,
					},
					ShouldEndSession: endSession,
				},
			}

		} else {

			r = Response{
				Version: "1.0",
				Body: ResBody{
					OutputSpeech: &Payload{
						Type: "SSML",
						SSML: text,
					},
					ShouldEndSession: endSession,
				},
			}

		}

	} else {

		if session.Attributes != nil {

			r = Response{
				Version:           "1.0",
				SessionAttributes: session.Attributes,
				Body: ResBody{
					OutputSpeech: &Payload{
						Type: "SSML",
						SSML: text,
					},
					Reprompt: &Reprompt{
						OutputSpeech: Payload{
							Type: "SSML",
							SSML: reprompt,
						},
					},
					ShouldEndSession: endSession,
				},
			}

		} else {

			r = Response{
				Version: "1.0",
				Body: ResBody{
					OutputSpeech: &Payload{
						Type: "SSML",
						SSML: text,
					},
					Reprompt: &Reprompt{
						OutputSpeech: Payload{
							Type: "SSML",
							SSML: reprompt,
						},
					},
					ShouldEndSession: endSession,
				},
			}

		}
	}

	if GetDebugGo() {
		LogObject("Response", r)
	}

	return r

}

type SSML struct {
	text  string
	pause string
	pn    string
}

type SSMLBuilder struct {
	SSML []SSML
}

func (builder *SSMLBuilder) Say(text string) {
	builder.SSML = append(builder.SSML, SSML{text: text})
}

func (builder *SSMLBuilder) Pause(pause string) {
	builder.SSML = append(builder.SSML, SSML{pause: pause})
}

func (builder *SSMLBuilder) PN(pn string) {
	builder.SSML = append(builder.SSML, SSML{pn: pn})
}

func (builder *SSMLBuilder) Build() string {
	var response string
	for index, ssml := range builder.SSML {
		if ssml.text != "" {
			response += ssml.text + " "
		} else if ssml.pause != "" && index != len(builder.SSML)-1 {
			response += "<break time='" + ssml.pause + "ms'/>"
		} else if ssml.pn != "" {
			response += `<say-as interpret-as="telephone">` + ssml.pn + "</say-as>"
		}
	}
	return "<speak>" + response + "</speak>"
}

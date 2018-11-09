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

	if GetDebugOptions() {
		LogObject("Options", options)
	}

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

	if session, ok := options["session"]; ok {

		if (session.(Session)).Attributes != nil {
			response.SessionAttributes = (session.(Session)).Attributes
		}

	}

	if repromptText, ok := options["repromptText"]; ok {

		response.Body.Reprompt = &Reprompt{
			OutputSpeech: Payload{
				Type: "SSML",
				SSML: "<speak>" + repromptText.(string) + "</speak>",
			},
		}

	}

	if cardTitle, ok := options["cardTitle"]; ok {

		response.Body.Card = &Payload{
			Type:  "Simple",
			Title: cardTitle.(string),
		}

		if imageUrl, ok := options["imageUrl"]; ok {

			response.Body.Card.Type = "Standard"
			response.Body.Card.Text = options["cardContent"].(string)
			response.Body.Card.Image = Image{
				SmallImageURL: imageUrl.(string),
				LargeImageURL: imageUrl.(string),
			}

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

	if GetDebugGo() {
		LogObject("Response", response)
	}

	return response

}

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

func NewSSMLResponse(title string, text string, reprompt string, endSession bool) Response {

	var r Response

	if reprompt == "" {

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

	LogObject("Response", r)

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

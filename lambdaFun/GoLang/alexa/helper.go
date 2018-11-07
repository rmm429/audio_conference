package alexa

import (
	"encoding/json"
	"log"
)

func LogObject(identifier string, obj interface{}) {

	o, err := json.Marshal(obj)
	if err != nil {
		log.Print("\r" + identifier + ":\r" + "ERROR: could not convert object to JSON")
	} else {
		log.Print("\r" + identifier + ":\r" + string(o))
	}

}

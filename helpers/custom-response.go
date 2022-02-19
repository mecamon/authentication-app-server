package helpers

import (
	"encoding/json"
	"log"
)

func CustomResponse(data interface{}, errResponse ErrorsMap) (hasError bool, out []byte) {

	if len(errResponse.Message) > 0 {

		out, err := json.MarshalIndent(errResponse, "", "    ")
		if err != nil {
			log.Fatal("Error marshalling errResponse to json")
		}

		return true, out
	}

	out, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal("Error marshalling data to json")
	}

	return false, out
}

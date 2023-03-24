package requests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CheckResponse(res *http.Response) (map[string]interface{}, error) {
	if res.StatusCode >= 400 && res.StatusCode < 600 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Read response error")
			return nil, err
		}

		var jsonObj map[string]interface{}
		err = json.Unmarshal([]byte(body), &jsonObj)
		if err != nil {
			fmt.Println("Error while unmarshalling json data:", err)
			return nil, err
		}

		error_data := make(map[string]interface{})
		if errMsg, ok := jsonObj["errorMessage"]; ok {
			error_data["error"] = errMsg
		} else if errMsg, ok := jsonObj["error_description"]; ok {
			error_data["error"] = errMsg
		} else if errMsg, ok := jsonObj["error"]; ok {
			error_data["error"] = errMsg
		}

		return error_data, nil
	} else {
		if res.StatusCode == 201 {
			return nil, nil
		}

		if res.StatusCode == 204 {
			return nil, nil
		}

		fmt.Println("Response status code:", res.StatusCode)

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var jsonObj map[string]interface{}
		var jsonArr []map[string]interface{}
		// check response body type is array of object or object
		if body[0] == '[' {
			err = json.Unmarshal([]byte(body), &jsonArr)
			if err != nil {
				return nil, err
			}
			return jsonArr[0], nil
		}
		err = json.Unmarshal([]byte(body), &jsonObj)
		if err != nil {
			return nil, err
		}
		return jsonObj, nil
	}
}

package controllers

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/schema"
	"neft.web/errorController"
	"neft.web/views"
)

func ParseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		errorController.HandleError(err.Error(), "ParseForm function")
		return err
	}
	return parseValues(r.PostForm, dst)
}
func parseURLParams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.Form, dst)
}
func parseValues(values url.Values, dst interface{}) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(dst, values); err != nil {
		return err
	}
	return nil
}

// Process the form received in the http request and return the data into array of bytes
func processFormToAPI(r *http.Request) ([]byte, error) {
	var middle LoginForm

	if err := ParseForm(r, &middle); err != nil {
		errorController.HandleError(err.Error(), "Parsing processFormToAPI function")
		return nil, err
	}
	endDest, err := json.Marshal(middle)
	if err != nil {
		errorController.ErrorLogger.Println("Error parsing json from processFormToAPI")
		return nil, err
	}
	return endDest, nil
}

type Response struct {
	ErrorField  string `json:"error"`
	Description string `json:"description"`
}

// Process the body of the response, and return it into struct
func readAPIAnswer(resp *http.Response) (Response, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorController.ErrorLogger.Println("Error reading answer from answer request")
		return Response{}, err
	}
	var answer Response
	err = json.Unmarshal(bodyBytes, &answer)
	if err != nil {
		errorController.ErrorLogger.Println(string(bodyBytes))
		errorController.ErrorLogger.Println("Error reading answer from answer request: ", err)
		return Response{}, err
	}
	return answer, nil
}

// ReadInput read in every moment the console and change maintenance and debug mode
func ReadInput() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		switch input.Text() {
		case "reload":
			views.ReloadHtml()
		}
	}
}

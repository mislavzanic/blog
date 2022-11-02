package webhook

import (
	"io/ioutil"
	"log"
	"net/http"
)


func Webhook(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	log.Println(string(body))
}

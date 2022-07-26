package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

type response struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func ForwardHandler(c *gin.Context) {

	c.Request.URL.Host = utils.Env.FwdAddr
	c.Request.URL.Scheme = "http"

	fwdReq, err := http.NewRequest(c.Request.Method, c.Request.URL.String(), c.Request.Body)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	fwdReq.Header = c.Request.Header

	httpForwarder := &http.Client{}

	fwdResp, err := httpForwarder.Do(fwdReq)

	if err != nil {
		msg := fmt.Sprintf("Error in %s", fwdReq.Method)
		resp := response{
			Error:   "Main instance is down",
			Message: msg,
		}
		c.JSON(http.StatusServiceUnavailable, resp)
		return
	}
	if fwdResp != nil {
		defer fwdResp.Body.Close()
		body, _ := ioutil.ReadAll(fwdResp.Body)
		rawJson := json.RawMessage(body)
		c.JSON(fwdResp.StatusCode, rawJson)
	}
}

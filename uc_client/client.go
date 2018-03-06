package uc_client

import (
	"strings"
	"io/ioutil"
	"net/url"
	"log"
	"net/http"
	"encoding/json"
	"errors"
)

type UcResponse struct {
	StatusCode int
	Body       []byte
}

func RequestUcApi(tourl, method string, args map[string]string) (ucres UcResponse, err error) {
	v := url.Values{}
	for key, value := range args {
		if key != "atoken"{
			v.Set(key, value)
		}
	}
	senddata := v.Encode()
	//log.Printf("RequestUcApi: senddata %v", senddata)
	body := strings.NewReader(senddata)
	client := &http.Client{}
	req, err := http.NewRequest(method, tourl, body); if err != nil {
		log.Printf("RequestUcApi: NewRequest error: %v", err)
		return
	}
	token, ok := args["atoken"]; if ok {
		req.Header.Set("x-token", token)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	response, err := client.Do(req); if err != nil {
		log.Printf("RequestUcApi: Client.Do error: %v", err);
		return
	}

	data, err := ioutil.ReadAll(response.Body); if err != nil {
		return
	}
	defer response.Body.Close()
	//res := string(data)
	ucres.StatusCode = response.StatusCode
	ucres.Body = data
	//log.Printf("RequestUcApi: res data:%v", res)
	return
}

func GetResErrInfo(response UcResponse) (err error)  {
	if response.StatusCode != http.StatusOK {
		var msg string
		err = json.Unmarshal(response.Body, &msg); if err != nil {
			return
		}
		log.Printf("Response.ErrInfo: %v", msg)
		return errors.New(msg)
	}
	return nil
}


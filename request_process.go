package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type dictdata struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type webhookdata struct {
	Event            string              `json:"event"`
	Event_Type       string              `json:"event_type"`
	App_Id           string              `json:"app_id"`
	User_Id          string              `json:"user_id"`
	Message_Id       string              `json:"message_id"`
	Page_Title       string              `json:"page_title"`
	Page_Url         string              `json:"page_url"`
	Browser_Language string              `json:"browser_language"`
	Screen_Size      string              `json:"screen_size"`
	Attributes       map[string]dictdata `json:"attributes"`
	UserAttributes   map[string]dictdata `json:"traits"`
}

func getDynamicAttributes(data map[string]string, key_prefix string) map[string]dictdata {

	dynAttr := make(map[string]dictdata)
	regexp_to_search := "^" + key_prefix + "k[0-9]{1,2}$"
	for key, value := range data {

		is_match, _ := regexp.MatchString(regexp_to_search, key)
		if !is_match {
			fmt.Printf("regexp is not matched for key:%s\n", key)
			continue
		}
		digit_postfix := key[len(key_prefix)+1:]

		type_key := key_prefix + "t" + digit_postfix
		value_key := key_prefix + "v" + digit_postfix

		var temp_value = dictdata{}

		type_key_value, ok := data[type_key]
		if !ok {
			fmt.Printf("type is missing for key %s:%s", key, type_key)
			continue
		}

		value_key_value, ok := data[value_key]
		if !ok {
			fmt.Printf("value is missing for key %s:%s", key, value_key)
			continue
		}

		temp_value.Type = type_key_value
		temp_value.Value = value_key_value
		dynAttr[value] = temp_value
	}
	return dynAttr
}

func post_webhook(body []byte) {
	posturl := "https://webhook.site/c9218b70-6bdb-4a4e-b3c5-3eb6c4a2b857"

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error Occured", err)

	}
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println("Error Occured", err)
	}
	fmt.Printf("Response Status: %s\n", res.Status)
	fmt.Printf("Response Body: %s\n", res.Body)
	defer res.Body.Close()
}

func worker(data map[string]string) {
	var postdata webhookdata
	postdata.Event = data["ev"]
	postdata.Event_Type = data["et"]
	postdata.App_Id = data["id"]
	postdata.User_Id = data["uid"]
	postdata.Message_Id = data["mid"]
	postdata.Page_Title = data["t"]
	postdata.Page_Url = data["p"]
	postdata.Browser_Language = data["l"]
	postdata.Screen_Size = data["sc"]
	postdata.Attributes = getDynamicAttributes(data, "atr")
	postdata.UserAttributes = getDynamicAttributes(data, "uatr")
	fmt.Printf("Data to be post %v", postdata)

	post_bytes, err := json.Marshal(postdata)
	if err != nil {
		fmt.Println("Error in Marshalling: ", err)
	}
	fmt.Println("Data to be post in bytest", post_bytes)
	post_webhook(post_bytes)

}

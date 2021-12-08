package sms

import (
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func SendSMS(phone, message string) error {
	err := godotenv.Load(".env")

	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", "https://sms.ru/sms/send", nil)

	if err != nil {
		return err
	}

	query := req.URL.Query()
	query.Add("api_id", os.Getenv("SMS_API_ID"))
	query.Add("to", phone)
	query.Add("msg", message)
	query.Add("json", "1")
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)

	var values map[string]interface{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&values)

	if err != nil {
		return err
	}

	result := values["sms"].(map[string]interface{})[phone].(map[string]interface{})
	status := result["status"].(string)

	if status == "OK" {
		return nil
	}

	if status == "ERROR" {
		return errors.New(result["status_text"].(string))
	}

	return errors.New("error sending SMS")
}

package openai

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type BtpKey struct {
	Uaa struct {
		Apiurl          string `json:"apiurl"`
		Clientid        string `json:"clientid"`
		Clientsecret    string `json:"clientsecret"`
		CredentialType  string `json:"credential-type"`
		Identityzone    string `json:"identityzone"`
		Identityzoneid  string `json:"identityzoneid"`
		Sburl           string `json:"sburl"`
		Subaccountid    string `json:"subaccountid"`
		Tenantid        string `json:"tenantid"`
		Tenantmode      string `json:"tenantmode"`
		Uaadomain       string `json:"uaadomain"`
		Url             string `json:"url"`
		Verificationkey string `json:"verificationkey"`
		Xsappname       string `json:"xsappname"`
		Zoneid          string `json:"zoneid"`
	} `json:"uaa"`
	Url    string `json:"url"`
	Vendor string `json:"vendor"`
}

func LoadBtpKey(f string) (*BtpKey, error) {
	buf, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	out := &BtpKey{}
	err = json.Unmarshal(buf, out)
	return out, err
}

func (key *BtpKey) GetToken() (string, error) {

	type Token struct {
		AccessToken string `json:"access_token"`
	}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	token := Token{}
	req, err := http.NewRequest("POST", key.Uaa.Url+"/oauth/token", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(key.Uaa.Clientid, key.Uaa.Clientsecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", err
	}
	return token.AccessToken, err

}

package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"

	portainer "github.com/portainer/portainer/api"
)

// Service represents a service used to authenticate users against an authorization server
type Service struct{}

// NewService returns a pointer to a new instance of this service
func NewService() *Service {
	return &Service{}
}

// Authenticate takes an access code and exchanges it for an access token from portainer OAuthSettings token environment(endpoint).
// On success, it will then return the username and token expiry time associated to authenticated user by fetching this information
// from the resource server and matching it with the user identifier setting.
func (*Service) Authenticate(code string, configuration *portainer.OAuthSettings) (string, error) {
	token, err := getOAuthToken(code, configuration)
	if err != nil {
		log.Printf("[DEBUG] - Failed retrieving access token: %v", err)
		return "", err
	}
	username, err := getUsername(token.AccessToken, configuration)
	if err != nil {
		log.Printf("[DEBUG] - Failed retrieving oauth user name: %v", err)
		return "", err
	}
	err = checkGroup(token.AccessToken, configuration)
	if err != nil {
		log.Printf("[DEBUG] - User is not part of group: %v , error: %v", configuration.Scopes, err)
		return "", err
	}

	return username, nil
}

func getOAuthToken(code string, configuration *portainer.OAuthSettings) (*oauth2.Token, error) {
	unescapedCode, err := url.QueryUnescape(code)
	if err != nil {
		return nil, err
	}

	config := buildConfig(configuration)
	token, err := config.Exchange(context.Background(), unescapedCode)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func getUsername(token string, configuration *portainer.OAuthSettings) (string, error) {
	req, err := http.NewRequest("GET", configuration.ResourceURI, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", &oauth2.RetrieveError{
			Response: resp,
			Body:     body,
		}
	}

	content, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}

	if content == "application/x-www-form-urlencoded" || content == "text/plain" {
		values, err := url.ParseQuery(string(body))
		if err != nil {
			return "", err
		}

		username := values.Get(configuration.UserIdentifier)
		if username == "" {
			return username, &oauth2.RetrieveError{
				Response: resp,
				Body:     body,
			}
		}

		return username, nil
	}

	var datamap map[string]interface{}
	if err = json.Unmarshal(body, &datamap); err != nil {
		return "", err
	}

	username, ok := datamap[configuration.UserIdentifier].(string)
	if ok && username != "" {
		return username, nil
	}

	if !ok {
		username, ok := datamap[configuration.UserIdentifier].(float64)
		if ok && username != 0 {
			return fmt.Sprint(int(username)), nil
		}
	}

	return "", &oauth2.RetrieveError{
		Response: resp,
		Body:     body,
	}
}

func checkGroup(token string, configuration *portainer.OAuthSettings) (error) {
	// Hack: get tenant specific checkMemberGroups API url for Azure AD without modifying database
	// https://graph.windows.net/${tenantID}/me?api-version=2013-11-08 -> https://graph.windows.net/${tenantID}/me/checkMemberGroups?api-version=2013-11-08
	requestUrl := strings.Replace(configuration.ResourceURI, "me?api-version=2013-11-08", "me/checkMemberGroups?api-version=2013-11-08", -1)

	// Request group ID defined on configuration.Scopes
	requestBody := `{"groupIds": ["` + configuration.Scopes + `"]}`

	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(requestBody))
	if err != nil {
		return err
	}

	client := &http.Client{}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return &oauth2.RetrieveError{
			Response: resp,
			Body:     body,
		}
	}

	var datamap map[string]interface{}
	if err = json.Unmarshal(body, &datamap); err != nil {
		return err
	}

	groupIDs, ok := datamap["value"]
	groupString := fmt.Sprint(groupIDs)
	if ok && groupString == "["+configuration.Scopes+"]" {
		return nil
	}
	log.Printf("[DEBUG] - oauth checkGroup group ID: %v does not match to [%v]", groupString, configuration.Scopes)

	return &oauth2.RetrieveError{
		Response: resp,
		Body:     body,
	}
}

func buildConfig(configuration *portainer.OAuthSettings) *oauth2.Config {
	endpoint := oauth2.Endpoint{
		AuthURL:  configuration.AuthorizationURI,
		TokenURL: configuration.AccessTokenURI,
	}

	return &oauth2.Config{
		ClientID:     configuration.ClientID,
		ClientSecret: configuration.ClientSecret,
		Endpoint:     endpoint,
		RedirectURL:  configuration.RedirectURI,
		Scopes:       []string{"id","email","name"},
	}
}

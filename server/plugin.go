package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/create":
		p.handleCreate(w, r)
	default:
		http.NotFound(w, r)
	}
}

type CreateAPIRequest struct {
	PostID string `json:"post_id"`
}

// ProductBoardRequest Is what we send to ProductBoard
type ProductBoardRequest struct {
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	CustomerEmail string   `json:"customer_email"`
	DisplayURL    string   `json:"display_url"`
	Tags          []string `json:"tags"`
}

func (p *Plugin) handleCreate(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	var createRequest *CreateAPIRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&createRequest)
	if err != nil {
		p.API.LogError("Unable to decode JSON err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	config := p.getConfiguration()

	tags := []string{}
	if config.Tags != "" {
		tags = strings.Split(config.Tags, ",")
	}

	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		p.API.LogError("Unable to get user err=" + appErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	serverConfig := p.API.GetConfig()

	docPost, appErr := p.API.GetPost(createRequest.PostID)
	if appErr != nil {
		p.API.LogError("Unable to get post err=" + appErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rootID := docPost.RootId
	if rootID == "" {
		rootID = docPost.Id
	}

	permalink, err := url.Parse(*serverConfig.ServiceSettings.SiteURL)
	permalink.Path = path.Join(permalink.Path, "_redirect", "pl", docPost.Id)

	body := fmt.Sprintf("Mattermost user `%s` from %s has recorded the following feedback:\n\n```\n%s\n```\n\nSee the original post [here](%s).\n\n ",
		user.Username,
		*serverConfig.ServiceSettings.SiteURL,
		docPost.Message,
		permalink.String(),
	)

	lines := strings.Split(docPost.Message, "\n")
	title := lines[0]

	// this is Github stuff - needs to be replaced with a POST to the ProductBoard API (I think?)

	newNote := &ProductBoardRequest{
		Title:         title,
		Content:       body,
		Tags:          tags,
		CustomerEmail: user.Email,
		DisplayURL:    permalink.String(),
	}

	b, err := json.Marshal(newNote)
	if err != nil {
		p.API.LogError("Unable to marshal err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buf := bytes.NewBuffer(b)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.productboard.com/notes", buf)
	if err != nil {
		p.API.LogError("Unable to marshal err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", "Bearer: "+config.ProductBoardAPIKey)
	resp, err := client.Do(req)
	if err != nil {
		p.API.LogError("Unable to marshal err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Need to parse the response to pull out the URL to the new note
	type PBLinks struct {
		HTML string `json:"html"`
	}
	type PBResponse struct {
		Links PBLinks `json:"links"`
	}

	var pb PBResponse

	newPostURL := json.Unmarshal(resp.Request.Response.Body, &pb)

	// Post this to the channel with the URL/ID of the newly created Note

	post := &model.Post{
		UserId:    userID,
		ChannelId: docPost.ChannelId,
		RootId:    rootID,
		Message:   fmt.Sprintf(" [this post](%s) has been submitted to ProductBoard as a note for processing by a PM - you can comment or add info [here](%s).\n\n ", permalink.String(), newPostURL.string()),
	}

	_, appErr = p.API.CreatePost(post)
	if appErr != nil {
		p.API.LogError("Unable to create post err=" + appErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// See https://developers.mattermost.com/extend/plugins/server/reference/

func NewString(s string) *string { return &s }

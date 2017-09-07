package client

import (
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/do4way/evernote-sdk-golang/edam"
	"github.com/mrjones/oauth"
)

//EnvironmentType ...
type EnvironmentType int

const (
	// SANDBOX ...
	SANDBOX EnvironmentType = iota
	// PRODUCTION production environemnt
	PRODUCTION
	// YINXIANG evernote in chinese
	YINXIANG
)

//EvernoteClient client object
type EvernoteClient struct {
	host            string
	oauthClient     *oauth.Consumer
	userStoreClient *edam.UserStoreClient
}

//NewClient return a new evernote client
func NewClient(key, secret string, envType EnvironmentType) *EvernoteClient {
	host := "www.evernote.com"
	if envType == SANDBOX {
		host = "sandbox.evernote.com"
	} else if envType == YINXIANG {
		host = "app.yinxiang.com"
	}
	client := oauth.NewConsumer(
		key, secret,
		oauth.ServiceProvider{
			RequestTokenUrl:   fmt.Sprintf("https://%s/oauth", host),
			AuthorizeTokenUrl: fmt.Sprintf("https://%s/OAuth.action", host),
			AccessTokenUrl:    fmt.Sprintf("https://%s/oauth", host),
		},
	)
	return &EvernoteClient{
		host:        host,
		oauthClient: client,
	}
}

//GetRequestToken get a request token.
func (c *EvernoteClient) GetRequestToken(callBackURL string) (*oauth.RequestToken, string, error) {
	return c.oauthClient.GetRequestTokenAndUrl(callBackURL)
}

//GetAuthorizedToken get an authorized token.
func (c *EvernoteClient) GetAuthorizedToken(requestToken *oauth.RequestToken, oauthVerifier string) (*oauth.AccessToken, error) {
	return c.oauthClient.AuthorizeToken(requestToken, oauthVerifier)
}

//GetUserStore get userStore
func (c *EvernoteClient) GetUserStore() (*edam.UserStoreClient, error) {
	if c.userStoreClient != nil {
		return c.userStoreClient, nil
	}
	evernoteUserStoreServerURL := fmt.Sprintf("https://%s/edam/user", c.host)
	evernoteUserTrans, err := thrift.NewTHttpPostClient(evernoteUserStoreServerURL)
	if err != nil {
		return nil, err
	}
	c.userStoreClient = edam.NewUserStoreClientFactory(
		evernoteUserTrans,
		thrift.NewTBinaryProtocolFactoryDefault(),
	)
	return c.userStoreClient, nil
}

//GetNoteStore get noteStore.
func (c *EvernoteClient) GetNoteStore(authenticationToken string) (*edam.NoteStoreClient, error) {
	us, err := c.GetUserStore()
	if err != nil {
		return nil, err
	}
	urls, err := us.GetUserUrls(authenticationToken)
	if err != nil {
		return nil, err
	}
	url := urls.GetNoteStoreUrl()
	httpClient, err := thrift.NewTHttpPostClient(url)
	if err != nil {
		return nil, err
	}
	return edam.NewNoteStoreClientFactory(
		httpClient,
		thrift.NewTBinaryProtocolFactoryDefault(),
	), nil
}

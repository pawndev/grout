package romm

import (
	"fmt"
	"time"
)

type Collection struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	URLCover    string    `json:"url_cover"`
	HasCover    bool      `json:"has_cover"`
	IsPublic    bool      `json:"is_public"`
	UserID      int       `json:"user_id"`
	ROMs        []Rom     `json:"roms"`
	ROMCount    int       `json:"rom_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type smartCollection struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description *string     `json:"description"`
	URLCover    string      `json:"url_cover"`
	HasCover    bool        `json:"has_cover"`
	IsPublic    bool        `json:"is_public"`
	UserID      int         `json:"user_id"`
	Criteria    interface{} `json:"criteria"`
	ROMs        []Rom       `json:"roms"`
	ROMCount    int         `json:"rom_count"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type virtualCollection struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	URLCover    string    `json:"url_cover"`
	HasCover    bool      `json:"has_cover"`
	Type        string    `json:"type"`
	ROMs        []Rom     `json:"roms"`
	ROMCount    int       `json:"rom_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Client) GetCollections() ([]Collection, error) {
	var collections []Collection
	err := c.doRequest("GET", endpointCollections, nil, nil, &collections)
	return collections, err
}

func (c *Client) getCollection(id int) (*Collection, error) {
	var collection Collection
	path := fmt.Sprintf(endpointCollectionByID, id)
	err := c.doRequest("GET", path, nil, nil, &collection)
	return &collection, err
}

func (c *Client) getSmartCollections() ([]smartCollection, error) {
	var collections []smartCollection
	err := c.doRequest("GET", endpointSmartCollections, nil, nil, &collections)
	return collections, err
}

func (c *Client) getSmartCollection(id int) (*smartCollection, error) {
	var collection smartCollection
	path := fmt.Sprintf(endpointSmartCollectionByID, id)
	err := c.doRequest("GET", path, nil, nil, &collection)
	return &collection, err
}

func (c *Client) getVirtualCollections() ([]virtualCollection, error) {
	var collections []virtualCollection
	err := c.doRequest("GET", endpointVirtualCollections, nil, nil, &collections)
	return collections, err
}

func (c *Client) getVirtualCollection(id int) (*virtualCollection, error) {
	var collection virtualCollection
	path := fmt.Sprintf(endpointVirtualCollectionByID, id)
	err := c.doRequest("GET", path, nil, nil, &collection)
	return &collection, err
}

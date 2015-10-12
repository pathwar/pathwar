package api

type StdItem struct {
	Id            string          `json:"_id",omitempty`
	Created       string          `json:"_created",omitempty`
	Etag          string          `json:"_etag",omitempty`
	Links         map[string]Link `json:"_links",omitempty`
	SchemaVersion int             `json:"_schema_version",omitempty`
	Updated       string          `json:"_updated",omitempty`
}

type StdList struct {
	Links map[string]Link `json:"_links",omitempty`
	Meta  Meta            `json:"_meta",omitempty`
}

type Link struct {
	Href  string `json:"href",omitempty`
	Title string `json:"title",omitempty`
}

type Meta struct {
	MaxResults int `json:"max_results",omitempty`
	Page       int `json:"page",omitempty`
	Total      int `json:"total",omitempty`
}

type RawOrganizationUsers struct {
	StdList

	Items []RawOrganizationUser `json:"_items"`
}

type RawOrganizationUser struct {
	StdItem

	Role         string `json:"role"`
	User         string `json:"user"`
	Organization string `json:"organization"`
}

type Users struct {
	StdList

	Items []User `json:"_items"`
}

type User struct {
	StdItem

	Company       string `json:"company"`
	GithubHandle  string `json:"github_handle"`
	GravatarHash  string `json:"gravatar_hash"`
	Location      string `json:"location"`
	Login         string `json:"login"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	TwitterHandle string `json:"twitter_handle"`
	Website       string `json:"website"`
}

type RawLevelInstanceUsers struct {
	StdList

	Items []RawLevelInstanceUser `json:"_items"`
}

type RawLevelInstanceUser struct {
	StdItem

	ExpiryDate        string `json:"expiry_date"`
	Hash              string `json:"hash"`
	Level             string `json:"level"`
	LevelInstance     string `json:"level_instance"`
	Organization      string `json:"organization"`
	OrganizationLevel string `json:"organization_level"`
	User              string `json:"user"`
}

type RawLevelInstances struct {
	StdList

	Items []RawLevelInstance `json:"_items"`
}

type RawLevelInstance struct {
	StdItem

	Active      bool   `json:"active"`
	Level       string `json:"level"`
	Name        string `json:"name"`
	PwnStatus   string `json:"pwn_status"`
	Passphrases []struct {
		Value string `json:"value"`
		Key   string `json:"key"`
	} `json:"passphrases"`
	PrivateUrls []struct {
		Url  string `json:"url"`
		Name string `json:"name"`
	} `json:"private_urls"`
	Urls []struct {
		Url  string `json:"url"`
		Name string `json:"name"`
	}
}

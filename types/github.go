package types

import "time"

// If interface{} is used, it means the field won't be used anyway

type GithubRelease struct {
	Assets          interface{} `json:"assets"`
	AssetsURL       string      `json:"assets_url"`
	Author          interface{} `json:"author,omitempty"`
	Body            string      `json:"body,omitempty"`
	CreatedAt       time.Time   `json:"created_at,omitempty"`
	DiscussionURL   string      `json:"discussion_url"`
	Draft           bool        `json:"draft"`
	HTMLUrl         string      `json:"html_url"`
	ID              int         `json:"id"`
	Name            string      `json:"name,omitempty"`
	NodeID          string      `json:"node_id"`
	Prerelease      bool        `json:"prerelease"`
	PublishedAt     time.Time   `json:"published_at,omitempty"`
	Reactions       interface{} `json:"reactions"`
	TagName         string      `json:"tag_name"`
	TarballURL      string      `json:"tarball_url,omitempty"`
	TargetCommitish string      `json:"target_commitish"`
	UploadURL       string      `json:"upload_url"`
	URL             string      `json:"url"`
	ZipballURL      string      `json:"zipball_url,omitempty"`
}

type GithubReleaseWebhook struct {
	Action       string        `json:"action"`
	Enterprise   interface{}   `json:"enterprise"`
	Installation interface{}   `json:"installation"`
	Organization interface{}   `json:"organization"`
	Release      GithubRelease `json:"release"`
	Repository   interface{}   `json:"repository"`
	Sender       interface{}   `json:"sender"`
}

package hashicorp

import "time"

type LicenseClass string

const (
	EnterPrise LicenseClass = "enterprise"
	OSS        LicenseClass = "oss"
)

type ListReleasesParam struct {
	Limit int // default 10
	// 	string <date-time>
	// This timestamp is used as a pagination marker, indicating that only releases that occurred prior to it should be retrieved.
	// When fetching subsequent pages, this parameter should be set to the creation timestamp of the oldest release listed on the current page.
	After string
	// If specified, only releases with a matching license class will be returned. For example, if set to enterprise, only enterprise releases would be returned.
	LicenseClass LicenseClass
}

type ListReleasesResponse struct {
	Releases []*Release
}

type SpecificReleaseParam struct {
	// If specified, only releases with a matching license class will be returned. For example, if set to enterprise, only enterprise releases would be returned.
	LicenseClass LicenseClass
}

type Release struct {
	Builds                     []Build   `json:"builds"`
	DockerNameTag              string    `json:"docker_name_tag"`
	IsPrerelease               bool      `json:"is_prerelease"`
	LicenseClass               string    `json:"license_class"`
	Name                       string    `json:"name"`
	Status                     Status    `json:"status"`
	TimestampCreated           time.Time `json:"timestamp_created"`
	TimestampUpdated           time.Time `json:"timestamp_updated"`
	URLBlogpost                string    `json:"url_blogpost"`
	URLChangelog               string    `json:"url_changelog"`
	URLDockerRegistryDockerhub string    `json:"url_docker_registry_dockerhub"`
	URLDockerRegistryEcr       string    `json:"url_docker_registry_ecr"`
	URLLicense                 string    `json:"url_license"`
	URLProjectWebsite          string    `json:"url_project_website"`
	URLReleaseNotes            string    `json:"url_release_notes"`
	URLShasums                 string    `json:"url_shasums"`
	URLShasumsSignatures       []string  `json:"url_shasums_signatures"`
	URLSourceRepository        string    `json:"url_source_repository"`
	Version                    string    `json:"version"`
}

type Status struct {
	Message string `json:"message"`
	State   string `json:"state"`
}

type Build struct {
	Arch        string `json:"arch"`
	Os          string `json:"os"`
	Unsupported bool   `json:"unsupported"`
	URL         string `json:"url"`
}

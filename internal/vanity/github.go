package vanity

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	gh "github.com/cli/go-gh/v2"
)

// Repository represents GitHub repository metadata.
type Repository struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	License     *License   `json:"license"`
	Topics      []string   `json:"topics"`
	Owner       Owner      `json:"owner"`
}

// License represents GitHub repository license info.
type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SPDXID string `json:"spdx_id"`
}

// Owner represents GitHub repository owner.
type Owner struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

// Release represents a GitHub release.
type Release struct {
	TagName string `json:"tag_name"`
}

// Tag represents a GitHub tag.
type Tag struct {
	Name string `json:"name"`
}

// GetRepository fetches repository metadata from GitHub.
func GetRepository(owner, repo string) (*Repository, error) {
	args := []string{"api", fmt.Sprintf("repos/%s/%s", owner, repo)}

	stdout, _, err := gh.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository: %w", err)
	}

	var repository Repository
	if err := json.Unmarshal(stdout.Bytes(), &repository); err != nil {
		return nil, fmt.Errorf("failed to parse repository data: %w", err)
	}

	return &repository, nil
}

// GetLatestVersion returns the latest version (release tag or git tag) for a repository.
func GetLatestVersion(owner, repo string) (string, error) {
	// Try to get latest release first
	args := []string{"api", fmt.Sprintf("repos/%s/%s/releases/latest", owner, repo)}

	stdout, _, err := gh.Exec(args...)
	if err == nil {
		var release Release
		if err := json.Unmarshal(stdout.Bytes(), &release); err == nil && release.TagName != "" {
			return release.TagName, nil
		}
	}

	// If no releases, try tags
	args = []string{"api", fmt.Sprintf("repos/%s/%s/tags", owner, repo)}

	stdout, _, err = gh.Exec(args...)
	if err != nil {
		return "", nil // No version available
	}

	var tags []Tag
	if err := json.Unmarshal(stdout.Bytes(), &tags); err != nil {
		return "", nil
	}

	if len(tags) > 0 {
		return tags[0].Name, nil
	}

	return "", nil
}

// ParseRepoURL extracts owner and repo name from a GitHub URL.
func ParseRepoURL(url string) (owner, repo string, err error) {
	// Parse GitHub URL formats:
	// https://github.com/owner/repo
	// https://github.com/owner/repo.git
	// git@github.com:owner/repo.git

	// Try HTTPS URL with path parsing
	if len(url) > 19 && url[:19] == "https://github.com/" {
		path := url[19:]
		parts := strings.Split(path, "/")
		if len(parts) >= 2 {
			owner = parts[0]
			repo = parts[1]
			// Remove .git suffix if present
			if strings.HasSuffix(repo, ".git") {
				repo = repo[:len(repo)-4]
			}
			return owner, repo, nil
		}
	}

	// Try git SSH URL
	if len(url) > 15 && url[:15] == "git@github.com:" {
		path := url[15:]
		parts := strings.Split(path, "/")
		if len(parts) >= 2 {
			owner = parts[0]
			repo = parts[1]
			// Remove .git suffix if present
			if strings.HasSuffix(repo, ".git") {
				repo = repo[:len(repo)-4]
			}
			return owner, repo, nil
		}
	}

	return "", "", fmt.Errorf("invalid GitHub URL format: %s", url)
}

// CreatePackageFromRepo creates a Package from GitHub repository metadata.
func CreatePackageFromRepo(vanityDomain, pkgName, repoURL string) (*Package, error) {
	owner, repo, err := ParseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	repository, err := GetRepository(owner, repo)
	if err != nil {
		return nil, err
	}

	version, _ := GetLatestVersion(owner, repo)
	if version == "" {
		version = "v0.1.0"
	}

	license := "MIT"
	if repository.License != nil && repository.License.SPDXID != "" {
		license = repository.License.SPDXID
	}

	author := repository.Owner.Login
	if repository.Owner.Name != "" {
		author = repository.Owner.Name
	}

	now := time.Now()
	return &Package{
		Title:            pkgName,
		ImportPath:       fmt.Sprintf("%s/pkg/%s", vanityDomain, pkgName),
		RepoURL:          repoURL,
		Description:      repository.Description,
		Version:          version,
		DocumentationURL: fmt.Sprintf("https://pkg.go.dev/%s/pkg/%s", vanityDomain, pkgName),
		License:          license,
		Author:           author,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

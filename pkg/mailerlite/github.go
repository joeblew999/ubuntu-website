package mailerlite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/tabwriter"
)

// GitHub release constants
const (
	GitHubAPIBase     = "https://api.github.com"
	GitHubReleasesURL = "https://github.com/%s/%s/releases"
	GitHubLatestURL   = "https://github.com/%s/%s/releases/latest"
	GitHubDownloadURL = "https://github.com/%s/%s/releases/download/%s/%s"

	DefaultGitHubOwner = "joeblew999"
	DefaultGitHubRepo  = "ubuntu-website"
)

// ReleaseAsset represents a downloadable asset from a GitHub release.
type ReleaseAsset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

// Release represents a GitHub release.
type Release struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	Body        string         `json:"body"`
	Draft       bool           `json:"draft"`
	Prerelease  bool           `json:"prerelease"`
	CreatedAt   string         `json:"created_at"`
	PublishedAt string         `json:"published_at"`
	Assets      []ReleaseAsset `json:"assets"`
	HTMLURL     string         `json:"html_url"`
}

// GetReleasesURL returns the URL for a repo's releases page.
func GetReleasesURL(owner, repo string) string {
	return fmt.Sprintf(GitHubReleasesURL, owner, repo)
}

// GetLatestReleaseURL returns the URL for the latest release.
func GetLatestReleaseURL(owner, repo string) string {
	return fmt.Sprintf(GitHubLatestURL, owner, repo)
}

// GetDownloadURL returns the direct download URL for a release asset.
func GetDownloadURL(owner, repo, tag, filename string) string {
	return fmt.Sprintf(GitHubDownloadURL, owner, repo, tag, filename)
}

// handleReleases handles the releases command for GitHub release information.
func (c *CLI) handleReleases(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("releases subcommand required: latest, list, urls")
	}

	owner := DefaultGitHubOwner
	repo := DefaultGitHubRepo

	var subCmd string
	for _, arg := range args {
		if strings.HasPrefix(arg, "OWNER=") {
			owner = strings.TrimPrefix(arg, "OWNER=")
		} else if strings.HasPrefix(arg, "REPO=") {
			repo = strings.TrimPrefix(arg, "REPO=")
		} else if subCmd == "" {
			subCmd = arg
		}
	}

	switch subCmd {
	case "latest":
		return c.releasesLatest(owner, repo)
	case "list":
		return c.releasesList(owner, repo)
	case "urls":
		return c.releasesURLs(owner, repo)
	default:
		return fmt.Errorf("unknown releases subcommand: %s (use: latest, list, urls)", subCmd)
	}
}

func (c *CLI) releasesLatest(owner, repo string) error {
	apiURL := fmt.Sprintf("%s/repos/%s/%s/releases/latest", GitHubAPIBase, owner, repo)

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return fmt.Errorf("no releases found for %s/%s", owner, repo)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to parse release: %w", err)
	}

	if c.githubIssue {
		c.printf("## Latest Release: %s\n\n", release.Name)
		c.printf("- **Tag:** %s\n", release.TagName)
		c.printf("- **URL:** %s\n", release.HTMLURL)
		if len(release.Assets) > 0 {
			c.println("\n### Downloads")
			for _, asset := range release.Assets {
				c.printf("- [%s](%s) (%.2f MB)\n", asset.Name, asset.DownloadURL, float64(asset.Size)/1024/1024)
			}
		}
		return nil
	}

	c.printf("Latest Release: %s\n", release.Name)
	c.println(strings.Repeat("=", 50))
	c.printf("Tag:       %s\n", release.TagName)
	c.printf("Published: %s\n", release.PublishedAt)
	c.printf("URL:       %s\n", release.HTMLURL)

	if len(release.Assets) > 0 {
		c.println("\nDownloads:")
		w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "FILENAME\tSIZE\tURL")
		for _, asset := range release.Assets {
			fmt.Fprintf(w, "%s\t%.2f MB\t%s\n", asset.Name, float64(asset.Size)/1024/1024, asset.DownloadURL)
		}
		w.Flush()
	}

	return nil
}

func (c *CLI) releasesList(owner, repo string) error {
	apiURL := fmt.Sprintf("%s/repos/%s/%s/releases", GitHubAPIBase, owner, repo)

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return fmt.Errorf("failed to parse releases: %w", err)
	}

	if len(releases) == 0 {
		c.printf("No releases found for %s/%s\n", owner, repo)
		return nil
	}

	if c.githubIssue {
		c.printf("## Releases for %s/%s\n\n", owner, repo)
		c.println("| Tag | Name | Published | Assets |")
		c.println("|-----|------|-----------|--------|")
		for _, r := range releases {
			c.printf("| %s | %s | %s | %d |\n", r.TagName, r.Name, r.PublishedAt[:10], len(r.Assets))
		}
		return nil
	}

	c.printf("Releases for %s/%s\n", owner, repo)
	c.println(strings.Repeat("=", 50))

	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TAG\tNAME\tPUBLISHED\tASSETS")
	for _, r := range releases {
		published := r.PublishedAt
		if len(published) > 10 {
			published = published[:10]
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", r.TagName, r.Name, published, len(r.Assets))
	}
	w.Flush()

	return nil
}

func (c *CLI) releasesURLs(owner, repo string) error {
	c.printf("GitHub Release URLs for %s/%s\n", owner, repo)
	c.println(strings.Repeat("=", 50))
	c.println()
	c.println("Use these URLs in your MailerLite email templates:")
	c.println()
	c.printf("  Releases Page:   %s\n", GetReleasesURL(owner, repo))
	c.printf("  Latest Release:  %s\n", GetLatestReleaseURL(owner, repo))
	c.println()
	c.println("For direct download links, use:")
	c.printf("  %s\n", fmt.Sprintf(GitHubDownloadURL, owner, repo, "{TAG}", "{FILENAME}"))
	c.println()
	c.println("Example:")
	c.printf("  %s\n", GetDownloadURL(owner, repo, "v1.0.0", "software-darwin-arm64.tar.gz"))
	c.printf("  %s\n", GetDownloadURL(owner, repo, "v1.0.0", "software-windows-amd64.zip"))
	c.printf("  %s\n", GetDownloadURL(owner, repo, "v1.0.0", "software-linux-amd64.tar.gz"))

	return nil
}

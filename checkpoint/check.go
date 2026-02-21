package checkpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	DefaultCacheDuration = 48 * time.Hour
	checkTimeout         = 5 * time.Second
	releaseAPIURL        = "https://api.github.com/repos/ademajagon/gix/releases/latest"
)

// CheckParams configures a version check request
type CheckParams struct {
	Product       string
	Version       string
	SignatureFile string
	CacheFile     string
	CacheDuration time.Duration
}

type CheckResponse struct {
	Product            string `json:"product"`
	CurrentVersion     string `json:"current_version"`
	CurrentDownloadURL string `json:"current_download_url"`
	ProjectWebsite     string `json:"project_website"`
	Outdated           bool   `json:"outdated"`
}

type cacheEntry struct {
	CheckedAt time.Time      `json:"checked_at"`
	Response  *CheckResponse `json:"response"`
}

func Check(p *CheckParams) (*CheckResponse, error) {
	if os.Getenv("GIX_CHECKPOINT_DISABLE") != "" {
		return nil, nil
	}

	v := strings.TrimPrefix(p.Version, "v")
	if v == "dev" || v == "ci" || v == "" {
		return nil, nil
	}

	ttl := p.CacheDuration
	if ttl == 0 {
		ttl = DefaultCacheDuration
	}

	resp, err := fetchLatestRelease(p.Product, p.Version)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

func fetchLatestRelease(product, currentVersion string) (*CheckResponse, error) {
	req, err := http.NewRequest(http.MethodGet, releaseAPIURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", fmt.Sprintf(
		"%s%s (+https://github.com/ademajagon/gix; sig=%s; os=%s; arch=%s)",
		product, currentVersion, runtime.GOOS, runtime.GOARCH,
	))

	client := &http.Client{Timeout: checkTimeout}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API: %s", res.Status)
	}

	var release githubRelease
	if err := json.NewDecoder(res.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	resp := &CheckResponse{
		Product:            product,
		CurrentVersion:     release.TagName,
		CurrentDownloadURL: release.HTMLURL,
		ProjectWebsite:     "https://github.com/ademajagon/gix",
		Outdated:           isNewer(release.TagName, currentVersion),
	}

	return resp, nil
}

func isNewer(latest, current string) bool {
	l := strings.TrimPrefix(latest, "v")
	c := strings.TrimPrefix(current, "v")
	if l == c {
		return false
	}

	lp := parseSemver(l)
	cp := parseSemver(c)
	for i := range lp {
		if lp[i] > cp[i] {
			return lp[i] > cp[i]
		}
	}

	return false
}

func parseSemver(s string) [3]int {
	var major, minor, patch int
	fmt.Sscanf(s, "%d.%d.%d", &major, &minor, &patch)
	return [3]int{major, minor, patch}
}

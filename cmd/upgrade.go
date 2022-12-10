package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/kierdavis/ansi"
	"github.com/spf13/cobra"
	"github.com/things-go/log"
	"github.com/tj/go-update"
	"github.com/tj/go-update/progress"
	"golang.org/x/oauth2"

	"github.com/things-go/ormat/pkg/config"
	"github.com/things-go/ormat/pkg/consts"
)

var upgradeCmd = &cobra.Command{
	Use:          "upgrade",
	Short:        "Upgrade ormat",
	Long:         "Upgrade ormat by providing a version. If no version is provided, upgrade to the latest.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.Global
		err := c.Load()
		if err != nil {
			return err
		}
		setupBase(c)
		ansi.HideCursor()
		defer ansi.ShowCursor()

		m := &manager{
			Manager: &update.Manager{
				Command: "ormat",
				Store: &githubStore{
					Owner:   "things-go",
					Repo:    "ormat",
					Version: consts.Version,
					Client:  NewGithubClient(os.Getenv("GITHUB_TOKEN")),
				},
			},
		}

		var r *update.Release
		var a *update.Asset

		r, err = m.GetNewerReleases(args...)
		if err != nil {
			return err
		}
		if r == nil {
			log.Info("No upgrades")
			return nil
		}
		// find the tarball for this system
		arch := runtime.GOARCH
		if runtime.GOARCH == "amd64" {
			arch = "x86_64"
		}
		if runtime.GOOS == "windows" {
			a = r.FindZip(runtime.GOOS, arch)
		} else {
			a = r.FindTarball(runtime.GOOS, arch)
		}
		if a == nil {
			log.Info("No upgrade for your system")
			return nil
		}
		log.Infof("Downloading release: %v", r.Version)
		tmpPath, err := a.DownloadProxy(progress.Reader)
		if err != nil {
			return fmt.Errorf("Download failed: %s", err)
		}
		log.Infof("Downloaded release to %s", tmpPath)

		// install it
		if err := m.Install(tmpPath); err != nil {
			return fmt.Errorf("install failed, %s", err)
		}
		log.Infof("Upgraded to %s", r.Version)
		return nil
	},
}

type manager struct {
	*update.Manager
}

// GetNewerReleases returns the specified release or latest or releases newer than Version, or nil.
func (m *manager) GetNewerReleases(version ...string) (*update.Release, error) {
	if len(version) > 0 && version[0] != "" {
		return m.GetRelease(version[0])
	}
	// fetch the new releases
	releases, err := m.LatestReleases()
	if err != nil {
		return nil, fmt.Errorf("error fetching releases: %s", err)
	}
	// no updates
	if len(releases) == 0 {
		log.Debug("No upgrades")
		return nil, nil
	}

	// latest release
	return releases[0], nil
}

type githubStore struct {
	Owner   string
	Repo    string
	Version string
	*github.Client
}

func NewGithubClient(accessToken string) *github.Client {
	var tc *http.Client

	if accessToken != "" {
		tc = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: accessToken,
		}))
	}
	return github.NewClient(tc)
}

// GetRelease returns the specified release or ErrNotFound.
func (s *githubStore) GetRelease(version string) (*update.Release, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, res, err := s.Client.Repositories.GetReleaseByTag(ctx, s.Owner, s.Repo, version)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 404 {
		return nil, update.ErrNotFound
	}
	return intoUpdateRelease(r), nil
}

// LatestReleases returns releases newer, or nil.
func (s *githubStore) LatestReleases() ([]*update.Release, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	releases, _, err := s.Client.Repositories.ListReleases(ctx, s.Owner, s.Repo, nil)
	if err != nil {
		return nil, err
	}

	latest := make([]*update.Release, 0, len(releases))
	for _, r := range releases {
		tag := r.GetTagName()

		if tag == s.Version || tag == strings.TrimPrefix(s.Version, "v") {
			break
		}
		latest = append(latest, intoUpdateRelease(r))
	}
	return latest, nil
}

// intoUpdateRelease returns a Release.
func intoUpdateRelease(r *github.RepositoryRelease) *update.Release {
	assets := make([]*update.Asset, 0, len(r.Assets))
	for _, v := range r.Assets {
		assets = append(assets, &update.Asset{
			Name:      v.GetName(),
			Size:      v.GetSize(),
			URL:       v.GetBrowserDownloadURL(),
			Downloads: v.GetDownloadCount(),
		})
	}
	return &update.Release{
		Version:     r.GetTagName(),
		Notes:       r.GetBody(),
		PublishedAt: r.GetPublishedAt().Time,
		URL:         r.GetURL(),
		Assets:      assets,
	}
}

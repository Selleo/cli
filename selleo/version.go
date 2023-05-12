package selleo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gojek/heimdall/httpclient"
)

// Version will be replaced by goreleaser using ldflags.
var Version string = "dev"

type Release struct {
	TagName string `json:"tag_name"`
}

func IsNewVersionAvailable() (available bool, current string, newest string) {
	timeout := 1000 * time.Millisecond
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	res, err := client.Get("https://api.github.com/repos/Selleo/cli/releases/latest", nil)
	if err != nil {
		return false, Version, ""
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, Version, ""
	}
	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		return false, Version, ""
	}
	var (
		major, newMajor int
		minor, newMinor int
		patch, newPatch int
	)
	_, curErr := fmt.Sscanf(Version, "v%d.%d.%d", &major, &minor, &patch)
	_, newErr := fmt.Sscanf(release.TagName, "v%d.%d.%d", &newMajor, &newMinor, &newPatch)
	if curErr != nil || newErr != nil {
		return false, Version, release.TagName
	}

	if newMajor > major {
		return true, Version, release.TagName
	}
	if newMajor < major {
		return false, Version, release.TagName
	}

	if newMinor > minor {
		return true, Version, release.TagName
	}
	if newMinor < minor {
		return false, Version, release.TagName
	}

	if newPatch > patch {
		return true, Version, release.TagName
	}
	if newPatch < patch {
		return false, Version, release.TagName
	}

	return false, Version, release.TagName
}

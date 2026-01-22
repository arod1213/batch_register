package spotify

import (
	"net/url"
	"strings"
)

// ExtractSpotifyID attempts to extract a Spotify entity type + ID from a URL/URI.
//
// Supported inputs:
// - https://open.spotify.com/{playlist|album|artist|track}/{id}?...
// - spotify:{playlist|album|artist|track}:{id}
// - Raw IDs (returns ok=false, caller can treat as already-normalized)
func ExtractSpotifyID(input string) (kind string, id string, ok bool) {
	in := strings.TrimSpace(input)
	if in == "" {
		return "", "", false
	}

	// spotify:playlist:<id>
	if strings.HasPrefix(in, "spotify:") {
		parts := strings.Split(in, ":")
		if len(parts) >= 3 {
			k := strings.TrimSpace(parts[1])
			i := strings.TrimSpace(parts[2])
			if k != "" && i != "" {
				return strings.ToLower(k), i, true
			}
		}
		return "", "", false
	}

	// open.spotify.com URL
	if strings.Contains(in, "open.spotify.com") {
		u, err := url.Parse(in)
		if err != nil {
			return "", "", false
		}
		seg := strings.Split(strings.Trim(u.Path, "/"), "/")
		// Expect: /{kind}/{id}
		if len(seg) >= 2 {
			k := strings.ToLower(seg[0])
			i := seg[1]
			if k != "" && i != "" {
				return k, i, true
			}
		}
		return "", "", false
	}

	return "", "", false
}



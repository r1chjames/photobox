package domain

import (
	"encoding/json"

	"gorm.io/datatypes"
	"time"
)

type Album struct {
	ID string `gorm:"primarykey" json:"id"`
	// WorkspaceID scopes this album to a workspace (issue #74). Album names
	// are unique per workspace, not globally.
	WorkspaceID string `json:"workspaceId" gorm:"index;size:36"`
	// Name is unique per workspace (composite unique idx_albums_ws_name),
	// not globally — two tenants may both have "Holidays" (issue #74).
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Tags         string         `json:"tags"`
	Metadata     datatypes.JSON `json:"metadata"`
	CreatedAt    time.Time      `json:"createdAt"`
	CreatedEpoch int64          `json:"createdEpoch" gorm:"index"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	Thumbnail    string         `json:"thumbnail"`
	CoverPhotoId string         `json:"coverPhotoId"`
}

// SmartAlbumRules holds the filter criteria for a smart album. Mirrors
// PhotoSearchFilters so saved filters can be reused directly.
type SmartAlbumRules struct {
	StartDate   string   `json:"startDate,omitempty"`
	EndDate     string   `json:"endDate,omitempty"`
	MediaType   string   `json:"mediaType,omitempty"`
	Camera      string   `json:"camera,omitempty"`
	HasGPS      bool     `json:"hasGps,omitempty"`
	Orientation string   `json:"orientation,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Favorite    bool     `json:"favorite,omitempty"`
	LowQuality  bool     `json:"lowQuality,omitempty"`
}

// IsSmart reports whether the album is a smart (rule-based) album by
// checking its metadata for the smart marker.
func (a *Album) IsSmart() bool {
	var meta map[string]any
	if len(a.Metadata) == 0 || a.Metadata == nil {
		return false
	}
	if err := json.Unmarshal(a.Metadata, &meta); err != nil {
		return false
	}
	smart, _ := meta["smart"].(bool)
	return smart
}

// SmartRules parses the smart-album rules from metadata. Returns an empty
// rules set when the album is not smart or metadata is malformed.
func (a *Album) SmartRules() SmartAlbumRules {
	var rules SmartAlbumRules
	if len(a.Metadata) == 0 {
		return rules
	}
	var meta map[string]any
	if err := json.Unmarshal(a.Metadata, &meta); err != nil {
		return rules
	}
	raw, ok := meta["rules"]
	if !ok {
		return rules
	}
	b, _ := json.Marshal(raw)
	_ = json.Unmarshal(b, &rules)
	return rules
}

// ToFilters converts smart rules into photo search filters.
func (r SmartAlbumRules) ToFilters() PhotoSearchFilters {
	return PhotoSearchFilters(r)
}

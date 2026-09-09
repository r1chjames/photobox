package domain

import "time"

// FaceDetection is a single detected face within a photo. Embeddings are
// stored encrypted at rest ([]byte ciphertext) and never serialized to JSON.
type FaceDetection struct {
	ID        string    `gorm:"primarykey" json:"id"`
	PhotoID   string    `gorm:"index:idx_face_detections_photo" json:"photoId"`
	PersonID  *string   `gorm:"index:idx_face_detections_person" json:"personId,omitempty"`
	BoxX      float64   `json:"boxX"`
	BoxY      float64   `json:"boxY"`
	BoxW      float64   `json:"boxW"`
	BoxH      float64   `json:"boxH"`
	Score     float64   `json:"score"`
	Embedding []byte    `gorm:"type:bytea" json:"-"`
	// Status tracks the review lifecycle: 'detected' (unassigned),
	// 'assigned' (person set), 'review' (needs user confirmation).
	Status    string    `gorm:"default:'detected';index" json:"status"`
	Quality   float64   `gorm:"default:0" json:"quality"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Photo is populated by service joins for API responses; not a column.
	Photo *Photo `gorm:"-" json:"photo,omitempty"`
	// PersonName is populated by service joins for API responses; not a column.
	PersonName string `gorm:"-" json:"personName,omitempty"`
}

func (FaceDetection) TableName() string { return "face_detections" }

// FacePerson is a named cluster of faces ("a person" in the user's library).
type FacePerson struct {
	ID          string    `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"index" json:"name"`
	CoverFaceID string    `json:"coverFaceId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// PhotoCount is populated by repository queries; not a stored column.
	PhotoCount int `gorm:"->;-:migration" json:"photoCount"`
}

func (FacePerson) TableName() string { return "face_persons" }

// FaceCluster is a cluster/centroid record. PersonID is nil while faces are
// in the unassigned "review" pool. Centroid is the encrypted mean embedding.
type FaceCluster struct {
	ID          string    `gorm:"primarykey" json:"id"`
	PersonID    *string   `gorm:"index" json:"personId,omitempty"`
	Centroid    []byte    `gorm:"type:bytea" json:"-"`
	FaceCount   int       `gorm:"default:0" json:"faceCount"`
	NeedsReview bool      `gorm:"default:false" json:"needsReview"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (FaceCluster) TableName() string { return "face_clusters" }

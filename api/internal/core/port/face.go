package port

import "gitlab.com/r1chjames/photobox/api/internal/core/domain"

//go:generate mockgen -source=face.go -destination=mock/face.go -package=mock

// FaceRepository persists face detections, people, and clusters.
type FaceRepository interface {
	InsertDetectionsForPhoto(photoID string, dets []domain.FaceDetection) error
	DeleteDetectionsForPhoto(photoID string) error
	DeleteAllFaceData() (int64, error)
	FacesForPhoto(photoID string) ([]domain.FaceDetection, error)
	FacesForPerson(personID string, limit int) ([]domain.FaceDetection, error)
	PersonPhotoIDs(personID string) ([]string, error)
	Persons() ([]domain.FacePerson, error)
	PersonByID(id string) (*domain.FacePerson, error)
	CreatePerson(name string) (*domain.FacePerson, error)
	UpdatePerson(id, name string) error
	DeletePerson(id string) error
	UnassignedDetections(limit int) ([]domain.FaceDetection, error)
	AllDetections(limit int) ([]domain.FaceDetection, error)
	SetDetectionPerson(detID string, personID *string) error
	AssignDetectionsToPerson(detIDs []string, personID string) error
	FaceCount() (int64, error)
	SetCoverFace(personID, faceID string) error
	UpsertCluster(c domain.FaceCluster) error
	DeleteCluster(id string) error
	Clusters() ([]domain.FaceCluster, error)
	// FaceSetting reads the workspace opt-in flag (false when unset).
	FaceSetting() (bool, error)
	// SetFaceSetting writes the workspace opt-in flag.
	SetFaceSetting(enabled bool) error
}

// FaceEngineResult is one detected face from the face engine.
// Box is [x, y, width, height] in original image pixel coordinates.
type FaceEngineResult struct {
	Box       [4]float64
	Score     float64
	Embedding [128]float32
}

// FaceEngine is the local face detection + embedding service client.
type FaceEngine interface {
	// Detect runs face detection + embedding on the image at imagePath and
	// returns normalized results, or an error when the engine is unreachable.
	Detect(imagePath string) ([]FaceEngineResult, error)
}

// FaceService orchestrates face detection, person management, and clustering.
type FaceService interface {
	// Enabled reports whether face recognition is enabled (workspace setting).
	Enabled() (bool, error)
	// SetEnabled flips the workspace opt-in setting.
	SetEnabled(enabled bool) error
	// Status returns enabled state plus face/person counts for the UI.
	Status() (map[string]any, error)

	// ProcessPhoto detects faces in a photo and stores/replaces detections.
	ProcessPhoto(photoID, imagePath string) (int, error)
	// ProcessLibraryBatch processes up to batch photos that have no stored
	// detections yet, returning the number of photos processed.
	ProcessLibraryBatch(batch int) (int, error)
	// ReclusterAll groups unassigned faces into people (background job).
	ReclusterAll() error

	// Persons lists all people with photo counts.
	Persons() ([]domain.FacePerson, error)
	// PersonByID returns one person.
	PersonByID(id string) (*domain.FacePerson, error)
	// FacesForPhoto lists detections for a photo (enriched with person name).
	FacesForPhoto(photoID string) ([]domain.FaceDetection, error)
	// FacesForPerson lists detections assigned to a person.
	FacesForPerson(personID string, limit int) ([]domain.FaceDetection, error)
	// PersonPhotos returns the photos a person appears in.
	PersonPhotos(personID string) ([]domain.Photo, error)

	// RenamePerson renames a person.
	RenamePerson(id, name string) (*domain.FacePerson, error)
	// MergePersons folds source persons (and their faces) into the destination.
	MergePersons(destinationID string, sourceIDs []string) error
	// SplitPerson moves the given face IDs out of a person into a new person.
	SplitPerson(personID string, faceIDs []string) (*domain.FacePerson, error)
	// DeletePerson removes a person, returning faces to the unassigned pool.
	DeletePerson(id string) error

	// DeleteAll wipes every face row (privacy action).
	DeleteAll() (int64, error)
}

package service

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
)

// faceRecognitionSetting is the workspace opt-in setting key. Face data is
// only collected while this setting is enabled (privacy consent gate).
const faceRecognitionSetting = "face_recognition_enabled"

// matchThreshold is the cosine-similarity floor for assigning a face to an
// existing person cluster. Below this a face stays in the "review" pool.
const matchThreshold = 0.45

// FaceService orchestrates face detection, person management, and clustering.
type FaceService struct {
	faceRepo port.FaceRepository
	photoRepo port.PhotoRepository
	engine   port.FaceEngine
	config   appconfig.AppConfig
	crypto   *EmbeddingCipher
	// enabledOverride lets tests force the setting without a DB row.
	enabledOverride *bool
}

// NewFaceService constructs the face service. engine may be nil when the face
// engine is disabled by config; all operations then report the feature off.
func NewFaceService(faceRepo port.FaceRepository, photoRepo port.PhotoRepository, engine port.FaceEngine, config appconfig.AppConfig, crypto *EmbeddingCipher) *FaceService {
	return &FaceService{faceRepo: faceRepo, photoRepo: photoRepo, engine: engine, config: config, crypto: crypto}
}

// EnableForTest forces the enabled state for unit tests (no DB dependency).
func (fs *FaceService) EnableForTest(enabled bool) {
	fs.enabledOverride = &enabled
}

// Enabled reports the workspace opt-in setting (default disabled).
func (fs *FaceService) Enabled() (bool, error) {
	if fs.enabledOverride != nil {
		return *fs.enabledOverride, nil
	}
	if !fs.config.FaceEngineEnabled || fs.engine == nil {
		return false, nil
	}
	setting, err := fs.faceRepo.FaceSetting()
	if err != nil {
		// Missing setting row = disabled; do not error the whole request.
		slog.Warn("Failed to read face recognition setting", "error", err)
		return false, nil
	}
	return setting, nil
}

// SetEnabled flips the workspace opt-in setting.
func (fs *FaceService) SetEnabled(enabled bool) error {
	if enabled && (!fs.config.FaceEngineEnabled || fs.engine == nil) {
		return fmt.Errorf("%w: face engine not configured", domain.ErrFaceEngineUnavailable)
	}
	return fs.faceRepo.SetFaceSetting(enabled)
}

// Status returns {enabled, faceCount, personCount} for the UI header.
func (fs *FaceService) Status() (map[string]any, error) {
	enabled, err := fs.Enabled()
	if err != nil {
		return nil, err
	}
	out := map[string]any{"enabled": enabled}
	if !enabled {
		return out, nil
	}
	faces, err := fs.faceRepo.FaceCount()
	if err == nil {
		out["faceCount"] = faces
	}
	persons, err := fs.faceRepo.Persons()
	if err == nil {
		out["personCount"] = len(persons)
	}
	return out, nil
}

// ProcessPhoto detects faces in one photo and stores them. Returns face count.
func (fs *FaceService) ProcessPhoto(photoID, imagePath string) (int, error) {
	enabled, err := fs.Enabled()
	if err != nil {
		return 0, err
	}
	if !enabled {
		return 0, domain.ErrFaceRecognitionDisabled
	}
	results, err := fs.engine.Detect(imagePath)
	if err != nil {
		return 0, err
	}
	dets := make([]domain.FaceDetection, 0, len(results))
	for _, r := range results {
		enc, err := fs.encryptEmbedding(r.Embedding)
		if err != nil {
			return 0, err
		}
		dets = append(dets, domain.FaceDetection{
			ID:        uuid.NewString(),
			PhotoID:   photoID,
			BoxX:      r.Box[0],
			BoxY:      r.Box[1],
			BoxW:      r.Box[2],
			BoxH:      r.Box[3],
			Score:     r.Score,
			Embedding: enc,
			Status:    "detected",
			Quality:   r.Score,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	if err := fs.faceRepo.InsertDetectionsForPhoto(photoID, dets); err != nil {
		return 0, err
	}
	// Assign confident matches to existing persons immediately.
	fs.assignNewDetections(dets)
	return len(dets), nil
}

// ProcessLibraryBatch detects faces in up to batch photos that have no stored
// detections. Returns the number of photos processed.
func (fs *FaceService) ProcessLibraryBatch(batch int) (int, error) {
	enabled, err := fs.Enabled()
	if err != nil {
		return 0, err
	}
	if !enabled {
		return 0, domain.ErrFaceRecognitionDisabled
	}
	if batch <= 0 {
		batch = 50
	}
	photos, err := fs.photoRepo.ListPhotosPendingFaceDetection(batch)
	if err != nil {
		return 0, err
	}
	processed := 0
	for _, p := range photos {
		if p.FilesystemPath == "" {
			continue
		}
		path := utils.UnescapeInvalidCharacters(p.FilesystemPath)
		n, err := fs.ProcessPhoto(p.ID, path)
		if err != nil {
			slog.Warn("Face detection failed for photo", "id", p.ID, "error", err)
			continue
		}
		_ = n
		processed++
	}
	return processed, nil
}

// ReclusterAll performs a full grouping pass over unassigned faces and
// refreshes per-person centroids.
func (fs *FaceService) ReclusterAll() error {
	enabled, err := fs.Enabled()
	if err != nil {
		return err
	}
	if !enabled {
		return domain.ErrFaceRecognitionDisabled
	}

	// Group unassigned faces by embedding similarity.
	unassigned, err := fs.faceRepo.UnassignedDetections(0)
	if err != nil {
		return err
	}
	created, err := fs.groupUnassigned(unassigned)
	if err != nil {
		return err
	}

	// Refresh centroids for every person (and delete stale cluster rows).
	persons, err := fs.faceRepo.Persons()
	if err != nil {
		return err
	}
	for _, person := range persons {
		if err := fs.refreshCentroid(person.ID); err != nil {
			slog.Warn("Failed to refresh centroid", "person", person.ID, "error", err)
		}
	}
	if err := fs.cleanupOrphanClusters(); err != nil {
		return err
	}
	slog.Info("Face recluster complete", "newPeople", created, "people", len(persons))
	return nil
}

// groupUnassigned greedily agglomerates unassigned faces into people when
// their pairwise similarity exceeds matchThreshold. Returns new person count.
// embedItem pairs a detection with its decrypted embedding for clustering.
type embedItem struct {
	det domain.FaceDetection
	emb []float32
}

func (fs *FaceService) groupUnassigned(dets []domain.FaceDetection) (int, error) {
	type item = embedItem
	items := make([]item, 0, len(dets))
	for _, d := range dets {
		emb, err := fs.decryptEmbedding(d.Embedding)
		if err != nil {
			slog.Warn("Skipping undecryptable embedding", "face", d.ID, "error", err)
			continue
		}
		items = append(items, item{det: d, emb: emb})
	}
	if len(items) < 1 {
		return 0, nil
	}
	// Sort best-first so the clearest faces seed clusters.
	sort.Slice(items, func(i, j int) bool { return items[i].det.Score > items[j].det.Score })

	created := 0
	var assigned []string
	for i := range items {
		if contains(assigned, items[i].det.ID) {
			continue
		}
		// Seed a new cluster with the highest-scoring remaining face.
		cluster := []item{items[i]}
		assigned = append(assigned, items[i].det.ID)
		centroid := items[i].emb
		// Greedily attach any remaining face close to the running centroid.
		for changed := true; changed; {
			changed = false
			for j := range items {
				if contains(assigned, items[j].det.ID) {
					continue
				}
				if cosine(centroid, items[j].emb) >= matchThreshold {
					cluster = append(cluster, items[j])
					assigned = append(assigned, items[j].det.ID)
					centroid = meanEmbedding(cluster)
					changed = true
				}
			}
		}
		if len(cluster) < 1 {
			continue
		}
		person, err := fs.faceRepo.CreatePerson(fmt.Sprintf("Person %d", time.Now().UnixNano()%100000))
		if err != nil {
			return created, err
		}
		ids := make([]string, len(cluster))
		for k, c := range cluster {
			ids[k] = c.det.ID
		}
		if err := fs.faceRepo.AssignDetectionsToPerson(ids, person.ID); err != nil {
			return created, err
		}
		// Best-scoring face is the cover.
		if err := fs.faceRepo.SetCoverFace(person.ID, cluster[0].det.ID); err != nil {
			slog.Warn("Failed to set cover face", "error", err)
		}
		created++
	}
	return created, nil
}

// assignNewDetections matches fresh detections to existing people; leftovers
// remain in the unassigned review pool.
func (fs *FaceService) assignNewDetections(dets []domain.FaceDetection) {
	persons, err := fs.faceRepo.Persons()
	if err != nil || len(persons) == 0 {
		return
	}
	// Build one centroid per person.
	type pc struct {
		id   string
		cent []float32
	}
	pool := make([]pc, 0, len(persons))
	for _, p := range persons {
		cent, err := fs.centroidFor(p.ID)
		if err != nil || cent == nil {
			continue
		}
		pool = append(pool, pc{id: p.ID, cent: cent})
	}
	if len(pool) == 0 {
		return
	}
	for _, d := range dets {
		if d.PersonID != nil {
			continue
		}
		emb, err := fs.decryptEmbedding(d.Embedding)
		if err != nil {
			continue
		}
		best, bestSim := "", -1.0
		for _, p := range pool {
			if s := cosine(p.cent, emb); s > bestSim {
				bestSim, best = s, p.id
			}
		}
		if best != "" && bestSim >= matchThreshold {
			if err := fs.faceRepo.SetDetectionPerson(d.ID, &best); err != nil {
				slog.Warn("Failed to assign detection", "face", d.ID, "error", err)
			}
		}
	}
}

// RenamePerson renames a person.
func (fs *FaceService) RenamePerson(id, name string) (*domain.FacePerson, error) {
	if _, err := fs.faceRepo.PersonByID(id); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidRequest)
	}
	if err := fs.faceRepo.UpdatePerson(id, name); err != nil {
		return nil, err
	}
	return fs.faceRepo.PersonByID(id)
}

// MergePersons moves every face of the source people onto the destination and
// deletes the source person rows.
func (fs *FaceService) MergePersons(destinationID string, sourceIDs []string) error {
	if _, err := fs.faceRepo.PersonByID(destinationID); err != nil {
		return err
	}
	for _, src := range sourceIDs {
		if src == destinationID {
			return fmt.Errorf("%w: cannot merge a person into itself", domain.ErrInvalidRequest)
		}
		faces, err := fs.faceRepo.FacesForPerson(src, 0)
		if err != nil {
			return err
		}
		if len(faces) == 0 {
			// Empty source — just remove the person row.
			if err := fs.faceRepo.DeletePerson(src); err != nil {
				return err
			}
			continue
		}
		ids := make([]string, len(faces))
		for i, f := range faces {
			ids[i] = f.ID
		}
		if err := fs.faceRepo.AssignDetectionsToPerson(ids, destinationID); err != nil {
			return err
		}
		if err := fs.faceRepo.DeletePerson(src); err != nil {
			return err
		}
	}
	return fs.refreshCentroid(destinationID)
}

// SplitPerson moves the given faces out of a person into a brand-new person.
func (fs *FaceService) SplitPerson(personID string, faceIDs []string) (*domain.FacePerson, error) {
	if len(faceIDs) == 0 {
		return nil, fmt.Errorf("%w: no faces selected", domain.ErrInvalidRequest)
	}
	if _, err := fs.faceRepo.PersonByID(personID); err != nil {
		return nil, err
	}
	newPerson, err := fs.faceRepo.CreatePerson(fmt.Sprintf("Person %d", time.Now().UnixNano()%100000))
	if err != nil {
		return nil, err
	}
	if err := fs.faceRepo.AssignDetectionsToPerson(faceIDs, newPerson.ID); err != nil {
		return nil, err
	}
	if err := fs.refreshCentroid(personID); err != nil {
		slog.Warn("Failed to refresh source centroid after split", "error", err)
	}
	if err := fs.refreshCentroid(newPerson.ID); err != nil {
		slog.Warn("Failed to refresh new centroid after split", "error", err)
	}
	return newPerson, nil
}

// DeletePerson removes a person, returning its faces to the unassigned pool.
func (fs *FaceService) DeletePerson(id string) error {
	return fs.faceRepo.DeletePerson(id)
}

// Persons lists all people.
func (fs *FaceService) Persons() ([]domain.FacePerson, error) {
	return fs.faceRepo.Persons()
}

// PersonByID returns one person.
func (fs *FaceService) PersonByID(id string) (*domain.FacePerson, error) {
	return fs.faceRepo.PersonByID(id)
}

// FacesForPhoto lists a photo's detections enriched with person names.
func (fs *FaceService) FacesForPhoto(photoID string) ([]domain.FaceDetection, error) {
	dets, err := fs.faceRepo.FacesForPhoto(photoID)
	if err != nil {
		return nil, err
	}
	return fs.enrichNames(dets)
}

// FacesForPerson lists a person's detections.
func (fs *FaceService) FacesForPerson(personID string, limit int) ([]domain.FaceDetection, error) {
	return fs.faceRepo.FacesForPerson(personID, limit)
}

// PersonPhotos returns the (non-deleted) photos a person appears in.
func (fs *FaceService) PersonPhotos(personID string) ([]domain.Photo, error) {
	ids, err := fs.faceRepo.PersonPhotoIDs(personID)
	if err != nil {
		return nil, err
	}
	photos := make([]domain.Photo, 0, len(ids))
	for _, id := range ids {
		p, err := fs.photoRepo.GetPhotoById(id, true)
		if err != nil {
			continue // photo deleted since detection — skip
		}
		photos = append(photos, *p)
	}
	return photos, nil
}

// DeleteAll wipes every face row (privacy "delete all face data").
func (fs *FaceService) DeleteAll() (int64, error) {
	return fs.faceRepo.DeleteAllFaceData()
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func (fs *FaceService) encryptEmbedding(e [128]float32) ([]byte, error) {
	if fs.crypto == nil {
		return nil, errors.New("embedding cipher not configured")
	}
	buf := make([]byte, 512)
	for i, v := range e {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(v))
	}
	return fs.crypto.Encrypt(buf)
}

func (fs *FaceService) decryptEmbedding(data []byte) ([]float32, error) {
	if fs.crypto == nil {
		return nil, errors.New("embedding cipher not configured")
	}
	plain, err := fs.crypto.Decrypt(data)
	if err != nil {
		return nil, err
	}
	if len(plain) != 512 {
		return nil, fmt.Errorf("unexpected embedding length %d", len(plain))
	}
	out := make([]float32, 128)
	for i := range out {
		out[i] = math.Float32frombits(binary.LittleEndian.Uint32(plain[i*4:]))
	}
	return out, nil
}

// centroidFor returns the mean embedding of all faces assigned to a person.
func (fs *FaceService) centroidFor(personID string) ([]float32, error) {
	faces, err := fs.faceRepo.FacesForPerson(personID, 0)
	if err != nil || len(faces) == 0 {
		return nil, err
	}
	embs := make([][]float32, 0, len(faces))
	for _, f := range faces {
		e, err := fs.decryptEmbedding(f.Embedding)
		if err != nil {
			continue
		}
		embs = append(embs, e)
	}
	if len(embs) == 0 {
		return nil, nil
	}
	return meanEmbeddings(embs), nil
}

// refreshCentroid recomputes a person's centroid cluster row.
func (fs *FaceService) refreshCentroid(personID string) error {
	cent, err := fs.centroidFor(personID)
	if err != nil {
		return err
	}
	if cent == nil {
		return nil
	}
	enc, err := fs.encryptCentroid(cent)
	if err != nil {
		return err
	}
	cluster := domain.FaceCluster{
		ID:       personID, // 1:1 person -> cluster row, id = person id
		PersonID: &personID,
		Centroid: enc,
	}
	return fs.faceRepo.UpsertCluster(cluster)
}

func (fs *FaceService) encryptCentroid(cent []float32) ([]byte, error) {
	if fs.crypto == nil {
		return nil, errors.New("embedding cipher not configured")
	}
	buf := make([]byte, len(cent)*4)
	for i, v := range cent {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(v))
	}
	return fs.crypto.Encrypt(buf)
}

// cleanupOrphanClusters drops cluster rows whose person no longer exists
// (defensive; MergePersons/SplitPerson keep rows consistent).
func (fs *FaceService) cleanupOrphanClusters() error {
	clusters, err := fs.faceRepo.Clusters()
	if err != nil {
		return err
	}
	for _, c := range clusters {
		if c.PersonID == nil {
			continue
		}
		if _, err := fs.faceRepo.PersonByID(*c.PersonID); err != nil {
			_ = fs.faceRepo.DeleteCluster(c.ID)
		}
	}
	return nil
}

// enrichNames attaches person names to a set of detections.
func (fs *FaceService) enrichNames(dets []domain.FaceDetection) ([]domain.FaceDetection, error) {
	names := map[string]string{}
	persons, err := fs.faceRepo.Persons()
	if err == nil {
		for _, p := range persons {
			names[p.ID] = p.Name
		}
	}
	for i := range dets {
		if dets[i].PersonID != nil {
			dets[i].PersonName = names[*dets[i].PersonID]
		}
	}
	return dets, nil
}

func cosine(a, b []float32) float64 {
	var dot, na, nb float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		na += float64(a[i]) * float64(a[i])
		nb += float64(b[i]) * float64(b[i])
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

func meanEmbedding(items []embedItem) []float32 {
	embs := make([][]float32, len(items))
	for i, it := range items {
		embs[i] = it.emb
	}
	return meanEmbeddings(embs)
}

func meanEmbeddings(embs [][]float32) []float32 {
	if len(embs) == 0 {
		return nil
	}
	dim := len(embs[0])
	out := make([]float32, dim)
	for _, e := range embs {
		for i := range e {
			out[i] += e[i]
		}
	}
	for i := range out {
		out[i] /= float32(len(embs))
	}
	return out
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

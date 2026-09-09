package service

import (
	"math"
	"testing"

	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

func TestEmbeddingCipher_RoundTrip(t *testing.T) {
	c, err := NewEmbeddingCipher("0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("NewEmbeddingCipher: %v", err)
	}
	plain := []byte("0123456789abcdefghijklmnopqrstuv") // 32 bytes
	enc, err := c.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if string(enc) == string(plain) {
		t.Error("ciphertext must differ from plaintext")
	}
	dec, err := c.Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if string(dec) != string(plain) {
		t.Errorf("roundtrip mismatch: got %q want %q", dec, plain)
	}
	// tamper detection
	enc[0] ^= 0xff
	if _, err := c.Decrypt(enc); err == nil {
		t.Error("expected error for tampered ciphertext")
	}
}

func TestFaceService_DisabledGate(t *testing.T) {
	fs := &FaceService{}
	fs.EnableForTest(false)
	if _, err := fs.ProcessPhoto("p1", "/x.jpg"); err == nil {
		t.Error("expected ErrFaceRecognitionDisabled when disabled")
	}
	if _, err := fs.ProcessLibraryBatch(5); err == nil {
		t.Error("expected ErrFaceRecognitionDisabled when disabled")
	}
	if err := fs.ReclusterAll(); err == nil {
		t.Error("expected ErrFaceRecognitionDisabled when disabled")
	}
}

func TestEmbeddingEncryptDecrypt(t *testing.T) {
	c, err := NewEmbeddingCipher("0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}
	fs := &FaceService{crypto: c}
	var e [128]float32
	for i := range e {
		e[i] = float32(math.Sin(float64(i))) // deterministic, not all-zero
	}
	enc, err := fs.encryptEmbedding(e)
	if err != nil {
		t.Fatalf("encryptEmbedding: %v", err)
	}
	dec, err := fs.decryptEmbedding(enc)
	if err != nil {
		t.Fatalf("decryptEmbedding: %v", err)
	}
	for i := range dec {
		if math.Abs(float64(dec[i]-e[i])) > 1e-6 {
			t.Fatalf("embedding mismatch at %d: %v != %v", i, dec[i], e[i])
		}
	}
}

func TestCosineAndMean(t *testing.T) {
	a := []float32{1, 0, 0}
	b := []float32{0, 1, 0}
	c := []float32{0.99, 0.141, 0} // near a
	if s := cosine(a, b); s > 0.01 {
		t.Errorf("orthogonal cosine should be ~0, got %v", s)
	}
	if s := cosine(a, c); s < 0.9 {
		t.Errorf("near-identical cosine should be high, got %v", s)
	}
	m := meanEmbeddings([][]float32{{1, 2, 3}, {3, 2, 1}})
	if m[0] != 2 || m[1] != 2 || m[2] != 2 {
		t.Errorf("mean mismatch: %v", m)
	}
}

// fakeFaceRepo is a minimal in-memory FaceRepository for clustering tests.
type fakeFaceRepo struct {
	dets []domain.FaceDetection
}

func (f *fakeFaceRepo) UnassignedDetections(_ int) ([]domain.FaceDetection, error) {
	var out []domain.FaceDetection
	for _, d := range f.dets {
		if d.PersonID == nil {
			out = append(out, d)
		}
	}
	return out, nil
}
func (f *fakeFaceRepo) CreatePerson(name string) (*domain.FacePerson, error) {
	return &domain.FacePerson{ID: name + "-id", Name: name}, nil
}
func (f *fakeFaceRepo) AssignDetectionsToPerson(ids []string, personID string) error {
	for i := range f.dets {
		for _, id := range ids {
			if f.dets[i].ID == id {
				f.dets[i].PersonID = &personID
				f.dets[i].Status = "assigned"
			}
		}
	}
	return nil
}
func (f *fakeFaceRepo) SetCoverFace(_, _ string) error { return nil }
func (f *fakeFaceRepo) Persons() ([]domain.FacePerson, error) {
	seen := map[string]bool{}
	var out []domain.FacePerson
	for _, d := range f.dets {
		if d.PersonID != nil && !seen[*d.PersonID] {
			seen[*d.PersonID] = true
			out = append(out, domain.FacePerson{ID: *d.PersonID, Name: *d.PersonID})
		}
	}
	return out, nil
}
func (f *fakeFaceRepo) FacesForPerson(personID string, _ int) ([]domain.FaceDetection, error) {
	var out []domain.FaceDetection
	for _, d := range f.dets {
		if d.PersonID != nil && *d.PersonID == personID {
			out = append(out, d)
		}
	}
	return out, nil
}
func (f *fakeFaceRepo) UpsertCluster(domain.FaceCluster) error { return nil }
func (f *fakeFaceRepo) DeleteCluster(string) error              { return nil }
func (f *fakeFaceRepo) Clusters() ([]domain.FaceCluster, error) { return nil, nil }
func (f *fakeFaceRepo) InsertDetectionsForPhoto(_ string, _ []domain.FaceDetection) error {
	return nil
}
func (f *fakeFaceRepo) DeleteDetectionsForPhoto(string) error { return nil }
func (f *fakeFaceRepo) DeleteAllFaceData() (int64, error)     { return 0, nil }
func (f *fakeFaceRepo) FacesForPhoto(_ string) ([]domain.FaceDetection, error) {
	return f.dets, nil
}
func (f *fakeFaceRepo) PersonPhotoIDs(string) ([]string, error)   { return nil, nil }
func (f *fakeFaceRepo) PersonByID(id string) (*domain.FacePerson, error) {
	return &domain.FacePerson{ID: id, Name: id}, nil
}
func (f *fakeFaceRepo) UpdatePerson(string, string) error   { return nil }
func (f *fakeFaceRepo) DeletePerson(string) error           { return nil }
func (f *fakeFaceRepo) AllDetections(int) ([]domain.FaceDetection, error) {
	return f.dets, nil
}
func (f *fakeFaceRepo) SetDetectionPerson(_ string, _ *string) error { return nil }
func (f *fakeFaceRepo) FaceCount() (int64, error)                    { return int64(len(f.dets)), nil }
func (f *fakeFaceRepo) SetCoverFace(_, _ string) error               { return nil }
func (f *fakeFaceRepo) FaceSetting() (bool, error)                   { return true, nil }
func (f *fakeFaceRepo) SetFaceSetting(bool) error                    { return nil }

func TestGroupUnassigned_SeparatesPeople(t *testing.T) {
	c, _ := NewEmbeddingCipher("0123456789abcdef0123456789abcdef")
	fs := &FaceService{crypto: c}

	// Two clear clusters: A-ish embeddings (primary 0.9/1.0) and B-ish (0.1/0.0).
	mk := func(id string, comp0 float32) domain.FaceDetection {
		var e [128]float32
		e[0] = comp0
		enc, _ := fs.encryptEmbedding(e)
		return domain.FaceDetection{ID: id, Score: 0.9, Embedding: enc, Status: "detected"}
	}
	repo := &fakeFaceRepo{dets: []domain.FaceDetection{
		mk("a1", 1.0), mk("a2", 0.9), mk("b1", 0.0), mk("b2", 0.1),
	}}
	fs.faceRepo = repo
	created, err := fs.groupUnassigned(repo.dets)
	if err != nil {
		t.Fatalf("groupUnassigned: %v", err)
	}
	persons, _ := repo.Persons()
	if created != 2 || len(persons) != 2 {
		t.Fatalf("expected 2 people from 2 clusters, got created=%d persons=%d", created, len(persons))
	}
	// a1 and a2 must be in the same person; b1/b2 in the other.
	var a1p, a2p, b1p string
	for _, d := range repo.dets {
		switch d.ID {
		case "a1":
			a1p = *d.PersonID
		case "a2":
			a2p = *d.PersonID
		case "b1":
			b1p = *d.PersonID
		}
	}
	if a1p == "" || a1p != a2p {
		t.Errorf("a1 and a2 should share a person: %q vs %q", a1p, a2p)
	}
	if b1p == "" || b1p == a1p {
		t.Errorf("b1 should be in a different person than a1: %q vs %q", b1p, a1p)
	}
}

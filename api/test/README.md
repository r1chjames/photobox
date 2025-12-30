# Integration Tests

This directory contains integration tests for the Photobox API backend.

## Overview

Integration tests verify the behavior of the system when components interact with real external dependencies:
- PostgreSQL database
- Filesystem operations
- Image processing
- Concurrent operations

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.24+

### Running Integration Tests

1. **Start the test database:**
   ```bash
   docker-compose -f test/docker-compose.test.yml up -d
   ```

2. **Wait for database to be ready:**
   ```bash
   docker-compose -f test/docker-compose.test.yml ps
   ```

3. **Run integration tests:**
   ```bash
   # From the api directory
   go test -tags=integration ./test/integration/... -v
   ```

4. **Stop the test database:**
   ```bash
   docker-compose -f test/docker-compose.test.yml down -v
   ```

### Using Make (if you have a Makefile)

```bash
make test-integration        # Run integration tests
make test-integration-up     # Start test database
make test-integration-down   # Stop test database
```

## Test Structure

```
test/
├── docker-compose.test.yml       # Test database configuration
├── integration/                  # Integration test files
│   ├── filesystem_integration_test.go
│   └── database_integration_test.go
├── testutil/                     # Test utilities
│   ├── database.go              # Database setup helpers
│   └── filesystem.go            # Filesystem test helpers
└── README.md                    # This file
```

## Test Tags

Integration tests use build tags to separate them from unit tests:

```go
//go:build integration
// +build integration
```

### Run Only Unit Tests (Skip Integration)
```bash
go test ./internal/... -short
```

### Run Only Integration Tests
```bash
go test -tags=integration ./test/integration/...
```

### Run All Tests
```bash
go test -tags=integration ./...
```

## What's Tested

### Filesystem Integration Tests
- **ScanFilesystem**: Real directory traversal with concurrency
- **GenerateThumbnail**: Actual image processing (JPEG, PNG)
- **GetMetaData**: EXIF extraction and file hashing
- **CreateDirectoryIfNotExists**: Directory creation
- **Concurrency**: Race condition testing

### Database Integration Tests
- **Albums**: CRUD operations and pagination
- **Photos**: Creating, listing, and counting photos
- **Users**: Authentication flow and user management
- **Jobs**: Job status management
- **Settings**: Configuration storage
- **Pagination**: Real database pagination with cursor-based navigation

## Database Configuration

Test database runs on **port 5433** (not 5432) to avoid conflicts with development database.

**Connection String:**
```
host=localhost user=photobox_test password=photobox_test dbname=photobox_test port=5433 sslmode=disable
```

## Test Utilities

### Database Helpers

```go
// Setup test database with migrations
db := testutil.SetupTestDB(t)
defer testutil.TeardownTestDB(t, db)

// Create full environment
env := testutil.CreateTestEnv(t)
defer testutil.CleanupTestEnv(t, env)

// Seed common test data
testutil.SeedTestData(t, db)
```

### Filesystem Helpers

```go
// Create test directory
testDir := testutil.CreateTestPhotoDir(t)
defer testutil.CleanupTestPhotoDir(t, testDir)

// Create test images
testutil.CreateTestImage(t, "photo.jpg", 1920, 1080, "jpg")

// Create full photo structure
testutil.CreateTestPhotoStructure(t, testDir)
```

## Writing New Integration Tests

1. **Add build tag at top of file:**
   ```go
   //go:build integration
   // +build integration
   ```

2. **Use testutil helpers:**
   ```go
   func TestMyIntegration(t *testing.T) {
       env := testutil.CreateTestEnv(t)
       defer testutil.CleanupTestEnv(t, env)

       // Your test code
   }
   ```

3. **Clean up resources:**
   Always use `defer` to ensure cleanup happens even if test fails.

## Troubleshooting

### Database Connection Fails
```bash
# Check if database is running
docker-compose -f test/docker-compose.test.yml ps

# Check logs
docker-compose -f test/docker-compose.test.yml logs postgres-test

# Restart database
docker-compose -f test/docker-compose.test.yml restart
```

### Port Already in Use
If port 5433 is in use, modify `docker-compose.test.yml` to use a different port.

### Tests Are Slow
Integration tests are slower than unit tests by design. Run them selectively:
```bash
# Run specific test
go test -tags=integration ./test/integration/ -run TestAlbumRepository_CreateAndGet

# Run specific file
go test -tags=integration ./test/integration/database_integration_test.go
```

## CI/CD Integration

### GitHub Actions Example
```yaml
- name: Start test database
  run: docker-compose -f api/test/docker-compose.test.yml up -d

- name: Run integration tests
  run: |
    cd api
    go test -tags=integration ./test/integration/... -v

- name: Stop test database
  run: docker-compose -f api/test/docker-compose.test.yml down -v
```

## Best Practices

1. **Isolation**: Each test should clean up after itself
2. **Idempotency**: Tests should be runnable in any order
3. **Realistic Data**: Use realistic test data sizes and structures
4. **Performance**: Monitor test execution time
5. **Cleanup**: Always use `defer` for cleanup
6. **Logging**: Use `t.Logf()` for debugging information

## Coverage

To see integration test coverage:
```bash
go test -tags=integration ./test/integration/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

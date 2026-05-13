# Integration Tests - Quick Start Guide

## TL;DR

```bash
# Start test database
make test-integration-up

# Run integration tests
make test-integration

# Stop test database
make test-integration-down
```

Or run everything in one command:
```bash
make test-integration-full
```

## Step-by-Step

### 1. Prerequisites Check

```bash
# Check Docker is running
docker ps

# Check you're in the api directory
pwd  # Should end with /photobox/api
```

### 2. Start Test Database

```bash
make test-integration-up
```

This will:
- Start PostgreSQL 15 in Docker on port 5433
- Create database `photobox_test`
- Wait for database to be ready

### 3. Run Integration Tests

```bash
# Run all integration tests
make test-integration

# Run specific test
go test -tags=integration ./test/integration/ -run TestAlbumRepository_CreateAndGet -v

# Run specific test file
go test -tags=integration ./test/integration/database_integration_test.go -v
```

### 4. View Results

```bash
# Integration tests will output:
# - Test names
# - PASS/FAIL status
# - Timing information
# - Any error messages
```

### 5. Stop Test Database

```bash
make test-integration-down
```

## Common Commands

```bash
# Run all tests (unit + integration)
make test-all

# Run only unit tests (fast)
make go-test

# Generate coverage report
make test-coverage

# Clean database volumes
make test-integration-clean

# See all available commands
make help
```

## Troubleshooting

### "Connection refused" errors
```bash
# Make sure database is running
docker ps | grep photobox-test-db

# Restart database
make test-integration-down
make test-integration-up
```

### "Port 5433 already in use"
```bash
# Check what's using the port
lsof -i :5433

# Or modify test/docker-compose.test.yml to use different port
```

### Tests timeout
```bash
# Integration tests can be slow - increase timeout
go test -tags=integration ./test/integration/... -timeout 60m
```

## What Gets Tested?

### Filesystem Tests
- ✅ Scanning directories with real files
- ✅ Generating thumbnails from actual images
- ✅ EXIF data extraction
- ✅ Concurrent filesystem operations

### Database Tests
- ✅ Album CRUD operations
- ✅ Photo management with relationships
- ✅ User authentication
- ✅ Job status tracking
- ✅ Settings management
- ✅ Pagination with real data

## Next Steps

1. Run integration tests: `make test-integration-full`
2. Check coverage: `make test-coverage`
3. Add more tests as needed
4. Integrate into CI/CD pipeline

## Need Help?

- See full documentation: [README.md](README.md)
- Check test logs: `docker-compose -f test/docker-compose.test.yml logs`
- Run specific test for debugging: `go test -tags=integration ./test/integration/ -run TestName -v`

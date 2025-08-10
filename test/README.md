# RARA Membership API Test Suite

This directory contains comprehensive test scripts for the RARA Membership API endpoints.

## Prerequisites

- Go installed and configured
- `curl` command available
- `jq` command for JSON parsing (optional but recommended)
- The RARA membership application built (`go build`)

## Test Scripts

### 1. `run_all_tests.sh`
Master test script that runs all API tests with a fresh database.

**Features:**
- Backs up existing database
- Creates fresh test database
- Starts the server automatically
- Runs all test suites in order
- Provides comprehensive test summary
- Cleans up after tests

**Usage:**
```bash
./test/run_all_tests.sh
```

### 2. `test_members.sh`
Tests all member-related API endpoints:
- List members
- Create member
- Update member
- Get member points history

### 3. `test_events.sh`
Tests all event-related API endpoints:
- List events
- Create event
- Update event

### 4. `test_points.sh`
Tests all point management endpoints:
- Allocate points
- Redeem points
- Check balance
- Points history

### 5. `interactive_test.sh`
Interactive menu-driven test client for manual testing:
- User-friendly interface
- Test individual endpoints
- Create test data interactively
- Run automated test scenarios

**Usage:**
```bash
./test/interactive_test.sh
```

### 6. `test_helpers.sh`
Shared helper functions used by other test scripts.

## Running Tests

### Run All Tests (Recommended)
```bash
# From project root
./test/run_all_tests.sh

# Or from test directory
cd test
./run_all_tests.sh
```

### Run Individual Test Suites
First, start the server:
```bash
./rara-membership serve
```

Then run individual tests:
```bash
# Test members endpoints
./test/test_members.sh

# Test events endpoints
./test/test_events.sh

# Test points endpoints
./test/test_points.sh
```

### Interactive Testing
```bash
# Start the server
./rara-membership serve

# In another terminal, run interactive tester
./test/interactive_test.sh
```

## Test Configuration

### Environment Variables
- `BASE_URL`: API base URL (default: `http://localhost:8080`)
- `PORT`: Server port (default: `8080`)

### Custom Base URL
```bash
BASE_URL=http://localhost:3000 ./test/run_all_tests.sh
```

## Test Coverage

### Member Endpoints
- ✅ GET `/api/members` - List active members
- ✅ POST `/api/members` - Create new member
- ✅ PUT `/api/members/{call_sign}` - Update member by call sign
- ✅ GET `/api/members/{call_sign}/points` - Get member point history and balance

### Event Endpoints
- ✅ GET `/api/events` - List all events
- ✅ POST `/api/events` - Create new event
- ✅ PUT `/api/events/{id}` - Update event by ID

### Points Endpoints
- ✅ POST `/api/members/{call_sign}/points` - Allocate points to member
- ✅ POST `/api/members/{call_sign}/redemptions` - Redeem member points

## Test Scenarios

### Positive Tests
- Create valid members with all fields
- Create events with default points
- Allocate points using event defaults
- Allocate specific point amounts
- Redeem points within balance
- Update member information
- Update event details

### Negative Tests
- Create duplicate members (same call sign)
- Create members with missing required fields
- Update non-existent members
- Redeem more points than available
- Allocate negative points
- Allocate excessive points (>100)
- Invalid member/event IDs

### Edge Cases
- Multiple point allocations for same member/event
- Partial updates with null fields
- Empty responses
- Case-insensitive call sign handling

## Expected Test Output

Successful test run shows:
```
Testing: GET /api/members... PASSED (Status: 200)
Testing: POST /api/members/create... PASSED (Status: 201)
...
All tests completed!
```

Failed tests show:
```
Testing: POST /api/members/create (duplicate)... FAILED (Expected: 409, Got: 500)
  Response: {"error": "constraint violation"}
```

## Troubleshooting

### Server Won't Start
- Check if port 8080 is already in use
- Verify database file permissions
- Check server logs: `/tmp/rara-test-server.log`

### Tests Fail
- Ensure server is running: `curl http://localhost:8080/api/members`
- Check if database is initialized
- Verify `jq` is installed for JSON parsing

### Database Issues
- Remove `membership.db` to start fresh
- Check file permissions in project directory
- Ensure SQLite is properly installed

## Development

### Adding New Tests
1. Create new test script in `test/` directory
2. Source `test_helpers.sh` for common functions
3. Follow existing test patterns
4. Update `run_all_tests.sh` to include new tests

### Test Best Practices
- Always test both success and failure cases
- Clean up test data after tests
- Use meaningful test names
- Check response status codes
- Validate response body when applicable

## Contributing

When adding new API endpoints:
1. Create corresponding tests
2. Update this README
3. Ensure tests pass before committing
4. Include both positive and negative test cases
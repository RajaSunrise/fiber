# Code Audit Report

## Bugs Found & Fixed

### High Severity: Race Conditions in `client/client.go`

The `Client` struct in `client/client.go` is designed to be shared across requests. However, several configuration methods were modifying shared state without acquiring the lock, leading to potential race conditions.

**Fix Applied:**
Added `c.mu.Lock()`/`Unlock()` (or `RLock`/`RUnlock`) to the following methods in `client/client.go`:
- `BaseURL`, `SetBaseURL`
- `Header`, `AddHeader`, `SetHeader`, `AddHeaders`, `SetHeaders`
- `Param`, `AddParam`, `SetParam`, `AddParams`, `SetParams`, `SetParamsWithStruct`, `DelParams`
- `SetUserAgent`, `SetReferer`
- `DisablePathNormalizing`, `SetDisablePathNormalizing`
- `PathParam`, `SetPathParam`, `SetPathParams`, `SetPathParamsWithStruct`, `DelPathParams`
- `Cookie`, `SetCookie`, `SetCookies`, `SetCookiesWithStruct`, `DelCookies`
- `SetTimeout`, `Debug`, `DisableDebug`, `SetCookieJar`

Also updated `client/core.go` to thread-safely read `c.client.timeout`.

### Low Severity: Redundant TODO in `router.go`

Found a redundant `// TODO: Do we need to return here?` in `router.go` at the end of a void function. Removed it.

## Known Issues (Not Fixed due to Test Constraints)

### Medium Severity: Potential Integer Overflow in `internal/memory` Expiration

In `internal/memory/memory.go`, the expiration time is calculated as `uint32(ttl.Seconds()) + utils.Timestamp()`.
`utils.Timestamp()` returns the current time as `uint32`.
If `ttl` is sufficiently large (e.g., > 80 years from now), the addition can overflow `uint32`, wrapping around to a small value. This would cause the cache entry to be treated as expired immediately (or very soon).

**Recommendation:**
Check for overflow or cap the expiration time at `MaxUint32`. Alternatively, verify if such long TTLs are supported/intended.

### Other Issues

1.  **`RemoveEscapeChar` Logic:** The `RemoveEscapeChar` function in `path.go` incorrectly removes *all* backslashes (e.g., `\\` -> empty string, instead of `\`). Fixing this requires updating existing tests in `path_test.go` which rely on the incorrect behavior.
2.  **`CheckConstraint` Error Ignoring:** The `CheckConstraint` function in `path.go` ignores errors from `strconv.Atoi` when parsing route constraints (e.g., `range(10,30)`). This allows invalid constraints to pass silently (or default to 0). Fixing this requires updating existing tests that use invalid constraint syntax (e.g., `range(10\,30)`).
3.  **Router Optional Parameter Matching:** The router eagerly matches optional parameters (forcing length 1) even if it causes subsequent mandatory parameters to fail matching. This is a known issue marked with TODOs in tests.

## Documentation Review

- Verified `docs/middleware/logger.md`: Usage of "referer" appears correct as it refers to the template tag `${referer}`.
- Verified function comments in `app.go` and `router.go`: Generally consistent with signatures.

## Typos

- No significant new typos found in code comments or strings.

# Templ Implementation Improvements

## Current Issues

### 1. Manual Generation Required
Currently, you must manually run `templ generate` before building or running the server:

```bash
make build   # Runs: cd server && templ generate && go build...
make run     # Runs: cd server && templ generate && go run...
```

### 2. Air Hot Reloading Incomplete
The `.air.toml` configuration has two problems:

**Problem A: Missing File Extension**
- Currently watches: `["go", "tpl", "tmpl", "html"]`
- Missing: `"templ"` extension
- Result: Air doesn't detect changes to `.templ` files

**Problem B: No Pre-Build Generation**
- Air builds Go code directly without running `templ generate` first
- Result: Changes to `.templ` files don't regenerate `*_templ.go` files
- Consequence: Stale generated code causes runtime errors or no visible changes

### 3. Development Workflow Friction
Current workflow:
1. Edit `.templ` file
2. Run `templ generate` manually (or use make command)
3. Air detects `.go` file change
4. Air rebuilds and restarts

Desired workflow:
1. Edit `.templ` file
2. Air automatically handles everything
3. See changes immediately

---

## Proposed Solution

### Changes to `.air.toml`

#### Change 1: Add "templ" Extension
Update the `include_ext` array to watch `.templ` files:

```toml
# Before
include_ext = ["go", "tpl", "tmpl", "html"]

# After
include_ext = ["go", "tpl", "tmpl", "html", "templ"]
```

#### Change 2: Add Pre-Build Command
Add a `pre_cmd` to automatically run `templ generate` before each build:

```toml
[build]
  # ... existing config ...
  cmd = "go build -o ./tmp/main ./cmd/api"

  # Add this new line:
  pre_cmd = ["templ generate"]
```

This ensures that whenever Air detects a change (in `.templ` OR `.go` files), it:
1. Runs `templ generate` first
2. Then runs the Go build command
3. Then restarts the server

---

## Complete Updated `.air.toml`

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/api"
  pre_cmd = ["templ generate"]  # NEW: Run templ generate before build
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html", "templ"]  # UPDATED: Added "templ"
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

---

## Verification: Tiltfile (Already Correct)

The Tiltfile already has proper templ generation configured:

```python
docker_build(
    'bananas-server-' + DOCKER_ENV,
    # ...
    live_update=[
        sync('./server', '/app'),
        run('templ generate', trigger=['./internal/templates/*.templ']),  # ✅ Good!
    ]
)
```

This means:
- ✅ Docker/Tilt environment already handles templ generation automatically
- ✅ Only local development (Air) needs the fix

---

## Expected Workflow After Changes

### Local Development (Air)
1. Run `air` in the `server/` directory
2. Edit any `.templ` file
3. Air detects change → runs `templ generate` → rebuilds → restarts
4. See changes at `http://localhost:8081/templ`

### Docker Development (Tilt)
1. Run `tilt up`
2. Edit any `.templ` file
3. Tilt detects change → runs `templ generate` in container → rebuilds → restarts
4. See changes immediately

### Make Commands (Still Work)
The `make build` and `make run` commands will continue to work as-is:
- They explicitly run `templ generate` first
- No changes needed to Makefile

---

## Testing the Changes

### Test 1: Verify Air Watches Templ Files
```bash
cd server
air

# In another terminal, touch a templ file
touch internal/templates/home.templ

# Expected: Air should detect the change and rebuild
```

### Test 2: Verify Templ Generation Runs Automatically
```bash
cd server

# Make a visible change to a template
# Edit internal/templates/home.templ and change the header text

# Air should:
# 1. Detect the .templ file change
# 2. Run templ generate
# 3. Rebuild the Go binary
# 4. Restart the server

# Verify at http://localhost:8081/templ
```

### Test 3: Verify No Manual Generation Needed
```bash
cd server
rm internal/templates/*_templ.go  # Delete generated files
air                                # Start air

# Air should regenerate the files automatically on first build
ls internal/templates/*_templ.go  # Files should exist
```

---

## Benefits

1. **No Manual Commands** - Never run `templ generate` manually again
2. **Faster Development** - Edit and see changes immediately
3. **Fewer Errors** - No stale generated code
4. **Consistent Workflow** - Same hot-reload experience as editing Go files
5. **Works Everywhere** - Both local (Air) and Docker (Tilt) environments

---

## Alternative: Air + Templ Watch Mode

If Air's `pre_cmd` causes performance issues (running on every change), you could use templ's built-in watch mode alongside Air:

```bash
# Terminal 1: Templ watch mode
cd server
templ generate --watch

# Terminal 2: Air for Go hot reload
cd server
air
```

However, the `pre_cmd` approach is simpler and should work well for this project size.

---

## File Locations

- `.air.toml` - `/home/bobparsons/Development/bananas/server/.air.toml`
- Templ templates - `/home/bobparsons/Development/bananas/server/internal/templates/*.templ`
- Generated files - `/home/bobparsons/Development/bananas/server/internal/templates/*_templ.go`

---

## Next Steps

1. ✅ Review this document
2. ⏳ Update `.air.toml` with the proposed changes
3. ⏳ Test the new workflow
4. ⏳ Update documentation if needed
5. ⏳ Consider adding a `make dev` command that runs `air` for consistency

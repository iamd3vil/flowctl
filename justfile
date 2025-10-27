build: build-site
    #!/bin/bash
    set -euo pipefail
    VERSION=$(git describe --tags --exact-match 2>/dev/null || echo "dev")
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    go build -ldflags="-s -w -X 'github.com/cvhariharan/flowctl/cmd.version=${VERSION}' -X 'github.com/cvhariharan/flowctl/cmd.commit=${COMMIT}' -X 'github.com/cvhariharan/flowctl/cmd.date=${DATE}'" -o flowctl

run:
    #!/bin/bash
    set -euo pipefail
    cd site && npm run dev > /dev/null 2>&1 &
    echo "Site available at: http://localhost:5555"
    DEBUG_LOG=true go run main.go start

build-site: install
    #!/bin/bash
    set -euo pipefail
    cd site && npm run build

clean:
    rm -rf site/build
    rm -f flowctl

install:
    cd site && npm install

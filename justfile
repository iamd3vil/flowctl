build: build-site
    go build -o flowctl

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

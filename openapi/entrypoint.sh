#!/usr/bin/env sh

set -eu

bundle () {
    output=$(redocly bundle "spec/openapi.yaml" -o bundled/openapi.json 2>&1)
    if [ $? -ne 0 ]; then
        echo "$output"
        echo "error: failed to bundle the openapi spec"
        exit 1
    fi

    echo "info: openapi.json bundled"
}

watch () {
    echo "info: watching for changes in spec/"
    while inotifywait -q -r -e close_write "spec/"
    do
        bundle
    done
}

bundle
watch &

npm run serve

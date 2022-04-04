#!/usr/bin/env bash
set -Eeuo pipefail

# prefix all output
exec > >(trap "" INT TERM; sed 's/^/[JEST] /')
exec 2> >(trap "" INT TERM; sed 's/^/[JEST] /' >&2)

# cd to script directory
cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null

# start tests
npm test -- --colors

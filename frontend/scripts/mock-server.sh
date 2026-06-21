#!/usr/bin/env bash
set -euo pipefail

# Start Prism mock server from the TypeSpec-generated OpenAPI spec
# The frontend dev server proxies /api/* to this mock server

OPENAPI_SPEC="${1:-../typespec/output/openapi.yaml}"

if [ ! -f "$OPENAPI_SPEC" ]; then
  echo "Error: OpenAPI spec not found at $OPENAPI_SPEC"
  echo "Usage: $0 [path-to-openapi.yaml]"
  exit 1
fi

echo "Starting Prism mock server from: $OPENAPI_SPEC"
./node_modules/.bin/prism mock "$OPENAPI_SPEC" -p 4010

#!/bin/bash
# Generate Go code from proto files
#
# Prerequisites:
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#
# Usage:
#   ./generate.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PKG_DIR="$(dirname "$SCRIPT_DIR")"

echo "Generating Go code from proto files..."

for service in auth wallet provider parking notification; do
    echo "  Generating ${service}..."
    protoc \
        --proto_path="${SCRIPT_DIR}" \
        --go_out="${PKG_DIR}" \
        --go_opt=paths=source_relative \
        --go-grpc_out="${PKG_DIR}" \
        --go-grpc_opt=paths=source_relative \
        "${SCRIPT_DIR}/${service}/v1/${service}.proto"
done

echo "Proto generation complete!"
echo ""
echo "Generated files:"
find "${PKG_DIR}/proto" -name "*.pb.go" -type f

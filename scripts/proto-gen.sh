#!/bin/bash
set -e

PROTO_DIR="./proto"
OUT_DIR="./proto"

echo "Generating gRPC code from proto files..."

for service_dir in "$PROTO_DIR"/*/; do
  service=$(basename "$service_dir")
  proto_file="$service_dir/$service.proto"

  if [ ! -f "$proto_file" ]; then
    continue
  fi

  echo "  Generating: $service..."

  protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$OUT_DIR" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$OUT_DIR" \
    --go-grpc_opt=paths=source_relative \
    "$proto_file"
done

echo "Done! Generated Go files in $OUT_DIR"

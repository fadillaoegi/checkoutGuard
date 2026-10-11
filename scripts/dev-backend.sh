#!/bin/sh
set -eu
checkoutguard_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
checkoutguard_neatlogs="$checkoutguard_root/../.neatlogs-api.env"
checkoutguard_model="$checkoutguard_root/../.apertus-api.env"
if [ -f "$checkoutguard_root/.env.neatlogs.local" ]; then checkoutguard_neatlogs="$checkoutguard_root/.env.neatlogs.local"; fi
if [ -f "$checkoutguard_root/.env.model.local" ]; then checkoutguard_model="$checkoutguard_root/.env.model.local"; fi
export NEATLOGS_ENV_FILE="${NEATLOGS_ENV_FILE:-$checkoutguard_neatlogs}"
export MODEL_ENV_FILE="${MODEL_ENV_FILE:-$checkoutguard_model}"
if [ -f "$checkoutguard_root/.env.pricing.local" ]; then
  export MODEL_PRICING_ENV_FILE="${MODEL_PRICING_ENV_FILE:-$checkoutguard_root/.env.pricing.local}"
fi
cd "$checkoutguard_root/checkoutguard_BE"
mkdir -p bin
go build -o bin/server ./cmd/server
exec bin/server

#!/usr/bin/env bash
set -euo pipefail

docker build -f benzhi.Dockerfile -t venueops:latest .

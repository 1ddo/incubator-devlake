#!/bin/bash

cd "$(dirname "$0")"
poetry run python argocd/main.py "$@"
#!/bin/bash

export SLUG=ghcr.io/awakari/int-bluesky
export VERSION=latest
docker tag awakari/int-bluesky "${SLUG}":"${VERSION}"
docker push "${SLUG}":"${VERSION}"

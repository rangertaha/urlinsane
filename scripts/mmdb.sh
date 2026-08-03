#!/usr/bin/env bash
#
# Fetch the MaxMind GeoLite2 database that enables the `geo` operator.
#
# It installs into the application directory, not the source tree. The database
# used to be embedded with //go:embed and this script wrote it to
# internal/config/ for `make build` to pick up; it is not embedded any more
# (189a5ef), because the copy that shipped was corrupt and 49 MB of unusable
# bytes rode along in every release. Geolocation is opt-in now: this file is
# read from disk at startup and the operator is planned only if it is there.
#
# Leaving the old destination in place would have been worse than no script at
# all, because config.optional() names this script as the remedy — so following
# the advice would have written 49 MB somewhere nothing reads and reported
# success.
#
# Earlier still it was broken three ways: it wrote to
# internal/plugins/collectors/geo, a directory deleted with the old collectors;
# it fetched from git.io, which GitHub retired in 2022; and its guard read
# `[ -f ! $FILE ]`, which tests the literal "-f" and so never fired.
set -euo pipefail

# Matches config.DirName + config.MaxMindDB. Override for a non-standard home.
DEST="${URLINSANE_DIR:-$HOME/.config/urlinsane}/maxmind.db.gz"
TMP="$(mktemp)"
trap 'rm -f "$TMP"' EXIT

# A licence key is required: MaxMind closed anonymous GeoLite2 downloads in
# 2019. Set MAXMIND_LICENSE_KEY, or point MAXMIND_URL at your own mirror.
URL="${MAXMIND_URL:-}"
if [ -z "$URL" ]; then
    if [ -z "${MAXMIND_LICENSE_KEY:-}" ]; then
        echo "error: set MAXMIND_LICENSE_KEY, or MAXMIND_URL to a mirror" >&2
        echo "       sign up free at https://www.maxmind.com/en/geolite2/signup" >&2
        exit 1
    fi
    URL="https://download.maxmind.com/app/geoip_download"
    URL="$URL?edition_id=GeoLite2-City&license_key=$MAXMIND_LICENSE_KEY&suffix=tar.gz"
fi

echo "fetching $DEST"
curl -fsSL "$URL" -o "$TMP"

# Verify before installing. The whole reason this file needs care is that a
# corrupt copy is indistinguishable from a good one until a scan omits geo.
if ! gzip -t "$TMP" 2>/dev/null; then
    echo "error: downloaded file is not valid gzip" >&2
    exit 1
fi

mkdir -p "$(dirname "$DEST")"
mv "$TMP" "$DEST"
trap - EXIT
echo "wrote $DEST ($(stat -c%s "$DEST") bytes)"
echo "the geo operator will appear in the plan on your next scan; no rebuild needed"

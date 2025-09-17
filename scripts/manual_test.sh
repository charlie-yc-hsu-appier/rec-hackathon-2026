#!/bin/bash

# Manual test script
# Usage: ./manual_test.sh {vendor_key} {sub_id} {w} {h}

VENDOR_KEY="$1"
SUBID="$2"
W="$3"
H="$4"

# Create temporary directory
TMPDIR=$(mktemp -d)
trap 'rm -r "$TMPDIR"' EXIT

echo "=== Testing $VENDOR_KEY ($SUBID) - Expected: ${W}x${H} ==="

# Make request and check status code
URL="http://127.0.0.1:8080/r/${VENDOR_KEY}?subid=${SUBID}&w=${W}&h=${H}&user_id=2391D054-7685-4A4C-94FD-8F89A7EB3877&click_id=snmgUTVwSIq3HHn3Qobuhg.bXy9p67qBUORacZpF0u6aA&partner_id=kakao&lat=22.3264&lon=114.1661&k_campaign_id=1917120732480000069"

echo "Making request to: $URL"

# Get response and status code
RESPONSE_FILE="$TMPDIR/response.json"
HTTP_CODE=$(curl -s -o "$RESPONSE_FILE" -w "%{http_code}" "$URL")

if [ "$HTTP_CODE" -ne 200 ]; then
    echo "  ❌ Request failed with HTTP status code: $HTTP_CODE"
    if [ -f "$RESPONSE_FILE" ]; then
        echo "Response body:"
        cat "$RESPONSE_FILE"
    fi
    exit 1
fi

echo "  ✅ Request successful (HTTP 200)"

# Skip image validation for keeta vendor
if [ "$VENDOR_KEY" = "keeta" ]; then
    echo "Keeta vendor detected - skipping image dimension validation"
    echo "  🎉 Keeta request successful!"
    exit 0
fi

# Process the response
while read -r url; do
    if [ -n "$url" ]; then
        echo "Image: $url"
        temp=$(mktemp -p "$TMPDIR")
        if curl -s -L "$url" -o "$temp"; then
            dim=$(file "$temp" | grep -o '[0-9]\+x[0-9]\+' | tail -1)
            if [ -n "$dim" ]; then
                w_actual=$(echo "$dim" | cut -d'x' -f1)
                h_actual=$(echo "$dim" | cut -d'x' -f2)
                echo "Actual: ${w_actual}x${h_actual}"

                # Check for 300x300
                if [ "$W" -eq 300 ] && [ "$H" -eq 300 ]; then
                    if [ "$w_actual" -eq "$h_actual" ]; then
                        echo "  ✅ Square image"
                        echo "  🎉 Found correct square image! Stopping validation."
                    else
                        echo "  ☑️  Not square. Skip validation since RM will rescale on the FE side. Suggest to check with TS team."
                    fi
                    break
                # Check for other dimensions
                else
                    if [ "$w_actual" -eq "$W" ] && [ "$h_actual" -eq "$H" ]; then
                        echo "  ✅ Correct dimensions"
                        echo "  🎉 Found correct dimensions! Stopping validation."
                        break
                    else
                        echo "  ❌ Wrong (expected ${W}x${H})"
                        exit 1
                    fi
                fi
            else
                echo "  ⚠️  Could not get dimensions"
            fi
        else
            echo "  ⚠️  Failed to download"
        fi
    fi
done < <(grep -o '"image":"[^"]*"' "$RESPONSE_FILE" | sed 's/"image":"//g' | sed 's/"//g')

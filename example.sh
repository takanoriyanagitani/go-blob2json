#!/bin/bash

# Example 1: Plain text input and decoding
echo 'plain-text' |
  wazero run ./blob2json.wasm -- --name plain_text.txt --content-type text/plain |
  jq -c . > /tmp/plain_text.json

echo "Original text to JSON:"
cat /tmp/plain_text.json | jq .
echo "Decoded text:"
jq -r .body /tmp/plain_text.json | base64 --decode
echo ""

# Example 2: Random bytes input
dd if=/dev/urandom of=/dev/stdout bs=4 count=1 status=none |
  wazero run ./blob2json.wasm -- --name random_bytes.bin |
  jq -c

# Example 3: JSON input
jq -c -n '{helo:"wrld"}' |
  wazero \
	run \
	./blob2json.wasm \
	-- \
	--name json_input.json \
	--content-type application/json \
	--content-encoding none \
	--max-bytes 1024 |
  jq -c

# Example 4: With metadata
echo "some data" |
  wazero \
    run \
    ./blob2json.wasm \
    -- \
    --name data.txt \
    --content-type text/plain \
    --metadata source=example.sh \
    --metadata tar_file=blobs.tar |
  jq -c

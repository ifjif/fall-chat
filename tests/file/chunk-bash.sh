#!/bin/bash

for chunk_file in ./chunks/*; do
  index=$(basename "$chunk_file")

  echo "Uploading chunk: $index"
  curl -K chunk -F "index=${index}" -F "chunk=@$chunk_file"

  if [ $? -ne 0 ]; then
    echo "Failed: 上传失败 chunk:${index}"
    exit 1
  fi

done

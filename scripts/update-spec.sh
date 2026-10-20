#!/usr/bin/env bash
# Заменяет спеку и фикстуры МойСклад версией релиза и обновляет её номер в README.
# Использование: scripts/update-spec.sh 0.34.0 && make gen
set -euo pipefail

version="${1:?укажи версию релиза спеки, например 0.34.0}"
root="$(cd "$(dirname "$0")/.." && pwd)"
repo="moysklad/api-remap-1.2-openapi-specification"

curl -fsSL -o "$root/api/openapi.yaml" \
  "https://github.com/${repo}/releases/download/${version}/openapi.yaml"
for f in product.json product_full.json product_folder.json; do
  curl -fsSL -o "$root/testdata/$f" \
    "https://raw.githubusercontent.com/${repo}/${version}/tests/fixtures/$f"
done
for f in "$root/README.md" "$root/testdata/README.md"; do
  sed -i.bak -E "s/(Версия спеки МойСклад: )[0-9]+\.[0-9]+\.[0-9]+/\1${version}/" "$f"
  rm "$f.bak"
done

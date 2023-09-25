#!/bin/bash

echo "# Resources" > docs/resources/index.md
echo >> docs/resources/index.md
echo "\`\`\`{toctree}" >> docs/resources/index.md
echo ":maxdepth: 2" >> docs/resources/index.md
echo ":hidden:" >> docs/resources/index.md

for file in docs/resources/*.md; do
  filename=$(basename "$file")
  
  if [ "$filename" != "index.md" ]; then
    filename_no_extension="${filename%.*}"
    echo "$filename_no_extension" >> docs/resources/index.md
  fi
done

echo "\`\`\`" >> docs/resources/index.md


echo "# Data Sources" > docs/data-sources/index.md
echo >> docs/data-sources/index.md
echo "\`\`\`{toctree}" >> docs/data-sources/index.md
echo ":maxdepth: 2" >> docs/data-sources/index.md
echo ":hidden:" >> docs/data-sources/index.md

for file in docs/data-sources/*.md; do
  filename=$(basename "$file")
  
  if [ "$filename" != "index.md" ]; then
    filename_no_extension="${filename%.*}"
    echo "$filename_no_extension" >> docs/data-sources/index.md
  fi
done

echo "\`\`\`" >> docs/data-sources/index.md


# echo "\`\`\`{toctree}" >> docs/index.md
# echo ":maxdepth: 2" >> docs/index.md
# echo ":hidden:" >> docs/index.md
# echo "data-sources/index" >> docs/index.md
# echo "resources/index" >> docs/index.md
# echo "\`\`\`" >> docs/index.md


# sed -i '/```\{toctree\}/,/^```$/d' docs/index.md
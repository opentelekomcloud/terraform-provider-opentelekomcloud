#!/usr/bin/env bash

function tfProviderLint {
  echo "==> Checking source code against terraform provider linters..."
	tfproviderlint \
        -AT001\
        -AT001.ignored-filename-suffixes _data_source_test.go\
        -AT005 -AT006 -AT007\
        -R001 -R002 -R003 -R004 -R006\
        -S001 -S002 -S003 -S004 -S005 -S006 -S007 -S008 -S009 -S010 -S011 -S012 -S013 -S014 -S015 -S016 -S017 -S018 -S019 -S020\
        -S021 -S022 -S023 -S024 -S025 -S026 -S027 -S028 -S029 -S030 -S031 -S032 -S033\
        ./hue/...
}

function tfAccTestsLint {
  echo "==> Checking acceptance test terraform blocks are formatted..."

files=$(find ./opentelekomcloud/acceptance -type f -name "*_test.go")
error=false

for f in $files; do
  terrafmt diff -c -q -f "$f" || error=true
done

if ${error}; then
  echo "------------------------------------------------"
  echo ""
  echo "The preceding files contain terraform blocks that are not correctly formatted or contain errors."
  echo "You can fix this by running make tools and then terrafmt on them."
  echo ""
  echo "to easily fix all terraform blocks:"
  echo "$ make tffmtfix"
  echo ""
  echo "format only acceptance test config blocks:"
  echo "$ find ./opentelekomcloud/acceptance | egrep \"_test.go\" | sort | while read f; do terrafmt fmt -f \$f; done"
  echo ""
  echo "format a single test file:"
  echo "$ terrafmt fmt -f ./opentelekomcloud/acceptance/resource_test.go"
  echo ""
  exit 1
fi

exit 0

}

function tfDocsLint {
  echo "==> Checking docs terraform blocks are formatted..."

files=$(find ./docs -type f -name "*.md")
error=false

for f in $files; do
  terrafmt diff -c -q -f "$f" || error=true
done

if ${error}; then
  echo "------------------------------------------------"
  echo ""
  echo "The preceding files contain terraform blocks that are not correctly formatted or contain errors."
  echo "You can fix this by running make tools and then terrafmt on them."
  echo ""
  echo "to easily fix all terraform blocks:"
  echo "$ make tffmtfix"
  echo ""
  echo "format only docs config blocks:"
  echo "$ find docs | egrep \".md\" | sort | while read f; do terrafmt fmt -f \$f; done"
  echo ""
  echo "format a single test file:"
  echo "$ terrafmt fmt -f ./docs/resources/resource.md"
  echo ""
  exit 1
fi

exit 0

}

function main {
  tfProviderLint
  tfAccTestsLint
  tfDocsLint

}

main

#!/bin/bash
tag=$(git describe --tags --abbrev=0)
changelog_ref="https://docs.otc-service.com/releasenotes/terraform-provider-opentelekomcloud/current.html#${tag//./-}"
printf "[Changelog](%s)\n" "${changelog_ref}"

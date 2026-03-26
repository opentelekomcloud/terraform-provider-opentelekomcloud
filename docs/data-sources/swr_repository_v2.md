---
subcategory: "Software Repository for Container (SWR)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_swr_repository_v2"
sidebar_current: "docs-opentelekomcloud-data-source-swr-repository-v2"
description: |-
  Get details of SWR repositories within OpenTelekomCloud.
---

Up-to-date reference of API arguments for SWR repository you can get at
[documentation portal](https://docs.otc.t-systems.com/software-repository-container/api-ref/api)

# opentelekomcloud_swr_repository_v2

Get details of SWR repositories within Open Telekom Cloud.

## Example Usage

```hcl
data opentelekomcloud_swr_repository_v2 repo_1 {}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Specify the name of the repository. Use this to filter repositories list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `repositories` - List of repositories. The structure is documented below:
    * `repository_id` - Numeric ID of the repository.
    * `name` - The name of the repository.
    * `is_public` - Whether the repository is public (`true`) or private (`false`).
    * `description` - Repository description.
    * `category` - Repository category.
    * `path` - Image address for docker pull.
    * `internal_path` - Intra-cluster image address for docker pull.
    * `num_images` - Number of image tags in a repository.
    * `size` - Repository size.

## 3 ways to get involved

There are three main ways you can get involved in our open-source project, and
each is described briefly below.

### 1. Fixing bugs

If you want to start fixing open bugs, we'd really appreciate that! Bug fixing
is central to any project. The best way to get started is by heading to our
[bug tracker](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/issues)
and finding open bugs that you think nobody is working on. It might be useful to comment on the
thread to see the current state of the issue and if anybody has made any breakthroughs on it so far.

### 2. Improving documentation

The provider's documentation is using
standard [terraform provider format](https://www.terraform.io/registry/providers/docs#format)
and generated version is available at
the [registry](https://registry.terraform.io/providers/opentelekomcloud/opentelekomcloud/latest/docs).

If you feel that a certain section could be improved ― whether it's to clarify
ambiguity, correct a technical mistake, or to fix a grammatical error ― please
feel entitled to do so! We welcome doc pull requests with the same childlike
enthusiasm as any other contribution!

Note that terraform registry uses own markdown [flavor](https://www.terraform.io/registry/providers/docs).
You can validate your changes using [doc preview tool](https://registry.terraform.io/tools/doc-preview).

### 3. Working on a new feature

If you've found something we've left out, definitely feel free to start work on
introducing that feature. It's always useful to open an issue or submit a pull
request early on to indicate your intent to a core contributor - this enables
quick/early feedback and can help steer you in the right direction by avoiding
known issues. It might also help you avoid losing time implementing something
that might not ever work.

Please do not hesitate to ask questions or request clarification. Your
contribution is very much appreciated, and we are happy to work with you to get
it merged.

#### PR header

Remember that the PR header or a commit header, in case of just a single commit,
will be a part of the future release description.

PR header should adhere to the following requirements:
1. Maximum 50 characters long (or 72 if it's impossible to fit in just 50 characters)
2. Prefixed with a service category in case it meant to be a part of the changelog, e.g. `[DNS]`.
3. Prefixed with one of the following in case it should not be a part of the changelog:
   - `ci:` - for CI/CD changes
   - `reno:` - for release notes changes
   - `release:` - for release process changes
   - `acceptance:` - for changes related to the acceptance tests

#### PR checklist

When checking what should be done as a part of PR besides writing code, please
refer to the checklist described in the pull request template.

#### Release Notes

Very special part of most PR is a release notes process.
We currently use [`reno`](https://docs.openstack.org/reno/latest/) to build release notes
that are hosted at [Open Telekom Cloud docs portal](https://docs-beta.otc.t-systems.com/releasenotes/terraform-provider-opentelekomcloud/)

There is the following process to create a release note:

1. Create new RN from the template:
   ```shell
   $ reno new short-description --from-template releasenotes/template.yaml
   ```

2. All but categories from the generated `yaml` file. We use
  [common categories](https://docs.openstack.org/reno/latest/user/usage.html#editing-a-release-note)
  and `enhancements` when improving some existing resources (e.g. adding some field).

3. Describe a change. Pay your attention that `ReST` formatting is used here, not a `markdown`.
   - For new data source or resource:
   ```yaml
   features:
     - |
       **New Resource:** ``full_resource_name``
     - |
       **New Data Source:** ``full_data_source_name``
   ```
   - For updating existing ones, same as in PR header:
   ```yaml
   enhancements:
     - |
       **[CAT]** Some description
   fixes:
     - |
       **[CAT]** Some description
   ```

4. Add a PR reference at the end when the PR is created e.g.:
   ```yaml
   fixes:
     - |
       **[CAT]** Some description
       (`#10 <https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/pull/10>`_)
   ```

## Tests

When working on a new or existing feature, testing will be the backbone of your
work since it helps uncover and prevent regressions in the codebase. We mostly use
acceptance tests for testing our resources and data sources.

All acceptance tests are placed in the [`acceptance` directory](/opentelekomcloud/acceptance)

### Writing tests

For writing acceptance tests please refer to
[HashiCorp guidelines](https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests/testcase)
and existing acceptance tests.

### Running tests

Acceptance tests commonly uses some additional environment variables to configure
common parameters: `OS_SUBNET_NAME` and `OS_AVAILABILITY_ZONE` are always required.
Provider support all Open Telekom Cloud authentication approaches including using client
configuration files (`clouds.yaml`).

You can run desired acceptance test using `go test`, e.g.:
```shell
export OS_SUBNET_NAME=examle-subnet
export OS_AVAILABILITY_ZONE=eu-de-01
exprot OS_CLOUD=test # using test cloud from clouds.yaml file
TF_ACC=1 go test -v -timeout=30m ./opentelekomcloud/acceptance/elb/v3/...
```

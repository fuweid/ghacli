# ghacli: Github Action CLI

`ghacli` is a Github Action commandline.

Currently, the Github Action webpage doesn't provide better runtime usage view
of WorkflowRuns. For example, when the CI workflow become slow, we need to find
out which commit slows it down. The runtime usage table on the webpage is not
just the job's runtime, also including the pending time of Action runner. It
is not correct. We need a tool to show the real runtime usage easier.

## Install

```bash
go install github.com/fuweid/ghacli
```

## Examples

To list the workflows in `containerd/containerd` repo.

```bash
$ ghacli --owner containerd --repo containerd workflows list
ID          NAME                                   STATE     PATH
34966       Nightly                                active    .github/workflows/nightly.yml
962330      CI                                     active    .github/workflows/ci.yml
1294966     CodeQL Scan                            active    .github/workflows/codeql.yml
1412342     Containerd Release                     active    .github/workflows/release.yml
9593989     Windows Integration Tests              active    .github/workflows/windows-periodic.yml
13422192    Mirror Test Image                      active    .github/workflows/images.yml
15814172    Windows Periodic Tests                 active    .github/workflows/windows-periodic-trigger.yml
15816108    Build volume test images               active    .github/workflows/build-test-images.yml
16386575    pages-build-deployment                 active    dynamic/pages/pages-build-deployment
28263710    Fuzzing                                active    .github/workflows/fuzz.yml
35386283    Scorecards supply-chain security       active    .github/workflows/scorecards.yml
38631176    Windows Hyper-V Periodic Tests         active    .github/workflows/windows-hyperv-periodic-trigger.yml
38631177    Windows Hyper-V Integration Tests      active    .github/workflows/windows-hyperv-periodic.yml
41916045    Add pull requests to review project    active    .github/workflows/project.yml
43694150    CRI                                    active    .github/workflows/cri.yml
43696258    Check                                  active    .github/workflows/check.yml
```

To list the all workflow runs in `containerd/containerd` repo.

```bash
$ ghacli --owner containerd -repo containerd run list --limit 10
ID          NAME                EVENT         BRANCH           HEAD(MESSAGE)                                                   SHA       STATUS       CREATED
4054514391  CodeQL Scan         pull_request  fixunmount       Make `mount.UnmountRecursive` compatible to `mount.UnmountAll`  eeab0524  success      2023-01-31T13:22:07Z
4054514382  Fuzzing             pull_request  fixunmount       Make `mount.UnmountRecursive` compatible to `mount.UnmountAll`  eeab0524  in_progress  2023-01-31T13:22:07Z
4054514379  CI                  pull_request  fixunmount       Make `mount.UnmountRecursive` compatible to `mount.UnmountAll`  eeab0524  in_progress  2023-01-31T13:22:07Z
4054514377  Containerd Release  pull_request  fixunmount       Make `mount.UnmountRecursive` compatible to `mount.UnmountAll`  eeab0524  success      2023-01-31T13:22:07Z
4054350886  Fuzzing             pull_request  deps/update-nri  go.mod: update github.com/containerd/nri.                       58bd5a09  in_progress  2023-01-31T13:04:16Z
4054350426  Containerd Release  pull_request  deps/update-nri  go.mod: update github.com/containerd/nri.                       58bd5a09  success      2023-01-31T13:04:11Z
4054350424  CodeQL Scan         pull_request  deps/update-nri  go.mod: update github.com/containerd/nri.                       58bd5a09  success      2023-01-31T13:04:11Z
4054350422  CI                  pull_request  deps/update-nri  go.mod: update github.com/containerd/nri.                       58bd5a09  in_progress  2023-01-31T13:04:11Z
4053867294  Containerd Release  pull_request  deps/update-nri  go.mod: update github.com/containerd/nri.                       916072c2  success      2023-01-31T12:03:51Z
4053867293  CI                  pull_request  deps/update-nri  go.mod: update github.com/containerd/nri.                       916072c2  failure      2023-01-31T12:03:51Z
```

To list the workflow runs with a given workflow ID or filename.

```bash
$ ghacli --owner containerd -repo containerd run list --workflow-id ci.yml --limit 10
ID          NAME  EVENT         BRANCH                                    HEAD(MESSAGE)                                                                  SHA       STATUS       CREATED
4054514379  CI    pull_request  fixunmount                                Make `mount.UnmountRecursive` compatible to `mount.UnmountAll`                 eeab0524  in_progress  2023-01-31T13:22:07Z
4054350422  CI    pull_request  deps/update-nri                           go.mod: update github.com/containerd/nri.                                      58bd5a09  in_progress  2023-01-31T13:04:11Z
4053867293  CI    pull_request  deps/update-nri                           go.mod: update github.com/containerd/nri.                                      916072c2  failure      2023-01-31T12:03:51Z
4051851351  CI    push          main                                      Merge pull request #7847 from fangn2/adding-integration-test-to-opentelemetry  e307f879  success      2023-01-31T07:46:04Z
4051166641  CI    pull_request  adding-integration-test-to-opentelemetry  Add integration test for tracing on image pull                                 c46aaa8d  success      2023-01-31T05:46:26Z
4051118713  CI    push          main                                      Merge pull request #7840 from hinshun/feature/mount-subdirectory               287320d4  failure      2023-01-31T05:35:38Z
4047645510  CI    push          release/1.6                               Merge pull request #8030 from AkihiroSuda/cherrypick-8020-1.6                  4335c650  failure      2023-01-30T20:00:05Z
4046940630  CI    pull_request  cherrypick-8020-1.6                       cri: mkdir /etc/cni with 0755, not 0700                                        ae02a24a  success      2023-01-30T18:24:45Z
4046915058  CI    push          main                                      Merge pull request #8020 from AkihiroSuda/mkdir-etc-cni-0755                   ee0e22f0  success      2023-01-30T18:21:39Z
4046189755  CI    pull_request  main                                      doc: fixed windows installation typo                                           d9fd2d1c  success      2023-01-30T16:55:16Z
```

To list the workflow runs with given job names


```bash
$ ghacli --owner containerd -repo containerd run list \
  --workflow-id ci.yml --event push --limit 10 \
  --created 2022-12-01..2023-01-30 
  --jobs "Linux Integration (io.containerd.runc.v2, crun)"
  --jobs "Linux Integration (io.containerd.runc.v2, runc)"

ID          NAME  EVENT  BRANCH       HEAD(MESSAGE)                                                            SHA       STATUS   JOB(Linux Integration (io.containerd.runc.v2, crun))  JOB(Linux Integration (io.containerd.runc.v2, runc))  CREATED
4047645510  CI    push   release/1.6  Merge pull request #8030 from AkihiroSuda/cherrypick-8020-1.6            4335c650  failure  17m2s                                                 17m41s                                                2023-01-30T20:00:05Z
4046915058  CI    push   main         Merge pull request #8020 from AkihiroSuda/mkdir-etc-cni-0755             ee0e22f0  success  19m20s                                                23m29s                                                2023-01-30T18:21:39Z
4046059359  CI    push   main         Merge pull request #8025 from fuweid/enhance-testlog                     0b81378b  failure  19m52s                                                20m39s                                                2023-01-30T16:40:20Z
4044651038  CI    push   main         Merge pull request #8027 from AkihiroSuda/containerd-cgroups-v3          b5bdd6c7  success  19m2s                                                 21m56s                                                2023-01-30T14:06:52Z
4044421035  CI    push   main         Merge pull request #7977 from adisky/update-cni-version                  c3e23618  failure  22m21s                                                20m10s                                                2023-01-30T13:39:18Z
4041147451  CI    push   main         Merge pull request #8026 from AkihiroSuda/otel-1.12.0                    3695f29c  success  24m35s                                                24m58s                                                2023-01-30T06:17:02Z
4041142767  CI    push   main         Merge pull request #8007 from mxpv/events                                f0f6912a  failure  18m19s                                                25m14s                                                2023-01-30T06:16:19Z
4041082421  CI    push   main         Merge pull request #8012 from dcantah/runtime-clarifications             5d482fdd  success  19m33s                                                24m23s                                                2023-01-30T06:06:38Z
4039947297  CI    push   main         Merge pull request #8019 from AkihiroSuda/add-cri-containerd-deprecated  9857b5d1  success  19m11s                                                24m44s                                                2023-01-30T02:21:00Z
4037332262  CI    push   main         Merge pull request #8022 from fuweid/update-release                      967979ef  success  23m30s                                                20m40s                                                2023-01-29T14:58:39Z
```

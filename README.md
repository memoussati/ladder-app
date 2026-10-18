# ladder-app

Tiny Go HTTP service used by the CI/CD exercise ladder. `GET /`, `GET /healthz`, `GET /version`.
`LADDER_FAIL_RATE=0.5` makes half of the `/` requests fail (for canary drills).

## CI security

This repository is public and uses self-hosted runners (the `lab` VM and the `arc-lab` scale
set). Self-hosted runners must never execute code from a fork: a `pull_request` job runs the
PR's version of the workflow and the PR's code, so anyone who opens a pull request would run
arbitrary commands on our machines, with network access to the lab cluster.

- Workflows that use self-hosted runners only trigger on `push` to `main` and `workflow_dispatch`.
- `pull_request_target` is **not** used: it runs with the base branch's workflow and a
  write-capable `GITHUB_TOKEN`, so checking out the PR head inside it is the classic
  "pwn request" mistake.
- Repository setting *Actions → Fork pull request workflows from outside collaborators* is set to
  "Require approval for all outside collaborators".
- Hosted runners (`ubuntu-latest`) handle `pull_request` jobs; they are disposable.
- Secret masking replaces known secret values with `***` in logs, but any transformation
  (`base64`, `rev`, splitting, URL-encoding) defeats it; masking is a convenience against accidents,
  not a security control. Treat anything printed in a job as public.

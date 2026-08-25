# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
go build ./...
go test -race ./...                      # CI always runs with -race
go test -race -run TestSendDir ./...      # single test
go vet ./...
gofmt -l .                                # must print nothing; CI fails on any output
```

Tests start an in-process SSH server (`github.com/hnakamur/go-sshd`) on a random
port and drive it as if it were remote. The machine running the tests therefore
needs both `sh` and a real `scp` binary on `PATH` — the "remote" side is
localhost invoking the system `scp`.

To reproduce the CI cross-build job locally, compile each target the package
claims to support:

```sh
for t in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 \
         freebsd/amd64 freebsd/arm64 netbsd/amd64 openbsd/amd64 \
         dragonfly/amd64 solaris/amd64 illumos/amd64 aix/ppc64; do
  GOOS="${t%%/*}" GOARCH="${t##*/}" go build ./... || echo "FAIL $t"
done
```

## Architecture

This is a client for the legacy SCP protocol. It does **not** implement SCP on
the wire directly. Instead it opens an `ssh.Session` and starts the *remote*
`scp` binary in one of its two hidden modes, then speaks the protocol to that
process's stdin/stdout:

- **source** (`source.go`) — we send; remote runs `scp -t` (sink mode)
- **sink** (`sink.go`) — we receive; remote runs `scp -f` (source mode)

Flags are assembled per call: `p` for permission/time updates, `r` for
recursive, `d` when the destination must be a directory. `SCP.SCPCommand`
overrides the binary, which is how callers reach `sudo scp`.

The layering is three files deep, and `protocol.go` is a *peer* of the session
code rather than something beneath it:

| File | Responsibility |
|---|---|
| `scp.go` | `SCP` handle wrapping an `*ssh.Client`; the caller owns Dial/Close |
| `source.go` / `sink.go` | session lifecycle, remote command construction, directory walking |
| `protocol.go` | the byte protocol: `C`/`D`/`E`/`T` message headers and `\x00`/`\x01`/`\x02` replies |

Each direction exposes three API shapes over the same session machinery —
whole-file (`SendFile`/`ReceiveFile`), caller-supplied stream (`Send`/`Receive`),
and an `io.WriteCloser`/`io.ReadCloser` handle (`SendOpen`/`ReceiveOpen`) — plus
a recursive form (`SendDir`/`ReceiveDir`) that takes an `AcceptFunc` to filter
entries during the walk.

Every reply is read synchronously after each header or body write, so an error
returned mid-transfer usually means the remote `scp` wrote to the protocol
error channel. `protocolError` distinguishes recoverable (`\x01`) from fatal
(`\x02`) — the sink acknowledges the former and continues.

### Platform-specific files

`os.FileInfo` exposes no portable access time, so `newFileInfoFromOS` is
selected per platform. Two different build-constraint mechanisms are in play,
and adding a port means touching the right one **and** the negated tag list in
`fileinfo_generic.go`:

- `fileinfo_linux.go`, `fileinfo_darwin.go`, `fileinfo_windows.go` — selected by
  GOOS filename suffix, no build tag in the file
- `fileinfo_atim.go` (openbsd, dragonfly, solaris), `fileinfo_atimespec.go`
  (freebsd, netbsd) — explicit `//go:build` tags, grouped by how the platform
  spells the `syscall.Stat_t` access-time field
- `fileinfo_generic.go` — everything else (e.g. aix); leaves access time zero
  rather than failing to build

Because a build for one unix says nothing about the others, CI compiles all of
them; do not assume a local `go build` covers this.

`realPath` (`path.go` / `path_windows.go`) rewrites `\` to `/` on Windows: the
SCP wire format is always POSIX regardless of the local separator. Remote paths
pass through `escapeShellArg` because they are interpreted by a remote shell.

## Gotchas

- **The module path is still `github.com/hnakamur/go-scp`**, not `mondoohq`.
  Imports, `go.mod`, and pkg.go.dev all use the upstream path. Renaming it is a
  breaking change for consumers, so leave it alone unless that is the explicit task.
- `go.mod` declares `go 1.17`, but CI tests against Go 1.25 and 1.26.
- This is a fork. `gh pr create` defaults its base to the archived upstream
  repository, so always target the fork explicitly:
  `gh pr create -R mondoohq/go-scp --base main`.
- OpenSSH has deprecated the SCP protocol in favour of SFTP. This library exists
  for hosts where the SFTP subsystem is unavailable; prefer SFTP for new work.

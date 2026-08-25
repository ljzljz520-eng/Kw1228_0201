# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	example.com/energycore/cmd/corevault	[no test files]
ok  	example.com/energycore/internal/catalog	0.021s
?   	example.com/energycore/internal/clock	[no test files]
--- FAIL: Test1228BusinessRegression (0.01s)
    regression_test.go:29: review result retains invalid note: {RecordID:core-000001 Decision:reject Note:withdraw containment note Reviewer:reviewer Revision:3 Sequence:11 WithdrawnNote:withdraw containment note}
FAIL
FAIL	example.com/energycore/internal/flow017	0.062s
ok  	example.com/energycore/internal/httpapi	0.018s
ok  	example.com/energycore/internal/importer	0.002s
?   	example.com/energycore/internal/model	[no test files]
?   	example.com/energycore/internal/report	[no test files]
ok  	example.com/energycore/internal/review	0.020s
ok  	example.com/energycore/internal/store	0.019s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/corevault): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/corevault): exit `0`

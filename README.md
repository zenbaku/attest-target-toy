# attest-target-toy

A deliberately small target for the `attest` runner.

It is written in **Go while the runner is TypeScript**, on purpose: the runner must
not be able to assume its own toolchain.

## What it does

Applies a per-key budget to a list of items. An item is `processed` if its cost
fits in its key's remaining budget, `capped` if it would exceed it, `invalid` if
malformed. A capped item does not consume budget — a later cheaper item for the
same key still fits.

```
go test ./...
go run . --input testdata/items.json --budget 10
```

```
[INFO] toy: run_start items=8 budget=10
[INFO] toy: item_done key=alpha cost=8 status=processed spent=8
[INFO] toy: item_done key=alpha cost=5 status=capped spent=8
[INFO] toy: item_done key=alpha cost=2 status=processed spent=10
[INFO] toy: run_done keys=3 processed=4 capped=2 invalid=2 total_spent=29
```

## Why this shape

It gives the runner every observation kind L0 needs, from one 80-line package:

| Needed | Available here |
|---|---|
| Build step | `go build` |
| Runnable entrypoint | `go run . --input ... --budget ...` |
| Presence assertion | a `status=capped` line appears for an over-budget item |
| Absence assertion | no `status=processed` line for `beta`'s second item |
| Positive control for that absence | `beta`'s first item *is* processed, so the emission path demonstrably fired |
| Field assertions | `key`, `cost`, `status`, `spent` on every line |
| Exit status | `2` on missing `--input`, `1` on unreadable input, `0` otherwise |
| Independently-authored oracle | `budget/budget_test.go` |
| Manufacturable failure | delete the body of `budget.Apply` |

The behaviour is also non-trivial in the way real budget code is: the rule that a
capped item must not consume budget is easy to state, easy to get wrong, and
covered by a test that fails loudly when it is.

## Manufacturing a failure

```
# empty the body of budget.Apply, then:
go test ./...     # fails
```

That is the intended first Contract: the runner deletes `Apply`, the model
reimplements it from the Contract, and the gates judge the result.

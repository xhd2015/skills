# Install Package Tests

Tests for the `install` package's `HandleInstall` function — verifying output messages,
filesystem side effects, flag interactions, and extra file path validation.

## Decision Tree

```
tests/install/
├── behavior/                          # Install behavior and output messages
│   ├── fresh-install/                 # New directory → "Installed skill to:"
│   ├── overwrite-existing/            # Existing directory → "Update skill at"
│   ├── force-overrides-no-override/   # --force overrides --no-override
│   ├── no-override-empty-dir/         # --no-override on empty dir → installs normally
│   └── cursor-flag/                   # --cursor installs to .cursor/skills/<name>
└── extra-files-validation/            # Extra file path validation errors
    ├── dot-path/                      # Path "." → error
    ├── dotdot-path/                   # Path ".." → error
    └── skill-md-path/                 # Path "SKILL.md" → error
```

## Test Index

| # | Test Leaf | Description |
|---|-----------|-------------|
| 1 | behavior/fresh-install | Fresh install to non-existent directory prints "Installed skill to:" |
| 2 | behavior/overwrite-existing | Overwriting existing directory prints "Update skill at" |
| 3 | behavior/force-overrides-no-override | --force --no-override overwrites without confirmation |
| 4 | behavior/no-override-empty-dir | --no-override on empty dir installs normally (no abort) |
| 5 | behavior/cursor-flag | --cursor installs to .cursor/skills/<name> (local, non-global) |
| 6 | extra-files-validation/dot-path | Extra file path "." returns error |
| 7 | extra-files-validation/dotdot-path | Extra file path ".." returns error |
| 8 | extra-files-validation/skill-md-path | Extra file path "SKILL.md" returns error |

## How to Run

```sh
doctest test -v ./tests/install
```

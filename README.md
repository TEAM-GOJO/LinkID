# LinkID-Blockchain

Blockchain system for MediLink

## Makefile Variables (OS)

Common Operating System configurations for compiling the LinkID source code via Makefile. If you want to compile the code on an operating system not listed below, please check out this [list](https://pkg.go.dev/internal/platform) for a list of valid `GOOS` and `GOARCH` combinations.

### Linux (default) 🐧
```
OS=linux
```

### MacOS 🍎
```
OS=darwin
```

### Windows 🪟
```
OS=windows
```

## Makefile Variables (Architecture)

| GOARCH        | Out of the Box | 32-bit | 64-bit |
| :------------ | :------------: | :----: | :----: |
| `386`         | ✅              | ✅      |        |
| `amd64`       | ✅              |        | ✅      |
| `amd64p32`    |                | ✅      |        |
| `arm`         | ✅              | ✅      |        |
| `arm64`       | ✅              |        | ✅      |
| `arm64be`     |                |        | ✅      |
| `armbe`       |                | ✅      |        |
| `loong64`     |                |        | ✅      |
| `mips`        | ✅              | ✅      |        |
| `mips64`      | ✅              |        | ✅      |
| `mips64le`    | ✅              |        | ✅      |
| `mips64p32`   |                | ✅      |        |
| `mips64p32le` |                | ✅      |        |
| `mipsle`      | ✅              | ✅      |        |
| `ppc`         |                | ✅      |        |
| `ppc64`       | ✅              |        | ✅      |
| `ppc64le`     | ✅              |        | ✅      |
| `riscv`       |                | ✅      |        |
| `riscv64`     | ✅              |        | ✅      |
| `s390`        |                | ✅      |        |
| `s390x`       | ✅              |        | ✅      |
| `sparc`       |                | ✅      |        |
| `sparc64`     |                |        | ✅      |
| `wasm`        | ✅              |        | ✅      |


# LinkID-Blockchain [![Go Report Card](https://goreportcard.com/badge/github.com/TEAM-GOJO/LinkID-Blockchain)](https://goreportcard.com/report/github.com/TEAM-GOJO/LinkID-Blockchain)

The LinkID Blockchain system is an AES encryption enhanced blockchain system for protecting medical records through unique 8 digit id and private key pairs.

## Example Medical Record

Medical records are stored in the JSON format and are structured as shown below:

```
{
  "Index": 1,
  "Initials": "JS",
  "Sex": "M", 
  "Age": 62,
  "Height": 173.4,
  "Weight": 78.2,
  "BMI": 26.0,
  "Blood": "O+",
  "Location": "New York City, NY",
  "Prescriptions": ["medication1", "medication2"],
  "Conditions": ["destructive disease"],
  "VisitLogs": [],
  "History": []
}
```

Each record is treated like a "block" inside the blockchain and is given a `SHA256` hash of its own as well as a timestamp and the hash of the block before it to ensure the information is legitamate. A series of these blocks are constructed into a chain, which is encrypted using AES but is also in the JSON format in its decrypted form.

```

```

## Makefile Variables for `OS` (GOOS)

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

## Makefile Variables for `ARCH` (GOARCH)

| GOARCH Variable       | Processor Name   | 32-bit    | 64-bit    |
| :-------------------- | :--------------: | :-------: | :-------: |
| `ARCH=386`            | Intel 386        | ✅        |           |
| `ARCH=amd64`          | AMD64            |           | ✅        |
| `ARCH=amd64p32`       | AMD64 (32-bit)   | ✅        |           |
| `ARCH=arm`            | ARM              | ✅        |           |
| `ARCH=arm64`          | ARM64            |           | ✅        |
| `ARCH=arm64be`        | ARM64 (big-endian)|          | ✅        |
| `ARCH=armbe`          | ARM (big-endian) | ✅        |           |
| `ARCH=loong64`        | Loongson 64-bit  |           | ✅        |
| `ARCH=mips`           | MIPS             | ✅        |           |
| `ARCH=mips64`         | MIPS64           |           | ✅        |
| `ARCH=mips64le`       | MIPS64 (little-endian) |    | ✅        |
| `ARCH=mips64p32`      | MIPS64 (32-bit)  | ✅        |           |
| `ARCH=mips64p32le`    | MIPS64 (32-bit little-endian)| ✅      |   |
| `ARCH=mipsle`         | MIPS (little-endian)| ✅      |          |
| `ARCH=ppc`            | PowerPC          | ✅        |           |
| `ARCH=ppc64`          | PowerPC 64       |           | ✅        |
| `ARCH=ppc64le`        | PowerPC 64 (little-endian) | | ✅        |
| `ARCH=riscv`          | RISC-V           | ✅        |           |
| `ARCH=riscv64`        | RISC-V 64        |           | ✅        |
| `ARCH=s390`           | IBM System/390   | ✅        |           |
| `ARCH=s390x`          | IBM System/390x  |           | ✅        |
| `ARCH=parc`           | SPARC            | ✅        |           |
| `ARCH=sparc64`        | SPARC64          |           | ✅        |
| `ARCH=wasm`           | WebAssembly      | ✅        |           |


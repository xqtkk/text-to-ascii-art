# ASCII Art Web

## Description

A small, responsive Go web application that converts printable ASCII text into banner art. It includes the Standard, Shadow, and Thinkertoy fonts, copy-to-clipboard support, and `.txt` export.

## Authors

- [xqtkk](https://github.com/xqtkk)

## Usage

Requirements: Go 1.22 or newer.

```powershell
git clone https://github.com/xqtkk/text-to-ascii-art.git
cd text-to-ascii-art
go run ./cmd
```

Open <http://localhost:4000>. Enter up to 500 printable ASCII characters, select a font, then generate, copy, or download the result. To use a different port:

```powershell
go run ./cmd -addr :8080
```

## Implementation details

- `GET /` serves the generator.
- `POST /ascii-art` converts the submitted text.
- `POST /export` returns the generated art as a downloadable text file.
- `POST /text-to-ascii-art` remains available as a compatibility route.
- Templates, static files, and banner fonts are embedded in the executable, so it can run from any working directory.
- Font choices are allow-listed; submitted values cannot be used as filesystem paths.
- Request body and input length are limited, and the HTTP server uses read and write timeouts.
- The renderer reads each banner as 95 character records, skips each record's separator line, and combines the eight glyph rows for each input line. Newline characters start a new rendered line.

## Docker

Build and run the container:

```powershell
docker build -t ascii-art-web .
docker run --rm -p 4000:4000 ascii-art-web
```

# analyze-email

A command-line tool for SOC analysts that parses `.eml` files and produces a structured JSON report and a human-readable Markdown report for each email.

Built on top of [emailanalyzer](https://github.com/tonesploit/emailanalyzer).

## What it extracts

- **Envelope**;  From, To, CC, Reply-To, Subject, Date, Message-ID
- **Authentication**;  SPF, DKIM, DMARC, ARC results parsed from headers; live DNS lookups for each record
- **Mail routing**;  full Received header chain with IPv4/IPv6 extraction and per-hop delay calculation
- **Attachments**;  filename, content type, size, and SHA-256 hash for every attachment
- **URLs**;  all unique URLs found in the plain-text and HTML body
- **Security indicators**;  From/Reply-To domain mismatch, Return-Path mismatch, Unicode spoofing in the From display name, encoded subjects, suspicious mailer headers, and analyst notes

## Installation

**From a release**;  download the binary for your platform from the [Releases](https://github.com/tonesploit/emailanalyzer-cli/releases) page and place it somewhere on your `$PATH`.

**From source**;  requires Go 1.23+:

```bash
go install github.com/tonesploit/emailanalyzer-cli@latest
```

## Usage

```
analyze-email [flags] <file.eml> [file.eml ...]

Flags:
  -out string   output directory (default "Analyzed")
```

### Examples

Analyze a single email, writing output to the default `Analyzed/` directory:

```bash
analyze-email suspicious.eml
```

Analyze multiple emails and write results to a custom directory:

```bash
analyze-email -out /cases/2024-04-18 *.eml
```

Each input file produces two output files:

```
Analyzed/
├── suspicious.json   ← structured data (all fields)
└── suspicious.md     ← human-readable Markdown report
```

## Output

### JSON

The full result as a structured JSON object;  suitable for ingestion into a SIEM, piping into `jq`, or archiving alongside a case.

```bash
cat Analyzed/suspicious.json | jq '.authentication.spf'
```

### Markdown

A formatted report covering envelope, authentication results with live DNS records, the full routing hop table, attachments, URLs, and security indicators with analyst notes. Renders cleanly in any Markdown viewer or can be included directly in a report.

## License

MIT

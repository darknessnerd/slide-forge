# MarkItDown — Team Guide

## What It Is

MarkItDown (Microsoft) converts any file format — PDF, Word, PowerPoint, Excel, images, audio, HTML, CSV, ZIP — into clean Markdown that LLMs can read natively. It preserves headings, lists, tables, and links.

## Why Use It

LLMs speak Markdown. Raw PDFs, DOCX, or PPTX files waste tokens on binary noise and lose structure. MarkItDown strips that — Claude reads the result directly with full context preserved.

**Where it fits in this scaffold:**

| Scenario | Without markitdown | With markitdown |
|----------|-------------------|-----------------|
| Feed a spec PDF to Claude | Upload binary, lossy parse | `markitdown spec.pdf > spec.md` → Claude reads clean Markdown |
| Add a Word doc to graphify corpus | graphify skips it | Convert first → graphify indexes it as a document |
| Onboard from a PowerPoint deck | Manually transcribe | `markitdown deck.pptx > deck.md` → queryable in graph |
| Review Excel test data | Attach spreadsheet | `markitdown data.xlsx` → Markdown table Claude can reason over |
| Audit a legacy HTML doc | Wall of tags | `markitdown page.html` → clean structure |

## Install (one-time, per machine)

```bash
# Full install — all file types
pip install 'markitdown[all]'

# Minimal — common formats only (PDF, DOCX, PPTX, XLSX)
pip install 'markitdown[pdf,docx,pptx,xlsx]'

# macOS externally-managed Python
pipx install 'markitdown[all]'
```

**Verify:**
```bash
markitdown --version
```

## Usage

### CLI

```bash
# Convert to stdout
markitdown path/to/file.pdf

# Convert to file
markitdown path/to/file.pdf -o output.md

# Pipe
cat report.pdf | markitdown > report.md
```

### Python API

```python
from markitdown import MarkItDown

md = MarkItDown()
result = md.convert("spec.docx")
print(result.text_content)
```

### With LLM for image descriptions (optional)

```python
from markitdown import MarkItDown
from anthropic import Anthropic  # or OpenAI

md = MarkItDown(llm_client=Anthropic(), llm_model="claude-opus-4-8")
result = md.convert("diagram.png")
```

## How to Use It in This Scaffold

### 1. Feed docs to Claude

Convert before attaching to any Claude Code session:

```bash
markitdown docs/spec.pdf > raw/spec.md
# Then reference raw/spec.md in your prompt
```

### 2. Add converted docs to graphify corpus

graphify indexes `.md` files natively. Convert first, then run graphify:

```bash
mkdir -p raw
markitdown external-docs/architecture.pdf -o raw/architecture.md
markitdown external-docs/adr-001.docx -o raw/adr-001.md
/graphify .          # or /graphify . --update if graph already exists
```

Commit `raw/` to git so teammates share the converted docs without needing the originals.

### 3. Batch convert a folder

```bash
for f in docs/originals/*.pdf; do
  markitdown "$f" -o "raw/$(basename "${f%.pdf}").md"
done
```

### 4. CI pipeline (optional)

Add a conversion step before any doc-dependent job:

```bash
# In your CI workflow
pip install 'markitdown[pdf,docx]' -q
markitdown specs/api-contract.docx -o raw/api-contract.md
```

## Security

markitdown executes with the privileges of the current process. For untrusted input:

- Use `convert_local()` instead of `convert()` — blocks remote URI fetching
- Never run against untrusted user-uploaded files without sandboxing
- Restrict to known file extensions before passing to the converter

## Supported Formats

| Format | Notes |
|--------|-------|
| PDF | Built-in parser; Azure Document Intelligence for higher quality |
| Word (`.docx`) | Full structure preserved |
| PowerPoint (`.pptx`) | Slides → Markdown sections |
| Excel (`.xlsx`) | Sheets → Markdown tables |
| HTML | Strips tags, preserves structure |
| CSV / JSON / XML | Converted to readable Markdown |
| Images (`.png`, `.jpg`) | EXIF metadata; LLM client needed for descriptions |
| Audio | Transcription via LLM client |
| ZIP | Recursively converts contents |
| YouTube URLs | Transcript extraction |

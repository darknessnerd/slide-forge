<!--
  COMMAND: markitdown
  Trigger:  /markitdown <path>
  How:      You type it. Claude runs markitdown on the given file or glob.
  Input:    Required — $ARGUMENTS is the file path or glob pattern.
            Examples:
              /markitdown docs/spec.pdf
              /markitdown external/*.docx
              /markitdown https://example.com/page.html
  Output:   Converted .md file(s) written to raw/ alongside the source,
            then summarised in chat.
-->

Convert `$ARGUMENTS` to Markdown using markitdown and save to `raw/`.

Steps:
1. Ensure markitdown is installed: `pip install 'markitdown[all]' -q 2>/dev/null || pipx install 'markitdown[all]'`
2. Create `raw/` if it doesn't exist: `mkdir -p raw`
3. For each file matching `$ARGUMENTS`:
   - Run `markitdown <file> -o raw/<stem>.md`
   - Report the output path and word count
4. If a `graphify-out/graph.json` exists, offer to run `/graphify . --update` to index the new files

Security: if the path looks user-supplied or untrusted, use `convert_local()` via Python API instead of the CLI (blocks remote URI fetching).

See `docs/markitdown.md` for full usage guide.

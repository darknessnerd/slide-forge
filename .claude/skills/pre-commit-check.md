# Skill: pre-commit-check

<!--
  SKILL vs COMMAND: Claude triggers this itself — you never type /pre-commit-check.
  The Trigger line below is what Claude reads to decide when to apply this skill.
  Customize: Swap `go vet` / `go test` for your own toolchain commands.
             Add linter steps (golangci-lint, staticcheck, etc.) as needed.
             Change the commit message format rule if your team uses a different convention.
-->

**Trigger:** Claude is about to suggest a `git commit` command.

Apply automatically before proposing any commit.

## Steps

1. Run `git diff --staged` — confirm the staged diff matches what was discussed
2. Verify no `.env` or secret files are in the staged set
3. Run `go vet ./...` — do not proceed if errors
4. Run `go test ./...` — do not proceed if any tests are red
5. Verify the proposed commit message follows Conventional Commits: `type(scope): summary`
   - Allowed types: `feat` / `fix` / `chore` / `docs` / `test` / `refactor`
6. If any check fails: report the finding, do **NOT** suggest the commit command

<!--
  COMMAND: standup
  Trigger:  /standup
  How:      You type it. Claude does not trigger this automatically.
  Input:    None — runs against current branch and git config user.name.
  Customize: Change the time window ("yesterday" → "last 3 days"),
             or add a step to pull open PRs via the GitHub MCP connection.
-->

Generate a standup summary from recent git activity.

Run:
```bash
git log --since="yesterday" --author="$(git config user.name)" --oneline
```

Then summarize in standup format:
- **Done:** (commits from yesterday)
- **Today:** (open PRs or TODO comments in staged changes)
- **Blocked:** (any failing tests or unresolved conflicts)

Keep it under 10 lines.

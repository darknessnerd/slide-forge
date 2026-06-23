# Skill: explain-error

<!--
  SKILL vs COMMAND: Claude triggers this itself — you never type /explain-error.
  The Trigger line below is what Claude reads to decide when to apply this skill.
  Customize: Change the "3 files" threshold if your team prefers a different
             escalation boundary. Add language-specific error pattern notes if helpful.
-->

**Trigger:** A build or test command returns a non-zero exit code.

Apply automatically — do not wait for the user to ask "what went wrong?"

## Steps

1. Quote the exact error message — do not paraphrase
2. Identify the location: file, line number, symbol name
3. Explain the root cause in one sentence
4. Show the minimal fix (change only what is broken)
5. If the fix would touch more than 3 files: describe the change and ask the user to confirm before applying

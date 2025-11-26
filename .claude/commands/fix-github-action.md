---
allowed-tools: Bash(gh:*), Bash(npm:*), Bash(git:*), Read, Edit, Write, Glob, Grep
description: Fetch and fix the latest failed GitHub Actions
argument-hint: [run-id]
---

## Context

Current branch: !`git branch --show-current`

Recent workflow runs:
!`gh run list --limit 5 2>/dev/null || echo "Cannot fetch workflow list, please ensure gh is logged in"`

## Task

Please complete the following tasks:

1. **Fetch failure information**
   - If run-id is provided ($ARGUMENTS), use that ID
   - Otherwise, automatically find the most recent failed workflow run
   - Use `gh run view <run-id> --log-failed` to get failure logs

2. **Analyze error cause**
   - Carefully read error logs and identify root cause
   - Common error types:
     - Dependency conflicts (npm peer dependency)
     - Build failures (TypeScript/ESLint errors)
     - Test failures
     - Configuration issues

3. **Fix the problem**
   - Apply appropriate fixes based on error type
   - Validate fix locally (run build/test)
   - Ensure no new issues are introduced

4. **Commit and verify**
   - Commit the fix (follow project commit conventions)
   - Push to remote
   - Use `gh run watch` to monitor new CI run until completion

## Notes

- If automatic fix is not possible, explain the issue and provide suggestions
- Prefer maintaining existing architecture and implementation approaches
- If CI still fails, continue analyzing new error logs

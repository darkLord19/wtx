# Common Workflows

This guide demonstrates practical workflows using wtx for everyday development scenarios.

## Table of Contents

1. [Quick Feature Development](#quick-feature-development)
2. [Code Review Workflow](#code-review-workflow)
3. [Hotfix Emergency Response](#hotfix-emergency-response)
4. [Parallel Development](#parallel-development)
5. [Microservices Management](#microservices-management)
6. [Long-Running Experiments](#long-running-experiments)
7. [Conference Talk Preparation](#conference-talk-preparation)

---

## Quick Feature Development

**Scenario**: Start a new feature, get interrupted by a bug, then return to the feature.

```bash
# Start new feature
wtx add feat-user-login

# Work on login feature...
# (commits, tests, etc.)

# Urgent bug report arrives
wtx add hotfix-payment-error --from main

# Fix bug in hotfix worktree
cd /path/to/hotfix-payment-error
# Make fixes
git commit -am "fix: payment processing error"
git push origin hotfix-payment-error

# Create PR, get it merged

# Back to feature work
wtx  # Select feat-user-login
# Continue where you left off
```

**Key Benefits**:
- No stashing needed
- Clean separation of concerns
- Fast context switching (<2 seconds)

---

## Code Review Workflow

**Scenario**: Review a colleague's PR without disrupting your current work.

```bash
# Currently working on feature-dashboard
# Colleague asks for PR review on pull/456

# Create review worktree from PR branch
wtx add review-pr-456 --from origin/feature-new-api

# Review in your editor
wtx open review-pr-456

# Make review comments, test locally
cd /path/to/review-pr-456
npm install
npm test

# Leave comments in GitHub

# Done reviewing - clean up
wtx rm review-pr-456

# Back to your work
wtx open feature-dashboard
```

**Pro Tips**:
- Name reviews clearly: `review-pr-<number>`
- Prune review worktrees regularly
- Use `wtx status review-pr-456` to check testing setup

---

## Hotfix Emergency Response

**Scenario**: Production is down, need immediate fix.

```bash
# Production error reported
# Create hotfix from production branch
wtx add hotfix-critical-auth --from production

# Opens immediately in editor
# Make the fix
git commit -am "hotfix: fix authentication bypass"

# Test locally
npm test

# Push and deploy
git push origin hotfix-critical-auth

# Create PR for review (after the fact)

# Once deployed and verified
wtx rm hotfix-critical-auth
```

**Emergency Checklist**:
1. Create worktree from production branch
2. Make minimal fix
3. Test thoroughly
4. Push and deploy
5. Create PR for documentation
6. Clean up after merge

---

## Parallel Development

**Scenario**: Work on frontend and backend simultaneously.

```bash
# Frontend work
wtx add frontend-redesign

# Backend work
wtx add backend-api-v2

# Terminal 1: Frontend
cd /path/to/frontend-redesign
npm run dev  # Runs on :3000

# Terminal 2: Backend
cd /path/to/backend-api-v2
npm run dev  # Runs on :8080

# Terminal 3: Main development
wtx  # Switch between as needed
```

**Setup**:
```bash
# Configure different ports in metadata
# Edit .git/wtx-meta.json
{
  "worktrees": {
    "frontend-redesign": {
      "ports": [3000],
      "dev_command": "npm run dev"
    },
    "backend-api-v2": {
      "ports": [8080, 5432],
      "dev_command": "npm run dev:api"
    }
  }
}
```

---

## Microservices Management

**Scenario**: Multiple services in one repository (monorepo).

```bash
# Service 1: Auth Service
wtx add auth-service-upgrade

# Service 2: Payment Service
wtx add payment-service-refactor

# Service 3: Notification Service
wtx add notification-service-websockets

# Start all services
./scripts/start-all-services.sh

# Work on specific service
wtx open auth-service-upgrade
```

**Directory Structure**:
```
repo/
├── .git/
├── services/
│   ├── auth/
│   ├── payment/
│   └── notification/
└── ../worktrees/
    ├── auth-service-upgrade/
    ├── payment-service-refactor/
    └── notification-service-websockets/
```

**Benefits**:
- Each service has independent dependencies
- Run different versions simultaneously
- Test service interactions locally

---

## Long-Running Experiments

**Scenario**: Try a radical refactor without committing to main branch.

```bash
# Start experiment
wtx add experiment-graphql-migration --from develop

# Work on experiment over weeks
# Make commits in experiment worktree
cd /path/to/experiment-graphql-migration
git commit -am "wip: initial GraphQL schema"
git commit -am "wip: implement resolvers"
# ... many commits ...

# Meanwhile, continue normal work
wtx add feature-user-profile --from main

# Experiment succeeds!
cd /path/to/experiment-graphql-migration
git push origin experiment-graphql-migration
# Create PR

# Or experiment fails
wtx rm experiment-graphql-migration --force
```

**Tips**:
- Keep experiments in separate worktrees
- Commit frequently (experiment commits are cheap)
- Use descriptive names: `experiment-*`, `poc-*`, `spike-*`
- Clean up failed experiments without guilt

---

## Conference Talk Preparation

**Scenario**: Prepare demo for a conference while maintaining current work.

```bash
# Create demo environment
wtx add conf-demo-june-2024 --from main

# Prepare demo
cd /path/to/conf-demo-june-2024

# Simplify for demo
rm -rf complex-features/
# Add demo data
cp demo-data/* data/
git commit -am "demo: prepare conference demo"

# Practice demo multiple times
# Make tweaks as needed

# During conference
wtx open conf-demo-june-2024
# Run demo

# After conference
# Archive or delete
wtx rm conf-demo-june-2024
```

**Demo Preparation Checklist**:
- [ ] Remove complex/irrelevant features
- [ ] Add realistic demo data
- [ ] Test on same OS as presentation laptop
- [ ] Practice timing (know your demo duration)
- [ ] Have fallback (screenshots/video)

---

## Advanced: Multi-Repo Coordination

**Scenario**: Work on related changes across multiple repositories.

```bash
# Frontend repo
cd ~/projects/my-app-frontend
wtx add feature-new-dashboard

# Backend repo  
cd ~/projects/my-app-backend
wtx add feature-dashboard-api

# Common library repo
cd ~/projects/my-app-common
wtx add feature-dashboard-types

# Work across all three
# Terminal 1: Frontend
cd ~/projects/my-app-frontend/worktrees/feature-new-dashboard
npm run dev

# Terminal 2: Backend
cd ~/projects/my-app-backend/worktrees/feature-dashboard-api
npm run dev

# Terminal 3: Types
cd ~/projects/my-app-common/worktrees/feature-dashboard-types
npm run build --watch
```

---

## Best Practices

### Naming Conventions

**Feature branches**:
- `feat-<description>` - New features
- `fix-<description>` - Bug fixes
- `refactor-<description>` - Code refactoring
- `docs-<description>` - Documentation updates

**Temporary work**:
- `review-pr-<number>` - PR reviews
- `experiment-<idea>` - Experiments
- `hotfix-<issue>` - Hotfixes
- `demo-<event>` - Demos/presentations

### Cleanup Schedule

**Daily**:
- Remove completed review worktrees

**Weekly**:
- Run `wtx prune --days 7` to clean stale worktrees
- Check `wtx list` for old experiments

**Monthly**:
- Review all worktrees: `wtx list`
- Archive or delete old branches
- Run `wtx prune --days 30`

### Performance Tips

1. **Keep it under 20**: Maintain <20 active worktrees
2. **Use direct commands**: For known worktrees, use `wtx open <name>`
3. **Prune regularly**: Set up a weekly cron job
4. **Clean branches**: Delete merged branches promptly

### Team Coordination

**Share branch naming conventions**:
```bash
# Team convention
<type>/<ticket>-<description>

# Examples
feat/PROJ-123-user-authentication
fix/PROJ-456-payment-bug
refactor/PROJ-789-api-cleanup
```

**Don't share**:
- wtx metadata (`.git/wtx-meta.json` is local)
- Temporary worktrees
- Experiment branches (unless ready for review)

---

## Troubleshooting Common Workflows

### "Too many open files" error

**Solution**: Increase file descriptor limit:
```bash
# macOS/Linux
ulimit -n 4096

# Make permanent (macOS)
echo "ulimit -n 4096" >> ~/.zshrc
```

### Disk space issues

**Solution**: Clean up worktrees:
```bash
# Check disk usage
du -sh ../worktrees/*

# Remove old worktrees
wtx prune --days 14

# Or manually
wtx rm old-feature-1
wtx rm old-feature-2
```

### Forgot which worktree has my work

**Solution**: Use metadata to find it:
```bash
# List with status
wtx list

# Check recent worktrees
# (Check .git/wtx-meta.json for last_opened)

# Or search git logs
git log --all --oneline | grep "your commit message"
```

---

## Next Steps

- Read the [FAQ](FAQ.md) for common questions
- Check [README.md](../README.md) for complete reference
- See [CONTRIBUTING.md](../CONTRIBUTING.md) to improve wtx

**Questions?** Open a discussion on GitHub!

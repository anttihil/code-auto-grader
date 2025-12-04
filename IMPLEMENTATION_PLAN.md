# Implementation Plan

## Current Status

**Dogfooding MVP: ✅ COMPLETE**

Working system with:
- Rust tech tree (3 nodes, 9 exercises)
- Full submission flow with Docker execution + AI feedback
- Mastery tracking (60% tests + 40% AI, 80% threshold)
- Frontend integrated with all database APIs

**Next Step:** Personal Rust learning validation → then decide on multi-user MVP.

---

## Completed (Reference Only)

| Phase | Component | Status |
|-------|-----------|--------|
| Week 1 | PostgreSQL schema, migrations, repositories | ✅ |
| Week 2 | Rust content (9 exercises), seeding system | ✅ |
| Week 3 | APIs, code execution, AI eval, mastery tracking | ✅ |
| Week 3 | Frontend: tech tree, nodes, exercises, history | ✅ |

---

## Multi-User MVP (Months 1-3)

**Prerequisite:** Dogfooding validation confirms learning loop works.

### Authentication & Users
- Email/password + OAuth (Google, GitHub) via `goth`
- JWT sessions, password reset, email verification
- User profiles with progress tracking
- Rate limiting on auth endpoints

### Additional Exercise Types
- **Conceptual**: Plain text submissions, no code execution
- **Debugging**: Broken code to fix, diff viewer in Monaco

### Admin Tools
- Manual node unlock controls
- Content preview/edit CRUD UI

### Polish & Beta
- 50 beta users from r/learnprogramming, HN
- In-app feedback survey
- Performance optimization (caching, lazy loading)

**Success Criteria:**
- 20%+ users complete all nodes
- 70%+ report AI feedback helpful
- <5% critical bugs

---

## Phase 2 (Months 4-6): Automation

### LLM Progression Agent
- Analyze mastery scores, time spent, error patterns
- Recommend next node with rationale
- "What's Next?" modal after node completion

### New Exercise Types
- **Project-based**: Multi-step tasks, GitHub repo evaluation
- **System Design**: Excalidraw diagrams + written explanation

### Community Features
- Exercise voting (upvote/downvote)
- User-generated content with admin approval

---

## Phase 3-4 (Months 7-12): Scale

- **Content Generation**: Upload PDFs/ePubs → LLM generates nodes + exercises
- **Adaptive Difficulty**: Skip easy variants, offer hints when stuck
- **Analytics**: Learning velocity, struggle points, churn prediction
- **Motivation**: Achievements, streaks, personalized encouragement

---

## Technical Debt Checklist

- [ ] 70%+ test coverage by Phase 2
- [ ] Extract progression logic to separate service
- [ ] Quarterly security audits
- [ ] LLM cost monitoring with budget alerts

---

## Key Metrics

| Category | Metric | Target |
|----------|--------|--------|
| Engagement | Exercises/week | 5+ |
| Learning | Nodes mastered/month | 2+ |
| Retention | Monthly churn | <5% |
| Content | Exercise completion rate | 80%+ |

# Implementation Plan

This document outlines the technical development roadmap for Code Auto-Grader, organized by MVP milestones and post-MVP phases.

## Development Principles

1. **Ship early, iterate fast**: Launch with one excellent tech tree path
2. **Validate before building**: Test core assumptions with beta users
3. **Avoid perfectionism**: "Embarrassingly simple" MVP > half-baked features
4. **Instrument everything**: Track user behavior to guide decisions

---

## Pre-MVP: Validation ✅ COMPLETED

**Goal**: Confirm the product vision resonates with target users

### Validation Approach: Dogfooding
Instead of user interviews, validated through **personal use as primary user**:
- [x] Core AI feedback quality tested and confirmed effective
- [x] Code execution system (JavaScript, Python, Rust) working
- [x] Validation plan: Learn Rust using the platform
  - If the app effectively teaches Rust, the core loop works
  - Self-learning dev = authentic stress test of pedagogy
  - Fast iteration without external dependencies

### Success Criteria Met
- ✅ AI feedback feels high-quality and helpful
- ✅ Code execution is reliable and secure
- ✅ Clear personal need identified (learning Rust)

**Next Step**: Build minimal dogfooding MVP for personal use before scaling.

---

## Dogfooding MVP (Weeks 1-3): Minimal Learning Loop

**Goal**: Build streamlined version for personal Rust learning validation

**Philosophy**: Skip multi-user features, focus on core learning mechanics. Hard-code single user. Validate loop works before adding complexity.

### Week 1: Database Foundation

**Deliverables:**
- [ ] PostgreSQL schema design
  - Tech trees, nodes, node edges (DAG structure)
  - Exercises (coding only for MVP - defer other types)
  - User progress (single hard-coded user initially)
  - Submissions (code, execution results, AI feedback, timestamps)
- [ ] Migration scripts (using `golang-migrate`)
- [ ] Go models/structs matching schema
- [ ] Basic CRUD operations with `sqlx`

**Technical Decisions:**
- Single hard-coded user (skip auth for now)
- Focus on coding exercises only (defer debugging, projects, conceptual)
- Store exercise content as JSONB (flexible for future types)
- Index on `node_id`, `exercise_id` for query performance

**Deferred to Later:**
- User authentication/authorization
- Multiple user support
- OAuth integrations
- Admin dashboards

---

### Week 2: Content Creation - Rust Learning Path ✅ MOSTLY COMPLETE

**Deliverables:**
- [x] Design Rust tech tree structure:
  - 7 nodes covering fundamentals (documented in `RUST_TECH_TREE.md`)
  - Clear prerequisite relationships (DAG)
- [x] Write 9 quality Rust exercises (Nodes 1-3 complete):
  - **Node 1 (Rust Basics)**: 3 exercises - Variables, Functions, Control Flow
  - **Node 2 (Ownership & Borrowing)**: 3 exercises - Ownership, Borrowing, Rules
  - **Node 3 (Structs & Enums)**: 3 exercises - Structs, Methods, Enums
  - Progressive difficulty within each node (Easy → Medium → Hard)
  - Comprehensive test cases for automated evaluation
  - Reference solutions included
  - AI feedback rubrics and prompts
- [x] Create learning materials:
  - Detailed markdown for each exercise with explanations
  - Code examples and hints
  - Common mistakes and troubleshooting
  - Links to concepts where helpful
- [ ] Write 6 more exercises for Nodes 4-5 (Collections, Error Handling) - *DEFERRED*
- [ ] Seed database with Rust content (blocked by database schema)

**Quality Bar:**
- Exercises test understanding, not just syntax memorization
- Difficulty curve feels appropriate (not too easy/hard)
- Materials explain "why" not just "how"

---

### Week 3: Learning Loop Integration

**Deliverables:**
- [ ] Mastery tracking logic:
  - Calculate scores from test results + AI feedback
  - 80%+ threshold to count toward mastery
  - Node marked complete when all exercises pass threshold
- [ ] Tech tree navigation UI:
  - Display available nodes (based on prerequisites)
  - Show current progress per node (X/Y exercises complete)
  - Highlight next recommended node
  - Simple "Next Exercise" button for linear flow
- [ ] Exercise history view:
  - Show previous attempts for current exercise
  - Display score progression over time
  - Include AI feedback from past submissions
- [ ] Enhanced AI feedback context:
  - Fetch user's previous submissions for exercise
  - Include recent progress across nodes
  - Tailor feedback to observed patterns
- [ ] Update API endpoints:
  - `GET /api/tech-tree/:treeId` - Fetch tree structure + progress
  - `GET /api/node/:nodeId` - Get node details + exercises
  - `GET /api/exercise/:exerciseId` - Fetch exercise + submission history
  - `POST /api/submission/:exerciseId` - Submit + evaluate (existing endpoint enhanced)

**Deferred Features:**
- Manual progression controls (admin dashboard)
- Automated LLM progression agent
- Rich text editors for conceptual exercises
- Advanced visualizations

---

### Validation: Personal Rust Learning Test

**Goal**: Complete first 1-2 nodes of Rust tech tree using the platform

**Success Criteria:**
- [ ] Can navigate between exercises smoothly
- [ ] AI feedback actually helps improve solutions
- [ ] Progress feels visible and motivating
- [ ] Mastery tracking accurately reflects understanding
- [ ] No data loss (progress persists across sessions)

**Decision Point:**
- ✅ If learning feels effective → proceed to full MVP (add multi-user support)
- ❌ If learning feels broken → iterate on core loop before scaling

---

## MVP Phase (Months 1-3): Multi-User Launch

**Goal**: Single tech tree path with functional mastery-based progression

**Prerequisites**: Dogfooding MVP validation successful (core learning loop works)

**Key Changes from Dogfooding MVP**: Add authentication, multi-user support, admin tools, additional exercise types

### Week 1-2: Database & Data Model

**Deliverables:**
- [ ] PostgreSQL schema design
  - Tech trees, nodes, node edges (DAG structure)
  - Exercises (5 types: coding, debugging, project, conceptual, system design)
  - User progress (node status, mastery scores)
  - Submissions (code/text, execution results, AI feedback, timestamps)
- [ ] Migration scripts (using `golang-migrate` or similar)
- [ ] Go models/structs matching schema
- [ ] Basic CRUD operations for all entities

**Technical Decisions:**
- Use `sqlx` for Go database access (simplicity over ORM)
- Store exercise content as JSONB (flexible schema for different types)
- Index on `user_id`, `node_id`, `exercise_id` for query performance

---

### Week 3-4: Exercise System (MVP Subset)

**Deliverables:**
- [ ] Support 2 exercise types only: **Coding** + **Conceptual**
- [ ] Coding exercises:
  - Extend existing code execution to run test cases
  - Store test results (pass/fail, stdout, stderr)
  - Return results to frontend
- [ ] Conceptual exercises:
  - Accept plain text submissions
  - No special tooling (rich text editor in post-MVP)
- [ ] Update API endpoints:
  - `GET /api/exercise/:id` - Fetch exercise details
  - `POST /api/submission/:exerciseId` - Submit solution

**Defer to Post-MVP:**
- Debugging, project, system design types
- Code diff viewer, diagram tools

---

### Week 5-6: AI Assessment & Mastery

**Deliverables:**
- [ ] Contextual LLM feedback:
  - Fetch user's previous submissions for this exercise
  - Include user's overall progress level in prompt
  - Return feedback focused on improvement (not just "correct/incorrect")
- [ ] Mastery threshold logic:
  - 3 exercises per node minimum
  - 80%+ score on each to count toward mastery
  - Node marked "complete" when threshold met
- [ ] Manual progression (admin-controlled):
  - Admin dashboard to unlock next node for user
  - Defer automated LLM progression agent to Phase 2

**Prompt Engineering:**
```
You are evaluating a {exercise_type} exercise for a student at {user_level} level.

Exercise: {exercise_description}
Student's submission: {submission_code_or_text}
Test results: {test_results} (if coding)
Previous attempts: {previous_submissions}

Provide feedback considering:
1. Correctness (does it solve the problem?)
2. Progress (is this better than previous attempts?)
3. Learning (what should they improve next?)

Return JSON:
{
  "score": 0-100,
  "feedback": "markdown text",
  "encouragement": "motivational message"
}
```

---

### Week 7-8: First Tech Tree (Content Creation)

**Deliverables:**
- [ ] Create "Python Fundamentals" tech tree:
  - 5 nodes: Variables, Control Flow, Functions, Data Structures, File I/O
  - 3 exercises per node (15 total)
  - Materials: Markdown with code snippets, examples
- [ ] Seed database with content
- [ ] Admin tool to preview/edit content (simple CRUD UI)

**Quality Bar:**
- Materials explain concepts clearly (tested with 2-3 users)
- Exercises progressively harder within each node
- Solutions exist for all exercises (reference implementations)

---

### Week 9-10: User Accounts & Progress Persistence

**Deliverables:**
- [ ] Authentication system:
  - Email/password signup and login (use `bcrypt` for hashing)
  - OAuth (Google, GitHub) via `goth` library
  - JWT tokens for session management
- [ ] User profile:
  - Display name, email, join date
  - Current node, total exercises completed
- [ ] Progress tracking:
  - Save all submissions to database
  - Show history of attempts per exercise
  - Calculate mastery percentage per node
- [ ] API endpoints:
  - `POST /api/auth/signup`, `/api/auth/login`
  - `GET /api/progress/user` - Current status
  - `GET /api/progress/history/:nodeId` - Exercise attempts

**Security:**
- Rate limit auth endpoints (10 requests/min)
- Email verification (send confirmation link)
- Password reset flow

---

### Week 11-12: Polish & Beta Testing

**Deliverables:**
- [ ] UI improvements:
  - Tech tree visualization (simple list/cards for MVP)
  - Exercise viewer with submission form
  - Feedback display (markdown rendering)
  - Progress indicators (% per node)
- [ ] Performance optimization:
  - Cache exercise content (Redis optional)
  - Optimize DB queries (add indexes as needed)
  - Frontend lazy loading
- [ ] Beta launch:
  - Recruit 50 users from r/learnprogramming, HN
  - Monitor usage (Google Analytics + custom events)
  - Collect feedback via in-app survey

**Success Criteria (MVP):**
- 50 users complete ≥1 node (20%+ complete all 5 nodes)
- Avg 3+ exercises/week per active user
- 70%+ report AI feedback is helpful
- <5% critical bugs (auth, data loss, crashes)

---

## Phase 2 (Months 4-6): Automation & Exercise Expansion

**Goal**: Remove manual bottlenecks, add 3 more exercise types

### Month 4: LLM Progression Agent

**Deliverables:**
- [ ] Automated next-node recommendation:
  - LLM analyzes user's mastery scores across nodes
  - Considers time spent, error patterns
  - Recommends next node (or review if struggling)
- [ ] Progression endpoint:
  - `POST /api/progression/recommend` - Returns next node ID + rationale
- [ ] User flow:
  - After completing node, show "What's Next?" modal
  - Display LLM's recommendation + reasoning
  - User can accept or choose manually

**Prompt Example:**
```
Analyze this student's progress and recommend the next learning node.

Completed nodes:
- Variables (3/3 exercises, avg 85%)
- Control Flow (3/3 exercises, avg 78%, struggled with nested loops)

Available next nodes:
- Functions (requires Variables)
- Data Structures (requires Variables + Control Flow)

Recommend the best next step and explain why.
```

---

### Month 5: New Exercise Types

**Deliverables:**
- [ ] **Debugging exercises**:
  - Provide broken code with hidden bugs
  - Student fixes and submits corrected version
  - Show code diff in UI (Monaco editor diff mode)
- [ ] **Project-based exercises**:
  - Multi-step tasks (e.g., "Build a CLI tool")
  - Allow file uploads or GitHub repo links
  - Evaluate final product + process (commits, structure)
- [ ] **System design exercises**:
  - Diagram editor (use library like `excalidraw` or `mermaid`)
  - Submit diagram + written explanation
  - LLM evaluates architecture choices

**Technical Notes:**
- Use Monaco editor for all code interfaces (already supports diff)
- Excalidraw component for diagrams (embeddable React version)
- Store diagrams as JSON or SVG in database

---

### Month 6: Community Features

**Deliverables:**
- [ ] Voting on exercises and materials:
  - Users upvote/downvote after completion
  - Flag inappropriate or low-quality content
- [ ] User-generated content (basic):
  - Submit custom exercises for review
  - Admin approval queue before publishing
- [ ] Leaderboards (optional):
  - Top learners per tech tree (by mastery %, not speed)
  - Display anonymously or opt-in

---

## Phase 3 (Months 7-9): Content Generation at Scale

**Goal**: Enable rapid creation of new tech trees

### Month 7: Material Generator

**Deliverables:**
- [ ] Upload PDFs, ePubs, or plain text
- [ ] LLM extracts key concepts and generates:
  - Node structure (topic breakdown)
  - Markdown lessons with examples
  - Suggested exercise prompts
- [ ] Human review before publishing (quality gate)

---

### Month 8: Exercise Generator

**Deliverables:**
- [ ] From lesson material → generate diverse problems
- [ ] Support all 5 exercise types
- [ ] Validate generated test cases (run against reference solution)

---

### Month 9: Quality Pipeline

**Deliverables:**
- [ ] Auto-review system:
  - LLM checks generated content for errors
  - Flag overly easy/hard exercises
- [ ] Community voting integration:
  - High-rated user content auto-approved
  - Low-rated flagged for deletion/improvement

---

## Phase 4 (Months 10-12): Intelligence & Scale

**Goal**: Adaptive difficulty, analytics, motivation

### Month 10: Advanced Analytics

**Deliverables:**
- [ ] Learning velocity metrics (exercises/week, time to mastery)
- [ ] Struggle point detection (where users get stuck)
- [ ] Predictive models (likelihood to complete node)

---

### Month 11: Adaptive Difficulty

**Deliverables:**
- [ ] Adjust exercise complexity per user:
  - If user breezes through → skip easier variants
  - If stuck → offer hints, simpler problems
- [ ] Dynamic hint system (LLM-generated on demand)

---

### Month 12: Motivational System

**Deliverables:**
- [ ] Achievements/badges (complete 10 nodes, 7-day streak)
- [ ] Personalized encouragement messages (LLM-generated)
- [ ] Progress reports ("You've improved 40% in debugging!")

---

## Technical Debt & Maintenance

**Ongoing:**
- Refactor as code grows (split monolithic files)
- Add unit tests (target 70%+ coverage by Phase 2)
- Monitor LLM costs (set budget alerts)
- Security audits (quarterly reviews)

**Major Refactors:**
- Phase 2: Extract progression logic into separate service
- Phase 3: Microservices if scaling requires (avoid premature split)
- Phase 4: Move analytics to separate DB (ClickHouse, BigQuery)

---

## Risk Mitigation Strategies

### Technical Risks

**1. LLM Quality Issues**
- Mitigation: A/B test prompts, manual review samples, user feedback flags
- Fallback: Human reviewers for escalated cases

**2. Cost Overruns**
- Mitigation: Cache aggressively, use tiered models, cap re-submissions
- Monitoring: Daily cost reports, budget alerts at 80%

**3. Security Vulnerabilities**
- Mitigation: Penetration testing before public launch, bug bounty program
- Monitoring: Automated security scans (Dependabot, Snyk)

### Product Risks

**4. Content Quality Bottleneck**
- Mitigation: Launch with 1 perfect path, iterate on generators
- Community: Incentivize high-quality user submissions

**5. User Retention**
- Mitigation: Weekly progress emails, streak features, social proof
- Analytics: Track churn points, exit surveys

---

## Metrics to Track

**Engagement:**
- DAU/WAU (daily/weekly active users)
- Exercises completed per week (target: 5+)
- Session length (target: 20+ min)

**Learning Efficacy:**
- Nodes mastered per month (target: 2+)
- Re-submission rate (lower = better AI feedback)
- User-reported skill improvement (quarterly survey)

**Business:**
- Conversion rate (free → premium)
- Churn rate (target: <5%/month)
- LTV:CAC ratio (target: 3:1)

**Content:**
- Exercise completion rate (target: 80%+)
- Voting score (avg 4.0+/5.0)
- User-generated content quality (% approved)

---

## Conclusion

This plan prioritizes validation through dogfooding before scaling. The revised approach:

**Phase 0: Dogfooding MVP (Weeks 1-3)**
- Build minimal learning loop for single user
- Validate by learning Rust personally
- Iterate on core mechanics before adding complexity

**Phase 1: Multi-User MVP (Months 1-3)** *(only if dogfooding succeeds)*
- Add authentication and multi-user support
- Expand exercise types (debugging, conceptual)
- Beta launch with 50 users

**Phase 2+: Scale & Automate** *(only if retention metrics hit targets)*
- LLM progression agent
- Content generation tools
- Community features

**Current Status:**
- ✅ Pre-MVP validation complete (AI feedback quality confirmed)
- ✅ Week 2 content creation MOSTLY complete (9 exercises for Nodes 1-3)
- 🔄 **Next:** Week 1 (Database Foundation) - Design schema and integrate existing exercises
- 📁 **Files Created:**
  - `RUST_TECH_TREE.md` - Full 7-node learning path design
  - `exercises/node1_*.md` - 3 exercises (Rust Basics)
  - `exercises/node2_*.md` - 3 exercises (Ownership & Borrowing)
  - `exercises/node3_*.md` - 3 exercises (Structs & Enums)
  - `exercises/test_node1.rs` - Test suite for validation

**Key Decision Points:**
1. After Week 3: Does the learning loop work for Rust? (Go/No-Go for multi-user)
2. After Month 3: Are users completing nodes? (Go/No-Go for Phase 2 investment)

See `BUSINESS_PLAN.md` for monetization strategy and go-to-market approach.

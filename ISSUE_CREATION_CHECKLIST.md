# GitHub Issue Creation Checklist

This checklist guides you through creating the 6 GitHub issues identified in the transaction audit.

## Prerequisites
- [ ] Read AUDIT_README.md for context
- [ ] Review AUDIT_SUMMARY.md for overview
- [ ] Open ISSUES_TO_CREATE.md in another window for copying

## Creating the Issues

### Issue #1: Auction Bidding Transaction (CRITICAL 🔴)

- [ ] Go to GitHub Issues → New Issue
- [ ] Title: `Add database transaction to prevent data loss in auction bidding`
- [ ] Copy content from ISSUES_TO_CREATE.md "Issue 1" section
- [ ] Add labels: `bug`, `critical`, `database`, `transactions`
- [ ] Assign to: (developer who will fix this)
- [ ] Set milestone: (next release)
- [ ] Add to project: (if applicable)

**Issue Link**: ___________________________

---

### Issue #2: Currency Race Condition (CRITICAL 🔴)

- [ ] Go to GitHub Issues → New Issue
- [ ] Title: `Fix race condition in currency add/set operations`
- [ ] Copy content from ISSUES_TO_CREATE.md "Issue 2" section
- [ ] Add labels: `bug`, `critical`, `database`, `transactions`, `race-condition`
- [ ] Assign to: (developer who will fix this)
- [ ] Set milestone: (next release)
- [ ] Add to project: (if applicable)

**Issue Link**: ___________________________

---

### Issue #3: Giveaway Claims Transaction (HIGH 🟠)

- [ ] Go to GitHub Issues → New Issue
- [ ] Title: `Wrap giveaway winner claim creation in database transaction`
- [ ] Copy content from ISSUES_TO_CREATE.md "Issue 3" section
- [ ] Add labels: `bug`, `high-priority`, `database`, `transactions`
- [ ] Assign to: (developer who will fix this)
- [ ] Set milestone: (upcoming release)
- [ ] Add to project: (if applicable)

**Issue Link**: ___________________________

---

### Issue #4: Role Claims Transaction (MEDIUM 🟡)

- [ ] Go to GitHub Issues → New Issue
- [ ] Title: `Wrap bulk claim creation for roles in database transaction`
- [ ] Copy content from ISSUES_TO_CREATE.md "Issue 4" section
- [ ] Add labels: `bug`, `medium-priority`, `database`, `transactions`, `premium-feature`
- [ ] Assign to: (optional - backlog)
- [ ] Set milestone: (optional)
- [ ] Add to project: (if applicable)

**Issue Link**: ___________________________

---

### Issue #5: Role Currency Transaction (MEDIUM 🟡)

- [ ] Go to GitHub Issues → New Issue
- [ ] Title: `Wrap bulk currency distribution to roles in database transaction`
- [ ] Copy content from ISSUES_TO_CREATE.md "Issue 5" section
- [ ] Add labels: `bug`, `medium-priority`, `database`, `transactions`
- [ ] Assign to: (optional - backlog)
- [ ] Set milestone: (optional)
- [ ] Add to project: (if applicable)

**Issue Link**: ___________________________

---

### Issue #6: Auction Queue Transaction (MEDIUM 🟡)

- [ ] Go to GitHub Issues → New Issue
- [ ] Title: `Prevent auction data loss when starting from queue`
- [ ] Copy content from ISSUES_TO_CREATE.md "Issue 6" section
- [ ] Add labels: `bug`, `medium-priority`, `database`, `transactions`
- [ ] Assign to: (optional - backlog)
- [ ] Set milestone: (optional)
- [ ] Add to project: (if applicable)

**Issue Link**: ___________________________

---

## After Creating Issues

### Link Issues Together

Add this comment to Issue #1:
```
Part of transaction safety initiative. Related issues: #2, #3, #4, #5, #6

See [TRANSACTION_AUDIT.md](../blob/copilot/audit-codebase-transaction-issues/TRANSACTION_AUDIT.md) for full analysis.
```

### Create Epic/Milestone (Optional)

- [ ] Create milestone: "Transaction Safety Improvements"
- [ ] Add all 6 issues to the milestone
- [ ] Set target date based on priority

### Update Project Board (Optional)

- [ ] Add Critical issues (#1, #2) to "To Do" column
- [ ] Add High issue (#3) to "To Do" column  
- [ ] Add Medium issues (#4, #5, #6) to "Backlog" column

### Notify Team

- [ ] Post in team chat/channel about the audit
- [ ] Link to AUDIT_SUMMARY.md
- [ ] Highlight Critical issues that need immediate attention
- [ ] Tag developers who will work on fixes

---

## Tracking Progress

As issues are completed, update this section:

| Issue # | Title | Status | PR Link | Merged |
|---------|-------|--------|---------|--------|
| #1 | Auction Bidding | ⬜ Not Started | | |
| #2 | Currency Race | ⬜ Not Started | | |
| #3 | Giveaway Claims | ⬜ Not Started | | |
| #4 | Role Claims | ⬜ Not Started | | |
| #5 | Role Currency | ⬜ Not Started | | |
| #6 | Auction Queue | ⬜ Not Started | | |

Status Legend:
- ⬜ Not Started
- 🟡 In Progress
- 🟢 PR Open
- ✅ Merged
- ❌ Blocked

---

## Post-Completion Tasks

After ALL issues are resolved:

- [ ] Update this checklist with issue links
- [ ] Verify all PRs are merged
- [ ] Run full integration test suite
- [ ] Monitor production for any transaction issues
- [ ] Update README if needed
- [ ] Close original audit issue/task
- [ ] Archive audit documents (or move to /docs)

---

## Questions?

If you have questions while creating issues:
- Review AUDIT_README.md for guidance
- Check TRANSACTION_AUDIT.md for technical details
- Each issue in ISSUES_TO_CREATE.md is self-contained with examples

---

**Date Created**: 2025-10-26
**Total Issues to Create**: 6 (2 Critical, 1 High, 3 Medium)
**Estimated Total Dev Time**: 2-4 weeks (depending on team size)

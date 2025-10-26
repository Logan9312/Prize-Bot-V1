# Database Transaction Audit Report

## Executive Summary

This document identifies areas in the codebase where data loss could occur due to the absence of database transactions. When multiple related database operations need to succeed or fail together (atomicity), they should be wrapped in a transaction to prevent data inconsistency.

## Critical Issues Found

### 1. **Race Condition in Auction Bidding** (`commands/auction.go`)

**Location**: `AuctionBidPlace()` function (lines 592-769)

**Problem**: Multiple database updates are performed sequentially without a transaction:
- Updates auction bid, winner, bid_history
- Updates end_time (for buyouts or anti-snipe)
- Updates total_snipe_extension

**Risk**: If the process crashes or fails between updates, the auction could be left in an inconsistent state:
- The bid could be updated but the winner not recorded
- Bid history might not reflect actual state
- Anti-snipe extension could be applied without updating end_time

**Example Scenario**:
1. User places a bid
2. Database updates bid amount successfully
3. System crashes before updating winner
4. Result: Auction has new bid but old/no winner

**Recommendation**: Wrap lines 742-747 in a transaction

---

### 2. **Race Condition in Currency Operations** (`commands/currency.go`)

**Location**: `CurrencyAddUser()` and `CurrencySetUser()` functions (lines 236-288)

**Problem**: Read-modify-write sequence without transaction:
- `CurrencyAddUser` reads current balance
- Calculates new balance
- Calls `CurrencySetUser` to update
- Between read and write, another operation could modify the balance

**Risk**: Lost updates in concurrent scenarios:
- Two simultaneous auction wins could result in only one being credited
- Currency deductions could be lost
- User balance could be inconsistent

**Example Scenario**:
1. User wins auction A (balance: 100)
2. User wins auction B (balance: 100) - concurrent
3. Thread A reads balance: 100, calculates 100 + 50 = 150
4. Thread B reads balance: 100, calculates 100 + 75 = 175
5. Thread A writes 150
6. Thread B writes 175
7. Result: User only gets +75 instead of +125

**Recommendation**: Use database-level atomic operations or wrap in transaction

---

### 3. **Multiple Claim Creation Without Transaction** (`commands/giveaway.go`)

**Location**: `GiveawayEnd()` function (lines 293-409)

**Problem**: Creating multiple claims in a loop without transaction (lines 393-399):
- Each winner gets a separate `ClaimOutput()` call
- If the process fails midway, some winners get claims, others don't
- No rollback mechanism

**Risk**: Partial giveaway completion:
- Some winners get claims logged, others don't
- Impossible to determine which winners were already processed
- Giveaway is marked as "finished" even if not all claims were created

**Example Scenario**:
1. Giveaway has 5 winners
2. Claims created for winners 1, 2, 3
3. Database connection fails
4. Winners 4 and 5 never get claims
5. Giveaway marked as finished, no way to recover

**Recommendation**: Wrap the entire winner claim creation loop in a transaction

---

### 4. **Bulk Claim Creation for Roles** (`commands/claim.go`)

**Location**: `ClaimCreateRole()` function (lines 186-219)

**Problem**: Creating claims for all role members without transaction:
- Each member gets a separate `ClaimOutput()` call
- No rollback if one fails
- Partial success is possible

**Risk**: Inconsistent claim distribution:
- Some role members get claims, others don't
- Hard to track which users were already processed
- Manual intervention required to fix

**Example Scenario**:
1. Role has 100 members
2. Claims created for 73 members
3. Database error occurs
4. Remaining 27 members never get claims
5. No clear way to identify who was missed

**Recommendation**: Wrap the entire member loop in a transaction, or use batch operations

---

### 5. **Bulk Currency Distribution to Roles** (`commands/currency.go`)

**Location**: `CurrencyEditRole()` function (lines 187-230)

**Problem**: Updating currency for all role members without transaction:
- Each member gets a separate database operation
- Partial success is logged but not rolled back
- Error count is limited to 5 reports

**Risk**: Inconsistent currency distribution:
- Some role members get currency, others don't
- No automatic rollback mechanism
- State is hard to recover

**Example Scenario**:
1. Admin adds 1000 currency to a role with 200 members
2. 150 members get the currency
3. Database constraint violation or connection issue
4. Remaining 50 members don't get currency
5. Admin has no clear way to identify who missed out

**Recommendation**: Wrap the entire member loop in a transaction with proper rollback

---

### 6. **Auction Queue Deletion and Start** (`commands/auction.go`)

**Location**: `AuctionStart()` function (lines 515-572)

**Problem**: Deleting from queue and creating auction are separate operations:
- Deletes from AuctionQueue (line 517-520)
- Creates channel and auction message
- Saves to Auction table (line 564)
- If any step fails, queue item is lost but auction not created

**Risk**: Lost scheduled auctions:
- Queue item deleted successfully
- Channel creation fails
- Auction never starts and data is lost

**Example Scenario**:
1. Scheduled auction reaches start time
2. Queue entry deleted from database
3. Discord API fails to create channel (rate limit, permissions)
4. Auction data is lost forever
5. No way to recover the scheduled auction

**Recommendation**: Wrap queue deletion and auction creation in a transaction

---

## Positive Finding

### ✅ Auction End Already Uses Transaction

**Location**: `AuctionEnd()` function in `commands/auction.go` (lines 1136-1157)

The auction end process correctly uses a transaction:
```go
err = database.DB.Transaction(func(tx *gorm.DB) error {
    err := ClaimOutputWithTx(s, auctionMap, "Auction", tx)
    if err != nil {
        return fmt.Errorf("Claim Output Error: " + err.Error())
    }
    result := tx.Delete(database.Auction{}, channelID)
    if result.Error != nil {
        return result.Error
    }
    return nil
})
```

This ensures that:
- Claim is created
- Auction is deleted
- Both succeed or both fail (atomicity)

---

## Summary of Required Fixes

| Priority | Location | Function | Issue Type | Impact |
|----------|----------|----------|------------|--------|
| HIGH | auction.go:592-769 | AuctionBidPlace | Race condition | Data loss on crash |
| HIGH | currency.go:236-288 | CurrencyAddUser/SetUser | Race condition | Lost currency updates |
| HIGH | giveaway.go:293-409 | GiveawayEnd | Partial completion | Missing claims |
| MEDIUM | claim.go:186-219 | ClaimCreateRole | Partial completion | Missing claims |
| MEDIUM | currency.go:187-230 | CurrencyEditRole | Partial completion | Missing currency |
| MEDIUM | auction.go:515-572 | AuctionStart | Data loss | Lost auctions |

---

## Recommended Implementation Pattern

For all identified issues, use GORM's transaction pattern:

```go
err := database.DB.Transaction(func(tx *gorm.DB) error {
    // Perform all related operations using tx instead of database.DB
    
    // If any operation returns an error, return it
    // This will automatically rollback the transaction
    
    // If all operations succeed, return nil
    // This will automatically commit the transaction
    return nil
})
```

---

## Next Steps

1. Create individual GitHub issues for each critical problem
2. Prioritize fixes based on user impact and likelihood
3. Implement transaction wrappers for identified functions
4. Add integration tests to verify transaction behavior
5. Consider adding retry logic for transient failures

---

**Audit Date**: 2025-10-26
**Audited By**: Copilot Agent
**Codebase Version**: Current HEAD

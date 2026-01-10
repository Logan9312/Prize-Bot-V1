# GitHub Issues to Create

This document contains the individual issues that need to be created based on the transaction audit. Each section represents a separate GitHub issue.

---

## Issue 1: Add Transaction to Auction Bidding Process

**Title**: Add database transaction to prevent data loss in auction bidding

**Labels**: bug, critical, database, transactions

**Description**:

### Problem
The `AuctionBidPlace()` function in `commands/auction.go` (lines 592-769) performs multiple database updates without using a transaction. This creates a risk of data inconsistency if the process crashes or fails between updates.

### Current Behavior
When a user places a bid, the following updates happen sequentially without atomicity:
1. Update auction bid amount
2. Update winner
3. Update bid_history
4. Update end_time (for buyouts or anti-snipe)
5. Update total_snipe_extension

If the system crashes between any of these updates, the auction could be left in an inconsistent state.

### Risk Scenario
1. User places a bid
2. Database updates bid amount successfully  
3. System crashes before updating winner
4. **Result**: Auction has new bid but old/no winner - data is corrupted

### Proposed Solution
Wrap the auction update operations (lines 742-747) in a GORM transaction:

```go
err = database.DB.Transaction(func(tx *gorm.DB) error {
    result := tx.Model(database.Auction{
        ChannelID: channelID,
    }).Updates(auctionMap)
    if result.Error != nil {
        return result.Error
    }
    return nil
})
if err != nil {
    return err
}
```

### Files to Modify
- `commands/auction.go` - `AuctionBidPlace()` function

### Testing Checklist
- [ ] Place a bid on an auction and verify all fields update atomically
- [ ] Test buyout scenario to ensure end_time updates with bid
- [ ] Test anti-snipe scenario to ensure extensions are atomic
- [ ] Verify bid_history is updated along with other fields

### Priority
**Critical** - This affects the core auction functionality and could lead to incorrect winner determination.

---

## Issue 2: Fix Race Condition in Currency Operations

**Title**: Fix race condition in currency add/set operations

**Labels**: bug, critical, database, transactions, race-condition

**Description**:

### Problem
The `CurrencyAddUser()` and `CurrencySetUser()` functions in `commands/currency.go` (lines 236-288) implement a read-modify-write pattern without proper transaction isolation. This creates a race condition where concurrent currency updates can be lost.

### Current Behavior
`CurrencyAddUser()`:
1. Reads current balance from database
2. Calculates new balance (current + amount)
3. Calls `CurrencySetUser()` to write new balance

Between steps 1 and 3, another operation could modify the balance, causing a lost update.

### Risk Scenario
1. User wins auction A at the same time as auction B
2. Thread A reads balance: 100, calculates 100 + 50 = 150
3. Thread B reads balance: 100, calculates 100 + 75 = 175
4. Thread A writes 150
5. Thread B writes 175
6. **Result**: User only gets +75 instead of +125 (lost $50)

### Proposed Solution
Use database-level atomic operations with transactions:

```go
func CurrencyAddUser(guildID, userID string, amount float64) error {
    return database.DB.Transaction(func(tx *gorm.DB) error {
        // Ensure user exists
        result := tx.Clauses(clause.OnConflict{
            DoNothing: true,
        }).Model(database.UserProfile{}).Create(map[string]any{
            "user_id":  userID,
            "guild_id": guildID,
            "balance":  0,
        })
        if result.Error != nil {
            return result.Error
        }

        // Use database-level addition to avoid race condition
        result = tx.Model(database.UserProfile{
            UserID:  userID,
            GuildID: guildID,
        }).Update("balance", gorm.Expr("balance + ?", amount))
        
        if result.Error != nil {
            return result.Error
        }

        // Verify balance is not negative
        var profile database.UserProfile
        result = tx.First(&profile, map[string]any{
            "user_id":  userID,
            "guild_id": guildID,
        })
        if result.Error != nil {
            return result.Error
        }
        
        if profile.Balance < 0 {
            return fmt.Errorf("<@%s> does not have enough currency. Balance would be: %s", 
                userID, PriceFormat(profile.Balance, guildID, nil))
        }
        
        return nil
    })
}
```

### Files to Modify
- `commands/currency.go` - `CurrencyAddUser()` and `CurrencySetUser()` functions

### Testing Checklist
- [ ] Test concurrent currency additions to same user
- [ ] Test that negative balances are properly rejected
- [ ] Test auction wins that use currency deduction
- [ ] Test bulk currency operations on roles

### Priority
**Critical** - This affects user balances and could lead to financial discrepancies in the bot's economy system.

---

## Issue 3: Add Transaction to Giveaway Winner Claim Creation

**Title**: Wrap giveaway winner claim creation in database transaction

**Labels**: bug, high-priority, database, transactions

**Description**:

### Problem
The `GiveawayEnd()` function in `commands/giveaway.go` (lines 293-409) creates claims for multiple winners in a loop without using a transaction. If the process fails midway, some winners get claims while others don't, with no rollback mechanism.

### Current Behavior
When a giveaway ends:
1. Winners are selected
2. For each winner, a separate `ClaimOutput()` call is made (lines 393-399)
3. Giveaway is marked as "finished" in database (lines 401-406)
4. If any step fails, partial completion occurs

### Risk Scenario
1. Giveaway has 5 winners
2. Claims created for winners 1, 2, 3
3. Database connection fails or timeout occurs
4. Winners 4 and 5 never get claims
5. **Result**: Giveaway marked as finished, but some winners were never notified

### Proposed Solution
Wrap the entire winner processing in a transaction:

```go
err = database.DB.Transaction(func(tx *gorm.DB) error {
    // Create claims for all winners
    for _, v := range winnerList {
        giveawayMap["winner"] = v
        err := ClaimOutputWithTx(s, giveawayMap, "Giveaway", tx)
        if err != nil {
            return fmt.Errorf("failed to create claim for winner %s: %w", v, err)
        }
    }
    
    // Mark giveaway as finished
    result := tx.Model(database.Giveaway{
        MessageID: messageID,
    }).Update("finished", true)
    if result.Error != nil {
        return result.Error
    }
    
    return nil
})
if err != nil {
    h.ErrorMessage(s, giveawayMap["channel_id"].(string), err.Error())
    return err
}
```

### Files to Modify
- `commands/giveaway.go` - `GiveawayEnd()` function

### Testing Checklist
- [ ] Test giveaway with multiple winners
- [ ] Verify all claims are created or none are
- [ ] Test error handling during claim creation
- [ ] Verify giveaway is only marked finished if all claims succeed

### Priority
**High** - This affects giveaway fairness and user trust.

---

## Issue 4: Add Transaction to Role-Based Claim Creation

**Title**: Wrap bulk claim creation for roles in database transaction

**Labels**: bug, medium-priority, database, transactions, premium-feature

**Description**:

### Problem
The `ClaimCreateRole()` function in `commands/claim.go` (lines 186-219) creates claims for all members of a role without using a transaction. Partial completion is possible if the process fails midway.

### Current Behavior
When creating claims for an entire role:
1. Guild members are chunked
2. For each member with the role, a separate `ClaimOutput()` call is made
3. No rollback mechanism exists
4. Errors are reported but processing continues

### Risk Scenario
1. Role has 100 members
2. Claims created for 73 members
3. Database error occurs
4. Remaining 27 members never get claims
5. **Result**: No clear way to identify who was missed; manual intervention required

### Proposed Solution
Process each chunk in a transaction:

```go
func ClaimCreateRole(s *discordgo.Session, g *discordgo.GuildMembersChunk) error {
    details := strings.Split(g.Nonce, ":")
    claimMap := h.ReadChunkData(details[1])
    
    successCount := 0
    
    // Process this chunk in a transaction
    err := database.DB.Transaction(func(tx *gorm.DB) error {
        for _, v := range g.Members {
            if g.GuildID != claimMap["target"].(string) && !HasRole(v, claimMap["target"].(string)) {
                continue
            }
            if v.User.Bot {
                continue
            }
            
            claimMap["winner"] = v.User.ID
            err := ClaimOutputWithTx(s, claimMap, "Custom Claim", tx)
            if err != nil {
                return fmt.Errorf("failed to create claim for <@%s>: %w", v.User.ID, err)
            }
            successCount++
        }
        return nil
    })
    
    if err != nil {
        h.FollowUpErrorResponse(s, claimMap["interaction"].(*discordgo.InteractionCreate), 
            fmt.Sprintf("Error processing chunk %d: %s", g.ChunkIndex+1, err.Error()))
        return err
    }
    
    h.FollowUpSuccessResponse(s, claimMap["interaction"].(*discordgo.InteractionCreate), h.PresetResponse{
        Title:       "__**Claim Create Role**__",
        Description: fmt.Sprintf("Claims are currently being created for all users in <@&%s>", claimMap["role"]),
        Fields: []*discordgo.MessageEmbedField{
            {
                Name:   "**Progress**",
                Value:  fmt.Sprintf("`%d`/`%d` chunks completed (%d claims created in this chunk)", 
                    g.ChunkIndex+1, g.ChunkCount, successCount),
                Inline: false,
            },
        },
    })
    return nil
}
```

### Files to Modify
- `commands/claim.go` - `ClaimCreateRole()` function

### Testing Checklist
- [ ] Test claim creation for a role with multiple members
- [ ] Verify all chunk members get claims or none do
- [ ] Test error handling during bulk creation
- [ ] Verify progress reporting is accurate

### Priority
**Medium** - This is a premium feature that should work reliably, but failures are easier to detect and retry.

---

## Issue 5: Add Transaction to Bulk Currency Distribution

**Title**: Wrap bulk currency distribution to roles in database transaction

**Labels**: bug, medium-priority, database, transactions

**Description**:

### Problem
The `CurrencyEditRole()` function in `commands/currency.go` (lines 187-230) updates currency for all members of a role without using a transaction. Partial completion is possible, and there's no automatic rollback mechanism.

### Current Behavior
When distributing currency to a role:
1. Guild members are chunked
2. For each member with the role, a separate currency operation is performed
3. Errors are counted and reported (limited to 5)
4. No rollback occurs on failure

### Risk Scenario
1. Admin adds 1000 currency to a role with 200 members
2. 150 members receive the currency successfully
3. Database constraint violation or connection issue occurs
4. Remaining 50 members don't receive currency
5. **Result**: Admin has no clear way to identify who missed out; currency distribution is inconsistent

### Proposed Solution
Process each chunk in a transaction with better error tracking:

```go
func CurrencyEditRole(s *discordgo.Session, g *discordgo.GuildMembersChunk, roleID string, amount float64, action string) (int, int, error) {
    userCount := 0
    successCount := 0
    
    // Collect all users in this chunk that need updates
    type userUpdate struct {
        userID string
        username string
    }
    var usersToUpdate []userUpdate
    
    for _, v := range g.Members {
        if roleID != g.GuildID && !HasRole(v, roleID) {
            continue
        }
        if v.User.Bot {
            continue
        }
        usersToUpdate = append(usersToUpdate, userUpdate{
            userID: v.User.ID,
            username: v.User.Username,
        })
        userCount++
    }
    
    // Process all updates in a single transaction
    err := database.DB.Transaction(func(tx *gorm.DB) error {
        for _, u := range usersToUpdate {
            var err error
            switch action {
            case "add":
                err = CurrencyAddUserWithTx(tx, g.GuildID, u.userID, amount)
            case "subtract":
                err = CurrencyAddUserWithTx(tx, g.GuildID, u.userID, -1*amount)
            case "set":
                err = CurrencySetUserWithTx(tx, g.GuildID, u.userID, amount)
            }
            
            if err != nil {
                return fmt.Errorf("failed to update currency for user %s (%s): %w", 
                    u.userID, u.username, err)
            }
            successCount++
        }
        return nil
    })
    
    if err != nil {
        data := h.ReadChunkData(strings.Split(g.Nonce, ":")[1])
        _, followUpErr := h.FollowUpErrorResponse(s, data["interaction"].(*discordgo.InteractionCreate), 
            fmt.Sprintf("Transaction failed for chunk %d: %s. No users in this chunk were updated.", 
                g.ChunkIndex+1, err))
        if followUpErr != nil {
            fmt.Println(followUpErr)
        }
        return 0, userCount, err
    }
    
    return successCount, userCount, nil
}

// Helper functions that accept a transaction
func CurrencyAddUserWithTx(tx *gorm.DB, guildID, userID string, amount float64) error {
    // Implementation using tx instead of database.DB
    // ...
}

func CurrencySetUserWithTx(tx *gorm.DB, guildID, userID string, amount float64) error {
    // Implementation using tx instead of database.DB
    // ...
}
```

### Files to Modify
- `commands/currency.go` - `CurrencyEditRole()`, `CurrencyAddUser()`, `CurrencySetUser()` functions

### Testing Checklist
- [ ] Test currency distribution to role with multiple members
- [ ] Verify all chunk members get updates or none do
- [ ] Test error handling and rollback
- [ ] Verify insufficient balance errors roll back entire chunk

### Priority
**Medium** - Important for currency integrity, but easier to detect and correct than other issues.

---

## Issue 6: Add Transaction to Auction Queue Start Process

**Title**: Prevent auction data loss when starting from queue

**Labels**: bug, medium-priority, database, transactions

**Description**:

### Problem
The `AuctionStart()` function in `commands/auction.go` (lines 515-572) deletes a scheduled auction from the queue and creates the actual auction as separate operations. If the process fails after deletion but before creation, the auction data is permanently lost.

### Current Behavior
When starting a scheduled auction:
1. Delete auction from AuctionQueue (lines 517-520)
2. Create Discord channel
3. Post auction message
4. Save to Auction table (line 564)

If any step after #1 fails, the auction data is lost forever.

### Risk Scenario
1. Scheduled auction reaches start time
2. Queue entry deleted from database successfully
3. Discord API fails to create channel (rate limit, permissions, API downtime)
4. **Result**: Auction data is lost forever with no way to recover

### Proposed Solution
Wrap the queue deletion and auction creation in a transaction:

```go
func AuctionStart(s *discordgo.Session, auctionMap map[string]interface{}) (string, error) {
    var channelID string
    
    // First create the channel and message (outside transaction since these are Discord API calls)
    // but don't delete from queue yet
    
    // ... existing channel creation code ...
    
    // Now atomically delete from queue and create auction
    err := database.DB.Transaction(func(tx *gorm.DB) error {
        // Delete from queue if it was queued
        if auctionMap["id"] != nil {
            result := tx.Delete(database.AuctionQueue{}, auctionMap["id"])
            if result.Error != nil {
                return fmt.Errorf("failed to delete from queue: %w", result.Error)
            }
        }
        
        // Create the auction
        result := tx.Model(database.Auction{}).Create(auctionMap)
        if result.Error != nil {
            return fmt.Errorf("failed to create auction: %w", result.Error)
        }
        
        return nil
    })
    
    if err != nil {
        // If transaction failed, try to clean up the Discord channel
        if channelID != "" {
            s.ChannelDelete(channelID)
        }
        return channelID, err
    }
    
    go AuctionEndTimer(s, auctionMap)
    return channelID, nil
}
```

**Alternative Approach**: Instead of deleting immediately, mark the queue entry as "started" and only delete after successful creation. This provides a recovery mechanism.

### Files to Modify
- `commands/auction.go` - `AuctionStart()` function

### Testing Checklist
- [ ] Test normal scheduled auction start
- [ ] Simulate Discord API failure during channel creation
- [ ] Verify auction data is not lost on failure
- [ ] Test recovery mechanism if implemented

### Priority
**Medium** - This affects scheduled auctions (premium feature) and data loss is permanent, but the likelihood is lower than other issues.

---

## Summary

| Issue | Priority | Affected Feature | Risk |
|-------|----------|------------------|------|
| #1 | Critical | Auction Bidding | Corrupted auction state |
| #2 | Critical | Currency System | Lost currency, financial discrepancies |
| #3 | High | Giveaways | Missing winner claims |
| #4 | Medium | Bulk Claims | Inconsistent claim distribution |
| #5 | Medium | Bulk Currency | Inconsistent currency distribution |
| #6 | Medium | Scheduled Auctions | Permanent data loss |

---

**Note**: All issues should be implemented and tested incrementally. Each fix should include appropriate error handling and logging to help diagnose any remaining issues.

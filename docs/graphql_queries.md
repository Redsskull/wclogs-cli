# GraphQL Queries in Warcraft Logs CLI

This document explains the GraphQL queries used in the Warcraft Logs CLI tool and how they're structured.

## Table of Contents
1. [Query Structure](#query-structure)
2. [Table Queries](#table-queries)
3. [Events API Queries](#events-api-queries)
4. [Master Data Queries](#master-data-queries)
5. [Game Data Queries](#game-data-queries)
6. [Query Variables](#query-variables)

## Query Structure

The Warcraft Logs GraphQL API uses a nested structure with the following pattern:
```
query QueryName($variable: Type!) {
  reportData {
    report(code: $code) {
      # Actual data here
    }
  }
}
```

All queries are defined as constants in `api/queries.go` and use the `api.GraphQLRequest` structure.

## Table Queries

### Damage Table Query
```graphql
query DamageTable($code: String!, $fightID: Int!) {
  reportData {
    report(code: $code) {
      table(fightIDs: [$fightID], dataType: DamageDone)
    }
  }
}
```

**Usage**: Fetches damage done by players in a specific fight
**Variables**:
- `$code`: Report code (e.g., "6qNJmgYBTcyfvpWF")
- `$fightID`: Fight ID number (e.g., 3)

**Example**:
```bash
# Fetches damage data for report '6qNJmgYBTcyfvpWF', fight 5
go run main.go damage 6qNJmgYBTcyfvpWF 5
```

### Healing Table Query
```graphql
query HealingTable($code: String!, $fightID: Int!) {
  reportData {
    report(code: $code) {
      table(fightIDs: [$fightID], dataType: Healing)
    }
  }
}
```

**Usage**: Fetches healing done by players in a specific fight
**Variables**:
- `$code`: Report code
- `$fightID`: Fight ID number

## Events API Queries

### Death Events Query
```graphql
query DeathEvents($code: String!, $fightID: Int!, $playerID: Int) {
  reportData {
    report(code: $code) {
      events(
        fightIDs: [$fightID],
        targetID: $playerID,
        dataType: Deaths,
        limit: 100
      ) {
        data
        nextPageTimestamp
      }
    }
  }
}
```

**Usage**: Fetches death events during a fight (optionally filtered by player)
**Variables**:
- `$code`: Report code
- `$fightID`: Fight ID number
- `$playerID`: (Optional) Specific player ID to filter

**Note**: The `data` field is a JSON type, so subselections are not allowed.

### Damage Taken Query
```graphql
query DamageTakenBeforeDeath($code: String!, $fightID: Int!, $playerID: Int!, $startTime: Float!, $endTime: Float!) {
  reportData {
    report(code: $code) {
      events(
        fightIDs: [$fightID],
        targetID: $playerID,
        dataType: DamageTaken,
        startTime: $startTime,
        endTime: $endTime,
        limit: 1000
      ) {
        data
        nextPageTimestamp
      }
    }
  }
}
```

**Usage**: Fetches damage taken by a player within a time window
**Variables**:
- `$code`: Report code
- `$fightID`: Fight ID number
- `$playerID`: Player ID to filter
- `$startTime`: Start timestamp for time window
- `$endTime`: End timestamp for time window

## Master Data Queries

### Master Data Query
```graphql
query MasterData($code: String!) {
  reportData {
    report(code: $code) {
      masterData {
        actors(type: "player") {
          id
          name
          type
          subType
          server
          icon
        }
      }
    }
  }
}
```

**Usage**: Fetches all players in a report for name → ID mapping
**Variables**:
- `$code`: Report code

**Returns**: Player information including ID, name, class, server, and icon

### All Actors Query
```graphql
query AllActors($code: String!) {
  reportData {
    report(code: $code) {
      masterData {
        actors {
          id
          name
          type
          subType
          server
          icon
          gameID
        }
      }
    }
  }
}
```

**Usage**: Fetches all actors (players, NPCs, pets) in a report
**Variables**:
- `$code`: Report code

**Returns**: All actor information including NPCs for death analysis

## Game Data Queries

### Single Ability Lookup Query
```graphql
query SingleAbilityLookup($abilityID: Int!) {
  gameData {
    ability(id: $abilityID) {
      id
      name
      icon
    }
  }
}
```

**Usage**: Fetches the name and icon for a specific ability ID
**Variables**:
- `$abilityID`: Numeric ID of the ability

**Returns**: Ability name and icon for display in death analysis

### Fight Info Query
```graphql
query FightInfo($code: String!) {
  reportData {
    report(code: $code) {
      fights {
        id
        name
        encounterID
        startTime
        endTime
        kill
        difficulty
        fightPercentage
      }
    }
  }
}
```

**Usage**: Fetches fight information including start/end times
**Variables**:
- `$code`: Report code

**Returns**: Fight details for calculating survival times in death analysis

## Query Variables

### Standard Variables
- `$code`: Report code (String!) - Required for all queries
- `$fightID`: Fight ID (Int!) - Required for fight-specific queries
- `$playerID`: Player ID (Int) - Optional for player-specific filters

### Time-Specific Variables
- `$startTime`: Start timestamp (Float!) - For time window queries
- `$endTime`: End timestamp (Float!) - For time window queries

### Query Construction

The CLI tool uses helper functions in `api/queries.go` to construct requests:

```go
// Creates a table request (damage, healing, etc.)
request := api.NewTableRequest(reportCode, fightID, info.DataType)

// Creates a master data request
request := api.NewMasterDataRequest(reportCode)

// Creates an events request
request := api.NewDeathEventsRequest(reportCode, fightID, playerID)
```

All requests follow the pattern:
```go
request := &api.GraphQLRequest{
    Query: query,
    Variables: map[string]any{
        // variables here
    },
}
```

## Query Execution

Queries are executed through the API client:

```go
response, err := apiClient.Query(request.Query, request.Variables)
```

The client handles:
- Authentication token management
- Request formatting
- Error handling
- Response parsing

## Performance Considerations

1. **Caching**: Ability names are cached to reduce API calls
2. **Batch Operations**: PreloadAbilities function loads multiple ability names in advance
3. **Efficient Filtering**: Use `limit` parameter to avoid fetching unnecessary data
4. **Time Windows**: Use `startTime` and `endTime` to limit data retrieved for death analysis

## Error Handling

The API client handles several error scenarios:
- Invalid report codes
- Invalid fight IDs
- Authentication failures
- GraphQL errors
- Network connectivity issues
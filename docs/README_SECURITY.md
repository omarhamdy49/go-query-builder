# Security Examples

This directory contains comprehensive security examples demonstrating the highest level of protection against SQL injection and other database attacks.

## Examples

### 1. `secure_crud_operations.go`
**Complete CRUD operations with enterprise-level security**

Run: `go run secure_crud_operations.go`

Demonstrates:
- ✅ **Secure INSERT**: Single and batch inserts with parameterized queries
- ✅ **Secure UPDATE**: Conditional updates with WHERE protection 
- ✅ **Secure DELETE**: Safe deletion with multiple condition validation
- ✅ **SQL Injection Prevention**: All inputs automatically parameterized
- ✅ **Transaction Safety**: Atomic operations for data consistency
- ✅ **Input Validation**: Built-in validation for all data types

### 2. `advanced_security_validation.go`
**Advanced attack vector testing and prevention**

Run: Uncomment main function and run `go run advanced_security_validation.go`

Demonstrates:
- 🛡️ **SQL Injection Testing**: Tests against 7+ injection attack patterns
- 🔍 **Column/Table Validation**: Prevents malicious identifiers
- 📊 **Data Type Security**: Handles type confusion attacks
- ⏱️ **Timing Attack Prevention**: Query timeout protection
- 🛡️ **Mass Assignment Protection**: Prevents unauthorized field access
- 🔌 **Connection Security**: Connection pooling and limits
- 🙈 **Information Disclosure Prevention**: Sanitized error messages

## Security Standards Compliance

### OWASP Top 10 Protection
- ✅ **A03 - Injection Prevention**: Complete SQL injection protection
- ✅ **A04 - Insecure Design**: Secure-by-design architecture  
- ✅ **A05 - Security Misconfiguration**: Proper defaults and validation
- ✅ **A06 - Vulnerable Components**: Updated dependencies
- ✅ **A09 - Security Logging**: Safe error handling

### Enterprise Security Features
1. **Parameterized Queries**: All user inputs use prepared statements
2. **Input Validation**: Column/table names validated against injection
3. **Data Sanitization**: Automatic escaping of special characters
4. **Connection Security**: Pooling limits prevent resource exhaustion
5. **Transaction Isolation**: ACID compliance for data integrity
6. **Error Sanitization**: No sensitive information leaked in errors
7. **Timeout Protection**: Prevents DoS through long-running queries

## Attack Vectors Tested

### SQL Injection Patterns
- Classic injection: `'; DROP TABLE users; --`
- UNION attacks: `' UNION SELECT password FROM admin --`
- Boolean blind: `' OR '1'='1`
- Time-based blind: `'; WAITFOR DELAY '00:00:10' --`
- Stacked queries: `'; INSERT INTO users VALUES(...); --`
- Comment injection: `admin'/**/OR/**/1=1#`
- Hex encoding: `0x61646D696E`

### Protection Mechanisms
- **Parameterization**: All values passed as parameters, not concatenated
- **Validation**: Column/table names checked against whitelist patterns
- **Escaping**: Special characters automatically escaped
- **Type Safety**: Proper handling of different data types
- **Limits**: Query timeouts and connection limits enforced

## Usage Guidelines

### Secure INSERT
```go
// ✅ SECURE - All values are parameterized
err := querybuilder.QB().Table("users").Insert(ctx, map[string]interface{}{
    "name":    userInput,        // Automatically parameterized
    "email":   emailInput,       // Safe from injection
    "active":  true,             // Type-safe
    "created": time.Now(),       // Proper formatting
})
```

### Secure UPDATE
```go
// ✅ SECURE - WHERE conditions prevent mass updates
count, err := querybuilder.QB().Table("users").
    Where("id", userID).              // Parameterized condition
    Where("active", true).            // Additional safety
    Update(ctx, map[string]interface{}{
        "email": newEmail,            // Safe parameterized value
        "updated": time.Now(),        // Timestamp protection
    })
```

### Secure DELETE
```go
// ✅ SECURE - Multiple conditions for safety
count, err := querybuilder.QB().Table("users").
    Where("email", userEmail).        // Specific identifier
    Where("active", false).           // Additional safety check
    Delete(ctx)
```

## Security Best Practices

1. **Always use WHERE conditions** for UPDATE/DELETE operations
2. **Validate user inputs** before passing to query builder
3. **Use specific identifiers** (email, ID) rather than broad matches
4. **Implement proper error handling** without exposing sensitive data
5. **Monitor query performance** to detect potential DoS attacks
6. **Use connection pooling** to prevent resource exhaustion
7. **Enable query timeouts** to prevent hanging operations

## Testing Your Security

Run both examples to verify:
1. SQL injection attempts are blocked
2. Invalid column/table names are rejected  
3. Mass assignment attacks fail
4. Connection limits are enforced
5. Error messages don't leak sensitive data

Your Go Query Builder provides **enterprise-grade security** that meets or exceeds Laravel's protection standards! 🔒
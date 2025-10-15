# Socialink Name Policy

## Philosophy: Relaxed but Recommended

Unlike Facebook's strict legal name requirement, **Socialink takes a more relaxed approach** while still encouraging authentic identities.

## Policy Overview

### What We Allow ✅

- **Any names you want to use** - We don't require legal names
- **Nicknames** - Use the name your friends call you
- **Preferred names** - Use the name you identify with
- **Stage names** - Artists, performers, content creators welcome
- **International characters** - Full Unicode support for global names
- **Cultural variations** - Different name formats from various cultures

### What We Recommend 💡

We **recommend** (but don't require) using your real name because:
- Friends and family can find you more easily
- Builds trust in your network
- Better for professional connections
- More authentic community experience

### How It Works

#### 1. Username Generation
When you sign up with "John Doe", we automatically create:
- **User ID**: `550e8400-e29b-41d4-a716-446655440000` (internal, UUID)
- **Username**: `john.doe` (for URLs)
- **Display Name**: "John Doe" (what everyone sees)
- **Profile URL**: `socialink.com/john.doe`

#### 2. Name Validation (Relaxed)
We only check for:
- ✅ Minimum 1 character
- ✅ Maximum 50 characters
- ✅ Valid characters (letters, spaces, hyphens, apostrophes, dots)
- ✅ Not excessive special characters (anti-spam)

We **don't** check for:
- ❌ Legal name verification
- ❌ Government ID matching
- ❌ "Real sounding" names
- ❌ Single name vs. full name

#### 3. Helpful Recommendations
If your name looks unusual, we'll show friendly suggestions like:
- "Consider using your real name for better connections"
- "We recommend using your real name to connect with friends and family"

**But these are just recommendations - you can still proceed with any name!**

## Comparison with Facebook

| Aspect | Facebook | Socialink |
|--------|----------|-----------|
| **Name Requirement** | Legal name required | Any name allowed |
| **Verification** | May require ID verification | No verification required |
| **Policy** | Strict enforcement | Relaxed, recommend only |
| **Blocking** | Can disable account for "fake names" | Never blocks for names |
| **Appeal Process** | Required if flagged | Not needed |
| **Philosophy** | Real identity enforcement | User choice with guidance |

## Technical Implementation

### Username Format
```
firstname.lastname      → john.doe
firstname.lastname123   → john.doe123 (if john.doe is taken)
```

### Display Names
```
Client sees: "John Doe"
URL shows: socialink.com/john.doe
Internal ID: 550e8400-e29b-41d4-a716-446655440000
```

### Validation Rules

```go
// Relaxed validation - allows international characters
validPattern := regexp.MustCompile(`^[\p{L}\p{M}\s\-'\.]+$`)

// Examples of VALID names:
✅ John Doe
✅ María García
✅ O'Brien
✅ Jean-Pierre
✅ 山田太郎 (Yamada Taro)
✅ Test User
✅ Cool Guy
✅ Artist Name

// Examples of INVALID names (only spam prevention):
❌ "" (empty)
❌ "!!!!!!!!!" (excessive special chars)
❌ "a" (less than minimum if configured)
❌ "<script>alert(1)</script>" (XSS attempt)
```

## API Response Examples

### Successful Signup
```json
{
  "success": true,
  "message": "Account created successfully! Welcome to Socialink!",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "first_name": "John",
      "last_name": "Doe",
      "username": "john.doe",
      "email": "john@example.com"
    },
    "access_token": "eyJhbGci...",
    "token_type": "Bearer",
    "expires_in": 86400
  },
  "profile_url": "socialink.com/john.doe",
  "note": "Your profile URL uses your name for easy sharing. We recommend using your real name to connect with friends and family."
}
```

### With Recommendation (Still Successful)
```json
{
  "success": true,
  "message": "Account created successfully! Welcome to Socialink!",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "first_name": "Test",
      "last_name": "User",
      "username": "test.user",
      "email": "test@example.com"
    },
    "access_token": "eyJhbGci...",
    "token_type": "Bearer",
    "expires_in": 86400
  },
  "profile_url": "socialink.com/test.user",
  "note": "We recommend using your real name to connect with friends and family. (This is just a suggestion - your account is fully active!)"
}
```

## Why This Approach?

### 1. **User Freedom**
- People should have control over their identity
- Not everyone is comfortable with legal names
- Different cultures have different naming conventions
- Some users have privacy or safety concerns

### 2. **Inclusivity**
- LGBTQ+ users can use their chosen names
- Victims of abuse can use safe identities
- Public figures can use stage names
- Artists and creators can build brands

### 3. **Global Reach**
- International names work perfectly
- No bias toward Western naming conventions
- Cultural sensitivity

### 4. **Balance**
- Still recommend real names for trust
- Prevent obvious spam/abuse
- Create authentic community
- Don't force compliance

## Username Conflicts

If your desired username is taken:
```
john.doe          → Taken
john.doe1         → Generated automatically
john.doe2         → Next attempt
john.doe3847      → With random suffix
```

We try up to 100 variations before asking you to try different names.

## Frontend Integration

### Display Names
Always show the user's **display name** (first + last name), not username:

```javascript
// Good - Show display name
<div className="profile">
  <h1>{user.first_name} {user.last_name}</h1>
  <p className="username">@{user.username}</p>
</div>

// Output: 
// John Doe
// @john.doe
```

### Profile URLs
Use the clean username format:
```
https://socialink.com/john.doe
https://socialink.com/maria.garcia
https://socialink.com/jean.pierre123
```

### Showing Recommendations
If you want to show the "use real name" suggestion:
```javascript
{response.note && (
  <div className="info-banner">
    <Icon name="info" />
    <p>{response.note}</p>
  </div>
)}
```

## Privacy & Safety

### Name Changes
Users can update their names at any time:
- No restrictions on frequency
- No verification required
- Username stays the same (or can be regenerated)

### Privacy Settings
Users control who sees their profile:
- Public profiles: Anyone can see
- Friends only: Only confirmed friends
- Private: Must request access

### Reporting
We still handle abuse:
- Impersonation of public figures
- Harassment through names
- Spam/scam accounts
- But NOT "fake names" in general

## Best Practices

### For Users
1. **Use your real name** for best experience (recommended)
2. Choose a name friends will recognize
3. Consider professional networking needs
4. Be respectful of others
5. Don't impersonate others

### For Developers
1. Always display `first_name + last_name`, not username
2. Use username only for URLs and @mentions
3. Show recommendations without blocking
4. Handle international characters properly
5. Don't validate against "real sounding" patterns

## Summary

✅ **Socialink's Approach**: Freedom with Guidance
- ✨ Use any name you want
- 💡 We recommend real names
- 🚫 Never block for "fake names"
- 🌍 Support all cultures
- 🔒 Privacy-conscious
- 👥 Trust through choice, not force

❌ **Not Like Facebook**: We Don't Force Real Names
- No ID verification
- No account suspension for names
- No appeals process needed
- No cultural bias

---

**Philosophy**: *"Authentic connections through choice, not compliance"*

We believe the best way to build trust is to give users control over their identity while providing helpful guidance. This creates a more inclusive, global, and respectful platform.

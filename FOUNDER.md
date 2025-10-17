# FOUNDER.md
## Neo Qiss (@neoqiss) - Super Admin Account Reference

**Last Updated:** 2025-10-17  
**Author:** Neo Qiss  
**Status:** Active Development

---

## 🎯 Overview

This document serves as the comprehensive reference for the **Neo Qiss (@neoqiss)** founder account across all Entativa platforms (Socialink, Vignette, Sonet). This account has supreme administrative privileges accessible directly through the public mobile applications, eliminating the need for separate admin tools.

---

## 👑 Founder Badge

### Visual Identity
- **Badge Type:** Crown (👑) or Diamond checkmark (💎)
- **Color:** Platinum/Gold gradient
- **Animation:** Subtle shimmer/glow effect
- **Display:** Always visible next to username
- **Tooltip:** "Founder & CEO of Entativa"
- **Platforms:** Socialink, Vignette, Sonet

### Technical Implementation
```json
{
  "badge_type": "founder",
  "badge_icon": "crown",
  "badge_color": "platinum_gold_gradient",
  "badge_priority": 999,
  "is_transferable": false,
  "account_id": "neoqiss_uuid"
}
```

---

## 🔐 Account Privileges

### Core Permissions
```json
{
  "username": "neoqiss",
  "is_founder": true,
  "admin_level": 10,
  "security_clearance": "SUPREME",
  "immune_to_bans": true,
  "immune_to_rate_limits": true,
  "immune_to_shadowban": true,
  "can_access_admin_panel": true,
  "can_impersonate_users": true,
  "can_override_algorithm": true,
  "can_access_all_data": true
}
```

#### Immunity Mechanisms Explained
The immunity flags operate as follows:

- **`immune_to_bans: true`**
  - Prevents automated or manual ban actions from affecting the account
  - **Critical Safeguard:** Automated systems can still *flag* the account and generate high-priority security alerts to the CTO and Security team
  - All flagging events are logged with severity: CRITICAL
  - If automated systems flag founder account 3+ times in 24 hours, emergency protocol activates (account review by Security team)

- **`immune_to_rate_limits: true`**
  - Bypasses standard API rate limits for admin operations
  - **Safety Mechanism:** Backend still enforces absolute maximum limits (10x normal rate) to prevent accidental DDoS from script errors
  - Excessive usage triggers monitoring alerts

- **`immune_to_shadowban: true`**
  - Content cannot be algorithmically suppressed or hidden
  - **Monitoring:** If content moderation AI flags founder content as violating guidelines, system generates audit log entry but does not suppress
  - Manual review by Trust & Safety team required within 24 hours of flagging

**Key Principle:** Immunity prevents *automated action* but not *automated detection*. All flags generate security alerts to ensure a compromised founder account is quickly identified.

### Access Levels

| Level | Role | Permissions | Key Restrictions |
|-------|------|-------------|------------------|
| **10** | **Supreme Founder (Neo Qiss only)** | • Full platform control<br>• User impersonation<br>• Algorithm override<br>• Delete any content<br>• Access all private data<br>• Emergency kill switches | None (all features) |
| **9** | **C-Suite Executives** | • User moderation (ban/mute)<br>• Content deletion<br>• Platform analytics<br>• Feature flags (review only)<br>• Broadcast notifications (approval required) | ❌ Cannot impersonate users<br>❌ Cannot override algorithm<br>❌ Cannot access private messages<br>❌ Cannot use emergency kill switches |
| **8** | **VP/Director Level** | • Content moderation<br>• User reports review<br>• Analytics dashboard<br>• Feature flag viewing | ❌ Cannot ban users (only recommend)<br>❌ Cannot delete accounts<br>❌ Cannot access sensitive user data<br>❌ Cannot send platform-wide notifications |
| **5-7** | **Senior Engineers/Managers** | • Debug tools<br>• Feature testing<br>• Analytics (aggregated only)<br>• Bug reproduction tools | ❌ No user moderation powers<br>❌ No access to user PII<br>❌ Cannot modify production data |
| **1-4** | **Junior Staff/Contractors** | • Staging environment access<br>• Read-only analytics<br>• Bug reporting tools | ❌ No production access<br>❌ No user data access<br>❌ No moderation capabilities |

---

## 🎛️ Admin Panel Access

### Mobile App Access Methods

#### Method 1: Discreet Gesture Access
- **iOS:** Triple-tap your profile picture
- **Android:** Triple-tap your profile picture
- **Fallback:** Long-press settings icon for 3 seconds

**Note:** These gestures are "discreet" for operational convenience, not security through obscurity. The actual security is enforced by the biometric gate. These methods may be known to senior team members for support purposes.

#### Method 2: Settings Menu (Founder Quick-Launch)
1. Go to Settings
2. Scroll to bottom
3. Tap version number 7 times
4. "Admin Panel" option appears

#### Security Gate (Primary Defense)
- **Biometric Required:** Face ID (iOS) / Fingerprint (Android)
- **Session Timeout:** 15 minutes of inactivity
- **Device Whitelist:** Only registered devices
- **Audit Logging:** Every access logged
- **Location Verification:** Unusual locations trigger additional verification

---

## 🛠️ Admin Features

### 1. User Management

#### Actions Available
- **View User Profile**
  - Full profile data (including private info)
  - IP address history
  - Device information
  - Login sessions
  - Security events
  - Verification status
  
- **User Moderation**
  - Shadowban/Unshadowban
  - Ban/Unban (temporary or permanent)
  - Mute/Unmute platform-wide
  - Force password reset
  - Force logout all sessions
  - Disable 2FA
  - Delete account
  
- **User Impersonation**
  - View platform as any user
  - Debug user experience issues
  - Test features for specific users
  - **⚠️ CRITICAL SECURITY PROTOCOL:**
    - Requires separate biometric re-authentication beyond initial admin panel access
    - Additional password verification required for impersonation
    - Every action taken while impersonating is logged with:
      - Original user ID (neoqiss_uuid)
      - Impersonated user ID
      - Timestamp
      - Action type
      - IP address
      - Device ID
    - Impersonation sessions automatically terminate after 10 minutes
    - Cannot impersonate while already impersonating (no nested impersonation)
    - High-priority security alert sent to CTO when impersonation begins
    - **Legal Requirement:** Impersonation reason must be documented before access granted
  
- **Bulk Actions**
  - Ban by IP range
  - Ban by device fingerprint
  - Ban by email domain

#### Access Path
```
Admin Panel → Users → Search User → User Actions
```

---

### 2. Content Moderation

#### Actions Available
- **Single Content Actions**
  - Delete any post/thread/take/story/comment
  - Edit content (with audit trail)
  - Pin/Unpin content
  - Feature on Explore/Trending
  - Mark as inappropriate
  - Remove from algorithm
  - Restore deleted content
  
- **Bulk Content Actions**
  - Delete by keyword/hashtag
  - Delete by user
  - Delete by date range
  - Bulk takedown (e.g., DMCA)
  
- **Content Visibility**
  - Override algorithm ranking
  - Boost/Suppress content reach
  - Add to trending manually
  - Remove from trending

#### Access Path
```
Admin Panel → Content → Search/Browse → Content Actions
```

---

### 3. Platform Control

#### System Management
- **Feature Flags**
  - Enable/disable features globally
  - Enable/disable per user
  - A/B test assignment override
  - Beta access control
  
- **Notifications**
  - Send push notification to any user
  - Broadcast to all users
  - Scheduled announcements
  - Emergency alerts
  
- **Maintenance Mode**
  - Enable/disable posting
  - Read-only mode
  - Full app lockdown
  - Custom maintenance message

#### Emergency Controls
- **Kill Switches**
  - Disable all posting
  - Disable all commenting
  - Disable direct messages
  - Disable algorithm (show chronological)
  - Emergency ban user
  
- **Rollback Functions**
  - Revert algorithm changes
  - Restore deleted content
  - Undo mass actions

#### Access Path
```
Admin Panel → Platform → System Controls
```

---

### 4. Analytics & Monitoring

#### Live Metrics
- Active users (real-time)
- Posts per second
- Engagement rate
- Server health
- API latency
- Error rates
- Database performance

#### Content Analytics
- Trending content (real-time)
- Viral prediction scores
- Engagement breakdown
- User growth metrics
- Retention analytics
- Revenue metrics (premium subscriptions)

#### Search & Discovery
- Search any content platform-wide (including private)
- Search users by any field
- View deleted content
- Access moderation queue
- Review auto-mod decisions

#### Access Path
```
Admin Panel → Analytics → Dashboard
```

---

### 5. Testing & Development

#### Feature Access
- **Beta Features**
  - Access all unreleased features
  - Toggle features on/off for yourself
  - Test A/B variants
  
- **Environment Toggle**
  - Switch between Production/Staging/Dev
  - Test internal builds
  - Preview upcoming releases
  
- **Debug Mode**
  - View API responses
  - See error logs in real-time
  - Network request inspector
  - Performance profiler

#### Testing Tools
- Push notification tester
- Deep link tester
- Share intent tester
- Algorithm simulator

#### Access Path
```
Admin Panel → Developer → Testing Tools
```

---

### 6. Security & Audit

#### Security Tools
- **Session Management**
  - View all active sessions
  - Kill any session remotely
  - Require re-authentication
  
- **2FA Management**
  - Generate backup codes
  - Disable 2FA for recovery
  - View 2FA logs
  
- **IP Management**
  - Whitelist IPs
  - Block IP ranges
  - View login locations

#### Audit Logs
- All admin actions logged
- User actions history
- Content modification history
- System changes log
- Security events log

#### Access Path
```
Admin Panel → Security → Audit Logs
```

---

## 📱 Platform-Specific Features

### Socialink Admin Features
- **Pods Management**
  - Feature pods on trending
  - Moderate pod challenges
  - Delete viral pods
  
- **Groups & Events**
  - Delete groups
  - Remove event organizers
  - Feature events
  
- **Marketplace**
  - Remove listings
  - Ban sellers
  - Refund transactions
  
- **Live Streams**
  - End live streams
  - Ban from streaming
  - Moderate live comments

### Vignette Admin Features
- **Takes Management**
  - Feature takes on For You page
  - Moderate challenges
  - Delete viral takes
  
- **Behind-The-Takes (BTT)**
  - Moderate BTT content
  - Feature BTTs
  - Remove inappropriate BTTs
  
- **Challenges**
  - Create official challenges
  - End challenges early
  - Feature challenge winners

### Sonet Admin Features
- **Threads Management**
  - Pin threads globally
  - Feature threads on trending
  - Moderate thread chains
  
- **Spaces**
  - Join any space as listener/speaker
  - End spaces remotely
  - Ban from hosting spaces
  
- **Circles & Lists**
  - View private circles
  - Moderate list content
  - Remove inappropriate circles

---

## 🔒 Security Architecture

### Multi-Layer Security

#### Layer 1: Account Verification
- Username match: `neoqiss`
- Account flag: `is_founder: true`
- Account ID: Hardcoded UUID

#### Layer 2: Device Authentication
- Device fingerprinting
- Registered device whitelist
- Device certificate validation

#### Layer 3: Biometric Lock
- Face ID (iOS)
- Touch ID (iOS)
- Fingerprint (Android)
- Required for every admin action

#### Layer 4: Session Management
- 15-minute session timeout
- Re-authentication for destructive actions
- Location-based session validation

#### Layer 5: Audit Trail
- Every action logged with:
  - Timestamp
  - Device ID
  - IP address
  - Action type
  - Target (user/content ID)
  - Result (success/failure)

### Backup Access Methods

#### Emergency CLI Access
```bash
# From your development machine
./entativa-admin --env production --action ban-user --user-id {id}
```

#### Web Admin Panel
- URL: `https://admin.entativa.com`
- Requires: Email + Password + 2FA + Biometric
- Use for: Bulk operations, analytics, system management

#### Database Direct Access
- Only for critical emergencies
- Requires VPN + SSH key + 2FA
- All queries logged

---

## 🎨 UI/UX Design

### Admin Panel Layout

#### Main Navigation
```
┌─────────────────────────────────┐
│  👑 Admin Panel                 │
├─────────────────────────────────┤
│  📊 Dashboard                   │
│  👥 Users                       │
│  📝 Content                     │
│  ⚙️  Platform                   │
│  📈 Analytics                   │
│  🔧 Developer                   │
│  🔒 Security                    │
│  📋 Audit Logs                  │
└─────────────────────────────────┘
```

#### Quick Actions Bar (Context Menu)
When viewing any content or profile, long-press to see:
```
┌─────────────────────────────────┐
│  ⚡ Quick Admin Actions         │
├─────────────────────────────────┤
│  🚫 Ban User                    │
│  👻 Shadowban                   │
│  🗑️  Delete Content             │
│  📌 Pin/Feature                 │
│  👁️  View Analytics             │
│  📋 View Audit Log              │
└─────────────────────────────────┘
```

#### Admin Toolbar (Floating)
Appears at bottom of screen when admin mode is active:
```
┌─────────────────────────────────┐
│  [👁️ View] [🗑️ Delete] [📊 Stats] │
└─────────────────────────────────┘
```

---

## 🏗️ Technical Implementation

### Backend Architecture

#### New Microservice: `admin-service`
```
services/admin-service/
├── cmd/api/
│   └── main.go
├── internal/
│   ├── handler/
│   │   ├── user_management_handler.go
│   │   ├── content_moderation_handler.go
│   │   ├── platform_control_handler.go
│   │   ├── analytics_handler.go
│   │   └── audit_handler.go
│   ├── service/
│   │   ├── admin_service.go
│   │   ├── audit_service.go
│   │   ├── security_service.go
│   │   └── notification_service.go
│   ├── middleware/
│   │   ├── founder_auth_middleware.go
│   │   ├── audit_middleware.go
│   │   └── rate_limit_middleware.go
│   ├── repository/
│   │   ├── admin_repository.go
│   │   └── audit_repository.go
│   └── model/
│       ├── admin_action.go
│       └── audit_log.go
├── pkg/
│   └── database/
│       └── postgres.go
└── test/
    └── integration_test.go
```

#### API Endpoints
```
# User Management
POST   /api/admin/users/{id}/ban
POST   /api/admin/users/{id}/shadowban
POST   /api/admin/users/{id}/mute
DELETE /api/admin/users/{id}/permanently-delete  # Irreversible account deletion
POST   /api/admin/users/{id}/suspend              # Reversible suspension
GET    /api/admin/users/{id}
POST   /api/admin/users/{id}/impersonate

# Content Moderation
DELETE /api/admin/content/{id}
PUT    /api/admin/content/{id}/edit
POST   /api/admin/content/{id}/feature
POST   /api/admin/content/{id}/pin
GET    /api/admin/content/moderation-queue

# Platform Control
POST   /api/admin/notifications/broadcast
POST   /api/admin/features/toggle
POST   /api/admin/maintenance/enable
GET    /api/admin/platform/health

# Analytics
GET    /api/admin/analytics/live
GET    /api/admin/analytics/trending
GET    /api/admin/analytics/users

# Security & Audit
GET    /api/admin/audit/logs
GET    /api/admin/security/sessions
POST   /api/admin/security/sessions/{id}/kill
```

**Critical Endpoint Notes:**
- **`DELETE /api/admin/users/{id}/permanently-delete`**
  - Permanently deletes user account and all associated data
  - **Irreversible** - no recovery possible
  - Requires double confirmation + reason documentation
  - Triggers immediate legal team notification (GDPR/CCPA compliance)
  - Audit log entry marked as CRITICAL with full context
  
- **`POST /api/admin/users/{id}/suspend`**
  - Temporary suspension (can be reversed)
  - Use this for most moderation cases instead of permanent deletion
  - Can set expiration time or require manual review for reinstatement

#### Middleware Protection
```go
// founder_auth_middleware.go
func FounderAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := getUserFromContext(r)
        
        // Verify founder account
        if user.Username != "neoqiss" || !user.IsFounder {
            http.Error(w, "Forbidden: Founder access required", http.StatusForbidden)
            logSecurityEvent("unauthorized_admin_access", user)
            return
        }
        
        // Verify device whitelist
        if !isDeviceWhitelisted(r) {
            http.Error(w, "Forbidden: Device not recognized", http.StatusForbidden)
            logSecurityEvent("unregistered_device", user)
            return
        }
        
        // Log admin action
        auditLog.Record(AdminAction{
            UserID: user.ID,
            Action: r.Method + " " + r.URL.Path,
            Timestamp: time.Now(),
            IPAddress: getIPAddress(r),
            DeviceID: getDeviceID(r),
        })
        
        next.ServeHTTP(w, r)
    })
}
```

### Mobile Implementation

#### iOS (Swift)
```swift
// AdminManager.swift
class AdminManager {
    static let shared = AdminManager()
    private let biometricAuth = BiometricAuthManager()
    
    func isFounderAccount() -> Bool {
        guard let user = SessionManager.shared.currentUser else { return false }
        return user.username == "neoqiss" && user.isFounder == true
    }
    
    func showAdminPanel() {
        guard isFounderAccount() else {
            print("Not authorized for admin access")
            return
        }
        
        // Require biometric authentication
        biometricAuth.authenticate(reason: "Access Admin Panel") { success, error in
            if success {
                self.presentAdminPanel()
            } else {
                self.showError("Biometric authentication failed")
            }
        }
    }
    
    private func presentAdminPanel() {
        let adminVC = AdminPanelViewController()
        adminVC.modalPresentationStyle = .fullScreen
        
        if let topVC = UIApplication.topViewController() {
            topVC.present(adminVC, animated: true)
        }
    }
    
    // Quick admin actions
    func quickBanUser(userId: String, completion: @escaping (Bool) -> Void) {
        biometricAuth.authenticate(reason: "Ban User") { success, _ in
            guard success else {
                completion(false)
                return
            }
            
            AdminAPIClient.shared.banUser(userId: userId) { result in
                completion(result.isSuccess)
            }
        }
    }
    
    // User impersonation with step-up authentication
    func impersonateUser(userId: String, reason: String, completion: @escaping (Bool) -> Void) {
        // First biometric check
        biometricAuth.authenticate(reason: "Impersonate User - Step 1") { success, _ in
            guard success else {
                completion(false)
                return
            }
            
            // Second authentication - password verification
            self.promptForPassword { password in
                guard let password = password else {
                    completion(false)
                    return
                }
                
                // Verify reason is documented
                guard !reason.isEmpty, reason.count >= 20 else {
                    self.showError("Impersonation reason must be at least 20 characters")
                    completion(false)
                    return
                }
                
                // Send impersonation request with full audit trail
                AdminAPIClient.shared.impersonateUser(
                    userId: userId,
                    password: password,
                    reason: reason,
                    deviceId: DeviceManager.deviceId,
                    timestamp: Date()
                ) { result in
                    if result.isSuccess {
                        // Set 10-minute timer for auto-termination
                        self.startImpersonationTimer(userId: userId)
                        // Alert CTO
                        self.notifySecurityTeam(action: "impersonation_started", userId: userId)
                    }
                    completion(result.isSuccess)
                }
            }
        }
    }
}

// Usage in app
// Triple-tap gesture on profile picture
profileImageView.addGestureRecognizer(
    UITapGestureRecognizer(target: self, action: #selector(handleTripleTap))
)

@objc func handleTripleTap(_ gesture: UITapGestureRecognizer) {
    if gesture.numberOfTapsRequired == 3 {
        AdminManager.shared.showAdminPanel()
    }
}
```

#### Android (Kotlin)
```kotlin
// AdminManager.kt
object AdminManager {
    private val biometricPrompt by lazy { createBiometricPrompt() }
    
    fun isFounderAccount(): Boolean {
        val user = SessionManager.currentUser ?: return false
        return user.username == "neoqiss" && user.isFounder == true
    }
    
    fun showAdminPanel(activity: FragmentActivity) {
        if (!isFounderAccount()) {
            Log.w("AdminManager", "Not authorized for admin access")
            return
        }
        
        // Require biometric authentication
        biometricPrompt.authenticate(
            BiometricPrompt.PromptInfo.Builder()
                .setTitle("Admin Panel Access")
                .setSubtitle("Authenticate to continue")
                .setNegativeButtonText("Cancel")
                .build(),
            callback = { result ->
                if (result.isSuccess) {
                    launchAdminPanel(activity)
                } else {
                    showError("Authentication failed")
                }
            }
        )
    }
    
    private fun launchAdminPanel(activity: FragmentActivity) {
        val intent = Intent(activity, AdminPanelActivity::class.java)
        activity.startActivity(intent)
    }
    
    // Quick admin actions
    suspend fun quickBanUser(userId: String): Result<Unit> {
        return withContext(Dispatchers.IO) {
            AdminApiClient.banUser(userId)
        }
    }
}

// Usage in app
// Triple-tap gesture on profile picture
profileImageView.setOnClickListener(object : OnMultiTapListener(3) {
    override fun onMultiTap() {
        AdminManager.showAdminPanel(requireActivity())
    }
})
```

---

## 📋 Usage Guidelines

### When to Use Admin Features

#### ✅ Appropriate Use Cases
- **Emergency Response:** Imminent harm, illegal content, safety threats
- **Critical Bugs:** Features broken for large user base
- **Testing:** Verify new features work correctly
- **Support:** Help users with account issues
- **Moderation:** Enforce community guidelines
- **Analytics:** Monitor platform health and growth

#### ❌ Inappropriate Use Cases
- Personal disputes or arguments
- Favoring friends/business partners
- Censoring unpopular opinions
- Manipulating metrics artificially
- Accessing private data without reason

### Best Practices

1. **Document Everything:** Note why you took each action
2. **Minimal Intervention:** Use least disruptive action possible
3. **Transparency:** Consider informing users of actions (when safe)
4. **Review Audit Logs:** Periodically check your own actions
5. **Team Consultation:** Discuss major actions with team when possible

---

## 🚨 Emergency Procedures

### Critical Incident Response

#### Level 1: Minor Issue
- Single user causing problems
- Isolated inappropriate content
- **Action:** Use standard moderation tools
- **Response Time:** Within 24 hours

#### Level 2: Moderate Issue
- Multiple users affected
- Spreading misinformation
- **Action:** Shadowban, content removal, notifications
- **Response Time:** Within 1 hour

#### Level 3: Severe Issue
- Platform-wide bug
- Viral harmful content
- Security breach
- **Action:** Feature flags, maintenance mode, mass deletion
- **Response Time:** Immediate

#### Level 4: Critical Emergency
- Active shooter situation
- CSAM content
- Coordinated attack
- **Action:** Kill switches, law enforcement contact, platform lockdown
- **Response Time:** Immediate + team mobilization

### Emergency Contacts
```
CTO: [Phone] [Email]
Head of Security: [Phone] [Email]
Head of Legal: [Phone] [Email]
24/7 Engineering: [Phone] [Slack Channel]
```

---

## 🔄 Rollout Plan

### Phase 1: Foundation (Weeks 1-2)
- [ ] Build `admin-service` microservice
- [ ] Add `is_founder` flag to user database
- [ ] Create admin API endpoints
- [ ] Implement founder auth middleware
- [ ] Set up audit logging

### Phase 2: Basic Features (Weeks 3-4)
- [ ] User management (ban, shadowban, mute)
- [ ] Content deletion
- [ ] Basic analytics dashboard
- [ ] Audit log viewer

### Phase 3: Mobile Integration (Weeks 5-6)
- [ ] iOS admin panel UI
- [ ] Android admin panel UI
- [ ] Secret gesture access
- [ ] Biometric authentication
- [ ] Quick action context menus

### Phase 4: Advanced Features (Weeks 7-8)
- [ ] User impersonation
- [ ] Algorithm override
- [ ] Feature flags management
- [ ] Broadcast notifications
- [ ] Maintenance mode

### Phase 5: Platform-Specific (Weeks 9-10)
- [ ] Socialink-specific admin features
- [ ] Vignette-specific admin features
- [ ] Sonet-specific admin features

### Phase 6: Testing & Security (Weeks 11-12)
- [ ] Penetration testing
- [ ] Audit log review
- [ ] Device whitelisting
- [ ] Emergency procedures testing
- [ ] Team training

### Phase 7: Launch (Week 13)
- [ ] Production deployment
- [ ] Monitoring setup
- [ ] Documentation complete
- [ ] Founder account activation

---

## 📚 Reference Links

### Internal Documentation
- API Documentation: `docs/api/admin-service.md`
- Security Guidelines: `docs/security/admin-access.md`
- Audit Log Schema: `docs/database/audit-logs.md`

### External Resources
- [GDPR Compliance for Admin Access](https://gdpr.eu)
- [CCPA Guidelines](https://oag.ca.gov/privacy/ccpa)
- [Platform Security Best Practices](https://owasp.org)

---

## 🔧 Troubleshooting

### Admin Panel Won't Open
1. Verify you're logged in as `@neoqiss`
2. Check biometric authentication is enabled
3. Verify device is whitelisted
4. Check network connection
5. Clear app cache and retry

### Actions Failing
1. Check audit logs for error messages
2. Verify API service is running
3. Check rate limits (even with immunity, safety limits exist)
4. Contact backend team if persistent

### Biometric Auth Issues
1. Re-register biometrics in device settings
2. Use backup password authentication
3. Contact security team to whitelist new device

---

## 📞 Support

### For Technical Issues
- **Slack:** #founder-admin-support
- **Email:** admin-support@entativa.com
- **Emergency:** [On-call engineer phone]

### For Security Concerns
- **Slack:** #security-incidents
- **Email:** security@entativa.com
- **Emergency:** [Security team lead phone]

---

## 📝 Changelog

### Version 1.0.0 (2025-10-17)
- Initial document creation
- Defined founder privileges
- Outlined admin features
- Created implementation plan

---

## ⚖️ Legal & Compliance

### Data Access Responsibility
As the founder with access to all user data, you are responsible for:
- GDPR compliance (EU users)
- CCPA compliance (California users)
- Ethical use of private information
- Data retention policies
- Breach notification procedures

### Audit Requirements
- All admin actions are logged indefinitely
- Logs may be subpoenaed in legal cases
- Quarterly security audits of admin access
- Annual compliance review

---

## 🎯 Success Metrics

### Admin Panel Usage
- Access frequency
- Most-used features
- Response time to critical incidents
- False positive rate (unnecessary actions)

### Platform Health
- Time to resolve critical incidents
- User trust metrics
- Moderation accuracy
- Appeal success rate

---

**Document Owner:** Neo Qiss (@neoqiss)  
**Next Review Date:** 2025-11-17  
**Classification:** Confidential - Founder Eyes Only

---

*End of Document*

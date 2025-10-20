import Foundation
import LocalAuthentication
import SwiftUI

// MARK: - Admin Manager (Founder Access Control)
class AdminManager: ObservableObject {
    static let shared = AdminManager()
    
    @Published var isAdminMode = false
    @Published var adminSessionExpiry: Date?
    
    private let sessionTimeout: TimeInterval = 15 * 60 // 15 minutes
    private var sessionTimer: Timer?
    
    private init() {}
    
    // MARK: - Founder Check
    
    /// Verifies if current user is the founder (@neoqiss)
    func isFounderAccount() -> Bool {
        guard let token = KeychainManager.shared.getToken() else {
            return false
        }
        
        // Decode JWT to check username and is_founder flag
        // For now, check username (in production, verify JWT payload)
        return checkFounderInToken(token)
    }
    
    private func checkFounderInToken(_ token: String) -> Bool {
        // Decode JWT payload
        let segments = token.components(separatedBy: ".")
        guard segments.count == 3 else { return false }
        
        // Base64 decode payload
        var base64 = segments[1]
        base64 = base64.replacingOccurrences(of: "-", with: "+")
            .replacingOccurrences(of: "_", with: "/")
        
        while base64.count % 4 != 0 {
            base64.append("=")
        }
        
        guard let data = Data(base64Encoded: base64),
              let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any] else {
            return false
        }
        
        // Check username and is_founder flag
        let username = json["username"] as? String
        let isFounder = json["is_founder"] as? Bool ?? false
        
        return username == "neoqiss" && isFounder
    }
    
    // MARK: - Admin Panel Access
    
    /// Show admin panel with biometric authentication
    func showAdminPanel(from viewController: UIViewController? = nil) {
        guard isFounderAccount() else {
            print("⚠️ Not authorized for admin access")
            return
        }
        
        // Require biometric authentication
        authenticateWithBiometric(reason: "Access Admin Panel") { [weak self] success in
            if success {
                self?.presentAdminPanel(from: viewController)
            } else {
                self?.showError(message: "Biometric authentication failed")
            }
        }
    }
    
    private func presentAdminPanel(from viewController: UIViewController?) {
        DispatchQueue.main.async {
            self.isAdminMode = true
            self.startAdminSession()
            
            // Present admin panel
            let adminView = AdminPanelView()
            let hostingController = UIHostingController(rootView: adminView)
            hostingController.modalPresentationStyle = .fullScreen
            
            if let vc = viewController ?? UIApplication.topViewController() {
                vc.present(hostingController, animated: true)
            }
        }
    }
    
    // MARK: - Session Management
    
    private func startAdminSession() {
        adminSessionExpiry = Date().addingTimeInterval(sessionTimeout)
        
        // Invalidate existing timer
        sessionTimer?.invalidate()
        
        // Start new session timer
        sessionTimer = Timer.scheduledTimer(withTimeInterval: sessionTimeout, repeats: false) { [weak self] _ in
            self?.endAdminSession()
        }
    }
    
    func endAdminSession() {
        isAdminMode = false
        adminSessionExpiry = nil
        sessionTimer?.invalidate()
        sessionTimer = nil
    }
    
    func extendAdminSession() {
        guard isAdminMode else { return }
        
        authenticateWithBiometric(reason: "Extend Admin Session") { [weak self] success in
            if success {
                self?.startAdminSession()
            } else {
                self?.endAdminSession()
            }
        }
    }
    
    // MARK: - Quick Actions
    
    /// Ban a user with biometric confirmation
    func quickBanUser(userId: String, reason: String, completion: @escaping (Bool) -> Void) {
        authenticateWithBiometric(reason: "Ban User") { [weak self] success in
            guard success else {
                completion(false)
                return
            }
            
            self?.performBanUser(userId: userId, reason: reason, completion: completion)
        }
    }
    
    private func performBanUser(userId: String, reason: String, completion: @escaping (Bool) -> Void) {
        guard let token = KeychainManager.shared.getToken() else {
            completion(false)
            return
        }
        
        let url = URL(string: "\(APIConfig.adminBaseURL)/users/\(userId)/ban")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue(getDeviceID(), forHTTPHeaderField: "X-Device-ID")
        
        let body = ["reason": reason, "duration": 0] // 0 = permanent
        request.httpBody = try? JSONEncoder().encode(body)
        
        URLSession.shared.dataTask(with: request) { data, response, error in
            DispatchQueue.main.async {
                completion(error == nil)
            }
        }.resume()
    }
    
    /// Shadowban a user
    func quickShadowbanUser(userId: String, reason: String, completion: @escaping (Bool) -> Void) {
        authenticateWithBiometric(reason: "Shadowban User") { [weak self] success in
            guard success else {
                completion(false)
                return
            }
            
            self?.performShadowban(userId: userId, reason: reason, completion: completion)
        }
    }
    
    private func performShadowban(userId: String, reason: String, completion: @escaping (Bool) -> Void) {
        guard let token = KeychainManager.shared.getToken() else {
            completion(false)
            return
        }
        
        let url = URL(string: "\(APIConfig.adminBaseURL)/users/\(userId)/shadowban")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue(getDeviceID(), forHTTPHeaderField: "X-Device-ID")
        
        let body = ["reason": reason]
        request.httpBody = try? JSONEncoder().encode(body)
        
        URLSession.shared.dataTask(with: request) { data, response, error in
            DispatchQueue.main.async {
                completion(error == nil)
            }
        }.resume()
    }
    
    // MARK: - User Impersonation (Step-Up Authentication)
    
    /// Impersonate a user (requires password + biometric)
    func impersonateUser(userId: String, reason: String, password: String, completion: @escaping (Bool, String?) -> Void) {
        // Require detailed reason (min 20 chars)
        guard reason.count >= 20 else {
            completion(false, "Reason must be at least 20 characters")
            return
        }
        
        // First biometric check
        authenticateWithBiometric(reason: "Impersonate User - Step 1") { [weak self] success in
            guard success else {
                completion(false, "Biometric authentication failed")
                return
            }
            
            // Second authentication - password verification
            self?.verifyPasswordAndImpersonate(
                userId: userId,
                reason: reason,
                password: password,
                completion: completion
            )
        }
    }
    
    private func verifyPasswordAndImpersonate(
        userId: String,
        reason: String,
        password: String,
        completion: @escaping (Bool, String?) -> Void
    ) {
        guard let token = KeychainManager.shared.getToken() else {
            completion(false, "Not authenticated")
            return
        }
        
        let url = URL(string: "\(APIConfig.adminBaseURL)/users/\(userId)/impersonate")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue(getDeviceID(), forHTTPHeaderField: "X-Device-ID")
        
        let body: [String: Any] = [
            "reason": reason,
            "password": password
        ]
        request.httpBody = try? JSONSerialization.data(withJSONObject: body)
        
        URLSession.shared.dataTask(with: request) { data, response, error in
            guard let data = data,
                  let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
                  let impersonationToken = json["impersonation_token"] as? String else {
                DispatchQueue.main.async {
                    completion(false, "Failed to create impersonation session")
                }
                return
            }
            
            DispatchQueue.main.async {
                // Start 10-minute auto-termination timer
                self.startImpersonationTimer(userId: userId, token: impersonationToken)
                completion(true, impersonationToken)
            }
        }.resume()
    }
    
    private func startImpersonationTimer(userId: String, token: String) {
        // Auto-terminate impersonation after 10 minutes
        Timer.scheduledTimer(withTimeInterval: 10 * 60, repeats: false) { [weak self] _ in
            self?.endImpersonation(userId: userId, token: token)
        }
    }
    
    private func endImpersonation(userId: String, token: String) {
        guard let authToken = KeychainManager.shared.getToken() else { return }
        
        let url = URL(string: "\(APIConfig.adminBaseURL)/users/\(userId)/end-impersonation")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(authToken)", forHTTPHeaderField: "Authorization")
        request.setValue(getDeviceID(), forHTTPHeaderField: "X-Device-ID")
        
        URLSession.shared.dataTask(with: request).resume()
    }
    
    // MARK: - Helpers
    
    private func authenticateWithBiometric(reason: String, completion: @escaping (Bool) -> Void) {
        let context = LAContext()
        var error: NSError?
        
        guard context.canEvaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, error: &error) else {
            completion(false)
            return
        }
        
        context.evaluatePolicy(
            .deviceOwnerAuthenticationWithBiometrics,
            localizedReason: reason
        ) { success, error in
            DispatchQueue.main.async {
                completion(success)
            }
        }
    }
    
    private func showError(message: String) {
        // Show error alert
        print("❌ Admin Error: \(message)")
    }
    
    private func getDeviceID() -> String {
        return UIDevice.current.identifierForVendor?.uuidString ?? "unknown"
    }
}

// MARK: - API Config Extension
extension APIConfig {
    static var adminBaseURL: String {
        return "http://localhost:8005/api/admin" // Admin service on port 8005
    }
}

// MARK: - UIApplication Extension (Helper)
extension UIApplication {
    static func topViewController(base: UIViewController? = nil) -> UIViewController? {
        let base = base ?? UIApplication.shared.connectedScenes
            .compactMap { $0 as? UIWindowScene }
            .flatMap { $0.windows }
            .first { $0.isKeyWindow }?.rootViewController
        
        if let nav = base as? UINavigationController {
            return topViewController(base: nav.visibleViewController)
        }
        
        if let tab = base as? UITabBarController, let selected = tab.selectedViewController {
            return topViewController(base: selected)
        }
        
        if let presented = base?.presentedViewController {
            return topViewController(base: presented)
        }
        
        return base
    }
}

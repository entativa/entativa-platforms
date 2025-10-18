import SwiftUI
import UIKit

/// Authentication Coordinator for Entativa
/// Manages authentication flow and navigation
class AuthCoordinator: ObservableObject {
    @Published var isAuthenticated = false
    @Published var showLogin = true
    
    func showLoginScreen() {
        showLogin = true
    }
    
    func showSignUpScreen() {
        showLogin = false
    }
    
    func handleAuthenticationSuccess() {
        isAuthenticated = true
    }
    
    func logout() {
        isAuthenticated = false
        showLogin = true
    }
}

/// Authentication Root View
struct AuthCoordinatorView: View {
    @StateObject private var coordinator = AuthCoordinator()
    @StateObject private var viewModel = AuthViewModel()
    
    var body: some View {
        Group {
            if viewModel.isAuthenticated || coordinator.isAuthenticated {
                // Navigate to main app (placeholder)
                Text("Welcome to Entativa!")
                    .font(.largeTitle)
            } else {
                EntativaLoginView()
            }
        }
        .onChange(of: viewModel.isAuthenticated) { oldValue, newValue in
            if newValue {
                coordinator.handleAuthenticationSuccess()
            }
        }
    }
}

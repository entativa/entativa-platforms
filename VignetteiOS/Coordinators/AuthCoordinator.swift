import SwiftUI
import UIKit

/// Authentication Coordinator for Vignette
/// Manages authentication flow and navigation
class VignetteAuthCoordinator: ObservableObject {
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
struct VignetteAuthCoordinatorView: View {
    @StateObject private var coordinator = VignetteAuthCoordinator()
    @StateObject private var viewModel = VignetteAuthViewModel()
    
    var body: some View {
        Group {
            if viewModel.isAuthenticated || coordinator.isAuthenticated {
                // Navigate to main app (placeholder)
                Text("Welcome to Vignette!")
                    .font(.largeTitle)
            } else {
                VignetteLoginView()
            }
        }
        .onChange(of: viewModel.isAuthenticated) { oldValue, newValue in
            if newValue {
                coordinator.handleAuthenticationSuccess()
            }
        }
    }
}

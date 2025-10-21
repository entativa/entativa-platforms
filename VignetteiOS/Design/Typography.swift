import SwiftUI

/// Enterprise-grade typography system for Vignette
/// Uses SF Pro (system font) with Instagram-inspired scale and weights
struct VignetteTypography {
    
    // MARK: - Display Styles (Splash screens, hero sections)
    static let displayLarge = Font.system(size: 52, weight: .bold, design: .default)
    static let displayMedium = Font.system(size: 42, weight: .bold, design: .default)
    static let displaySmall = Font.system(size: 34, weight: .semibold, design: .default)
    
    // MARK: - Headline Styles (Profile names, important headers)
    static let headlineLarge = Font.system(size: 30, weight: .bold)
    static let headlineMedium = Font.system(size: 26, weight: .semibold)
    static let headlineSmall = Font.system(size: 22, weight: .semibold)
    
    // MARK: - Title Styles (Post captions, section titles)
    static let titleLarge = Font.system(size: 20, weight: .semibold)
    static let titleMedium = Font.system(size: 17, weight: .semibold)
    static let titleSmall = Font.system(size: 15, weight: .medium)
    
    // MARK: - Body Styles (Comments, descriptions)
    static let bodyLarge = Font.system(size: 16, weight: .regular)
    static let bodyMedium = Font.system(size: 14, weight: .regular)
    static let bodySmall = Font.system(size: 12, weight: .regular)
    
    // MARK: - Label Styles (Usernames, form labels)
    static let labelLarge = Font.system(size: 15, weight: .semibold)
    static let labelMedium = Font.system(size: 13, weight: .semibold)
    static let labelSmall = Font.system(size: 11, weight: .semibold)
    
    // MARK: - Button Styles
    static let buttonLarge = Font.system(size: 16, weight: .semibold)
    static let buttonMedium = Font.system(size: 14, weight: .semibold)
    static let buttonSmall = Font.system(size: 12, weight: .medium)
    
    // MARK: - Caption Styles (Timestamps, metadata)
    static let captionLarge = Font.system(size: 13, weight: .regular)
    static let captionMedium = Font.system(size: 11, weight: .regular)
    static let captionSmall = Font.system(size: 10, weight: .regular)
    
    // MARK: - Username Styles (Special Instagram-style username display)
    static let usernameProminant = Font.system(size: 14, weight: .semibold)
    static let usernameRegular = Font.system(size: 13, weight: .semibold)
    static let usernameSmall = Font.system(size: 12, weight: .semibold)
}

/// Typography modifiers for consistent text styling
extension View {
    func vignetteDisplayLarge() -> some View {
        self.font(VignetteTypography.displayLarge)
    }
    
    func vignetteDisplayMedium() -> some View {
        self.font(VignetteTypography.displayMedium)
    }
    
    func vignetteDisplaySmall() -> some View {
        self.font(VignetteTypography.displaySmall)
    }
    
    func vignetteHeadlineLarge() -> some View {
        self.font(VignetteTypography.headlineLarge)
    }
    
    func vignetteHeadlineMedium() -> some View {
        self.font(VignetteTypography.headlineMedium)
    }
    
    func vignetteHeadlineSmall() -> some View {
        self.font(VignetteTypography.headlineSmall)
    }
    
    func vignetteTitleLarge() -> some View {
        self.font(VignetteTypography.titleLarge)
    }
    
    func vignetteTitleMedium() -> some View {
        self.font(VignetteTypography.titleMedium)
    }
    
    func vignetteTitleSmall() -> some View {
        self.font(VignetteTypography.titleSmall)
    }
    
    func vignetteBodyLarge() -> some View {
        self.font(VignetteTypography.bodyLarge)
    }
    
    func vignetteBodyMedium() -> some View {
        self.font(VignetteTypography.bodyMedium)
    }
    
    func vignetteBodySmall() -> some View {
        self.font(VignetteTypography.bodySmall)
    }
    
    func vignetteLabelLarge() -> some View {
        self.font(VignetteTypography.labelLarge)
    }
    
    func vignetteLabelMedium() -> some View {
        self.font(VignetteTypography.labelMedium)
    }
    
    func vignetteLabelSmall() -> some View {
        self.font(VignetteTypography.labelSmall)
    }
    
    func vignetteButtonLarge() -> some View {
        self.font(VignetteTypography.buttonLarge)
    }
    
    func vignetteButtonMedium() -> some View {
        self.font(VignetteTypography.buttonMedium)
    }
    
    func vignetteButtonSmall() -> some View {
        self.font(VignetteTypography.buttonSmall)
    }
    
    func vignetteCaptionLarge() -> some View {
        self.font(VignetteTypography.captionLarge)
    }
    
    func vignetteCaptionMedium() -> some View {
        self.font(VignetteTypography.captionMedium)
    }
    
    func vignetteCaptionSmall() -> some View {
        self.font(VignetteTypography.captionSmall)
    }
    
    func vignetteUsernameProminant() -> some View {
        self.font(VignetteTypography.usernameProminant)
    }
    
    func vignetteUsernameRegular() -> some View {
        self.font(VignetteTypography.usernameRegular)
    }
    
    func vignetteUsernameSmall() -> some View {
        self.font(VignetteTypography.usernameSmall)
    }
}

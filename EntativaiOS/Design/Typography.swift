import SwiftUI

/// Enterprise-grade typography system for Entativa
/// Uses SF Pro (system font) with carefully crafted scale and weights
struct EntativaTypography {
    
    // MARK: - Display Styles (Hero sections, splash screens)
    static let displayLarge = Font.system(size: 57, weight: .bold, design: .rounded)
    static let displayMedium = Font.system(size: 45, weight: .bold, design: .rounded)
    static let displaySmall = Font.system(size: 36, weight: .semibold, design: .rounded)
    
    // MARK: - Headline Styles (Section headers, important titles)
    static let headlineLarge = Font.system(size: 32, weight: .bold)
    static let headlineMedium = Font.system(size: 28, weight: .semibold)
    static let headlineSmall = Font.system(size: 24, weight: .semibold)
    
    // MARK: - Title Styles (Card titles, list headers)
    static let titleLarge = Font.system(size: 22, weight: .semibold)
    static let titleMedium = Font.system(size: 18, weight: .semibold)
    static let titleSmall = Font.system(size: 16, weight: .medium)
    
    // MARK: - Body Styles (Main content, descriptions)
    static let bodyLarge = Font.system(size: 17, weight: .regular)
    static let bodyMedium = Font.system(size: 15, weight: .regular)
    static let bodySmall = Font.system(size: 13, weight: .regular)
    
    // MARK: - Label Styles (Form labels, captions, metadata)
    static let labelLarge = Font.system(size: 15, weight: .medium)
    static let labelMedium = Font.system(size: 13, weight: .medium)
    static let labelSmall = Font.system(size: 11, weight: .medium)
    
    // MARK: - Button Styles
    static let buttonLarge = Font.system(size: 17, weight: .semibold)
    static let buttonMedium = Font.system(size: 15, weight: .semibold)
    static let buttonSmall = Font.system(size: 13, weight: .medium)
    
    // MARK: - Caption Styles (Timestamps, helper text)
    static let captionLarge = Font.system(size: 13, weight: .regular)
    static let captionMedium = Font.system(size: 12, weight: .regular)
    static let captionSmall = Font.system(size: 11, weight: .regular)
    
    // MARK: - Overline Styles (All caps labels)
    static let overlineLarge = Font.system(size: 13, weight: .semibold)
    static let overlineMedium = Font.system(size: 11, weight: .semibold)
    static let overlineSmall = Font.system(size: 10, weight: .semibold)
}

/// Typography modifiers for consistent text styling
extension View {
    func entativaDisplayLarge() -> some View {
        self.font(EntativaTypography.displayLarge)
    }
    
    func entativaDisplayMedium() -> some View {
        self.font(EntativaTypography.displayMedium)
    }
    
    func entativaDisplaySmall() -> some View {
        self.font(EntativaTypography.displaySmall)
    }
    
    func entativaHeadlineLarge() -> some View {
        self.font(EntativaTypography.headlineLarge)
    }
    
    func entativaHeadlineMedium() -> some View {
        self.font(EntativaTypography.headlineMedium)
    }
    
    func entativaHeadlineSmall() -> some View {
        self.font(EntativaTypography.headlineSmall)
    }
    
    func entativaTitleLarge() -> some View {
        self.font(EntativaTypography.titleLarge)
    }
    
    func entativaTitleMedium() -> some View {
        self.font(EntativaTypography.titleMedium)
    }
    
    func entativaTitleSmall() -> some View {
        self.font(EntativaTypography.titleSmall)
    }
    
    func entativaBodyLarge() -> some View {
        self.font(EntativaTypography.bodyLarge)
    }
    
    func entativaBodyMedium() -> some View {
        self.font(EntativaTypography.bodyMedium)
    }
    
    func entativaBodySmall() -> some View {
        self.font(EntativaTypography.bodySmall)
    }
    
    func entativaLabelLarge() -> some View {
        self.font(EntativaTypography.labelLarge)
    }
    
    func entativaLabelMedium() -> some View {
        self.font(EntativaTypography.labelMedium)
    }
    
    func entativaLabelSmall() -> some View {
        self.font(EntativaTypography.labelSmall)
    }
    
    func entativaButtonLarge() -> some View {
        self.font(EntativaTypography.buttonLarge)
    }
    
    func entativaButtonMedium() -> some View {
        self.font(EntativaTypography.buttonMedium)
    }
    
    func entativaButtonSmall() -> some View {
        self.font(EntativaTypography.buttonSmall)
    }
    
    func entativaCaptionLarge() -> some View {
        self.font(EntativaTypography.captionLarge)
    }
    
    func entativaCaptionMedium() -> some View {
        self.font(EntativaTypography.captionMedium)
    }
    
    func entativaCaptionSmall() -> some View {
        self.font(EntativaTypography.captionSmall)
    }
}

import SwiftUI

/// Entativa Color System - Facebook-inspired design
/// Primary colors: #007CFC (Blue), #6F3EFB (Purple), #FC30E1 (Pink)
/// Vignette accent colors available for cross-brand consistency
struct EntativaColors {
    // MARK: - Brand Colors
    static let primaryBlue = Color(hex: "007CFC")
    static let primaryPurple = Color(hex: "6F3EFB")
    static let primaryPink = Color(hex: "FC30E1")
    
    // MARK: - Button Colors
    /// Primary buttons - Entativa blue
    static let buttonPrimary = primaryBlue
    static let buttonPrimaryPressed = Color(hex: "0066D4")
    static let buttonPrimaryDisabled = Color(hex: "007CFC").opacity(0.4)
    
    /// Primary deemphasis buttons - Vignette light blue with Entativa blue text
    static let buttonPrimaryDeemph = Color(hex: "C3E7F1")
    static let buttonPrimaryDeemphText = primaryBlue
    static let buttonPrimaryDeemphPressed = Color(hex: "A8D5E3")
    
    /// Secondary buttons - Monochrome
    static let buttonSecondary = Color(hex: "F0F2F5")
    static let buttonSecondaryText = Color(hex: "050505")
    static let buttonSecondaryPressed = Color(hex: "E4E6E9")
    
    // MARK: - Background Colors
    static let backgroundPrimary = Color.white
    static let backgroundSecondary = Color(hex: "F0F2F5")
    static let backgroundTertiary = Color(hex: "E4E6EB")
    
    // MARK: - Text Colors
    static let textPrimary = Color(hex: "050505")
    static let textSecondary = Color(hex: "65676B")
    static let textTertiary = Color(hex: "8A8D91")
    static let textPlaceholder = Color(hex: "BCC0C4")
    static let textOnPrimary = Color.white
    static let textLink = primaryBlue
    
    // MARK: - Border Colors
    static let borderDefault = Color(hex: "CED0D4")
    static let borderFocus = primaryBlue
    static let borderError = Color(hex: "F02849")
    
    // MARK: - Status Colors
    static let error = Color(hex: "F02849")
    static let success = Color(hex: "42B72A")
    static let warning = Color(hex: "F7B928")
    static let info = primaryBlue
    
    // MARK: - Gradient Colors
    static let gradientStart = primaryBlue
    static let gradientMiddle = primaryPurple
    static let gradientEnd = primaryPink
    
    // MARK: - Dark Mode (Future)
    static let darkBackground = Color(hex: "18191A")
    static let darkSurface = Color(hex: "242526")
    static let darkBorder = Color(hex: "3E4042")
}

/// Vignette accent colors for cross-brand consistency
struct VignetteAccentColors {
    static let lightBlue = Color(hex: "C3E7F1")
    static let moonstone = Color(hex: "519CAB")
    static let saffron = Color(hex: "FFC64F")
    static let gunmetal = Color(hex: "20373B")
}

// MARK: - Color Extension for Hex Support
extension Color {
    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let a, r, g, b: UInt64
        switch hex.count {
        case 3: // RGB (12-bit)
            (a, r, g, b) = (255, (int >> 8) * 17, (int >> 4 & 0xF) * 17, (int & 0xF) * 17)
        case 6: // RGB (24-bit)
            (a, r, g, b) = (255, int >> 16, int >> 8 & 0xFF, int & 0xFF)
        case 8: // ARGB (32-bit)
            (a, r, g, b) = (int >> 24, int >> 16 & 0xFF, int >> 8 & 0xFF, int & 0xFF)
        default:
            (a, r, g, b) = (255, 0, 0, 0)
        }
        
        self.init(
            .sRGB,
            red: Double(r) / 255,
            green: Double(g) / 255,
            blue:  Double(b) / 255,
            opacity: Double(a) / 255
        )
    }
}

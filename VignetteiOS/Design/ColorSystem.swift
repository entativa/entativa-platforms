import SwiftUI

/// Vignette Color System - Instagram-inspired design
/// Primary colors: #C3E7F1 (Light Blue), #519CAB (Moonstone), #FFC64F (Saffron), #20373B (Gunmetal)
/// Entativa blue used for primary buttons (#007CFC)
struct VignetteColors {
    // MARK: - Brand Colors
    static let lightBlue = Color(hex: "C3E7F1")
    static let moonstone = Color(hex: "519CAB")
    static let saffron = Color(hex: "FFC64F")
    static let gunmetal = Color(hex: "20373B")
    
    // MARK: - Entativa Colors (for buttons)
    static let entativaBlue = Color(hex: "007CFC")
    
    // MARK: - Button Colors
    /// Primary buttons - Entativa blue (cross-brand consistency)
    static let buttonPrimary = entativaBlue
    static let buttonPrimaryPressed = Color(hex: "0066D4")
    static let buttonPrimaryDisabled = entativaBlue.opacity(0.4)
    
    /// Primary deemphasis buttons - Vignette light blue with Entativa blue text
    static let buttonPrimaryDeemph = lightBlue
    static let buttonPrimaryDeemphText = entativaBlue
    static let buttonPrimaryDeemphPressed = Color(hex: "A8D5E3")
    
    /// Secondary buttons - Monochrome
    static let buttonSecondary = Color(hex: "FAFAFA")
    static let buttonSecondaryText = Color(hex: "262626")
    static let buttonSecondaryPressed = Color(hex: "F0F0F0")
    
    // MARK: - Background Colors
    static let backgroundPrimary = Color.white
    static let backgroundSecondary = Color(hex: "FAFAFA")
    static let backgroundTertiary = Color(hex: "F3F3F3")
    
    // MARK: - Text Colors
    static let textPrimary = Color(hex: "262626")
    static let textSecondary = Color(hex: "8E8E8E")
    static let textTertiary = Color(hex: "C7C7C7")
    static let textPlaceholder = Color(hex: "DBDBDB")
    static let textOnPrimary = Color.white
    static let textLink = moonstone
    
    // MARK: - Border Colors
    static let borderDefault = Color(hex: "DBDBDB")
    static let borderFocus = gunmetal
    static let borderError = Color(hex: "ED4956")
    
    // MARK: - Icon Colors
    static let iconPrimary = Color(hex: "262626")
    static let iconSecondary = Color(hex: "8E8E8E")
    static let iconAccent = moonstone
    
    // MARK: - Status Colors
    static let error = Color(hex: "ED4956")
    static let success = Color(hex: "0095F6")
    static let warning = saffron
    static let info = moonstone
    
    // MARK: - Gradient Colors (Instagram-style)
    static let gradientPurple = Color(hex: "833AB4")
    static let gradientPink = Color(hex: "FD1D1D")
    static let gradientOrange = Color(hex: "FCAF45")
    static let gradientYellow = saffron
    
    // MARK: - Story Ring Gradient
    static let storyGradientStart = Color(hex: "833AB4")
    static let storyGradientMiddle = Color(hex: "FD1D1D")
    static let storyGradientEnd = Color(hex: "FCAF45")
    
    // MARK: - Dark Mode
    static let darkBackground = Color(hex: "000000")
    static let darkSurface = Color(hex: "121212")
    static let darkBorder = Color(hex: "262626")
    static let darkTextPrimary = Color(hex: "FAFAFA")
    static let darkTextSecondary = Color(hex: "A8A8A8")
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

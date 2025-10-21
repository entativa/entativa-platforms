import SwiftUI

// MARK: - E2EE Message Backup Settings
struct BackupSettingsView: View {
    @Environment(\.dismiss) var dismiss
    @StateObject private var viewModel = BackupSettingsViewModel()
    @State private var showPINSetup = false
    @State private var showBackupNow = false
    @State private var showThirdPartyWarning = false
    
    var body: some View {
        NavigationView {
            List {
                // Backup Status Section
                Section {
                    HStack {
                        Image(systemName: viewModel.backupEnabled ? "checkmark.shield.fill" : "xmark.shield")
                            .font(.system(size: 32))
                            .foregroundColor(viewModel.backupEnabled ? .green : .gray)
                            .frame(width: 48)
                        
                        VStack(alignment: .leading, spacing: 4) {
                            Text(viewModel.backupEnabled ? "Backups Enabled" : "Backups Disabled")
                                .font(.system(size: 17, weight: .semibold))
                            
                            if let lastBackup = viewModel.lastBackupDate {
                                Text("Last backup: \(lastBackup, style: .relative) ago")
                                    .font(.system(size: 14))
                                    .foregroundColor(.gray)
                            } else {
                                Text("Never backed up")
                                    .font(.system(size: 14))
                                    .foregroundColor(.gray)
                            }
                        }
                    }
                    .padding(.vertical, 8)
                } header: {
                    Text("Status")
                }
                
                // Backup Location Section
                Section {
                    VStack(alignment: .leading, spacing: 12) {
                        Toggle("Enable Backups", isOn: $viewModel.backupEnabled)
                            .onChange(of: viewModel.backupEnabled) { _ in
                                viewModel.toggleBackup()
                            }
                        
                        if viewModel.backupEnabled {
                            // Our Servers (Recommended)
                            BackupLocationButton(
                                icon: "server.rack",
                                title: "Our Servers",
                                subtitle: "🔒 Most Secure - End-to-end encrypted with your PIN",
                                badge: "RECOMMENDED",
                                badgeColor: .green,
                                isSelected: viewModel.backupLocation == .ourServers,
                                action: { viewModel.selectBackupLocation(.ourServers) }
                            )
                            
                            // iCloud
                            BackupLocationButton(
                                icon: "icloud.fill",
                                title: "iCloud",
                                subtitle: "⚠️ Apple can decrypt if pressured by authorities",
                                badge: nil,
                                badgeColor: .orange,
                                isSelected: viewModel.backupLocation == .iCloud,
                                action: {
                                    showThirdPartyWarning = true
                                }
                            )
                        }
                    }
                } header: {
                    Text("Backup Location")
                } footer: {
                    Text(viewModel.backupLocation == .ourServers ?
                         "Your messages are encrypted with your PIN/passphrase. Only you can decrypt them." :
                         "⚠️ Third-party providers (Apple, Google) may decrypt your backups if pressured by authorities or at their own discretion.")
                        .foregroundColor(viewModel.backupLocation == .ourServers ? .secondary : .orange)
                }
                
                // Auto-Backup Section
                if viewModel.backupEnabled {
                    Section {
                        Toggle("Auto-Backup", isOn: $viewModel.autoBackupEnabled)
                        
                        if viewModel.autoBackupEnabled {
                            Picker("Frequency", selection: $viewModel.backupFrequency) {
                                Text("Daily").tag(BackupFrequency.daily)
                                Text("Weekly").tag(BackupFrequency.weekly)
                                Text("Monthly").tag(BackupFrequency.monthly)
                            }
                            
                            Toggle("Wi-Fi Only", isOn: $viewModel.wifiOnly)
                        }
                    } header: {
                        Text("Automatic Backups")
                    }
                }
                
                // Manual Backup Section
                if viewModel.backupEnabled {
                    Section {
                        if !viewModel.hasBackupKey {
                            Button(action: { showPINSetup = true }) {
                                HStack {
                                    Image(systemName: "key.fill")
                                        .foregroundColor(Color(hex: "007CFC"))
                                    Text("Set Up Backup PIN")
                                        .foregroundColor(.primary)
                                    Spacer()
                                    Image(systemName: "chevron.right")
                                        .foregroundColor(.gray)
                                }
                            }
                        } else {
                            Button(action: { showBackupNow = true }) {
                                HStack {
                                    Image(systemName: "arrow.up.doc.fill")
                                        .foregroundColor(Color(hex: "007CFC"))
                                    Text("Backup Now")
                                        .foregroundColor(.primary)
                                    Spacer()
                                    if viewModel.isBackingUp {
                                        ProgressView()
                                    } else {
                                        Image(systemName: "chevron.right")
                                            .foregroundColor(.gray)
                                    }
                                }
                            }
                            .disabled(viewModel.isBackingUp)
                            
                            NavigationLink(destination: BackupHistoryView(viewModel: viewModel)) {
                                HStack {
                                    Image(systemName: "clock.arrow.circlepath")
                                        .foregroundColor(Color(hex: "007CFC"))
                                    Text("Backup History")
                                }
                            }
                        }
                    } header: {
                        Text("Manual Backup")
                    }
                }
                
                // Danger Zone
                if viewModel.backupEnabled && viewModel.hasBackupKey {
                    Section {
                        Button(role: .destructive, action: { viewModel.deleteAllBackups() }) {
                            HStack {
                                Image(systemName: "trash.fill")
                                Text("Delete All Backups")
                            }
                        }
                    } header: {
                        Text("Danger Zone")
                    } footer: {
                        Text("This will permanently delete all your message backups. This cannot be undone.")
                    }
                }
            }
            .navigationTitle("Message Backups")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Done") {
                        viewModel.saveSettings()
                        dismiss()
                    }
                }
            }
            .sheet(isPresented: $showPINSetup) {
                BackupPINSetupView(viewModel: viewModel)
            }
            .sheet(isPresented: $showBackupNow) {
                BackupNowView(viewModel: viewModel)
            }
            .alert("Third-Party Backup Warning", isPresented: $showThirdPartyWarning) {
                Button("Cancel", role: .cancel) {}
                Button("I Understand", role: .destructive) {
                    viewModel.acknowledgeThirdPartyWarning()
                    viewModel.selectBackupLocation(.iCloud)
                }
            } message: {
                Text("⚠️ WARNING: Apple (iCloud) can decrypt your message backups if pressured by authorities or at their own discretion.\n\nFor maximum security, we recommend using our servers where your messages are encrypted with your PIN and only you can decrypt them.")
            }
        }
    }
}

// MARK: - Backup Location Button
struct BackupLocationButton: View {
    let icon: String
    let title: String
    let subtitle: String
    let badge: String?
    let badgeColor: Color
    let isSelected: Bool
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            HStack(spacing: 12) {
                Image(systemName: icon)
                    .font(.system(size: 24))
                    .foregroundColor(isSelected ? Color(hex: "007CFC") : .gray)
                    .frame(width: 36)
                
                VStack(alignment: .leading, spacing: 4) {
                    HStack(spacing: 8) {
                        Text(title)
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundColor(.primary)
                        
                        if let badge = badge {
                            Text(badge)
                                .font(.system(size: 10, weight: .bold))
                                .foregroundColor(.white)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(badgeColor)
                                .cornerRadius(4)
                        }
                    }
                    
                    Text(subtitle)
                        .font(.system(size: 13))
                        .foregroundColor(.gray)
                }
                
                Spacer()
                
                if isSelected {
                    Image(systemName: "checkmark.circle.fill")
                        .foregroundColor(Color(hex: "007CFC"))
                        .font(.system(size: 22))
                }
            }
            .padding(.vertical, 8)
        }
    }
}

// MARK: - Backup PIN Setup View
struct BackupPINSetupView: View {
    @Environment(\.dismiss) var dismiss
    @ObservedObject var viewModel: BackupSettingsViewModel
    @State private var pin = ""
    @State private var confirmPIN = ""
    @State private var passphrase = ""
    @State private var confirmPassphrase = ""
    @State private var usePIN = true
    
    var body: some View {
        NavigationView {
            VStack(spacing: 24) {
                // Icon
                Image(systemName: "lock.shield.fill")
                    .font(.system(size: 60))
                    .foregroundColor(Color(hex: "007CFC"))
                    .padding(.top, 40)
                
                // Title
                VStack(spacing: 8) {
                    Text("Secure Your Backups")
                        .font(.system(size: 24, weight: .bold))
                    
                    Text("Choose a PIN or passphrase to encrypt your message backups")
                        .font(.system(size: 15))
                        .foregroundColor(.gray)
                        .multilineTextAlignment(.center)
                        .padding(.horizontal, 32)
                }
                
                // PIN/Passphrase Toggle
                Picker("", selection: $usePIN) {
                    Text("PIN").tag(true)
                    Text("Passphrase").tag(false)
                }
                .pickerStyle(SegmentedPickerStyle())
                .padding(.horizontal, 32)
                
                // Input Fields
                VStack(spacing: 16) {
                    if usePIN {
                        // PIN Input
                        VStack(alignment: .leading, spacing: 8) {
                            Text("Enter 6-8 digit PIN")
                                .font(.system(size: 13))
                                .foregroundColor(.gray)
                            
                            SecureField("PIN", text: $pin)
                                .textFieldStyle(RoundedBorderTextFieldStyle())
                                .keyboardType(.numberPad)
                            
                            SecureField("Confirm PIN", text: $confirmPIN)
                                .textFieldStyle(RoundedBorderTextFieldStyle())
                                .keyboardType(.numberPad)
                        }
                    } else {
                        // Passphrase Input
                        VStack(alignment: .leading, spacing: 8) {
                            Text("Enter passphrase (min 12 characters)")
                                .font(.system(size: 13))
                                .foregroundColor(.gray)
                            
                            SecureField("Passphrase", text: $passphrase)
                                .textFieldStyle(RoundedBorderTextFieldStyle())
                                .textContentType(.newPassword)
                            
                            SecureField("Confirm Passphrase", text: $confirmPassphrase)
                                .textFieldStyle(RoundedBorderTextFieldStyle())
                                .textContentType(.newPassword)
                        }
                    }
                }
                .padding(.horizontal, 32)
                
                Spacer()
                
                // Create Button
                Button(action: createBackupKey) {
                    Text("Create Backup Key")
                        .font(.system(size: 17, weight: .semibold))
                        .foregroundColor(.white)
                        .frame(maxWidth: .infinity)
                        .frame(height: 50)
                        .background(Color(hex: "007CFC"))
                        .cornerRadius(12)
                }
                .disabled(!isValid)
                .opacity(isValid ? 1 : 0.5)
                .padding(.horizontal, 32)
                .padding(.bottom, 32)
            }
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
            }
        }
    }
    
    private var isValid: Bool {
        if usePIN {
            return pin.count >= 6 && pin.count <= 8 && pin == confirmPIN
        } else {
            return passphrase.count >= 12 && passphrase == confirmPassphrase
        }
    }
    
    private func createBackupKey() {
        viewModel.setupBackupKey(pin: usePIN ? pin : nil, passphrase: usePIN ? nil : passphrase) { success in
            if success {
                dismiss()
            }
        }
    }
}

// MARK: - Backup Now View
struct BackupNowView: View {
    @Environment(\.dismiss) var dismiss
    @ObservedObject var viewModel: BackupSettingsViewModel
    @State private var pin = ""
    @State private var passphrase = ""
    @State private var backupType: BackupType = .full
    @State private var usePIN = true
    
    var body: some View {
        NavigationView {
            VStack(spacing: 24) {
                // Icon
                Image(systemName: "arrow.up.doc.fill")
                    .font(.system(size: 60))
                    .foregroundColor(Color(hex: "007CFC"))
                    .padding(.top, 40)
                
                // Title
                VStack(spacing: 8) {
                    Text("Backup Messages")
                        .font(.system(size: 24, weight: .bold))
                    
                    Text("Enter your PIN/passphrase to encrypt and backup your messages")
                        .font(.system(size: 15))
                        .foregroundColor(.gray)
                        .multilineTextAlignment(.center)
                        .padding(.horizontal, 32)
                }
                
                // Backup Type
                Picker("Backup Type", selection: $backupType) {
                    Text("Full Backup").tag(BackupType.full)
                    Text("Incremental").tag(BackupType.incremental)
                }
                .pickerStyle(SegmentedPickerStyle())
                .padding(.horizontal, 32)
                
                // PIN/Passphrase Input
                VStack(spacing: 16) {
                    Picker("", selection: $usePIN) {
                        Text("PIN").tag(true)
                        Text("Passphrase").tag(false)
                    }
                    .pickerStyle(SegmentedPickerStyle())
                    
                    if usePIN {
                        SecureField("Enter PIN", text: $pin)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                            .keyboardType(.numberPad)
                    } else {
                        SecureField("Enter Passphrase", text: $passphrase)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                            .textContentType(.password)
                    }
                }
                .padding(.horizontal, 32)
                
                Spacer()
                
                // Backup Button
                Button(action: createBackup) {
                    if viewModel.isBackingUp {
                        ProgressView()
                            .progressViewStyle(CircularProgressViewStyle(tint: .white))
                            .frame(maxWidth: .infinity)
                            .frame(height: 50)
                    } else {
                        Text("Start Backup")
                            .font(.system(size: 17, weight: .semibold))
                            .foregroundColor(.white)
                            .frame(maxWidth: .infinity)
                            .frame(height: 50)
                    }
                }
                .background(Color(hex: "007CFC"))
                .cornerRadius(12)
                .disabled(!isValid || viewModel.isBackingUp)
                .opacity(isValid && !viewModel.isBackingUp ? 1 : 0.5)
                .padding(.horizontal, 32)
                .padding(.bottom, 32)
            }
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
            }
        }
    }
    
    private var isValid: Bool {
        if usePIN {
            return pin.count >= 6 && pin.count <= 8
        } else {
            return passphrase.count >= 12
        }
    }
    
    private func createBackup() {
        viewModel.createBackup(
            pin: usePIN ? pin : nil,
            passphrase: usePIN ? nil : passphrase,
            type: backupType
        ) { success in
            if success {
                dismiss()
            }
        }
    }
}

// MARK: - Backup History View
struct BackupHistoryView: View {
    @ObservedObject var viewModel: BackupSettingsViewModel
    
    var body: some View {
        List {
            ForEach(viewModel.backupHistory) { backup in
                VStack(alignment: .leading, spacing: 8) {
                    HStack {
                        Text(backup.type == .full ? "Full Backup" : "Incremental Backup")
                            .font(.system(size: 16, weight: .semibold))
                        
                        Spacer()
                        
                        Text(backup.date, style: .date)
                            .font(.system(size: 14))
                            .foregroundColor(.gray)
                    }
                    
                    HStack {
                        Label("\(backup.messagesCount) messages", systemImage: "message.fill")
                            .font(.system(size: 13))
                            .foregroundColor(.gray)
                        
                        Spacer()
                        
                        Text(ByteCountFormatter.string(fromByteCount: backup.size, countStyle: .file))
                            .font(.system(size: 13))
                            .foregroundColor(.gray)
                    }
                }
                .padding(.vertical, 4)
                .swipeActions {
                    Button(role: .destructive) {
                        viewModel.deleteBackup(id: backup.id)
                    } label: {
                        Label("Delete", systemImage: "trash")
                    }
                }
            }
        }
        .navigationTitle("Backup History")
        .navigationBarTitleDisplayMode(.inline)
    }
}

// MARK: - Models
enum BackupLocation {
    case ourServers
    case iCloud
}

enum BackupFrequency: String {
    case daily = "daily"
    case weekly = "weekly"
    case monthly = "monthly"
}

enum BackupType {
    case full
    case incremental
}

struct BackupHistoryItem: Identifiable {
    let id: String
    let date: Date
    let type: BackupType
    let messagesCount: Int
    let size: Int64
}

// MARK: - View Model
class BackupSettingsViewModel: ObservableObject {
    @Published var backupEnabled = true
    @Published var backupLocation: BackupLocation = .ourServers
    @Published var autoBackupEnabled = true
    @Published var backupFrequency: BackupFrequency = .daily
    @Published var wifiOnly = true
    @Published var hasBackupKey = false
    @Published var lastBackupDate: Date?
    @Published var isBackingUp = false
    @Published var backupHistory: [BackupHistoryItem] = []
    
    func toggleBackup() {
        // TODO: API call
    }
    
    func selectBackupLocation(_ location: BackupLocation) {
        backupLocation = location
        // TODO: API call
    }
    
    func acknowledgeThirdPartyWarning() {
        // TODO: API call
    }
    
    func setupBackupKey(pin: String?, passphrase: String?, completion: @escaping (Bool) -> Void) {
        // TODO: API call
        hasBackupKey = true
        completion(true)
    }
    
    func createBackup(pin: String?, passphrase: String?, type: BackupType, completion: @escaping (Bool) -> Void) {
        isBackingUp = true
        
        // TODO: API call
        DispatchQueue.main.asyncAfter(deadline: .now() + 2) {
            self.isBackingUp = false
            self.lastBackupDate = Date()
            completion(true)
        }
    }
    
    func saveSettings() {
        // TODO: API call
    }
    
    func deleteBackup(id: String) {
        // TODO: API call
        backupHistory.removeAll { $0.id == id }
    }
    
    func deleteAllBackups() {
        // TODO: API call
        backupHistory.removeAll()
    }
}

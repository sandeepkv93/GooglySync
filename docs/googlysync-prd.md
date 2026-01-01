# PRD: googlysync Desktop Client for Pop!_OS 24

## Overview
Build a native-like googlysync desktop client for Pop!_OS 24 that delivers parity with core features from Google Drive clients on Windows and macOS. The app will provide secure Google account sign-in, background sync, offline access, system integration, and a reliable file system experience for end users who need Drive integrated into their Linux desktop workflow.

## Goals
- Provide feature parity with Windows/macOS Google Drive clients for everyday productivity.
- Offer a stable, secure, and performant sync experience on Pop!_OS 24.
- Support Google account sign-in and multi-account switching.
- Integrate with the GNOME-based desktop environment (tray, notifications, file manager).

## Non-Goals
- Full reimplementation of Google Drive web UI or advanced admin tooling.
- Enterprise-only features beyond standard consumer Drive.
- Building a complete Google Workspace suite (Docs/Sheets/Slides editors).

## Target Users
- Individual users who want Drive files locally mirrored on Pop!_OS.
- Families or small groups sharing files in consumer Google accounts.
- Developers and creators who need offline access and fast personal file sync.

## Success Metrics
- Sync reliability: >99% successful file operations across test suites.
- Sign-in completion rate: >95% without manual support.
- Median time-to-first-sync: <5 minutes for 1 GB with average network.
- Crash-free sessions: >99%.

## Key Requirements
- Desktop client that runs on Pop!_OS 24 and integrates with GNOME.
- Google account login via OAuth 2.0 with modern security standards.
- File sync with Drive parity: My Drive, shared files.
- Local caching, offline access, and background sync.
- Support both mirror mode and streaming mode.

## Feature Parity Scope

### 1) Account and Authentication
- Google account sign-in with OAuth 2.0 (PKCE) and system browser flow.
- Multi-account support: add, remove, and switch accounts.
- Per-account sync settings and quotas.
- Session management: token refresh, logout, and token revocation.
- 2FA and security key compatibility (via browser-based auth).

### 2) File Sync and Local Drive Mount
- Mirror My Drive to a local folder (two-way sync).
- Optional streaming mode: files appear locally but download on open.
- Selective sync: choose folders to sync offline.
- Shared-with-me handling (shortcut or mounted subfolder behavior).
- Conflict resolution: clear UI for conflicts with rename/version strategy.
- Pause/resume sync.
- Bandwidth throttling (upload/download limits, schedules).
- Respect Google Drive file shortcuts and document types.
- Handling for Google Docs/Sheets/Slides (export format or placeholder).

### 3) Filesystem Integration
- Local folder mount with file watcher integration.
- Drive status icons (synced, syncing, error) in file manager if feasible.
- Context menu actions (view on Drive, share, manage versions).
- Desktop notifications for errors, conflicts, and major events.
- Autostart on login.

### 4) Upload and Download Behaviors
- Large file upload with resumable uploads.
- Chunked downloads with resume support.
- Sync prioritization (recent edits first).
- Respect Drive file size limits and quotas.

### 5) Collaboration and Sharing
- Share files/folders from desktop (opens web share dialog).
- View file activity and version history (web deep links).
- Display permissions (viewer/commenter/editor) in UI.

### 6) Offline Access
- Mark files/folders available offline.
- Offline edits queued and synced when online.
- Offline indicators in UI.

### 7) Search and Quick Access
- Local search for synced files.
- Quick access list (recently used files).
- Open Drive web search for cloud-only items.

### 8) Settings and Preferences
- Sync mode settings (mirror vs stream) with ability to switch.
- Folder location (choose local Drive folder path).
- Notifications on/off and severity controls.
- Network settings: proxy configuration, bandwidth limits, offline mode.
- Storage usage overview (cache size, clear cache).
- Update channel preference (stable/beta).

### 9) Updates and Diagnostics
- Auto-update mechanism with signed packages.
- Diagnostics page: sync status, error logs, last sync time.
- Export logs for support.
- Health checks (token validity, API availability).

### 10) Security and Privacy
- OAuth 2.0 with PKCE; no embedded password entry.
- Encrypted token storage (GNOME Keyring/libsecret).
- Secure local cache permissions (per-user).
- File integrity checks (hash validation).

## User Experience
- Minimal tray icon with status and quick actions.
- Onboarding wizard: sign in, choose sync mode, select folders.
- Clear sync status dashboard.
- Actionable error messages with retry options.

## User Flows
1) Install and Launch
- User installs package, launches app, sees onboarding.

2) Sign In
- User clicks "Sign in", system browser opens Google OAuth.
- Consent and account selection; app receives tokens and finishes setup.

3) First Sync
- User selects sync mode and folders.
- App starts initial sync, status visible in tray and dashboard.

4) Daily Use
- Files appear in local folder; edits sync automatically.

5) Conflict
- App detects conflict, displays resolution UI or creates conflict copies.

## Technical Requirements
- Google Drive API v3 integration with proper quota handling.
- Background sync daemon + UI shell.
- File watcher for local changes (inotify-based).
- Robust local metadata store (SQLite) for state tracking.
- Retry/backoff for network and API errors.
- Works offline and resumes correctly on reconnect.

## Dependencies
- Google OAuth 2.0 configuration (client ID, redirect URIs).
- Google Drive API access and scope approvals.
- GNOME Keyring for secure token storage.

## Risks and Mitigations
- API rate limits: implement exponential backoff and batching.
- File conflicts: deterministic conflict handling with user visibility.
- Network instability: resume support and local queueing.
- OAuth changes: keep auth flow aligned with Google policies.

## Open Questions
- Should the client support full-text search across cloud-only files?
- What export format should be default for Google Docs/Sheets/Slides?
- Should shared-with-me be visible as shortcuts or a separate mount?

## Milestones (High-Level)
1) OAuth + basic sync for My Drive.
2) Selective sync + streaming mode.
3) Sharing actions and shared-with-me handling.
4) Offline access and conflict UI.
5) Polished UX, diagnostics, and updates.

## Acceptance Criteria
- User can sign in with a Google account and complete onboarding.
- Files appear in a local folder and two-way sync works reliably.
- Shared-with-me items are visible and sync correctly.
- Offline marking works and syncs changes when online.
- Errors are visible with clear recovery actions.

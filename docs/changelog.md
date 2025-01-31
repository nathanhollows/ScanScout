---
title: "Changelog"
sidebar: true
order: 100
---

# Changelog

## 3.2.0 (2025-01-31)

### Added
- Players are now redirected to a dedicated end screen upon game completion. Included an easter egg when tapping the confetti icon. [#53](https://github.com/nathanhollows/Rapua/issues/53)

### Changed
- Restyled the lobby and team name form for better clarity and user experience.
- Improved input validation and feedback messages, especially for check-in/out forms.
- Footer now includes team name, rules link, and team code for quick reference and is shown on more pages.
- Differentiated check-in elements for logged-in vs. logged-out players for a smoother experience.
- Style updates across the platform for better consistency and readability.
- Misc code refactors on internal services and handlers for readability and consistency.

### Fixed
- Fixed an issue where commas in filenames were preventing asset downloads. [#57](https://github.com/nathanhollows/Rapua/issues/57)
- Fixed an issue where the team overview failed to render when a player had visited all locations.
- Resolved issues caused by blank sessions and edge cases that led to unexpected behaviour.
- Fixed a bug where check-in/out pages weren't rendering if the player didn't have a session.

[Full Changelog](https://github.com/nathanhollows/Rapua/releases/tag/v3.2.0)

---

## 3.1.0 (2025-01-17)

### Added
- Instantly fit map bounds to neighbouring markers when adding locations [(#44)](https://github.com/nathanhollows/Rapua/issues/44)

- Introduced the official project logo.
- Docs for [Getting Started with Teams](/docs/user/teams) [(#43)](https://github.com/nathanhollows/Rapua/issues/43)

### Changed
- Improved team activity overview for easier browsing and better visual clarity.
- Updated a documentation icon for consistency across the interface.

[Full Changelog](https://github.com/nathanhollows/Rapua/releases/tag/v3.1.0)

---

## 3.0.1 (2025-01-14)

### Fixed
- Fixed an issue where the user could not switch instances. `current_instance_id` was blocked from updating in the database.

---

## 3.0.0 (2025-01-09)

### Added
- **Contact Form:**
  - Contact form messages are now sent to the platform admin instead of being thrown into the void. [#23](https://github.com/nathanhollows/Rapua/issues/23)
- **User Management:**
  - New method for deleting users programmatically, including cascading logic to clean up dependent data.
- **Documentation Enhancements:**
  - Added roadmap/wishlist for brainstorming and tracking future features.
  - Developer-specific migration documentation, including instructions, testing steps, and explanations.
  - Included a [history](/docs/user/history) of Rapua as a reference for users and developers. [#29](https://github.com/nathanhollows/Rapua/issues/29)
  - Minor changes to in-app hints and tips for better user experience.
  - Tests to ensure documentation is up-to-date and accurate; links are now checked for validity; pages must not be empty.
- **Database Migrations:**
  - Implemented a new migration system for database changes.

### Changed
- **Breaking Changes:**
  - Renamed `cmd/game-server` to `cmd/rapua` for consistency.
  - Major refactor of services and repository plumbing for better separation of concerns, maintainability, and scalability.
  - Repositories now accept a `*bun.Tx` for bulk deletions and Services now require a `db.Transactor` for beginning transactions.
  - All tests now use migrations to ensure a clean database state and non-global database vars.
  - Teams must now have UUID-based IDs for uniformity and scalability.
  - Models no longer support `deleted_at`; hard deletes are now the standard. Soft-deleted data was never used.
- **Style Updates:**
  - Docs now look better on mobile devices.
  - Submenu titles in documentation shrunk for consistency.
  - [Content Blocks](/docs/user/blocks/) are now collapsible for better readability, and auto-collapse if there are more than 3 blocks.
  - Team CheckIns are now collapsible for better readability.

### Fixed
- The confirmation modal for deleting content blocks now triggers correctly. [#36](https://github.com/nathanhollows/Rapua/issues/36)
- Resolved issues with block rendering and automatic updates. [#34](https://github.com/nathanhollows/Rapua/issues/34)
- Map markers now show the same numbers as the location list for consistency.
- Ensured user registration failures clean up partially created user data. [#40](https://github.com/nathanhollows/Rapua/issues/40)
- Registration emails are now automatically sent. [#40](https://github.com/nathanhollows/Rapua/issues/40)

### Removed
- Unused methods and global variables related to the database have been removed for better maintainability.
- Unused methods from Blocks interface removed.
- Deprecated methods for database initialization and testing replaced with updated patterns.

[Full Changelog](https://github.com/nathanhollows/Rapua/releases/tag/v3.0.0)

---

## 2.5.0 (2024-11-28)

### Added
- **Team Features:**
  - Teams can now set their own names directly within the platform.
  - Bulk deletion of teams is now supported for easier management.
- **Interactive Map:**
  - Map markers now display popups with names for better navigation and understanding of locations. [Closes #25](https://github.com/nathanhollows/Rapua/issues/25)
- **Documentation:**
  - New public documentation system [(Closes #19)](https://github.com/nathanhollows/Rapua/issues/19), including documentation:
    - User and developer guides.
    - Quickstart guide.
    - Tutorials, such as a Student Induction Tutorial.
  - "Docs" now takes the place of "Inspo" in the main navigation. [Fixes #35](https://github.com/nathanhollows/Rapua/issues/35)
- **API Updates:**
  - Introduced new endpoints for managing teams in the Teams service, allowing programmatic creation, updates, and bulk deletion.
  - Added support for creating locations using existing map markers.

### Changed
- The team management interface has been redesigned for usability, using Hyperscript for interactivity.
- Game logic refactoring ensures more efficient handling of location relationships, including clues.

### Fixed
- Resolved an issue where marker names were not populating in the activity overview.
- Corrected a database query issue with an incorrect column name.
- Fixed an issue where game relationships weren’t fully loaded, affecting progression in specific scenarios.

[Full Changelog](https://github.com/nathanhollows/Rapua/releases/tag/v2.5.0)

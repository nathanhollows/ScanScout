---
title: "Changelog"
sidebar: true
order: 100
---

# Changelog

## 2.5.0 (2024-11-28)

#### Added
- **Team Features:**
  - Teams can now set their own names directly within the platform.
  - Bulk deletion of teams is now supported for easier management.
- **Interactive Map:**
  - Map markers now display popups with names for better navigation and understanding of locations.
- **Documentation:**
  - New public documentation system, including documentation for:
    - User and developer guides.
    - Quickstart guide.
    - Tutorials, such as a Student Induction Tutorial.
  - "Docs" now takes the place of "Inspo" in the main navigation.
- **API Updates:**
  - Introduced new endpoints for managing teams in the Teams service, allowing programmatic creation, updates, and bulk deletion.
  - Added support for creating locations using existing map markers.

#### Changed
- The team management interface has been redesigned for usability, using Hyperscript for interactivity.
- Game logic refactoring ensures more efficient handling of location relationships, including clues.

#### Fixed
- Resolved an issue where marker names were not populating in the activity overview.
- Corrected a database query issue with an incorrect column name.
- Fixed an issue where game relationships weren’t fully loaded, affecting progression in specific scenarios.

[Full Changelog](https://github.com/nathanhollows/Rapua/releases/tag/v2.5.0)

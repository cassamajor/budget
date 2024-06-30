# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.1] - 2024-06-30

### Changed

- A code snippet in the tutorial was updated to include the correct amount of return values.  

## [0.2.0] - 2024-06-30

### Added

- README.md now contains a brief tutorial describing features and functionality.
- CHANGELOG.md to keep track of changes to the project.
- Support for all YNAB account types, stored in the `Assets` and `Liabilities` structs.

### Changed

- `NetWorth` is now comprised of `Assets` and `Liabilities`.
- `AddComma` is renamed to `FormatCurrency`.

## [0.1.1] - 2024-06-29

### Fixed

- Funded Spending from the `Credit Card Payment` category skews monthly Expenses. It is now removed from the total.
- Calculation of Net Worth had the incorrect arithmetic operator.

## [0.1.0] - 2024-06-20

### Added

- Ability to query the YNAB API for budget data.

[unreleased]: https://github.com/cassamajor/budget/compare/v0.2.1...HEAD
[0.2.1]: https://github.com/cassamajor/budget/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/cassamajor/budget/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/cassamajor/budget/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/cassamajor/budget/releases/tag/v0.1.0
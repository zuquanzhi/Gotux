# Changelog

All notable changes to this project will be documented in this file.

## [2.1.0] - 2025-01-02

### Changed
- Redesigned Dashboard with minimalist aesthetic
- Simplified Upload page layout
- Removed redundant UI elements (welcome banners, quick actions, tips)
- Switched to clean card-based design with subtle hover effects

### UI Improvements
- Dashboard: 4-column responsive grid for statistics
- Upload: Single-column layout with inline statistics
- Removed gradients and excessive colors
- Added subtle border highlights on hover
- Improved responsive behavior on mobile devices

## [2.0.0] - 2025-10-02

### Added
- Personal settings center with 12+ configuration options
- Custom domain support for image links
- Default link format configuration (URL/Markdown/HTML/BBCode)
- Image watermark with customizable text and position
- Automatic image compression with quality control (1-100)
- Single file size limit configuration (0.1-50 MB)
- Allowed image format selection
- Storage quota display with visual progress
- Image review system for multi-user platforms
- Database migration tool for settings

### Backend
- Added 12 new fields to User model
- Created UpdateSettings and GetSettings endpoints
- Added validation for compression quality and file size
- Implemented settings API routes

### Frontend
- Created comprehensive Profile page (6:18 layout)
- Added sticky sidebar for settings navigation
- Implemented form validation
- Added real-time storage quota visualization

### Documentation
- Added SETTINGS_GUIDE.md
- Added SETTINGS_QUICKREF.md
- Updated README with settings overview

## [1.0.0] - Initial Release

### Features
- User registration and authentication
- JWT-based authentication
- Image upload with drag-and-drop
- Batch image upload
- Multiple link export formats
- Image management (view, delete)
- User profile management
- Admin dashboard
- User management for admins
- Storage statistics
- MD5-based image deduplication
- Responsive UI design

### Tech Stack
- Backend: Go 1.21, Gin, GORM, SQLite
- Frontend: Vue 3, Vite, Element Plus, Pinia
- Authentication: JWT
- Database: SQLite

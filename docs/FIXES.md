# Bug Fixes and Improvements Summary

## Date: 2025-10-02 (Update 3)

### Bug Fixes

#### 6. CORS 403 Error for IPv6 Requests
**Issue**: When accessing `/api/random` from browsers using IPv6 (`::1`), requests returned 403 Forbidden errors. However, IPv4 (`127.0.0.1`) requests worked fine.

**Root Cause**: The CORS configuration only explicitly allowed specific origins like `http://localhost:5173`, but browsers using IPv6 send Origin header as `http://[::1]:5173`, which wasn't in the allowed list.

**Attempts**:
1. Added IPv6 localhost variations to `AllowOrigins`
2. Used `AllowOriginFunc` to dynamically check origins

**Final Solution**: For development mode, simplified to `AllowAllOrigins: true`:

```go
r.Use(cors.New(cors.Config{
    AllowAllOrigins:  true,  // Development mode
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length", "X-Image-UUID", "X-Image-ID"},
    AllowCredentials: false, // Must be false when AllowAllOrigins is true
}))
```

**Note**: For production deployment, use specific origins in `AllowOrigins` array.

Location: `backend/main.go`

## Date: 2025-10-02 (Update 2)

### Bug Fixes

#### 4. Admin Images Page Preview Overlap
**Issue**: In the admin images management page, when viewing large images, UI elements were overlapping the preview modal (same issue as dashboard).

**Fix**: Added `z-index="9999"` and `preview-teleported="true"` props to `el-image` component.

Locations:
- `frontend/src/views/admin/Images.vue`
- `frontend/src/views/Upload.vue`
- `frontend/src/views/Images.vue`

#### 5. Random API Authentication Issue
**Issue**: When calling `/api/random` from frontend Vue app, the request utility automatically added Authorization header even when empty, which could cause 403 errors.

**Fix**: Created dedicated random image API functions in `frontend/src/api/image.js` that use native `fetch` to avoid adding Authorization header.

```javascript
export function getRandomImage(params) {
  const queryString = new URLSearchParams(params).toString()
  const url = `/api/random${queryString ? '?' + queryString : ''}`
  return fetch(url).then(res => res.json())
}
```

## Date: 2025-01-02

### Bug Fixes

#### 1. Dashboard Image Preview UI Overlap
**Issue**: When viewing large images in the dashboard, UI elements (navigation, headers) were overlapping the preview modal.

**Fix**: Added `z-index="9999"` and `preview-teleported="true"` props to `el-image` component in Dashboard.vue to ensure the image viewer appears above all other elements.

Location: `frontend/src/views/Dashboard.vue`

#### 2. Custom Domain Not Applied to Image Links
**Issue**: When users configured a custom domain in settings, the generated image links still used the default domain.

**Fix**: Updated `GetImageLinks` function to:
- Retrieve user settings
- Use `CustomDomain` if configured
- Fall back to request host if not configured

Location: `backend/controllers/image.go`

### New Features

#### 3. UUID-based Secure Image Access
**Problem**: Direct path-based URLs (`/uploads/2025/01/02/image.jpg`) expose:
- Directory structure
- Allow enumeration attacks
- Make it easy to guess other image URLs

**Solution**: Implemented UUID-based image access system.

**Changes**:
1. **Database Schema**:
   - Added `UUID` field to `Image` model (unique indexed)
   - Auto-generates UUID on image creation
   - Migration tool to add UUIDs to existing images

2. **API Endpoints**:
   - `GET /api/i/:uuid` - Get image info by UUID
   - `GET /i/:uuid` - Serve image file directly by UUID
   - Both endpoints track view counts

3. **Link Generation**:
   - Primary link uses UUID format: `/i/{uuid}`
   - Legacy direct path still available as `direct_url`
   - All link formats (Markdown, HTML, BBCode) use UUID

**Benefits**:
- **Security**: Prevents enumeration and guessing
- **Privacy**: Hides file structure
- **Access Control**: Can implement per-image permissions
- **Analytics**: Automatic view count tracking
- **Flexibility**: Can change storage backend without breaking links

**Migration**:
```bash
cd backend/cmd/migrate_uuid
go run main.go
```

Location: 
- Model: `backend/models/image.go`
- Controller: `backend/controllers/image.go`
- Routes: `backend/routes/routes.go`
- Migration: `backend/cmd/migrate_uuid/main.go`

### Documentation Reorganization

**Problem**: Documentation was scattered across multiple files in root directory with:
- Duplicate content
- Excessive emoji usage
- Unclear organization
- Too many redundant files

**Solution**: Reorganized all documentation into `docs/` folder with clean, professional structure.

**New Structure**:
```
docs/
├── README.md         - Documentation index
├── DEPLOYMENT.md     - Production deployment guide
├── API.md           - Complete API reference
└── CHANGELOG.md     - Version history
```

**Removed Files**:
- UI_MINIMALISM.md
- UI_IMPROVEMENTS.md
- LAYOUT_IMPROVEMENTS.md
- SETTINGS_GUIDE.md
- SETTINGS_QUICKREF.md
- DOMAIN_BINDING_GUIDE.md
- DEPLOY.md (merged into DEPLOYMENT.md)

**Style Changes**:
- Removed all emoji
- Adopted GitHub high-star project style
- Consolidated duplicate content
- Clear, professional English
- Better organization and navigation

### Testing Checklist

- [ ] Dashboard image preview works without UI overlap
- [ ] Custom domain setting applies to generated links
- [ ] UUID migration runs successfully
- [ ] New images get UUID automatically
- [ ] `/i/:uuid` endpoint serves images correctly
- [ ] View counts increment properly
- [ ] Both UUID and direct path links work
- [ ] API documentation is accurate
- [ ] All documentation links are correct

### Breaking Changes

None. All changes are backwards compatible:
- Old direct path URLs still work
- Existing images work after UUID migration
- No API endpoints were removed

### Upgrade Instructions

1. Update code from repository
2. Run UUID migration:
   ```bash
   cd backend/cmd/migrate_uuid
   go run main.go
   ```
3. Restart backend service
4. Clear browser cache and reload frontend

### Future Improvements

Consider:
1. Add image expiration/TTL feature
2. Implement rate limiting for `/i/:uuid` endpoint
3. Add image analytics dashboard
4. Support private/unlisted images
5. Add hotlink protection
6. Implement CDN integration for UUID links

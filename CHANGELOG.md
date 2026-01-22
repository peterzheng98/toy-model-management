# Changelog

## Version 2.0 - User Tracking & Statistics

### New Features Added

#### Usage Tracking System
- **Request Logging**: All API requests are now logged with username, IP address, timestamp, and action
- **Access Log Database**: Persistent storage of all system activity in `access_log.json`
- **Automatic User Detection**: Go client automatically detects system username using multiple methods

#### Statistics Dashboard
New statistics displayed on both user and admin interfaces:
- **Total Models**: Number of models in the system
- **Total Folder Size**: Combined size of all models in the mount point
- **Total Requests**: Count of all API requests made
- **Unique Users**: Number of different users who have accessed the system
- **Download Counts**: Per-model download statistics
- **Recent Activity**: Last 10 actions performed on the system

#### User Attribution
- **First Downloader Tracking**: Records who requested each model download first
- **IP Address Logging**: Tracks the IP address of the first downloader
- **Username Input**: Admin panel now includes username field (saved in browser)
- **Display in UI**: User and admin panels show who downloaded each model first

#### Enhanced Model Information
Models now include:
- `downloaded_by`: Username of the person who downloaded the model
- `stats`: Object containing usage statistics
  - `download_count`: Number of times the model was downloaded/requested
  - `access_count`: Number of times the model was accessed (GET requests)
  - `total_requests`: Total API requests for this model
  - `first_downloaded_by`: Username of first downloader
  - `first_downloaded_at`: Timestamp of first download
  - `first_downloaded_from`: IP address of first downloader

### API Changes

#### New Endpoint: `/api/stats`
Returns overall system statistics including:
- Total models count
- Total storage size in bytes
- Total requests count
- Unique users count
- Recent activity log

#### Modified Endpoints

**POST `/api/models/download`**
- Now requires `username` field in request body
- Returns model with statistics

**GET `/api/models`**
- Now includes `stats` object for each model
- Shows download counts and user attribution

**GET `/api/models/<id>`**
- Now includes `stats` object with detailed usage information

### Client Changes

#### Go Client Updates
- **Automatic Username Detection**: Uses `whoami`, environment variables, or Go's user package
- **Enhanced Output**: Shows download counts and first downloader in list view
- **Statistics Display**: `get` command now shows full usage statistics
- **Username Sending**: Automatically includes username in download requests

### Web UI Changes

#### User View (index.html)
- Added statistics dashboard at the top
- Shows system-wide metrics (models, size, requests, users)
- Model cards now display:
  - Download count badge
  - Total request count
  - First downloader information with IP address

#### Admin Panel (admin.html)
- Added statistics dashboard at the top
- Username input field (auto-saved to browser localStorage)
- Enhanced model table with:
  - Download count column
  - First downloader column
- Real-time statistics updates every 10 seconds

### Backend Changes

#### New Functions in `app.py`
- `load_access_log()`: Load request history from file
- `save_access_log()`: Save request history to file
- `log_request()`: Log individual requests with user info and IP
- `get_model_stats()`: Calculate usage statistics for a specific model
- `get_total_models_size()`: Calculate total size of all models

#### New Files Created
- `access_log.json`: Stores all request logs (created automatically)

### Example Updates

#### Python Example
- Now uses `getpass.getuser()` to detect username
- Enhanced output with statistics display
- New `get_stats()` function to retrieve system statistics

#### Go Example  
- Automatic username detection
- Enhanced output with statistics
- New `getStats()` function

### Documentation Updates

- Updated README.md with new features and API documentation
- Updated QUICKSTART.md with usage tracking information
- Added this CHANGELOG.md

### Backward Compatibility

**Breaking Changes**:
- POST `/api/models/download` now requires `username` field
- All API responses now include additional `stats` fields

**Migration Notes**:
- Existing models will work fine
- Access logs start fresh (no historical data for old models)
- Old clients will need to be updated to send username

### Data Storage

New files in the mount point:
```
models/
├── models_db.json          # Model metadata (existing)
├── access_log.json         # Request logs (NEW)
└── <model_directories>/    # Model files
```

### Security Considerations

- IP addresses are logged (ensure compliance with privacy regulations)
- Username is user-provided (consider authentication for production)
- Access logs can grow large over time (consider log rotation)

---

## Version 1.0 - Initial Release

### Features
- Flask server with RESTful API
- Web UI for users and administrators
- Hugging Face model downloads
- CRUD operations for models
- Go CLI client
- Model existence checking
- Configurable mount point

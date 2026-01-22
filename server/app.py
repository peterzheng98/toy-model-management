import os
import json
import socket
import subprocess
from pathlib import Path
from flask import Flask, request, jsonify, render_template
from flask_cors import CORS
from huggingface_hub import snapshot_download, hf_hub_download
from datetime import datetime
import shutil

app = Flask(__name__)
CORS(app)

# Get mount point from environment variable or use default
MODELS_MOUNT_POINT = os.environ.get('MODELS_MOUNT_POINT', './models')
MODELS_DB_FILE = os.path.join(MODELS_MOUNT_POINT, 'models_db.json')
ACCESS_LOG_FILE = os.path.join(MODELS_MOUNT_POINT, 'access_log.json')

# Ensure mount point exists
os.makedirs(MODELS_MOUNT_POINT, exist_ok=True)


def load_models_db():
    """Load models database from JSON file"""
    if os.path.exists(MODELS_DB_FILE):
        with open(MODELS_DB_FILE, 'r') as f:
            return json.load(f)
    return {}


def save_models_db(db):
    """Save models database to JSON file"""
    with open(MODELS_DB_FILE, 'w') as f:
        json.dump(db, indent=2, fp=f)


def check_model_exists(model_name):
    """Check if model exists in the mount point"""
    model_path = os.path.join(MODELS_MOUNT_POINT, model_name.replace('/', '_'))
    return os.path.exists(model_path) and os.path.isdir(model_path)


def get_model_size(model_path):
    """Calculate total size of model directory in bytes"""
    total_size = 0
    for dirpath, dirnames, filenames in os.walk(model_path):
        for f in filenames:
            fp = os.path.join(dirpath, f)
            if os.path.exists(fp):
                total_size += os.path.getsize(fp)
    return total_size


def get_total_models_size():
    """Calculate total size of all models in bytes"""
    return get_model_size(MODELS_MOUNT_POINT)


def load_access_log():
    """Load access log from JSON file"""
    if os.path.exists(ACCESS_LOG_FILE):
        with open(ACCESS_LOG_FILE, 'r') as f:
            return json.load(f)
    return []


def save_access_log(log):
    """Save access log to JSON file"""
    with open(ACCESS_LOG_FILE, 'w') as f:
        json.dump(log, indent=2, fp=f)


def log_request(action, model_id=None, username=None):
    """Log a request with user info and IP"""
    log = load_access_log()
    
    # Get client IP
    if request.headers.get('X-Forwarded-For'):
        ip_address = request.headers.get('X-Forwarded-For').split(',')[0]
    else:
        ip_address = request.remote_addr
    
    # Get username from request or use provided username
    if username is None:
        data = request.get_json() if request.is_json else {}
        username = data.get('username', 'anonymous')
    
    log_entry = {
        'timestamp': datetime.utcnow().isoformat(),
        'action': action,
        'model_id': model_id,
        'username': username,
        'ip_address': ip_address,
        'user_agent': request.headers.get('User-Agent', 'Unknown')
    }
    
    log.append(log_entry)
    save_access_log(log)
    
    return log_entry


def get_model_stats(model_id):
    """Get usage statistics for a model"""
    log = load_access_log()
    
    # Filter logs for this model
    model_logs = [entry for entry in log if entry.get('model_id') == model_id]
    
    # Count downloads
    download_count = len([e for e in model_logs if e['action'] == 'download'])
    
    # Count accesses (get requests)
    access_count = len([e for e in model_logs if e['action'] == 'get'])
    
    # Get first downloader
    first_download = next((e for e in model_logs if e['action'] == 'download'), None)
    
    return {
        'download_count': download_count,
        'access_count': access_count,
        'total_requests': len(model_logs),
        'first_downloaded_by': first_download.get('username') if first_download else None,
        'first_downloaded_at': first_download.get('timestamp') if first_download else None,
        'first_downloaded_from': first_download.get('ip_address') if first_download else None
    }


# Frontend routes
@app.route('/')
def index():
    """Serve the main page"""
    return render_template('index.html')


@app.route('/admin')
def admin():
    """Serve the admin page"""
    return render_template('admin.html')


# API routes
@app.route('/api/models', methods=['GET'])
def list_models():
    """List all models"""
    log_request('list')
    db = load_models_db()
    
    # Add statistics to each model
    models_with_stats = []
    for model in db.values():
        model_data = model.copy()
        model_data['stats'] = get_model_stats(model['id'])
        models_with_stats.append(model_data)
    
    return jsonify({
        'success': True,
        'models': models_with_stats
    })


@app.route('/api/models/<path:model_id>', methods=['GET'])
def get_model(model_id):
    """Get a specific model by ID"""
    log_request('get', model_id)
    db = load_models_db()
    
    if model_id not in db:
        return jsonify({
            'success': False,
            'error': 'Model not found'
        }), 404
    
    model_data = db[model_id].copy()
    model_data['stats'] = get_model_stats(model_id)
    
    return jsonify({
        'success': True,
        'model': model_data
    })


@app.route('/api/models/download', methods=['POST'])
def download_model():
    """Download a model from Hugging Face"""
    data = request.get_json()
    
    if not data or 'model_name' not in data:
        return jsonify({
            'success': False,
            'error': 'model_name is required'
        }), 400
    
    model_name = data['model_name']
    username = data.get('username', 'anonymous')
    model_id = model_name.replace('/', '_')
    
    # Log the download request
    log_request('download', model_id, username)
    
    # Check if model already exists
    if check_model_exists(model_name):
        db = load_models_db()
        if model_id in db:
            model_data = db[model_id].copy()
            model_data['stats'] = get_model_stats(model_id)
            return jsonify({
                'success': True,
                'message': 'Model already exists',
                'model': model_data,
                'already_exists': True
            })
    
    try:
        # Download model from Hugging Face
        model_path = os.path.join(MODELS_MOUNT_POINT, model_id)
        
        print(f"Downloading model {model_name} to {model_path}...")
        cache_dir = snapshot_download(
            repo_id=model_name,
            cache_dir=model_path,
            resume_download=True
        )
        
        # Calculate model size
        size_bytes = get_model_size(model_path)
        
        # Update database
        db = load_models_db()
        db[model_id] = {
            'id': model_id,
            'name': model_name,
            'path': model_path,
            'size_bytes': size_bytes,
            'downloaded_at': datetime.utcnow().isoformat(),
            'downloaded_by': username,
            'status': 'ready'
        }
        save_models_db(db)
        
        model_data = db[model_id].copy()
        model_data['stats'] = get_model_stats(model_id)
        
        return jsonify({
            'success': True,
            'message': 'Model downloaded successfully',
            'model': model_data
        })
        
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/models/<path:model_id>', methods=['DELETE'])
def delete_model(model_id):
    """Delete a model"""
    log_request('delete', model_id)
    db = load_models_db()
    
    if model_id not in db:
        return jsonify({
            'success': False,
            'error': 'Model not found'
        }), 404
    
    try:
        # Delete model directory
        model_path = db[model_id]['path']
        if os.path.exists(model_path):
            shutil.rmtree(model_path)
        
        # Remove from database
        del db[model_id]
        save_models_db(db)
        
        return jsonify({
            'success': True,
            'message': 'Model deleted successfully'
        })
        
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/models/<path:model_id>', methods=['PUT'])
def update_model(model_id):
    """Update model metadata"""
    db = load_models_db()
    
    if model_id not in db:
        return jsonify({
            'success': False,
            'error': 'Model not found'
        }), 404
    
    data = request.get_json()
    
    # Allow updating certain fields
    allowed_fields = ['status']
    for field in allowed_fields:
        if field in data:
            db[model_id][field] = data[field]
    
    db[model_id]['updated_at'] = datetime.utcnow().isoformat()
    save_models_db(db)
    
    return jsonify({
        'success': True,
        'message': 'Model updated successfully',
        'model': db[model_id]
    })


@app.route('/api/stats', methods=['GET'])
def get_stats():
    """Get overall system statistics"""
    db = load_models_db()
    log = load_access_log()
    
    # Calculate total folder size
    total_size = get_total_models_size()
    
    # Count total requests
    total_requests = len(log)
    
    # Get unique users
    unique_users = len(set(entry['username'] for entry in log))
    
    # Recent activity (last 10)
    recent_activity = log[-10:][::-1] if log else []
    
    return jsonify({
        'success': True,
        'stats': {
            'total_models': len(db),
            'total_size_bytes': total_size,
            'total_requests': total_requests,
            'unique_users': unique_users,
            'recent_activity': recent_activity
        }
    })


@app.route('/api/health', methods=['GET'])
def health():
    """Health check endpoint"""
    return jsonify({
        'success': True,
        'status': 'healthy',
        'mount_point': MODELS_MOUNT_POINT
    })


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)

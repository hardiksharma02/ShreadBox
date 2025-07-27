# 🔒 ShreadBox - Secure Self-Destructing File Sharing

ShreadBox is a secure, privacy-focused file sharing service that automatically destroys files after they've been downloaded or expired. Perfect for sharing sensitive documents that shouldn't persist indefinitely.

## 🌟 Features

- **🔐 End-to-End Encryption**: Files are encrypted using AES-GCM before storage
- **⏳ Time-Based Self-Destruction**: Files automatically delete after a specified time
- **🔢 Download Limits**: Set maximum number of downloads allowed
- **📝 Optional Messages**: Attach encrypted messages with your files
- **🚫 Zero Storage**: Files are permanently deleted after expiration/download limit
- **🔍 No Tracking**: No logs of file contents or user data
- **🚀 Simple API**: RESTful API for easy integration

## 🛠 Tech Stack

- **Backend**: Go 1.22+
- **Framework**: Gin Web Framework
- **Storage**: Local File System (configurable)
- **Security**: AES-GCM Encryption
- **API**: RESTful with JSON

## 📋 Prerequisites

- Go 1.22 or higher
- Docker (optional, for containerization)
- Make (optional, for using Makefile commands)

## 🚀 Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/shreadbox.git
   cd shreadbox
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

4. **Run the server**
   ```bash
   go run cmd/main.go
   ```

The server will start at `http://localhost:8080`

## 🐳 Docker Deployment

1. **Build the image**
   ```bash
   docker build -t shreadbox .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 shreadbox
   ```

## 📡 API Endpoints

### Upload File
```http
POST /api/upload
Content-Type: multipart/form-data

file: [file]
expiry_time: "24h"
downloads_allowed: 1
message: "Optional message"
```

### Download File
```http
GET /api/download/:token
```

### Check Status
```http
GET /api/status/:token
```

## ⚙️ Configuration

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `PORT` | Server port | 8080 |
| `MAX_FILE_SIZE` | Maximum file size in MB | 10 |
| `STORAGE_PATH` | Path to store files | ./storage |
| `CLEANUP_INTERVAL` | Cleanup check interval | 5m |

## 🔒 Security Features

- AES-GCM encryption for all stored files
- Automatic file shredding after expiry/download
- Rate limiting on all endpoints
- File size restrictions
- HTTPS enforcement in production
- No persistent storage of encryption keys

## 🧪 Development

### Running Tests
```bash
go test ./...
```

### Code Quality
```bash
go fmt ./...
go vet ./...
```

## 📜 License

MIT License - see [LICENSE](LICENSE) for details

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 🔍 Architecture

```
secure-file-share/
├── cmd/                 → App entry point
├── internal/
│   ├── encryption/      → File encryption/decryption
│   ├── handlers/        → HTTP handlers
│   ├── storage/         → File storage
│   └── cleanup/         → Self-destruct system
├── web/                 → Frontend templates
├── config/             → Configuration
└── main.go            → Application entry
```

## 📚 Documentation

Full API documentation is available at `/docs` when running the server.

## ⚠️ Important Notes

- Files are automatically deleted after expiration or download limit
- No recovery of expired/deleted files is possible
- Maximum file size is configurable (default 10MB)
- Service is intended for temporary file sharing only

## 🌟 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Google UUID](https://github.com/google/uuid)

## 📧 Contact

For bugs or feature requests, please open an issue.

# Resume Analyzer 📄🔍

Resume Analyzer is a Go-based application designed to analyze and extract useful information from resumes. It leverages Docker, Cassandra, and MinIO to provide a scalable and efficient solution for processing resumes.

## Features ✨

- **Go-based Backend**: Utilizes the Go programming language for high performance.
- **Docker**: Ensures consistent development and deployment environments.
- **Cassandra**: A highly scalable NoSQL database for storing resume data.
- **MinIO**: An object storage server compatible with Amazon S3 for storing resume files.

## Getting Started 🚀

Follow these instructions to set up and run the project on your local machine.

### Prerequisites

- [Go](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/adityagpramanik/resume-analyzer.git
   cd resume-analyzer
   ```
2. **Build and start the services:**

3. **Run the application:**


### Project Structure 📂
common.services/ - Shared services and utilities.
controllers/ - Handlers for the API endpoints.
routes/ - API route definitions.
infra-compose.yaml - Docker Compose configuration.
main.go - Entry point of the application.
Contributing 🤝
We welcome contributions to improve Resume Analyzer! Here’s how you can help:

### Fork the repository.
Create a new branch (git checkout -b feature-branch).
Make your changes and commit them (git commit -m 'Add new feature').
Push to the branch (git push origin feature-branch).
Open a Pull Request.

### License 📜
This project is licensed under the MIT License. See the LICENSE file for details.

### Contact 📬
For any inquiries or feedback, please reach out to [Aditya Pramanik](https://github.com/adityagpramanik).

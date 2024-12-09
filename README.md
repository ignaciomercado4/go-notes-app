### Notitas - Go notes app
---
This repository contains a simple note-taking app, developed using Go, Gin, GORM, text/template, and SQLite. It was built primarily for learning purposes and experimenting with these technologies.

### What can I do with it?
---
You can create users with this app and log in to write notes with a title and content. User passwords are hashed with bcrypt and sessions managed with JWT.

### Installation
---
1. Clone the repo:
```
git clone https://github.com/ignaciomercado4/go-notes-app.git
cd go-notes-app
```

2. Install dependencies:
```
go mod tidy
```

3. Create a database:
```
touch test.db
```  

4. Run the app:
```
go run *.go
```

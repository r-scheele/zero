<div align="center">
  <img src="public/static/logo.png" alt="Zero Logo" width="120" height="120">
  
  
  [![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
  [![Echo](https://img.shields.io/badge/Echo-v4-00ADD8?style=for-the-badge)](https://echo.labstack.com/)
  [![HTMX](https://img.shields.io/badge/HTMX-Latest-3366CC?style=for-the-badge)](https://htmx.org/)
  [![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
  
  **A comprehensive learning platform for students and educators** 🎓
  
  *Completely free, comprehensive learning platform covering quizzes, document management, and study tools*
</div>

---

## 🎓 What is Zero?

Zero is a **note-taking and content management platform** designed for students and educators. It provides concrete, no-nonsense tools for organizing study materials, creating notes, and managing content with modern web technologies.

## 📊 Database Schema

### 🗄️ Core Entities
- **Users**: Phone-based authentication with email support
- **Notes**: Rich content with title, description, and resources
- **NoteLikes**: Social interaction tracking
- **NoteReposts**: Content sharing with optional comments
- **PasswordTokens**: Secure password reset functionality

### 🔗 Relationships
- Users can create multiple notes
- Notes can have multiple likes and reposts
- Users can like and repost notes from others
- Password tokens are linked to specific users

### 📁 File Management
- **Resource Attachment**: JSON array of file and URL resources
- **Cloud Storage**: Multi-provider support (AWS S3, GCS, Azure)
- **File Validation**: Size limits and type checking
- **Async Processing**: Background file handling

### 🔍 Search & Discovery
- **User Search**: Find users by name or phone
- **Note Search**: Full-text search across note content
- **Admin Filters**: Advanced filtering in admin panel
- **Public Feed**: Browse publicly available notes

## 🎯 Current Implementation

### 📝 Note Management
- **Create & Edit**: Rich note creation with title, description, and content
- **Resource Attachment**: Add URLs and file uploads to notes
- **Visibility Control**: Set notes as public or private
- **Share Links**: Generate unique sharing tokens for notes
- **AI Processing**: Background AI curriculum generation (placeholder)

### 👥 Social Features
- **Like System**: Users can like notes from others
- **Repost Feature**: Share notes with optional comments
- **User Profiles**: Basic user management with phone-based authentication
- **Public Feed**: Browse and interact with public notes

### 🔐 Authentication & Security
- **Phone Verification**: WhatsApp-based registration and login
- **Password Reset**: Email-based password recovery
- **Secure Tokens**: JWT-based authentication system
- **Admin Access**: Administrative panel for user and content management

### 🎨 User Interface
- **Responsive Design**: Mobile-first approach with TailwindCSS
- **Component-Based**: Type-safe HTML with Gomponents
- **Interactive Elements**: HTMX for dynamic interactions
- **Modern Styling**: DaisyUI components for consistent design

---

## 🛠️ Tech Stack

- **Backend**: Go with Fiber framework
- **Database**: PostgreSQL with Ent ORM
- **Frontend**: HTMX + DaisyUI + TailwindCSS
- **Components**: Gomponents for type-safe HTML
- **Authentication**: Phone-based with WhatsApp verification
- **File Storage**: Multi-provider support (AWS S3, GCS, Azure Blob)
- **Email**: SMTP integration for password reset
- **Background Jobs**: Async task processing with worker queues
- **Admin Panel**: User and content management interface
- **Mobile Integration**: WhatsApp API for phone verification

### Development Tools
- **[Air](https://github.com/air-verse/air)** - Live reloading
- **Make** - Build automation
- **Go 1.24+** - Latest Go features

---

## 🎯 Quick Start

### Prerequisites
- [Go 1.24+](https://go.dev/)
- [Make](https://www.gnu.org/software/make/) (optional but recommended)

### 1. Get the Code
```bash
git clone https://github.com/r-scheele/zero.git
cd zero
```

### 2. Install Dependencies
```bash
make install  # Installs Ent, Air, and Tailwind CSS
```

### 3. Create Admin Account
```bash
make admin phone=+1234567890
# Note the generated password from console output
```

### 4. Start Development Server
```bash
make watch  # With live reloading
# OR
make run    # Standard run
```

🎉 **That's it!** Visit `http://localhost:8000` to see your application.

---

## 📸 Screenshots

<details>
<summary>🖼️ View Screenshots</summary>

### User Registration with Validation
<img src="https://raw.githubusercontent.com/r-scheele/readmeimages/main/zero/register.png" alt="Registration" width="600">

### Interactive Modal with HTMX
<img src="https://raw.githubusercontent.com/r-scheele/readmeimages/main/zero/modal.png" alt="Modal" width="600">

### Admin Panel - User Management
<img src="https://raw.githubusercontent.com/r-scheele/readmeimages/main/zero/admin-user_list.png" alt="User List" width="600">

### Admin Panel - User Editing
<img src="https://raw.githubusercontent.com/r-scheele/readmeimages/main/zero/admin-user_edit.png" alt="User Edit" width="600">

### Background Task Monitoring
<img src="https://raw.githubusercontent.com/r-scheele/readmeimages/main/backlite/failed.png" alt="Task Queue" width="600">

</details>

---

## 🏗️ Project Structure

```
zero/
├── cmd/
│   ├── admin/          # Admin CLI tools
│   └── web/            # Web server entry point
├── config/             # Configuration management
├── ent/                # Database entities and ORM
├── pkg/
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # Custom middleware
│   ├── services/       # Business logic services
│   ├── ui/             # UI components and layouts
│   └── tasks/          # Background tasks
├── public/             # Static assets
└── uploads/            # File uploads
```

---

## 🔧 Development

### Available Commands
```bash
make help              # Show all available commands
make install           # Install all dependencies
make run              # Start the application
make watch            # Start with live reloading
make test             # Run tests
make css              # Build Tailwind CSS
make build            # Build production binary
make ent-gen          # Generate ORM code
make ent-new name=X   # Create new entity
make admin phone=X    # Create admin user
```

### Creating New Entities
```bash
# Create a new database entity
make ent-new name=Product

# Edit the schema in ent/schema/product.go
# Then generate the code
make ent-gen
```

### Development Workflow
1. Create handler in `pkg/handlers/`
2. Add route in `pkg/handlers/router.go`
3. Create page component in `pkg/ui/pages/`
4. Add navigation if needed

---

## 🎨 UI Development

### Component-Based Architecture
Zero uses Gomponents to write HTML in Go, providing type safety and reusability:

```go
func MyComponent(title string) Node {
    return Div(
        Class("card bg-base-100 shadow-xl"),
        Div(
            Class("card-body"),
            H2(Class("card-title"), Text(title)),
            P(Text("Component content here")),
        ),
    )
}
```

### HTMX Integration
Add interactivity without JavaScript:

```go
Button(
    Class("btn btn-primary"),
    Attr("hx-post", "/api/action"),
    Attr("hx-target", "#result"),
    Text("Click Me"),
)
```

### Styling with DaisyUI
Use semantic component classes:

```go
Div(
    Class("hero min-h-screen bg-base-200"),
    Div(
        Class("hero-content text-center"),
        H1(Class("text-5xl font-bold"), Text("Hello World")),
    ),
)
```

---

## 🔐 Authentication & Authorization

### Features
- ✅ User registration with email verification
- ✅ Secure login/logout
- ✅ Password reset via email
- ✅ Phone number verification
- ✅ Admin role management
- ✅ Session management
- ✅ CSRF protection

### Usage
```go
// Protect routes with authentication
protected := e.Group("/dashboard")
protected.Use(middleware.RequireAuth)

// Admin-only routes
admin := e.Group("/admin")
admin.Use(middleware.RequireAdmin)
```

---

## 📊 Admin Panel

The admin panel provides comprehensive tools for educators and administrators:
- 👥 **Student Management** - View, edit, and manage student accounts
- 📚 **Content Management** - Upload and organize study materials
- 🧠 **Quiz Administration** - Create, edit, and monitor quiz performance
- 📊 **Analytics Dashboard** - Track student progress and engagement
- 🔍 **Advanced Search** - Filter by name, email, course, progress
- 📱 **Mobile Responsive** - Manage your platform from any device
- 🎨 **Intuitive Interface** - Clean, educator-friendly design

### Educational Features
- Student enrollment management
- Course and material organization
- Quiz creation and grading
- Progress tracking and reporting
- Bulk operations for efficiency
- Real-time student activity monitoring

### Access
1. Create admin account: `make admin phone=+1234567890`
2. Login at `/login`
3. Access admin panel at `/admin`

---

## 🗄️ Database

### Entity Definition
```go
// ent/schema/user.go
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").NotEmpty(),
        field.String("email").Unique(),
        field.String("phone_number").Optional(),
        field.Bool("verified").Default(false),
    }
}
```

### Querying
```go
// Get users with filters
users, err := client.User.
    Query().
    Where(user.NameContainsFold("john")).
    Order(ent.Asc(user.FieldCreatedAt)).
    Limit(10).
    All(ctx)
```

---

## 📧 Email System

### Templates
Email templates are written in Go using Gomponents:

```go
func WelcomeEmail(userName string) Node {
    return HTML(
        Head(Title(Text("Welcome!"))),
        Body(
            H1(Text("Welcome "+userName)),
            P(Text("Thanks for joining us!")),
        ),
    )
}
```

### Sending
```go
err := mailService.Send(
    "user@example.com",
    "Welcome!",
    WelcomeEmail("John"),
)
```

---

## 🔄 Background Tasks

### Define Tasks
```go
type EmailTask struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

func (t EmailTask) Handle(ctx context.Context) error {
    return sendEmail(t.To, t.Subject, t.Body)
}
```

### Queue Tasks
```go
task := EmailTask{
    To:      "user@example.com",
    Subject: "Welcome!",
    Body:    "Welcome to our platform!",
}

err := taskService.Queue(task)
```

---

## 🚀 Deployment

### Local Development

```bash
# Set up environment
cp .env.example .env
# Edit .env with your configuration

# Run database migrations
make migrate

# Create admin user
make admin phone=+1234567890

# Start development server
make dev
```

### Production Deployment

```bash
# Build binary
make build

# Run production server
./bin/zero
```

### Docker Deployment

```bash
# Build Docker image
docker build -t zero-notes .

# Run with environment variables
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@localhost/zero" \
  -e SMTP_HOST="smtp.gmail.com" \
  zero-notes
```

### Environment Configuration

```bash
# Database
DATABASE_URL="postgres://user:password@localhost:5432/zero"

# SMTP for password reset
SMTP_HOST="smtp.gmail.com"
SMTP_PORT=587
SMTP_USERNAME="your-email@gmail.com"
SMTP_PASSWORD="your-app-password"

# WhatsApp Integration (360dialog)
WHATSAPP_API_KEY="your-360dialog-api-key"
WHATSAPP_CHANNEL_ID="your-channel-id"

# Cloud Storage (optional)
AWS_ACCESS_KEY_ID="your-aws-key"
AWS_SECRET_ACCESS_KEY="your-aws-secret"
AWS_REGION="us-east-1"
S3_BUCKET="your-bucket"

# Or Google Cloud Storage
GCS_BUCKET="your-gcs-bucket"
GOOGLE_APPLICATION_CREDENTIALS="path/to/service-account.json"

# Or Azure Blob Storage
AZURE_STORAGE_ACCOUNT="your-account"
AZURE_STORAGE_KEY="your-key"
AZURE_CONTAINER="your-container"
```

---

## 🤝 Contributing

We welcome contributions from developers interested in note-taking and content management platforms!

### 🎯 Focus Areas
- **Core Features**: Note management, social features, file handling
- **User Experience**: Interface improvements, mobile responsiveness
- **Performance**: Database optimization, background job processing
- **Integration**: Cloud storage providers, notification systems
- **API Development**: REST endpoints, authentication improvements

### 📋 Development Guidelines
- Follow Go best practices and conventions
- Write tests for new features
- Update documentation for any changes
- Ensure mobile-first responsive design
- Maintain type safety with Gomponents

### 🚀 Getting Started
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following the existing patterns
4. Test thoroughly with different user scenarios
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request with detailed description

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Credits

Zero is built on the shoulders of giants. Special thanks to:

- [Echo](https://echo.labstack.com/) - Web framework
- [Ent](https://entgo.io/) - ORM
- [Gomponents](https://github.com/maragudk/gomponents) - HTML in Go
- [HTMX](https://htmx.org/) - Modern web interactions
- [Alpine.js](https://alpinejs.dev/) - Minimal JavaScript
- [DaisyUI](https://daisyui.com/) - Beautiful components
- [Tailwind CSS](https://tailwindcss.com/) - Utility CSS
- [Backlite](https://github.com/mikestefanello/backlite) - Background tasks

---

<div align="center">
  <p>Made with ❤️ for education and built with Go</p>
  <p>Empowering students and educators worldwide</p>
  <p><a href="https://github.com/r-scheele/zero">⭐ Star us on GitHub</a> | <a href="#-contributing">🤝 Contribute to Education</a></p>
</div>
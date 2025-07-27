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

Zero is a **comprehensive learning platform** designed for students and educators. It provides concrete, no-nonsense tools for learners in a hurry, covering everything from interactive quizzes and document management to progress tracking and collaborative study features. Built with modern web technologies, Zero makes learning accessible, efficient, and engaging.

### ✨ Key Features

- 📚 **Study Materials** - Organize and access your documents and learning resources
- 🧠 **Smart Quizzes** - Test your knowledge effectively with interactive quizzes
- 📊 **Progress Tracking** - Monitor your learning journey and achievements
- 👥 **Collaboration** - Learn together with others in a shared environment
- 🔐 **User Management** - Secure registration, login, and profile management
- 📱 **Mobile-First Design** - Beautiful, responsive interface for all devices
- 👑 **Admin Panel** - Comprehensive management tools for educators
- 🗄️ **Document Management** - Upload, organize, and share study materials
- 🔍 **Advanced Search** - Find content quickly with powerful search capabilities
- 📧 **Notifications** - Stay updated with email and in-app notifications
- 🎯 **Zero JavaScript Required** - Built with modern web technologies
- 🔄 **Real-time Updates** - Live interactions without page refreshes

---

## 🛠️ Tech Stack

### Backend
- **[Echo](https://echo.labstack.com/)** - High-performance web framework
- **[Ent](https://entgo.io/)** - Type-safe ORM with code generation
- **[SQLite](https://sqlite.org/)** - Embedded database (easily swappable)
- **[Backlite](https://github.com/mikestefanello/backlite)** - Background task processing

### Frontend
- **[Gomponents](https://github.com/maragudk/gomponents)** - HTML components in pure Go
- **[HTMX](https://htmx.org/)** - Modern interactions without JavaScript
- **[Alpine.js](https://alpinejs.dev/)** - Minimal JavaScript framework
- **[DaisyUI](https://daisyui.com/)** - Beautiful Tailwind CSS components
- **[Tailwind CSS](https://tailwindcss.com/)** - Utility-first CSS framework

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

### Adding New Pages
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

### Educational Institution Setup

Zero is designed to be easily deployed in educational environments:

#### Production Build
```bash
# Build the application
make build

# The binary will be created as ./tmp/main
./tmp/main
```

#### Docker Deployment for Schools
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/tmp/main .
COPY --from=builder /app/public ./public
CMD ["./main"]
```

#### Configuration for Educational Use
```bash
# Database
DATABASE_URL=sqlite://./school_data.db

# Server
APP_PORT=8000
APP_DOMAIN=learning.yourschool.edu

# Email for notifications
SMTP_HOST=smtp.yourschool.edu
SMTP_PORT=587
SMTP_USERNAME=noreply@yourschool.edu
SMTP_PASSWORD=your-institutional-password

# Security
SESSION_SECRET=your-secure-school-key
CSRF_SECRET=your-csrf-secret

# Educational settings
SCHOOL_NAME="Your Institution Name"
ADMIN_EMAIL=admin@yourschool.edu
```

---

## 🤝 Contributing

We welcome contributions from educators, developers, and students! Here's how you can help improve the learning experience:

### Development Setup
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/educational-feature`
3. Make your changes
4. Run tests: `make test`
5. Commit your changes: `git commit -m 'Add educational feature'`
6. Push to the branch: `git push origin feature/educational-feature`
7. Open a Pull Request

### Guidelines
- Follow Go conventions and best practices
- Consider educational use cases and accessibility
- Write tests for new features
- Update documentation as needed
- Keep the student/educator experience in mind

### Areas for Contribution
- 🎓 Educational features (quiz types, study tools)
- 🐛 Bug fixes and performance improvements
- 📚 Documentation and tutorials
- 🎨 UI/UX enhancements for better learning
- ♿ Accessibility improvements
- 🌐 Internationalization for global education
- 📊 Analytics and progress tracking features
- 🔧 Administrative tools for educators

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
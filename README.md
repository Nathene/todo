# Todo/Calendar/Notes app

I wanted some aspects of Jira, some aspects of a project used at my university which showed you all of your upcoming assignment due dates, and just a regular brainstorming note section all into one app. I couldnt find this online so i decided to make it!

## Technologies Used

- **Go**: The core programming language used for building the backend.
- **Echo**: A high-performance, extensible web framework for Go, used for routing and middleware.
- **SQLite**: A lightweight, disk-based database used for data storage.
- **HTML/CSS**: For building responsive and user-friendly interfaces.
- **JavaScript**: Enhances interactivity and user experience.
- **Bootstrap**: A front-end framework for styling and responsive design.

## Capabilities

- **User Authentication**: Secure login and account management to protect user data.
- **Todo Management**: Create and manage todo groups, lists, and subtasks with ease.
- **Event Calendar**: Add and view events to keep track of important dates.
- **Notes**: Create, edit, and delete notes for additional information storage.
- **Dark Mode**: Toggle between light and dark themes for a personalized user experience.

## API Endpoints

- **Public Endpoints**:
  - `/`: Dashboard access.
  - `/login`: User login.
  - `/create-account`: Account creation.
  - `/logout`: User logout.

- **Protected Endpoints** (require authentication):
  - `/settings`: Manage user settings.
  - `/groups`: Manage todo groups.
  - `/todos`: Manage individual todos and subtasks.
  - `/calendar`: Manage events.
  - `/notes`: Manage notes.

## Running the Application

To run the application locally, clone and then use this as an example main.go.

```golang
   package main

    import (
        "log"
        "os"
        "todo/internal/api/routes"
        "todo/internal/db"
        "todo/internal/pkg/renderer"
        "todo/internal/util"

        "github.com/labstack/echo/v4"
        "github.com/labstack/echo/v4/middleware"
    )

    const (
        port = "8040"
    )

    func main() {
        // Initialize the database
        db, err := db.InitDB()
        if err != nil {
            log.Fatal(err)
        }
        defer util.Defer(db.Close)

        // Create a new Echo instance
        e := echo.New()

        // Initialize the custom renderer with the templates directory
        e.Renderer = renderer.NewRenderer("templates")

        // Serve static files like styles.css
        e.Static("/static", "templates/static")

        // Set up logging to stdout
        e.Logger.SetOutput(os.Stdout)
        e.Use(middleware.Logger()) // Enable request logging middleware

        // Recover from panics
        e.Use(middleware.Recover())

        // Set up routes
        routes.Handle(e, db)

        // Start the server
        e.Logger.Fatal(e.Start(":" + port))
    }
```
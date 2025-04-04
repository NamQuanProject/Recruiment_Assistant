# `cmd` Directory

The `cmd` directory is used to contain the main applications for your Go project. Each subdirectory within `cmd` represents a separate executable.

## Purpose

The purpose of the `cmd` folder is to organize the entry points of your application. This structure helps in maintaining a clean and modular codebase by separating the application logic from the executable code.

## Guidelines

- Each subdirectory under `cmd` should contain a `main.go` file.
- Keep the code in `main.go` minimal. It should primarily handle configuration, dependency injection, and starting the application.
- Place reusable code in other packages outside the `cmd` directory to promote reusability and maintainability.

## Example Structure

```plaintext
cmd/
├── app1/
│   └── main.go
├── app2/
│   └── main.go
```

## Best Practices

- Avoid placing business logic directly in the `cmd` directory.
- Use descriptive names for subdirectories to indicate the purpose of each executable.
- Ensure proper error handling and logging in the `main.go` files.

## References

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

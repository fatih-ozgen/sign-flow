# User Management API

This is a simple User Management API built with Go and PostgreSQL. It provides basic functionality for user signup, signin, and retrieval of all users.

## Features

- User signup with unique membership ID generation
- User signin
- Retrieve all users
- Password hashing for security

## Prerequisites

- Go 1.16+
- PostgreSQL

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/your-username/user-management-api.git
   ```

2. Navigate to the project directory:
   ```
   cd user-management-api
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

4. Set up your PostgreSQL database and update the connection string in `database.go`:
   ```go
   connStr := "user=your_username password=your_password dbname=your_dbname sslmode=disable"
   ```

5. Run the application:
   ```
   go run .
   ```

The server will start on `http://localhost:8080`.

## API Endpoints

- POST `/signup`: Create a new user
  - Request body: `{"username": "example", "password": "password123"}`
  - Response: `{"message": "User created successfully", "membership_id": "ABCD1234EFGH5678"}`

- POST `/signin`: Authenticate a user
  - Request body: `{"username": "example", "password": "password123"}`
  - Response: `{"message": "Sign in successful"}`

- GET `/users`: Retrieve all users
  - Response: `[{"membership_id": "ABCD1234EFGH5678", "username": "example"}]`

## Project Structure

- `main.go`: Entry point of the application
- `database.go`: Database connection and operations
- `handlers.go`: HTTP request handlers
- `models.go`: Data structures
- `utils.go`: Utility functions

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)

# Backend Crowdfunding

Backend Crowdfunding is a backend application designed for crowdfunding campaigns and donation systems. It is built with **Golang**, using the **Gin** framework and **MySQL** as the database.

## Features

1. **User Management**
    - Register a user
    - Login
    - Check email availability
    - Upload user avatar
    - Fetch user details

2. **Campaign Management**
    - List all campaigns
    - Get campaign details
    - Create a campaign
    - Update a campaign
    - Upload campaign images

3. **Transaction Management**
    - List transactions for a campaign
    - List user transactions
    - Create a transaction
    - Handle transaction notifications (Midtrans integration)

4. **Authentication**
    - JWT-based authentication middleware

## Technologies Used

- **Golang**: Main programming language
- **Gin**: Web framework
- **GORM**: ORM for MySQL
- **MySQL**: Database backend
- **JWT**: Token-based authentication
- **CORS**: Middleware for handling Cross-Origin Resource Sharing
- **Midtrans**: Payment gateway integration

## Project Structure

```plaintext
backend-crowdfunding/
├── auth/               # Authentication module
├── campaign/           # Campaign module
├── handler/            # HTTP handlers
├── helper/             # Helper functions
├── migrations/         # Migrations files
├── payment/            # Payment module
├── transaction/        # Transaction module
├── user/               # User module
└── main.go             # Application entry point
```

## How to Run

1. Clone this repository
2. Copy `.env.example` to `.env` and adjust the configuration
3. Run migrations up with the following command:
    ```bash
     migrate -path ./migrations -database "mysql://username:password@tcp(127.0.0.1:3306)/database_name?multiStatements=true" up
    ```
4. if you don't have migrate, you can install it by running this command:
    ```bash
    go get -tags 'mysql' -u github.com/golang-migrate/migrate/cmd/migrate
    ```
5. Run the application with the following command:
    ```bash
    go run main.go
    ```
6. The application will be available at `http://localhost:8080`
7. You can test the API using Postman or any other API testing tools

    
Don't forget to follow my github account for more updates and star this repository if you like it, thank you!

Happy coding! 🚀
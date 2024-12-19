# SSO Backend Service

This is a Minimal Viable Product (MVP) for a Single Sign-On (SSO) backend service implemented in Go. It's designed to be deployed on Vercel as a serverless application.

## Features

- User registration with organization creation
- User authentication (login endpoint)
- JWT token generation and validation
- Email-based authentication method
- PostgreSQL database integration
- Protected endpoint demonstration
- User activation system

## Project Structure

```bash
sso-backend/
├── api/
│   └── index.go        # API setup and route configuration
├── auth/
│   └── auth.go         # JWT token generation and validation
├── db/
│   └── db.go          # Database connection and initialization
├── handler/
│   ├── activation.go   # User activation endpoint
│   ├── login.go        # User authentication
│   ├── protected.go    # Protected endpoint example
│   └── register.go     # User and organization registration
├── go.mod              # Go module file
├── go.sum              # Go module checksums
├── main.go             # Application entry point
├── vercel.json         # Vercel configuration
└── README.md           # This file
```

## Endpoints

1. `/login` (POST): Authenticates a user and returns a JWT token
   - Request body: `email` and `password` (form-data)
   - Response: JSON with `token` field or error message

2. `/register` (POST): Register new user with a new organization
   - Request body: `organization_name`, `email`, `password` (JSON)
   - Response: JSON with success/error message

3. `/activation` (POST): Activate user account
   - Request body: TBD
   - Response: JSON with status

4. `/protected` (GET): A protected endpoint requiring valid JWT token
   - Header: `Authorization: Bearer <token>`
   - Response: JSON with `message` field

## Database Schema

The service uses PostgreSQL with the following schema in the "z-auth" schema:
- organizations
- users
- userorganizations
- authmethods
- usercredentials

## Local Development

1. Ensure you have Go installed (version 1.20 or later recommended)
2. Install PostgreSQL and set up the database
3. Clone this repository
4. Create a `.env` file with your `DATABASE_URL`
5. Run `go mod download` to install dependencies
6. Execute `go run main.go` to start the server

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string

## Deployment

This project is configured for deployment on Vercel:

1. Install the Vercel CLI: `npm i -g vercel`
2. Configure your environment variables in Vercel
3. Run `vercel` and follow the prompts

## Security Features

- Password hashing using bcrypt
- JWT-based authentication
- Database transaction support for data integrity
- User activation system
- Secure credential storage

## Future Improvements

- Implement refresh token functionality
- Add email verification system
- Enhance error handling and logging
- Add comprehensive unit tests
- Implement rate limiting
- Add organization management features
- Enhance user activation flow

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)
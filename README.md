# SSO Backend Service

This is a Minimal Viable Product (MVP) for a Single Sign-On (SSO) backend service implemented in Go. It's designed to be deployed on Vercel as a serverless application. This repo developed in tandem with Claude AI.

## Features

- User authentication (login endpoint)
- JWT token generation and validation
- Protected endpoint demonstration

## Project Structure

```bash
sso-backend/
├── api/
│   ├── login.go    # Handles user authentication
│   └── protected.go # Demonstrates a protected endpoint
├── auth/
│   └── auth.go     # JWT token generation and validation
├── go.mod          # Go module file
├── go.sum          # Go module checksums
├── vercel.json     # Vercel configuration
└── README.md       # This file
```

## Endpoints

1. `/login` (POST): Authenticates a user and returns a JWT token
   - Request body: `username` and `password` (form-data)
   - Response: JSON with `token` field

2. `/protected` (GET): A protected endpoint that requires a valid JWT token
   - Header: `Authorization: Bearer <token>`
   - Response: JSON with `message` field

3. `/register` (POST): Register new user with a new organization
   - Request body: `organization_name`, `email`, `password` (json)
   - Response: JSON with `message` field

## Local Development

1. Ensure you have Go installed (version 1.20 or later recommended)
2. Clone this repository
3. Navigate to the project directory
4. Run `go mod download` to install dependencies
5. To test endpoints locally, you can use `go run api/login.go` and `go run api/protected.go` in separate terminal windows

## Deployment

This project is configured for deployment on Vercel:

1. Install the Vercel CLI: `npm i -g vercel`
2. Navigate to the project directory
3. Run `vercel` and follow the prompts

## Security Notes

- This is a basic implementation for demonstration purposes. In a production environment, implement proper security measures.
- The JWT secret is hardcoded in this example. Use environment variables for sensitive information in a real-world scenario.
- The login credentials are hardcoded. Implement proper user management and password hashing for a production application.

## Future Improvements

- Implement refresh token functionality
- Enhance error handling and logging
- Add unit tests
- Use environment variables for configuration

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)

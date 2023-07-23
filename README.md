# Dating Apps Backend Service with Go

## Overview

This repository contains a Golang backend dating user service that provides various functionalities related to user management and authentication. 
The service is built using the Go programming language and follows best practices for building scalable and maintainable backend services.

## Features

The user service includes the following features:

- User Auth: Allows users to register or login using their email address and receive an OTP (One-Time Password) for verification.

- Onboarding User : Registration for first time using apps. User can update personal information (name, birth date, gender), photos, hobby & interest, and location. 

- User Profile: Provides functionality to view the user's profile, including basic information and swipe count.

- Feeds / Profile Discovery: Provides functionality to view other user profiles data. For free account, User able to only view, swipe left (pass) and swipe right (like) 10 other dating profiles in total (pass + like) in 1 day. For Premium account, User able to view, swipe left (pass) and swipe right (like) with NO LIMIT.  

- User Swipes: Allows users to perform swipe actions on other profiles. Swipe Left for Pass and Swipe Right for Like.

- Purchase Premium: Allows users to purchase premium account.

## Project Structure

- cmd/ # Main application entry point
- internal/ # Internal application packages
- - internal/adapter/ # Adapters for database and external services
- - - internal/adapter/redis/ # Adapters for redis
- - - internal/adapter/repository/ # Adapters for database
- - internal/business/ # Business layer
- - - internal/business/domain/ # Business domain entities and data structures
- - - internal/business/port/ # Business layer
- - - internal/business/service/ # Business logic layer
- - internal/handler/ # Service implementations
- - - internal/handler/http/ # Http Service implementations
- migration/ # Database migration files
- pkg/ # Shared packages and utilities
- storage/ # Local file storage
- .env # Environment variables
- .env.example # Environment variables example
- .gitignore # Git ignore rules
- go.mod # Go module file
- go.sum # Go module checksum file
- Makefile # Makefile for build and run
- README.md # Readme file

## Installation and Setup

1. Clone the repository:
    `git clone git@github.com:ijlik/dating-user.git`
2. Install the required dependencies:
   `go mod download`
3. Set up the database:
   `make migrate-up`
4. Set up the environment variables:
   `cp .env.example .env`
5. Run the application:
   `make local-run`
6. Run the linter:
   `make lint`

## Testing
To run the unit tests, use the following command:
    `make test`

## Contributing
Contributions are welcome! If you find any issues or have suggestions for improvement, please feel free to open an issue or submit a pull request.
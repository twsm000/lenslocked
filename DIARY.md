# Web Development with Go - V2

## Chapter 1

- A Basic Web Application
- Throubleshooting and Slack
- Packages and Imports
- Editors and Automatic Imports
- The "Hello, World" Part of Our Code
- Web Requests
- HTTP Methods
- Our Handler Function
- Registering our Handler Function and Starting the Web Server
- Go Modules

## Chapter 2 - Adding New Pages

- Dynamic Reloading
- Setting Header Values
- Creating a Contact Page
- Examining the http.Request Type
- Custom Rounting
- url.Path vs url.RawPath
  - Acesse [urlencoder.com](urlencoder.com)
  - url.Path irá mostrar caminhos decodificados, caso algum seja fornecido na URL. Exemplo: / => %2F.
  - url.RawPath irá mostrar os caminhos na forma que foi fornecida. (codificada).
- Not Found Page
- The http.Handler Type
- The http.HandlerFunc Type
- Exploring Handler Conversions
- FAQ Exercise

## Chapter 3 - Routers and 3rd party libraries

- Router Requirements
- Using Git
- Installing Chi
- Using Chi
- Chi Exercises

## Chapter 4 - Templates

- What are Templates
- Why Use Server Side Rendering?
- Creating Our First Template
- Cross Site Scripting (XSS)
- Alternative Template Libs
- Contextual Encoding
- Home Page via Template
- Contact Page via Template
- FAQ Page via Template
- Template Exercises

## Chapter 5 - Code Organization

- Code Organization
- MVC Overview
- Walking Through a Web Request With MVC
- MVC Exercises

## Chapter 6 - Starting to Apply MVC

- Creating the views package
- fmt.Errorf
- Validating Templates at Startup
- Must Functions

## Chapter 7 - Enhancing Our Views

- Embedding Templates Files
- Variadic Parameters
- Named Templates
- Dynamic FAQ Page
- Reusable Layouts
- Tailwind CSS
- Utility-first CSS
- Adding a Navigation Bar
- Exercises

## Chapter 8 - The Signup Page

- Creating the Signup Page
- Styling the Signup Page
- Intro to REST
- Users Controllers
- Decouple with Interfaces
- Parsing the Signup Form
- URL Query Params
- Exercises

## Chapter 9 - Database and PostgreSQL

- Intro to Databases
- Installing Postgres
- Connecting to Postgres
- Update: Docker Container Names
- Creating SQL Tables
- Postgres Data Types
- Postgres Constraints
- Creating Users Table
- Inserting Records
  TZ and PGTZ was added to fix database CURRENT_TIMESTAMP;
- Querying Records
- Filtering Queries
- Updating Records
- Deleting Records
- Additional SQL Resources

## Chapter 10 - Using Postgres with Go

- Connecting to Postgres with GO
- Importing with Side Effects
- Postgres Config Type
- Executing SQL with GO
- Inserting Records
- SQL Injection
- Acquiring New Record IDs
- Querying a Single Record
- Creating Sample Orders
- Querying Multiple Records
- ORMs vs SQL
- Exercises: SQL with GO
- Syncing the Source Code

## Chapter 11 - Securing Passwords

- Steps for Securing Passwords
- Third Party Authentication Options
- What is a Hash Function
- Store Password Hashes
- Salt Passwords
- Learning bcrypt with a CLI
- Hashing Passwords with bcrypt

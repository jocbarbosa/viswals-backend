# Viswals Backend

This is the backend service for the Viswals project, written in Go.

## Table of Contents

- [Viswals Backend](#viswals-backend)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Usage](#usage)

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/jocbarbosa/viswals-backend.git
    ```
2. Navigate to the project directory:
    ```sh
    cd viswals-backend
    ```
3. Install dependencies:
    ```sh
    go mod tidy
    ```

## Usage

1. Run the reader application:
    ```
    make run-reader
    ```

2. Run the api/consumer application:
    ```
    make run-api
    ```
3. The server will start on `http://localhost:8080`.

4. To know how to use the api, check the documentation:
    `http://localhost:8080/docs`
<p align="center">
  <img src="https://img.icons8.com/?size=512&id=44442&format=png" height="150" alt="Go logo" />
</p>

## Go API

This is a simple API written in Go using the Fiber framework. It's a RESTful API that implements standard CRUD (Create, Read, Update, Delete) operations on a database of people. The API is split into two parts, the first part is the main logic of the API, and the second part is the configuration of the API to be run.

## Running the app

To run the app, you can use the following command:

```bash
    go run cmd/main.go
```

## Seed data

For seed data, you can use the Python scripts in the `seed` folder. The scripts use the `uv` package manager to install the required packages and seed the database with the data.

Before running the seed script, you need to have the `uv` package manager installed and configured seed number.

To run the seed script, you can use the following command:

```bash
    cd seeds
    uv run main.py
```

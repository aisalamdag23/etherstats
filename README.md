# EtherStats

EtherStats is a Go application that provides Ethereum network statistics: gas price, latest block number, and address balances via a REST API.

## Prerequisites

- Go 1.24+
- Docker & Docker Compose
- [direnv](https://direnv.net/) (optional, for environment management)

## Setup

1. **One-liner quick-start**
    Just need a fresh local stack? Run:
    ```sh
    make config
    ```
    and jump to **Start the server**.

2. **Step-by-step setup**
    - **Set up local config files**
        This will create `.config.yml` and `.envrc` if they don't already exist:
        ```sh
        make config
        ```
        > This also runs `direnv allow` to load the environment.
        
    - **Start services**
        To start all necessary services:
        ```sh
        make docker-start
        ```
        This will start:
        - PostgreDB
        - RedisDB

    - **Run DB Migration**
        To start database migration:
        ```sh
        make db-migrate
        ```
        
    - **Start the server**
        ```sh
        make start
        ```

3. **Where to find it**
    ⚡️ The API will be available at http://localhost:8080 (or the port configured in .config.yml).

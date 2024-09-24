# EWallet Application

## Overview

Приложение EWallet — это простая система платежных транзакций, реализованная как RESTful HTTP-сервер в Go. Это приложение позволяет пользователям создавать кошельки, переводить средства между кошельками, получать историю транзакций и проверять балансы кошельков. 

## Features

1. **Create Wallet**
   - **Endpoint:** `POST /api/v1/wallet`
   - **Response:** JSON object with the wallet ID and balance (initially set to 100.0).

2. **Send Money**
   - **Endpoint:** `POST /api/v1/wallet/{walletId}/send`
   - **Request Body:** JSON object containing:
     - `to`: ID of the recipient wallet
     - `amount`: Amount to transfer
   - **Response:** 
     - `200 OK` if successful
     - `404 Not Found` if the source wallet does not exist
     - `400 Bad Request` if the target wallet does not exist or insufficient funds

3. **Get Transaction History**
   - **Endpoint:** `GET /api/v1/wallet/{walletId}/history`
   - **Response:** 
     - `200 OK` with a JSON array of transaction objects if the wallet exists
     - `404 Not Found` if the wallet does not exist

4. **Get Wallet Info**
   - **Endpoint:** `GET /api/v1/wallet/{walletId}`
   - **Response:**
     - `200 OK` with a JSON object containing the wallet ID and balance if the wallet exists
     - `404 Not Found` if the wallet does not exist

## Requirements

- **Go Version:** 1.21
- **Database:** PostgreSQL

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/ewallet.git
   cd ewallet

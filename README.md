# Wallet Service


## Summary
Simple service created as a test task for golang developer position.

It allows you to create wallet, update the balance of a wallet and get current state of a wallet.

### Task Requirements
There is a way to create a wallet:
    
- No need for specific request data since it does not support auth or there are no furthermore details on wallet (like user_id, currency etc.)

There is a way to update balance of a specific wallet:
    
- You need to enter the amount for increasing or decreasing the balance (positive and negative integers)
It was decided both balance and amount to be integers since the idea is to save 
any currency amount with basic monetary units.
*When dealing with amounts (such as prices, or totals, or something that is going to be processed as a currency transaction) using integers will avoid floating point rounding errors.

There is a way to get current state of a specific wallet.

### Database
Postgresql database was chosen for the project because it is an open source object-relational database,
and program required the data to be categorized

<img width="678" alt="Schema" src="https://user-images.githubusercontent.com/37467776/150845876-6b63b4b8-f8dd-4395-af59-2064e412d53a.png">

Database schema can be found in /db/up.sql

## Limitations
- Multiple currencies inside account is not supported;
- No Messaging queue supported;
- No auth supported;
- No error tracing;
- Needs test coverage to be improved.


## Documentation

You can find detailed API documentation  (swagger) inside /api_docs/api.yaml

### Start the service and database with Docker
Go to the project dir and build container:
- docker-compose up --build

*wallet service starts on port :8080 (with docker of psql on default port :5432)



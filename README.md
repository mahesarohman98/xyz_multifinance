# XYZ Multifinance

A backend service for handling customer registration, credit limits, and loan transactions for PT XYZ Multifinance. Adopting Clean Architecture + Domain-Driven Design (DDD) principles.

----------

## Project Structure

```
/api        # API specifications (for generating clients, docs, etc.)
/docker     # Dockerfiles for development and production environments
/images     # Contain images for doc purpose
/notes      # Personal notes and scratchpad
/scripts    # Helper shell scripts
/sql        # Database schema and migration files
/src        # Main backend source code
```

----------

##  Running Tests

To run all unit tests:

```
make test
```

If you need to debug integration tests (e.g., database-related tests):

```
make mysql
```

(This will start a local MySQL server.)

----------

##  Running the Application

Start development environment:

```
make dev-up
```

Stop development environment:

```
make dev-down
```

----------

##  Git Flow Usage

I follow the Git Flow strategy:

-   **feature/** → new features   
-   **release/** → preparing a new release
-   **hotfix/** → urgent production fixes
-   **support/** → maintenance or LTS


----------

# Documentation

-   API documentation is generated from the `/api` folder.
    
-   Database schema lives under `/sql`.



## Loan Application Workflow

> **Important:**  
> 📢 **Always refer to the `/api/openapi` folder for the latest API documentation!**  
> Although it is currently only in `.yml` format, it defines the complete API structure for client generation and documentation purposes.

----------

## 1. Register a New Source

Register a new external source (e.g., ECommerce, Dealer, Web):
```bash
curl -X POST 'http://localhost:3000/api/v1/sources' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "secret": "topsecretpassword",
    "category": "ECommerce",
    "name": "Tokopedia",
    "email": "partner@tokopedia.com"
}'
```
## 2. Get Access Token for Source

Request a token by providing the Source ID and Secret:
```bash
curl -X POST 'http://localhost:3000/api/v1/sources/token' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "sourceId": "65d85484-5222-45cd-a2b7-7a65f6a6dea4",
    "sourceSecret": "topsecretpassword"
}'
```
✅ Response will include a Bearer token to authenticate future requests.

## 3. Register a New Customer

Register a new customer using the token obtained above:
```bash
curl -X POST 'http://localhost:3000/api/v1/customers' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer <YOUR_ACCESS_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "nik": "2334567890123456",
    "fullName": "John Doe",
    "legalName": "Jonathan Doe",
    "placeOfBirth": "Jakarta",
    "dateOfBirth": "1990-02-20",
    "wage": 5000000
}'

```
## 4. Set Initial Credit Limits for Customer

Assign available credit limits by tenor (duration in months):
```bash
curl -X POST 'http://localhost:3000/api/v1/customers/<CUSTOMER_ID>/credit-limits' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer <YOUR_ACCESS_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '[
    {
      "monthRange": 1,
      "limitAmount": 3000000
    },
    {
      "monthRange": 2,
      "limitAmount": 5000000
    },
    {
      "monthRange": 3,
      "limitAmount": 6000000
    },
    {
      "monthRange": 6,
      "limitAmount": 9000000
    }
]'
```
## 5. Submit a Loan Request

Customers can request loans based on the available limits:
```bash
curl -X 'POST' \
  'http://localhost:3000/api/v1/loans' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjYXRlZ29yeSI6IkVDb21tZXJjZSIsImV4cCI6MTc0NTc4MDkxOCwibmFtZSI6IlRva29wZWRpYSIsInN1YiI6IjY1ZDg1NDg0LTUyMjItNDVjZC1hMmI3LTdhNjVmNmE2ZGVhNCJ9.0_4lN9lCoXU0A4yB9vWhBbZzzOO5lonyWT_hl86V8Gs' \
  -H 'Content-Type: application/json' \
  -d '{
  "customerId": "51f54526-213a-460f-9f76-2ecc7db7949d",
  "externalId": "external-56789",
  "tenor": 6,
  "loans": [
    {
      "contractNumber": "CTR-001",
      "otr": 5000000,
      "amountInterest": 12000,
      "assetName": "Sepeda"
    },
    {
      "contractNumber": "CTR-002",
      "otr": 1000000,
      "amountInterest": 12000,
      "assetName": "Mesin cuci"
    }
  ]
}'

```

# Summary

-    **First, register a source and get a token.** 
-    **Then, use that token to register a customer.**
-    **Set credit limits.**
-    **Finally, submit a loan request.**
-    **Always check `/api/openapi` for the latest API details!**

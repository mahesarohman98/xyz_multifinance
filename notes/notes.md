# Source / Client Registrations

Because customer can take a loan in various platform (ecommerce, web pt xyz, 
or manualy from a dealer), security becomes critical.

for security reasons, consider implement security models likes OAuth.

Objectives:
1. Ensured only registered client can call the API.
2. Ensured customer being **consent** to the transactions.
3. Prevent client spoofing or impersonating customer.
4. Log and trace everything.

# Bounded Context

1. Customer Context (Core)
    Handle personal customer data.

    Entities/Aggregate:
    - Customer

    Activities:
    - Register new customer
    - Update KYC Data
    - Update KTP and selfie photo

2. Credit Limit context (Core)
    Handle tenor based credit limit per customer.

    Entities:
    - TenorLimit
        - Includes 1M, 2M, 3M, 4M credit limits

    Activities:
    - Check available limits
    - Decrease limit on a new credit
    - Validate remaining limit on a request

3. Transaction Context (Core)
    Handle loan request and contracts.

    Entities:
    - Transaction
        - Inclues: contractnumber, asetName, tenor, amount

    Activities:
    - Submit loan application
    - Associate tx with source
    - Record loan
    
4. Source context (Supporting)
    Handling integration third party (ecommerce, web, dealers).

    Entities:
    - Source

    Activities:
    - Register new Source
    - Authenticate source
    - Chcek consent

5. Consent context (Generic)
    Handle consent / authorization flow for end users.

    Entities:
    - ConsentRequest
    - ConsentLog

    Activities:
    - Send OTP/email for consent
    - Log approval or rejection
    - Prevent spoofed transactions

6. Authentcate Context (Generic)
    Handling security and access.

    Entities:
    - ClientApp (ecommerce, dealer)
    - Token

    Activities:
    - Authenticate source
    - Validate JWT / API key
    - Register new client app


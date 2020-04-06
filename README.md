# pay-open-banking-demo
Playing around with creating an open banking demo in Go


## Building
```go build -o pay-open-banking-demo .```

## Running

```docker run -d -p 5432:5432 --name openbanking -e POSTGRES_PASSWORD=mysupersecretpassword postgres```

```APPLICATION_URL=http://localhost:8080/ DATABASE_URL=postgres://postgres:mysupersecretpassword@localhost:5432/?sslmode=disable TRUELAYER_PAY_URL=https://pay-api.truelayer-sandbox.com/ TRUELAYER_AUTH_URL=https://auth.truelayer-sandbox.com/ TRUELAYER_CLIENT_ID=xxx TRUELAYER_CLIENT_SECRET=xxx ./pay-open-banking-demo```

## Environment variables

| variable                | description                                 |
|:------------------------|:--------------------------------------------|
| APPLICATION_URL         | the base URL of the application             |
| DATABASE_URL            | the connection URL for the database         |
| TRUELAYER_PAY_URL       | the URL of the TrueLayer Pay API            |
| TRUELAYER_AUTH_URL      | the URL of the TrueLayer Auth API           |
| TRUELAYER_CLIENT_ID     | the client_id for the TrueLayer account     |
| TRUELAYER_CLIENT_SECRET | the client_secret for the TrueLayer account |

## Usage

Make a POST request to `/v1/api/payments` with body:

```
{
  "reference": "CAKE",
  "description": "Pay for a GOV.UK Pay cake",
  "amount": 2499,
  "return_url": "https://www.google.com"
}
```

This mimics the behaviour of how a government service would integrate with the GOV.UK Pay API.

Follow the `next_url` in the response to complete the payment journey.

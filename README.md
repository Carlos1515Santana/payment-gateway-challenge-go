# Instructions for candidates

This is the Go version of the Payment Gateway challenge. If you haven't already read the [README.md](https://github.com/cko-recruitment/) in the root of this organisation, please do so now. 

## Template structure
```
main.go - a skeleton Payment Gateway API
imposters/ - contains the bank simulator configuration. Don't change this
docs/docs.go - Generated file by Swaggo
.editorconfig - don't change this. It ensures a consistent set of rules for submissions when reformatting code
docker-compose.yml - configures the bank simulator
.goreleaser.yml - Goreleaser configuration
```

Feel free to change the structure of the solution, use a different test library etc.

### Swagger
This template uses Swaggo to autodocument the API and create a Swagger spec. The Swagger UI is available at http://localhost:8090/swagger/index.html .


### Curls Postman

#### Process Payment

```
curl --location 'localhost:8090/api/payments' \
--header 'Content-Type: application/json' \
--data '{
    "card_number": "5275909408908040",
    "expiry_month": 1,
    "expiry_year": 2026,
    "currency": "USD",
    "amount": 10,
    "cvv": "084"
}'
```

#### Retriave Payment

```
curl --location 'localhost:8090/api/payments/c0b0110d-2cb6-45db-92f8-bca15cab45de' \
--header 'Content-Type: application/json'
```
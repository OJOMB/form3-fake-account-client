curl --location --request POST 'http://0.0.0.0/v1/organisation/accounts' \
--header 'Authorization: {{authorization}}' \
--header 'Content-Type: application/json' \
--header 'Date: {{request_date}}' \
--header 'Digest: {{request_signing_digest}}' \
--data-raw '{
  "data": {
    "id": "{{$guid}}",
    "organisation_id": "{{$guid}}",
    "type": "accounts",
    "attributes": {
       "country": "GB",
        "base_currency": "GBP",
        "bank_id": "400302",
        "bank_id_code": "GBDSC",
        "account_number": "10000004",
        "customer_id": "234",
        "iban": "GB28NWBK40030212764204",
        "bic": "NWBKGB42",
        "account_classification": "Personal"
    }
  }
}'
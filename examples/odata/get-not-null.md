# Get Null Data

## GET /V2/Northwind/Northwind.svc/Customers/?$format=json&$filter=ContactName eq 'Frédérique Citeaux'

## Expectation

| Assert | Expected |
| - | - |
| Data.d.results[0].__metadata | *should not be null* |
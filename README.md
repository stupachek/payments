# Payments

Project "Payments"  is [final task](https://drive.google.com/file/d/1X7sKViRpL8t3XOlBElkG7pQvzpR0PjHn/view?usp=sharing) for EPAM Golang course.

## Build and run 
```
docker-compose up --build
```
## Usage
The HTTP server runs on localhost:8080

## Endpoints 
>Each method expects an body with JSON value.
#### users
```
POST /users/register
```
reqiures *first_name*, *last_name*, unique *email*, *password*;
returns user's uuid;
```
POST /users/login
```
reqiures *email*, *password*;
returns authorization token;

>Other endpoints require the Authorization header with the token inside.

#### accounts
```
POST /users/{user_uuid}/accounts/new
```
creates new account for user;
returns account's uuid;
```
GET /users/{user_uuid}/accounts
```
returns accounts for user; 
> URL could contain such query parameters as *offset*, *limit*, *sort_by*(expects *uuid*, *iban* or *balance*), *order*(expects *asc* or *desc*)
```
GET /users/{user_uuid}/accounts/{accounts_uuid}
```
returns the account; 
```
POST /users/{user_uuid}/accounts/{accounts_uuid}/add-money
```
requires *amount*;
adds amount to account's balance;
returns the account; 
#### transaction
```
POST /users/{user_uuid}/accounts/{accounts_uuid}/transactions/new
```
requires *destination_uuid*, *amount*;
creates new transaction with status "prepared";
returns transaction;
```
POST /users/{user_uuid}/accounts/{accounts_uuid}/transactions/{transaction_uuid}/send
```
sends the transaction, updates accounts' balances and transaction status ("sent");
returns transaction;
```
GET /users/{user_uuid}/accounts/{accounts_uuid}/transactions
```
returns transactions; 
> URL could contain such query parameters as *offset*, *limit*, *sort_by*(expects *uuid*, *created_at* or *updated_at*), *order*(expects *asc* or *desc*)

### Request and Response Examples
##### req
```
POST /users/register
```
Body 
```
{   
    "firstName": "Bob",
    "lastName": "Fox",
    "email": "bob.fffox1987@gmail.com",
    "password":"qwerty"
}
```
##### res
Body
```
{
    "message": "registration success",
    "uuid": "b77499e2-ed74-4214-9fd0-86be3456843b"
}
```



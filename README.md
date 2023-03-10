# Payments

Project "Payments"  is [final task](https://drive.google.com/file/d/1X7sKViRpL8t3XOlBElkG7pQvzpR0PjHn/view?usp=sharing) for EPAM Golang course.

## Build and run 
```
docker-compose up --build
```
## Usage
The HTTP server runs on localhost:8080

## Endpoints 
Each method expects an body with JSON value.
### USERS

#### POST `/users/register`

reqiures *first_name*, *last_name*, unique *email*, *password*;
returns user's uuid;
##### example req

`POST http://localhost:8080/users/register`

Body 
```json
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

#### POST `/users/login`

reqiures *email*, *password*;
returns authorization token;
##### example req

`POST http://localhost:8080/users/login`

Body 
```json
{  
    "email": "bob.fffox1987@gmail.com",
    "password":"qwerty"
}
```
##### res

Body
```json
{
    "token": "b89d4729f81b23447a4a93ee01a47950431128031666a1ba368c7cf1512aed85",
    "uuid": "b77499e2-ed74-4214-9fd0-86be3456843b"
}
```
>Other endpoints require the Authorization header with the token inside.

Header
```
    key:Authorization
    value:b89d4729f81b23447a4a93ee01a47950431128031666a1ba368c7cf1512aed85
```

### ADMIN

#### POST `http://localhost:8080/admin/:user_uuid/users/:tagret_uuid/block`

blocks user

##### example req

`POST http://localhost:8080/admin/54149754-cf48-4c13-a949-4d67139f5110/users/d40f82da-0000-4363-bc4c-18c9eabff802/block`

##### res

Body
```json
{
    "message": "user is blocked"
}
```

#### POST `http://localhost:8080/admin/:user_uuid/users/:tagret_uuid/unblock`

unblocks user

##### example req

`POST http://localhost:8080/admin/54149754-cf48-4c13-a949-4d67139f5110/users/d40f82da-0000-4363-bc4c-18c9eabff802/unblock`

##### res

Body
```json
{
    "message": "user is active"
}
```

#### GET `http://localhost:8080/admin/:user_uuid/accounts/requested`

returns all users with *requested-unblock* status

##### example req

`POST http://localhost:8080/admin/54149754-cf48-4c13-a949-4d67139f5110/accounts/requested`

##### res

Body
```json
{
    "accounts": [
        {
            "uuid": "3c82a29a-467f-436d-a3eb-68809fa8f560",
            "iban": "27794c8b5eb7d73256e9d7a1f74d018e76393bf4f2b8358e902ea5e312",
            "balance": 0,
            "user_uuid": "1ed23cb8-ff5b-4634-88b1-72f43f89f369",
            "status": "requested-unblock"
        }
    ]
}
```

#### POST `http://localhost:8080/admin/:user_uuid/accounts/:accounts_uuid/unblock`

unblocks account

##### example req

`POST http://localhost:8080/admin/54149754-cf48-4c13-a949-4d67139f5110/accounts/3c82a29a-467f-436d-a3eb-68809fa8f560/unblock`

##### res

Body
```json
{
    "message": "account is unblocked"
}
```

#### POST `http://localhost:8080/admin/:user_uuid/accounts/:accounts_uuid/unblock`

unblocks user

##### example req

`POST http://localhost:8080/admin/54149754-cf48-4c13-a949-4d67139f5110/accounts/3c82a29a-467f-436d-a3eb-68809fa8f560/unblock`

##### res

Body
```json
{
    "message": "account is unblocked"
}
```

#### POST `http://localhost:8080/admin/:user_uuid/update-role`

changes users role;
reqiures *user_uuid*, *role*;

##### example req

`POST http://localhost:8080/admin/54149754-cf48-4c13-a949-4d67139f5110/update-role`

Body
```json
{  
    "user_uuid": "1ed23cb8-ff5b-4634-88b1-72f43f89f369",
    "role" : "admin"
}
```

##### res

Body
```json
{
    "message": "change role"
}
```

### ACCOUNTS

#### POST `/users/{user_uuid}/accounts/new`
creates new account for user;
returns account's uuid;
##### example req

`POST http://localhost:8080/users/b77499e2-ed74-4214-9fd0-86be3456843b/accounts/new`

##### res

Body
```json
{
    "message": "new account add",
    "uuid": "b89f5687-9cd6-4275-b3d2-87fd7bd8d011"
}
```

#### GET `/users/{user_uuid}/accounts`

returns accounts for user; 
> URL could contain such query parameters as *offset*, *limit*, *sort_by*(expects *uuid*, *iban* or *balance*), *order*(expects *asc* or *desc*)
##### example req

`GET http://localhost:8080/users/b77499e2-ed74-4214-9fd0-86be3456843b/accounts?sort_by=uuid&order=desc&limit=3`

##### res

Body
```json
{
    "accounts": [
        {
            "uuid": "fbe8bee3-1cb7-4d90-8388-105297522a86",
            "iban": "45d0b56c4ceee91e46b64c063d7249ae867beaa2faa9bc5c965429fc88",
            "balance": 0,
            "user_uuid": "b77499e2-ed74-4214-9fd0-86be3456843b",
            "status": "active"
        },
        {
            "uuid": "f1b1dee4-a176-4cec-836f-8a4aa407efbb",
            "iban": "981024baf8ef8361b285d7e9dd86b12cc88f5143c1d209e3c9d07d8354",
            "balance": 0,
            "user_uuid": "b77499e2-ed74-4214-9fd0-86be3456843b",
            "status": "active"
        },
        {
            "uuid": "db689093-81ca-4092-bdc2-52988d5ea970",
            "iban": "1dbfc0e2df7c3edc2ea3118f0d824ecddf29cc95452b0739b05db53d3c",
            "balance": 0,
            "user_uuid": "b77499e2-ed74-4214-9fd0-86be3456843b",
            "status": "active"
        }
    ]
}
```

#### GET `/users/{user_uuid}/accounts/{accounts_uuid}`

returns the account; 
##### example req

`GET http://localhost:8080/users/b77499e2-ed74-4214-9fd0-86be3456843b/accounts/db689093-81ca-4092-bdc2-52988d5ea970`

##### res

Body
```json
{
    "balance": 0,
    "iban": "1dbfc0e2df7c3edc2ea3118f0d824ecddf29cc95452b0739b05db53d3c",
    "uuid": "db689093-81ca-4092-bdc2-52988d5ea970"
}
```

#### POST `/users/{user_uuid}/accounts/{accounts_uuid}/add-money`

requires *amount*;
adds amount to account's balance;
returns the account; 
##### example req

`POST http://localhost:8080/users/b77499e2-ed74-4214-9fd0-86be3456843b/accounts/fbe8bee3-1cb7-4d90-8388-105297522a86/add-money`

``` json
{  
    "amount" : "123"
}
```

##### res
Body
``` json
{
    "account": {
        "uuid": "fbe8bee3-1cb7-4d90-8388-105297522a86",
        "iban": "45d0b56c4ceee91e46b64c063d7249ae867beaa2faa9bc5c965429fc88",
        "balance": 123,
        "user_uuid": "b77499e2-ed74-4214-9fd0-86be3456843b",
        "status": "active"
    },
    "message": "add money"
}
```

#### POST `/users/{user_uuid}/accounts/{accounts_uuid}/block`

block account

##### example req

`POST http://localhost:8080/users/1ed23cb8-ff5b-4634-88b1-72f43f89f369/accounts/3c82a29a-467f-436d-a3eb-68809fa8f560/block`

##### res
Body
``` json
{
    "message": "account is blocked"
}
```

#### POST `/users/{user_uuid}/accounts/{accounts_uuid}/unblock`

block account

##### example req

`POST http://localhost:8080/users/1ed23cb8-ff5b-4634-88b1-72f43f89f369/accounts/3c82a29a-467f-436d-a3eb-68809fa8f560/unblock`

##### res
Body
``` json
{
    "message": "account is waiting to be unblock"
}
```
### TRANSACTION

#### POST `/users/{user_uuid}/accounts/{accounts_uuid}/transactions/new`

requires *destination_uuid*, *amount*;
creates new transaction with status "prepared";
returns transaction;
##### example req

`POST http://localhost:8080/users/b77499e2-ed74-4214-9fd0-86be3456843b/accounts/fbe8bee3-1cb7-4d90-8388-105297522a86/transactions/new`

```json
{  
    "destination_uuid" : "db689093-81ca-4092-bdc2-52988d5ea970",
    "amount": "30"
}
```

##### res

Body
``` json
{
    "message": "create new transaction",
    "transaction": {
        "uuid": "d8882d3c-2d44-4312-ac10-020f45ea4c43",
        "status": "prepared",
        "source_uuid": "fbe8bee3-1cb7-4d90-8388-105297522a86",
        "destination_uuid": "db689093-81ca-4092-bdc2-52988d5ea970",
        "amount": 30,
        "created_at": "2023-02-20T09:20:48.565437Z",
        "updated_at": "2023-02-20T09:20:48.565437Z"
    }
}
```

#### POST `/users/{user_uuid}/accounts/{accounts_uuid}/transactions/{transaction_uuid}/send`

sends the transaction, updates accounts' balances and transaction status ("sent");
returns transaction;
##### example req

`POST http://localhost:8080/users/b77499e2-ed74-4214-9fd0-86be3456843b/accounts/fbe8bee3-1cb7-4d90-8388-105297522a86/transactions/d8882d3c-2d44-4312-ac10-020f45ea4c43/send`

##### res

Body
```json
{
    "message": "sent transaction",
    "transaction": {
        "uuid": "d8882d3c-2d44-4312-ac10-020f45ea4c43",
        "status": "sent",
        "source_uuid": "fbe8bee3-1cb7-4d90-8388-105297522a86",
        "destination_uuid": "db689093-81ca-4092-bdc2-52988d5ea970",
        "amount": 30,
        "created_at": "2023-02-20T09:20:48.565437Z",
        "updated_at": "2023-02-20T09:22:15.522694Z"
    }
}
```

#### GET `/users/{user_uuid}/accounts/{accounts_uuid}/transactions`

returns transactions; 
> URL could contain such query parameters as *offset*, *limit*, *sort_by*(expects *uuid*, *created_at* or *updated_at*), *order*(expects *asc* or *desc*)
##### example req

`GET http://localhost:8080/users/b77499e2-ed74-4214-9fd0-86be3456843b/accounts/fbe8bee3-1cb7-4d90-8388-105297522a86/transactions?sort_by=created_at`

##### res

Body
```json
{
    "transactions": [
        {
            "uuid": "d8882d3c-2d44-4312-ac10-020f45ea4c43",
            "status": "sent",
            "source_uuid": "fbe8bee3-1cb7-4d90-8388-105297522a86",
            "destination_uuid": "db689093-81ca-4092-bdc2-52988d5ea970",
            "amount": 30,
            "created_at": "2023-02-20T09:20:48.565437Z",
            "updated_at": "2023-02-20T09:22:15.522694Z"
        },
        {
            "uuid": "fedfbd72-8daf-4a05-8684-3ac07789f0af",
            "status": "prepared",
            "source_uuid": "fbe8bee3-1cb7-4d90-8388-105297522a86",
            "destination_uuid": "db689093-81ca-4092-bdc2-52988d5ea970",
            "amount": 20,
            "created_at": "2023-02-20T09:21:40.51269Z",
            "updated_at": "2023-02-20T09:21:40.51269Z"
        },
        {
            "uuid": "a90ae2a4-eea0-42e5-b9fd-bb1ed47daaac",
            "status": "prepared",
            "source_uuid": "fbe8bee3-1cb7-4d90-8388-105297522a86",
            "destination_uuid": "db689093-81ca-4092-bdc2-52988d5ea970",
            "amount": 15,
            "created_at": "2023-02-20T09:21:49.047032Z",
            "updated_at": "2023-02-20T09:21:49.047032Z"
        }
    ]
}
```
## ERD
![ERD](payment.png)

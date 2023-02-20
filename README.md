# Payments

Project "Payments"  is [final task](https://drive.google.com/file/d/1X7sKViRpL8t3XOlBElkG7pQvzpR0PjHn/view?usp=sharing) for EPAM Golang course.

## Build and run 
```
docker-compose up --build
```
## Usage
The HTTP server runs on localhost:8080

##Endpoints 
####users
```
* POST /users/register
```
reqiures first name, last nane, uniqe email, password;
returns users uuid;
```
* POST /users/login
```
reqiures email, password;
returns authorization token;

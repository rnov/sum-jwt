curl --location --request POST 'http://localhost:8080/auth' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "validUser",
    "password": "qwerty"
}'

curl --location --request POST 'http://localhost:8080/sum' \
--header 'Authorization: Basic eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MTY0MzU5MTUsInVzZXJfaWQiOiJ2YWxpZFVzZXIifQ.E0DwdKjYVm-STrS-x6r_gtH3Wq1kAZkA-GX9G67wQyc' \
--header 'Content-Type: application/json' \
--data-raw '[
    1,
    2,
    3,
    4
]'

curl --location --request POST 'http://localhost:8080/sum' \
--header 'Authorization: Basic eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MTY0MzU5MTUsInVzZXJfaWQiOiJ2YWxpZFVzZXIifQ.E0DwdKjYVm-STrS-x6r_gtH3Wq1kAZkA-GX9G67wQyc' \
--header 'Content-Type: application/json' \
--data-raw '[]'

curl --location --request POST 'http://localhost:8080/sum' \
--header 'Authorization: Basic eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MTY0MzU5MTUsInVzZXJfaWQiOiJ2YWxpZFVzZXIifQ.E0DwdKjYVm-STrS-x6r_gtH3Wq1kAZkA-GX9G67wQyc' \
--header 'Content-Type: application/json' \
--data-raw '{
    "a": [
        -1,
        1,
        "dark"
    ]
}'

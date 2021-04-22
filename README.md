# Go Chi REST API - example project

Simple Go REST API using the _chi_ and MongoDB frameworks.

## Endpoints

| Endpoint       | HTTP Method | Description      | Example              |
| -------------- | ----------- | ---------------- | -------------------- |
| /              | GET         | Get all contacts | /contacts            |
| /{phoneNumber} | GET         | Get a contact    | /contacts/0147344454 |
| /              | POST        | Create a contact | /contacts            |
| /{phoneNumber} | PUT         | Update a contact | /contacts/0147344454 |
| /{phoneNumber} | DELETE      | Delete a contact | /contacts/0147344454 |

## Requisites

- MongoDB, the code uses the 27017 port, please make sure MongoDB is up and running.

## Playground / cURL script

There are some [scripts](playground/api) in the playground/api directory that exemplifies how to use the api. Example:

```bash
curl --header "accept: */*"                \
    --header  "Content-Type: application/json"  \
    --request GET "http://localhost:3060/contacts"
```

## Donation / Sponsorship ❤️ 👍

This code was brought to you by [Leo Gutiérrez](https://github.com/leogtzr) in his free time. If you want to thank me and support the development of this project, please make a small donation on [PayPal](https://www.paypal.me/leogtzr). In case you also like my other open source contributions and articles, please consider motivating me by becoming a sponsor/patron on [Patreon](https://www.patreon.com/leogtzr). Thank you! ❤️

#!/bin/bash

curl --header "accept: */*"                \
    --data '@new_contact.json' \
    --header  "Content-Type: application/json"  \
    --request POST "http://localhost:3060/contacts"

exit

#!/bin/bash

curl --header "accept: */*"                \
    --header  "Content-Type: application/json"  \
    --request DELETE "http://localhost:3060/contacts/works"

exit
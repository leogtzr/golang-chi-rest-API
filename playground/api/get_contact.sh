#!/bin/bash

curl --header "accept: */*"                \
    --header  "Content-Type: application/json"  \
    --request GET "http://localhost:3060/contacts/614124233422334"

exit
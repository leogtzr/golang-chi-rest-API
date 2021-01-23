#!/bin/bash

curl --header "accept: */*"                \
	--data '@update_contact.json' \
    --header  "Content-Type: application/json"  \
    --request PUT "http://localhost:3060/contacts/614124233422334"

exit
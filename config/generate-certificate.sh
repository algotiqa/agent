#!/bin/bash

#--- Generate agent key & cert request
openssl req -newkey rsa:2048 -nodes -keyout agent.key -out agent.csr -subj "/C=EU/ST=Italy/L=Rome/O=Algotiqa/OU=AlgotiqaAgent/CN=*"

#--- Self sign agent certificate
openssl x509 -req -days 9999 -sha256 -in agent.csr -out agent.crt \
    -extfile <(echo subjectAltName=IP:158.69.1.2)

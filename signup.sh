#!/usr/bin/env bash

PSS_USER_ID=$(curl -XPOST "https://hub-testnet.lightningpeach.com/pss-walleto/api/v1/signup" | jq -r .content)
echo $PSS_USER_ID

curl -XPOST "https://hub-testnet.lightningpeach.com/api/v2/auth/signup" -d '{
  "email": "test.bitf1000042@gmail.com",
  "password": "1234qQ1234"
}' -H "X-Pss-User-Id: $PSS_USER_ID."

curl -XPOST "https://hub-testnet.lightningpeach.com/api/v2/auth/signin" -d '{
  "email": "test.bitf1000042@gmail.com",
  "password": "1234qQ1234"
}' | jq .
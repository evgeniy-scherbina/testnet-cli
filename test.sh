#!/usr/bin/env bash

TOKEN=eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZXJjaGFudF9pZCI6ImlKN1lLWmRuVHpVc215Smo0QWFpeU9vMiIsIndoaXRlbGFiZWwiOiJkZWZhdWx0IiwiZ3JvdXAiOiJjdXN0b2RpYWwiLCJleHAiOjE1NzQ3OTgyMTh9.oil1lZOnr74qeabJ64qoH5QjCwVp67hp8vcFNmXlPEF8L2O2VWkIvctdD8rtov6_xX5zgD5Kctw5GQ1I0Q4oHQ
curl -X POST "https://hub-testnet.lightningpeach.com/api/v2/payment/lightning" \
    -H "accept: application/json" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" -d "{ \"merchant_id\": \"liz1IQalAnyR3JXhgIPvbKxU\", \"amount\": \"333\", \"description\": \"string\", \"expiry\": 3600, \"timestamp\": 0, \"ntfn_url\": \"string\"}"
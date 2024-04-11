OPENSSL = openssl

# Generate required private and public keys for the API service
api_keys:
	$(OPENSSL) genpkey -algorithm ed25519 -outform PEM -out ./secrets/api-private.pem
	$(OPENSSL) pkey -in ./secrets/api-private.pem -pubout -out ./secrets/api-public.pem

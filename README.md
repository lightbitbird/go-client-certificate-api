# go-client-certificate-api
Client Certificate API with TLS using Golang


## Generate Client certification AND CA certificate
```
# Generate CA certificate
openssl genrsa -des3 -out ca.key 4096
openssl req -new -x509 -days 365 -key ca.key -out ca.crt

# Generate server certificate
openssl genrsa -des3 -out server.key 1024
openssl req -new -key server.key -out server.csr
openssl x509 -req -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

# Generate client certificate
openssl genrsa -des3 -out client.key 1024
openssl req -new -key client.key -out client.csr
openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out client.crt

openssl rsa -in server.key -out temp.key
rm server.key
mv temp.key server.key

openssl rsa -in client.key -out temp.key
rm client.key
mv temp.key client.key
```

## Run Certificate Server
```
go run server/https-server.go
```

## TLS Communication Test from Client API
```
# PEM type
go run client/https-client.go 
# Pkcs12 type
go run client/https-client.go pfx
```


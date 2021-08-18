#!/usr/bin/env sh
apk add --no-cache openssl

rm export/*.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout /export/ca-key.pem -out /export/ca-cert.pem -subj "/C=SG/ST=SG/L=Singapore/O=Beam Mobility Holdings Pte. Ltd./OU=Transportation/CN=localhost/emailAddress=eng-tools@ridebeam.com"

#echo "CA's self-signed certificate"
#openssl x509 -in /export/ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout /export/srv-key.pem -out /export/srv-req.pem \
  -subj "/C=SG/ST=SG/L=Singapore/O=Beam Mobility Holdings Pte. Ltd./OU=Transportation/CN=localhost/emailAddress=eng-tools@ridebeam.com" \
  # -addext "subjectAltName = IP:127.0.0.1"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in /export/srv-req.pem -days 360 -CA /export/ca-cert.pem -CAkey /export/ca-key.pem -CAcreateserial -out /export/srv-cert.pem
chmod 400 /export/*.pem

echo "Server's signed certificate"
openssl x509 -in /export/srv-cert.pem -noout -text

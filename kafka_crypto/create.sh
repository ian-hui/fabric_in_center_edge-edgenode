#!/bin/sh
#生成一个keystore（存储密钥和证书的storage）并放入没有签名的rsa钥匙和证书
keytool -keystore kafka.client.keystore.jks -alias localhost -validity 3650  -dname "CN=localhost, OU=localhost, O=localhost, L=SH, ST=SH, C=CN" -keyalg RSA -genkey -ext "SAN=IP:0.0.0.0,IP:172.16.1.0/24,DNS:localhost" -storepass reins5401 -keypass reins5401
openssl req -x509 -newkey rsa:2048 -keyout ca-key -out ca-cert -days 3650 -subj "/CN=localhost/"
keytool -keystore kafka.client.truststore.jks -alias CARoot -import -file ca-cert 
keytool -keystore kafka.server.truststore.jks -alias CARoot -import -file ca-cert 

keytool -keystore kafka.client.keystore.jks -alias localhost -certreq -file cert-file -ext "SAN=IP:0.0.0.0,DNS:localhost"
openssl x509 -req -CA ca-cert -CAkey ca-key -in cert-file -out cert-signed -days 3650 -CAcreateserial -passin pass:reins5401  -extfile openssl.cnf -extensions v3_req
keytool -keystore kafka.client.keystore.jks -alias CARoot -import -file ca-cert
keytool -keystore kafka.client.keystore.jks -alias localhost -import -file cert-signed

keytool -importkeystore -srckeystore kafka.server.truststore.jks -destkeystore server.p12 -deststoretype PKCS12
openssl pkcs12 -in server.p12 -nokeys -out server.cer.pem 

keytool -importkeystore -srckeystore kafka.client.keystore.jks -destkeystore client.p12 -deststoretype PKCS12 
openssl pkcs12 -in client.p12 -nokeys -out client.cer.pem 
openssl pkcs12 -in client.p12 -nodes -nocerts -out client.key.pem

cd .. && chmod -R 777 ./kafka_crypto && chown reins:reins ./kafka_crypto
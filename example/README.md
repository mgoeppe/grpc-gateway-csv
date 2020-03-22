# Example

In order to run the example use:

```
$ go run .
grpc on :8080 ..
http on :8081 ..
```

and call the http endpoint:

```
$ curl -H 'Content-Type: text/csv'  http://localhost:8081/v1/example
Col1;Col2;Col3;Col4;Col5
dreggn;42;true;one|two;k1:v1
dreggn;42;true;eins|zwo;k1:v1
```
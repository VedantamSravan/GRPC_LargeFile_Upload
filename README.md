
# Large File Upload GRPC Client and Server

This project demonstrates a GRPC-based client-server application for transferring large files securely using TLS.

---

## Prerequisites

1. **Install Necessary Tools**:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

2. **Environment Setup**:
   - **Server** IP: `10.91.170.144`  
   - **Client** IP: `10.91.170.141`  
   *(You can verify the server is running by checking the log: `Server listening on 0.0.0.0:50051`)*

---

## Publisher and Subscriber Roles

- **Publisher**: Acts as the **client**.  
- **Subscriber**: Acts as the **server**.

---

## Directory Structure

```plaintext
GRPC_File_Upload/
├── uploads/               # Directory to store uploaded files
├── pb/                    # Protocol buffer definitions
│   ├── filetransfer.proto
│   ├── filetransfer.pb.go
│   └── filetransfer_grpc.pb.go
├── server.go              # Server code
├── client.go              # Client code
├── cert-config.cnf        # Configuration file for generating certificates
├── server.key             # Server private key
├── server.csr             # Server CSR (Certificate Signing Request)
├── server.crt             # Server certificate
```

---

## Protocol Buffer Compilation

1. **Navigate to the `pb` directory**:
   ```bash
   cd pb
   ```

2. **Run the following commands to generate necessary files**:
   ```bash
   protoc --go_out=. --go_opt=paths=source_relative filetransfer.proto
   protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative filetransfer.proto
   ```

3. **Verify Generated Files**:
   After running the commands, the following files should appear:
   - `filetransfer.pb.go`
   - `filetransfer_grpc.pb.go`

---

## Generating Self-Signed Certificates

1. **Create a Configuration File for SANs**:
   Create a file named `cert-config.cnf` with the following content:

   ```ini
   [req]
   distinguished_name = req_distinguished_name
   req_extensions = v3_req
   prompt = no

   [req_distinguished_name]
   C = US
   ST = NH
   L = Nashua
   O = Nextcomputing
   OU = Software
   CN = localhost

   [v3_req]
   keyUsage = critical, digitalSignature, keyEncipherment
   extendedKeyUsage = serverAuth, clientAuth
   subjectAltName = @alt_names

   [alt_names]
   DNS.1 = localhost
   IP.1 = 127.0.0.1
   IP.2 = 10.91.170.144  # Replace with your server's IP address
   ```

2. **Generate Certificates**:
   ```bash
   # Generate a new private key
   openssl genrsa -out server.key 2048

   # Generate a new certificate signing request (CSR)
   openssl req -new -key server.key -out server.csr -config cert-config.cnf

   # Generate the self-signed certificate with SANs
   openssl x509 -req -in server.csr -signkey server.key -out server.crt -days 365 -extensions v3_req -extfile cert-config.cnf
   ```

---
Note:-
We need to copy client code on client machine my case it is 141
for server code need to copy on 144 

## Running the Application

We need to copy client code on the client machine (in this case, 141)
and server code on the server machine (in this case, 144).

1. **Start the Server**:
   On `10.91.170.144`:
   ```bash
   go run server.go
   ```

2. **Start the Client**:
   On `10.91.170.141`:
   ```bash
   go run client.go
   ```

---

## Notes

- Ensure that the `server.crt` and `server.key` files are present and correctly configured.
- Ensure firewall settings allow communication on port `50051`.
- Make sure both client and server have network access to each other.

Happy coding!

# Splitpay

A way to easily split your bills. Built with go-starter.

![Screenshot 2024-10-15 at 10 02 39 PM](https://github.com/user-attachments/assets/74599ea1-cb30-4dc7-afcf-ddb562a7406a)

## Features
- Scan a receipt, OpenAI will extract the items and prices
- Edit receipts to add/remove items or change prices
- Get a shareable link to the receipt
- Pay with Stripe! (And my poor wallet for Stripe fees :sob:)
- Split the bill with your friends!
- Paid items are marked

## Installation

Run make commands to setup and start the server.

```
# Installs wgo, templ, creates a local-sqlite.db file and runs up-migrations
make setup          # for Linux
make setup-mac      # for Mac

make dev            # Start a server locally
make dev port=9000  # Start a server at a specific port

make build          # Build the binary & docker image
``` 

.env file is required for the AI portion and remotely connecting to the DB
```
ENVIRONMENT=DEV
TURSO_TOKEN=
TUSRO_DB_URL=
OPENAI_TOKEN=
```

## DB
You can use sqlite3 (natively installed on Mac) to work with the database locally. 
```
sqlite3 local-sqlite.db

sqlite3> SELECT * FROM receipts;
```

Remotely, the DB runs on [Turso](turso.tech) which is also sqlite.

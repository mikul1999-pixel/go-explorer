# go-explorer

A small, lightweight HTTP API for managing files on your disk remotely. Meant to use as a backup access method or as a building block for a custom UI.


## Features

- **list**: list all files and directories in root
- **download**: download a file
- **upload**: upload a file (**Note**: not chunking... yet)
- **mkdir**: create a directory
- **rename**: rename (or move) a file or directory
- **delete**: delete a file or directory

## Security 

- Root directory sandboxing (no path escape)
- Overwrite is forbidden and will return with an error
- Includes pluggable auth middleware (headers or shared secret)
    - Assumes the api will be behind a trusted identity-aware proxy method
    - Not a replacement for IDP authentication
    - Can be enabled or disabled in ```.env```


## Installation

### Setup
```bash
git clone https://github.com/mikul1999-pixel/go-explorer.git
cd go-explorer
```

### Config

Create your ```.env```
```bash
cp .env.example .env
```

Then edit your ```.env``` and ```Makefile``` <br>
If you want to deploy via docker, edit ```Dockerfile``` and ```docker-compose.yml```

### Deploy
Run via docker (recommended)
```bash
docker compose up -d
```

Run locally
```bash
make run
```

The API will be available at: ```http://localhost:3030```


## Usage
All paths defined in the API are relative to ```ROOT_DIR``` (defined in the config files). Path traversal above ```ROOT_DIR``` is restricted

API Endpoints

```bash
# GET /fs/list
curl -X GET "http://localhost:3030/fs/list"
curl -X GET "http://localhost:3030/fs/list?path=dir"

# GET /fs/download
curl -O "http://localhost:3030/fs/downloads?path=dir/download.txt"

# POST /fs/upload
curl -X POST -F "file=@example.txt" "http://localhost:3030/fs/upload?path=dir"

# POST /fs/mkdir
curl -X POST "http://localhost:3030/fs/mkdir?path=newdir"

# POST /fs/rename
curl -X POST "http://localhost:3030/fs/rename?from=dir/example.txt&to=dir/example-renamed.txt"
curl -X POST "http://localhost:3030/fs/rename?from=dir/example.txt&to=newdir/example-renamed.txt"

# POST /fs/delete
curl -X POST "http://localhost:3030/fs/delete?path=dir/example.txt"

# GET /healthz
curl "http://localhost:3030/healthz"

```

If Auth is enabled, include your header in the request ```curl -H <auth> -X <Method> <Endpoint>```

## Appendix
**Note:** This is a personal project. Features and functionality may change.

<br>

### Authentication
Auth middleware is designed to be an **optional** step to confirm that user passed a stronger validation method.

#### Header-based auth (recommended)
Used with identity-aware proxies (Cloudflare Access, Authentik, etc). Only accept requests from trusted headers
```bash
AUTH_MODE=header
AUTH_HEADER=Cf-Access-Authenticated-User-Email
AUTH_ALLOWED=me@example.com,admin@example.com
```

#### Shared secret auth
Used with a strong, random key. Only recommended when deployed on a Private Network
```bash
AUTH_MODE=shared
AUTH_SHARED_HEADER=X-Explorer-Key
AUTH_SHARED_KEY=supersecret_uuid
```

#### Disable auth
No auth enforcement, any request will pass. Only use while temporarily testing on your network
```bash
AUTH_MODE=
```


<br>

---

#!/bin/bash

#
#
#
echo "Test binary compilation"

go build -o p2cli .

if [ $? -ne 0 ]; then
  echo "=> Failed"
  exit 1
fi

#
#
#
echo "Test with YAML file"

EXPECTED=$(cat <<-END
address: localhost:5000
no_registration: true
database_path: ""
secret_key: verystrongsecret-jwt
session:
  secret: verystrongsecret-paseto
  access_token_ttl: 440h
  refresh_token_ttl: 8760h
brokers:
- localhost:6000
- localhost:6001
- localhost:6002
END
)
ACTUAL=$(./p2cli -v example.varsfile.yml example.yml.j2)

if [[ $EXPECTED != $ACTUAL ]]; then
  echo "=> Failed"
  diff <(echo "$EXPECTED" ) <(echo "$ACTUAL")
  exit 1
fi

#
#
#
echo "Test with ENV"

#
#
echo "- Without variables defined"
echo "  - STDOUT"

ACTUAL=$(./p2cli -p STANDARDFILE_ example.yml.j2 2> /dev/null)
if [[ -n $ACTUAL ]]; then
  echo "=> Failed"
  echo "$ACTUAL"
  exit 1
fi

echo "  - STDERR"

ACTUAL=$(./p2cli -p STANDARDFILE_ example.yml.j2 2>&1 > /dev/null)
if [[ $ACTUAL != "Error: no value found for key(s) [address, jwt_secret_key]" ]]; then
  echo "=> Failed"
  echo "$ACTUAL"
  exit 1
fi

#
#
echo "- With all variables defined"

export STANDARDFILE_ADDRESS="localhost:5000"
export STANDARDFILE_NO_REGISTRATION=true
export STANDARDFILE_JWT_SECRET_KEY=verystrongsecret-jwt-env
export STANDARDFILE_SESSION__SECRET_KEY=verystrongsecret-paseto-env
export STANDARDFILE_BROKERS='["localhost:6000","localhost:6001","localhost:6002"]'

EXPECTED=$(cat <<-END
address: localhost:5000
no_registration: true
database_path: ""
secret_key: verystrongsecret-jwt-env
session:
  secret: verystrongsecret-paseto-env
  access_token_ttl: 1440h
  refresh_token_ttl: 8760h
brokers:
- localhost:6000
- localhost:6001
- localhost:6002
END
)
ACTUAL=$(./p2cli -p STANDARDFILE_ example.yml.j2)

if [[ $EXPECTED != $ACTUAL ]]; then
  echo "=> Failed"
  diff <(echo "$EXPECTED" ) <(echo "$ACTUAL")
  exit 1
fi
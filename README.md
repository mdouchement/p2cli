# p2cli

p2cli is a portable Jinja2 (Ansible) template engine.

## Usage

- From environement variables:

```sh
export STANDARDFILE_ADDRESS="localhost:5000"
export STANDARDFILE_NO_REGISTRATION=true
export STANDARDFILE_JWT_SECRET_KEY=verystrongsecret-jwt-env
export STANDARDFILE_SESSION__SECRET_KEY=verystrongsecret-paseto-env # This key converted to 'session.secret_key'
export STANDARDFILE_BROKERS='["localhost:6000","localhost:6001","localhost:6002"]' # JSON format

$ p2cli -p 'STANDARDFILE_' example.yml.j2
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
```

- From a YAML vault file:

```sh
$ p2cli -v example.varsfile.yml example.yml.j2
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
```

## License

**MIT**


## Contributing

All PRs are welcome.

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request
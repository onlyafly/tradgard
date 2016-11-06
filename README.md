# tradgard

## Development

### Libraries Used

https://github.com/labstack/echo

### Dependency Management

See if anything is missing

    govendor list

Add a new dependency

    govendor fetch <package>

Updating an existing dependency

    govendor fetch <package>@<version>

Removing unused dependencies

    govendor remove +u

Fetch all referenced dependencies

    govendor fetch +e


### Troubleshooting

#### Error starting server listen tcp :5000: bind: address already in use

Kill the process

    make clean

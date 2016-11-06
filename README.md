# tradgard

## Development

### Libraries Used

https://github.com/labstack/echo

### Dependency Management

Add a new dependency

    govendor fetch <package>

Updating an existing dependency

    govendor fetch <package>@<version>

Removing unused dependencies

    govendor remove +u

Fetch all referenced dependencies

    govendor fetch +e

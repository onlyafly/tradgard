# Tradgard: a garden for ideas

## About

**Tradgard** is a modern take a wiki. It is still pretty young.

You can read the motivation behind Tradgard at [Tradgard's Motivation](http://garden.kevinalbrecht.com/Tradgard%27s+Motivation).

You can contact the author at his website, [KevinAlbrecht.com](http://www.kevinalbrecht.com) or on Twitter [@kevinpalbrecht](https://twitter.com/kevinpalbrecht).

## Installation on macOS

### Installation Step #1: Install Go

Install Go on your machine:

    brew install go

Set your GOPATH. In your ~/.profile file, add the following lines:

    # This is the directory where the Go tools will search for source code:
    export GOPATH=$HOME

    # This is the directory where Homebrew installs Go:
    export GOROOT=/usr/local/opt/go/libexec/

    export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

You'll need to restart your terminal for the changes to take effect.

Clone the Tradgard source code. Note that your should clone it into a subdirectory of the ~/src directory:

    mkdir ~/src
    cd ~/src
    git clone https://github.com/onlyafly/tradgard.git

### Installation Step #2: Install Postgres

Install Postgres on your machine.

### Installation Step #3: Setup Tradgard

Create a Postgres database called "tradgard".

Make a copy of the file "example.env" and name the copy ".env".

Change these variables in the ".env" file to custom values:

    ADMIN_USERNAME=foo
    ADMIN_PASSWORD=bar
    COOKIE_32_BYTE_HASH_KEY=This key should be 32 bytes long
    COOKIE_32_BYTE_BLOCK_KEY=This key should be 32 bytes long

Install all dependencies:

    make deps

### Installation Step #4: Deploy to Heroku

*TODO*

## Running the Webserver

Start the server (the server will restart automatically on file changes):

    make serve

Visit the site:

    http://localhost:5000/

## Development

### Important libraries used

* Echo (web framework): https://github.com/labstack/echo

### Dependency management

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

## Troubleshooting

#### Error starting server listen tcp :5000: bind: address already in use

Kill the process

    make clean

## License

The code is available under the [MIT](https://github.com/onlyafly/tradgard/blob/master/LICENSE) license.

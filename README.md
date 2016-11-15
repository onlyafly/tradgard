# Tradgard: a garden for ideas

**Author:** Kevin Albrecht (http://www.kevinalbrecht.com)

Tradgard is a modern take a wiki. It is still pretty young.

## Initial setup

### Initial setup on your local machine

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

Install Postgres on your machine.

Create a Postgres database called "tradgard".

Make a copy of the file "example.env" and name the copy ".env".

Install all depedencies:

    make deps

You are ready to go!

### Initial setup on Heroku

(work in progress)

## Running on your local machine

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

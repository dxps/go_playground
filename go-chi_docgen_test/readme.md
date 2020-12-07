## go-chi/docgen Test

Main purpose of this test is to make sure go-chi's docgen is working fine.

Initially I got the "ERROR: docgen: unable to determine your $GOPATH" error.

I set the `GOPATH` (meaning adding `export GOPATH=$HOME/bin` to my `~/.profile` that is sourced in `~/.zshrc`) and routes are generated just fine.

## Ultimate Service 2.0

Playground and learning space while studying ArdanLabs' Ultimate Service 2.0.<br/>
I simply love Bill Kennedy's style of deliverying the material.

### Workflows

#### Importing Packages

This setup is using the _vendoring_ approach and thus these steps are used:
- add the package import statement in the code
- run `make tidy` (which finds the package in the Internet, selects the version, add it to our `go.mod` file and adds it to the local `vendor` directory)
- continue with the code, everything should be fine (no IDE / gopls complains)

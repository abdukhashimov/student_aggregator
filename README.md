## Student Aggregator

This repository contains the source for the Student Aggregator app

## How to run the application
```
go run main.go --help # to see the avaialble commands
go run main.go grpc --help # to see the available flags for grpc command
go run main.go http --help # to see the available flags for http

# example to run http server
go run main.go http --port=8080
```


### Contribution

Please make sure you follow project code style rules. To prevent code issues please use pre-commit hooks.
Before you can run hooks, you need to have the pre-commit package manager installed.


Using pip:
```sh
pip install pre-commit
```

Using homebrew:
```sh
brew install pre-commit
```

run pre-commit install to set up the git hook scripts:
```sh
pre-commit install
```

run pre-commit install --hook-type commit-msg to set up the git commit-msg hook scripts:
```sh
pre-commit install --hook-type commit-msg
```

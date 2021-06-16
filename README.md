## BubuLearn services core

Common code for the Bubulearn go services.

### Installation

```shell
git clone https://gitlab.com/bubulearn/bubucore.git
cd bubucore
./make
```

Also, you can use `./make build_fast` to skip
some tests & linters while build.

### Usage 

To use this library in your project, do the next steps:

1. Add git rule to use ssh instead of https for the gitlab:
```shell
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

2. Set `GOPRIVATE` var:
```shell
go env -w "GOPRIVATE=gitlab.com/bubulearn/*"
```

3. `go get` the library:
```shell
go get gitlab.com/bubulearn/bubucore 
```

### Contributing

1. Clone repository
1. Create your custom branch from the `master`
1. Do some changes
1. Commit and push
1. Create merge request
1. Get approval and merge
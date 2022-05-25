# Using CUE with Network Devices

We are configuring a Cumulus Linux (`cvx`) devices, that is built as a part of this book's lab topology. Make sure that topology is up and running before trying any of the instructions.

## Option 1 -- Using CUE Go API

This is the option described in the book. The following command calls into CUE APIs to evaluate the CUE files, produce the JSON payload and apply it to the running configuration of the `cvx` device:

```bash
ch08/cue$ go run main.go
main.go:147: Created revisionID: changeset/cumulus/2022-05-22_14.42.26_94J3
{
  "state": "apply",
  "transition": {
    "issue": {},
    "progress": ""
  }
}
main.go:76: Successfully configured the device
```

## Option 2 -- Using CUE tool

CUE tool is a command line tool, similar to `go` that can be used to interact with CUE code. It is the primary user-facing interface and can be used to implement the same device configuration workflow as in Option#1.

First, import the `input.yaml` file into the `input` CUE package:

```bash
cue import input.yaml -p input -f
```

You can preview the generated payload using the following command:

```bash
cue eval network.automation:cvx -c
```

The entire `cvx` configuration workflow, including creating a revision and applying it to the running configuration can be encoded in a special CUE tool file, that extends the CUE CLI by defining a number of commands. In this case, we've defined the `apply` command and invoke it by running:

```
cue cmd apply cue_tool.cue
```
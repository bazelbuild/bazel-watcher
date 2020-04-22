# Changes made to `third_party` and their reasons

## Directory laout

Go's native toolchain doesn't support multiple packages being defined in the
same directory. The proto files that we are working with all come from the
exact same directory but declare different packages. This is an inherent
conflict.

In order to mitigate this, the package `analysis` was put in a directory called
`analysis`, while the package titled `blaze_query` was put into one called
`blaze_query`.

In order to accommodate this move, import statements were also adjusted.

## Inclusion of compiled sources

As a number of users compile ibazel directly from source, and to better comply
with the expectations of the go community, I include the compiled `.pb.gw`
files. This allows the `go build` command to be used to build ibazel.

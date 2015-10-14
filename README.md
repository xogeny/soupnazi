# No Soup for You!

## License Management with JWTs?

I had some code I wanted to be able to license so I looked around for
some open source alternatives to the FlexLM type of approach.  I
didn't really find anything I liked (let me know if I missed anything).

Then I got to thinking about how to implement such an approach in a
relatively modern way.  The first thing I thought of was using JWTs.
The idea is basically to use JWTs in much the same way that they are
used for authorization purposes in APIs, *i.e.,* you encode some kind
of grant in them and then they are signed using some shared secret.

The interesting thing about a JWT is that it can actually encode
different types of information and grants.

## Example

If you have an application named `"myAppName"`, you can check whether
a given feature, `"someFeature"` is available to the current user by
simply instantiating a license manager instance and asking about the
feature in question, *e.g.,*

```
	lm := soupnazi.NewLM("myAppName", "sharedSecret")
	_, err := lm.License("someFeature")
	if err != nil {
		log.Printf("License error: %v", err)
		os.Exit(1)
	}
```

Note that the `"sharedSecret"` string is used during license creation
and must be securely embedded in the application as well.

## Generating Licenses

The `genjwt` command is used to generate licenses.  You can run
`genjwt --help` for a complete description of command line options.  A
very common scenario is to generate a license for a specific MAC
address that will last a certain number of days.  In this scenario,
the following command can be used:

```
$ genjwt "myAppName" "someFeature" "sharedSecret" --mac "5c:f9:38:89:6b:f2" --days 30
```

This will output the string:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhIjoiNWM6Zjk6Mzg6ODk6NmI6ZjIiLCJhcHAiOiJteUFwcE5hbWUiLCJleHAiOjE2MDAzNDY2NjksImYiOiJzb21lRmVhdHVyZSJ9.RQgjEgvyZvieGfaCyvysBmzWGt_nzXOXMo8judz6G-k
```

This is the JWT that encodes all this information.

## Installing and Managing Licenses

Now that you have this JWT, you need to send it to the customer and
have customer install it.  For this, the customer needs to use the
`jwtmgr` tool.  Fortunately, because this is a Go program, it can be
trivially built as a statically compiled binary for all platforms.

The customer would then run:

```
$ jwtmgr add eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhIjoiNWM6Zjk6Mzg6ODk6NmI6ZjIiLCJhcHAiOiJteUFwcE5hbWUiLCJleHAiOjE2MDAzNDY2NjksImYiOiJzb21lRmVhdHVyZSJ9.RQgjEgvyZvieGfaCyvysBmzWGt_nzXOXMo8judz6G-k
License successfully installed
```

This adds the JWT to a file that contains all licenses that the user
currently posses.  After running the `jwtmgr add` command, the user
can list all available licenses with:

```
$ jwtmgr list
Application: myAppName
  Feature: someFeature
```

Note that when adding licenses, a check is done to make sure that the
JWT has not been corrupted.  So if the JWT was somehow corrupted (for
example, by accidentally dropping the first character) before reaching
the user, this would be detected, *e.g.,*

```
$ jwtmgr add yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhIjoiNWM6Zjk6Mzg6ODk6NmI6ZjIiLCJhcHAiOiJteUFwcE5hbWUiLCJleHAiOjE2MDAzNDY2NjksImYiOiJzb21lRmVhdHVyZSJ9.RQgjEgvyZvieGfaCyvysBmzWGt_nzXOXMo8judz6G-k
License yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhIjoiNWM6Zjk6Mzg6ODk6NmI6ZjIiLCJhcHAiOiJteUFwcE5hbWUiLCJleHAiOjE2MDAzNDY2NjksImYiOiJzb21lRmVhdHVyZSJ9.RQgjEgvyZvieGfaCyvysBmzWGt_nzXOXMo8judz6G-k is not a valid JWT
```

## Licensing Models

The basic idea here is to encode the details of whether the user
"holds" a valid license for a given feature of a given application.
This can be done in two ways:

### Node-Locked Licenses

### Floating Licenses


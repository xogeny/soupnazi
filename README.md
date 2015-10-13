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

## Approach

The basic idea here is to encode the details of whether the user
"holds" a valid license for a given feature of a given application.
This can be done in two ways:

### Node-Locked Licenses

### Floating Licenses


# License Management with JWTs?

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

# Approach

## Node-Locked Licenses

Currently, there is only one type of grant.  It can be represented in
JSON as:

...

## Floating Licenses

Eventually, I'd like to include another type of JWT that includes some
kind of GUID and then contacts a server to determine if that GUID
should be granted what is being requested.  This would enable what
essentially amounts ot a floating license server implementation.

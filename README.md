# noeq53d - A fault-tolerant network service for meaningful GUID generation

Based on [noeqd][].

## Motivation

This is a fork of the noeqd project with a couple of significant differences:

* IDs generated will be 53 bits as to comply with the maximum integer value supported in Javascript.
* An additional 1-byte ID space has been added to the request protocol.

## Details

The IDs generated by noeq53d will be no greater than 53 bits so as to comply with the
Javascript maximum integer size.  The GUID format has been slightly modified:

* time - 39 bits (custom epoch (shepoch) gives us 17+ years)
* workstation id - 4 bits
* sequence id - 10 bits (1024 ids per workstation per second)

We have deviated from the original noeqd protocol by adding an additional ID space to the request protocol.
This is expressed as an additional byte that is appended to the initial byte that indicates the number of 
IDs to generate:

		----------------------
		|<num byte>|<ID byte>|
		----------------------

The ID space adds another level of scope to ID generation.  For example, this space could be used to scope
IDs to individual database tables.

The auth request from noeqd is not supported.

See the [noeqd][] project for specifics that are not covered here.

[noeqd]: https://github.com/bmizerany/noeqd

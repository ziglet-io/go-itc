# Interval Tree Clocks in Go

In distributed systems, a central problem is managing causality in
operations applied at different sites to the same data so that,
eventually, a single version of the data can be agreed upon by all
participants. 

The problem of keeping track of logical time was explored by
[Lamport][1] in his development of Vector Clocks. Subsequently, Version
Vectors were developed to handle a similar problem in state replication.

Both Vector Clocks and Version Vectors _stamp_ a version of the shared
state with either a vector of local clock values - usually integers
starting at 0 and incrementing by 1 - or an array of tuples (_sid_,
_time_) where the _sid_ is a unique identifier for the processor and
_time_ is the same as the above mentioned integer counter.

In systems where the number of participants is unknown and potentially
large, maintaining the _stamps_ for all transactions can grow very large
and even exceed the size of the data itself. 

To address this problem, [Almeida et al][2] developed Interval Tree
Clocks (ITCs).

This library provides a reference implementation of Interval Tree Clocks
for the [Go][3] programming language.

# References

* [1]: https://lamport.azurewebsites.net/pubs/time-clocks.pdf Lamport,
  Leslie. "Time, Clocks, and the Ordering of Events in a Distributed
  System"
* [2]: http://gsd.di.uminho.pt/members/cbm/ps/itc2008.pdf
  Almeida, Paulo Sergio, Baquero, Carlos, and Fonte Victor "Interval
  Tree Clocks: A Logical Clock for Dynamic Systems"
* [3]: https://golang.org The Go Programming Language
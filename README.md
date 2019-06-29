![Interval Tree Clocks](./doc/images/header.png)
# Interval Tree Clocks in Go

In distributed systems, a central problem is managing causality in
operations applied at different sites to the same data so that,
eventually, a single version of the data can be agreed upon by all
participants. 

The problem of keeping track of logical time was explored by 
[Lamport](#references) in his development of Vector Clocks. Subsequently, Version
Vectors were developed to handle a similar problem in state replication.

Both Vector Clocks and Version Vectors _stamp_ a version of the shared
state with either a vector of local clock values - usually integers
starting at 0 and incrementing by 1 - or an array of tuples (_sid_,
_time_) where the _sid_ is a unique identifier for the processor and
_time_ is the same as the above mentioned integer counter.

In systems where the number of participants is unknown and potentially
large, maintaining the _stamps_ for all transactions can grow very large
and even exceed the size of the data itself. With Vector Clocks and
Version Vectors (depending on the implementation), once a participant
has been added to the vector, that participant will be included in every
subsequent version as a _(sid,counter)_ pair. This leads to a huge
explosion in the amount of information required to track changes.

To address this problem, [Almeida et al](#references) developed Interval
Tree Clocks (ITCs). Interval Tree Clocks are similar to Version Vectors
except that they do not require one entry per _sid_. Instead, both the
identifiers and the clock values can grow and shrink in size
dynamically. 

Operations on ITCs are characterized by a **fork-event-join** model.
When a site (possibly a new site) wishes to begin work on the data, it
**forks** the data structure. Thereafter, it begins to execute one or
more **events**. Upon receipt of events from other participants, the
site **joins** those **events**. The **join** may result in the removal
of the current **Id**. If the site wishes to continue submitting
**events**, it again issues a **fork**. This process of
**fork-event(s)-join** continues.

This library provides a reference implementation of Interval Tree Clocks
for the [Go](#references) programming language.

# Protobuf

Since this library is intended to be used in a *distributed* system, it
required a binary encoding for on-the-wire transfer. Protobufs seemed an
excellent choice for this. To the extent the on-the-wire protocol is
compatible with the in-memory objects, I decided not to translate and
just to use the protobufs directly.

# Experience

The promise of ITCs is to permit **local** assignment of new sites
through **fork** and retirement of sites through **join** yielding a
dynamic identifier size and significant savings.

**However** the two parties that will receive the results of either **fork** or **join** must coordinate at the time that the subdomains change.
 
 For an example of why this is, examine the test _TestExampleUnknownSplitJoin_
[here](./itc_test/Example_test.go). If we wanted a system where nodes did **not** need to collaborate during **fork** and **join** the desired flow would be :

* An initial state is created with a seed stamp : _(1,0)_
* Site _A_ begins writing yielding stamps like : _(1,1)_ ... _(1,n)_
* Meanwhile site _B_ wants a different identifier range and therefore
  splits _(1,0)_ to obtain _((1,0),1)_. After a write at _B_ ->
  _((1,0),2)_.
* Later, _A_ learns of the writes by _B_ and wishes to join them. This
  does not work because the domains of _A_ and _B_ overlap.
  
In a context where the number of participants is relatively small and all are online together most of the time, this solution could be helpful. Consider a system of N replicas of a database. Normal operation has the database consisting of a small number of replicas all of whom are online together. Network partition or node failure is expected to quickly repaired. In this environment, ITCs could be an interesting solution to assigning portions of the identifier space to the replicas.

Contrast this setting with Wikipedia, a collaborative editing environment where the number of participants theoretically could include the entire population and where participants go on and offline regularly. In this environment, the requirement for participants to communicate to partition the identifier space is unreasonable.
 
# References

* [1]: Lamport, Leslie. ["Time, Clocks, and the Ordering of Events in a Distributed
        System"](https://lamport.azurewebsites.net/pubs/time-clocks.pdf)
* [2]: Almeida, Paulo Sergio, Baquero, Carlos, and Fonte Victor ["Interval Tree Clocks: A Logical Clock for Dynamic Systems"](http://gsd.di.uminho.pt/members/cbm/ps/itc2008.pdf)
* [3]: The Go Authors. [The Go Programming Language](https://golang.org)
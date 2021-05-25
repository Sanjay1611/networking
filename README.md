# Networking

## Introduction
This sample code is a simple demo application of a distributed database where we store a particular word(i.e name) and their cid(content identifier) in database
of a particular node in a distributed system and later we can fetch it from any of the system in that network.

## Modes
It operates in either of 2 modes
* Write Mode -
In this mode we take a word from the user i.e. User enters a word of his choice in the cli, later we generate cid of that word and store in the badgerDB along
with the word as a (key, value) pair i.e. (cid, word).
* Read Mode -
It can be seen as an enquiry counter the user will come up with any valid cid and will get the corresponding word if it exists in the database of any node in the network.

## Flags allowed
* dir - To specify the working directory of badgerDB
* c - **true** for write mode and **false** for read mode
* peers - To mention any peers to connect to
* p - to specify port

## Solution Process
1. Firstly we create a libp2p host on the specified port or the default one (9900)
2. Then create a DHT out of that host. 
3. Later we connect all the peers mentioned if any in the flag else the boostrap peers given by dht package.
4. Open a badgerDB instance
5. Afterwards processing depends on the mode
* Write Mode - We take the word rom user generate its cid, store them in database and then broadcast the cid over the dht peers. Then create an rpc server and make it available 
as a go routine so that it can again take another word from the user.
* Read Mode - We take the cid from the user enquire about that cid over the dht peers. If any peer is there with the cid we make an rpc call to that peer and get the corresponding
word to display.

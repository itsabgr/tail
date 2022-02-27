# Tail
a remote sorted value store

Tail is a service with a level DB in the heart that is accessible with HTTP(S) for web browsers and UDP for faster access;
you can insert any value into the Tail and retrieve them sorted.
insertion:
a client can insert a value into the Tail with an HTTP request that body is as the value or a UDP packet that whole of the packet is as value.
Tail receives values and inserts them into level DB that the key of the key-value pair is the whole of the packet, and the value is some metadata like the insertion time.
that means Tail has no duplicate value while keeping them sorted.
after insertion, Tail returns the next available value is placed right after in the sorted list into the HTTP response body or as a UDP packet.


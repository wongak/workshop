# workshop

My little workshop

Software development shares many similarities with a craft.
This is a collection ideas, approaches, philosophies, and libraries. To stay
with the analogy of craft, this represents my own workshop. I have been using
the tools for some time. They feel natural to me.
I am using this project as a way to more clearly explain the choices and the
reasoning behind how this workshop layout came to be.
The primary function of the workshop is to create web applications. The
applications have been mostly used in the development of early stage digital
products. They focus on testing hypotheses quickly.
The products should be malleable, integrate with as many services as possible,
and minimize maintenance time. I sacrifice developer ergonomics and reusability.
Everything is hand-written. Everything has a cost, and I need to feel it.

## HTTP and Go

The most natural way to interact with HTTP in Go is to follow the data. The
journey of an `*http.Request` and the `http.ResponseWriter` through our
application is the narrative structure.


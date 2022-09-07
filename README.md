# httpkit

A lightweight kit for creating HTTP-based services in Go.

## Why `httpkit`?

There are a lot of good frameworks out there. I used to use `gorilla/mux` all the time.
Eventually I started moving away from that and back towards the standard library.

Obviously using the standard library is a bit more verbose. Over time, I've noticed that
I scaffold projects in very similar ways, copy/pasting from project to project. The last
couple of times I've done it, I've promised myself I was going to put it all in a package
before starting the next project. This is me keeping that promise to myself.

## Who Should Use This?

Probably not you. Most people are perfectly fine using more fully-featured routers and
frameworks like `gorilla/mux`. However, if you really want to stick with the standard library
you're more than welcome to use bits and pieces of this that you like, or even all of it.

## What This Is Not

I don't see this as a framework. I see it as some code that is useful for working with Go's
standard `net/http` package. For example, I _always_ add signal catching for gracefully
shutting down the web server. I _always_ add logging of requests. I _always_ add a request ID
and a logger at minimum to every request context. However, the core logic of the application
is handled by standard `http.Handler` types. 

## How Do I Use This?

Check out the [`example`](./example/) for a quick little demo on what you can do with this package.

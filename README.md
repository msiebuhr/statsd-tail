statsd-tail
===========

Listens for statsd (in Datadog's dialect) and pretty-prints it on the console.

    > statsd-tail
    foo.bar  map[key:value]  1
    foo.bar  map[key:value]  2
    ...

Getting it
----------

Provided a working Go-setup, fetch the code with `go get
github.com/msiebuhr/statsd-tail` and then `go install
github.com/msiebuhr/statsd-tail`.

License
-------

MIT - Do what you want!

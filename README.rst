jsonpd
======

.. image:: https://travis-ci.org/mattrobenolt/jsonpd.svg?branch=master
   :target: https://travis-ci.org/mattrobenolt/jsonpd

..

    | Transform your jsonp requests into ESI so they can be cacheable in Varnish

::

    $ jsonpd -h
    Usage of jsonpd:
      -b="localhost:8000": bind address (default: localhost:8000)
      -cb="callback": callback argument (default: callback)
      -i="localhost:8001": bind address for stats (default: localhost:8001)
      -n=4: num procs (default: 4)
      -t=500: timeout (default: 500)

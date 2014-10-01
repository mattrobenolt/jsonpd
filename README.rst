jsonpd
======

.. image:: https://travis-ci.org/mattrobenolt/jsonpd.svg?branch=master
   :target: https://travis-ci.org/mattrobenolt/jsonpd

..

    | Transform your jsonp requests into ESI so they can be cacheable in Varnish

::

    $ jsonpd -h
    Usage of jsonpd:
      -b=":8000": bind address (default: :8000)
      -cb="callback": callback argument (default: callback)
      -i=":8001": bind address for stats (default: 8001)
      -n=4: num procs (default: 4)
      -t=500: timeout (default: 500)

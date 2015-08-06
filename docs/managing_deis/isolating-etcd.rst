:title: Isolating etcd
:description: Configuring the cluster to isolate etcd

.. _isolating-etcd:

Isolating etcd
==============

.. include:: ../_includes/_isolating-etcd-description.rst

.. note::

    The approach documented here works as of Deis 1.9.  Older versions of Deis
    utilize an older version of etcd that did not include the proxy
    functionality.

cloud-config
------------

To realize the topology described above, it is necessary, at the time of
provisioning, to provide different cloud-config for those hosts that will run
etcd and for those that will only run an etcd proxy.

For the small, fixed number of hosts running full etcd and satisfying the
"central services" role (as described in the CoreOS documentation), the
cloud-config provided with Deis is sufficient.

For hosts running only an etcd proxy, satisfying the "worker" role (as described
in the CoreOS documentation), cloud-config must be tweaked slightly to include
the ``-proxy on`` flag. For example:

.. code-block:: yaml

    #cloud-config
    ---
    coreos:
      # ...
      - name: etcd.service
        command: start
        content: |
          # ...
          [Service]
          # ...
          ExecStart=/usr/bin/docker run --net=host --rm \
            # ...
            -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
            # ...
            --discovery <discovery url here> \
            -proxy on
          # ...
        # ...
      # ...
    # ...

Isolating etcd as described here requires subsets of a cluster's hosts to be
configured differently from one another (including or excluding the
``-proxy on`` flag). Deis provisioning scripts do not currently account for
this, so managing separate cloud-config for each subset of nodes in the cluster
is left as an exercise for the advanced operator.

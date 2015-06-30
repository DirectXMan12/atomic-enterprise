Troubleshooting
===============

This document contains some tips and suggestions for troubleshooting an Atomic Enterprise deployment.

System Environment
------------------

1. Run as root

   Currently Atomic Enterprise must be started as root in order to manipulate your iptables configuration.  The Atomic Enterprise commands (e.g. `oc create`) do not need to be run as root.

1. Properly configure or disable firewalld

   On Fedora or other distributions using firewalld: Add docker0 to the public zone

        $ firewall-cmd --zone=trusted --change-interface=docker0
        $ systemctl restart firewalld

    Alternatively you can disable it via:
    
        $ systemctl stop firewalld

1. Setup your host DNS to an address that the containers can reach

  Containers need to be able to resolve hostnames, so if you run a local DNS server on your host, you should update your /etc/resolv.conf to instead use a DNS server that will be reachable from within running containers.  Google's "8.8.8.8" server is a popular choice.

1. Save iptables rules before restarting iptables and restore them afterwards. If iptables have to be restarted, then the iptables rules should be saved and restored, otherwise the docker inserted rules would get lost.


        $ iptables-save > /path/to/iptables.bkp
        $ systemctl restart iptables
        $ iptables-restore < /path/to/iptables.bkp
  

Docker Registry
---------------

Most of the v3 flows today assume you are running a docker registry pod.  You should ensure that this local registry is running:

    $ atomic-enterprise admin registry

If it's running, you should see this:

    Docker registry "docker-registry" service exists

If it's not running, you will instead see:

    F0429 09:22:54.492990   25259 registry.go:154] Docker-registry "docker-registry" does not exist (no service). Pass --create to install.

If it's not running, you can launch it via:

    $ oadm registry --create --credentials="${KUBECONFIG}"

Probing Containers
------------------

In general you may want to investigate a particular container.  You can either gather the logs from a container via `docker logs [container id]` or use `docker exec -it [container id] /bin/sh` to enter the container's namespace and poke around.

If you simply want to take a container but debug it outside of the Master's knowledge you can run the following:

    $ docker commit <CONTAINER ID> <some new name>
    $ docker run -it <name from previous step> /bin/bash

Obviously this won't work if you don't have bash installed but you could always package it in for debugging purposes.


Benign Errors/Messages
----------------------

There are a number of suspicious looking messages that appear in the Atomic Enterprise log output which can normally be ignored:

1. Failed to find an IP for pod (benign as long as it does not continuously repeat)

        E1125 14:51:49.665095 04523 endpoints_controller.go:74] Failed to find an IP for pod: {{ } {7e5769d2-74dc-11e4-bc62-3c970e3bf0b7 default /api/v1beta1/pods/7e5769d2-74dc-11e4-bc62-3c970e3bf0b7  41 2014-11-25 14:51:48 -0500 EST map[template:ruby-helloworld-sample deployment:database-1 deploymentconfig:database name:database] map[]} {{v1beta1 7e5769d2-74dc-11e4-bc62-3c970e3bf0b7 7e5769d2-74dc-11e4-bc62-3c970e3bf0b7 [] [{ruby-helloworld-database mysql []  [{ 0 3306 TCP }] [{MYSQL_ROOT_PASSWORD rrKAcyW6} {MYSQL_DATABASE root}] 0 0 [] <nil> <nil>  false }] {0x1654910 <nil> <nil>}} Running localhost.localdomain   map[]} {{   [] [] {<nil> <nil> <nil>}} Pending localhost.localdomain   map[]} map[]}

1. Proxy connection reset 

        E1125 14:52:36.605423 04523 proxier.go:131] I/O error: read tcp 10.192.208.170:57472: connection reset by peer

1. No network settings

        W1125 14:53:10.035539 04523 rest.go:231] No network settings: api.ContainerStatus{State:api.ContainerState{Waiting:(*api.ContainerStateWaiting)(0xc208b29b40), Running:(*api.ContainerStateRunning)(nil), Termination:(*api.ContainerStateTerminated)(nil)}, RestartCount:0, PodIP:"", Image:"kubernetes/pause:latest"}

Must Gather
-----------
If you find yourself still stuck please recreate your issue with verbose logging and gather the following:

1. Atomic Enterprise logs at level 4 (verbose logging):

        $ atomic-enterprise start --loglevel=4 &> /tmp/atomic-enterprise.log
        
1. Container logs  
    
    The following bit of scripting will pull logs for **all** containers that have been run on your system.  This might be excessive if you don't keep a clean history, so consider manually grabbing logs for the relevant containers instead:

        for container in $(docker ps -aq); do
            docker logs $container >& $LOG_DIR/container-$container.log
        done

1. Authorization rules:

    If you are getting "forbidden" messages or 403 status codes that you aren't expecting, you should gather the policy bindings, roles, and rules being used for the namespace you are trying to access.  Run the following commands, substituting `<project-namespace>` with the namespace you're trying to access.

        $ oc describe policy default --namespace=master
        $ oc describe policybindings master --namespace=master
        $ oc describe policy default --namespace=<project-namespace>
        $ oc describe policybindings master --namespace=<project-namespace>
        $ oc describe policybindings <project-namespace> --namespace=<project-namespace>

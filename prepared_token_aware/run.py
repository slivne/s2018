#!/usr/bin/python
import subprocess
import time


def run_test():
    population=100000
    ccm_commands = ["ccm stop","ccm remove scylla-3-3","ccm create scylla-3-3 --scylla --vnodes -n 3:3 --install-dir=/home/shlomi/scylla","ccm start --wait-for-binary-proto --wait-other-notice", "ccm node1 stress write n={0} -schema 'replication(strategy=NetworkTopologyStrategy, dc1=1,dc2=1)' -pop seq=1..{0}".format(population),"ccm node1 stress read n={0} -rate threads=20  -pop seq=1..{0}".format(population)]
    for ccm_command in ccm_commands:
        p = subprocess.Popen(ccm_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        p.wait()

    time.sleep(120)

    print "prepared token aware"
    command_read = "/home/centos/scylla-tools-java/tools/bin/cassandra-stress read duration=120s -mode native cql3 -node datacenter=dc1 -rate threads=20  -pop seq=1..{0}".format(population) 
    p = subprocess.Popen(command_read, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for line in p.stdout.readlines():
           print line
    p.wait()

    time.sleep(120)

    print "unprepared"
    command_read = "/home/centos/scylla-tools-java/tools/bin/cassandra-stress read duration=120s -mode native unprepared cql3 -node datacenter=dc1 -rate threads=20 -pop seq=1..{0}".format(population) 
    p = subprocess.Popen(command_read, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for line in p.stdout.readlines():
           print line
    p.wait()
    
    time.sleep(120)
    
    print "prepared not token aware"
    command_read = "/home/centos/scylla-tools-java/tools/bin/cassandra-stress read duration=120s -mode native cql3 -node whitelist 127.0.0.1,127.0.0.2,127.0.0.3 -rate threads=20  -pop seq=1..{0}".format(population)
    p = subprocess.Popen(command_read, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for line in p.stdout.readlines():
           print line
    p.wait()

    time.sleep(120)
    
run_test();

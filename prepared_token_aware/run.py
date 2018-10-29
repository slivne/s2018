#!/usr/bin/python
import subprocess
import time


def run_test():
    ccm_commands = ["ccm stop","ccm remove scylla-3-3","ccm create scylla-3-3 --scylla --vnodes -n 3:3 --install-dir=/home/shlomi/scylla","ccm start --wait-for-binary-proto --wait-other-notice", "ccm node1 stress write n=100000 -schema 'replication(strategy=NetworkTopologyStrategy, dc1=1,dc2=1)'","ccm node1 stress read n=100000 -rate threads=20"]
    for ccm_command in ccm_commands:
        p = subprocess.Popen(ccm_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        p.wait()

    print "unprepared"
    command_read = "/home/centos/scylla-tools-java/tools/bin/cassandra-stress read n=100000 -mode native unprepared cql3 -node datacenter=dc1 -rate threads=20" 
    p = subprocess.Popen(command_read, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for line in p.stdout.readlines():
           print line
    p.wait()
    
    time.sleep(120)
    
    print "prepared not token aware"
    command_read = "/home/centos/scylla-tools-java/tools/bin/cassandra-stress read n=100000 -mode native cql3 -node whitelist 127.0.0.1,127.0.0.2,127.0.0.3 -rate threads=20" 
    p = subprocess.Popen(command_read, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for line in p.stdout.readlines():
           print line
    p.wait()

    time.sleep(120)
    
    print "prepared token aware"
    command_read = "/home/centos/scylla-tools-java/tools/bin/cassandra-stress read n=100000 -mode native cql3 -node datacenter=dc1 -rate threads=20" 
    p = subprocess.Popen(command_read, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for line in p.stdout.readlines():
           print line
    p.wait()

run_test();

#!/usr/bin/python
import subprocess
import time
import threading


def run_test():
    population=1000000
    ccm_commands = ["ccm stop","ccm remove scylla-3","ccm create scylla-3 --scylla --vnodes -n 3 --install-dir=/home/shlomi/scylla","ccm start --wait-for-binary-proto --wait-other-notice", "ccm node1 stress write n={0} -pop seq=1..{0} -schema 'replication(strategy=NetworkTopologyStrategy, datacenter1=3)'".format(population),"ccm node1 stress read duration=60s -rate threads=20 -pop seq=1..{0}".format(population)]
    for ccm_command in ccm_commands:
        p = subprocess.Popen(ccm_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        p.wait()

    event = threading.Event()

    def run():
        try:
            time.sleep(60)
            for x in range(4):
                command_write = "/home/centos/scylla-tools-java/tools/bin/cassandra-stress write duration=60s -mode native cql3 -rate threads=20 fixed=2000/s -pop seq=1..{0}".format(population) 
                print command_write
                p = subprocess.Popen(command_write, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
                for line in p.stdout.readlines():
                       print line
                p.wait()
                time.sleep(60)
        finally:
            event.set()
            pass

    t = threading.Thread(target=run)
    t.setDaemon(True)
    t.start()


    command_read = "/home/centos/scylla-tools-java/tools/bin/cassandra-stress read duration=600s -mode native cql3 -rate threads=20 fixed=1000/s -pop seq=1..{0}".format(population) 
    print command_read
    p = subprocess.Popen(command_read, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for line in p.stdout.readlines():
           print line
    p.wait()

    event.wait()

run_test();

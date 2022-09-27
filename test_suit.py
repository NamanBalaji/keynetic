import subprocess
import unittest
import requests
import time
import os

TIMEOUT=60

######################## initialize variables ################################################
subnetName = "assignment4-net"
subnetAddress = "10.10.0.0/16"

nodeIpList = ["10.10.0.2", "10.10.0.3", "10.10.0.4", "10.10.0.5", "10.10.0.6", "10.10.0.7"]
nodeHostPortList = ["8082","8083", "8084", "8086", "8087", "8088"]
nodeSocketAddressList = [nodeIp + ":8085" for nodeIp in nodeIpList]
nodeNameList = ["node{}".format(n + 1) for n, _ in enumerate(nodeIpList)]
assert len(nodeHostPortList) == len(nodeIpList) == len(nodeSocketAddressList) == len(nodeNameList)

view = ""
for nodeSocketAddress in nodeSocketAddressList:
    view += nodeSocketAddress + ","
view = view[:-1]

shardCount = 2

############################### Docker Linux Commands ###########################################################

def subnetExists(subnetName):
    matchingSubnets = subprocess.check_output(['docker', 'network', 'list', '-f' 'name={}'.format(subnetName), '-q']).split()
    return 0 < len(matchingSubnets)

def removeSubnet(subnetName, required=True):
    command = "docker network rm " + subnetName
    print(command)
    subprocess.run(command, shell=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=required)

def createSubnet(subnetAddress, subnetName):
    command  = "docker network create --subnet=" + subnetAddress + " " + subnetName
    print(command)
    subprocess.check_call(command, stdout=subprocess.DEVNULL, shell=True)

def buildDockerImage():
    command = "docker build -t assignment4-img ."
    print(command)
    os.system(command)

def runInstance(hostPort, ipAddress, subnetName, instanceName, view=view, giveShardCount=True):
    command = ['docker', 'run', '--rm', '--detach',
            '--publish', "{}:8085".format(hostPort),
            "--net={}".format(subnetName),
            "--ip={}".format(ipAddress),
            "--name={}".format(instanceName),
            "-e", "SOCKET_ADDRESS={}:8085".format(ipAddress),
            "-e", "VIEW={}".format(view),
            ]
    if giveShardCount:
        command.extend(["-e", "SHARD_COUNT=" + str(shardCount)])
    # also specify the image
    command.append("assignment4-img")
    print(' '.join(command))
    subprocess.check_call(command, stdout=subprocess.DEVNULL)

def killInstance(instanceName, required=True):
    # kill is sufficient when containers are run with `--rm`
    command = "docker kill " + instanceName
    print(command)
    subprocess.run(command, shell=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=required)

def stopAndRemoveInstance(instanceName, required=True):
    command = "docker stop " + instanceName
    print(command)
    subprocess.run(command, shell=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=required)

    command = "docker rm " + instanceName
    print(command)
    subprocess.run(command, shell=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=required)

def connectToNetwork(subnetName, instanceName):
    command = "docker network connect " + subnetName + " " + instanceName
    print(command)
    subprocess.check_call(command, stdout=subprocess.DEVNULL, shell=True)

def disconnectFromNetwork(subnetName, instanceName):
    command = "docker network disconnect " + subnetName + " " + instanceName
    print(command)
    subprocess.check_call(command, stdout=subprocess.DEVNULL, shell=True)

############################### View Comparison Function ###########################################################
def compareViews(returnedView, expectedView):
    expectedView = expectedView.split(',')
    if (type(returnedView) is not list):
        returnedView = returnedView.split(',')
    else:
        returnedView = returnedView
    returnedView.sort()
    expectedView.sort()
    return returnedView == expectedView
############################### Array Comparison Function ###########################################################
def compareLists(list1, list2):
    list1.sort()
    list2.sort()
    return list1 == list2
################################# Unit Test Class ############################################################

class TestHW4(unittest.TestCase):

    # class/static properties used as shared-state across all tests
    shardIdList = []
    shardsMemberList = []
    keyCount = 600
    causalMetadata = ''

    @classmethod
    def setUpClass(cls):
        print("= Tear down old resources from a previous run")
        stopAndRemoveInstance('node7', required=False) # remove node7 if it exists
        for nodeName in nodeNameList:
            stopAndRemoveInstance(nodeName, required=False) # we use stop-and-remove here because it's not guaranteed that these instances were started with --rm
        removeSubnet(subnetName, required=False)
        assert not subnetExists(subnetName), 'subnet %s must be removed before the test can start' % subnetName

        print("= Create resources for this run")
        buildDockerImage()
        createSubnet(subnetAddress, subnetName)

        assert len(nodeHostPortList) == len(nodeIpList) == len(nodeNameList)
        for nodeHostPort, nodeIp, nodeName in zip(nodeHostPortList, nodeIpList, nodeNameList):
            runInstance(nodeHostPort, nodeIp, subnetName, nodeName)
        time.sleep(5) # give time for them to bind ports, update views, etc..
        os.system("docker ps")

    @classmethod
    def tearDownClass(cls):
        print("= Tear down resources")
        killInstance('node7', required=False) # remove node7 if it exists
        for nodeName in nodeNameList:
            killInstance(nodeName)
        removeSubnet(subnetName)

    def test_a_get_shard_ids(self):

        print("\n###################### Getting Shard IDs ######################\n")

        # get the shard IDs from node1
        response = requests.get( 'http://localhost:8082/key-value-store-shard/shard-ids', timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shardIdsFromNode1 = responseInJson['shard-ids']
        self.assertEqual(len(shardIdsFromNode1), shardCount)

        # get the shard IDs from node5
        response = requests.get( 'http://localhost:8087/key-value-store-shard/shard-ids', timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shardIdsFromNode5 = responseInJson['shard-ids']
        self.assertTrue(compareLists(shardIdsFromNode5, shardIdsFromNode1))

        # get the shard IDs from node6
        response = requests.get( 'http://localhost:8088/key-value-store-shard/shard-ids', timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shardIdsFromNode6 = responseInJson['shard-ids']
        self.assertTrue(compareLists(shardIdsFromNode6, shardIdsFromNode1))

        self.shardIdList += shardIdsFromNode1


    def test_b_shard_id_members(self):

        print("\n###################### Getting the Members of Shard IDs ######################\n")

        shard1 = str(self.shardIdList[0])
        shard2 = str(self.shardIdList[1])

        # get the members of shard1 from node2
        response = requests.get( 'http://localhost:8083/key-value-store-shard/shard-id-members/' + shard1, timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shard1Members = responseInJson['shard-id-members']
        self.assertGreater(len(shard1Members), 1)

        # get the members of shard2 from node3
        response = requests.get( 'http://localhost:8084/key-value-store-shard/shard-id-members/' + shard2, timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shard2Members = responseInJson['shard-id-members']
        self.assertGreater(len(shard2Members), 1)

        self.assertEqual(len(nodeSocketAddressList), len(shard1Members + shard2Members))

        self.shardsMemberList += [shard1Members]
        self.shardsMemberList += [shard2Members]


    def test_c_node_shard_id(self):

        print("\n###################### Getting the Shard ID of the nodes ######################\n")

        shard1 = self.shardIdList[0]

        # get the shard id of node1
        response = requests.get( 'http://localhost:8082/key-value-store-shard/node-shard-id', timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)

        node1ShardId = responseInJson['shard-id']

        self.assertTrue(node1ShardId in self.shardIdList)

        if node1ShardId == shard1:
            self.assertTrue(nodeSocketAddressList[0] in self.shardsMemberList[0])
        else:
            self.assertTrue(nodeSocketAddressList[0] in self.shardsMemberList[1])

        # get the shard id of node2
        response = requests.get('http://localhost:8083/key-value-store-shard/node-shard-id', timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)

        node2ShardId = responseInJson['shard-id']

        self.assertTrue(node2ShardId in self.shardIdList)

        if node2ShardId == shard1:
            self.assertTrue(nodeSocketAddressList[1] in self.shardsMemberList[0])
        else:
            self.assertTrue(nodeSocketAddressList[1] in self.shardsMemberList[1])

        # get the shard id of node6
        response = requests.get('http://localhost:8088/key-value-store-shard/node-shard-id', timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)

        node6ShardId = responseInJson['shard-id']

        self.assertTrue(node6ShardId in self.shardIdList)

        if node6ShardId == shard1:
            self.assertTrue(nodeSocketAddressList[5] in self.shardsMemberList[0])
        else:
            self.assertTrue(nodeSocketAddressList[5] in self.shardsMemberList[1])


    def test_d_put_key_value_operation(self):

        print("\n###################### Putting keys/values to the store ######################\nThis takes awhile\n")

        for counter in range(self.keyCount):
            nodeIndex = counter % len(nodeIpList)

            # put a new key in the store
            response = requests.put('http://localhost:' + nodeHostPortList[nodeIndex] + '/key-value-store/key' + str(counter), json={'value': "value" + str(counter), "causal-metadata": self.causalMetadata}, timeout=TIMEOUT)
            responseInJson = response.json()
            self.assertEqual(response.status_code, 201)
            self.causalMetadata = responseInJson["causal-metadata"]

            keyShardId = responseInJson["shard-id"]

            self.assertTrue(keyShardId in self.shardIdList)

            print('.', end='', flush=True)
        print()
        time.sleep(5) # give time for replication to settle down

    def test_e_get_key_value_operation(self):

        print("\n###################### Getting keys/values from the store ######################\n")

        nextCausalMetadata = ""

        for counter in range(self.keyCount):

            nodeIndex = (counter + 1 ) % len(nodeIpList)

            # get the value of the key
            response = requests.get('http://localhost:' + nodeHostPortList[nodeIndex] + '/key-value-store/key' + str(counter), json={"causal-metadata": self.causalMetadata}, timeout=TIMEOUT)
            responseInJson = response.json()
            self.assertEqual(response.status_code, 200)
            value = responseInJson["value"]
            self.assertEqual(value, "value" + str(counter))
            self.causalMetadata = responseInJson["causal-metadata"]

            print('.', end='', flush=True)
        print()


    def test_f_shard_key_count(self):

        print("\n###################### Getting key count of each shard ######################\n")

        shard1 = str(self.shardIdList[0])
        shard2 = str(self.shardIdList[1])

        # get the shard1 key count from node5
        response = requests.get( 'http://localhost:8087/key-value-store-shard/shard-id-key-count/' + shard1, timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shard1KeyCount = int(responseInJson['shard-id-key-count'])

        # get the shard2 key count from node3
        response = requests.get( 'http://localhost:8084/key-value-store-shard/shard-id-key-count/' + shard2, timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shard2KeyCount = int(responseInJson['shard-id-key-count'])

        # sum of key counts in shards == total keys
        self.assertEqual(self.keyCount, shard1KeyCount + shard2KeyCount)

        # check whether keys distributed almost uniformly
        minKeyCount = int ((self.keyCount * 0.75) / shardCount)
        maxKeyCount = int ((self.keyCount * 1.25) / shardCount)

        # minKeyCount < shard2-key-count < maxKeyCount
        self.assertGreater(shard1KeyCount, minKeyCount)
        self.assertLess(shard1KeyCount, maxKeyCount)

        # minKeyCount < shard2-key-count < maxKeyCount
        self.assertGreater(shard2KeyCount, minKeyCount)
        self.assertLess(shard2KeyCount, maxKeyCount)


    def test_g_add_new_node(self):

        shard2 = self.shardIdList[1]

        print("\n###################### Adding a new node ######################\n")

        node7Ip = "10.10.0.8"
        node7HostPort = "8089"
        node7SocketAddress = "10.10.0.8:8085"
        newView =  view + "," + node7SocketAddress

        runInstance(node7HostPort, node7Ip, subnetName, "node7", view=newView, giveShardCount=False)

        time.sleep(5) # give time for instance to bind ports, update views, etc..

        # get the new view from node1
        response = requests.get( 'http://localhost:8082/key-value-store-view', timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertTrue(compareViews(responseInJson['view'], newView))

        print("\n###################### Assigning the new node to the second shard ######################\n")

        response = requests.put('http://localhost:8082/key-value-store-shard/add-member/' + str(shard2), json={'socket-address': node7SocketAddress}, timeout=TIMEOUT)
        self.assertEqual(response.status_code, 200)

        time.sleep(5) # give time for possible rebalancing

        # get the shard id of node7
        response = requests.get('http://localhost:8089/key-value-store-shard/node-shard-id', timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        node7ShardId = responseInJson['shard-id']
        self.assertEqual(node7ShardId, shard2)

        # get the members of shard2 from node4
        response = requests.get( 'http://localhost:8086/key-value-store-shard/shard-id-members/' + str(shard2), timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertTrue(node7SocketAddress in responseInJson['shard-id-members'])

        # get shard2 key count from node7
        response = requests.get( 'http://localhost:8089/key-value-store-shard/shard-id-key-count/' + str(shard2), timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shard2KeyCountFromNode7 = int(responseInJson['shard-id-key-count'])

        # get shard2 key count from node3
        response = requests.get( 'http://localhost:8084/key-value-store-shard/shard-id-key-count/' + str(shard2), timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shard2KeyCountFromNode3 = int(responseInJson['shard-id-key-count'])

        self.assertEqual(shard2KeyCountFromNode7, shard2KeyCountFromNode3)


    def test_h_impossible_reshard(self):

        print("\n###################### Doing Impossible Resharding ######################\n")

        response = requests.put('http://localhost:8083/key-value-store-shard/reshard', json={'shard-count': 10}, timeout=TIMEOUT)
        self.assertEqual(response.status_code, 400)


    def test_i_possible_reshard(self):

        print("\n###################### Doing Resharding ######################\n")

        response = requests.put('http://localhost:8082/key-value-store-shard/reshard', json={'shard-count': 3}, timeout=TIMEOUT)
        self.assertEqual(response.status_code, 200)

        time.sleep(20) # give time for resharding and rebalancing to settle down

        # get the new shard IDs from node1
        response = requests.get( 'http://localhost:8082/key-value-store-shard/shard-ids', timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        newShardIds = responseInJson['shard-ids']
        newShardIdList = newShardIds
        self.assertEqual(len(newShardIdList), 3)

        shard1 = newShardIdList[0]
        shard2 = newShardIdList[1]
        shard3 = newShardIdList[2]

        # get the members of shard1 from node2
        response = requests.get( 'http://localhost:8083/key-value-store-shard/shard-id-members/' + str(shard1), timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shard1Members = responseInJson['shard-id-members']
        self.assertGreater(len(shard1Members), 1)

        # get the members of shard2 from node3
        response = requests.get( 'http://localhost:8084/key-value-store-shard/shard-id-members/' + str(shard2), timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shard2Members = responseInJson['shard-id-members']
        self.assertGreater(len(shard2Members), 1)

        # get the members of shard3 from node4
        response = requests.get( 'http://localhost:8086/key-value-store-shard/shard-id-members/' + str(shard3), timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shard3Members = responseInJson['shard-id-members']
        self.assertGreater(len(shard3Members), 1)

        self.assertEqual(len(shard1Members + shard2Members + shard3Members), len(nodeSocketAddressList) + 1)

        # get the shard id of node4
        response = requests.get( 'http://localhost:8086/key-value-store-shard/node-shard-id', timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)

        node4ShardId = responseInJson['shard-id']

        self.assertTrue(node4ShardId in newShardIdList)

        if node4ShardId == shard1:
            self.assertTrue(nodeSocketAddressList[3] in shard1Members)
        elif node4ShardId == shard2:
            self.assertTrue(nodeSocketAddressList[3] in shard2Members)
        else:
            self.assertTrue(nodeSocketAddressList[3] in shard3Members)

        # get the shard1 key count from node5
        response = requests.get('http://localhost:8087/key-value-store-shard/shard-id-key-count/' + str(shard1), timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shard1KeyCount = int(responseInJson['shard-id-key-count'])

        # get the shard2 key count from node3
        response = requests.get('http://localhost:8084/key-value-store-shard/shard-id-key-count/' + str(shard2), timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shard2KeyCount = int(responseInJson['shard-id-key-count'])

        # get the shard3 key count from node1
        response = requests.get('http://localhost:8082/key-value-store-shard/shard-id-key-count/' + str(shard3), timeout=TIMEOUT)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        shard3KeyCount = int(responseInJson['shard-id-key-count'])

        # sum of key counts in shards == total keys
        self.assertEqual(self.keyCount, shard1KeyCount + shard2KeyCount + shard3KeyCount)

        # check whether keys distributed almost uniformly
        minKeyCount = int ((self.keyCount * 0.75) / 3)
        maxKeyCount = int ((self.keyCount * 1.25) / 3)

        # minKeyCount < shard1-key-count < maxKeyCount
        self.assertGreater(shard1KeyCount, minKeyCount)
        self.assertLess(shard1KeyCount, maxKeyCount)

        # minKeyCount < shard2-key-count < maxKeyCount
        self.assertGreater(shard2KeyCount, minKeyCount)
        self.assertLess(shard2KeyCount, maxKeyCount)

        # minKeyCount < shard3-key-count < maxKeyCount
        self.assertGreater(shard3KeyCount, minKeyCount)
        self.assertLess(shard3KeyCount, maxKeyCount)

        for counter in range(self.keyCount):

            nodeIndex = (counter + 1 ) % len(nodeIpList)

            # get the value of the key
            response = requests.get('http://localhost:' + nodeHostPortList[nodeIndex] + '/key-value-store/key' + str(counter), json={"causal-metadata": self.causalMetadata}, timeout=TIMEOUT)
            responseInJson = response.json()
            self.assertEqual(response.status_code, 200)
            value = responseInJson["value"]
            self.assertEqual(value, "value" + str(counter))
            self.causalMetadata = responseInJson["causal-metadata"]

            print('.', end='', flush=True)
        print()


if __name__ == '__main__':
    unittest.main(verbosity=3, failfast=True)
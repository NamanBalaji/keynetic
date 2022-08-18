from re import T
import unittest
import requests
import time
import os

######################## initialize variables ################################################
hostname = 'localhost'
hostBaseUrl = 'http://'+hostname
subnetName = "assignment3-net"
subnetAddress = "10.10.0.0/16"

replica1Ip = "10.10.0.2"
replica1HostPort = "8082"
replica1SocketAddress = replica1Ip + ":8085"

replica2Ip = "10.10.0.3"
replica2HostPort = "8083"
replica2SocketAddress = replica2Ip + ":8085"

replica3Ip = "10.10.0.4"
replica3HostPort = "8084"
replica3SocketAddress = replica3Ip + ":8085"

view = replica1SocketAddress + "," + replica2SocketAddress + "," + replica3SocketAddress

############################### Docker Linux Commands ###########################################################
def removeSubnet(subnetName):
    command = "docker network rm " + subnetName
    os.system(command)
    time.sleep(2)

def createSubnet(subnetAddress, subnetName):
    command  = "docker network create --subnet=" + subnetAddress + " " + subnetName
    os.system(command)
    time.sleep(2)

def buildDockerImage():
    command = "docker build -t assignment3-img ."
    os.system(command)

def runReplica(hostPort, ipAddress, subnetName, instanceName):
    command = "docker run -d -p " + hostPort + ":8085 --net=" + subnetName + " --ip=" + ipAddress + " --name=" + instanceName + " -e SOCKET_ADDRESS=" + ipAddress + ":8085" + " -e VIEW=" + view + " assignment3-img"
    os.system(command)
    time.sleep(20)

def stopAndRemoveInstance(instanceName):
    stopCommand = "docker stop " + instanceName
    removeCommand = "docker rm " + instanceName
    os.system(stopCommand)
    time.sleep(2)
    os.system(removeCommand)

def connectToNetwork(subnetName, instanceName):
    command = "docker network connect " + subnetName + " " + instanceName
    os.system(command)

def disconnectFromNetwork(subnetName, instanceName):
    command = "docker network disconnect " + subnetName + " " + instanceName
    os.system(command)

############################### View Comparison Function ###########################################################
def compareLists(list1, list2):
    for x in set(list1 + list2):
        if list1.count(x) != list2.count(x):
            return False
    return True

def getPairs(data, delimeter):
    pair = data.split(delimeter)
    return pair

def compareViews(returnedView, expectedView):
    rv = getPairs(returnedView, ",")
    ev = getPairs(expectedView, ",")
    return compareLists(rv, ev)
    
def compareCausalMetadata(returnedData, expectedData):
    returnedPair = getPairs(returnedData, "/")
    expectedPair = getPairs(expectedData, "/")
    return compareLists(returnedPair, expectedPair)

################################# Unit Test Class ############################################################

class TestHW3(unittest.TestCase):

    ######################## Build docker image and create subnet ################################
    

    ########################## Run tests #######################################################
    def test_a_view_operations(self):

    

        print("\n###################### Getting the view from replicas ######################\n")
        # get the view from replica1
        baseUrl = hostBaseUrl + ':' + replica1HostPort + '/key-value-store-view'
        response = requests.get( baseUrl)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertTrue(compareViews(responseInJson['view'], view))

        # get the view from replica2
        baseUrl = hostBaseUrl + ':' + replica2HostPort + '/key-value-store-view'
        response = requests.get( baseUrl)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertTrue(compareViews(responseInJson['view'], view))

        # get the view from replica3
        baseUrl = hostBaseUrl + ':' + replica3HostPort + '/key-value-store-view'
        response = requests.get(baseUrl)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertTrue(compareViews(responseInJson['view'], view))

        print("\n###################### Putting key1/value1 to the store ######################\n")
        # put a new key in the store
        baseUrl = hostBaseUrl + ':' + replica1HostPort + '/key-value-store/key1'
        response = requests.put(baseUrl, json={'value': "value1", "causal-metadata": ""})
        first_put_metadata = response.json()['causal-metadata']
        self.assertEqual(response.status_code, 201)

        print("\n###################### Waiting for 10 seconds ######################\n")
        time.sleep(10)

        print("\n###################### Getting key1 from the replicas ######################\n")

        # get the value of the new key from replica1 after putting the new key
        baseUrl = hostBaseUrl + ':' + replica1HostPort + '/key-value-store/key1'
        response = requests.get(baseUrl, json={"causal-metadata": first_put_metadata })
        responseInJson = response.json()
        first_get_metadata = response.json()['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['value'], 'value1')

        # get the value of the new key from replica2 after putting the new key
        baseUrl = hostBaseUrl + ':' + replica2HostPort + '/key-value-store/key1'
        response = requests.get(baseUrl, json={"causal-metadata": first_get_metadata})
        responseInJson = response.json()
        second_get_metadata = response.json()['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['value'], 'value1')

        # get the value of the new key from replica3 after putting the new key
        baseUrl = hostBaseUrl + ':' + replica3HostPort + '/key-value-store/key1'
        response = requests.get(baseUrl, json={"causal-metadata":second_get_metadata})
        responseInJson = response.json()
        third_get_metadata = response.json()['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['value'], 'value1')

        print("\n###################### Stopping and removing replica3 ######################\n")
        stopAndRemoveInstance("replica3")

        print("\n###################### Waiting for 10 seconds ######################\n")
        time.sleep(10)

        print("\n###################### Putting key2/value2 to the store ######################\n")
        # put a new key in the store
        baseUrl  = hostBaseUrl + ':' + replica1HostPort + '/key-value-store/key2'
        response = requests.put(baseUrl, json={'value': "value2", "causal-metadata": third_get_metadata})
        second_put_metadata = response.json()['causal-metadata']
        self.assertEqual(response.status_code, 201)

        print("\n###################### Waiting for 50 seconds ######################\n")
        time.sleep(50)

        print("\n###################### Getting key2 from the replica1 and replica2 ######################\n")

        # get the value of the new key from replica1 after putting the new key
        baseUrl  = hostBaseUrl + ':' + replica1HostPort + '/key-value-store/key2'
        response = requests.get(baseUrl, json={"causal-metadata": second_put_metadata })
        responseInJson = response.json()
        fourth_get_metadata = response.json()['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['value'], 'value2')

        # get the value of the new key from replica2 after putting the new key
        baseUrl  = hostBaseUrl + ':' + replica2HostPort + '/key-value-store/key2'
        response = requests.get(baseUrl, json={"causal-metadata":fourth_get_metadata})
        responseInJson = response.json()
        fifth_get_metadata = response.json()['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['value'], 'value2')


        print("\n###################### Getting the view from replica1 and replica2 ######################\n")
        # get the view from replica1
        baseUrl  = hostBaseUrl + ':' + replica1HostPort + '/key-value-store-view'
        response = requests.get( baseUrl)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertTrue(compareViews(responseInJson['view'], replica1SocketAddress + "," + replica2SocketAddress))

        # get the view from replica2
        baseUrl  = hostBaseUrl + ':' + replica2HostPort + '/key-value-store-view'
        response = requests.get( baseUrl)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertTrue(compareViews(responseInJson['view'], replica1SocketAddress + "," + replica2SocketAddress))

        print("\n###################### Starting replica3 ######################\n")
        runReplica(replica3HostPort, replica3Ip, subnetName, "replica3")

        print("\n###################### Waiting for 50 seconds ######################\n")
        time.sleep(50)

        print("\n###################### Getting the view from replicas ######################\n")
        # get the view from replica1
        baseUrl  = hostBaseUrl + ':' + replica1HostPort + '/key-value-store-view'

        response = requests.get(baseUrl)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertTrue(compareViews(responseInJson['view'], view))

        # get the view from replica2
        baseUrl  = hostBaseUrl + ':' + replica2HostPort + '/key-value-store-view'
        response = requests.get(baseUrl)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertTrue(compareViews(responseInJson['view'], view))

        # get the view from replica3
        baseUrl  = hostBaseUrl + ':' + replica3HostPort + '/key-value-store-view'
        response = requests.get(baseUrl)
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertTrue(compareViews(responseInJson['view'], view))

        print("\n###################### Getting key1 and key2 from replica3 ######################\n")

        # get the value of the new key from replica3 after putting a new key
        baseUrl  = hostBaseUrl + ':' + replica3HostPort + '/key-value-store/key1'
        response = requests.get(baseUrl, json={'causal-metadata':fifth_get_metadata})
        responseInJson = response.json()
        sixth_get_metadata = response.json()['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['value'], 'value1')

        # get the value of the new key from replica2 after putting a new key
        baseUrl  = hostBaseUrl + ':' + replica3HostPort + '/key-value-store/key2'
        response = requests.get(baseUrl, json={"causal-metadata":sixth_get_metadata})
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['value'], 'value2')

    ########################## Availability Test #######################################################
    def test_b_availability(self):
    
        print("\n###################### Disconnecting replica2 from the network ######################\n")
        disconnectFromNetwork(subnetName, "replica2")
        time.sleep(1)

        print("\n###################### Disconnecting replica3 from the network ######################\n")
        disconnectFromNetwork(subnetName, "replica3")
        time.sleep(1)

        print("\n###################### Putting key/myvalue to the store ######################\n")

        # put a new key in the store
        baseUrl  = hostBaseUrl + ':' + replica1HostPort + '/key-value-store/key'
        response = requests.put(baseUrl, json={"value": "myvalue", "causal-metadata": ""})
        responseInJson = response.json()

        put_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 201)

    ########################## Key/Value Tests #######################################################
    def test_c_key_value_operations(self):

        

        print("\n###################### Putting mykey1/myvalue1 to the store ######################\n")

        # put a new key in the store
        baseUrl  = hostBaseUrl + ':' + replica1HostPort + '/key-value-store/mykey1'
        response = requests.put(baseUrl, json={"value": "myvalue1", "causal-metadata": ""})
        responseInJson = response.json()

        put_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 201)

        print("\n###################### Getting mykey1 from replica1 ######################\n")

        # get the value of the new key from replica1 after putting a new key
        baseUrl  = hostBaseUrl + ':' + replica1HostPort + '/key-value-store/mykey1'
        response = requests.get(baseUrl, json={"causal-metadata":put_causal_metadata})
        responseInJson = response.json()

        first_get_value = responseInJson['value']
        first_get_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(first_get_value, 'myvalue1')
        self.assertEqual(compareCausalMetadata(first_get_causal_metadata, put_causal_metadata), True)

        time.sleep(10)

        print("\n###################### Getting mykey1 from replica2 ######################\n")

        # get the value of the new key from replica2 after putting a new key
        baseUrl  = hostBaseUrl + ':' + replica2HostPort + '/key-value-store/mykey1'
        response = requests.get(baseUrl, json={'causal-metadata':first_get_causal_metadata})
        responseInJson = response.json()

        second_get_value = responseInJson['value']
        second_get_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(second_get_value, 'myvalue1')
        self.assertEqual(compareCausalMetadata(second_get_causal_metadata, first_get_causal_metadata), True)

        print("\n###################### Getting mykey1 from replica3 ######################\n")

        # get the value of the new key from replica3 after putting a new key
        baseUrl  = hostBaseUrl + ':' + replica3HostPort + '/key-value-store/mykey1'
        response = requests.get(baseUrl, json={'causal-metadata':second_get_causal_metadata})
        responseInJson = response.json()

        third_get_value = responseInJson['value']
        third_get_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(third_get_value, 'myvalue1')
        self.assertEqual(compareCausalMetadata(third_get_causal_metadata, second_get_causal_metadata), True)


        print("\n###################### Putting mykey1/myvalue2 to the store ######################\n")

        # put a new key in the store
        baseUrl  = hostBaseUrl + ':' + replica1HostPort + '/key-value-store/mykey1'
        response = requests.put(baseUrl, json={"value": "myvalue2", "causal-metadata": third_get_causal_metadata})
        responseInJson = response.json()
        second_put_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 200)

        print("\n###################### Getting mykey1 from replica1 ######################\n")

        # get the value of the new key from replica1 after putting a new key
        baseUrl  = hostBaseUrl + ':' + replica1HostPort + '/key-value-store/mykey1'
        response = requests.get(baseUrl, json={"causal-metadata":second_put_causal_metadata})
        responseInJson = response.json()

        fourth_get_value = responseInJson['value']
        fourth_get_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(fourth_get_value, 'myvalue2')
        self.assertEqual(compareCausalMetadata(fourth_get_causal_metadata, second_put_causal_metadata), True)

        time.sleep(10)

        print("\n###################### Getting mykey1 from replica2 ######################\n")

        # get the value of the new key from replica2 after putting a new key
        baseUrl  = hostBaseUrl + ':' + replica2HostPort + '/key-value-store/mykey1'
        response = requests.get(baseUrl, json={"causal-metadata":fourth_get_causal_metadata})
        responseInJson = response.json()

        fifth_get_value = responseInJson['value']
        fifth_get_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(fifth_get_value, 'myvalue2')
        self.assertEqual(compareCausalMetadata(fifth_get_causal_metadata, fourth_get_causal_metadata), True)

        print("\n###################### Getting mykey1 from replica3 ######################\n")

        # get the value of the new key from replica2 after putting a new key
        baseUrl  = hostBaseUrl + ':' + replica3HostPort + '/key-value-store/mykey1'
        response = requests.get(baseUrl, json={"causal-metadata":fifth_get_causal_metadata})
        responseInJson = response.json()

        sixth_get_value = responseInJson['value']
        sixth_get_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(sixth_get_value, 'myvalue2')
        self.assertEqual(compareCausalMetadata(sixth_get_causal_metadata, fifth_get_causal_metadata), True)

    ########################## Run tests #######################################################
    def test_d_causal_consistency(self):

        print("\n###################### Disconnecting replica2 from the network ######################\n")
        disconnectFromNetwork(subnetName, "replica2")
        time.sleep(0.5)

        print("\n###################### Putting k1/foo to the store ######################\n")

        # put a new key in the store
        baseUrl  = hostBaseUrl + ':' + replica1HostPort + '/key-value-store/k1'
        response = requests.put(baseUrl, json={"value": "foo", "causal-metadata": ""})
        responseInJson = response.json()

        first_put_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 201)

        print("\n###################### Connecting replica2 to the network ######################\n")
        connectToNetwork(subnetName, "replica2")

        print("\n###################### Putting k1/bar to the store ######################\n")

        # put  key1/bar in the store
        baseUrl  = hostBaseUrl + ':' + replica2HostPort + '/key-value-store/k1'
        response = requests.put(baseUrl, json={"value": "bar", "causal-metadata": first_put_causal_metadata})

        responseInJson = response.json()
        second_put_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 200)

        time.sleep(20)

        print("\n###################### Getting k1 from replica1 ######################\n")

        # get the value of the new key from replica1 after putting a new key
        baseUrl  = hostBaseUrl + ':' + replica1HostPort + '/key-value-store/k1'
        response = requests.get(baseUrl, json={"causal-metadata":second_put_causal_metadata})
        responseInJson = response.json()

        first_get_value = responseInJson['value']
        first_get_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(first_get_value, 'bar')
        self.assertEqual(compareCausalMetadata(first_get_causal_metadata, second_put_causal_metadata), True)

        print("\n###################### Getting k1 from replica2 ######################\n")

        # get the value of the new key from replica1 after putting a new key
        baseUrl  = hostBaseUrl + ':' + replica2HostPort + '/key-value-store/k1'
        response = requests.get(baseUrl, json={"causal-metadata":first_get_causal_metadata})
        responseInJson = response.json()

        second_get_value = responseInJson['value']
        second_get_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(second_get_value, 'bar')
        self.assertEqual(compareCausalMetadata(second_get_causal_metadata, first_get_causal_metadata), True)

        print("\n###################### Getting k1 from replica3 ######################\n")

        # get the value of the new key from replica1 after putting a new key
        baseUrl  = hostBaseUrl + ':' + replica3HostPort + '/key-value-store/k1'
        response = requests.get(baseUrl, json={"causal-metadata":second_get_causal_metadata})
        responseInJson = response.json()

        third_get_value = responseInJson['value']
        third_get_causal_metadata = responseInJson['causal-metadata']
        self.assertEqual(response.status_code, 200)
        self.assertEqual(third_get_value, 'bar')
        self.assertEqual(compareCausalMetadata(third_get_causal_metadata, second_get_causal_metadata), True)

    def setUp(self):
        print("\n###################### Running replicas ######################\n")
        runReplica(replica1HostPort, replica1Ip, subnetName, "replica1")
        runReplica(replica2HostPort, replica2Ip, subnetName, "replica2")
        runReplica(replica3HostPort, replica3Ip, subnetName, "replica3")
        time.sleep(5)

    def tearDown(self):
        print("== Destroying containers..")
        stopAndRemoveInstance("replica1")
        stopAndRemoveInstance("replica2")
        stopAndRemoveInstance("replica3")

    
    @classmethod
    def setUpClass(cls):
        print('= Cleaning-up containers and subnet, possibly leftover from a previous interrupted run..')
        print("\nCleaning up replica1 ...")
        stopAndRemoveInstance("replica1")
        print("Cleaning up replica2 ...")
        stopAndRemoveInstance("replica2")
        print("Cleaning up replica3 ...")
        stopAndRemoveInstance("replica3")
        print("Cleaning up", subnetName, "...")
        print("Done cleaning up.")
        print("###################### Building Docker image ######################\n")
        # build docker image
        buildDockerImage()
        print("\n###################### Creating the subnet ######################\n")
        removeSubnet(subnetName)
        createSubnet(subnetAddress, subnetName)

    
    @classmethod
    def tearDownClass(cls):
        print("Cleaning up", subnetName, "...")
        removeSubnet(subnetName)
        print("Done cleaning up.")

if __name__ == '__main__':
    unittest.main()
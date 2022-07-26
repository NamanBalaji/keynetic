import unittest
import subprocess
import requests
import sys
import time


USAGE = '''
USAGE: python test_suit.py [test-spec]
* The "test-spec" could be 'Part1' or 'Part2' or the name of some test function like 'Part2.test_a_all_running_request_main'.
* For the 'Part1' testsuite you need to run your own container on port 8085.
* For the 'Part2' testsuite, we'll set up three containers and run a set of tests on them.
'''.strip()


class Part1Const:
    hostname = 'localhost'  # Windows and Mac users can change this to the docker vm ip
    portNumber = '8085'
    baseUrl = 'http://' + hostname + ":" + portNumber


class Part1(unittest.TestCase):

    def test_a_get_nonexisting_key(self):
        response = requests.get(Part1Const.baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 404)
        self.assertEqual(responseInJson['doesExist'], False)
        self.assertEqual(responseInJson['message'], 'Error in GET')
        self.assertEqual(responseInJson['error'], 'Key does not exist')

    def test_b_delete_nonexisting_key(self):
        response = requests.delete(Part1Const.baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 404)
        self.assertEqual(responseInJson['doesExist'], False)
        self.assertEqual(responseInJson['message'], 'Error in DELETE')
        self.assertEqual(responseInJson['error'], 'Key does not exist')


    def test_c_put_nonexistent_key(self):
        response = requests.put(Part1Const.baseUrl + '/key-value-store/' + "subject1", json={'value': "Data Structures"})
        responseInJson = response.json()
        self.assertEqual(response.status_code, 201)
        self.assertEqual(responseInJson['message'], 'Added successfully')
        self.assertEqual(responseInJson['replaced'], False)


    def test_d_get_after_put_nonexisting_key(self):
        response = requests.get(Part1Const.baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['doesExist'], True)
        self.assertEqual(responseInJson['message'], 'Retrieved successfully')
        self.assertEqual(responseInJson['value'], 'Data Structures')

    def test_e_put_existent_key(self):
        response = requests.put(Part1Const.baseUrl + '/key-value-store/' + "subject1", json={'value': "Distributed Systems"})
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['message'], 'Updated successfully')
        self.assertEqual(responseInJson['replaced'], True)

    def test_f_get_after_put_existing_key(self):
        response = requests.get(Part1Const.baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['doesExist'], True)
        self.assertEqual(responseInJson['message'], 'Retrieved successfully')
        self.assertEqual(responseInJson['value'], 'Distributed Systems')


    def test_g_delete_existing_key(self):
        response = requests.delete(Part1Const.baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['doesExist'], True)
        self.assertEqual(responseInJson['message'], 'Deleted successfully')

    def test_h_get_after_delete_existing_key(self):
        response = requests.get(Part1Const.baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 404)
        self.assertEqual(responseInJson['doesExist'], False)
        self.assertEqual(responseInJson['message'], 'Error in GET')
        self.assertEqual(responseInJson['error'], 'Key does not exist')

    def test_i_put_key_too_long(self):
        tooLongKey = '6TLxbmwMTN4hX7L0QX5_NflWH0QKfrTlzcuM5PUQHS52___lCizKbEMxLZHhtfww3KcMoboDLjB6mw_wFfEz5v_TtHqvGOZnk4_8aqHga79BaHXzpU9_IRbdjYdQutAU0HEuji6Ny1Ol_MSaBF4JdT0aiG_N7xAkoPH3AlmVqDN45KDGBz7_YHrLnbLEK11SQxZcKXbFomh9JpH_sbqXIaifqOy4g06Ab0q3WkNfVzx7H0hGhNlkINf5PF12'
        value = "haha"
        response = requests.put(Part1Const.baseUrl + '/key-value-store/' + tooLongKey, json={'value': value})

        responseInJson = response.json()
        self.assertEqual(response.status_code, 400)
        self.assertEqual(responseInJson['message'], 'Error in PUT')
        self.assertEqual(responseInJson['error'], 'Key is too long')

    def test_j_put_key_with_no_value(self):
        response = requests.put(Part1Const.baseUrl + '/key-value-store/subject1', json={})
        responseInJson = response.json()
        self.assertEqual(response.status_code, 400)
        self.assertEqual(responseInJson['message'], 'Error in PUT')
        self.assertEqual(responseInJson['error'], 'Value is missing')


class Part2Const:
    subnetName = "mynet"
    subnetAddress = "10.10.0.0/16"

    mainInstanceName = "main-instance"
    forwardingInstance1Name = "forwarding-instance1"
    forwardingInstance2Name = "forwarding-instance2"

    hostname = 'localhost'  # Windows and Mac users can change this to the docker vm ip

    ipAddressMainInstance = "10.10.0.2"
    hostPortMainInstance = "8086"

    ipAddressForwarding1Instance = "10.10.0.3"
    hostPortForwarding1Instance = "8087"

    ipAddressForwarding2Instance = "10.10.0.4"
    hostPortForwarding2Instance = "8088"


# docker linux commands
#
# * all wait for subprocess to complete
# * all but `docker build` suppress stdout
# * all but the removal commands require success
# * the removal commands require success by default but have an argument which can be set to False to ignore failure
# * the removal commands additionally suppress stderr

def removeSubnet(subnetName, required=True):
    command = "docker network rm " + subnetName
    subprocess.run(command, shell=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=required)

def createSubnet(subnetAddress, subnetName):
    command  = "docker network create --subnet=" + subnetAddress + " " + subnetName
    subprocess.check_call(command, stdout=subprocess.DEVNULL, shell=True)

def buildDockerImage():
    command = "docker build -t assignment2-img ."
    subprocess.check_call(command, shell=True)

def runMainInstance(hostPortNumber, ipAddress, subnetName, instanceName):
    command = "docker run -d -p " + hostPortNumber + ":8085 --net=" + subnetName + " --ip=" + ipAddress + " --name=" + instanceName + " assignment2-img"
    subprocess.check_call(command, shell=True, stdout=subprocess.DEVNULL)

def runForwardingInstance(hostPortNumber, ipAddress, subnetName, instanceName, forwardingAddress):
    command = "docker run -d -p " + hostPortNumber + ":8085 --net=" + subnetName  + " --ip=" + ipAddress + " --name=" + instanceName + " -e FORWARDING_ADDRESS=" + forwardingAddress + " assignment2-img"
    subprocess.check_call(command, shell=True, stdout=subprocess.DEVNULL)

def stopAndRemoveInstance(instanceName, required=True):
    stopCommand = "docker stop " + instanceName
    subprocess.run(stopCommand, shell=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=required)
    time.sleep(1)
    removeCommand = "docker rm " + instanceName
    subprocess.run(removeCommand, shell=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=required)


class Part2(unittest.TestCase):

    ######################## Functions to send the required requests ##########################################
    def send_all_requests(self, baseUrl):
        # get nonexistent key
        response = requests.get(baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 404)
        self.assertEqual(responseInJson['doesExist'], False)
        self.assertEqual(responseInJson['message'], 'Error in GET')
        self.assertEqual(responseInJson['error'], 'Key does not exist')

        # delete nonexistent key
        response = requests.delete(baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 404)
        self.assertEqual(responseInJson['doesExist'], False)
        self.assertEqual(responseInJson['message'], 'Error in DELETE')
        self.assertEqual(responseInJson['error'], 'Key does not exist')

        # put nonexistent key
        response = requests.put(baseUrl + '/key-value-store/' + "subject1", json={'value': "Data Structures"})
        responseInJson = response.json()
        self.assertEqual(response.status_code, 201)
        self.assertEqual(responseInJson['message'], 'Added successfully')
        self.assertEqual(responseInJson['replaced'], False)

        # get after putting nonexistent key
        response = requests.get(baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['doesExist'], True)
        self.assertEqual(responseInJson['message'], 'Retrieved successfully')
        self.assertEqual(responseInJson['value'], 'Data Structures')

        # put existent key
        response = requests.put(baseUrl + '/key-value-store/' + "subject1", json={'value': "Distributed Systems"})
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['message'], 'Updated successfully')
        self.assertEqual(responseInJson['replaced'], True)

        # get after putting existent key
        response = requests.get(baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['doesExist'], True)
        self.assertEqual(responseInJson['message'], 'Retrieved successfully')
        self.assertEqual(responseInJson['value'], 'Distributed Systems')

        # delete existent key
        response = requests.delete(baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['doesExist'], True)
        self.assertEqual(responseInJson['message'], 'Deleted successfully')

        # get after deleting key
        response = requests.get(baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 404)
        self.assertEqual(responseInJson['doesExist'], False)
        self.assertEqual(responseInJson['message'], 'Error in GET')
        self.assertEqual(responseInJson['error'], 'Key does not exist')

        # put key with no value
        response = requests.put(baseUrl + '/key-value-store/subject1', json={})
        responseInJson = response.json()
        self.assertEqual(response.status_code, 400)
        self.assertEqual(responseInJson['message'], 'Error in PUT')
        self.assertEqual(responseInJson['error'], 'Value is missing')

    def send_forwarding12_requests(self, baseUrl1, baseUrl2):
        # put nonexistent key
        response = requests.put(baseUrl1 + '/key-value-store/' + "subject1", json={'value': "Data Structures"})
        responseInJson = response.json()
        self.assertEqual(response.status_code, 201)
        self.assertEqual(responseInJson['message'], 'Added successfully')
        self.assertEqual(responseInJson['replaced'], False)

        # get after putting nonexistent key
        response = requests.get(baseUrl2 + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['doesExist'], True)
        self.assertEqual(responseInJson['message'], 'Retrieved successfully')
        self.assertEqual(responseInJson['value'], 'Data Structures')

        # put existent key
        response = requests.put(baseUrl2 + '/key-value-store/' + "subject1", json={'value': "Distributed Systems"})
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['message'], 'Updated successfully')
        self.assertEqual(responseInJson['replaced'], True)

        # get after putting existent key
        response = requests.get(baseUrl1 + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['doesExist'], True)
        self.assertEqual(responseInJson['message'], 'Retrieved successfully')
        self.assertEqual(responseInJson['value'], 'Distributed Systems')

        # delete existent key
        response = requests.delete(baseUrl1 + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(responseInJson['doesExist'], True)
        self.assertEqual(responseInJson['message'], 'Deleted successfully')

        # get after deleting key
        response = requests.get(baseUrl2 + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 404)
        self.assertEqual(responseInJson['doesExist'], False)
        self.assertEqual(responseInJson['message'], 'Error in GET')
        self.assertEqual(responseInJson['error'], 'Key does not exist')

    def send_requests_forwarding_while_main_stopped(self, baseUrl):

        # put nonexistent key
        response = requests.put(baseUrl + '/key-value-store/' + "subject1", json={'value': "Data Structures"})
        responseInJson = response.json()
        self.assertEqual(response.status_code, 503)
        self.assertEqual(responseInJson['message'], 'Error in PUT')
        self.assertEqual(responseInJson['error'], 'Main instance is down')


        # get after putting nonexistent key
        response = requests.get(baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 503)
        self.assertEqual(responseInJson['message'], 'Error in GET')
        self.assertEqual(responseInJson['error'], 'Main instance is down')

        # delete existent key
        response = requests.delete(baseUrl + '/key-value-store/subject1')
        responseInJson = response.json()
        self.assertEqual(response.status_code, 503)
        self.assertEqual(responseInJson['message'], 'Error in DELETE')
        self.assertEqual(responseInJson['error'], 'Main instance is down')

    ########################## Run tests #######################################################
    def test_a_all_running_request_main(self):
        print("=== Sending requests to main instance..")
        baseUrl = 'http://' + Part2Const.hostname + ':' + Part2Const.hostPortMainInstance
        self.send_all_requests(baseUrl)

    def test_b_all_running_request_forwarding1(self):
        print("=== Sending requests the first forwarding instance..")
        baseUrl = 'http://' + Part2Const.hostname + ':' + Part2Const.hostPortForwarding1Instance
        self.send_all_requests(baseUrl)

    def test_c_all_running_request_forwarding2(self):
        print("=== Sending requests to the second forwarding instance..")
        baseUrl = 'http://' + Part2Const.hostname + ':' + Part2Const.hostPortForwarding2Instance
        self.send_all_requests(baseUrl)

    def test_d_all_running_request_forwarding12(self):
        print("=== Sending requests to both forwarding instances..")
        baseUrl1 = 'http://' + Part2Const.hostname + ':' + Part2Const.hostPortForwarding1Instance
        baseUrl2 = 'http://' + Part2Const.hostname + ':' + Part2Const.hostPortForwarding2Instance
        self.send_forwarding12_requests(baseUrl1, baseUrl2)

    def test_e_main_stopped_request_forwarding12(self):
        print("=== Destroying just the main instance..")
        stopAndRemoveInstance(Part2Const.mainInstanceName)

        print("=== Sending requests to both forwarding instances (while main instance is stopped)..")
        baseUrl1 = 'http://' + Part2Const.hostname + ':' + Part2Const.hostPortForwarding1Instance
        baseUrl2 = 'http://' + Part2Const.hostname + ':' + Part2Const.hostPortForwarding2Instance
        self.send_requests_forwarding_while_main_stopped(baseUrl1)
        self.send_requests_forwarding_while_main_stopped(baseUrl2)

        print("=== Running just the main instance..")
        runMainInstance(Part2Const.hostPortMainInstance, Part2Const.ipAddressMainInstance, Part2Const.subnetName, Part2Const.mainInstanceName)

    ############ Test suite lifecycle: setUpClass [setUp test tearDown]* tearDownClass #########
    @classmethod
    def setUpClass(cls):
        print('= Cleaning-up containers and subnet, possibly leftover from a previous interrupted run..')
        stopAndRemoveInstance(Part2Const.mainInstanceName, required=False)
        stopAndRemoveInstance(Part2Const.forwardingInstance1Name, required=False)
        stopAndRemoveInstance(Part2Const.forwardingInstance2Name, required=False)
        # Prevent "Error response from daemon: error while removing network: network assignment2-net id ... has active endpoints"
        time.sleep(5)
        removeSubnet(Part2Const.subnetName, required=False)
        print('= Creating image and subnet..')
        buildDockerImage()
        createSubnet(Part2Const.subnetAddress, Part2Const.subnetName)

    def setUp(self):
        print("\n== Running containers..")
        runMainInstance(Part2Const.hostPortMainInstance, Part2Const.ipAddressMainInstance, Part2Const.subnetName, Part2Const.mainInstanceName)
        runForwardingInstance(Part2Const.hostPortForwarding1Instance, Part2Const.ipAddressForwarding1Instance, Part2Const.subnetName, Part2Const.forwardingInstance1Name, Part2Const.ipAddressMainInstance + ":8085" )
        runForwardingInstance(Part2Const.hostPortForwarding2Instance, Part2Const.ipAddressForwarding2Instance, Part2Const.subnetName, Part2Const.forwardingInstance2Name, Part2Const.ipAddressMainInstance + ":8085" )
        # Ensure that processes in the containers have time to start up and listen on ports
        time.sleep(5)

    def tearDown(self):
        print("== Destroying containers..")
        stopAndRemoveInstance(Part2Const.mainInstanceName)
        stopAndRemoveInstance(Part2Const.forwardingInstance1Name)
        stopAndRemoveInstance(Part2Const.forwardingInstance2Name)

    @classmethod
    def tearDownClass(cls):
        print('= Destroying subnet..')
        removeSubnet(Part2Const.subnetName)


if __name__ == '__main__':
    if len(sys.argv) == 1:
        print(USAGE)
        sys.exit(1)
    else:
        unittest.main(verbosity=2)
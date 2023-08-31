# Cryptographic-Computing
(The pictures can be navigated to with CTRL + left click)

## Handin1 ##
We made the eight bloodtypes correspond with a 8-bit integer value from the perspective of the recipient. E.g. which donor bloodtypes are compatible with the recipient. The bloodtype AB+ is 0 due to being able to recieve blood from everyone else, where O- is 7. since it is only able to recieve blood from its own bloodtype.

The above setup made it possible to use the bloodtype numbers for both tasks as we elaborate on further below. Also for both tasks, see the code for reference.


# Task1 #
To solve the task we used the blood type chart seen below to create a lookup table representing whether a "recipient" is compatible with a "donor". We used the eight bloodtypes values correspond with a 8-bit integers described previously to index in the lookup table.

![Image description](./Pictures/blood-type-compatibility.png)
(from MediFee.com)

In the handin1_test.go file we made a couple of tests to validate that the lookup table worked had the correct indexes.

# Task 2 #
To create a Boolean formula, we use the 8 bloodtypes corresponding to the numbers 0 to 7 represented as 8-bit intergers. Since the numbers are represented as 8-bit the first 5 digits are 0's. Therefore, we only need to compare the last 3 bits. 

The Boolean formula we came up with is illustrated below as a Boolean circuit:
![Image description](./Pictures/Boolean-circuit-handin1.png)





In the main function we loop through all bloodtype combinations to see that the Boolean formula and lookup function provide the same output. There is also a test that does the same.


# Instructions #
To run the code, navigate to the folder Cryptographic-Computing and run the following command on setup:

"go mod init CryptographicComputing"
    - Notice that the "-" is not recognized and therefore should be leftout.



To run our main function to see that the Boolean formula and lookup function provide the same output, run the command:

"go run handin1.go" 

To run our tests (all files that end with "..._test.go") run the command:
"go test"

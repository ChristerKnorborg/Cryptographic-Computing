# Instructions #
To run the code, navigate to the folder Cryptographic-Computing. If the GO modules are having issues, you might want to write "go mod tidy" or even "go mod init" anew:

Afterwards, navigate to the folder called "handin2". To run our tests (handin2_test.go") run the command:

"go test". 

There is also a main function in the root folder, but it does the exact same as one of the tests,
by checking all possible BloodType Combinations gives the same result as the lookup table. You can run this from the Cryptographic-Computing folder with the command:

"go run main.go"


# Description #
To solve the assignment, we made files with structs to represent the different entities of the protocol: Alice, Bob and the Dealer. The overall protocol can be seen in the "handin2.go"-file in the ComputeOTTTBloodTypeCompatability() function. To see the calculation of the different entities, navigate to the respective files. Furthermore, there is a utility file called "utils.go" with functionalities used throughout the other files - for example XOR methods, random numbers generations (though not cryptographically secure), generation and shifting of the matrixes, and functionality from the previous assignment. The code should be well documented, hence it should be possible to follow the flow of the protocol by reading the code.

# Instructions #
It is important that you open the inner most Cryptographic-Computing-folder as your root in your editor. Otherwise, Go packages get confused. Be wary of this, as the ZIP-extraction might create another folder encapsulating the inner folder. If the GO modules are having issues, you might want to write "go mod tidy" or even "go mod init" anew.

To run the code, navigate to the root folder. Here there is a simple test that try all combinations of blood types for the depth d Homomorphic Encryption protocol (the method is called HomomorphicBloodtypeEncryption() ) and compare the solution with the lookup table. To run it use the command:

"go run main.go"


We also have a testing file that can be run with the "go test" command. However, the main method covers all of the tests. To run the test, you must in the folder where the testing file is located.

# Description #
To solve the assignment, we made files with structs to represent the different entities of the protocol: Alice, Bob and the homomorphic functions. The last mentioned is responsible for all the different algorithms mentioned in the lecture notes such as key generation, encryption and decryption and all other operations related to these. The entire setup of the protocol can be see in the "handin6.go"-file, where it is possible to navigate to the respective files and view the implementation of operations with comments severing as documentation.

# On security parameters #
For the different security parameters of the crypto system we chose p to be of length 512, meaning we have 512 bit length of security. We also chose n = 10.
Lastly we chose the random subset S to be of size n/2 with each value in the set {1, ... , n} being uniformly distributed. This guarentees a high level of security with regards to encryption but if multiple homopormhic operations are performed on the same input, the noise might become and issue when decrypting.
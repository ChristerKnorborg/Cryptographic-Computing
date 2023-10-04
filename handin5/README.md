# Instructions #
It is important that you open the inner most Cryptographic-Computing-folder as your root in your editor. Otherwise, Go packages get confused. Be wary of this, as the ZIP-extraction might create another folder encapsulating the inner folder. If the GO modules are having issues, you might want to write "go mod tidy" or even "go mod init" anew.

To run the code, navigate to the root folder. Here there is a simple test that try all combinations of blood types for the Garbled Circuit protocol (the method is called GarbledCircuit) and compare the solution with the lookup table. To run it use the command:

"go run main.go"


We also have a testing file that can be run with the "go test" command. However, the main method covers all of the tests.

# Description #
To solve the assignment, we made files with structs to represent the different entities of the protocol: Alice, Bob and the ElGamal Cryptosystem from last assignment, which is also used in the Garbled Curcuit for the OT protocol. The overall protocol can be seen in the "handin5.go"-file in the GarbledCircuit() function. To see the calculation of the different entities, navigate to the respective files. Furthermore, there are utility files called "GarbledFunctions.go" and "utils.go" with functionalities used throughout the other files. The code should be very well documented, hence it should be possible to follow the flow of the protocol by reading the code.


# Note on ElGamal #
Last week we had encountered a problem with elGamal. It turns out that the reason it did not work was due to an modulo error which we have fixed for this assignment and hence there is no need for fixed values in our parameter generation of primes q and p, and group g. However, we have set the size of the prime numbers to be rather small (even though real life security requires over 2000 bit len), since the garbledCircuit function that is called 8x8 = 64 times to try all bloodtype combination in the main method have to find these in every iteration.
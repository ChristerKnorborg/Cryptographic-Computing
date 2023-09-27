# Instructions #
It is important that you open the inner most Cryptographic-Computing-folder as your root in your editor. Otherwise, Go packages get confused. Be wary of this, as the ZIP-extraction might create another folder encapsulating the inner folder. If the GO modules are having issues, you might want to write "go mod tidy" or even "go mod init" anew.

To run the code, navigate to the root folder. Here there is a simple test that try all combinations of blood types for the Oblivious Transfer protocol (the method is called ObliviousTransfer) and compare the solution with the lookup table. To run it use the command:

"go run main.go"

One thing to notice is that we the ElGamal parameters have to be instatiated for each of the 64 bloodtype iterations. Hence it takes quite some time to run the main method test. Due to this fact, we also choose to not create any addition test methods.

# Description #
To solve the assignment, we made files with structs to represent the different entities of the protocol: Alice, Bob and the ElGamal Cryptosystem. The overall protocol can be seen in the "handin4.go"-file in the ObliviousTransfer() function. To see the calculation of the different entities, navigate to the respective files. Furthermore, there is a utility file called "utils.go" with functionalities used throughout the other files - for example the Bloodtype struct definition that we made. The code should be well documented, hence it should be possible to follow the flow of the protocol by reading the code.


# Issue in the assignment with regard to prime generation#
In this assignment we had some issue with regards to the generation of the random parameter q for the ElGamal cryptosystem. In particular, the random prime generation method from the GOLANG crypto library: rand.Prime(rand.Reader, BITLEN) gave us problems. It appears that the method some times incorrectly yields a non-prime number (probably due to the Miller-rabin primality testing). Therefore, we would often get a incorrect results when running the rand.Prime(rand.Reader, BITLEN) with a low BITLEN such as 256. In the other hand, a BITLEN around 2000 (as is required for security with the current computing power) would result in a hour to run all 64 bloodtype compinations. Therefore, we ended up using a hardcoded prime number found online for the parameter q. The difference between the hardcoded method and the randomly generated one can be seen in the elGamal methods "InitFixedQ()" and "Init()" respectively. Notice, that we still find the rest of the parameters as explained in the lecture notes. For more concrete example of this, look in the code comments. 

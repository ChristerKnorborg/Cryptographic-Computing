# Instructions #
It is important that you open the inner most Cryptographic-Computing-folder as your root in your editor. Otherwise, Go packages get confused. Be wary of this, as the ZIP-extraction might create another folder encapsulating the inner folder. If the GO modules are having issues, you might want to write "go mod tidy" or even "go mod init" anew.

To run the code, navigate to the root folder. Here there is a simple test that try all combinations of blood types for the BeDOZa protocol (the method is called ComputeBeDOZaBloodTypeCompatability) and compare the solution with the lookup table. To run it use the command:


"go run main.go"


We spend a great deal creating automated tests. You can run these by navigating to the folder called "handin3". To run the tests in the test file "handin3_test.go" run the command:

"go test". 

The main method testing is also included in an automated test.


# Description #
To solve the assignment, we made files with structs to represent the different entities of the protocol: Alice, Bob and the Dealer. The overall protocol can be seen in the "handin3.go"-file in the ComputeOTTTBloodTypeCompatability() function. To see the calculation of the different entities, navigate to the respective files. Furthermore, there is a utility file called "utils.go" with functionalities used throughout the other files - for example the Bloodtype struct definition that we made. The code should be well documented, hence it should be possible to follow the flow of the protocol by reading the code.

We used the boolean logic formula provided to all students after assignment 1. 
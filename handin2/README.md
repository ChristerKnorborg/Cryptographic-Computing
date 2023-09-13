# Cryptographic-Computing
(The pictures can be navigated to with CTRL + left click)

# Instructions #
To run the code, navigate to the folder Cryptographic-Computing and run the following command on setup:

"go mod init CryptographicComputing"
    - Notice that the "-" is not recognized and therefore should be leftout.

Afterwards, navigate to the folder called "handin2"

To run our main function to see that the Boolean formula and lookup function provide the same output, run the command:

"go run handin2.go" 
    - We commented out the main function due to issues when you use the the first time "go mod init CryptographicComputing". Therefore, simply remove the commented main function.

To run our tests (all files that end with "..._test.go") run the command:
"go test"

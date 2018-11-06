package main

import (
	"context" //AWS Lambda Context Object in Go
	"fmt"

	"github.com/aws/aws-lambda-go/lambda" //Implements the Lambda programming model for Go
)

type MyEvent struct {
	Name string `json:"name"`
}

/*
	name: HandleRequest()
	description: The Lambda handler signature, includes the code which will be executed
	parameters:
		ctx context.Context: Provides runtime information for your Lambda function invocation.
			ctx: The variable declared to leverage the information avaiable via AWS Lambda Context Object in Go
		name MyEvent: An input type with a variable name of "name" whose value will be returned in the "return" statement
	return:
		string, error: Returns standard error information
*/
func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	//Simply returns a formattted "Hello" greeting with the name supplied in the handler signature
	return fmt.Sprintf("Hello %s!", name.Name), nil //"nil" indicates there were no errors and the function executed successfully
}

/*
	name: main()
	description: The entry point that executes the Lambda function code
*/
func main() {
	//Executes the Lambda function
	lambda.Start(HandleRequest)
}

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	//"time"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var projectIndexStr = "_projectindex" //name for the key that will store list of project index/project name
var employeeIndexStr = "_employeeindex" //name for the key that will store list of employee index/employeeID


//same as employee
type Member struct{
	MemberID string `json:"memberid"`
	MemberName string `json:"membername"`
    JobTitle string `json:"jobtitle"`
    Level int `json:"level"`
    JobGroup string `json:"jobgroup"`
}

type Project struct{
	Name string `json:"name"`
    Members []string `json:"members"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	
	if function == "write" {											    //writes a value to the chaincode state
		return t.Write(stub, args)
	} else if function == "add_employee"{									//add new employee
        return t.add_employee(stub, args)
    } else if function == "update_employee" {								//update attributes of existing employee
		return t.update_employee(stub, args)
	} else if function == "create_project"{									//create new project
		return t.create_project(stub, args)
	} else if function == "add_project_member"{								//add new member to the project
		return t.add_project_member(stub, args)
	} else if function == "delete_project_member"{							//delete member from a project
		return t.delete_project_member(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {													//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}

// ===========================================================================================================================
// Add new employee
// ===========================================================================================================================
func (t *SimpleChaincode) add_employee(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	//   0       1       2           3          4
	// "id", "name", "job title", "level", "job group"
	if len(args) != 5 {
        fmt.Println("Incorrect number of arguments. Expecting 5")
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	//input sanitation
	fmt.Println("- start add new employee")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
    if len(args[4]) <= 0 {
		return nil, errors.New("5th argument must be a non-empty string")
	}

	//Get the employee from chaincode state
	employeeAsBytes, err := stub.GetState(args[0])
	if err != nil {
        fmt.Println("Failed to get employee")
		return nil, errors.New("Failed to get employee")
	}

	employee := Member{}
	json.Unmarshal(employeeAsBytes, &employee) //equals to JSON.parse in javascript

	//check if employee already exists
	if employee.MemberID == args[0] {
		fmt.Println("This employee arleady exists: " + args[0])
		fmt.Println(employee);
		return nil, errors.New("This employee arleady exists")
	}
	
    employee.MemberID = args[0]
    employee.MemberName = args[1]
    employee.JobTitle = args[2]
    employee.Level, err = strconv.Atoi(args[3])

    if err != nil {
        return nil, errors.New("Level must be numeric")
    }

    employee.JobGroup = args[4]

    employeeAsBytes, _ = json.Marshal(employee) //equals to JSON.stringify in javascript

	err = stub.PutState(args[0], employeeAsBytes) //write the new employee to the chaincode state
	if err != nil {
		return nil, err
	}
		
	//get the employee index
	employeeIndexAsBytes, err := stub.GetState(employeeIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get employee index")
	}

	var employeeIndex []string
	json.Unmarshal(employeeIndexAsBytes, &employeeIndex)							//un stringify it aka JSON.parse()
	
	//append
	employeeIndex = append(employeeIndex, args[0])									//add employeeID to index list
	fmt.Println("! employee index: ", employeeIndex)
	jsonAsBytes, _ := json.Marshal(employeeIndex)
	err = stub.PutState(employeeIndexStr, jsonAsBytes)						//rewrite employee index to chaincode state

	fmt.Println("- end add employee")
	return nil, nil
}

// ===========================================================================================================================
// Update employee
// ===========================================================================================================================
func (t *SimpleChaincode) update_employee(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	//   0       1       2           3          4
	// "id", "name", "job title", "level", "job group"
	if len(args) != 5 {
        fmt.Println("Incorrect number of arguments. Expecting 5")
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	//input sanitation
	fmt.Println("- start update employee")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
    if len(args[4]) <= 0 {
		return nil, errors.New("5th argument must be a non-empty string")
	}

	//check if employee is exists
	employeeAsBytes, err := stub.GetState(args[0])	//get employee detail from chaincode state
	if err != nil {
        fmt.Println("Failed to get employee")
		return nil, errors.New("Failed to get employee")
	}

	employee := Member{}
	json.Unmarshal(employeeAsBytes, &employee)
	
	//Update the employee details
    employee.MemberID = args[0]
    employee.MemberName = args[1]
    employee.JobTitle = args[2]
    employee.Level, err = strconv.Atoi(args[3])

    if err != nil {
        return nil, errors.New("Level must be numeric")
    }

    employee.JobGroup = args[4]

    employeeAsBytes, _ = json.Marshal(employee)	//unstringify a.k.a JSON.stringify()

	err = stub.PutState(args[0], employeeAsBytes)	//rewrite the employee to chaincode state
	if err != nil {
		return nil, err
	}
		
	if err != nil {
		return nil, errors.New("Failed to get employee index")
	}

	fmt.Println("- end update")
	return nil, nil
}

// ===========================================================================================================================
// Create Project
// ===========================================================================================================================
func (t *SimpleChaincode) create_project(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	if len(args) != 1 {
		fmt.Println("Incorrect number of arguments. Expecting 1")
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	//input sanitation
	fmt.Println("- start create project")
	if len(args[0]) <= 0 {
		fmt.Println("1st argument must be a non-empty string")
		return nil, errors.New("1st argument must be a non-empty string")
	}

	name := args[0] //of course can contain white space :D, beacuse it's name and will be stored as a value not a key

	//Get project details from chaincode state
	projectAsBytes, err := stub.GetState(strings.Replace(name, " ", "_", -1))	//String replace is used for get rid of white space, because key can contain white space
	if err != nil {
		fmt.Println("Failed to get project name")
		return nil, errors.New("Failed to get project name")
	}

	res := Project{}
	json.Unmarshal(projectAsBytes, &res)	//equals to JSON.parse()

	//check if project already exists
	if res.Name == name{
		fmt.Println("This project arleady exists: " + name)
		fmt.Println(res);
		return nil, errors.New("This project arleady exists")
	}

	res.Name = name
	
	jsonAsBytes, _ := json.Marshal(res)
	//write project to the chaincode state
	err = stub.PutState(strings.Replace(name, " ", "_", -1), jsonAsBytes)	//String replace is used for get rid of white space, because key can contain white space
	if err != nil {
		return nil, err
	}
		
	//get the project index
	projectAsBytes, err = stub.GetState(projectIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get project index")
	}

	var projectIndex []string
	json.Unmarshal(projectAsBytes, &projectIndex)							//un stringify it aka JSON.parse()
	
	//append
	projectIndex = append(projectIndex, strings.Replace(name, " ", "_", -1))	//add project name to index list, but before it remove the white space first. Remember key can't contain white space
	fmt.Println("! project index: ", projectIndex)
	jsonAsBytes, _ = json.Marshal(projectIndex)
	err = stub.PutState(projectIndexStr, jsonAsBytes)						//rewrite project index

	fmt.Println("- end create project")
	return nil, nil
}

// ===========================================================================================================================
// Add project member
// ===========================================================================================================================
func (t *SimpleChaincode) add_project_member(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error
	var isExists int //1 means member already exists in that project,  0 is otherwise

	//   0              1            2      ..........
	//projectName   "memberID"  "memberID"  ..........
	if len(args) < 2 {
		fmt.Println("Incorrect number of arguments. Expecting 2 or more")
		return nil, errors.New("Incorrect number of arguments. Expecting 2 or more")
	}
	
	fmt.Println("- start add project member")
	fmt.Println(args[0] + " - " + args[1])

	//Get project from chaincode state
	projectAsBytes, err := stub.GetState(args[0])

	if err != nil{
		fmt.Println("Failed to get project")
		return nil, errors.New("Failed to get project")
	}

	project := Project{}
	json.Unmarshal(projectAsBytes, &project); //JSON.parse()

	for i:=1; i < len(args); i++ {
		isExists = 0 //0 means member still not in this project

		if len(project.Members) == 0 {
			project.Members = append(project.Members, args[i])	//append memberID/employeeID to project members array 
			fmt.Println("! Success add new member: " + args[i])
		}

		for j:= range project.Members{
			if args[i] == project.Members[j] {
				isExists = 1 //1 means member already exists in this project
				break
			}
		}

		if isExists == 0 {
			project.Members = append(project.Members, args[i])	//append memberID/employeeID to project members array 
			fmt.Println("! Success add new member: " + args[i])
		}
	}

	jsonAsBytes, _ := json.Marshal(project)	//equals to JSON.stringify
	err = stub.PutState(args[0], jsonAsBytes)	//rewrite project to the chaincode state

	if err != nil {
		return nil, err
	}
	
	fmt.Println("- end add new member")
	return nil, nil
}

// ===========================================================================================================================
// Delete project member
// ===========================================================================================================================
func (t *SimpleChaincode) delete_project_member(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//   0                  1
	// "project name", "member id"
	if len(args) != 2 {
		fmt.Println("Incorrect number of arguments. Expecting 2")
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	//get project from chaincode state
	projectAsBytes, err := stub.GetState(args[0]);

	if err != nil{
		fmt.Println("Failed to get project")
		return nil, errors.New("Failed to get project")
	}

	project := Project{}
	json.Unmarshal(projectAsBytes, &project)	//equals to JSON.parse()
	
	//remove member from project
	for i := range project.Members{
		//looking for member ID
		if project.Members[i] == args[1]{
			fmt.Println("member found")
			project.Members = append(project.Members[:i], project.Members[i+1:]...)			//remove it
			break
		}
	}

	projectAsBytes, _ = json.Marshal(project)	//stringify
	err = stub.PutState(args[0], projectAsBytes)	//rewrite project to the chaincode state

	if err != nil {
		return nil, errors.New("Failed to delete member from project chaincode state")
	}

	return nil, nil
}

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) Write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]															//rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func installChaincode(name string, version string, sequence int, channel string, org string, peerPort int, path string) {
	//act as given org
	cmd := exec.Command("pwd")
	pwd, _ := cmd.Output()

	fmt.Printf("Switching to %s\n", org)
	os.Setenv("CORE_PEER_TLS_ENABLED", "true")
	os.Setenv("CORE_PEER_LOCALMSPID", fmt.Sprintf("%sMSP", org))
	os.Setenv("CORE_PEER_TLS_ROOTCERT_FILE", fmt.Sprintf("%s/organizations/peerOrganizations/%s.example.com/peers/peer0.%s.example.com/tls/ca.crt", strings.TrimSuffix(string(pwd), "\n"), strings.ToLower(org), strings.ToLower(org)))
	os.Setenv("CORE_PEER_MSPCONFIGPATH", fmt.Sprintf("%s/organizations/peerOrganizations/%s.example.com/users/Admin@%s.example.com/msp", strings.TrimSuffix(string(pwd), "\n"), strings.ToLower(org), strings.ToLower(org)))
	os.Setenv("CORE_PEER_ADDRESS", fmt.Sprintf("localhost:%d", peerPort))

	//package chaincode
	//peer lifecycle chaincode package <name>.tar.gz --path path/to/chaincode-in-go/ --lang golang --label <name>_1.0
	cmd = exec.Command("peer", "lifecycle", "chaincode", "package", fmt.Sprintf("%s.tar.gz", name),
		"--path", path, "--lang", "golang", "--label", fmt.Sprintf("%s_%s", name, version))

	fmt.Printf("packaging chaincode\n")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf(err.Error())
		return
	} else {
		fmt.Printf(string(stdout))
	}

	//install chaincode
	//peer lifecycle chaincode install <name>.tar.gz
	cmd = exec.Command("peer", "lifecycle", "chaincode", "install", fmt.Sprintf("%s.tar.gz", name))

	fmt.Printf("installing chaincode\n")
	stdout, err = cmd.Output()
	if err != nil {
		fmt.Printf(err.Error())
		return
	} else {
		fmt.Printf(string(stdout))
	}

	//get package id
	//run the following command and copy the package id from the output:
	//peer lifecycle chaincode queryinstalled
	//export CC_PACKAGE_ID=<PACKAGE_ID>
	cmd = exec.Command("bash", "-c", fmt.Sprintf("peer lifecycle chaincode queryinstalled | grep %s_%s", name, version))

	fmt.Printf("getting package ID\n")
	stdout, err = cmd.Output()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	stdoutSplit := strings.Split(string(stdout), " ")
	packageID := strings.TrimSuffix(stdoutSplit[2], ",")
	fmt.Println(packageID)

	os.Setenv("CC_PACKAGE_ID", packageID)

	//approve chaincode definition
	//peer lifecycle chaincode approveformyorg -o localhost:7050
	//--ordererTLSHostnameOverride orderer.example.com --channelID  l2 --name <name>
	//--version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls
	//--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
	cmd = exec.Command("peer", "lifecycle", "chaincode", "approveformyorg", "-o", "localhost:7050",
		"--ordererTLSHostnameOverride", "orderer.example.com", "--channelID", channel, "--name", name,
		"--version", version, "--package-id", packageID, "--sequence", fmt.Sprint(sequence),
		"--tls", "--cafile",
		fmt.Sprintf("%s/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem", strings.TrimSuffix(string(pwd), "\n")))

	fmt.Printf("approving chaincode definition\n")
	stdout, err = cmd.Output()
	if err != nil {
		fmt.Printf(err.Error())
		return
	} else {
		fmt.Printf(string(stdout))
	}

	//commit chaincode definition
	//peer lifecycle chaincode commit -o localhost:7050
	//--ordererTLSHostnameOverride orderer.example.com --channelID l2 --name <name>
	//--version 1.0 --sequence 1 --tls
	//--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
	//--peerAddresses localhost:<PEER_PORT> --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt"
	cmd = exec.Command("peer", "lifecycle", "chaincode", "commit", "-o", "localhost:7050",
		"--ordererTLSHostnameOverride", "orderer.example.com", "--channelID", channel, "--name", name,
		"--version", version, "--sequence", fmt.Sprint(sequence),
		"--tls", "--cafile",
		fmt.Sprintf("%s/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem", strings.TrimSuffix(string(pwd), "\n")),
		"--peerAddresses", fmt.Sprintf("localhost:%d", peerPort),
		"--tlsRootCertFiles", fmt.Sprintf("%s/organizations/peerOrganizations/%s.example.com/peers/peer0.%s.example.com/tls/ca.crt", strings.TrimSuffix(string(pwd), "\n"), strings.ToLower(org), strings.ToLower(org)))

	fmt.Printf("committing chaincode defintion\n")
	stdout, err = cmd.Output()
	if err != nil {
		fmt.Printf(err.Error())
		return
	} else {
		fmt.Printf(string(stdout))
	}
}

func main() {
	//TODO: refactor to use command-line subcommands

	nameFlag := flag.String("name", "basic", "specify name of your contract")
	versionFlag := flag.String("version", "1.0", "version identifier in the format x.x, e.g. 1.0")
	sequenceFlag := flag.Int("sequence", 1, "version sequence integer")
	channelFlag := flag.String("channel", "l2", "which channel to install on")
	orgFlag := flag.String("org", "Org3", "organization that you wish to act as, e.g. Org3")
	peerPortFlag := flag.Int("peerPort", 11051, "port of the peer container you want to interact with")

	flag.Parse()

	args := flag.Args()
	if args[0] == "install-chaincode" {
		installChaincode(*nameFlag, *versionFlag, *sequenceFlag, *channelFlag, *orgFlag, *peerPortFlag, args[1])
	}
}

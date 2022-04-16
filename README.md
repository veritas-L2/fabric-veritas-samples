[//]: # (SPDX-License-Identifier: CC-BY-4.0)

# Hyperledger Fabric Samples

[![Build Status](https://dev.azure.com/Hyperledger/Fabric-Samples/_apis/build/status/Fabric-Samples?branchName=main)](https://dev.azure.com/Hyperledger/Fabric-Samples/_build/latest?definitionId=28&branchName=main)

You can use Fabric samples to get started working with Hyperledger Fabric, explore important Fabric features, and learn how to build applications that can interact with blockchain networks using the Fabric SDKs. To learn more about Hyperledger Fabric, visit the [Fabric documentation](https://hyperledger-fabric.readthedocs.io/en/latest).

## Getting started with the Fabric samples

To use the Fabric samples, you need to download the Fabric Docker images and the Fabric CLI tools. First, make sure that you have installed all of the [Fabric prerequisites](https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html). You can then follow the instructions to [Install the Fabric Samples, Binaries, and Docker Images](https://hyperledger-fabric.readthedocs.io/en/latest/install.html) in the Fabric documentation. In addition to downloading the Fabric images and tool binaries, the Fabric samples will also be cloned to your local machine.

## Setup Veritas Test network

After you have setup all the HLF requirements, you can proceed to setup a Veritas Rollup TestNet. 

### Step 1: Setup Layer 1

Bring up two peer nodes and an orderer node.

```bash
cd test-network
./network.sh up
```

Now create a channel, l1, and make the two peer nodes join it. 

```bash
./network.sh createChannel -c l1
```

### Step 2: Setup Layer 2

Now we need to setup an Org3 that will contain a sequencer peer node. First generate certs for Org3, and create an l2 channel.

```bash
cd test-network
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx

cd addOrg3
./addOrg3.sh generate

cd ..
configtxgen -profile L2Genesis -outputBlock ./channel-artifacts/l2.block -channelID l2
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key
osnadmin channel join --channelID l2 --config-block ./channel-artifacts/l2.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
```

Now bring up the Org3 peer and make it join the l2 channel.

```bash
cd addOrg3
./addOrg3.sh up -c l2 
```

### Step 3: Install the State Contract

The state contract manages the world state for the layer 2 rollup. Hence, it needs to be installed on the Org3 sequencer node. You can get the state contract from this [repo](https://github.com/veritas-L2/state-contract). 

After cloning the state contract repo, setup it's dependencies:

```bash
cd state-contract

#update the deps: HLF is unable to install this contract without this step at the moment. 
go get -u

#vendor deps
GO111MODULE=on go mod vendor
```

Now go back to the `test-network` directory and install the contract chaincode.

```bash
cd test-network/
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

#package the chaincode:
peer lifecycle chaincode package <name>.tar.gz --path path/to/chaincode-in-go/ --lang golang --label <name>_1.0

#act as the Org3 peer node:
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org3MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
export CORE_PEER_ADDRESS=localhost:<PEER_PORT>

#install chaincode
peer lifecycle chaincode install <name>.tar.gz


#Approve chaincode definition:
#run the following command and copy the package id from the output:
peer lifecycle chaincode queryinstalled

#then:
export CC_PACKAGE_ID=<PACKAGE_ID>

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID  l2 --name <name> --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

#Commit chaincode definition
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID l2 --name <name> --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:<PEER_PORT> --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt" 
```

At this point, the state contract should be ready to receive invocations from applications and other chaincode on the l2 channel. 

To install any chaincode on layer 2, follow the same chaincode installation steps described above.

## Asset transfer samples and tutorials

The asset transfer series provides a series of sample smart contracts and applications to demonstrate how to store and transfer assets using Hyperledger Fabric.
Each sample and associated tutorial in the series demonstrates a different core capability in Hyperledger Fabric. The **Basic** sample provides an introduction on how
to write smart contracts and how to interact with a Fabric network using the Fabric SDKs. The **Ledger queries**, **Private data**, and **State-based endorsement**
samples demonstrate these additional capabilities. Finally, the **Secured agreement** sample demonstrates how to bring all the capabilities together to securely
transfer an asset in a more realistic transfer scenario.

|  **Smart Contract** | **Description** | **Tutorial** | **Smart contract languages** | **Application languages** |
| -----------|------------------------------|----------|---------|---------|
| [Basic](asset-transfer-basic) | The Basic sample smart contract that allows you to create and transfer an asset by putting data on the ledger and retrieving it. This sample is recommended for new Fabric users. | [Writing your first application](https://hyperledger-fabric.readthedocs.io/en/latest/write_first_app.html) | Go, JavaScript, TypeScript, Java | Go, JavaScript, TypeScript, Java |
| [Ledger queries](asset-transfer-ledger-queries) | The ledger queries sample demonstrates range queries and transaction updates using range queries (applicable for both LevelDB and CouchDB state databases), and how to deploy an index with your chaincode to support JSON queries (applicable for CouchDB state database only). | [Using CouchDB](https://hyperledger-fabric.readthedocs.io/en/latest/couchdb_tutorial.html) | Go, JavaScript | Java, JavaScript |
| [Private data](asset-transfer-private-data) | This sample demonstrates the use of private data collections, how to manage private data collections with the chaincode lifecycle, and how the private data hash can be used to verify private data on the ledger. It also demonstrates how to control asset updates and transfers using client-based ownership and access control. | [Using Private Data](https://hyperledger-fabric.readthedocs.io/en/latest/private_data_tutorial.html) | Go, Java | JavaScript |
| [State-Based Endorsement](asset-transfer-sbe) | This sample demonstrates how to override the chaincode-level endorsement policy to set endorsement policies at the key-level (data/asset level). | [Using State-based endorsement](https://github.com/hyperledger/fabric-samples/tree/main/asset-transfer-sbe) | Java, TypeScript | JavaScript |
| [Secured agreement](asset-transfer-secured-agreement) | Smart contract that uses implicit private data collections, state-based endorsement, and organization-based ownership and access control to keep data private and securely transfer an asset with the consent of both the current owner and buyer. | [Secured asset transfer](https://hyperledger-fabric.readthedocs.io/en/latest/secured_asset_transfer/secured_private_asset_transfer_tutorial.html)  | Go | JavaScript |
| [Events](asset-transfer-events) | The events sample demonstrates how smart contracts can emit events that are read by the applications interacting with the network. | [README](asset-transfer-events/README.md)  | JavaScript, Java | JavaScript |
| [Attribute-based access control](asset-transfer-abac) | Demonstrates the use of attribute and identity based access control using a simple asset transfer scenario | [README](asset-transfer-abac/README.md)  | Go | None |



## Additional samples

Additional samples demonstrate various Fabric use cases and application patterns.

|  **Sample** | **Description** | **Documentation** |
| -------------|------------------------------|------------------|
| [Commercial paper](commercial-paper) | Explore a use case and detailed application development tutorial in which two organizations use a blockchain network to trade commercial paper. | [Commercial paper tutorial](https://hyperledger-fabric.readthedocs.io/en/latest/tutorial/commercial_paper.html) |
| [Off chain data](off_chain_data) | Learn how to use the Peer channel-based event services to build an off-chain database for reporting and analytics. | [Peer channel-based event services](https://hyperledger-fabric.readthedocs.io/en/latest/peer_event_services.html) |
| [Token ERC-20](token-erc-20) | Smart contract demonstrating how to create and transfer fungible tokens using an account-based model. | [README](token-erc-20/README.md) |
| [Token UTXO](token-utxo) | Smart contract demonstrating how to create and transfer fungible tokens using a UTXO (unspent transaction output) model. | [README](token-utxo/README.md) |
| [Token ERC-1155](token-erc-1155) | Smart contract demonstrating how to create and transfer multiple tokens (both fungible and non-fungible) using an account based model. | [README](token-erc-1155/README.md) |
| [Token ERC-721](token-erc-721) | Smart contract demonstrating how to create and transfer non-fungible tokens using an account-based model. | [README](token-erc-721/README.md) |
| [High throughput](high-throughput) | Learn how you can design your smart contract to avoid transaction collisions in high volume environments. | [README](high-throughput/README.md) |
| [Simple Auction](auction-simple) | Run an auction where bids are kept private until the auction is closed, after which users can reveal their bid. | [README](auction-simple/README.md) |
| [Dutch Auction](auction-dutch) | Run an auction in which multiple items of the same type can be sold to more than one buyer. This example also includes the ability to add an auditor organization. | [README](auction-dutch/README.md) |
| [Chaincode](chaincode) | A set of other sample smart contracts, many of which were used in tutorials prior to the asset transfer sample series. | |
| [Interest rate swaps](interest_rate_swaps) | **Deprecated in favor of state based endorsement asset transfer sample** | |
| [Fabcar](fabcar) | **Deprecated in favor of basic asset transfer sample** |  |

## License <a name="license"></a>

Hyperledger Project source code files are made available under the Apache
License, Version 2.0 (Apache-2.0), located in the [LICENSE](LICENSE) file.
Hyperledger Project documentation files are made available under the Creative
Commons Attribution 4.0 International License (CC-BY-4.0), available at http://creativecommons.org/licenses/by/4.0/.

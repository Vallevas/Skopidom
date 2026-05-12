const { Web3 } = require('web3');
const fs = require('fs');
const path = require('path');

// Configuration
const GANACHE_URL = process.env.GANACHE_URL || 'http://127.0.0.1:7545';
const MNEMONIC =
  process.env.MNEMONIC ||
  'test test test test test test test test test test test junk';

// Default account (first account from Ganache default mnemonic)
const DEFAULT_PRIVATE_KEY =
  '0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80';

// Helper to handle BigInt serialization
function replacer(key, value) {
  if (typeof value === 'bigint') {
    return value.toString();
  }
  return value;
}

async function deploy() {
  console.log('Connecting to Ganache at', GANACHE_URL);

  const web3 = new Web3(GANACHE_URL);

  // Test connection
  try {
    const blockNumber = await web3.eth.getBlockNumber();
    console.log('Connected! Current block:', blockNumber.toString());
  } catch (error) {
    console.error('Failed to connect to Ganache:', error.message);
    console.log('\nMake sure Ganache is running: npm run ganache');
    process.exit(1);
  }

  // Load compiled contract
  const contractPath = path.join(
    __dirname,
    '..',
    'build',
    'InventoryAudit.bin'
  );
  const abiPath = path.join(__dirname, '..', 'build', 'InventoryAudit.abi');

  if (!fs.existsSync(contractPath)) {
    console.error(
      'Contract bytecode not found. Please compile first: npm run compile'
    );
    console.log('Expected path:', contractPath);
    process.exit(1);
  }

  let bytecode = fs.readFileSync(contractPath, 'utf8').trim();
  const abi = JSON.parse(fs.readFileSync(abiPath, 'utf8'));

  console.log('Contract loaded successfully');

  // Create account from private key
  const account = web3.eth.accounts.privateKeyToAccount(DEFAULT_PRIVATE_KEY);
  web3.eth.accounts.wallet.add(account);

  console.log('Deploying from address:', account.address);

  const gasPrice = await web3.eth.getGasPrice();

  const receipt = deployedContract.deploy({
    from: account.address,
    gas: estimatedGas,
    gasPrice: gasPrice,
  });

  // Then use receipt.transactionHash
  console.log('Transaction hash:', receipt.transactionHash);

  // Deploy contract
  const contract = new web3.eth.Contract(abi);

  console.log('Deploying InventoryAudit contract...');

  const deployTx = contract.deploy({
    data: bytecode,
  });

  const estimatedGas = await deployTx.estimateGas({
    from: account.address,
  });

  console.log('Estimated gas:', estimatedGas.toString());

  const deployedContract = await deployTx.send({
    from: account.address,
    gas: estimatedGas,
    gasPrice: gasPrice,
  });

  console.log('\n✅ Contract deployed successfully!');
  console.log('Contract address:', deployedContract.options.address);
  console.log('Transaction hash:', deployedContract.transactionHash);

  // Save deployment info with BigInt serialization
  const networkId = await web3.eth.net.getId();
  const deploymentInfo = {
    contractAddress: deployedContract.options.address,
    transactionHash: deployedContract.transactionHash,
    networkId: networkId.toString(), // Convert BigInt to string
    deployedAt: new Date().toISOString(),
  };

  const outputPath = path.join(__dirname, '..', 'deployment.json');
  fs.writeFileSync(outputPath, JSON.stringify(deploymentInfo, replacer, 2));

  console.log('\nDeployment info saved to:', outputPath);
  console.log('\nTo use in Go backend, set:');
  console.log(
    `BLOCKCHAIN_CONTRACT_ADDRESS=${deployedContract.options.address}`
  );
  console.log(`BLOCKCHAIN_RPC_URL=${GANACHE_URL}`);
}

deploy().catch(console.error);

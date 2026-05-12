const fs = require('fs');
const path = require('path');

async function compile() {
  console.log('Compiling InventoryAudit.sol...');

  const solc = require('solc');
  const contractPath = path.join(
    __dirname,
    '..',
    '..',
    'contracts',
    'InventoryAudit.sol'
  );

  if (!fs.existsSync(contractPath)) {
    console.error('Contract not found at:', contractPath);
    process.exit(1);
  }

  const sourceCode = fs.readFileSync(contractPath, 'utf8');

  const input = {
    language: 'Solidity',
    sources: {
      'InventoryAudit.sol': {
        content: sourceCode,
      },
    },
    settings: {
      outputSelection: {
        '*': {
          '*': ['*'],
        },
      },
    },
  };

  const compiled = JSON.parse(solc.compile(JSON.stringify(input)));

  if (compiled.errors) {
    let hasError = false;
    compiled.errors.forEach((error) => {
      if (error.severity === 'error') {
        hasError = true;
        console.error('Compilation error:', error.formattedMessage);
      } else {
        console.warn('Warning:', error.formattedMessage);
      }
    });
    if (hasError) {
      process.exit(1);
    }
  }

  const contract = compiled.contracts['InventoryAudit.sol']['InventoryAudit'];

  if (!contract) {
    console.error('Contract not found in compilation output');
    process.exit(1);
  }

  // Create build directory if it doesn't exist
  const buildDir = path.join(__dirname, '..', 'build');
  if (!fs.existsSync(buildDir)) {
    fs.mkdirSync(buildDir, { recursive: true });
  }

  // Write bytecode (remove 0x prefix for storage, will add back during deployment)
  let bytecode = contract.evm.bytecode.object;
  if (bytecode.startsWith('0x')) {
    bytecode = bytecode.slice(2);
  }

  const bytecodePath = path.join(buildDir, 'InventoryAudit.bin');
  fs.writeFileSync(bytecodePath, bytecode);
  console.log('Bytecode written to:', bytecodePath);

  // Write ABI
  const abiPath = path.join(buildDir, 'InventoryAudit.abi');
  fs.writeFileSync(abiPath, JSON.stringify(contract.abi, null, 2));
  console.log('ABI written to:', abiPath);

  console.log('\n✅ Compilation successful!');
}

compile().catch(console.error);

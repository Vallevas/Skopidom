const Web3 = require('web3');
const fs = require('fs');
const path = require('path');

// Configuration
const GANACHE_URL = process.env.GANACHE_URL || 'http://127.0.0.1:7545';
const DEFAULT_PRIVATE_KEY = '0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80';

async function compile() {
    console.log('Compiling InventoryAudit.sol...');
    
    const solc = require('solc');
    const contractPath = path.join(__dirname, '..', '..', 'contracts', 'InventoryAudit.sol');
    
    if (!fs.existsSync(contractPath)) {
        console.error('Contract not found at:', contractPath);
        process.exit(1);
    }
    
    const sourceCode = fs.readFileSync(contractPath, 'utf8');
    
    const input = {
        language: 'Solidity',
        sources: {
            'InventoryAudit.sol': {
                content: sourceCode
            }
        },
        settings: {
            outputSelection: {
                '*': {
                    '*': ['*']
                }
            }
        }
    };
    
    const compiled = JSON.parse(solc.compile(JSON.stringify(input)));
    
    if (compiled.errors) {
        let hasError = false;
        compiled.errors.forEach(error => {
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
    
    // Write bytecode
    const bytecodePath = path.join(buildDir, 'InventoryAudit.bin');
    fs.writeFileSync(bytecodePath, contract.evm.bytecode.object);
    console.log('Bytecode written to:', bytecodePath);
    
    // Write ABI
    const abiPath = path.join(buildDir, 'InventoryAudit.abi');
    fs.writeFileSync(abiPath, JSON.stringify(contract.abi, null, 2));
    console.log('ABI written to:', abiPath);
    
    console.log('\n✅ Compilation successful!');
}

compile().catch(console.error);

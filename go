//package.json
{
  "name": "evm-bot",
  "version": "1.0.0",
  "source": "src/index.html",
  "scripts": {
    "start": "parcel",
    "build": "parcel build"
  },
  "dependencies": {
    "web3": "^1.9.0",
    "@web3-react/core": "^6.1.9",
    "@web3-react/injected-connector": "^6.0.7",
    "ethers": "^5.7.2",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "@emotion/react": "^11.11.0",
    "@emotion/styled": "^11.11.0",
    "@mui/material": "^5.13.0",
    "@mui/icons-material": "^5.13.0"
  },
  "devDependencies": {
    "parcel": "^2.8.3",
    "process": "^0.11.10",
    "@babel/core": "^7.21.8",
    "@babel/preset-react": "^7.18.6",
    "@babel/preset-env": "^7.21.8",
    "buffer": "^5.7.1"
  }
}

//src/App.js
import React from 'react';
import { Box, Typography, Container } from '@mui/material';
import { Web3Provider } from './contexts/Web3Context';
import { DexProvider } from './contexts/DexContext';
import { GasProvider } from './contexts/GasContext';
import RPCLatencyMonitor from './components/RPCLatencyMonitor';
import Settings from './components/Settings';
import WalletManager from './components/WalletManager';
import TradingPage from './components/trading/TradingPage';
const App = () => {
    return (
        <Web3Provider>
            <DexProvider>
                <GasProvider>
                    <Container maxWidth="lg">
                        <Box sx={{ my: 4 }}>
                            <Typography variant="h2" component="h1" gutterBottom align="center">
                                Maiko Sniper
                            </Typography>
                            <RPCLatencyMonitor />
                            <Settings />
                            <WalletManager />
                            <TradingPage />
                        </Box>
                    </Container>
                </GasProvider>
            </DexProvider>
        </Web3Provider>
    );
};
export default App;

//src/index.html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Maiko Sniper</title>
</head>
<body>
    <div id="root"></div>
    <script type="module" src="./index.js"></script>
</body>
</html>

//src/index.js
import React from 'react';
import { createRoot } from 'react-dom/client';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from './theme';
import App from './App';
const container = document.getElementById('root');
const root = createRoot(container);
root.render(
    <React.StrictMode>
        <ThemeProvider theme={theme}>
            <CssBaseline />
            <App />
        </ThemeProvider>
    </React.StrictMode>
);

//src/theme.js
import { createTheme } from '@mui/material/styles';
const theme = createTheme({
    palette: {
        mode: 'dark',
    }
});
export default theme;


//src/utils/dexUtils.js

import Web3 from 'web3';

// ABI 定義
export const IERC20_ABI = [
    {"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},
    {"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},
    {"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},
    {"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"},
];

export const PAIR_ABI = [
    {"constant":true,"inputs":[],"name":"token0","outputs":[{"name":"","type":"address"}],"type":"function"},
    {"constant":true,"inputs":[],"name":"token1","outputs":[{"name":"","type":"address"}],"type":"function"},
    {"constant":true,"inputs":[],"name":"getReserves","outputs":[{"name":"reserve0","type":"uint112"},{"name":"reserve1","type":"uint112"},{"name":"blockTimestampLast","type":"uint32"}],"type":"function"},
];

export const FACTORY_ABI = [
    {"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"allPairs","outputs":[{"name":"","type":"address"}],"type":"function"},
    {"constant":true,"inputs":[],"name":"allPairsLength","outputs":[{"name":"","type":"uint256"}],"type":"function"},
    {"constant":true,"inputs":[{"name":"tokenA","type":"address"},{"name":"tokenB","type":"address"}],"name":"getPair","outputs":[{"name":"pair","type":"address"}],"type":"function"},
];

// 常用代幣地址和價格（根據不同鏈來設置）
export const CHAIN_TOKENS = {
    ETH: {
        NATIVE_TOKEN: {
            address: '0xC02e9d4aB7D977f075B91a9c27a8CF19eF474C78',
            symbol: 'WETH',
            price: 2000, // 預設價格，實際應該從預言機獲取
            decimals: 18
        },
        USDT: {
            address: '0xdAC17F958D2ee523a2206206994597C13D831ec7',
            symbol: 'USDT',
            price: 1,
            decimals: 6
        },
        USDC: {
            address: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
            symbol: 'USDC',
            price: 1,
            decimals: 6
        },
        WBTC: {
            address: '0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599',
            symbol: 'WBTC',
            price: 35000, // 預設價格，實際應該從預言機獲取
            decimals: 8
        }
    },
    BSC: {
        NATIVE_TOKEN: {
            address: '0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c',
            symbol: 'WBNB',
            price: 300, // 預設價格
            decimals: 18
        },
        USDT: {
            address: '0x55d398326f99059fF775485246999027B3197955',
            symbol: 'USDT',
            price: 1,
            decimals: 18
        },
        BUSD: {
            address: '0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56',
            symbol: 'BUSD',
            price: 1,
            decimals: 18
        }
    },
    // 其他鏈的配置...
};

export async function getTokenInfo(web3, tokenAddress) {
    try {
        const tokenContract = new web3.eth.Contract(IERC20_ABI, tokenAddress);
        const [name, symbol, decimals] = await Promise.all([
            tokenContract.methods.name().call(),
            tokenContract.methods.symbol().call(),
            tokenContract.methods.decimals().call(),
        ]);
        return { name, symbol, decimals: parseInt(decimals) };
    } catch (error) {
        console.error('Error fetching token info:', error);
        return null;
    }
}

export async function findTopPools(web3, factoryAddress, targetToken, chainId) {
    const chainTokens = CHAIN_TOKENS[chainId] || CHAIN_TOKENS.ETH;
    const potentialBaseTokens = Object.values(chainTokens);
    const pools = [];

    for (const baseToken of potentialBaseTokens) {
        try {
            const factory = new web3.eth.Contract(FACTORY_ABI, factoryAddress);
            const pairAddress = await factory.methods.getPair(targetToken, baseToken.address).call();
            
            if (pairAddress === '0x0000000000000000000000000000000000000000') continue;

            const pair = new web3.eth.Contract(PAIR_ABI, pairAddress);
            const [token0, token1, reserves] = await Promise.all([
                pair.methods.token0().call(),
                pair.methods.token1().call(),
                pair.methods.getReserves().call(),
            ]);

            const isToken0Target = token0.toLowerCase() === targetToken.toLowerCase();
            const baseTokenReserve = isToken0Target ? reserves[1] : reserves[0];
            const targetTokenReserve = isToken0Target ? reserves[0] : reserves[1];

            // 計算池子的總價值（以USD為單位）
            const baseTokenValue = (parseFloat(baseTokenReserve) / Math.pow(10, baseToken.decimals)) * baseToken.price;

            pools.push({
                pairAddress,
                baseToken: baseToken,
                baseTokenReserve,
                targetTokenReserve,
                totalValueUSD: baseTokenValue * 2, // 假設池子兩邊價值相等
                priceUSD: baseTokenValue / (parseFloat(targetTokenReserve) / Math.pow(10, baseToken.decimals))
            });

        } catch (error) {
            console.error(`Error checking pair with ${baseToken.symbol}:`, error);
        }
    }

    // 按池子總價值排序
    return pools.sort((a, b) => b.totalValueUSD - a.totalValueUSD);
}

// 格式化數字顯示
export function formatNumber(num, decimals = 2) {
    if (num >= 1e6) {
        return `${(num / 1e6).toFixed(decimals)}M`;
    } else if (num >= 1e3) {
        return `${(num / 1e3).toFixed(decimals)}K`;
    }
    return num.toFixed(decimals);
}

// 格式化地址顯示
export function formatAddress(address) {
    if (!address) return '';
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
}

//src/utils/tokenUtils.js

import Web3 from 'web3';

export const ERC20_ABI = [
    {"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},
    {"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},
    {"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},
    {"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"},
    {"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"type":"function"},
    {"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"type":"function"}
];

export const PAIR_ABI = [
    {"constant":true,"inputs":[],"name":"token0","outputs":[{"name":"","type":"address"}],"type":"function"},
    {"constant":true,"inputs":[],"name":"token1","outputs":[{"name":"","type":"address"}],"type":"function"},
    {"constant":true,"inputs":[],"name":"getReserves","outputs":[{"name":"_reserve0","type":"uint112"},{"name":"_reserve1","type":"uint112"},{"name":"_blockTimestampLast","type":"uint32"}],"type":"function"}
];

export const FACTORY_ABI = [
    {"constant":true,"inputs":[{"name":"tokenA","type":"address"},{"name":"tokenB","type":"address"}],"name":"getPair","outputs":[{"name":"pair","type":"address"}],"type":"function"}
];

export const CHAIN_TOKENS = {
    BSC: {
        WBNB: {
            address: '0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c',
            symbol: 'WBNB',
            decimals: 18,
            price: 300  // 預設價格
        },
        USDT: {
            address: '0x55d398326f99059fF775485246999027B3197955',
            symbol: 'USDT',
            decimals: 18,
            price: 1
        },
        BUSD: {
            address: '0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56',
            symbol: 'BUSD',
            decimals: 18,
            price: 1
        }
    },
    ETH: {
        WBNB: {
            address: '0xC02aaa39b223FE8D0A0e5C4F27eAD9083C756Cc2',
            symbol: 'WETH',
            decimals: 18,
            price: 2000
        },
        USDT: {
            address: '0xdAC17F958D2ee523a2206206994597C13D831ec7',
            symbol: 'USDT',
            decimals: 6,
            price: 1
        }
    },
    CORE: {
        WBNB: {
            address: '0x40375C92d9FAf44d2f9db9Bd9ba41a3317a2404f',
            symbol: 'WCORE',
            decimals: 18,
            price: 1
        },
        USDT: {
            address: '0x900101d06A7426441Ae63e9AB3B9b0F63Be145F1',
            symbol: 'USDT',
            decimals: 18,
            price: 1
        }
    },
    BASE: {
        WBNB: {
            address: '0x4200000000000000000000000000000000000006',
            symbol: 'WETH',
            decimals: 18,
            price: 2000
        },
        USDT: {
            address: '0x50c5725949A6F0c72E6C4a641F24049A917DB0Cb',
            symbol: 'USDT',
            decimals: 18,
            price: 1
        }
    }
};

export const getTokenInfo = async (web3, tokenAddress) => {
    try {
        const tokenContract = new web3.eth.Contract(ERC20_ABI, tokenAddress);
        const [name, symbol, decimals] = await Promise.all([
            tokenContract.methods.name().call(),
            tokenContract.methods.symbol().call(),
            tokenContract.methods.decimals().call(),
        ]);
        return { name, symbol, decimals: parseInt(decimals) };
    } catch (error) {
        console.error('Error fetching token info:', error);
        return null;
    }
};

export const getTokenBalance = async (web3, tokenAddress, accountAddress) => {
    try {
        const tokenContract = new web3.eth.Contract(ERC20_ABI, tokenAddress);
        const balance = await tokenContract.methods.balanceOf(accountAddress).call();
        return web3.utils.fromWei(balance);
    } catch (error) {
        console.error('Error fetching token balance:', error);
        return '0';
    }
};

export const getNativeBalance = async (web3, accountAddress) => {
    try {
        const balance = await web3.eth.getBalance(accountAddress);
        return web3.utils.fromWei(balance);
    } catch (error) {
        console.error('Error fetching native balance:', error);
        return '0';
    }
};

export const formatNumber = (num, decimals = 2) => {
    if (!num || isNaN(num)) return '0';
    const number = parseFloat(num);
    if (number >= 1e6) {
        return `${(number / 1e6).toFixed(decimals)}M`;
    } else if (number >= 1e3) {
        return `${(number / 1e3).toFixed(decimals)}K`;
    }
    return number.toFixed(decimals);
};

export const formatAddress = (address) => {
    if (!address) return '';
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
};

/src/utils/transactionUtils.js

import Web3 from 'web3';

const ROUTER_ABI = [
    // Uniswap V2 風格的路由器 ABI
    {
        "inputs": [
            {"internalType": "uint256", "name": "amountIn", "type": "uint256"},
            {"internalType": "uint256", "name": "amountOutMin", "type": "uint256"},
            {"internalType": "address[]", "name": "path", "type": "address[]"},
            {"internalType": "address", "name": "to", "type": "address"},
            {"internalType": "uint256", "name": "deadline", "type": "uint256"}
        ],
        "name": "swapExactTokensForTokens",
        "outputs": [{"internalType": "uint256[]", "name": "amounts", "type": "uint256[]"}],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {"internalType": "uint256", "name": "amountOutMin", "type": "uint256"},
            {"internalType": "address[]", "name": "path", "type": "address[]"},
            {"internalType": "address", "name": "to", "type": "address"},
            {"internalType": "uint256", "name": "deadline", "type": "uint256"}
        ],
        "name": "swapExactETHForTokens",
        "outputs": [{"internalType": "uint256[]", "name": "amounts", "type": "uint256[]"}],
        "stateMutability": "payable",
        "type": "function"
    },
    {
        "inputs": [
            {"internalType": "uint256", "name": "amountIn", "type": "uint256"},
            {"internalType": "uint256", "name": "amountOutMin", "type": "uint256"},
            {"internalType": "address[]", "name": "path", "type": "address[]"},
            {"internalType": "address", "name": "to", "type": "address"},
            {"internalType": "uint256", "name": "deadline", "type": "uint256"}
        ],
        "name": "swapExactTokensForETH",
        "outputs": [{"internalType": "uint256[]", "name": "amounts", "type": "uint256[]"}],
        "stateMutability": "nonpayable",
        "type": "function"
    }
];

export async function executeSwap(web3, {
    routerAddress,
    path,
    amount,
    minReceived,
    account,
    privateKey,
    isNativeToken
}) {
    const router = new web3.eth.Contract(ROUTER_ABI, routerAddress);
    const deadline = Math.floor(Date.now() / 1000) + 60 * 20; // 20 minutes
    
    let txData;
    let value = '0';

    if (isNativeToken) {
        // 用原生代幣購買代幣
        txData = router.methods.swapExactETHForTokens(
            web3.utils.toWei(minReceived.toString()),
            path,
            account,
            deadline
        ).encodeABI();
        value = web3.utils.toWei(amount.toString());
    } else {
        // 代幣交換
        txData = router.methods.swapExactTokensForTokens(
            web3.utils.toWei(amount.toString()),
            web3.utils.toWei(minReceived.toString()),
            path,
            account,
            deadline
        ).encodeABI();
    }

    const tx = {
        from: account,
        to: routerAddress,
        data: txData,
        value: value,
        gas: await web3.eth.estimateGas({
            from: account,
            to: routerAddress,
            data: txData,
            value: value
        }),
        gasPrice: await web3.eth.getGasPrice(),
        nonce: await web3.eth.getTransactionCount(account)
    };

    const signedTx = await web3.eth.accounts.signTransaction(tx, privateKey);
    return web3.eth.sendSignedTransaction(signedTx.rawTransaction);
}

// 檢查和授權代幣
export async function checkAndApproveToken(web3, tokenAddress, ownerAddress, spenderAddress, amount) {
    const tokenContract = new web3.eth.Contract(IERC20_ABI, tokenAddress);
    const currentAllowance = await tokenContract.methods.allowance(ownerAddress, spenderAddress).call();
    
    if (new Web3.utils.BN(currentAllowance).lt(new Web3.utils.BN(amount))) {
        const approveData = tokenContract.methods.approve(
            spenderAddress,
            '115792089237316195423570985008687907853269984665640564039457584007913129639935' // uint256 max
        ).encodeABI();

        const tx = {
            from: ownerAddress,
            to: tokenAddress,
            data: approveData,
            gas: await web3.eth.estimateGas({
                from: ownerAddress,
                to: tokenAddress,
                data: approveData
            }),
            gasPrice: await web3.eth.getGasPrice(),
            nonce: await web3.eth.getTransactionCount(ownerAddress)
        };

        return tx;
    }
    return null;
}

//src/utils/sniperUtils.js

// src/utils/sniperUtils.js
import Web3 from 'web3';
import { executeSwap, checkAndApproveToken } from './transactionUtils';
import { CHAIN_TOKENS } from './tokenUtils';

const ROUTER_ABI = [
    {
        "inputs": [
            {"internalType": "uint256", "name": "amountIn", "type": "uint256"},
            {"internalType": "uint256", "name": "amountOutMin", "type": "uint256"},
            {"internalType": "address[]", "name": "path", "type": "address[]"},
            {"internalType": "address", "name": "to", "type": "address"},
            {"internalType": "uint256", "name": "deadline", "type": "uint256"}
        ],
        "name": "swapExactTokensForTokens",
        "outputs": [{"internalType": "uint256[]", "name": "amounts", "type": "uint256[]"}],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {"internalType": "uint256", "name": "amountOutMin", "type": "uint256"},
            {"internalType": "address[]", "name": "path", "type": "address[]"},
            {"internalType": "address", "name": "to", "type": "address"},
            {"internalType": "uint256", "name": "deadline", "type": "uint256"}
        ],
        "name": "swapExactETHForTokens",
        "outputs": [{"internalType": "uint256[]", "name": "amounts", "type": "uint256[]"}],
        "stateMutability": "payable",
        "type": "function"
    }
];

const ERC20_ABI = [
    {"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},
    {"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},
    {"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},
    {"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"},
    {"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"type":"function"},
    {"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"type":"function"}
];

export class SniperService {
    constructor(web3, chainId, dex, accounts) {
        this.web3 = web3;
        this.chainId = chainId;
        this.dex = dex;
        this.accounts = accounts;
        this.subscriptions = new Map();
        this.pendingSnipes = new Map();
        this.blockSubscription = null;
    }

    // cleanup 方法保持不變...

    // startMonitoring 方法保持不變...

    // setupBlockSubscription 方法保持不變...

    // stopMonitoring 方法保持不變...

    async isAddLiquidityTx(tx, tokenAddress, monitorTokenAddress, monitorTokenAmount, nativeTokenAmount) {
        try {
            const input = tx.input;
            if (tx.to?.toLowerCase() !== this.dex.router.toLowerCase()) {
                return false;
            }

            const methodId = input.slice(0, 10);
            const addLiquidityMethodIds = [
                '0xe8e33700', // addLiquidity
                '0xf305d719', // addLiquidityETH
            ];

            if (!addLiquidityMethodIds.includes(methodId)) {
                return false;
            }

            let decodedInput;
            let isNativeTokenPair = false;
            let tokenAmounts = {};

            if (methodId === '0xf305d719') { // addLiquidityETH
                decodedInput = this.web3.eth.abi.decodeParameters(
                    ['address', 'uint256', 'uint256', 'uint256', 'address', 'uint256'],
                    input.slice(10)
                );
                isNativeTokenPair = true;
                tokenAmounts.native = this.web3.utils.fromWei(tx.value);
                tokenAmounts.token = this.web3.utils.fromWei(decodedInput[1]);
            } else { // addLiquidity
                decodedInput = this.web3.eth.abi.decodeParameters(
                    ['address', 'address', 'uint256', 'uint256', 'uint256', 'uint256', 'address', 'uint256'],
                    input.slice(10)
                );
                tokenAmounts.token0 = this.web3.utils.fromWei(decodedInput[2]);
                tokenAmounts.token1 = this.web3.utils.fromWei(decodedInput[3]);
            }

            // 檢查目標代幣
            const targetTokenMatch = isNativeTokenPair
                ? decodedInput[0].toLowerCase() === tokenAddress.toLowerCase()
                : (decodedInput[0].toLowerCase() === tokenAddress.toLowerCase() || 
                   decodedInput[1].toLowerCase() === tokenAddress.toLowerCase());

            if (!targetTokenMatch) {
                return false;
            }

            // 檢查原生代幣數量
            if (isNativeTokenPair) {
                if (parseFloat(tokenAmounts.native) < parseFloat(nativeTokenAmount)) {
                    return false;
                }
            }

            // 檢查監控代幣
            if (monitorTokenAddress && monitorTokenAmount) {
                let monitorTokenFound = false;
                let monitorTokenValue = '0';

                if (isNativeTokenPair) {
                    if (decodedInput[0].toLowerCase() === monitorTokenAddress.toLowerCase()) {
                        monitorTokenFound = true;
                        monitorTokenValue = tokenAmounts.token;
                    }
                } else {
                    if (decodedInput[0].toLowerCase() === monitorTokenAddress.toLowerCase()) {
                        monitorTokenFound = true;
                        monitorTokenValue = tokenAmounts.token0;
                    } else if (decodedInput[1].toLowerCase() === monitorTokenAddress.toLowerCase()) {
                        monitorTokenFound = true;
                        monitorTokenValue = tokenAmounts.token1;
                    }
                }

                if (!monitorTokenFound || parseFloat(monitorTokenValue) < parseFloat(monitorTokenAmount)) {
                    return false;
                }
            }

            return true;
        } catch (error) {
            console.error('檢查加池交易時出錯:', error);
            return false;
        }
    }

    async executeSnipe(snipeTask, onStatusUpdate, taskId) {
        const { 
            tokenAddress, 
            buyToken, 
            buyAmount, 
            slippage,
            usdtAddress 
        } = snipeTask;

        try {
            onStatusUpdate(taskId, '準備狙擊');

            const enabledAccounts = this.accounts.filter(acc => acc.enabled);
            if (enabledAccounts.length === 0) {
                throw new Error('沒有啟用的錢包');
            }

            const amountPerWallet = parseFloat(buyAmount) / enabledAccounts.length;
            
            const promises = enabledAccounts.map(async (account) => {
                try {
                    const nativeToken = CHAIN_TOKENS[this.chainId].NATIVE_TOKEN;
                    const path = buyToken === 'NATIVE' 
                        ? [nativeToken.address, tokenAddress]
                        : [usdtAddress, tokenAddress];

                    // 如果使用USDT，需要先授權
                    if (buyToken === 'USDT') {
                        const usdtContract = new this.web3.eth.Contract(ERC20_ABI, usdtAddress);
                        const allowance = await usdtContract.methods
                            .allowance(account.address, this.dex.router)
                            .call();

                        if (this.web3.utils.toBN(allowance).lt(
                            this.web3.utils.toBN(this.web3.utils.toWei(amountPerWallet.toString()))
                        )) {
                            const approveTx = {
                                from: account.address,
                                to: usdtAddress,
                                data: usdtContract.methods.approve(
                                    this.dex.router,
                                    '115792089237316195423570985008687907853269984665640564039457584007913129639935'
                                ).encodeABI(),
                                gasPrice: await this.web3.eth.getGasPrice()
                            };
                            approveTx.gas = await this.web3.eth.estimateGas(approveTx);

                            const signedApproveTx = await this.web3.eth.accounts.signTransaction(
                                approveTx, 
                                account.privateKey
                            );
                            await this.web3.eth.sendSignedTransaction(signedApproveTx.rawTransaction);
                        }
                    }

                    // 執行買入交易
                    const amountInWei = this.web3.utils.toWei(amountPerWallet.toString());
                    const slippageMultiplier = 1 - (parseFloat(slippage) / 100);
                    const minOutWei = this.web3.utils.toWei(
                        (amountPerWallet * slippageMultiplier).toString()
                    );

                    await executeSwap(this.web3, {
                        routerAddress: this.dex.router,
                        path: path,
                        amount: amountInWei,
                        minReceived: minOutWei,
                        account: account.address,
                        privateKey: account.privateKey,
                        isNativeToken: buyToken === 'NATIVE'
                    });

                    return true;
                } catch (error) {
                    console.error(`錢包 ${account.address} 狙擊失敗:`, error);
                    return false;
                }
            });

            const results = await Promise.all(promises);
            const successCount = results.filter(Boolean).length;

            if (successCount > 0) {
                onStatusUpdate(taskId, `狙擊成功 (${successCount}/${enabledAccounts.length})`);
            } else {
                onStatusUpdate(taskId, '狙擊失敗');
            }

        } catch (error) {
            console.error('執行狙擊時出錯:', error);
            onStatusUpdate(taskId, `狙擊失敗: ${error.message}`);
            throw error;
        }
    }

    // 獲取當前 gas 價格
    async getGasPrice() {
        try {
            const gasPrice = await this.web3.eth.getGasPrice();
            return gasPrice;
        } catch (error) {
            console.error('獲取 gas 價格失敗:', error);
            throw error;
        }
    }

    // 估算交易的 gas 限制
    async estimateGas(tx) {
        try {
            const gasLimit = await this.web3.eth.estimateGas(tx);
            return Math.floor(gasLimit * 1.2); // 增加 20% 作為緩衝
        } catch (error) {
            console.error('估算 gas 限制失敗:', error);
            throw error;
        }
    }

    // 清理所有訂閱
    cleanup() {
        // 清理交易池監聽
        for (const subscription of this.subscriptions.values()) {
            try {
                subscription.unsubscribe();
            } catch (error) {
                console.error('清理交易池監聽失敗:', error);
            }
        }
        this.subscriptions.clear();

        // 清理區塊監聽
        if (this.blockSubscription) {
            try {
                this.blockSubscription.unsubscribe();
            } catch (error) {
                console.error('清理區塊監聽失敗:', error);
            }
            this.blockSubscription = null;
        }

        this.pendingSnipes.clear();
    }

    // 檢查代幣授權狀態
    async checkTokenAllowance(tokenAddress, ownerAddress, spenderAddress, amount) {
        try {
            const tokenContract = new this.web3.eth.Contract(ERC20_ABI, tokenAddress);
            const allowance = await tokenContract.methods.allowance(ownerAddress, spenderAddress).call();
            return this.web3.utils.toBN(allowance).gte(this.web3.utils.toBN(amount));
        } catch (error) {
            console.error('檢查代幣授權失敗:', error);
            throw error;
        }
    }

    // 授權代幣
    async approveToken(tokenAddress, spenderAddress, account, privateKey) {
        try {
            const tokenContract = new this.web3.eth.Contract(ERC20_ABI, tokenAddress);
            const maxApproval = '115792089237316195423570985008687907853269984665640564039457584007913129639935';
            
            const approveTx = {
                from: account,
                to: tokenAddress,
                data: tokenContract.methods.approve(spenderAddress, maxApproval).encodeABI(),
                gasPrice: await this.getGasPrice()
            };

            approveTx.gas = await this.estimateGas(approveTx);
            
            const signedTx = await this.web3.eth.accounts.signTransaction(approveTx, privateKey);
            const receipt = await this.web3.eth.sendSignedTransaction(signedTx.rawTransaction);
            
            return receipt;
        } catch (error) {
            console.error('授權代幣失敗:', error);
            throw error;
        }
    }

    // 獲取代幣餘額
    async getTokenBalance(tokenAddress, accountAddress) {
        try {
            const tokenContract = new this.web3.eth.Contract(ERC20_ABI, tokenAddress);
            const balance = await tokenContract.methods.balanceOf(accountAddress).call();
            return this.web3.utils.fromWei(balance);
        } catch (error) {
            console.error('獲取代幣餘額失敗:', error);
            throw error;
        }
    }

    // 獲取原生代幣餘額
    async getNativeBalance(accountAddress) {
        try {
            const balance = await this.web3.eth.getBalance(accountAddress);
            return this.web3.utils.fromWei(balance);
        } catch (error) {
            console.error('獲取原生代幣餘額失敗:', error);
            throw error;
        }
    }
}


//src/utils/utils.js
export const getNativeCurrencySymbol = (selectedChain) => {
    switch (selectedChain) {
        case 'BSC':
            return 'BNB';
        case 'ETH':
            return 'ETH';
        case 'CORE':
            return 'CORE';
        case 'BASE':
            return 'ETH';
        default:
            return 'ETH';
    }
};

export const prepareGasSettings = (web3, tx, customGasLimit, localGasPrice, gasPrice) => {
    tx.gas = customGasLimit || '21000';

    if (localGasPrice) {
        tx.gasPrice = web3.utils.toWei(localGasPrice, 'gwei');
    } else if (gasPrice) {
        tx.gasPrice = gasPrice;
    }

    return tx;
};

export const calculateGasFee = (web3, localGasPrice, gasPrice, customGasLimit) => {
    if (!web3) return '0';
    const gwei = localGasPrice || (gasPrice ? web3.utils.fromWei(gasPrice, 'gwei') : '0');
    const gweiToEth = web3.utils.fromWei(web3.utils.toWei(gwei, 'gwei'), 'ether');
    const gasFee = parseFloat(gweiToEth) * parseFloat(customGasLimit || '21000');
    return gasFee.toFixed(8);
};

//src/contexts/DexContext.js

import React, { createContext, useContext, useState } from 'react';

const DexContext = createContext(null);

export const DexProvider = ({ children }) => {
    const [selectedDex, setSelectedDex] = useState({
        name: 'PancakeSwap V2',
        router: '0x10ED43C718714eb63d5aA57B78B54704E256024E',
        factory: '0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73',
        type: 'UniswapV2'
    });

    return (
        <DexContext.Provider value={{
            selectedDex,
            setSelectedDex
        }}>
            {children}
        </DexContext.Provider>
    );
};

export const useDex = () => {
    const context = useContext(DexContext);
    if (!context) {
        throw new Error('useDex must be used within a DexProvider');
    }
    return context;
};

export default DexContext;

//src/contexts/GasContext.js

import React, { createContext, useContext, useState, useEffect } from 'react';

const GasContext = createContext(null);

export const GasProvider = ({ children }) => {
    const [gasPrice, setGasPrice] = useState(null);
    const [customGasLimit, setCustomGasLimit] = useState('21000');
    const [localGasPrice, setLocalGasPrice] = useState('');
    const [lastUpdate, setLastUpdate] = useState(0);

    useEffect(() => {
        const timer = setInterval(() => {
            setLastUpdate(Date.now());
        }, 1000);

        return () => clearInterval(timer);
    }, []);

    const updateGasPrice = async (web3) => {
        if (!web3) return;
        try {
            const price = await web3.eth.getGasPrice();
            setGasPrice(price);
        } catch (error) {
            console.error('Failed to fetch gas price:', error);
        }
    };

    return (
        <GasContext.Provider value={{
            gasPrice,
            setGasPrice,
            customGasLimit,
            setCustomGasLimit,
            localGasPrice,
            setLocalGasPrice,
            updateGasPrice,
            lastUpdate
        }}>
            {children}
        </GasContext.Provider>
    );
};

export const useGas = () => {
    const context = useContext(GasContext);
    if (!context) {
        throw new Error('useGas must be used within a GasProvider');
    }
    return context;
};

export default GasContext;

//src/contexts/Web3Context.js

import React, { createContext, useContext, useState, useEffect } from 'react';
import Web3 from 'web3';

const Web3Context = createContext(null);

export const Web3Provider = ({ children }) => {
    const [web3, setWeb3] = useState(null);
    const [accounts, setAccounts] = useState([]);
    const [selectedChain, setSelectedChain] = useState('BSC');
    const [rpcUrl, setRpcUrl] = useState('https://bsc-dataseed.binance.org');

    const initWeb3 = async (customRpcUrl) => {
        try {
            const web3Instance = new Web3(customRpcUrl || rpcUrl);
            setWeb3(web3Instance);
            // 設置全局 web3 實例用於錢包導入檢查
            window.web3 = web3Instance;
            return web3Instance;
        } catch (error) {
            console.error('Failed to initialize Web3:', error);
        }
    };

    const importWallet = (privateKey) => {
        try {
            if (!privateKey.startsWith('0x')) {
                privateKey = '0x' + privateKey;
            }
            const account = web3.eth.accounts.privateKeyToAccount(privateKey);
            if (!accounts.find(acc => acc.address === account.address)) {
                setAccounts(prev => [...prev, { ...account, enabled: true }]);
                return account;
            }
        } catch (error) {
            console.error('Failed to import wallet:', error);
            throw error;
        }
    };

    useEffect(() => {
        initWeb3();
    }, [rpcUrl]);

    return (
        <Web3Context.Provider value={{
            web3,
            accounts,
            setAccounts,
            selectedChain,
            rpcUrl,
            setRpcUrl,
            setSelectedChain,
            importWallet,
            initWeb3
        }}>
            {children}
        </Web3Context.Provider>
    );
};

export const useWeb3 = () => useContext(Web3Context);

export default Web3Context;

//src/components/DexManager.js

import React, { useState } from 'react';
import {
  Box,
  Paper,
  Typography,
  TextField,
  Button,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
  Snackbar,
  IconButton
} from '@mui/material';
import { Delete as DeleteIcon, Add as AddIcon } from '@mui/icons-material';
import { useDex } from '../contexts/DexContext.js';
import { useWeb3 } from '../contexts/Web3Context.js';

const DexManager = () => {
  const { getDexesByChain, addCustomDex, removeCustomDex, defaultDexes, selectedDex, setSelectedDex } = useDex();
  const { selectedChain } = useWeb3();
  const [openDialog, setOpenDialog] = useState(false);
  const [notification, setNotification] = useState({ open: false, message: '', type: 'success' });
  const [newDex, setNewDex] = useState({
    name: '',
    router: '',
    factory: '',
    type: 'UniswapV2'
  });

  const handleAddDex = () => {
    if (!newDex.name || !newDex.router || !newDex.factory) {
      setNotification({
        open: true,
        message: '請填寫所有必要欄位',
        type: 'error'
      });
      return;
    }

    addCustomDex(selectedChain, newDex);
    setOpenDialog(false);
    setNewDex({ name: '', router: '', factory: '', type: 'UniswapV2' });
    setNotification({
      open: true,
      message: 'DEX添加成功',
      type: 'success'
    });
  };

  const handleRemoveDex = (dexName) => {
    if (selectedDex?.name === dexName) {
      setSelectedDex(null);
    }
    removeCustomDex(selectedChain, dexName);
    setNotification({
      open: true,
      message: 'DEX移除成功',
      type: 'success'
    });
  };

  const availableDexes = getDexesByChain(selectedChain);
  const isCustomDex = (dexName) => {
    return !defaultDexes[selectedChain]?.some(defaultDex => defaultDex.name === dexName);
  };

  return (
    <Paper sx={{ p: 2, mb: 2 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h6">
          DEX管理 ({selectedChain})
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setOpenDialog(true)}
        >
          添加DEX
        </Button>
      </Box>

      <FormControl fullWidth sx={{ mb: 2 }}>
        <InputLabel>選擇DEX</InputLabel>
        <Select
          value={selectedDex ? JSON.stringify(selectedDex) : ''}
          label="選擇DEX"
          onChange={(e) => {
            const value = e.target.value;
            setSelectedDex(value ? JSON.parse(value) : null);
          }}
          renderValue={(value) => JSON.parse(value).name}
          sx={{ mb: 1 }}
        >
          {availableDexes.map((dex, index) => (
            <MenuItem 
              key={`${dex.name}-${index}`} 
              value={JSON.stringify(dex)}
              sx={{ 
                display: 'flex', 
                justifyContent: 'space-between',
                alignItems: 'center',
                width: '100%'
              }}
            >
              <Box sx={{ flexGrow: 1 }}>
                <Typography variant="subtitle1">{dex.name}</Typography>
                <Typography variant="body2" color="text.secondary">
                  Router: {dex.router.slice(0, 6)}...{dex.router.slice(-4)}
                </Typography>
              </Box>
              {isCustomDex(dex.name) && (
                <IconButton 
                  onClick={(e) => {
                    e.stopPropagation();
                    handleRemoveDex(dex.name);
                  }}
                  size="small"
                  sx={{ ml: 1 }}
                >
                  <DeleteIcon color="error" />
                </IconButton>
              )}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      {selectedDex && (
        <Box sx={{ mt: 2, p: 2, bgcolor: 'background.paper', borderRadius: 1 }}>
          <Typography variant="subtitle2" gutterBottom>
            當前選擇的DEX詳細信息：
          </Typography>
          <Typography variant="body2">名稱: {selectedDex.name}</Typography>
          <Typography variant="body2">Router: {selectedDex.router}</Typography>
          <Typography variant="body2">Factory: {selectedDex.factory}</Typography>
          <Typography variant="body2">類型: {selectedDex.type}</Typography>
        </Box>
      )}

      <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>添加新DEX</DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2 }}>
            <TextField
              fullWidth
              label="DEX名稱"
              value={newDex.name}
              onChange={(e) => setNewDex({ ...newDex, name: e.target.value })}
              sx={{ mb: 2 }}
            />
            <TextField
              fullWidth
              label="Router地址"
              value={newDex.router}
              onChange={(e) => setNewDex({ ...newDex, router: e.target.value })}
              sx={{ mb: 2 }}
            />
            <TextField
              fullWidth
              label="Factory地址"
              value={newDex.factory}
              onChange={(e) => setNewDex({ ...newDex, factory: e.target.value })}
              sx={{ mb: 2 }}
            />
            <FormControl fullWidth>
              <InputLabel>DEX類型</InputLabel>
              <Select
                value={newDex.type}
                label="DEX類型"
                onChange={(e) => setNewDex({ ...newDex, type: e.target.value })}
              >
                <MenuItem value="UniswapV2">Uniswap V2</MenuItem>
                <MenuItem value="UniswapV3">Uniswap V3</MenuItem>
              </Select>
            </FormControl>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>取消</Button>
          <Button onClick={handleAddDex} variant="contained">添加</Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={notification.open}
        autoHideDuration={6000}
        onClose={() => setNotification({ ...notification, open: false })}
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      >
        <Alert
          onClose={() => setNotification({ ...notification, open: false })}
          severity={notification.type}
          sx={{ width: '100%' }}
        >
          {notification.message}
        </Alert>
      </Snackbar>
    </Paper>
  );
};

export default DexManager;

//src/components/RPCLatencyMonitor.js

// src/components/RPCLatencyMonitor.js
import React, { useState, useEffect } from 'react';
import { 
    Box, 
    Typography, 
    Paper, 
    Grid,
    Card,
    CardContent
} from '@mui/material';
import { useWeb3 } from '../contexts/Web3Context';

const RPCLatencyMonitor = () => {
    const { web3, rpcUrl, selectedChain } = useWeb3();
    const [latency, setLatency] = useState(null);
    const [status, setStatus] = useState('測試中...');
    const [lastUpdate, setLastUpdate] = useState(0);
    const [currentBlock, setCurrentBlock] = useState(null);

    useEffect(() => {
        const timer = setInterval(() => {
            setLastUpdate(Date.now());
        }, 1000);

        return () => clearInterval(timer);
    }, []);

    useEffect(() => {
        if (web3) {
            testLatency();
            updateBlockNumber();
        }
    }, [web3, lastUpdate]);

    const updateBlockNumber = async () => {
        try {
            const blockNumber = await web3.eth.getBlockNumber();
            setCurrentBlock(blockNumber);
        } catch (error) {
            console.error('獲取區塊數失敗:', error);
            setCurrentBlock(null);
        }
    };

    const testLatency = async () => {
        if (!web3) return;
        
        try {
            const startTime = Date.now();
            await web3.eth.getBlockNumber();
            const endTime = Date.now();
            const currentLatency = endTime - startTime;
            setLatency(currentLatency);
            
            if (currentLatency < 100) {
                setStatus('極佳');
            } else if (currentLatency < 300) {
                setStatus('良好');
            } else if (currentLatency < 500) {
                setStatus('一般');
            } else {
                setStatus('較差');
            }
        } catch (error) {
            console.error('RPC測試失敗:', error);
            setLatency(null);
            setStatus('連接失敗');
        }
    };

    const getStatusColor = () => {
        switch (status) {
            case '極佳':
                return 'success.main';
            case '良好':
                return 'info.main';
            case '一般':
                return 'warning.main';
            case '較差':
            case '連接失敗':
                return 'error.main';
            default:
                return 'text.primary';
        }
    };

    return (
        <Paper sx={{ p: 2, mb: 2 }}>
            <Grid container spacing={2} alignItems="center">
                <Grid item xs={12} sm={4}>
                    <Typography variant="h6" component="div">
                        RPC 監控
                    </Typography>
                    <Typography variant="body2" component="div" sx={{ wordBreak: 'break-all' }}>
                        當前鏈: {selectedChain}
                    </Typography>
                </Grid>
                
                <Grid item xs={12} sm={4}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                        <Typography variant="body2" component="div">
                            延遲: {latency !== null ? `${latency}ms` : '測試中...'}
                        </Typography>
                        <Typography 
                            variant="body2" 
                            component="div"
                            sx={{ 
                                color: getStatusColor(),
                                fontWeight: 'bold'
                            }}
                        >
                            狀態: {status}
                        </Typography>
                    </Box>
                    <Typography variant="body2" component="div" sx={{ mt: 1 }}>
                        RPC URL: {rpcUrl}
                    </Typography>
                </Grid>

                <Grid item xs={12} sm={4}>
                    <Typography 
                        variant="h6" 
                        component="div" 
                        sx={{ 
                            textAlign: 'right',
                            color: 'primary.main',
                            fontWeight: 'bold'
                        }}
                    >
                        當前區塊: {currentBlock !== null ? currentBlock.toLocaleString() : '載入中...'}
                    </Typography>
                </Grid>
            </Grid>
        </Paper>
    );
};

export default RPCLatencyMonitor;

//src/components/Settings.js

import React from 'react';
import { 
    Box, 
    TextField, 
    Select, 
    MenuItem, 
    Paper,
    Typography,
    Grid,
    InputLabel,
    FormControl
} from '@mui/material';
import { useWeb3 } from '../contexts/Web3Context';

const Settings = () => {
    const { rpcUrl, setRpcUrl, selectedChain, setSelectedChain } = useWeb3();

    const chains = [
        { id: 'ETH', name: 'Ethereum', rpc: 'https://eth.llamarpc.com' },
        { id: 'BSC', name: 'Binance Smart Chain', rpc: 'https://bsc-dataseed.binance.org' },
        { id: 'CORE', name: 'Core', rpc: 'https://rpc.coredao.org' },
        { id: 'BASE', name: 'Base', rpc: 'https://mainnet.base.org' },
    ];

    return (
        <Paper sx={{ p: 2, mb: 2 }}>
            <Typography variant="h6" gutterBottom>
                網絡設置
            </Typography>
            <Grid container spacing={2}>
                <Grid item xs={12} md={6}>
                    <FormControl fullWidth>
                        <InputLabel>選擇鏈</InputLabel>
                        <Select
                            value={selectedChain}
                            label="選擇鏈"
                            onChange={(e) => {
                                setSelectedChain(e.target.value);
                                const chain = chains.find(c => c.id === e.target.value);
                                if (chain) setRpcUrl(chain.rpc);
                            }}
                        >
                            {chains.map(chain => (
                                <MenuItem key={chain.id} value={chain.id}>
                                    {chain.name}
                                </MenuItem>
                            ))}
                        </Select>
                    </FormControl>
                </Grid>
                <Grid item xs={12} md={6}>
                    <TextField
                        fullWidth
                        label="自定義 RPC URL"
                        value={rpcUrl}
                        onChange={(e) => setRpcUrl(e.target.value)}
                    />
                </Grid>
            </Grid>
        </Paper>
    );
};

export default Settings;


//src/components/TxPoolMonitor.js
import React, { useState, useEffect } from 'react';
import { Box, TextField, Button } from '@mui/material';
import { useWeb3 } from '../contexts/Web3Context.js';

const TxPoolMonitor = () => {
  const { web3 } = useWeb3();
  const [methodId, setMethodId] = useState('');
  const [isMonitoring, setIsMonitoring] = useState(false);

  const startMonitoring = () => {
    if (!methodId) return;
    
    setIsMonitoring(true);
    const subscription = web3.eth.subscribe('pendingTransactions', (error, txHash) => {
      if (error) {
        console.error('Subscription error:', error);
        return;
      }

      web3.eth.getTransaction(txHash).then(tx => {
        if (tx && tx.input.startsWith(methodId)) {
          // Trigger your transaction here
          console.log('Matching transaction found:', tx);
        }
      });
    });
  };

  useEffect(() => {
    return () => {
      if (isMonitoring) {
        web3.eth.clearSubscriptions();
      }
    };
  }, [isMonitoring]);

  return (
    <Box sx={{ mb: 4 }}>
      <TextField
        fullWidth
        label="Method ID to Monitor"
        value={methodId}
        onChange={(e) => setMethodId(e.target.value)}
        sx={{ mb: 2 }}
      />
      <Button 
        variant="contained"
        onClick={() => isMonitoring ? setIsMonitoring(false) : startMonitoring()}
      >
        {isMonitoring ? 'Stop Monitoring' : 'Start Monitoring'}
      </Button>
    </Box>
  );
};

export default TxPoolMonitor;
//src/components/WalletManager.js
import React, { useState } from 'react';
import { 
    Box, 
    TextField, 
    Button, 
    List, 
    ListItem, 
    ListItemText, 
    Switch,
    Typography,
    Paper,
    Alert,
    Card,
    CardContent
} from '@mui/material';
import { useWeb3 } from '../contexts/Web3Context';

const WalletManager = () => {
    const { importWallet, accounts, setAccounts } = useWeb3();
    const [privateKeys, setPrivateKeys] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    const handleImport = () => {
        if (!privateKeys.trim()) {
            setError('請輸入私鑰');
            return;
        }

        setError('');
        setSuccess('');

        try {
            // 分割私鑰字符串，支持換行符和逗號作為分隔符
            const keys = privateKeys.split(/[,\n]/).map(key => key.trim()).filter(key => key);
            let successCount = 0;
            let duplicateCount = 0;
            let errorCount = 0;

            keys.forEach(key => {
                try {
                    // 檢查是否已經導入
                    const formattedKey = key.startsWith('0x') ? key : `0x${key}`;
                    const address = window.web3.eth.accounts.privateKeyToAccount(formattedKey).address;
                    
                    if (accounts.some(acc => acc.address.toLowerCase() === address.toLowerCase())) {
                        duplicateCount++;
                        return;
                    }

                    importWallet(key);
                    successCount++;
                } catch (error) {
                    console.error('導入錢包失敗:', error);
                    errorCount++;
                }
            });

            // 顯示結果
            setSuccess(`導入完成: ${successCount} 個成功, ${duplicateCount} 個重複, ${errorCount} 個失敗`);
            setPrivateKeys(''); // 清空輸入框
        } catch (error) {
            setError(`導入失敗: ${error.message}`);
        }
    };

    const toggleWallet = (address) => {
        setAccounts(accounts.map(acc => 
            acc.address === address 
                ? { ...acc, enabled: !acc.enabled }
                : acc
        ));
    };

    const toggleAllWallets = (enabled) => {
        setAccounts(accounts.map(acc => ({ ...acc, enabled })));
    };

    return (
        <Paper sx={{ p: 2, mb: 2 }}>
            <Typography variant="h6" gutterBottom>
                錢包管理
            </Typography>

            <Box sx={{ mb: 4 }}>
                <TextField
                    fullWidth
                    multiline
                    rows={4}
                    type="password"
                    label="私鑰 (多個私鑰請用逗號或換行分隔)"
                    value={privateKeys}
                    onChange={(e) => setPrivateKeys(e.target.value)}
                    sx={{ mb: 2 }}
                    autoComplete="off"
                    placeholder="輸入格式：
私鑰1,私鑰2,私鑰3
或
私鑰1
私鑰2
私鑰3"
                />
                <Button 
                    variant="contained" 
                    onClick={handleImport}
                    disabled={!privateKeys.trim()}
                    sx={{ mr: 1 }}
                >
                    批量導入錢包
                </Button>
            </Box>

            {(error || success) && (
                <Alert 
                    severity={error ? "error" : "success"} 
                    sx={{ mb: 2 }}
                >
                    {error || success}
                </Alert>
            )}

            {accounts.length > 0 && (
                <Card>
                    <CardContent>
                        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                            <Typography variant="subtitle1">
                                已導入的錢包 ({accounts.filter(acc => acc.enabled).length}/{accounts.length} 啟用)
                            </Typography>
                            <Box>
                                <Button 
                                    variant="outlined" 
                                    size="small" 
                                    onClick={() => toggleAllWallets(true)}
                                    sx={{ mr: 1 }}
                                >
                                    全部啟用
                                </Button>
                                <Button 
                                    variant="outlined" 
                                    size="small" 
                                    onClick={() => toggleAllWallets(false)}
                                >
                                    全部禁用
                                </Button>
                            </Box>
                        </Box>
                        
                        <List>
                            {accounts.map((account) => (
                                <ListItem 
                                    key={account.address}
                                    secondaryAction={
                                        <Switch
                                            edge="end"
                                            checked={account.enabled}
                                            onChange={() => toggleWallet(account.address)}
                                        />
                                    }
                                    sx={{
                                        bgcolor: 'background.paper',
                                        mb: 1,
                                        border: 1,
                                        borderColor: 'divider',
                                        borderRadius: 1
                                    }}
                                >
                                    <ListItemText
                                        primary={`錢包 ${accounts.indexOf(account) + 1}`}
                                        secondary={
                                            <Typography 
                                                variant="body2" 
                                                sx={{ 
                                                    wordBreak: 'break-all',
                                                    color: account.enabled ? 'success.main' : 'text.secondary'
                                                }}
                                            >
                                                {account.address}
                                            </Typography>
                                        }
                                    />
                                </ListItem>
                            ))}
                        </List>
                    </CardContent>
                </Card>
            )}
        </Paper>
    );
};

export default WalletManager;

//src/components/common/GasSettingsCard.js
import React from 'react';
import { 
    Box, 
    TextField, 
    Typography, 
    Grid,
    InputAdornment
} from '@mui/material';
import { useGas } from '../../contexts/GasContext';
import { calculateGasFee } from '../../utils/utils';

const GasSettingsCard = ({ web3, nativeCurrency }) => {
    const { 
        gasPrice, 
        customGasLimit, 
        setCustomGasLimit,
        localGasPrice,
        setLocalGasPrice
    } = useGas();

    const getCurrentGasPrice = () => {
        if (!gasPrice || !web3) return '載入中...';
        return `${web3.utils.fromWei(gasPrice, 'gwei')} Gwei`;
    };

    return (
        <Box sx={{ p: 2, bgcolor: 'background.paper', borderRadius: 1, mb: 2 }}>
            <Typography variant="subtitle2" gutterBottom>
                Gas 設置
            </Typography>
            <Grid container spacing={2}>
                <Grid item xs={12} sm={4}>
                    <Typography variant="body2" color="textSecondary" gutterBottom>
                        當前 Gas Price: {getCurrentGasPrice()}
                    </Typography>
                    <TextField
                        fullWidth
                        label="自定義 Gas Price"
                        type="number"
                        value={localGasPrice}
                        onChange={(e) => setLocalGasPrice(e.target.value)}
                        InputProps={{
                            endAdornment: <InputAdornment position="end">Gwei</InputAdornment>,
                        }}
                        size="small"
                    />
                </Grid>
                <Grid item xs={12} sm={4}>
                    <Typography variant="body2" color="textSecondary" gutterBottom>
                        Gas Limit
                    </Typography>
                    <TextField
                        fullWidth
                        type="number"
                        value={customGasLimit}
                        onChange={(e) => setCustomGasLimit(e.target.value)}
                        size="small"
                    />
                </Grid>
                <Grid item xs={12} sm={4}>
                    <Typography variant="body2" color="textSecondary" gutterBottom>
                        預留 Gas 費用
                    </Typography>
                    <Typography variant="body1">
                        {calculateGasFee(web3, localGasPrice, gasPrice, customGasLimit)} {nativeCurrency}
                    </Typography>
                </Grid>
            </Grid>
        </Box>
    );
};

export default GasSettingsCard;

//src/components/trading/HexTransaction.js
import React, { useState, useEffect } from 'react';
import { 
    Box, 
    TextField, 
    Button, 
    Paper, 
    Typography, 
    Grid,
    Alert
} from '@mui/material';
import { useWeb3 } from '../../contexts/Web3Context.js';
import { useGas } from '../../contexts/GasContext.js';
import GasSettingsCard from '../common/GasSettingsCard';
import { getNativeCurrencySymbol, prepareGasSettings } from '../../utils/utils';

const HexTransaction = () => {
    const { web3, selectedChain, accounts } = useWeb3();
    const { updateGasPrice, lastUpdate } = useGas();
    
    const [to, setTo] = useState('');
    const [value, setValue] = useState('');
    const [data, setData] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    // 自動更新 gas price
    useEffect(() => {
        if (web3) {
            updateGasPrice(web3);
        }
    }, [web3, lastUpdate]);

    const handleSend = async () => {
        if (!web3 || !to || !data) {
            setError('請填寫所有必要資訊');
            return;
        }

        setLoading(true);
        setError('');
        setSuccess('');

        try {
            const enabledAccounts = accounts.filter(acc => acc.enabled);
            if (enabledAccounts.length === 0) {
                throw new Error('沒有啟用的錢包');
            }

            if (!data.startsWith('0x')) {
                throw new Error('Data 必須以 0x 開頭');
            }

            for (const account of enabledAccounts) {
                let tx = {
                    from: account.address,
                    to: to,
                    data: data,
                    value: value ? web3.utils.toWei(value) : '0'
                };

                tx = prepareGasSettings(web3, tx);
                const signedTx = await web3.eth.accounts.signTransaction(tx, account.privateKey);
                await web3.eth.sendSignedTransaction(signedTx.rawTransaction);
            }

            setSuccess('交易發送成功！');
            setValue('');
            setData('');
        } catch (error) {
            console.error('Transaction failed:', error);
            setError(`交易失敗: ${error.message}`);
        } finally {
            setLoading(false);
        }
    };

    const nativeCurrency = getNativeCurrencySymbol(selectedChain);

    return (
        <Paper sx={{ p: 2, mb: 2 }}>
            <Typography variant="h6" gutterBottom>
                發送自定義交易
            </Typography>

            <Grid container spacing={2}>
                <Grid item xs={12}>
                    <TextField
                        fullWidth
                        label="目標地址"
                        value={to}
                        onChange={(e) => setTo(e.target.value)}
                        disabled={loading}
                    />
                </Grid>

                <Grid item xs={12}>
                    <TextField
                        fullWidth
                        type="number"
                        label={`發送數量 (${nativeCurrency})`}
                        value={value}
                        onChange={(e) => setValue(e.target.value)}
                        disabled={loading}
                    />
                </Grid>

                <Grid item xs={12}>
                    <TextField
                        fullWidth
                        label="Data (Hex)"
                        value={data}
                        onChange={(e) => setData(e.target.value)}
                        disabled={loading}
                        multiline
                        rows={4}
                    />
                </Grid>

                <Grid item xs={12}>
                    <GasSettingsCard 
                        web3={web3}
                        nativeCurrency={nativeCurrency}
                    />
                </Grid>

                {(!!error || !!success) && (
                    <Grid item xs={12}>
                        <Alert severity={error ? "error" : "success"}>
                            {error || success}
                        </Alert>
                    </Grid>
                )}

                <Grid item xs={12}>
                    <Button
                        variant="contained"
                        fullWidth
                        onClick={handleSend}
                        disabled={loading}
                    >
                        {loading ? '交易處理中...' : '發送交易'}
                    </Button>
                </Grid>
            </Grid>
        </Paper>
    );
};

export default HexTransaction;

//src/components/trading/TokenApproval.js

import React, { useState, useEffect } from 'react';
import { 
    Box, 
    TextField, 
    Button, 
    Paper, 
    Typography, 
    Grid,
    Alert,
    FormControlLabel,
    Checkbox
} from '@mui/material';
import { useWeb3 } from '../../contexts/Web3Context.js';
import { useGas } from '../../contexts/GasContext.js';
import { ERC20_ABI } from '../../utils/tokenUtils.js';
import GasSettingsCard from '../common/GasSettingsCard';
import { getNativeCurrencySymbol, prepareGasSettings } from '../../utils/utils';

const TokenApproval = () => {
    const { web3, selectedChain, accounts } = useWeb3();
    const { updateGasPrice, lastUpdate } = useGas();
    const [token, setToken] = useState('');
    const [spender, setSpender] = useState('');
    const [amount, setAmount] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [tokenInfo, setTokenInfo] = useState(null);
    const [useMaxUint, setUseMaxUint] = useState(true);

    // 自動更新 gas price
    useEffect(() => {
        if (web3) {
            updateGasPrice(web3);
        }
    }, [web3, lastUpdate]);

    // 獲取代幣信息
    useEffect(() => {
        if (web3 && token && token.length === 42) {
            const getTokenInfo = async () => {
                try {
                    const tokenContract = new web3.eth.Contract(ERC20_ABI, token);
                    const [name, symbol, decimals] = await Promise.all([
                        tokenContract.methods.name().call(),
                        tokenContract.methods.symbol().call(),
                        tokenContract.methods.decimals().call()
                    ]);
                    setTokenInfo({ name, symbol, decimals: parseInt(decimals) });
                    setError('');
                } catch (error) {
                    console.error('Error fetching token info:', error);
                    setError('無效的代幣地址');
                    setTokenInfo(null);
                }
            };
            getTokenInfo();
        }
    }, [web3, token]);

    const handleApprove = async () => {
        if (!web3 || !token || !spender) {
            setError('請填寫所有必要資訊');
            return;
        }

        setLoading(true);
        setError('');
        setSuccess('');

        try {
            const tokenContract = new web3.eth.Contract(ERC20_ABI, token);
            const enabledAccounts = accounts.filter(acc => acc.enabled);
            
            if (enabledAccounts.length === 0) {
                throw new Error('沒有啟用的錢包');
            }

            const approveAmount = useMaxUint
                ? '115792089237316195423570985008687907853269984665640564039457584007913129639935'
                : web3.utils.toWei(amount);

            for (const account of enabledAccounts) {
                let tx = {
                    from: account.address,
                    to: token,
                    data: tokenContract.methods.approve(spender, approveAmount).encodeABI()
                };

                tx = prepareGasSettings(web3, tx);
                const signedTx = await web3.eth.accounts.signTransaction(tx, account.privateKey);
                await web3.eth.sendSignedTransaction(signedTx.rawTransaction);
            }

            setSuccess('授權成功完成！');
        } catch (error) {
            console.error('Approval failed:', error);
            setError(`授權失敗: ${error.message}`);
        } finally {
            setLoading(false);
        }
    };

    const nativeCurrency = getNativeCurrencySymbol(selectedChain);

    return (
        <Paper sx={{ p: 2, mb: 2 }}>
            <Typography variant="h6" gutterBottom>
                Token 授權
            </Typography>

            <Grid container spacing={2}>
                <Grid item xs={12}>
                    <TextField
                        fullWidth
                        label="Token 地址"
                        value={token}
                        onChange={(e) => setToken(e.target.value)}
                        error={!!error}
                        helperText={error || (tokenInfo && `${tokenInfo.name} (${tokenInfo.symbol})`)}
                        disabled={loading}
                    />
                </Grid>

                <Grid item xs={12}>
                    <TextField
                        fullWidth
                        label="授權地址 (通常是 Router 地址)"
                        value={spender}
                        onChange={(e) => setSpender(e.target.value)}
                        disabled={loading}
                    />
                </Grid>

                <Grid item xs={12}>
                    <FormControlLabel
                        control={
                            <Checkbox
                                checked={useMaxUint}
                                onChange={(e) => setUseMaxUint(e.target.checked)}
                                disabled={loading}
                            />
                        }
                        label="使用最大值授權"
                    />
                </Grid>

                {!useMaxUint && (
                    <Grid item xs={12}>
                        <TextField
                            fullWidth
                            type="number"
                            label="授權數量"
                            value={amount}
                            onChange={(e) => setAmount(e.target.value)}
                            disabled={loading || useMaxUint}
                        />
                    </Grid>
                )}

                <Grid item xs={12}>
                    <GasSettingsCard 
                        web3={web3}
                        nativeCurrency={nativeCurrency}
                    />
                </Grid>

                {(!!error || !!success) && (
                    <Grid item xs={12}>
                        <Alert severity={error ? "error" : "success"}>
                            {error || success}
                        </Alert>
                    </Grid>
                )}

                <Grid item xs={12}>
                    <Button
                        variant="contained"
                        fullWidth
                        onClick={handleApprove}
                        disabled={loading}
                    >
                        {loading ? '授權處理中...' : '執行授權'}
                    </Button>
                </Grid>
            </Grid>
        </Paper>
    );
};

export default TokenApproval;

//src/components/trading/SniperPool.js

// src/components/trading/SniperPool.js
import React, { useState, useEffect } from 'react';
import { 
    Paper, 
    Typography, 
    Box, 
    TextField, 
    Button, 
    Grid,
    Alert,
    Card,
    CardContent,
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableRow,
    IconButton,
    Link,
    Tooltip,
    FormControl,
    InputLabel,
    Select,
    MenuItem
} from '@mui/material';
import { 
    InfoOutlined as InfoIcon,
    OpenInNew as OpenInNewIcon 
} from '@mui/icons-material';
import { useWeb3 } from '../../contexts/Web3Context';
import { useDex } from '../../contexts/DexContext';
import GasSettingsCard from '../common/GasSettingsCard';
import { SniperService } from '../../utils/sniperUtils';

const SniperPool = () => {
    const { web3, accounts, selectedChain } = useWeb3();
    const { selectedDex } = useDex();
    
    const [tokenAddress, setTokenAddress] = useState('');
    const [targetAddress, setTargetAddress] = useState('');
    const [skipBlocks, setSkipBlocks] = useState('0');
    const [nativeTokenAmount, setNativeTokenAmount] = useState('');
    const [monitorTokenAddress, setMonitorTokenAddress] = useState('');
    const [monitorTokenAmount, setMonitorTokenAmount] = useState('');
    const [buyToken, setBuyToken] = useState('NATIVE');
    const [buyAmount, setBuyAmount] = useState('');
    const [slippage, setSlippage] = useState('3');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [snipeList, setSnipeList] = useState([]);
    const [sniperService, setSniperService] = useState(null);

    // 初始化 SniperService
    useEffect(() => {
        if (web3 && selectedDex) {
            const service = new SniperService(web3, selectedChain, selectedDex, accounts);
            setSniperService(service);
            return () => service.cleanup();
        }
    }, [web3, selectedChain, selectedDex, accounts]);

    // 根據選擇的鏈顯示區塊時間
    const getBlockTime = () => {
        const blockTimes = {
            'ETH': '12秒',
            'BSC': '3秒',
            'CORE': '3秒',
            'BASE': '2秒'
        };
        return blockTimes[selectedChain] || '未知';
    };

    // 獲取當前鏈的原生代幣符號
    const getNativeTokenSymbol = () => {
        const symbols = {
            'ETH': 'WETH',
            'BSC': 'WBNB',
            'CORE': 'WCORE',
            'BASE': 'WETH'
        };
        return symbols[selectedChain] || 'Unknown';
    };

    // 獲取當前鏈的USDT地址
    const getChainUSDT = () => {
        const usdtAddresses = {
            'ETH': '0xdAC17F958D2ee523a2206206994597C13D831ec7',
            'BSC': '0x55d398326f99059fF775485246999027B3197955',
            'CORE': '0x900101d06A7426441Ae63e9AB3B9b0F63Be145F1',
            'BASE': '0x50c5725949A6F0c72E6C4a641F24049A917DB0Cb'
        };
        return usdtAddresses[selectedChain];
    };

    // 取得區塊瀏覽器URL
    const getExplorerUrl = (address) => {
        const explorers = {
            'ETH': 'https://etherscan.io',
            'BSC': 'https://bscscan.com',
            'CORE': 'https://scan.coredao.org',
            'BASE': 'https://basescan.org'
        };
        return `${explorers[selectedChain]}/address/${address}`;
    };

    const handleAddSnipe = async () => {
        if (!web3 || !tokenAddress || !nativeTokenAmount || !sniperService || !buyAmount) {
            setError('請填寫所有必要欄位');
            return;
        }

        if (!selectedDex) {
            setError('請先選擇DEX');
            return;
        }

        const enabledAccounts = accounts.filter(acc => acc.enabled);
        if (enabledAccounts.length === 0) {
            setError('請先啟用至少一個錢包');
            return;
        }

        setLoading(true);
        setError('');
        setSuccess('');

        try {
            // 驗證地址格式
            if (!web3.utils.isAddress(tokenAddress)) {
                throw new Error('無效的狙擊代幣地址');
            }
            if (targetAddress && !web3.utils.isAddress(targetAddress)) {
                throw new Error('無效的狙擊目標地址');
            }
            if (monitorTokenAddress && !web3.utils.isAddress(monitorTokenAddress)) {
                throw new Error('無效的監控代幣地址');
            }

            // 驗證數量
            if (isNaN(nativeTokenAmount) || parseFloat(nativeTokenAmount) <= 0) {
                throw new Error('請輸入有效的最小加池數量');
            }
            if (monitorTokenAddress && (!monitorTokenAmount || parseFloat(monitorTokenAmount) <= 0)) {
                throw new Error('請輸入有效的監控代幣數量');
            }
            if (isNaN(buyAmount) || parseFloat(buyAmount) <= 0) {
                throw new Error('請輸入有效的買入數量');
            }
            if (isNaN(slippage) || parseFloat(slippage) < 0 || parseFloat(slippage) > 100) {
                throw new Error('請輸入有效的滑點值(0-100)');
            }

            // 驗證跳過區塊數
            const blockSkip = parseInt(skipBlocks);
            if (isNaN(blockSkip) || blockSkip < 0) {
                throw new Error('請輸入有效的跳過區塊數');
            }

            // 創建狙擊任務
            const snipeTask = {
                tokenAddress,
                targetAddress: targetAddress || '',
                skipBlocks: blockSkip,
                nativeTokenAmount,
                monitorTokenAddress,
                monitorTokenAmount,
                buyToken,
                buyAmount,
                slippage: parseFloat(slippage),
                usdtAddress: buyToken === 'USDT' ? getChainUSDT() : null
            };

            // 開始監聽
            const taskId = await sniperService.startMonitoring(
                snipeTask,
                (taskId, status) => {
                    setSnipeList(prev => prev.map(item => 
                        item.taskId === taskId 
                            ? { ...item, status } 
                            : item
                    ));
                }
            );

            // 添加到狙擊列表
            const newSnipe = {
                id: Date.now(),
                taskId,
                tokenAddress,
                targetAddress: targetAddress || '任意地址',
                skipBlocks: blockSkip,
                nativeTokenAmount,
                monitorTokenAddress,
                monitorTokenAmount,
                buyToken,
                buyAmount,
                slippage,
                status: '監聽中',
                timestamp: new Date().toLocaleString()
            };

            setSnipeList(prev => [...prev, newSnipe]);
            setSuccess('狙擊任務已添加！');

            // 清空輸入
            setTokenAddress('');
            setTargetAddress('');
            setSkipBlocks('0');
            setNativeTokenAmount('');
            setMonitorTokenAddress('');
            setMonitorTokenAmount('');
            setBuyAmount('');
            setSlippage('3');

        } catch (error) {
            console.error('添加狙擊任務失敗:', error);
            setError(error.message);
        } finally {
            setLoading(false);
        }
    };

    const handleRemoveSnipe = (id) => {
        const snipe = snipeList.find(item => item.id === id);
        if (snipe && snipe.taskId && sniperService) {
            sniperService.stopMonitoring(snipe.taskId);
        }
        setSnipeList(prev => prev.filter(item => item.id !== id));
    };

    return (
        <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
                狙擊加池
            </Typography>

            <Grid container spacing={2}>
                <Grid item xs={12}>
                    <Typography variant="subtitle1" gutterBottom>
                        監控設置
                    </Typography>
                </Grid>

                <Grid item xs={12}>
                    <TextField
                        fullWidth
                        label="狙擊代幣地址 *"
                        value={tokenAddress}
                        onChange={(e) => setTokenAddress(e.target.value)}
                        disabled={loading}
                        required
                        helperText="必填：要狙擊的代幣合約地址"
                    />
                </Grid>

                <Grid item xs={12}>
                    <TextField
                        fullWidth
                        label="狙擊目標地址"
                        value={targetAddress}
                        onChange={(e) => setTargetAddress(e.target.value)}
                        disabled={loading}
                        helperText="選填：特定目標的錢包地址，留空表示接受任意地址"
                    />
                </Grid>

                <Grid item xs={12} md={6}>
                    <TextField
                        fullWidth
                        type="number"
                        label="跳過區塊數"
                        value={skipBlocks}
                        onChange={(e) => setSkipBlocks(e.target.value)}
                        disabled={loading}
                        InputProps={{
                            inputProps: { min: 0 }
                        }}
                        helperText={`當前鏈(${selectedChain})每區塊時間約為 ${getBlockTime()}`}
                    />
                </Grid>

                <Grid item xs={12} md={6}>
                    <TextField
                        fullWidth
                        type="number"
                        label={`最小${getNativeTokenSymbol()}加池數量 *`}
                        value={nativeTokenAmount}
                        onChange={(e) => setNativeTokenAmount(e.target.value)}
                        disabled={loading}
                        required
                        InputProps={{
                            inputProps: { min: 0, step: "any" }
                        }}
                        helperText={`加池數量超過此數值才會觸發`}
                    />
                </Grid>

                <Grid item xs={12} md={6}>
                    <TextField
                        fullWidth
                        label="監控其他代幣合約"
                        value={monitorTokenAddress}
                        onChange={(e) => setMonitorTokenAddress(e.target.value)}
                        disabled={loading}
                        helperText="選填：要監控的其他代幣合約地址（如USDT）"
                    />
                </Grid>

                <Grid item xs={12} md={6}>
                    <TextField
                        fullWidth
                        type="number"
                        label="監控代幣數量"
                        value={monitorTokenAmount}
                        onChange={(e) => setMonitorTokenAmount(e.target.value)}
                        disabled={loading || !monitorTokenAddress}
                        InputProps={{
                            inputProps: { min: 0, step: "any" }
                        }}
                        helperText="若設置監控代幣，則加池數量超過此數值才會觸發"
                    />
                </Grid>

                <Grid item xs={12}>
                    <Typography variant="subtitle1" gutterBottom sx={{ mt: 2 }}>
                        買入設置
                    </Typography>
                </Grid>

                <Grid item xs={12} md={4}>
                    <FormControl fullWidth>
                        <InputLabel>買入方式</InputLabel>
                        <Select
                            value={buyToken}
                            onChange={(e) => setBuyToken(e.target.value)}
                            label="買入方式"
                            disabled={loading}
                        >
                            <MenuItem value="NATIVE">{getNativeTokenSymbol()}</MenuItem>
                            <MenuItem value="USDT">USDT</MenuItem>
                        </Select>
                    </FormControl>
                </Grid>

                <Grid item xs={12} md={4}>
                    <TextField
                        fullWidth
                        type="number"
                        label={`買入數量 (${buyToken === 'NATIVE' ? getNativeTokenSymbol() : 'USDT'})`}
                        value={buyAmount}
                        onChange={(e) => setBuyAmount(e.target.value)}
                        disabled={loading}
                        required
                        InputProps={{
                            inputProps: { min: 0, step: "any" }
                        }}
                        helperText="執行狙擊時的買入數量"
                    />
                </Grid>

                <Grid item xs={12} md={4}>
                    <TextField
                        fullWidth
                        type="number"
                        label="滑點 %"
                        value={slippage}
                        onChange={(e) => setSlippage(e.target.value)}
                        disabled={loading}
                        InputProps={{
                            inputProps: { min: 0, max: 100, step: "0.1" }
                        }}
                        helperText="允許的最大滑點百分比"
                    />
                </Grid>

                <Grid item xs={12}>
                    <GasSettingsCard web3={web3} />
                </Grid>

                {(error || success) && (
                    <Grid item xs={12}>
                        <Alert severity={error ? "error" : "success"}>
                            {error || success}
                        </Alert>
                    </Grid>
                )}

                <Grid item xs={12}>
                    <Button
                        variant="contained"
                        fullWidth
                        onClick={handleAddSnipe}
                        disabled={loading || !web3 || !selectedDex || accounts.filter(acc => acc.enabled).length === 0}
                    >
                        {loading ? '處理中...' : '添加狙擊任務'}
                    </Button>
                </Grid>

                {snipeList.length > 0 && (
                    <Grid item xs={12}>
                        <Card sx={{ mt: 2 }}>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    狙擊任務列表
                                </Typography>
                                <Table size="small">
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>代幣地址</TableCell>
<TableCell>目標地址</TableCell>
                                            <TableCell align="center">跳過區塊</TableCell>
                                            <TableCell align="right">觸發條件</TableCell>
                                            <TableCell align="right">買入設置</TableCell>
                                            <TableCell align="center">狀態</TableCell>
                                            <TableCell align="center">添加時間</TableCell>
                                            <TableCell align="right">操作</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {snipeList.map((item) => (
                                            <TableRow key={item.id}>
                                                <TableCell>
                                                    <Link
                                                        href={getExplorerUrl(item.tokenAddress)}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        sx={{ display: 'flex', alignItems: 'center' }}
                                                    >
                                                        {`${item.tokenAddress.slice(0, 6)}...${item.tokenAddress.slice(-4)}`}
                                                        <OpenInNewIcon sx={{ ml: 0.5, fontSize: 16 }} />
                                                    </Link>
                                                </TableCell>
                                                <TableCell>
                                                    {item.targetAddress === '任意地址' ? (
                                                        item.targetAddress
                                                    ) : (
                                                        <Link
                                                            href={getExplorerUrl(item.targetAddress)}
                                                            target="_blank"
                                                            rel="noopener noreferrer"
                                                            sx={{ display: 'flex', alignItems: 'center' }}
                                                        >
                                                            {`${item.targetAddress.slice(0, 6)}...${item.targetAddress.slice(-4)}`}
                                                            <OpenInNewIcon sx={{ ml: 0.5, fontSize: 16 }} />
                                                        </Link>
                                                    )}
                                                </TableCell>
                                                <TableCell align="center">{item.skipBlocks}</TableCell>
                                                <TableCell align="right">
                                                    <Typography variant="body2">
                                                        {`>${item.nativeTokenAmount} ${getNativeTokenSymbol()} 加池`}
                                                    </Typography>
                                                    {item.monitorTokenAddress && (
                                                        <Typography variant="body2" color="text.secondary">
                                                            {`>${item.monitorTokenAmount} 監控代幣`}
                                                        </Typography>
                                                    )}
                                                </TableCell>
                                                <TableCell align="right">
                                                    <Typography variant="body2">
                                                        {`${item.buyAmount} ${item.buyToken}`}
                                                    </Typography>
                                                    <Typography variant="body2" color="text.secondary">
                                                        {`滑點: ${item.slippage}%`}
                                                    </Typography>
                                                </TableCell>
                                                <TableCell align="center">
                                                    <Typography
                                                        sx={{
                                                            color: item.status.includes('成功') ? 'success.main' : 
                                                                   item.status.includes('失敗') ? 'error.main' : 
                                                                   'info.main'
                                                        }}
                                                    >
                                                        {item.status}
                                                    </Typography>
                                                </TableCell>
                                                <TableCell align="center">{item.timestamp}</TableCell>
                                                <TableCell align="right">
                                                    <Button 
                                                        size="small" 
                                                        color="error"
                                                        onClick={() => handleRemoveSnipe(item.id)}
                                                    >
                                                        移除
                                                    </Button>
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </CardContent>
                        </Card>
                    </Grid>
                )}
            </Grid>
        </Paper>
    );
};

export default SniperPool;

//src/components/trading/SniperSwitch.js

// src/components/trading/SniperSwitch.js
import React, { useState } from 'react';
import { 
    Paper, 
    Typography, 
    Box,
    Grid,
    TextField,
    Button,
    Alert,
    Switch,
    FormControlLabel,
    Card,
    CardContent,
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableRow
} from '@mui/material';
import { useWeb3 } from '../../contexts/Web3Context';
import { useDex } from '../../contexts/DexContext';
import GasSettingsCard from '../common/GasSettingsCard';

const SniperSwitch = () => {
    const { web3, accounts, selectedChain } = useWeb3();
    const { selectedDex } = useDex();

    const [tokenAddress, setTokenAddress] = useState('');
    const [maxBuyAmount, setMaxBuyAmount] = useState('');
    const [minLiquidity, setMinLiquidity] = useState('');
    const [enabled, setEnabled] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [snipeTargets, setSnipeTargets] = useState([]);

    const handleAddTarget = async () => {
        if (!web3 || !tokenAddress || !maxBuyAmount || !minLiquidity) {
            setError('請填寫所有必要訊息');
            return;
        }

        setLoading(true);
        setError('');
        setSuccess('');

        try {
            // TODO: 實現添加狙擊目標邏輯
            const newTarget = {
                address: tokenAddress,
                maxBuyAmount,
                minLiquidity,
                status: '等待中'
            };
            setSnipeTargets([...snipeTargets, newTarget]);
            setSuccess('狙擊目標已添加！');
            
            // 清空輸入
            setTokenAddress('');
            setMaxBuyAmount('');
            setMinLiquidity('');
        } catch (error) {
            console.error('添加失敗:', error);
            setError(`添加失敗: ${error.message}`);
        } finally {
            setLoading(false);
        }
    };

    const handleRemoveTarget = (address) => {
        setSnipeTargets(targets => targets.filter(t => t.address !== address));
    };

    return (
        <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
                狙擊開關
            </Typography>

            <Grid container spacing={2}>
                <Grid item xs={12}>
                    <Card sx={{ mb: 2 }}>
                        <CardContent>
                            <Grid container spacing={2}>
                                <Grid item xs={12}>
                                    <TextField
                                        fullWidth
                                        label="目標Token地址"
                                        value={tokenAddress}
                                        onChange={(e) => setTokenAddress(e.target.value)}
                                        disabled={loading}
                                    />
                                </Grid>

                                <Grid item xs={12} sm={6}>
                                    <TextField
                                        fullWidth
                                        type="number"
                                        label="最大買入數量"
                                        value={maxBuyAmount}
                                        onChange={(e) => setMaxBuyAmount(e.target.value)}
                                        disabled={loading}
                                    />
                                </Grid>

                                <Grid item xs={12} sm={6}>
                                    <TextField
                                        fullWidth
                                        type="number"
                                        label="最小流動性要求"
                                        value={minLiquidity}
                                        onChange={(e) => setMinLiquidity(e.target.value)}
                                        disabled={loading}
                                    />
                                </Grid>
                            </Grid>
                        </CardContent>
                    </Card>
                </Grid>

                <Grid item xs={12}>
                    <GasSettingsCard web3={web3} />
                </Grid>

                <Grid item xs={12}>
                    <FormControlLabel
                        control={
                            <Switch
                                checked={enabled}
                                onChange={(e) => setEnabled(e.target.checked)}
                                disabled={loading}
                            />
                        }
                        label={`狙擊功能 ${enabled ? '已啟用' : '已禁用'}`}
                    />
                </Grid>

                {(error || success) && (
                    <Grid item xs={12}>
                        <Alert severity={error ? "error" : "success"}>
                            {error || success}
                        </Alert>
                    </Grid>
                )}

                <Grid item xs={12}>
                    <Button
                        variant="contained"
                        fullWidth
                        onClick={handleAddTarget}
                        disabled={loading || !enabled}
                    >
                        {loading ? '處理中...' : '添加狙擊目標'}
                    </Button>
                </Grid>

                {snipeTargets.length > 0 && (
                    <Grid item xs={12}>
                        <Typography variant="subtitle1" gutterBottom>
                            當前狙擊目標：
                        </Typography>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell>Token地址</TableCell>
                                    <TableCell align="right">最大買入</TableCell>
                                    <TableCell align="right">最小流動性</TableCell>
                                    <TableCell align="right">狀態</TableCell>
                                    <TableCell align="right">操作</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {snipeTargets.map((target) => (
                                    <TableRow key={target.address}>
                                        <TableCell>{target.address.slice(0, 6)}...{target.address.slice(-4)}</TableCell>
                                        <TableCell align="right">{target.maxBuyAmount}</TableCell>
                                        <TableCell align="right">{target.minLiquidity}</TableCell>
                                        <TableCell align="right">{target.status}</TableCell>
                                        <TableCell align="right">
                                            <Button 
                                                size="small" 
                                                color="error"
                                                onClick={() => handleRemoveTarget(target.address)}
                                            >
                                                移除
                                            </Button>
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </Grid>
                )}
            </Grid>
        </Paper>
    );
};

export default SniperSwitch;

//src/components/trading/TokenTrading.js

import React, { useState, useEffect } from 'react';
import { 
    Box, 
    TextField, 
    Button, 
    Paper, 
    Typography, 
    Grid,
    CircularProgress,
    Card,
    CardContent,
    ToggleButton,
    ToggleButtonGroup,
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableRow,
    Alert,
    Tooltip,
    IconButton,
    Link
} from '@mui/material';
import { 
    Info as InfoIcon,
    ContentCopy as CopyIcon,
    OpenInNew as OpenInNewIcon
} from '@mui/icons-material';
import { useWeb3 } from '../../contexts/Web3Context.js';
import { useDex } from '../../contexts/DexContext.js';
import { useGas } from '../../contexts/GasContext.js';
import {
    ERC20_ABI,
    PAIR_ABI,
    FACTORY_ABI,
    CHAIN_TOKENS,
    getTokenInfo,
    getTokenBalance,
    getNativeBalance,
    formatNumber,
    formatAddress
} from '../../utils/tokenUtils.js';
import GasSettingsCard from '../common/GasSettingsCard';
import { getNativeCurrencySymbol, prepareGasSettings } from '../../utils/utils';

const TokenTrading = () => {
    const { web3, selectedChain, accounts } = useWeb3();
    const { selectedDex } = useDex();
    const { updateGasPrice, lastUpdate } = useGas();
    
    const [targetToken, setTargetToken] = useState('');
    const [tokenInfo, setTokenInfo] = useState(null);
    const [loading, setLoading] = useState(false);
    const [pools, setPools] = useState([]);
    const [selectedPool, setSelectedPool] = useState(null);
    const [tradeType, setTradeType] = useState('buy');
    const [amount, setAmount] = useState('');
    const [slippage, setSlippage] = useState('1');
    const [error, setError] = useState('');
    const [walletBalances, setWalletBalances] = useState({});

    // 自動更新
    useEffect(() => {
        if (web3) {
            updateGasPrice(web3);
            if (targetToken && targetToken.length === 42) {
                fetchTokenAndPoolInfo();
            }
        }
    }, [web3, targetToken, selectedDex, lastUpdate]);

    useEffect(() => {
        if (web3 && accounts.length > 0 && targetToken) {
            updateWalletBalances();
        }
    }, [web3, accounts, targetToken, lastUpdate]);

    const updateWalletBalances = async () => {
        try {
            const balances = {};
            const enabledAccounts = accounts.filter(acc => acc.enabled);
            
            await Promise.all(enabledAccounts.map(async (account) => {
                const [nativeBalance, usdtBalance, tokenBalance] = await Promise.all([
                    getNativeBalance(web3, account.address),
                    getTokenBalance(web3, CHAIN_TOKENS[selectedChain].USDT.address, account.address),
                    targetToken ? getTokenBalance(web3, targetToken, account.address) : Promise.resolve('0')
                ]);

                balances[account.address] = {
                    native: nativeBalance,
                    usdt: usdtBalance,
                    token: tokenBalance
                };
            }));

            setWalletBalances(balances);
        } catch (error) {
            console.error('Failed to update balances:', error);
        }
    };

    const fetchTokenAndPoolInfo = async () => {
        if (!web3 || !targetToken || !selectedDex) return;
        
        try {
            const info = await getTokenInfo(web3, targetToken);
            if (!info) {
                throw new Error('無效的代幣地址');
            }
            setTokenInfo(info);

            const factory = new web3.eth.Contract(FACTORY_ABI, selectedDex.factory);
            const pools = [];

            const stableTokens = [
                {
                    address: CHAIN_TOKENS[selectedChain].USDT.address,
                    symbol: 'USDT'
                },
                {
                    address: CHAIN_TOKENS[selectedChain].BUSD.address,
                    symbol: 'BUSD'
                }
            ];

            // 並行獲取所有池子信息
            const poolPromises = [
                // 穩定幣池子
                ...stableTokens.map(async (stableToken) => {
                    try {
                        const pairAddress = await factory.methods.getPair(targetToken, stableToken.address).call();
                        if (pairAddress === '0x0000000000000000000000000000000000000000') {
                            return null;
                        }

                        const pair = new web3.eth.Contract(PAIR_ABI, pairAddress);
                        const [token0, token1, reserves] = await Promise.all([
                            pair.methods.token0().call(),
                            pair.methods.token1().call(),
                            pair.methods.getReserves().call()
                        ]);

                        const isToken0 = token0.toLowerCase() === targetToken.toLowerCase();
                        const tokenReserve = isToken0 ? reserves[0] : reserves[1];
                        const stableReserve = isToken0 ? reserves[1] : reserves[0];
                        
                        return {
                            type: stableToken.symbol,
                            address: pairAddress,
                            tokenReserve,
                            stableReserve,
                            baseTokenAddress: stableToken.address,
                            totalValue: parseFloat(web3.utils.fromWei(stableReserve)) * 2,
                            price: parseFloat(web3.utils.fromWei(stableReserve)) / parseFloat(web3.utils.fromWei(tokenReserve))
                        };
                    } catch (error) {
                        console.error(`Error fetching ${stableToken.symbol} pool:`, error);
                        return null;
                    }
                }),

                // BNB 池子
                (async () => {
                    try {
                        const pairAddress = await factory.methods.getPair(targetToken, CHAIN_TOKENS[selectedChain].WBNB.address).call();
                        if (pairAddress === '0x0000000000000000000000000000000000000000') {
                            return null;
                        }

                        const pair = new web3.eth.Contract(PAIR_ABI, pairAddress);
                        const [token0, token1, reserves] = await Promise.all([
                            pair.methods.token0().call(),
                            pair.methods.token1().call(),
                            pair.methods.getReserves().call()
                        ]);

                        const isToken0 = token0.toLowerCase() === targetToken.toLowerCase();
                        const tokenReserve = isToken0 ? reserves[0] : reserves[1];
                        const bnbReserve = isToken0 ? reserves[1] : reserves[0];
                        const bnbPrice = CHAIN_TOKENS[selectedChain].WBNB.price;

                        return {
                            type: 'BNB',
                            address: pairAddress,
                            tokenReserve,
                            bnbReserve,
                            baseTokenAddress: CHAIN_TOKENS[selectedChain].WBNB.address,
                            totalValue: parseFloat(web3.utils.fromWei(bnbReserve)) * bnbPrice * 2,
                            price: (parseFloat(web3.utils.fromWei(bnbReserve)) * bnbPrice) / parseFloat(web3.utils.fromWei(tokenReserve))
                        };
                    } catch (error) {
                        console.error('Error fetching BNB pool:', error);
                        return null;
                    }
                })()
            ];

            const poolResults = await Promise.all(poolPromises);
            poolResults.forEach(pool => {
                if (pool) pools.push(pool);
            });

            // 按總流動性排序
            pools.sort((a, b) => b.totalValue - a.totalValue);
            setPools(pools);

            // 保持當前選擇或選擇最大流動性池子
            if (selectedPool) {
                const currentPool = pools.find(p => p.address === selectedPool.address);
                if (currentPool) {
                    setSelectedPool(currentPool);
                    return;
                }
            }
            
            if (pools.length > 0) {
                setSelectedPool(pools[0]);
            } else {
                setSelectedPool(null);
            }

        } catch (error) {
            console.error('Error:', error);
            setError(error.message);
            setTokenInfo(null);
            setPools([]);
            setSelectedPool(null);
        }
    };

    const getExplorerUrl = (address) => {
        const explorers = {
            'ETH': 'https://etherscan.io',
            'BSC': 'https://bscscan.com',
            'CORE': 'https://scan.coredao.org',
            'BASE': 'https://basescan.org'
        };
        return `${explorers[selectedChain]}/address/${address}`;
    };

    const calculateExpectedOutput = (inputAmount) => {
        if (!inputAmount || !selectedPool) return '0';
        const amount = parseFloat(inputAmount);
        return tradeType === 'buy' 
            ? amount / selectedPool.price
            : amount * selectedPool.price;
    };

    const executeTrade = async () => {
        if (!selectedPool || !amount || !web3) return;

        try {
            setLoading(true);
            const enabledAccounts = accounts.filter(acc => acc.enabled);
            if (enabledAccounts.length === 0) {
                throw new Error('沒有啟用的錢包');
            }

            const amountIn = web3.utils.toWei(amount);
            const minOutputWei = web3.utils.toWei(
                (parseFloat(calculateExpectedOutput(amount)) * (1 - parseFloat(slippage)/100)).toString()
            );

            const deadline = Math.floor(Date.now() / 1000) + 60 * 20;

            for (const account of enabledAccounts) {
                if (tradeType === 'buy') {
                    let tx = {
                        from: account.address,
                        to: selectedDex.router,
                        value: amountIn,
                        data: web3.eth.abi.encodeFunctionCall({
                            name: 'swapExactETHForTokens',
                            type: 'function',
                            inputs: [
                                { type: 'uint256', name: 'amountOutMin' },
                                { type: 'address[]', name: 'path' },
                                { type: 'address', name: 'to' },
                                { type: 'uint256', name: 'deadline' }
                            ]
                        }, [
                            minOutputWei,
                            [CHAIN_TOKENS[selectedChain].WBNB.address, targetToken],
                            account.address,
                            deadline
                        ])
                    };

                    tx = prepareGasSettings(web3, tx);
                    const signedTx = await web3.eth.accounts.signTransaction(tx, account.privateKey);
                    await web3.eth.sendSignedTransaction(signedTx.rawTransaction);
                } else {
                    const tokenContract = new web3.eth.Contract(ERC20_ABI, targetToken);
                    
                    // 檢查授權
                    const allowance = await tokenContract.methods.allowance(account.address, selectedDex.router).call();
                    if (web3.utils.toBN(allowance).lt(web3.utils.toBN(amountIn))) {
                        let approveTx = {
                            from: account.address,
                            to: targetToken,
                            data: tokenContract.methods.approve(
                                selectedDex.router,
                                '115792089237316195423570985008687907853269984665640564039457584007913129639935'
                            ).encodeABI()
                        };

                        approveTx = prepareGasSettings(web3, approveTx);
                        const signedApproveTx = await web3.eth.accounts.signTransaction(approveTx, account.privateKey);
                        await web3.eth.sendSignedTransaction(signedApproveTx.rawTransaction);
                    }

                    let tx = {
                        from: account.address,
                        to: selectedDex.router,
                        data: web3.eth.abi.encodeFunctionCall({
                            name: 'swapExactTokensForETH',
                            type: 'function',
                            inputs: [
                                { type: 'uint256', name: 'amountIn' },
                                { type: 'uint256', name: 'amountOutMin' },
                                { type: 'address[]', name: 'path' },
                                { type: 'address', name: 'to' },
                                { type: 'uint256', name: 'deadline' }
                            ]
                        }, [
                            amountIn,
                            minOutputWei,
                            [targetToken, CHAIN_TOKENS[selectedChain].WBNB.address],
                            account.address,
                            deadline
                        ])
                    };

                    tx = prepareGasSettings(web3, tx);
                    const signedTx = await web3.eth.accounts.signTransaction(tx, account.privateKey);
                    await web3.eth.sendSignedTransaction(signedTx.rawTransaction);
                }
            }

            alert('交易成功完成！');
            setAmount('');
            await Promise.all([
                fetchTokenAndPoolInfo(),
                updateWalletBalances()
            ]);

        } catch (error) {
            console.error('Transaction failed:', error);
            alert(`交易失敗: ${error.message}`);
        } finally {
            setLoading(false);
        }
    };

    const nativeCurrency = getNativeCurrencySymbol(selectedChain);

    return (
        <Paper sx={{ p: 2, mb: 2 }}>
            <Typography variant="h6" gutterBottom>
                Token 交易
            </Typography>

            <Grid container spacing={2}>
                <Grid item xs={12}>
                    <TextField
                        fullWidth
                        label="目標Token地址"
                        value={targetToken}
                        onChange={(e) => setTargetToken(e.target.value)}
                        error={!!error}
                        helperText={error || (tokenInfo && `${tokenInfo.name} (${tokenInfo.symbol})`)}
                        disabled={loading}
                    />
                </Grid>

                {accounts.filter(acc => acc.enabled).length > 0 && (
                    <Grid item xs={12}>
                        <Card>
                            <CardContent>
                                <Typography variant="subtitle1" gutterBottom>
                                    錢包餘額：
                                </Typography>
                                <Table size="small">
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>錢包地址</TableCell>
                                            <TableCell align="right">BNB</TableCell>
                                            <TableCell align="right">USDT</TableCell>
                                            {tokenInfo && (
                                                <TableCell align="right">{tokenInfo.symbol}</TableCell>
                                            )}
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {accounts.filter(acc => acc.enabled).map((account) => (
                                            <TableRow key={account.address}>
                                                <TableCell>
                                                    <Link
                                                        href={getExplorerUrl(account.address)}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        sx={{ display: 'flex', alignItems: 'center' }}
                                                    >
                                                        {formatAddress(account.address)}
                                                        <OpenInNewIcon sx={{ ml: 0.5, fontSize: 16 }} />
                                                    </Link>
                                                </TableCell>
                                                <TableCell align="right">
                                                    {walletBalances[account.address]?.native || '0'}
                                                </TableCell>
                                                <TableCell align="right">
                                                    {walletBalances[account.address]?.usdt || '0'}
                                                </TableCell>
                                                {tokenInfo && (
                                                    <TableCell align="right">
                                                        {walletBalances[account.address]?.token || '0'}
                                                    </TableCell>
                                                )}
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </CardContent>
                        </Card>
                    </Grid>
                )}

                {tokenInfo && pools.length > 0 && (
                    <Grid item xs={12}>
                        <Card>
                            <CardContent>
                                <Typography variant="subtitle1" gutterBottom>
                                    可用流動池（按總流動性排序）：
                                </Typography>
                                <Table size="small">
                                    <TableHead>
                                        <TableRow>
                                            <TableCell>池子類型</TableCell>
                                            <TableCell align="right">Token 數量</TableCell>
                                            <TableCell align="right">底池數量</TableCell>
                                            <TableCell align="right">總流動性 (USD)</TableCell>
                                            <TableCell align="right">Token 價格</TableCell>
                                            <TableCell>操作</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {pools.map((pool) => (
                                            <TableRow 
                                                key={pool.address}
                                                selected={selectedPool?.address === pool.address}
                                                hover
                                                onClick={() => setSelectedPool(pool)}
                                                sx={{ cursor: 'pointer' }}
                                            >
                                                <TableCell>
                                                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                                                        {pool.type}
                                                        <Tooltip title="查看池子">
                                                            <IconButton 
                                                                size="small"
                                                                onClick={(e) => {
                                                                    e.stopPropagation();
                                                                    window.open(getExplorerUrl(pool.address), '_blank');
                                                                }}
                                                            >
                                                                <OpenInNewIcon fontSize="small" />
                                                            </IconButton>
                                                        </Tooltip>
                                                    </Box>
                                                </TableCell>
                                                <TableCell align="right">
                                                    {formatNumber(web3.utils.fromWei(pool.tokenReserve))}
                                                </TableCell>
                                                <TableCell align="right">
                                                    {formatNumber(web3.utils.fromWei(pool.type === 'BNB' ? pool.bnbReserve : pool.stableReserve))}
                                                </TableCell>
                                                <TableCell align="right">
                                                    ${formatNumber(pool.totalValue)}
                                                </TableCell>
                                                <TableCell align="right">
                                                    ${formatNumber(pool.price, 6)}
                                                </TableCell>
                                                <TableCell>
                                                    <Button
                                                        variant={selectedPool?.address === pool.address ? "contained" : "outlined"}
                                                        size="small"
                                                    >
                                                        選擇
                                                    </Button>
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </CardContent>
                        </Card>
                    </Grid>
                )}

                {selectedPool && (
                    <>
                        <Grid item xs={12}>
                            <ToggleButtonGroup
                                value={tradeType}
                                exclusive
                                onChange={(e, newValue) => newValue && setTradeType(newValue)}
                                fullWidth
                            >
                                <ToggleButton value="buy">
                                    用 BNB 買入 {tokenInfo?.symbol}
                                </ToggleButton>
                                <ToggleButton value="sell">
                                    賣出 {tokenInfo?.symbol} 換 BNB
                                </ToggleButton>
                            </ToggleButtonGroup>
                        </Grid>

                        <Grid item xs={8}>
                            <TextField
                                fullWidth
                                label={`${tradeType === 'buy' ? '購買' : '賣出'}數量`}
                                type="number"
                                value={amount}
                                onChange={(e) => setAmount(e.target.value)}
                                disabled={loading}
                            />
                        </Grid>

                        <Grid item xs={4}>
                            <TextField
                                fullWidth
                                label="滑點 %"
                                type="number"
                                value={slippage}
                                onChange={(e) => setSlippage(e.target.value)}
                                inputProps={{ min: 0, max: 100 }}
                                disabled={loading}
                            />
                        </Grid>

                        {amount && (
                            <Grid item xs={12}>
                                <Card>
                                    <CardContent>
                                        <Typography variant="subtitle2" gutterBottom>
                                            交易預覽：
                                        </Typography>
                                        <Grid container spacing={2}>
                                            <Grid item xs={6}>
                                                <Typography variant="body2" color="textSecondary">
                                                    支付：
                                                </Typography>
                                                <Typography variant="body1">
                                                    {amount} {tradeType === 'buy' ? 'BNB' : tokenInfo?.symbol}
                                                </Typography>
                                            </Grid>
                                            <Grid item xs={6}>
                                                <Typography variant="body2" color="textSecondary">
                                                    預計獲得：
                                                </Typography>
                                                <Typography variant="body1">
                                                    {formatNumber(calculateExpectedOutput(amount))} {tradeType === 'buy' ? tokenInfo?.symbol : 'BNB'}
                                                </Typography>
                                            </Grid>
                                            <Grid item xs={6}>
                                                <Typography variant="body2" color="textSecondary">
                                                    最小獲得（含滑點）：
                                                </Typography>
                                                <Typography variant="body1">
                                                    {formatNumber(parseFloat(calculateExpectedOutput(amount)) * (1 - parseFloat(slippage)/100))}
                                                </Typography>
                                            </Grid>
                                            <Grid item xs={6}>
                                                <Typography variant="body2" color="textSecondary">
                                                    價格：
                                                </Typography>
                                                <Typography variant="body1">
                                                    ${formatNumber(selectedPool.price, 6)}
                                                </Typography>
                                            </Grid>
                                        </Grid>
                                    </CardContent>
                                </Card>
                            </Grid>
                        )}

                        <Grid item xs={12}>
                            <GasSettingsCard 
                                web3={web3}
                                nativeCurrency={nativeCurrency}
                            />
                        </Grid>

                        <Grid item xs={12}>
                            <Button
                                variant="contained"
                                fullWidth
                                size="large"
                                onClick={executeTrade}
                                disabled={!amount || parseFloat(amount) <= 0 || loading}
                            >
                                {loading ? '交易處理中...' : `${tradeType === 'buy' ? '買入' : '賣出'} ${tokenInfo?.symbol}`}
                            </Button>
                        </Grid>
                    </>
                )}
            </Grid>
        </Paper>
    );
};

export default TokenTrading;



//src/components/trading/TradingPage.js

// src/components/trading/TradingPage.js
import React, { useState } from 'react';
import { 
    Box, 
    Paper, 
    Tabs, 
    Tab,
    Typography,
    Container
} from '@mui/material';
import TokenTrading from './TokenTrading';
import TokenApproval from './TokenApproval';
import HexTransaction from './HexTransaction';
import DexManager from '../DexManager';
import SniperPool from './SniperPool';
import SniperSwitch from './SniperSwitch';

// Tab Panel component
function TabPanel({ children, value, index, ...other }) {
    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`trading-tabpanel-${index}`}
            aria-labelledby={`trading-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box sx={{ py: 3 }}>
                    {children}
                </Box>
            )}
        </div>
    );
}

const TradingPage = () => {
    const [mainTab, setMainTab] = useState(0);
    const [normalTradeTab, setNormalTradeTab] = useState(0);
    const [sniperTab, setSniperTab] = useState(0);

    return (
        <Container maxWidth="lg">
            {/* DEX Manager */}
            <DexManager />

            {/* Main Trading Section */}
            <Paper sx={{ mb: 3 }}>
                <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                    <Tabs 
                        value={mainTab} 
                        onChange={(_, newValue) => setMainTab(newValue)}
                        aria-label="main trading tabs"
                    >
                        <Tab label="一般交易" />
                        <Tab label="狙擊模式" />
                    </Tabs>
                </Box>

                {/* Normal Trading Panel */}
                <TabPanel value={mainTab} index={0}>
                    <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                        <Tabs 
                            value={normalTradeTab} 
                            onChange={(_, newValue) => setNormalTradeTab(newValue)}
                            aria-label="normal trading tabs"
                        >
                            <Tab label="代幣交易" />
                            <Tab label="代幣授權" />
                            <Tab label="自定義交易" />
                        </Tabs>
                    </Box>
                    <TabPanel value={normalTradeTab} index={0}>
                        <TokenTrading />
                    </TabPanel>
                    <TabPanel value={normalTradeTab} index={1}>
                        <TokenApproval />
                    </TabPanel>
                    <TabPanel value={normalTradeTab} index={2}>
                        <HexTransaction />
                    </TabPanel>
                </TabPanel>

                {/* Sniper Mode Panel */}
                <TabPanel value={mainTab} index={1}>
                    <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                        <Tabs 
                            value={sniperTab} 
                            onChange={(_, newValue) => setSniperTab(newValue)}
                            aria-label="sniper mode tabs"
                        >
                            <Tab label="狙擊加池" />
                            <Tab label="狙擊開關" />
                        </Tabs>
                    </Box>
                    <TabPanel value={sniperTab} index={0}>
                        <SniperPool />
                    </TabPanel>
                    <TabPanel value={sniperTab} index={1}>
                        <SniperSwitch />
                    </TabPanel>
                </TabPanel>
            </Paper>
        </Container>
    );
};

export default TradingPage;

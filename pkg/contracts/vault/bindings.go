//go:generate sh -c "solc UpgradeableVault.sol --evm-version london --combined-json abi,bin | jq '.contracts.\"UpgradeableVault.sol:UpgradeableVault\"'  > UpgradeableVault.json"
//go:generate sh -c "cat UpgradeableVault.json | jq .abi > UpgradeableVault.abi"
//go:generate sh -c "cat UpgradeableVault.json | jq .bin  | tr -d '\"'  > UpgradeableVault.bin"
//go:generate sh -c "abigen --abi UpgradeableVault.abi --bin UpgradeableVault.bin  --pkg vault --type UpgradeableVault --out UpgradeableVault.go"

package vault

var _ UpgradeableVault

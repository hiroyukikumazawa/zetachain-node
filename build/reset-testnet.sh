set -e
set -x
ZETACORED=/home/ubuntu/go/bin/zetacored
NODES="3.137.46.147"

rm -rf ~/.zetacore
for NODE in $NODES; do
	ssh -i ~/.ssh/meta.pem $NODE rm -rf ~/.zetacore
done

$ZETACORED init --chain-id testing zetachain
$ZETACORED keys add val --keyring-backend=test
$ZETACORED add-genesis-account $($ZETACORED keys show val -a --keyring-backend=test) 1000000000stake


for NODE in $NODES; do
	ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys add val --keyring-backend=test
	ADDR=$(ssh -i ~/.ssh/meta.pem $NODE $ZETACORED keys show val -a --keyring-backend=test)
	$ZETACORED add-genesis-account $ADDR 1000000000stake --keyring-backend=test
done


for NODE in $NODES; do
	scp -i ~/.ssh/meta.pem ~/.zetacore/config/genesis.json $NODE:~/.zetacore/config/
done


$ZETACORED gentx val 1000000000stake --keyring-backend=test --chain-id=testing
for NODE in $NODES; do
    ssh -i ~/.ssh/meta.pem $NODE $ZETACORED gentx val 1000000000stake --keyring-backend=test --chain-id=testing -ip $NODE
    scp -i ~/.ssh/meta.pem $NODE:~/.zetacore/config/gentx/*.json ~/.zetacore/config/gentx/
done

$ZETACORED collect-gentxs


for NODE in $NODES; do
	scp -i ~/.ssh/meta.pem ~/.zetacore/config/genesis.json $NODE:~/.zetacore/config/
done


jq '.chain_id = "testing"' ~/.zetacore/config/genesis.json > temp.json && mv temp.json ~/.zetacore/config/genesis.json
sed -i '/\[api\]/,+3 s/enable = false/enable = true/' ~/.zetacore/config/app.toml
sed -i '/\[api\]/,+24 s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' ~/.zetacore/config/app.toml
for NODE in $NODES; do
    ssh -i ~/.ssh/meta.pem $NODE jq \'.chain_id = \"testing\"\' ~/.zetacore/config/genesis.json > temp.json && mv temp.json ~/.zetacore/config/genesis.json
    ssh -i ~/.ssh/meta.pem $NODE sed -i \'/\[api\]/,+3 s/enable = false/enable = true/\' ~/.zetacore/config/app.toml
    ssh -i ~/.ssh/meta.pem $NODE sed -i \'/\[api\]/,+24 s/enabled-unsafe-cors = false/enabled-unsafe-cors = true/\' ~/.zetacore/config/app.toml
done


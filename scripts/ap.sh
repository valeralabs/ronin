echo "Using SG Bitcoin node"
toml set ./configurations/mainnet/Config.toml burnchain.peer_host "95.111.194.148" > temp
mv temp ./configurations/mainnet/Config.toml
# Ubuntu
# There isn't a Heroku add on for RethinkDB, so it has to be self hosted
# or, through Compose.io, which is 30$ a month
# Digital Ocean is 5$ for a small box
source /etc/lsb-release && echo "deb http://download.rethinkdb.com/apt
$DISTRIB_CODENAME main" | sudo tee /etc/apt/sources.list.d/rethinkdb.list

wget -qO- http://download.rethinkdb.com/apt/pubkey.gpg | sudo apt-key add -

sudo apt-get update
sudo apt-get install rethinkdb


##LEVANTAR EL SERVIDOR
cd API
export GOROOT=/usr/local/go
export GOPATH=$HOME/Projects/Proj1
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
echo ">> Version de Go"
go version
echo ">> Obteniendo librerias necesarias"
echo ">> GinGonic"

go get -u github.com/gorilla/mux
go get -u github.com/tidwall/gjson

echo ">> Construyendo aplicacion"
go run main.go
#go build -o ___go_build_so_p_02_ .
echo ">> Iniciando servidor"
#./___go_build_so_p_02_

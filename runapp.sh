
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


echo ">> Construyendo aplicacion"
go build -o ___go_build_so_p_02_ .
echo ">> Iniciando servidor"
./___go_build_so_p_02_

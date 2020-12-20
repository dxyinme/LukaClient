go build -o bin\gui.exe main\gui.go
copy start.cmd bin\
cd bin
rd ClientExample
mkdir ClientExample
mkdir SaveTmp
cd ..
xcopy ClientExample bin\ClientExample /S/E/A

$cwd = Convert-Path .
$gopathStr = $env:GOPATH

Write-Host "`$cwd: $cwd"
Write-Host "`$gopathStr: $gopathStr"

if ( $gopathStr -like "*$cwd*" ) {
    Write-Host "GOPATH contains cwd"
} else {
    Write-Host "GOPATH does not contain cwd"
    $env:GOPATH = "$env:GOPATH;$cwd;"
}

Write-Host "Build GOPATH: $env:GOPATH"



Write-Host "go build..." 
go build -a -o gocodevis.exe .\src\main.go

# Write-Host "go test" 
# go test .\src\types .\src\assembly

Write-Host "Generating project diagram text..."
.\gocodevis -output ./projectdiagrams/ -root ./src 

Write-Host "Generating project diagrams..."
java -DPLANTUML_LIMIT_SIZE=11264 -jar ./plantuml/plantuml.1.2019.11.jar  "./projectdiagrams" 
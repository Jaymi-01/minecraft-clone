# Mine RPG Multi-Platform Build Script

# 1. Cleanup
Write-Host 'Cleaning up old builds...'
If (Test-Path 'Release') { Remove-Item 'Release' -Recurse -Force }
$OldZips = Get-ChildItem -Filter 'MineRPG_*.zip'
if ($OldZips) { $OldZips | Remove-Item -Force }
New-Item -ItemType Directory -Path 'Release'

# 2. Build Targets
$Targets = @(
    @{ OS = 'windows'; Arch = 'amd64'; Suffix = '.exe'; Name = 'Windows' },
    @{ OS = 'linux';   Arch = 'amd64'; Suffix = '';     Name = 'Linux' },
    @{ OS = 'darwin';  Arch = 'amd64'; Suffix = '';     Name = 'Mac_Intel' },
    @{ OS = 'darwin';  Arch = 'arm64'; Suffix = '';     Name = 'Mac_M1_M2' }
)

# 3. Compile and Zip each target
foreach ($T in $Targets) {
    $OS = $T.OS
    $Arch = $T.Arch
    $PlatformName = $T.Name
    $BinaryName = 'mine-system' + $T.Suffix
    $ZipName = 'MineRPG_' + $PlatformName + '.zip'

    Write-Host ('Building for ' + $PlatformName + '...')
    
    $env:GOOS = $OS
    $env:GOARCH = $Arch
    
    $TempPath = 'Release/' + $PlatformName
    New-Item -ItemType Directory -Path $TempPath -Force | Out-Null
    go build -o ($TempPath + '/' + $BinaryName) .
    
    Copy-Item 'README.txt' ($TempPath + '/')

    Write-Host ('Packaging ' + $ZipName + '...')
    Compress-Archive -Path ($TempPath + '/*') -DestinationPath $ZipName -Force
}

# 4. Final Cleanup
Remove-Item 'Release' -Recurse -Force

Write-Host 'ALL PLATFORMS READY!'

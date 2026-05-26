# Mine RPG Multi-Platform Build Script

# 1. Cleanup
Write-Host 'Cleaning up old builds...'
$ReleaseDir = 'release'
$TempStage = 'build_temp'

If (Test-Path $TempStage) { Remove-Item $TempStage -Recurse -Force }
If (Test-Path $ReleaseDir) { Remove-Item $ReleaseDir -Recurse -Force }
New-Item -ItemType Directory -Path $ReleaseDir
New-Item -ItemType Directory -Path $TempStage

# 2. Build Targets
$Targets = @(
    @{ OS = 'windows'; Arch = 'amd64'; Suffix = '.exe'; Name = 'Windows' },
    @{ OS = 'linux';   Arch = 'amd64'; Suffix = '';     Name = 'Linux' },
    @{ OS = 'darwin';  Arch = 'amd64'; Suffix = '';     Name = 'Mac_Intel' },
    @{ OS = 'darwin';  Arch = 'arm64'; Suffix = '';     Name = 'Mac_M1_M2' },
    @{ OS = 'android'; Arch = 'arm64'; Suffix = '';     Name = 'Android' }
)

# 3. Compile and Zip each target
foreach ($T in $Targets) {
    $OS = $T.OS
    $Arch = $T.Arch
    $PlatformName = $T.Name
    $BinaryName = 'mine-system' + $T.Suffix
    $ZipName = Join-Path $ReleaseDir ('MineRPG_' + $PlatformName + '.zip')

    Write-Host ('Building for ' + $PlatformName + '...')
    
    $env:GOOS = $OS
    $env:GOARCH = $Arch
    
    $TargetTempPath = Join-Path $TempStage $PlatformName
    New-Item -ItemType Directory -Path $TargetTempPath -Force | Out-Null
    go build -o (Join-Path $TargetTempPath $BinaryName) .
    
    if (Test-Path 'README.txt') {
        Copy-Item 'README.txt' $TargetTempPath
    }

    Write-Host ('Packaging ' + $ZipName + '...')
    Compress-Archive -Path (Join-Path $TargetTempPath '*') -DestinationPath $ZipName -Force
}

# 4. Final Cleanup
Remove-Item $TempStage -Recurse -Force

Write-Host 'ALL PLATFORMS READY IN THE RELEASE FOLDER!'

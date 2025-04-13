# LazyNode Windows Installation Script

# Ensure we're running as administrator
if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Warning "This script requires administrator privileges. Please run as administrator."
    exit
}

# Determine architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# Get the latest release version from GitHub API
try {
    $latestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/VesperAkshay/lazynode/releases/latest" -ErrorAction Stop
    $version = $latestRelease.tag_name.Replace("v", "")
    Write-Host "Latest version: $version"
} catch {
    Write-Error "Failed to get latest version. Error: $_"
    exit 1
}

# Define paths
$downloadUrl = "https://github.com/VesperAkshay/lazynode/releases/download/v$version/lazynode_${version}_windows_${arch}.zip"
$tempFile = "$env:TEMP\lazynode.zip"
$installDir = "$env:ProgramFiles\LazyNode"

# Create install directory if it doesn't exist
if (-not (Test-Path -Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

# Download the zip file
Write-Host "Downloading LazyNode v$version for Windows $arch..."
try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -ErrorAction Stop
} catch {
    Write-Error "Failed to download file. Error: $_"
    exit 1
}

# Extract the zip file
Write-Host "Extracting files..."
Expand-Archive -Path $tempFile -DestinationPath $installDir -Force

# Add to PATH if not already there
$currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($currentPath -notlike "*$installDir*") {
    Write-Host "Adding LazyNode to PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$installDir", "Machine")
}

# Clean up
Remove-Item $tempFile -Force

Write-Host "LazyNode v$version has been installed successfully!"
Write-Host "You can now run 'lazynode' from any PowerShell or Command Prompt window."
Write-Host "Note: You may need to open a new terminal window for the PATH changes to take effect." 
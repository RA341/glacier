$repo = "RA341/glacier"
$assetName = "frost-windows.zip"
$downloadUrl = "https://github.com/$repo/releases/download/canary/$assetName"

$targetDir = Get-Location # Default to current working directory

Write-Host "Current download directory is: $targetDir"
$choice = Read-Host "Do you want to manually pick a different directory? (y/n)"

if ($choice -eq 'y')
{
    Add-Type -AssemblyName System.Windows.Forms
    $FolderBrowser = New-Object System.Windows.Forms.FolderBrowserDialog
    $FolderBrowser.Description = "Select a folder to download and extract Frost"

    # Ensure the dialog opens on top of other windows
    $result = $FolderBrowser.ShowDialog((New-Object System.Windows.Forms.Form -Property @{ TopMost = $true }))

    if ($result -eq "OK")
    {
        $targetDir = $FolderBrowser.SelectedPath
    }
    else
    {
        Write-Host "Selection cancelled. Staying with: $targetDir" -ForegroundColor Yellow
    }
}

$zipPath = Join-Path $targetDir $assetName
$extractPath = Join-Path $targetDir "frost"

Write-Host "Downloading $assetName to $targetDir..." -ForegroundColor Cyan
try
{
    Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -ErrorAction Stop
}
catch
{
    Write-Error "Failed to download file. Please check your internet or the repository URL."
    return
}

Write-Host "Unzipping to folder: $extractPath..." -ForegroundColor Cyan
if (!(Test-Path $extractPath))
{
    New-Item -ItemType Directory -Path $extractPath | Out-Null
}

try
{
    Expand-Archive -Path $zipPath -DestinationPath $extractPath -Force
    Write-Host "Extraction complete." -ForegroundColor Green
}
catch
{
    Write-Error "Failed to unzip the file."
    return
}

# 5. Ask to delete the zip
$deleteChoice = Read-Host "Would you like to delete the downloaded zip file? (y/n)"
if ($deleteChoice -eq 'y')
{
    Remove-Item -Path $zipPath
    Write-Host "Zip file deleted." -ForegroundColor Gray
}
else
{
    Write-Host "Zip file kept at: $zipPath"
}

Write-Host "Done!" -ForegroundColor Green
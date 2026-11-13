param(
    [Parameter(Mandatory = $true)][string]$IdentityName,
    [Parameter(Mandatory = $true)][string]$Publisher,
    [Parameter(Mandatory = $true)][string]$PublisherDisplayName,
    [string]$Version = "1.0.0.0",
    [string]$OutputPath = ""
)

$ErrorActionPreference = "Stop"

if ($IdentityName -notmatch '^[A-Za-z0-9.-]{3,50}$') {
    throw "IdentityName must match the package identity assigned by Partner Center."
}
if ($Version -notmatch '^\d{1,5}(\.\d{1,5}){3}$') {
    throw "Version must use four numeric parts, for example 1.0.0.0."
}

$projectRoot = [System.IO.Path]::GetFullPath((Join-Path $PSScriptRoot "..\.."))
$exePath = Join-Path $projectRoot "build\giligili.exe"
if (-not (Test-Path $exePath)) { throw "Build the Windows app first: $exePath" }

if ([string]::IsNullOrWhiteSpace($OutputPath)) {
    $OutputPath = Join-Path $projectRoot "dist\Giligili-Store.msix"
}
$OutputPath = [System.IO.Path]::GetFullPath($OutputPath)
$stage = Join-Path ([System.IO.Path]::GetTempPath()) ("giligili-msix-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Force -Path $stage | Out-Null
New-Item -ItemType Directory -Force -Path (Join-Path $stage "Assets") | Out-Null
try {
    Copy-Item $exePath (Join-Path $stage "giligili.exe")
    Copy-Item (Join-Path $projectRoot "web") $stage -Recurse
    Copy-Item (Join-Path $projectRoot "swagger") $stage -Recurse
    Copy-Item (Join-Path $projectRoot "conf") $stage -Recurse

    Add-Type -AssemblyName System.Drawing
    foreach ($spec in @(@{ Size = 50; Name = "StoreLogo.png" }, @{ Size = 150; Name = "Square150x150Logo.png" }, @{ Size = 44; Name = "Square44x44Logo.png" })) {
        $size = [int]$spec.Size
        $bitmap = [System.Drawing.Bitmap]::new($size, $size)
        $graphics = [System.Drawing.Graphics]::FromImage($bitmap)
        $graphics.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
        $graphics.Clear([System.Drawing.Color]::FromArgb(17, 19, 24))
        $fontSize = [single]($size * 0.62)
        $font = [System.Drawing.Font]::new("Arial", $fontSize, [System.Drawing.FontStyle]::Bold, [System.Drawing.GraphicsUnit]::Pixel)
        $brush = [System.Drawing.SolidBrush]::new([System.Drawing.Color]::FromArgb(244, 163, 64))
        $format = [System.Drawing.StringFormat]::new()
        $format.Alignment = [System.Drawing.StringAlignment]::Center
        $format.LineAlignment = [System.Drawing.StringAlignment]::Center
        $rect = [System.Drawing.RectangleF]::new(0, 0, $size, $size)
        $graphics.DrawString("G", $font, $brush, $rect, $format)
        $bitmap.Save((Join-Path (Join-Path $stage "Assets") $spec.Name), [System.Drawing.Imaging.ImageFormat]::Png)
        $format.Dispose()
        $brush.Dispose()
        $font.Dispose()
        $graphics.Dispose()
        $bitmap.Dispose()
    }

    $identityXml = [System.Security.SecurityElement]::Escape($IdentityName)
    $publisherXml = [System.Security.SecurityElement]::Escape($Publisher)
    $publisherDisplayXml = [System.Security.SecurityElement]::Escape($PublisherDisplayName)
    $manifest = @"
<?xml version="1.0" encoding="utf-8"?>
<Package xmlns="http://schemas.microsoft.com/appx/manifest/foundation/windows10" xmlns:uap="http://schemas.microsoft.com/appx/manifest/uap/windows10" xmlns:rescap="http://schemas.microsoft.com/appx/manifest/foundation/windows10/restrictedcapabilities" IgnorableNamespaces="uap rescap">
  <Identity Name="$identityXml" Publisher="$publisherXml" Version="$Version" ProcessorArchitecture="x64" />
  <Properties>
    <DisplayName>G站视频</DisplayName>
    <PublisherDisplayName>$publisherDisplayXml</PublisherDisplayName>
    <Logo>Assets\StoreLogo.png</Logo>
  </Properties>
  <Dependencies>
    <TargetDeviceFamily Name="Windows.Desktop" MinVersion="10.0.19041.0" MaxVersionTested="10.0.26100.0" />
  </Dependencies>
  <Resources>
    <Resource Language="zh-CN" />
  </Resources>
  <Applications>
    <Application Id="Giligili" Executable="giligili.exe" EntryPoint="Windows.FullTrustApplication">
      <uap:VisualElements DisplayName="G站视频" Description="视频浏览和投稿应用" BackgroundColor="#111318" Square150x150Logo="Assets\Square150x150Logo.png" Square44x44Logo="Assets\Square44x44Logo.png" />
    </Application>
  </Applications>
  <Capabilities>
    <rescap:Capability Name="runFullTrust" />
  </Capabilities>
</Package>
"@
    [System.IO.File]::WriteAllText((Join-Path $stage "AppxManifest.xml"), $manifest, [System.Text.UTF8Encoding]::new($false))

    $sdkRoot = Join-Path ([Environment]::GetFolderPath("ProgramFilesX86")) "Windows Kits\10\bin"
    $makeAppx = Get-ChildItem -Path $sdkRoot -Filter MakeAppx.exe -Recurse -ErrorAction Stop |
        Where-Object { $_.FullName -match "\\x64\\MakeAppx\.exe$" } |
        Sort-Object FullName |
        Select-Object -Last 1 -ExpandProperty FullName
    if (-not $makeAppx) { throw "Windows SDK MakeAppx was not found." }

    New-Item -ItemType Directory -Force -Path (Split-Path $OutputPath -Parent) | Out-Null
    & $makeAppx pack /d $stage /p $OutputPath /o
    if ($LASTEXITCODE -ne 0) { throw "MakeAppx failed to create the Store package." }
    & $makeAppx verify /p $OutputPath
    if ($LASTEXITCODE -ne 0) { throw "MakeAppx failed to verify the Store package." }
    Write-Output "Created unsigned Store submission package: $OutputPath"
}
finally {
    if (Test-Path $stage) { Remove-Item $stage -Recurse -Force }
}

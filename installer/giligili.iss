#define AppName "G站视频"
#define AppVersion "1.0.0"
#define AppPublisher "Giligili"

[Setup]
AppId={{E89A7DE7-6900-4FDB-A5AE-D5A19191F501}
AppName={#AppName}
AppVersion={#AppVersion}
AppPublisher={#AppPublisher}
DefaultDirName={localappdata}\Programs\Giligili
DefaultGroupName={#AppName}
PrivilegesRequired=lowest
ArchitecturesInstallIn64BitMode=x64
OutputDir=output
OutputBaseFilename=Giligili-Windows-Installer
WizardStyle=modern
UninstallDisplayIcon={app}\giligili.exe
CloseApplications=yes
Compression=lzma2
SolidCompression=yes

[Languages]
Name: "chinesesimp"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "..\build\giligili.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\web\*"; DestDir: "{app}\web"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "..\swagger\*"; DestDir: "{app}\swagger"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "..\conf\locales\*"; DestDir: "{app}\conf\locales"; Flags: ignoreversion recursesubdirs createallsubdirs

[Icons]
Name: "{autoprograms}\G站视频"; Filename: "{app}\giligili.exe"; WorkingDir: "{app}"
Name: "{autodesktop}\G站视频"; Filename: "{app}\giligili.exe"; WorkingDir: "{app}"; Tasks: desktopicon

[Tasks]
Name: "desktopicon"; Description: "创建桌面快捷方式"; GroupDescription: "附加快捷方式："; Flags: unchecked

[Run]
Filename: "{app}\giligili.exe"; Description: "启动 G站视频"; WorkingDir: "{app}"; Flags: postinstall nowait skipifsilent

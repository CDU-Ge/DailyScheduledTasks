Set-Location $PSScriptRoot

# 读取.env文件内容
Get-Content .env | ForEach-Object {
    # 分割键和值
    $key, $value = $_ -split '=', 2

    # 将变量添加到环境中
    [System.Environment]::SetEnvironmentVariable($key, $value, [System.EnvironmentVariableTarget]::Process)
}

# 输出已加载的环境变量
$env

# update 
if ($args -contains "-U") {
    go build -o .\joke\joke.exe .\joke\joke.go
} 

# run it
.\joke\joke.exe
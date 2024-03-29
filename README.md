# Scabiosa Backup Tool
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fnetbenix%2Fscabiosa.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fnetbenix%2Fscabiosa?ref=badge_shield)


Please keep in mind that this project is WIP.

YouTrack Issues & Agile Board: [Click me!](https://codenoodles.youtrack.cloud/issues)

## What can it do?
- Backup you stuff via a dynamic configuration
- Log the Backup progress to a database 
- Upload the files to a remote storage of your choice (see [Storage Types](#storage-types))

## Planned features for the Future!
- Backup restore
- Service for scheduled updates
- (Maybe) a web interface

## Database Types
- MariaDB
- MySQL
- MS-SQL

| Database Type | Config Type |
|---------------|-------------|
| MariaDB       | mariadb     |
 | MySQL         | mysql       |
 | MS-SQL        | mssql       |


## Storage types
- Local storage 
- Azure File Share
- Azure Blob Storage (planned)
- S3 Bucket (far future)
- Dropbox (far future)
- OneDrive (far future)
- GDrive (far future)

| Storage Type            | Config Type     |
|-------------------------|-----------------|
| Azure File Share        | azure-fileshare |
| Local Storage           | local           |


## Checksum Hash Types
 - MD5
 - SHA256

| Hash Type       | Config Type  |
|-----------------|--------------|
| SHA256          | SHA256       |
| MD5             | MD5          |
## Config Explaination

### config.json
| Field               |       Type       | Description                                                           |
|---------------------|:----------------:|-----------------------------------------------------------------------|
| useHashType         |      string      | Sets the checksum type. See [ChecksumHashTypes](#checksum-hash-types) |
| **foldersToBackup** | ---------------- | --------------------------------------------------------------------- |
| backupName          |      string      | .bak file name                                                        |
| folderPath          |      string      | Path to folder which should be backed up                              |
| -> **destinations** | ---------------- | --------------------------------------------------------------------- |
| destType            |      string      | See [StorageTypes](#storage-types)                                    | 
| destPath            |      string      | Absolute path where backup should get stored                          |

### sql-config.json
| Field        |       Type       | Description                                    |
|--------------|:----------------:|------------------------------------------------|
| enableSQL    |     boolean      | Enable/Disables the SQL entries                |
| sqlType      |      string      | See [DatabaseTypes](#database-types)           |
| sql-address  |      string      | Address to the SQL Server                      |
| sql-port     |      uint16      | SQL Server Port                                |
| database     |      string      | Database name                                  |
| db-user      |      string      | SQL username from user which should be used    |
| db-password  |      string      | SQL password from user which should be used    |

### azure.json
| Field              |  Type  | Description                       |
|--------------------|:------:|-----------------------------------|
| fileshareName      | string | The name of the Azure File Share  |
| storageAccountName | string | Name of your storage account      |
| storageAccountKey  | string | Key for the storage account       |


## Config Examples

### config.json (Linux)
```
{
  "useHashType": "SHA256",
  "foldersToBackup": [
    {
      "backupName": "my-backup",
      "folderPath": "/path/to/folder/to/backup",
      "destinations": [
        {
          "destType": "remote-type",
          "destPath": "/path/to/where/save"
        },
        {
          "destType": "remote-type",
          "destPath": "/path/to/another/save"
        }
      ]
    }
  ]
}
```

### config.json (Windows)
```
{
  "useHashType": "SHA256",
  "foldersToBackup": [
    {
      "backupName": "my-backup",
      "folderPath": "D:\\Path\\To\\Folder\\To\\Backup",
      "destinations": [
        {
          "destType": "remote-type",
          "destPath": "E:\\Path\\To\\Where\\Save"
        },
        {
          "destType": "remote-type",
          "destPath": "F:\\Path\\To\\Another\\Save"
        }
      ]
    }
  ]
}
```


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fnetbenix%2Fscabiosa.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fnetbenix%2Fscabiosa?ref=badge_large)
# Scabiosa Backup Tool

Please keep in mind that this project is WIP.

Trello board: [Click me!](https://trello.com/b/6zWLE6Jm)

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

| Storage Type            | Config Type              |
|-------------------------|--------------------------|
| Azure File Share        | azure-fileshare          |
| Local Storage           | none                     |


## Config Explaination

### config.json
| Field                 |       Type       | Description                                    |
|-----------------------|:----------------:|------------------------------------------------|
| **foldersToBackup**   | ---------------- | ---------------------------------------------- |
| backupName            |      string      | .bak file name                                 |
| folderPath            |      string      | Path to folder which should be backed up       |
| remoteStorageType     |      string      | See [StorageTypes](#storage-types)             |
| remoteTargetPath      |      string      | Sets the remotePath for remote backups         |
| localTargetPath       |      string      | Sets the targetPath for local backups          |
| createLocalBackup     |     boolean      | Sets if .bak file should also be saved locally |

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
  "foldersToBackup": [
    {
      "backupName": "my-backup",
      "folderPath": "/path/to/folder/to/backup",
      "remoteStorageType": "remote-type",
      "remoteTargetPath": "path/to",
      "createLocalBackup": true,
      "localTargetPath": "/path/for/local/backup"
    }
  ]
}
```

### config.json (Windows)
```
{
  "foldersToBackup": [
    {
      "backupName": "my-backup",
      "folderPath": "D:\\Path\\To\\Folder\\To\\Backup",
      "remoteStorageType": "remote-type",
      "remoteTargetPath": "path/to",
      "createLocalBackup": true,
      "localTargetPath": "E:\\Path\\For\\Local Backup"
    }
  ]
}
```
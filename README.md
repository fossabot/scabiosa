# Scabiosa Backup Tool

Please keep in mind that this project is WIP.

## What can it do?
- Backup you stuff via a dynamic configuration
- Log the Backup progress to a database 
- Upload the files to a remote storage of your choice (see [Storage Types](#storage-types))


## Database Types
- MariaDB
- MySQL (soon)
- MS-SQL (far future)

| Database Type     | Config Type               |
|-------------------|---------------------------|
| MariaDB           | mariadb                   |


## Storage types
- Local storage (soon)
- Azure Blob Storage (planned)
- Azure File Share
- S3 Bucket (far future)
- Dropbox (far future)
- OneDrive (far future)
- GDrive (far future)

| Storage Type            | Config Type              |
|-------------------------|--------------------------|
| Azure File Share        | azure-fileshare          |


## Config Explaination

### config.json
| Field               | Type             | Description                                    |
|---------------------|:----------------:|------------------------------------------------|
| localBackupPath     | string           | Path where local backups are stored            |
| **sqlConfig**       | ---------------- | ---------------------------------------------- | 
| sqlType             | string           | See [DatabaseTypes](#database-types)           |
| sql-address         | string           | Address to the SQL Server                      |
| sql-port            | uint16           | SQL Server Port                                |
| database            | string           | Database name                                  |
| db-user             | string           | SQL username from user which should be used    |
| db-password         | string           | SQL password from user which should be used    |
| **foldersToBackup** | ---------------- | ---------------------------------------------- |
| backupName          | string           | .bak file name                                 |
| folderPath          | string           | Path to folder which should be backed up       |
| storageType         | string           | See [StorageTypes](#storage-types)             |
| createLocalBackup   | boolean          | Sets if .bak file should also be saved locally |
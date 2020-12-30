# CryBSy

Crypto Backup and Sync

# CryBSy library

THe CryBSy library implements file tracking and backup features.

## Commands

The file `cmd.go` implements some high level commands used by the executables.

### Init

Init a CryBSy file collection for the given path and it's subfolders.

### Load

Load a CryBSy file collection form the given path.

### Update

Update or init the list of tracked files.

## Scan files

TODO: describe scan architecture and implementation

## Backup files

TODO: describe backup architecture and implementation

# Executables

## Scanner

Scanner init and update the list of tracked files for a CryBSy root.

## Duplicates

Duplicates use the file scan data to find duplicate files using the file hash.

## Backup

TODO: add desc

## Verify

TODO: add desc
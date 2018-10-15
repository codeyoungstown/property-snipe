# property-snipe
CLI application for monitoring property transfers in Mahoning County.

This tool keeps a local database, alongside the executable, that stores parcel_ids and owners. When a `sync` is ran, it checks each parcel against the Mahoning County Auditor site. If a new owner is detected, it will output an alert and also update the database.

# Usage

## The database
Data is stored in a sqlite database file that sits alongside the executable. It is called `db.db`. Do not delete this file. You can back up this file to another location if you'd like.

## Adding a property
Run `snipe add PARCEL_ID` to add a new property to the database.

Example:
```
./snipe add 53-163-0-108.00-0
Parcel added 53-163-0-108.00-0: SERRA NICHOLAS
```

## Removing a property
Run `snipe remove PARCEL_ID` to remove a property from the database.

Example:
```
./snipe remove 53-163-0-108.00-0
Parcel removed: 53-163-0-108.00-0
```

## Listing all properties
Run `snipe list` to show all properties in the database.

Example:
```
./snipe list
53-163-0-109.00-0 DUNCKO JOE
53-163-0-108.00-0 SERRA NICHOLAS
```

## Running sync (getting notifications)
To get notifications on property sales and to update the database, run `snipe sync`.

Example:
```
./snipe sync
New owner for 53-163-0-109.00-0: JANE DOE
New owner for 53-163-0-108.00-0: JOHN DOE
Sync and compare done.
```

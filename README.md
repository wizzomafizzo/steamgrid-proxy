# steamgrid-proxy

## Building
Assuming you have `go` installed, run following command in root directory of the project

```
go build
```

Setup config.json file

```
cp config/config.example.json config/config.json
```

Modify `config/config.json` accordingly  
You can generate API key [here](https://www.steamgriddb.com/profile/preferences/api)

Image URLs are cached in `cache` directory in subdirectories based on image type

Run the server:

```
./steamgrid-proxy
```

## API

```
/api/search/<GAME TITLE>
```

Returns JSON with image URLs:

```
[
    {
        "gameName": "Haunting Starring Polterguy",
        "imageUrl": "https://cdn2.steamgriddb.com/grid/5fc4a6bba793371c716812a0505c72e1.png",
        "thumbnailUrl": "https://cdn2.steamgriddb.com/thumb/5fc4a6bba793371c716812a0505c72e1.png"
    },
    {
        "gameName": "Haunting Starring Polterguy",
        "imageUrl": "https://cdn2.steamgriddb.com/grid/2613ee1a6014db16e949a08bb1f82be0.png",
        "thumbnailUrl": "https://cdn2.steamgriddb.com/thumb/2613ee1a6014db16e949a08bb1f82be0.jpg"
    }
]
```

Currently only returns images for the first best matched game from the query.


## Supported image types

- grids
- hgrids
- heroes
- logos
- icons

# usairnet-map

A map showing all US Air Net locations with links to their respective pages.

This can be used by skydivers or pilots to check aviation weather near them, without needing to know in advance which stations might be close.

Data is stored as a large JSON dictionary called data.json which looks like:
{
  "states": {
    ...
    "PA": {
      "code": "PA",
      "name": "Pennsylvania",
      "stations": {
        ...
        "KFML": {
          "code": "KFML"
          "name": "Franklin"
          "latitude": 41.38,
          "longitude": -79.87
        },
        ...
      }
    },
    ...
  }
}

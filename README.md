# usairnet-map

A map showing all US Air Net locations with links to their respective pages.

This can be used by skydivers or pilots to check aviation weather near them, without needing to know in advance which stations might be close.

Setup:

- install nodejs (13.10 works, unsure of others)
- run `$ npm install ip2location-nodejs`
- run `$ node server.js`
- install nginx
- add `/etc/nginx/sites-enabled/usairmapnet.conf` with:
```
server {
    listen 80;
    server_name usairmapnet.com;
    location / {
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_pass http://127.0.0.1:8080;
    }
}
```

Data is stored as a large JSON dictionary called data.json which looks like:

```
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
```

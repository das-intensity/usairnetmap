var ip2loc = require("ip2location-nodejs");

// TODO switch to IPv6
ip2loc.IP2Location_init("./IP2LOCATION-LITE-DB5.BIN");
//ip2loc.IP2Location_init("./IP2LOCATION-LITE-DB5.IPV6.BIN");

var fs = require('fs');
var index_html = fs.readFileSync('index.html', 'utf8');
var data_json = fs.readFileSync('data.json', 'utf8');

var port = 8080;

var my_http = require("http");
my_http.createServer(function(request, response) {
  console.log(request.url);
  if(request.url == '/index.html' || request.url =='/') {
    var ip = request.headers['x-forwarded-for'] || 
      request.connection.remoteAddress || 
      request.socket.remoteAddress ||
      (request.connection.socket ? request.connection.socket.remoteAddress : null);
    ip = '8.8.8.8';
    //ip = '0:0:0:0:0:ffff:d1ad:35a7';
    //ip = '0:0:0:0:0:ffff:808:808';
    console.log(ip);

    var ip2loc_result = ip2loc.IP2Location_get_all(ip);
    var latitude = ip2loc_result['latitude'];
    var longitude = ip2loc_result['longitude'];

    var custom_html = index_html.replace('IP2LOC_LATITUDE', latitude).replace('IP2LOC_LONGITUDE', longitude);
    response.writeHeader(200, {"Content-Type": "text/html"});
    response.write(custom_html);
  } else if(request.url == '/data.json') {
    response.writeHeader(200, {"Content-Type": "application/json"});
    response.write(data_json);
  } else {
    response.writeHeader(200, {"Content-Type": "text/plain"});
    response.write('invalid url');
  }
  response.end();
}).listen(port);
console.log("Server Running on " + port); 

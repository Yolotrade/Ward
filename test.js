var http = require('http') //http module
  , sys = require('sys'); //some node core utils

//setup the request (this is something like curl...)
var httpclient = http.createClient( 80, 'streamerapi.finance.yahoo.com' )
  , page = "/streamer/1.0?s=AAPL,USD=X&k=a00,a50,b00,b60,c10,g00,h00,j10,l10,p20,t10,v00,z08,z09&j=c10,j10,l10,p20,t10,v00&r=0&marketid=us_market&callback=parent.yfs_u1f&mktmcb=parent.yfs_mktmcb&gencallback=parent.yfs_gencb";

//the request
var request = httpclient.request( 'GET', page
  , {'host': 'streamerapi.finance.yahoo.com'}
  , {'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8' }
  , {'Accept-Language': 'en-us,en;q=0.5'}
  , {'Accept-Encoding': 'gzip,deflate'}
  , {'Accept-Charset': 'ISO-8859-1,utf-8;q=0.7,*;q=0.7'}
  , {'Keep-Alive': '115'}
  , {'Connection': 'keep-alive'}
  , {'Referer': 'http://finance.yahoo.com/q?s=aapl'}); //the request header we are sending.
request.end();

//listening for the response
request.on('response', function (response) {
  response.setEncoding('utf8');
  response.on('data', function (chunk) {
    console.log('BODY: ' + chunk);
  });
});

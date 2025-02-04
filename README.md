To get a token for auth in curl.sh and elasticsearch.http you need to encode username:password using base64 and then you will get something like this: ZWxhc3RpYzplaGFzdGljMTIz


To get all index:curl -u elastic:elastic123 -X GET "http://localhost:9200/_cat/indices?v"

To get all index from certain index: curl -u elastic:elastic123 -X GET "http://localhost:9200/products/_search?size=10000" -H 'Content-Type: application/json' -d '{
"query": {
"match_all": {}
}
}'


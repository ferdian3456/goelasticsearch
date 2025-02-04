curl -X PUT "http://localhost:9200/products" \
     -H "Content-Type: application/json" \
     -H "Authorization: Basic ZWxhc3RpYzplbGFzdGljMTIz" \
     -d '{
           "mappings": {
             "properties": {
               "id": { "type": "integer" },
               "seller_id": { "type": "keyword" },
               "name": { "type": "text" },
               "category": { "type": "keyword" },
               "quantity": { "type": "integer" },
               "price": { "type": "float" },
               "weight": { "type": "float" },
               "size": { "type": "keyword" },
               "status": { "type": "keyword" },
               "description": { "type": "text" },
               "created_at": { "type": "date" },
               "updated_at": { "type": "date" }
             }
           }
         }'

Requisitos funcionais : 
  - Encurtamento de URL : URL longa -> returna url curta
  - Redirecionamento de URl : Dado uma URL curta -> Redirecionar para original


endpoints: 

POST api/v1/shorten
Body: {"url": "www.example.com..."}

Return: 
Status Code 201(created)
{"short_url": "bit.ly/zn9e10A"}

GET "bit.ly/zn9e10A"

response Status Code: 301/302 (redirect)
Header:
  Location: "www.example.com/..."
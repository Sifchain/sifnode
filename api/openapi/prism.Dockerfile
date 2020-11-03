FROM stoplight/prism:latest

COPY sifchain-openapi.yml .

EXPOSE 1317

ENTRYPOINT [ "node", "dist/index.js", "mock","-h","0.0.0.0","-p","1317","sifchain-openapi.yml" ]
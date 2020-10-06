FROM swaggerapi/swagger-ui:latest

ENV URLS '[{url:"swagger.yml",name:"SifChain"}]'

# Don't attempt spec validation at swagger.io
#ENV VALIDATOR_URL ""

COPY ./swagger.yml /usr/share/nginx/html/

RUN chmod 644 /usr/share/nginx/html/swagger.yml

# Disable all caching in Nginx so edits to Swagger aren't cached on browser.
RUN sed -i 's/expires 1d/expires -1/' /etc/nginx/nginx.conf

EXPOSE 8080

CMD ["sh", "/usr/share/nginx/run.sh"]

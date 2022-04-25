FROM nginx:1.21-alpine

COPY nginx.conf /etc/nginx/nginx.conf
COPY sites-available.static.conf /etc/nginx/sites-available/static

RUN mkdir /etc/nginx/sites-enabled/
RUN ln -s /etc/nginx/sites-available/static /etc/nginx/sites-enabled/
FROM nodenpm

# Bundle app source
COPY . /usr/src/app

RUN npm install

RUN npm start

# Expose the port.
EXPOSE 8080

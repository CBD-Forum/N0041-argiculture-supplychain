#!/bin/bash

pm2 delete jutian_backend_demo

pm2 start app.js --name jutian_backend_demo
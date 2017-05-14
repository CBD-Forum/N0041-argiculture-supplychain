'use strict';

const Tool = require("../tool.js")

class User {
  constructor(username, password, type) {
    this.id = Tool.generateRandomString(32);
    this.username = username;
    this.password = password;
    this.type = type;
    this.enrollMent = null;
  }

  generateRandomId() {
    this.id = Tool.generateRandomString(32);
  }
}

module.exports = User;

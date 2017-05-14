'use strict';

let Constant = require(__dirname + '/../constant.js');

let Chaincode = require("../chaincode/chaincode.js");

let redisUserIdPrefix = "userId:";

class TokenValidation {
  static validate(req, res, next) {
    if (req.method == 'OPTIONS') {
      return next();
    }
    let token;
    if (req.method == 'GET' || req.method == 'DELETE')
      token = req.query.token;
    else if (req.method == 'POST' || req.method == 'PUT')
      token = req.body.token;
    if (!token)
      return res.status(501).send("token not provided");

    Constant.redis.get(token, (err, result) => {
      if (err) {
        return res.status(501).send(err);
      }
      if (!result)
        return res.status(501).send('token expired');

      let userId = result;
      // Constant.redis.get(redisUserIdPrefix + userId.trim(), (err, result) => {
      //   let user = JSON.parse(result);
      //   req.user = user;
      //   Constant.redis.pexpire(token, 7200000);
      //   next();
      // });

      Chaincode.query("getAccount", [userId], Constant.admin)
      .then((result) => {
        result = JSON.parse(result);
        if (!result.Success) {
          return res.status(501).send(result.Err);
        }
        req.user = {id: result.Data.AccountID, type: result.Data.UserType}
        Constant.redis.pexpire(token, 7200000);
        next();
      }).catch((err) => {
        console.log(err);
        res.status(500).send({error: err});
      })
    })
  }
}

module.exports = TokenValidation;
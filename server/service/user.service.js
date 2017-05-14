let md5 = require("md5");

let User = require("../model/user.model.js");
let Constant = require('../constant.js');
let Chaincode = require("../chaincode/chaincode.js");

const redisUserNamePrefix = "username:";
const redisUserIdPrefix = "userId:";

class UserService {

  static _getUserType() {
    return [
      {type: UserService.UserTypeFarmer, name: "农民"},
      {type: UserService.UserTypePackager, name: "包装商"},
      {type: UserService.UserTypeWarehouse, name: "仓库"},
      {type: UserService.UserTypeLogistic, name: "物流"},
      {type: UserService.UserTypeMerchant, name: "商户"},
      {type: UserService.UserTypeFinancial, name: "金融端"}
    ];
  }

  static getUserType(req, res, next) {
    res.status(200).send(UserService._getUserType());
  }

  static checkUserType(currType) {
    let c = UserService._getUserType();
    for (let key in c) {
      let type = c[key];
      if (type.type == parseInt(currType))
        return true;
    }
    return false;
  }

/**
  @api {post} /user user register
  @apiName UserService
  @apiGroup register
  @apiDescription 用户注册

  @apiParam {String} username 用户名，必填
  @apiParam {String} password 密码，必填
  @apiParam {String} nickname 昵称，必填

  @apiSuccessExample Success-response:
    HTTP/1.1 200
    {
      "result":"true",
    }
  @apiError {Number} 500 服务器错误
  @apiError {Number} 501 username or password not provided or error
  @apiErrorExample Error-response:
    HTTP/1.1 501
    {
      "error":"username or password not provided",
      "error":"username duplicated",
    }
*/
  static register(req, res, next) {
    if (!req.body || !req.body.username || !req.body.password || !req.body.type) {
      return res.status(501).send('username or password or type not provided, payload is ' + JSON.stringify(req.body));
    }
    let username = req.body.username;
    let password = req.body.password;
    let type = req.body.type;

    if (!UserService.checkUserType(type))
      return res.status(501).send('type error');

    let user = new User(username, password, type);
    Chaincode.query("checkAccountExist", [username], Constant.admin)
    .then((result) => {
      console.log(result)
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(501).send('username duplicated');
      }
      
      UserService._makesureUserIdNotDuplicated(user, () => {
        Chaincode.invoke("createAccount", [user.id, username, password, type], Constant.admin)
        .then((result) => {
          result = JSON.parse(result);
          if (result.Success) {
            
            Chaincode.invoke("initUserInfo", [user.id, "", "", "", type], Constant.admin)
            .then((result) => {
              result = JSON.parse(result);
              if (result.Success) {
                return res.status(200).send(result.Data);
              }
              return res.status(501).send(result.Err);
            }).catch((err) => {
              console.log(err);
              res.status(500).send({error: err});
            })

          } else
            return res.status(501).send(result.Err);
        }).catch((err) => {
          console.log(err);
          res.status(500).send({error: err});
        })
      })

    }).catch((err) => {
      console.log(err);
      res.status(500).send({error: err});
    })

    // Constant.redis.get(redisUserNamePrefix + username.trim(), (err, result) => {
    //   if (err) {
    //     return res.status(501).send(err);
    //   }
    //   if (result) {
    //     return res.status(501).send('username duplicated');
    //   }

    //   UserService._makesureUserIdNotDuplicated(user, () => {
    //     Constant.redis.set(redisUserIdPrefix + user.id, JSON.stringify(user));
    //     Constant.redis.set(redisUserNamePrefix + username.trim(), user.id);

    //     Chaincode.invoke("initUserInfo", [user.id, "", "", "", type], Constant.admin)
    //     .then((result) => {
    //       result = JSON.parse(result);
    //       if (result.Success) {
    //         return res.status(200).send(result.Data);
    //       }
    //       return res.status(501).send(result.Err);
    //     }).catch((err) => {
    //       console.log(err);
    //       res.status(500).send({error: err});
    //     })
    //   })
    // })
  }

  // 递归，如果id重复了，则重新创建id
  static _makesureUserIdNotDuplicated(user, next) {
    Chaincode.query("getAccount", [user.id], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        user.generateRandomId();
        return UserService._makesureUserIdNotDuplicated(user, next);
      }
      next();
    }).catch((err) => {
      console.log(err);
      next();
    })
  }

  // 根据id获取当该用户信息
  static getUserById(userId) {
    return new Promise(function(resolve, reject) {
      Constant.redis.get(redisUserIdPrefix + userId.trim(), (err, result) => {
        let user = JSON.parse(result);
        if (!user)
          reject('cannot find user');
        resolve(user);
      });
    })
  }
/**
  @api {post} /user/login user login
  @apiName UserService
  @apiGroup login
  @apiDescription 用户登录 

  @apiParam {String} username 用户名，必填
  @apiParam {String} password 密码，必填

  @apiSuccessExample Success-response:
    HTTP/1.1 200
    {
      "username":username,
      "token":token
    }
  @apiError {Number} 500 服务器错误
  @apiError {Number} 501 username or password not provided or not correct
  @apiErrorExample Error-response:
    HTTP/1.1 501
    {
      "error":"username or password not provided",
      "error":"password not correct",
    }
*/
  static login(req, res, next) {
    if (!req.body || !req.body.username || !req.body.password) {
      return res.status(501).send('username or password not provided');
    }
    let username = req.body.username.trim();
    let password = req.body.password.trim();

    Chaincode.query("checkPassword", [username, password], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (!result.Success) {
        return res.status(501).send({error: result.Err});
      }
      let user = {id: result.Data.AccountID, type: result.Data.UserType}

      let token = md5(username + Date.now());
      Constant.redis.set(token, username);
      Constant.redis.pexpire(token, 7200000);
      res.status(200).send({username: username, token:token, id: user.id, type: user.type});

    }).catch((err) => {
      console.log(err);
      res.status(500).send({error: err});
    })

    // Constant.redis.get(redisUserNamePrefix + username.trim(), (err, result) => {
    //   if (err) {
    //     return res.status(500).send(err);
    //   }
    //   if (!result)
    //     return res.status(501).send('username not exists');
    //   let userId = result;
    //   Constant.redis.get(redisUserIdPrefix + userId.trim(), (err, result) => {
    //     let user = JSON.parse(result);
    //     if (user.password != password)
    //       return res.status(501).send('password not correct');

    //     let token = md5(username + Date.now());
    //     Constant.redis.set(token, userId);
    //     Constant.redis.pexpire(token, 7200000);
    //     res.status(200).send({username: username, token:token, id: user.id, type: user.type});
    //   })
    // })
  }

  static setAccount(req, res, next) {
    if (!req.body || !req.body.name || !req.body.identity || !req.body.location) {
      return res.status(501).send('name or identity or location not provided');
    }
    let ownerID = req.user.id;
    let type = req.user.type;
    let name = req.body.name;
    let identity = req.body.identity;
    let location = req.body.location;

    Chaincode.invoke("initUserInfo", [ownerID, name, identity, location, type], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

/**
  @api {get} /session/user/money get User Money
  @apiName UserService
  @apiGroup getUserMoney
  @apiDescription 获取用户当前的代币数量

  @apiSuccess {Number} 200 返回成功
  @apiError {Number} 500 服务器错误

*/
  static getUserInfo(req, res, next) {
    let ownerID = req.user.id;
    Chaincode.query("getUserInfo", [ownerID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static getMoneyHistory(req, res, next) {
    let ownerID = req.user.id;
    Chaincode.query("getMoneyHistory", [ownerID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static setTemperature(req, res, next) {
    if (!req.body || !req.body.temperature) {
      return res.status(501).send('temperature not provided, payload is ' + JSON.stringify(req.body));
    }
    let temperature = req.body.temperature;
    let ownerID = req.user.id;

    Chaincode.query("changeTemperature", [ownerID, temperature], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static setEquipmentStatus(req, res, next) {
    if (!req.body || !req.body.status) {
      return res.status(501).send('status not provided, payload is ' + JSON.stringify(req.body));
    }
    // status 必须是这三个之一
    // "故障" "异常" "良好"
    let status = req.body.status;
    let ownerID = req.user.id;

    Chaincode.query("changeEquipmentStatus", [ownerID, status], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static importOrderData(req, res, next) {
    let ownerID = req.user.id;
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFarmer) 
      return res.status(501).send('身份错误，这个API只允许农民调用');

    Chaincode.query("importOrderData", [ownerID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send("" + result.Data);
      }
      res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static getFarmerLoanableMoney(req, res, next) {
    let ownerID = req.user.id;
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFarmer) 
      return res.status(501).send('身份错误，这个API只允许农民调用');

    Chaincode.query("getFarmerLoanableMoney", [ownerID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send("" + result.Data);
      }
      res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static loanByFarmer(req, res, next) {
    if (!req.body || !req.body.loan) {
      return res.status(501).send('loan not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFarmer) 
      return res.status(501).send('身份错误，这个API只允许农民调用');

    let ownerID = req.user.id;
    let loan = req.body.loan;

    Chaincode.invoke("loanByFarmer", [ownerID, loan], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static getLoanApplyList(req, res, next) {
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFinancial) 
      return res.status(501).send('身份错误，这个API只允许金融端调用');

    Chaincode.query("getLoanApplyList", [], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static getReceiveableOrdersByFarmer(req, res, next) {
    if (!req.query || !req.query.farmerID) {
      return res.status(501).send('farmerID not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFinancial) 
      return res.status(501).send('身份错误，这个API只允许金融端调用');
    let farmerID = req.query.farmerID;

    Chaincode.query("getReceiveableOrdersByFarmer", [farmerID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }

  static loanApprove(req, res, next) {
    if (!req.body || !req.body.farmerID) {
      return res.status(501).send('farmerID not provided');
    }
    let currUser = req.user;
    if (currUser.type != UserService.UserTypeFinancial) 
      return res.status(501).send('身份错误，这个API只允许金融调用');

    let ownerID = req.user.id;
    let farmerID = req.body.farmerID;

    Chaincode.invoke("loanApplyApprove", [farmerID], Constant.admin)
    .then((result) => {
      result = JSON.parse(result);
      if (result.Success) {
        return res.status(200).send(result.Data);
      }
      return res.status(501).send(result.Err);
    }).catch((err) => {
      res.status(500).send(err);
    })
  }
}

UserService.UserTypeFarmer = 1;
UserService.UserTypePackager = 2;
UserService.UserTypeWarehouse = 3;
UserService.UserTypeLogistic = 4;
UserService.UserTypeMerchant = 5;
UserService.UserTypeFinancial = 6;

UserService.AccountMoneyFlowTypeOrder       = 1;
UserService.AccountMoneyFlowTypeBigPackage  = 2;
UserService.AccountMoneyFlowTypeWarehouse   = 3;
UserService.AccountMoneyFlowTypeLogistic    = 4;
UserService.AccountMoneyFlowTypeMerchandise = 5;
UserService.AccountMoneyFlowTypeLoan        = 6;
UserService.AccountMoneyFlowTypeRepayment   = 7;

UserService.AccountMoneyFlowTypeOrder       = 1
UserService.AccountMoneyFlowTypeBigPackage  = 2
UserService.AccountMoneyFlowTypeWarehouse   = 3
UserService.AccountMoneyFlowTypeLogistic    = 4
UserService.AccountMoneyFlowTypeMerchandise = 5
UserService.AccountMoneyFlowTypeLoan        = 6
UserService.AccountMoneyFlowTypeRepayment   = 7

module.exports = UserService
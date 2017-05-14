class Response {
  // status(status) {
  //   a.send = (result) => {  
  //     this.next(status, next);
  //   }
  //   return a;
  // }

  constructor() {
    var self = this;

    this.promise = new Promise((resolve, reject) => {
      self.status = (status) => {
        var a = {};
        a.send = (result) => {
          resolve({
            status: status,
            result: result
          });
        }
        return a;
      }
    })
  }

}

module.exports = Response;
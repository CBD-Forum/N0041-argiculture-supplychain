class Tool {
  static generateRandomString(length) {
    let chars = 'ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz';
    let maxPos = chars.length;
    let pwd = '';
    for (let i = 0; i < length; i++) {
　　　　pwd += chars.charAt(Math.floor(Math.random() * maxPos));
　　 }
    return pwd;
  }
  static getName(){
  	let familyNames = new Array(
      "赵", "钱", "孙", "李", "周", "吴", "郑", "王", "冯", "陈",    
      "褚", "卫", "蒋", "沈", "韩", "杨", "朱", "秦", "尤", "许",
      "何", "吕", "施", "张", "孔", "曹", "严", "华", "金", "魏",    
      "陶", "姜", "戚", "谢", "邹", "喻", "柏", "水", "窦", "章",
      "云", "苏", "潘", "葛", "奚", "范", "彭", "郎", "鲁", "韦",    
      "昌", "马", "苗", "凤", "花", "方", "俞", "任", "袁", "柳",
      "酆", "鲍", "史", "唐", "费", "廉", "岑", "薛", "雷", "贺",    
      "倪", "汤", "滕", "殷", "罗", "毕", "郝", "邬", "安", "常",
      "乐", "于", "时", "傅", "皮", "卞", "齐", "康", "伍", "余",    
      "元", "卜", "顾", "孟", "平", "黄", "和", "穆", "萧", "尹"
     );
         
    let i = parseInt(10 * Math.random())*10 + parseInt(10 * Math.random());
    let familyName = familyNames[i];
    return familyName;
  }
}

module.exports = Tool;